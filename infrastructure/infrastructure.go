package infrastructure

import (
	apiRouter "../../tools-go/api"
	"../../tools-go/config"
	infraConfig "../../tools-go/config/infrastructure"
	"../../tools-go/crypto"
	"../../tools-go/env"
	"../../tools-go/errors"
	"../../tools-go/logging"
	"../../tools-go/runtime"
	"../../tools-go/terminal"
	"../../tools-go/users"
	"../../tools-go/vars"
	"./ci/api"
	"./providers/gcloud"
	"./providers/gcloud/clusters"
	gcloudConfig "./providers/gcloud/config"
)

type F struct{}

func (f F) Router(cliArgs []string, routeAssistant bool) {
	f.checkSettings()
	var provider string = vars.InfraProviderActive
	var action string
	if routeAssistant || len(cliArgs) < 2 {
		action = terminal.GetUserSelection("Which "+runtime.F.GetCallerInfo(runtime.F{}, true)+" action do you want to start?", []string{"create", "inspect", "recreate", "delete"}, false, false)
	} else {
		action = cliArgs[1]
	}
	if provider == vars.InfraProviderGcloud {
		if action == "delete" || action == "del" {
			if env.F.ActiveFramework(env.F{}, vars.NeuraKubeNameRepo) {
				clusters.F.Delete(clusters.F{})
			} else {
				f.userConfirm(action)
			}
		} else {
			success := gcloud.F.Router(gcloud.F{}, action, cliArgs)
			if success {
				api.F.Create(api.F{})
			}
		}
	} else if provider == vars.InfraProviderSelfHosted {
		logging.Log([]string{"", vars.EmojiWarning, ""}, "You have currently configured a self hosted infrastructure provider which is not compatible yet with the infrastructure module!", 0)
	} else {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Invalid provider argument: "+provider, true, true, true)
	}
}

func (f F) checkSettings() {
	config.Setting("init", "infrastructure", "", "")
	if config.ValidSettings("infrastructure", "kube", false) {
		infraConfig.F.SetKubeConfig(infraConfig.F{})
	}
	if config.ValidSettings("infrastructure", vars.InfraProviderGcloud, false) {
		gcloudConfig.F.SetConfigs(gcloudConfig.F{})
	}
	if !config.ValidSettings("infrastructure", "cluster", false) {
		config.Setting("set", "infrastructure", "Spec.Cluster.Auth.Password", crypto.RandomString(64))
		if vars.InfraProviderActive != vars.InfraProviderSelfHosted {
			if !config.ValidSettings("infrastructure", "cluster", true) {
				infraConfig.F.SetCluster(infraConfig.F{})
			}
		}
	}
	if !config.ValidSettings("infrastructure", "kube", false) && !config.ValidSettings("infrastructure", vars.InfraProviderGcloud, false) {
		optionSelfhosted := "Self hosted kubernetes cluster (auth via existing kubeconfig)"
		optionGcloud := vars.InfraProviderGcloud
		selection := terminal.GetUserSelection("Please choose an infrastructure provider", []string{optionSelfhosted, optionGcloud}, false, false)
		if selection == optionSelfhosted {
			config.ValidSettings("infrastructure", "kube", true)
			infraConfig.F.SetKubeConfig(infraConfig.F{})
		} else if selection == optionGcloud {
			config.ValidSettings("infrastructure", vars.InfraProviderGcloud, true)
			gcloudConfig.F.SetConfigs(gcloudConfig.F{})
		}
	}
	if !config.ValidSettings("infrastructure", vars.NeuraKubeNameRepo, true) {
		infraConfig.F.SetNeuraKubeSpec(infraConfig.F{})
	}
}

func (f F) userConfirm(action string) {
	if action == "delete" || action == "del" {
		logging.ProgressSpinner("stop")
		logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiWarning}, "You're about to delete the entire cluster setup.", 0)
		logging.Log([]string{"", vars.EmojiWarning, vars.EmojiInfo}, "This will also delete all persistent (hostdir) volumes (storages) that are stored within it.\n", 0)
		sel := terminal.GetUserSelection("Do you really want to delete the entire infrastructure including all (hostdir) volumes?", []string{}, false, true)
		if sel == "Yes" {
			users.SetClusterRecentlyDeleted(true)
			apiRouter.F.Infrastructure(apiRouter.F{}, action)
		}
	}
}
