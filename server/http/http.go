package http

import (
	"net/http"

	"github.com/neurafuse/neurakube/server/requests"
	"github.com/neurafuse/tools-go/ci/api"
	"github.com/neurafuse/tools-go/crypto/rsa"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/timing"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

var apiPort string = api.F.GetContainerPorts(api.F{})[0][0]

func (f F) activateRoutes() {
	var routes []string = []string{"/inspect", "/inspect/healthcheck", "/inspect/version",
		"/inspect/api/init", "/users", "/user/create", "/user/infrastructure", "/user/infrastructure/setup", "/user/project",
		"/user/infrastructure/auth/kubeconfig/create", "/user/infrastructure/auth/gcloud/create",
		"/user/devconfig/create", "/develop", "/develop/delete", "/app", "/app/delete",
		"/app/inference", "/app/inference/delete", "/app/gpt/inference"}
	logging.Log([]string{"", vars.EmojiLink, vars.EmojiRoute}, "Active routes:\n", 0)
	for _, route := range routes {
		http.Handle(vars.RESTRoutePreamble+route, requests.F.Router(requests.F{}, route))
		logging.Log([]string{"", vars.EmojiRoute, ""}, route, 0)
	}
}

func (f F) StartServing() {
	logging.Log([]string{"", vars.EmojiAPI, vars.EmojiGlobe}, "Starting server..", 0)
	//serverTLSConf, _ := tls.GetCert(vars.NeuraKubeName)
	f.activateRoutes()
	server := &http.Server{
		Addr:    ":" + apiPort,
		Handler: nil,
		//TLSConfig:    serverTLSConf,
		//TLSConfig: tls.TLSConfig(),
		ReadTimeout:  timing.GetTimeDuration(10, "m"),
		WriteTimeout: timing.GetTimeDuration(10, "m"),
		//ErrorLog:     errorLogger(),
		//MaxHeaderBytes: 1 << 20,
	}
	publicKeyFilePath, privateKeyFilePath := rsa.GenerateKeys(vars.NeuraKubeName, env.F.GetAPIHTTPCertPath(env.F{}), false)
	f.selfCheck()
	logging.Log([]string{"", vars.EmojiLink, vars.EmojiSuccess}, vars.NeuraKubeName+" is now reachable.\n", 0)
	errors.Check(server.ListenAndServeTLS(publicKeyFilePath, privateKeyFilePath), runtime.F.GetCallerInfo(runtime.F{}, false), "Failed to start or continue serving!\nCheck if another program is blocking the port "+apiPort+".", false, true, true)
}

func (f F) selfCheck() {
	go f.checkInitStatusRoutine()
	logging.Log([]string{"", vars.EmojiGlobe, vars.EmojiAPI}, "Version: "+vars.NeuraKubeVersion, 0)
	logging.Log([]string{"", vars.EmojiLink, vars.EmojiAPI}, "Port: "+apiPort, 0)
}

func (f F) checkInitStatusRoutine() {
	var initStatusLast bool = true
	for {
		var initStatus bool = requests.F.GetInfrastructureInitStatus(requests.F{})
		if initStatus != initStatusLast {
			if initStatus {
				logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Infrastructure: initialized\n", 0)
			} else {
				logging.Log([]string{"", vars.EmojiAPI, vars.EmojiWarning}, "Infrastructure: uninitialized\n", 0)
			}
			initStatusLast = initStatus
		}
		timing.Sleep(1, "s")
	}
}
