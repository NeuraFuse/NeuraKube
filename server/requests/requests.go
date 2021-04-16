package requests

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"../../../tools-go/config"
	devConfig "../../../tools-go/config/dev"
	infraConfig "../../../tools-go/config/infrastructure"
	projectConfig "../../../tools-go/config/project"
	userConfig "../../../tools-go/config/user"
	"../../../tools-go/errors"
	"../../../tools-go/filesystem"
	"../../../tools-go/logging"
	"../../../tools-go/objects/strings"
	"../../../tools-go/projects"
	"../../../tools-go/readers/yaml"
	"../../../tools-go/runtime"
	"../../../tools-go/users"
	"../../../tools-go/vars"
	setup "../../infrastructure"
	"../../infrastructure/ci/app"
	"../../infrastructure/ci/develop/remote"
	"../../infrastructure/ci/inference"
	inferenceRequest "../../infrastructure/ci/inference/request"
	gconfig "../../infrastructure/providers/gcloud/config"
	"../auth"
	"../routines"
)

type F struct{}

func (f F) Router(route string) http.Handler {
	var handler http.Handler
	if route == "/inspect" {
		handler = auth.F.NoAuth(auth.F{}, f.getInspect)
	} else if route == "/inspect/healthcheck" {
		handler = auth.F.NoAuth(auth.F{}, f.getHealthCheck)
	} else if route == "/inspect/version" {
		handler = auth.F.NoAuth(auth.F{}, f.getVersion)
	} else if route == "/inspect/infrastructure/init" {
		handler = auth.F.NoAuth(auth.F{}, f.getInfrastructureInit)
	} else if route == "/users" {
		handler = auth.F.Check(auth.F{}, f.getUsers)
	} else if route == "/user/create" {
		handler = auth.F.Check(auth.F{}, f.createUser)
	} else if route == "/user/infrastructure" {
		handler = auth.F.Check(auth.F{}, f.infrastructure)
	} else if route == "/user/infrastructure/setup" {
		handler = auth.F.Check(auth.F{}, f.infraSetup)
	} else if route == "/user/project" {
		handler = auth.F.Check(auth.F{}, f.project)
	} else if route == "/user/infrastructure/auth/kubeconfig/create" {
		handler = auth.F.Check(auth.F{}, f.createInfraProviderAuthKubeConfig)
	} else if route == "/user/infrastructure/auth/gcloud/create" {
		handler = auth.F.Check(auth.F{}, f.createInfraProviderAuthGcloud)
	} else if route == "/user/devConfig/create" {
		handler = auth.F.Check(auth.F{}, f.createUserdevConfig)
	} else if route == "/develop/remote" {
		handler = auth.F.Check(auth.F{}, f.developRemote)
	} else if route == "/develop/remote/delete" {
		handler = auth.F.Check(auth.F{}, f.deleteDevelopRemote)
	} else if route == "/app" {
		handler = auth.F.Check(auth.F{}, f.app)
	} else if route == "/app/delete" {
		handler = auth.F.Check(auth.F{}, f.deleteApp)
	} else if route == "/app/inference" {
		handler = auth.F.Check(auth.F{}, f.appInference)
	} else if route == "/app/inference/delete" {
		handler = auth.F.Check(auth.F{}, f.appInferenceDelete)
	} else if route == "/app/gpt/inference" {
		handler = auth.F.Check(auth.F{}, f.appGPTInference)
	}
	return handler
}

func (f F) appGPTInference(wri http.ResponseWriter, req *http.Request) {
	body := f.readRequestBody(req)
	context := strings.BytesToString(body)
	fmt.Fprintf(wri, inferenceRequest.F.Router(inferenceRequest.F{}, context))
}

func (f F) getInspect(wri http.ResponseWriter, req *http.Request) {
	hc := "ok"
	version := vars.NeuraKubeVersion
	init := "uninitialized"
	if f.GetInfrastructureInitStatus() {
		init = "initialized"
	}
	fmt.Fprintf(wri, strings.Join([]string{version, hc, init}, ","))
}

