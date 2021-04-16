package ci

import (
	"../../../tools-go/config"
	"../../../tools-go/crypto"
	"../../../tools-go/env"
	"../../../tools-go/errors"
	"../../../tools-go/kubernetes/deployments"
	"../../../tools-go/kubernetes/resources"
	"../../../tools-go/kubernetes/services"
	"../../../tools-go/kubernetes/volumes"
	"../../../tools-go/logging"
	"../../../tools-go/objects/strings"
	"../../../tools-go/releases"
	"../../../tools-go/runtime"
	"../../../tools-go/terminal"
	"../../../tools-go/vars"
	"../providers/gcloud/nodepools"
)

type F struct{}

var contextLocal string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)

func (f F) Router(cliArgs []string, routeAssistant bool) {
	var module string
	var modules []string = []string{"releases"}
	if len(cliArgs) < 2 {
		module = terminal.GetUserSelection("Which "+contextLocal+" module do you want to start?", modules, false, false)
	} else {
		module = cliArgs[1]
	}
	switch module {
	case modules[0]:
		releases.F.Router(releases.F{}, cliArgs)
	}
}

func (f F) Exists(namespace, context string) bool {
	contextID := f.GetContextID(context) // TODO: Refactor
	if !volumes.F.Exists(volumes.F{}, namespace, contextID) {
		return false
	} else if !deployments.F.Exists(deployments.F{}, namespace, contextID) {
		return false
	} else if !services.F.Exists(services.F{}, namespace, contextID) {
		return false
	}
	return true
}

func (f F) Create(namespace, context, imageAddrs, accType, clusterIP string, volumesSpec, containerPorts [][]string) {
	if accType != "" {
		resources.Check(context, accType)
	}
	contextID := f.GetContextID(context) // TODO: Refactor
	volumes.F.Create(volumes.F{}, namespace, contextID, f.getServiceCluster(context), volumesSpec)
	repoAddrs := config.Setting("get", "dev", "Spec.Containers.Registry.Address", "")
	if repoAddrs != "" {
		imageAddrs = repoAddrs + "/" + imageAddrs
	}
	deployments.F.Create(deployments.F{}, namespace, contextID, imageAddrs, f.getServiceCluster(context), accType, volumesSpec, containerPorts)
	services.F.Create(services.F{}, namespace, contextID, clusterIP, containerPorts)
}

func (f F) NodeScheduling(context string) string {
	go f.CreateNodePool(context)
	return "success"
}

func (f F) CreateNodePool(context string) {
	if vars.InfraProviderActive != vars.InfraProviderSelfHosted {
		nodepools.F.Create(nodepools.F{}, f.getServiceCluster(context), f.GetType(context, false))
	} else {
		errors.Check(nil, contextLocal, "Unable to create nodepool for selfhosted setup!", true, false, true)
	}
}

func (f F) DeleteNodePool(context string) {
	if vars.InfraProviderActive != vars.InfraProviderSelfHosted {
		nodepools.F.Delete(nodepools.F{}, f.getServiceCluster(context))
	} else {
		errors.Check(nil, contextLocal, "Unable to delete nodepool for selfhosted setup!", true, false, true)
	}
}

func (f F) RecreateNodePool(context string) {
	logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiProcess}, "Recreating nodepool "+context+"..", 0)
	f.DeleteNodePool(context)
	f.CreateNodePool(context)
}

func (f F) Update(namespace, context, imageAddrs, accType string, volumes, containerPorts [][]string) {
	if accType != "" {
		resources.Check(context, accType)
	}
	contextID := f.GetContextID(context) // TODO: Refactor
	deployments.F.Delete(deployments.F{}, namespace, contextID)
	repoAddrs := config.Setting("get", "dev", "Spec.Containers.Registry.Address", "")
	if repoAddrs != "" {
		imageAddrs = repoAddrs + "/" + imageAddrs
	}
	deployments.F.Create(deployments.F{}, namespace, contextID, imageAddrs, f.getServiceCluster(context), accType, volumes, containerPorts)
}

func (f F) Delete(namespace, context string, volumesSpec [][]string) {
	contextID := f.GetContextID(context) // TODO: Refactor
	deployments.F.Delete(deployments.F{}, namespace, contextID)
	volumes.F.Delete(volumes.F{}, namespace, contextID, volumesSpec)
	services.F.Delete(services.F{}, namespace, contextID)
}

func (f F) getServiceCluster(context string) string {
	var serviceCluster string = vars.OrganizationNameRepo
	if context != vars.NeuraKubeNameRepo {
		dedicated := config.Setting("get", "infrastructure", "Spec."+strings.Title(context)+".NodePools.Dedicated", "")
		if dedicated == "true" {
			serviceCluster = serviceCluster + "-" + context + "-" + f.GetType(context, false)
		} else {
			serviceCluster = serviceCluster + "-" + f.GetType(context, false)
		}
	}
	return serviceCluster
}

func (f F) GetContextID(context string) string {
	return context + "-1"
}

func (f F) GetClusterIP(min, max int) string {
	baseIP := "10.24.0."
	return baseIP + strings.ToString(crypto.RandomInt(min, max))
}

func (f F) GetInitWaitDuration(context string) int {
	var waitDuration int
	var accType string = f.GetType(context, false)
	if accType == "tpu" {
		waitDuration = 20
	} else {
		waitDuration = 8
	}
	return waitDuration
}

func (f F) GetType(context string, upperCase bool) string {
	if context == "develop/remote" { // TODO: Refactor
		context = "remote"
	}
	var resType string = config.Setting("get", "infrastructure", "Spec."+strings.Title(context)+".Type", "")
	if upperCase {
		resType = strings.ToUpper(resType)
	}
	return resType
}
