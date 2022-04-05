package requests

import (
	"fmt"
	"io/ioutil"
	"net/http"

	infra "github.com/neurafuse/neurakube/infrastructure"
	"github.com/neurafuse/neurakube/server/auth"
	"github.com/neurafuse/neurakube/server/routines"
	"github.com/neurafuse/tools-go/ci/app"
	"github.com/neurafuse/tools-go/ci/develop"
	"github.com/neurafuse/tools-go/ci/inference"
	inferenceRequest "github.com/neurafuse/tools-go/ci/inference/request"
	gconfig "github.com/neurafuse/tools-go/cloud/providers/gcloud/config"
	"github.com/neurafuse/tools-go/config"
	devConfig "github.com/neurafuse/tools-go/config/dev"
	infraConfig "github.com/neurafuse/tools-go/config/infrastructure"
	projectConfig "github.com/neurafuse/tools-go/config/project"
	userConfig "github.com/neurafuse/tools-go/config/user"
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/filesystem"
	infraID "github.com/neurafuse/tools-go/infrastructures/id"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/readers/yaml"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/users"
	usersID "github.com/neurafuse/tools-go/users/id"
	"github.com/neurafuse/tools-go/vars"
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
	} else if route == "/inspect/api/init" {
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
	} else if route == "/user/devconfig/create" {
		handler = auth.F.Check(auth.F{}, f.createUserdevConfig)
	} else if route == "/develop" {
		handler = auth.F.Check(auth.F{}, f.develop)
	} else if route == "/develop/delete" {
		handler = auth.F.Check(auth.F{}, f.deleteDevelop)
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
	if f.GetInfrastructureInitStatus() {
		fmt.Fprintf(wri, "initialized")
	} else {
		fmt.Fprintf(wri, "uninitialized")
	}
}

func (f F) GetInfrastructureInitStatus() bool {
	return users.F.Existing(users.F{})
}

func (f F) getHealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func (f F) getVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, vars.NeuraKubeVersion)
}

func (f F) develop(wri http.ResponseWriter, req *http.Request) {
	var action string = f.getURLQueryParam(req, "action")
	if action == "prepare" {
		fmt.Fprintf(wri, "success/"+develop.F.Prepare(develop.F{}))
	} else if action == "create" {
		fmt.Fprintf(wri, "success/"+develop.F.Create(develop.F{}))
	} else {
		fmt.Fprintf(wri, "error/Unsupported action: "+action)
	}
}

func (f F) deleteDevelop(wri http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(wri, develop.F.Delete(develop.F{}))
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
	var action string = f.getURLQueryParam(req, "action")
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
	var createUser bool
	var userType string
	if !f.GetInfrastructureInitStatus() {
		createUser = true
		userType = "Admin user "
	} else {
		if config.Setting("get", "server", "Spec.Users.Admin", "") == usersID.F.GetActive(usersID.F{}) {
			if !users.F.Exists(users.F{}, usersID.F.GetActive(usersID.F{})) {
				createUser = true
				userType = "User "
			} else {
				fmt.Fprintf(wri, "Unable to create user "+usersID.F.GetActive(usersID.F{})+" because it already exists.")
			}
		} else {
			fmt.Fprintf(wri, "The user creation failed because this action is only available to admin users.")
		}
	}
	if createUser {
		users.F.Create(users.F{}, usersID.F.GetActive(usersID.F{}))
		err, _ := f.saveToFile(req, userConfig.F.GetFilePath(userConfig.F{}))
		var errMsg string = "Unable to interact with resource type user!"
		if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), errMsg, false, false, true) {
			config.Setting("init", "user", "Spec.", "")
			auth.F.Config(auth.F{}, req)
			logging.Log([]string{"", vars.EmojiUser, vars.EmojiSuccess}, userType+usersID.F.GetActive(usersID.F{})+" created.", 0)
			config.Setting("set", "server", "Spec.Users.Admin", usersID.F.GetActive(usersID.F{}))
			fmt.Fprintf(wri, "success")
		} else {
			fmt.Fprintf(wri, errMsg)
		}
	}
}