func (f F) getInfrastructureInit(wri http.ResponseWriter, req *http.Request) {
	auth.F.SetPaths(auth.F{})
	/*if config.Setting("get", "dev", "Spec.Status", "") == "active" {
		ci.RemovePackage("users")
	}*/
	if f.GetInfrastructureInitStatus() {
		fmt.Fprintf(wri, "initialized")
	} else {
		fmt.Fprintf(wri, "uninitialized")
	}
}

func (f F) GetInfrastructureInitStatus() bool {
	return users.Existing()
}

func (f F) getHealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func (f F) getVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, vars.NeuraKubeVersion)
}

func (f F) developRemote(wri http.ResponseWriter, req *http.Request) {
	action := f.getURLQueryParam(req, "action")
	if action == "prepare" {
		fmt.Fprintf(wri, remote.F.Prepare(remote.F{}))
	} else if action == "create" {
		fmt.Fprintf(wri, remote.F.Create(remote.F{}))
	} else {
		fmt.Fprintf(wri, "Error: Unsupported action: "+action)
	}
}

func (f F) deleteDevelopRemote(wri http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(wri, remote.F.Delete(remote.F{}))
}

func (f F) app(wri http.ResponseWriter, req *http.Request) {
	action := f.getURLQueryParam(req, "action")
	if action == "prepare" {
		fmt.Fprintf(wri, app.F.Prepare(app.F{}))
	} else if action == "create" {
		fmt.Fprintf(wri, app.F.Create(app.F{}))
	} else {
		fmt.Fprintf(wri, "Error: Unsupported action: "+action)
	}
}

func (f F) deleteApp(wri http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(wri, app.F.Delete(app.F{}))
}

func (f F) appInference(wri http.ResponseWriter, req *http.Request) {
	action := f.getURLQueryParam(req, "action")
	if action == "prepare" {
		fmt.Fprintf(wri, inference.F.Prepare(inference.F{}))
	} else if action == "create" {
		fmt.Fprintf(wri, inference.F.Create(inference.F{}))
	} else {
		fmt.Fprintf(wri, "Error: Unsupported action: "+action)
	}
}

func (f F) appInferenceDelete(wri http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(wri, inference.F.Delete(inference.F{}))
}

func (f F) createUser(wri http.ResponseWriter, req *http.Request) {
	createUser := false
	userType := ""
	if !f.GetInfrastructureInitStatus() {
		createUser = true
		userType = "Admin user "
	} else {
		if config.Setting("get", "server", "Spec.Users.Admin", "") == users.GetIDActive() {
			if !users.Exists(users.GetIDActive()) {
				createUser = true
				userType = "User "
			} else {
				fmt.Fprintf(wri, "Unable to create user "+users.GetIDActive()+" because it already exists.")
			}
		} else {
			fmt.Fprintf(wri, "The user creation failed because this action is only available to admin users.")
		}
	}
	if createUser {
		users.Create(users.GetIDActive())
		f.saveToFile(req, userConfig.F.GetFilePath(userConfig.F{}))
		config.Setting("init", "user", "Spec.", "")
		auth.F.Config(auth.F{}, req)
		logging.Log([]string{"", vars.EmojiUser, vars.EmojiSuccess}, userType+users.GetIDActive()+" created.", 0)
		config.Setting("set", "server", "Spec.Users.Admin", users.GetIDActive())
		fmt.Fprintf(wri, "success")
	}
}

func (f F) infrastructure(wri http.ResponseWriter, req *http.Request) {
	auth.F.SetPaths(auth.F{})
	update := f.saveToFile(req, infraConfig.F.GetFilePath(infraConfig.F{}))
	config.Setting("init", "infrastructure", "", "")
	success := false
	if !update {
		if !update {
			logging.Log([]string{"", vars.EmojiProject, vars.EmojiSuccess}, "Infrastructure created.", 0)
		}
		logging.Log([]string{"", vars.EmojiInfra, vars.EmojiSuccess}, "Infrastructure config synced.", 0)
		success = true
	} else {
		logging.Log([]string{"", vars.EmojiInfra, vars.EmojiSuccess}, "Infrastructure config updated.", 0)
		success = true
	}
	if success {
		fmt.Fprintf(wri, "success")
	}
}

