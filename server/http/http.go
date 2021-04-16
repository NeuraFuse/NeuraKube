package http

import (
	"net/http"

	"../../../tools-go/crypto/rsa"
	"../../../tools-go/env"
	"../../../tools-go/errors"
	"../../../tools-go/logging"
	"../../../tools-go/runtime"
	"../../../tools-go/timing"
	"../../../tools-go/vars"
	"../../infrastructure/ci/api"
	"../requests"
)

type F struct{}

var apiPort string = api.F.GetContainerPorts(api.F{})[0][0]

func (f F) activateRoutes() {
	routes := []string{"/inspect", "/inspect/healthcheck", "/inspect/version",
		"/inspect/infrastructure/init", "/users", "/user/create", "/user/infrastructure", "/user/infrastructure/setup", "/user/project",
		"/user/infrastructure/auth/kubeconfig/create", "/user/infrastructure/auth/gcloud/create",
		"/user/devconfig/create", "/develop/remote", "/develop/remote/delete", "/app", "/app/delete",
		"/app/inference", "/app/inference/delete", "/app/gpt/inference"}
	logging.Log([]string{"", vars.EmojiLink, vars.EmojiRoute}, "Active routes:\n", 0)
	for _, route := range routes {
		http.Handle(vars.RESTRoutePreamble+route, requests.F.Router(requests.F{}, route))
		logging.Log([]string{"", vars.EmojiRoute, ""}, route, 0)
	}
}

func (f F) StartServing() {
	logging.Log([]string{"", vars.EmojiAPI, vars.EmojiGlobe}, "Starting server..\n", 0)
	//serverTLSConf, _ := tls.GetCert(vars.NeuraKubeName)
	f.activateRoutes()
	s := &http.Server{
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
	errors.Check(s.ListenAndServeTLS(publicKeyFilePath, privateKeyFilePath), runtime.F.GetCallerInfo(runtime.F{}, false), "Failed to start or continue serving!", false, true, true)
}

func (f F) selfCheck() {
	var configured bool = false
	if requests.F.GetInfrastructureInitStatus(requests.F{}) {
		configured = true
	}
	if configured {
		logging.Log([]string{"", vars.EmojiSettings, vars.EmojiSuccess}, "Infrastructure: initialized\n", 0)
	} else {
		logging.Log([]string{"", vars.EmojiSettings, vars.EmojiWarning}, "Infrastructure: uninitialized\n", 0)
	}
	logging.Log([]string{"", vars.EmojiGlobe, vars.EmojiAPI}, "Version: "+vars.NeuraKubeVersion, 0)
	logging.Log([]string{"", vars.EmojiLink, vars.EmojiAPI}, "Port: "+apiPort, 0)
}
