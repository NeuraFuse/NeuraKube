package gcloud

import (
	"../../../../tools-go/config"
	"../../../../tools-go/errors"
	"../../../../tools-go/io"
	"../../../../tools-go/kubernetes/client/kubeconfig"
	"../../../../tools-go/kubernetes/namespaces"
	"../../../../tools-go/logging"
	"../../../../tools-go/runtime"
	"../../../../tools-go/vars"
	"./clusters"
	gconfig "./config"
)

type F struct{}

func (f F) Router(action string, cliArgs []string) bool { // nodePool("list") , nodePool("create") , nodePool("delete")
	if !config.ValidSettings("infrastructure", vars.InfraProviderGcloud, true) {
		gconfig.F.SetConfigs(gconfig.F{})
	}
	f.apiAvailability()
	logging.ProgressSpinner("start")
	success := false
	if action == "inspect" {
		f.inspect()
	} else if action == "create" {
		success = clusters.F.Create(clusters.F{})
		if success {
			kubeconfig.F.Create(kubeconfig.F{}, namespaces.Default)
		}
	} else if action == "recreate" {
		clusters.F.Delete(clusters.F{})
		success = clusters.F.Create(clusters.F{})
		if success {
			kubeconfig.F.Create(kubeconfig.F{}, namespaces.Default)
		}
	} else {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Invalid action argument: "+cliArgs[1], true, true, true)
	}
	if !success {
		f.inspect()
	}
	return success
}

func (f F) inspect() {
	logging.Log([]string{"", vars.EmojiInfra, vars.EmojiInspect}, "Starting "+runtime.F.GetCallerInfo(runtime.F{}, true)+" inspection..\n", 0)
	clusters.F.Get(clusters.F{}, true)
}

func (f F) apiAvailability() {
	apiURL := "status.cloud.google.com"
	if io.F.Reachable(io.F{}, apiURL) {
		logging.Log([]string{"", vars.EmojiInfra, vars.EmojiSuccess}, "Provider is reachable.", 0)
	} else {
		logging.Log([]string{"", vars.EmojiInfra, vars.EmojiWarning}, "Provider is not reachable!", 0)
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "There is probably an error with networking on your side or at "+vars.InfraProviderGcloud+"!", true, true, true)
	}
}