func (f F) project(wri http.ResponseWriter, req *http.Request) {
	auth.F.SetPaths(auth.F{})
	update := f.saveToFile(req, projectConfig.F.GetFilePath(projectConfig.F{}))
	config.Setting("init", "dev", "Spec.", "")
	success := false
	if !update {
		projects.F.Create(projects.F{}, projects.IDActive)
		if !update {
			logging.Log([]string{"", vars.EmojiProject, vars.EmojiSuccess}, "Project "+projects.IDActive+" created.", 0)
		}
		logging.Log([]string{"", vars.EmojiInfra, vars.EmojiSuccess}, "Project config synced.", 0)
		success = true
	} else {
		logging.Log([]string{"", vars.EmojiInfra, vars.EmojiSuccess}, "Project config updated.", 0)
		success = true
	}
	if success {
		fmt.Fprintf(wri, "success")
	}
}

func (f F) createInfraProviderAuthKubeConfig(wri http.ResponseWriter, req *http.Request) {
	filePath := config.Setting("get", "infrastructure", "Spec.Cluster.Auth.KubeConfigPath", "")
	if users.Exists(users.GetIDActive()) {
		f.saveToFile(req, filePath)
		infraConfig.F.SetKubeConfig(infraConfig.F{})
		logging.Log([]string{"", vars.EmojiInfra, vars.EmojiSuccess}, "Infrastructure auth kubeconfig synced.", 0)
		routines.F.Router(routines.F{})
		fmt.Fprintf(wri, "success")
	} else {
		fmt.Fprintf(wri, "Unable to register infrastructure auth kubeconfig: The user does not exist!")
	}
}

func (f F) createInfraProviderAuthGcloud(wri http.ResponseWriter, req *http.Request) {
	filePath := config.Setting("get", "infrastructure", "Spec.Gcloud.Auth.ServiceAccountJSONPath", "")
	if users.Exists(users.GetIDActive()) {
		f.saveToFile(req, filePath)
		gconfig.F.SetConfigs(gconfig.F{})
		logging.Log([]string{"", vars.EmojiInfra, vars.EmojiSuccess}, "Infrastructure auth "+vars.InfraProviderGcloud+" synced.", 0)
		routines.F.Router(routines.F{})
		fmt.Fprintf(wri, "success")
	} else {
		fmt.Fprintf(wri, "Unable to register infrastructure auth gcloud: The user does not exist!")
	}
}

func (f F) createUserdevConfig(wri http.ResponseWriter, req *http.Request) {
	if users.Exists(users.GetIDActive()) {
		f.saveToFile(req, devConfig.F.GetFilePath(devConfig.F{}))
		logging.Log([]string{"", vars.EmojiDev, vars.EmojiSuccess}, "devConfig synced.", 0)
		fmt.Fprintf(wri, "success")
	} else {
		fmt.Fprintf(wri, "Unable to register devConfig: The user does not exist!")
	}
}

func (f F) getUsers(wri http.ResponseWriter, r *http.Request) {
	if users.Existing() {
		fmt.Fprintf(wri, strings.Join(users.GetAllIDs(), "\n"))
	} else {
		fmt.Fprintf(wri, "nil")
	}
}

func (f F) infraSetup(wri http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(wri, "success")
	setup.F.Router(setup.F{}, []string{"infrastructure", "delete"}, false)
}

func (f F) getURLQueryParam(req *http.Request, key string) string {
	return req.URL.Query().Get(key)
}

func (f F) saveToFile(req *http.Request, filePath string) bool {
	body := f.readRequestBody(req)
	exists := false
	if !(filesystem.Exists(filePath)) {
		filesystem.CreateEmptyFile(filePath)
	} else {
		exists = true
	}
	if req.Header["Content-Type"][0] == "application/yaml" {
		body = yaml.ConvertJSON(body)
	}
	filesystem.SaveByteArrayToFile(body, filePath)
	return exists
}

func (f F) readRequestBody(req *http.Request) []byte {
	body, err := ioutil.ReadAll(req.Body)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to read the received file for saving!", false, false, true)
	return body
}
