package infrastructure

import (
	"github.com/neurafuse/tools-go/api/client"
	"github.com/neurafuse/tools-go/cloud/providers/gcloud"
	"github.com/neurafuse/tools-go/cloud/providers/gcloud/clusters"
	gcloudConfig "github.com/neurafuse/tools-go/cloud/providers/gcloud/config"
	"github.com/neurafuse/tools-go/config"
	infraConfig "github.com/neurafuse/tools-go/config/infrastructure"
	userConfig "github.com/neurafuse/tools-go/config/user"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/errors"
	kubeID "github.com/neurafuse/tools-go/kubernetes/client/id"
	kubeTools "github.com/neurafuse/tools-go/kubernetes/config"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/terminal"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

func (f F) Router(cliArgs []string, routeAssistant bool) {
	f.checkDependencies()
	var action string = f.getAction(cliArgs, routeAssistant)
	var providerID string = infraConfig.F.GetProviderIDActive(infraConfig.F{})
	if providerID == vars.InfraProviderGcloud {
		if action == "delete" {
			if env.F.IsFrameworkActive(env.F{}, vars.NeuraKubeNameID) {
				clusters.F.Delete(clusters.F{})
			} else {
				f.userConfirm(action)
			}
		} else {
			gcloud.F.Router(gcloud.F{}, action, cliArgs)
		}
	} else if providerID == vars.InfraProviderSelfHosted {
		logging.Log([]string{"", vars.EmojiWarning, ""}, "You have currently configured a self hosted infrastructure provider which is not compatible yet with the infrastructure module!", 0)
	} else {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Invalid provider ID argument: "+providerID, true, true, true)
	}
}

func (f F) CheckClusterID() {
	var id string
	if !config.ValidSettings("project", "infra/cluster", true) {
		id = terminal.GetUserSelection("What is the kubernetes (cluster) ID?", []string{"cluster-ai-1"}, true, false) // TODO: Implement existing infra. opts
		config.Setting("set", "project", "Spec.Infrastructure.Cluster.ID", id)
	} else {
		id = config.Setting("get", "project", "Spec.Infrastructure.Cluster.ID", "")
	}
	kubeID.F.SetActive(kubeID.F{}, id)
}

func (f F) getAction(cliArgs []string, routeAssistant bool) string {
	var action string
	if routeAssistant || len(cliArgs) < 2 {
		action = terminal.GetUserSelection("Which "+runtime.F.GetCallerInfo(runtime.F{}, true)+" action do you intend to start?", []string{"create", "inspect", "recreate", "delete"}, false, false)
	} else {
		action = cliArgs[1]
	}
	return action
}

func (f F) CheckDeploymentStatus() {
	f.checkDependencies()
	var checkInfra bool
	if config.DevConfigActive() {
		if config.APILocationCluster() {
			checkInfra = true
		}
	} else {
		checkInfra = true
	}
	if checkInfra {
		kubeTools.F.CheckResources(kubeTools.F{})
	}
}

func (f F) checkDependencies() {
	f.setInfraID()
	config.Setting("init", "infrastructure", "", "")
	f.SetProvider()
	if !config.ValidSettings("infrastructure", "cluster", false) {
		if infraConfig.F.ProviderIDIsActive(infraConfig.F{}, "gcloud") {
			if !config.ValidSettings("infrastructure", "cluster", true) {
				infraConfig.F.SetCluster(infraConfig.F{})
			}
		}
	}
	f.CheckClusterID()
	if !config.ValidSettings("infrastructure", vars.NeuraKubeNameID, true) {
		infraConfig.F.SetNeuraKubeSpec(infraConfig.F{})
	}
}

func (f F) setInfraID() {
	if config.ValidSettings("user", "defaults/infra", false) {
		userConfig.F.SetDefaultInfraID(userConfig.F{})
	} else {
		userConfig.F.SetDefaults(userConfig.F{})
	}
}

func (f F) SetProvider() {
	if config.ValidSettings("infrastructure", vars.InfraProviderGcloud, false) {
		infraConfig.F.SetProviderIDActive(infraConfig.F{}, "gcloud")
		gcloudConfig.F.SetConfigs(gcloudConfig.F{})
	} else {
		infraConfig.F.SetProviderIDActive(infraConfig.F{}, "selfhosted")
	}
}

func (f F) userConfirm(action string) {
	if action == "delete" {
		logging.ProgressSpinner("stop")
		logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiWarning}, "You're about to delete the entire cluster setup.", 0)
		logging.Log([]string{"", vars.EmojiWarning, vars.EmojiInfo}, "Data in hostPath and emptyDir storage pods will be deleted.", 0)
		logging.Log([]string{"", vars.EmojiWarning, vars.EmojiInfo}, "Any running containers will also be deleted.\n", 0)
		var clusterID string = config.Setting("get", "infrastructure", "Spec.Cluster.ID", "")
		var sel string = terminal.GetUserSelection("Do you really want to delete the cluster "+clusterID+"?", []string{}, false, true)
		if sel == "Yes" {
			infraConfig.F.SetClusterRecentlyDeleted(infraConfig.F{}, true)
			client.F.Infrastructure(client.F{}, action)
		}
	}
}