func (f F) infrastructure(wri http.ResponseWriter, req *http.Request) {
	auth.F.SetPaths(auth.F{})
	var infraIDReq string = auth.F.CookieReader(auth.F{}, req.Header["Cookie"], "infraID")
	infraID.F.SetActive(infraID.F{}, infraIDReq)
	err, _ := f.saveToFile(req, infraConfig.F.GetPath(infraConfig.F{}, true))
	var errMsg string = "Unable to interact with resource type infrastructure!"
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), errMsg, false, false, true) {
		config.Setting("init", "infrastructure", "", "")
		logging.Log([]string{"", vars.EmojiInfra, vars.EmojiSuccess}, "Infrastructure configuration synced.", 0)
		fmt.Fprintf(wri, "success")
	} else {
		fmt.Fprintf(wri, errMsg)
	}
}

func (f F) project(wri http.ResponseWriter, req *http.Request) {
	auth.F.SetPaths(auth.F{})
	var filePath string = projectConfig.F.GetFilePath(projectConfig.F{})
	err, _ := f.saveToFile(req, filePath)
	var errMsg string = "Unable to interact with resource type project!"
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), errMsg, false, false, true) {
		config.Setting("init", "dev", "Spec.", "")
		logging.Log([]string{"", vars.EmojiProject, vars.EmojiSuccess}, "Project configuration synced.", 0)
		fmt.Fprintf(wri, "success")
	} else {
		fmt.Fprintf(wri, errMsg)
	}
}

func (f F) createInfraProviderAuthKubeConfig(wri http.ResponseWriter, req *http.Request) {
	var filePath string = infraConfig.F.GetInfraKubeAuthPath(infraConfig.F{}, true)
	err, _ := f.saveToFile(req, filePath)
	var errMsg string = "Unable to synchronize cluster authentication!"
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), errMsg, false, false, true) {
		logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "Cluster authentication synced.", 0)
		routines.F.Router(routines.F{})
		fmt.Fprintf(wri, "success")
	} else {
		fmt.Fprintf(wri, errMsg)
	}
}

func (f F) createInfraProviderAuthGcloud(wri http.ResponseWriter, req *http.Request) {
	var filePath string = infraConfig.F.GetInfraGcloudAuthPath(infraConfig.F{})
	err, _ := f.saveToFile(req, filePath)
	var errMsg string = "Unable to synchronize " + vars.InfraProviderGcloud + " authentication!"
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), errMsg, false, false, true) {
		gconfig.F.SetConfigs(gconfig.F{})
		logging.Log([]string{"", vars.EmojiInfra, vars.EmojiSuccess}, vars.InfraProviderGcloud+" authentication synced.", 0)
		routines.F.Router(routines.F{})
		fmt.Fprintf(wri, "success")
	} else {
		fmt.Fprintf(wri, errMsg)
	}
}

func (f F) createUserdevConfig(wri http.ResponseWriter, req *http.Request) {
	if users.F.Exists(users.F{}, usersID.F.GetActive(usersID.F{})) {
		f.saveToFile(req, devConfig.F.GetFilePath(devConfig.F{}))
		logging.Log([]string{"", vars.EmojiDev, vars.EmojiSuccess}, "devConfig synced.", 0)
		fmt.Fprintf(wri, "success")
	} else {
		fmt.Fprintf(wri, "Unable to register devConfig: The user does not exist!")
	}
}

func (f F) getUsers(wri http.ResponseWriter, r *http.Request) {
	if users.F.Existing(users.F{}) {
		fmt.Fprintf(wri, strings.Join(users.F.GetAllIDs(users.F{}), "\n"))
	} else {
		fmt.Fprintf(wri, "nil")
	}
}

func (f F) infraSetup(wri http.ResponseWriter, req *http.Request) {
	go infra.F.Router(infra.F{}, []string{"infrastructure", "delete"}, false)
	fmt.Fprintf(wri, "success")
}

func (f F) getURLQueryParam(req *http.Request, key string) string {
	return req.URL.Query().Get(key)
}

func (f F) saveToFile(req *http.Request, filePath string) (error, bool) {
	var body []byte = f.readRequestBody(req)
	var exists bool
	if !filesystem.Exists(filePath) {
		filesystem.CreateEmptyFile(filePath)
	} else {
		exists = true
	}
	if req.Header["Content-Type"][0] == "application/yaml" {
		body = yaml.ConvertJSON(body)
	}
	var err error = filesystem.SaveByteArrayToFile(body, filePath)
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "", false, false, false) {
		if filesystem.FileContentEmpty(filePath) {
			err = errors.New("Received file has no content!")
		}
	}
	return err, exists
}

func (f F) readRequestBody(req *http.Request) []byte {
	body, err := ioutil.ReadAll(req.Body)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to read the received file for saving!", false, false, true)
	return body
}
