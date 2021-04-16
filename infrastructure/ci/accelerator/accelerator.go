package accelerator

import (
	ci "../"
	"../base"
	"../../../../tools-go/kubernetes/deployments"
	"../../../../tools-go/kubernetes/namespaces"
	"../../../../tools-go/kubernetes/pods"
	"../../../../tools-go/kubernetes/services"
	"../../../../tools-go/logging"
	"../../../../tools-go/vars"
	"../../../container"
)

type F struct{}

func (f F) Prepare(context, resType string) string {
	logging.Log([]string{"\n", vars.EmojiRemote, vars.EmojiProcess}, "Preparing "+resType+"..", 0)
	return ci.F.NodeScheduling(ci.F{}, context)
}

func (f F) Create(context, namespace, resType, imgAddrs, resources string, volumes [][]string) string {
	namespaces.F.Create(namespaces.F{}, namespace)
	if !deployments.F.Exists(deployments.F{}, namespace, context) {
		logging.Log([]string{"\n", vars.EmojiRemote, vars.EmojiProcess}, "Creating "+resType+"..", 0)
		ci.F.NodeScheduling(ci.F{}, context)
		ci.F.Create(ci.F{}, namespace, context, imgAddrs, resources, ci.F.GetClusterIP(ci.F{}, 10, 20), volumes, base.F.GetContainerPorts(base.F{}, context))
		logging.Log([]string{"", vars.EmojiRemote, vars.EmojiSuccess}, "Created "+resType+".\n", 0)
	} else {
		f.Update(context, namespace, resType, imgAddrs, resources, volumes)
	}
	contextID := ci.F.GetContextID(ci.F{}, context)
	if context == "remote" {
		pods.F.Logs(pods.F{}, namespace, contextID, container.F.GetProjectSyncWaitMsg(container.F{}), false, ci.F.GetInitWaitDuration(ci.F{}, context))
	}
	ip := services.F.GetLoadBalancerIP(services.F{}, namespace, contextID)
	logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Sending "+resType+" service IP to requesting "+vars.NeuraCLIName+" client..", 0)
	return ip
}

func (f F) Update(context, namespace, resType, imgAddrs, resources string, volumes [][]string) {
	logging.Log([]string{"\n", vars.EmojiRemote, vars.EmojiProcess}, "Updating "+resType+"..", 0)
	ci.F.NodeScheduling(ci.F{}, context)
	ci.F.Update(ci.F{}, namespace, context, imgAddrs, resources, volumes, base.F.GetContainerPorts(base.F{}, context))
	logging.Log([]string{"", vars.EmojiRemote, vars.EmojiSuccess}, "Updated "+resType+".\n", 0)
}

func (f F) Delete(context, namespace, resType string, volumes [][]string) string {
	logging.Log([]string{"\n", vars.EmojiRemote, vars.EmojiProcess}, "Deleting "+resType+"..", 0)
	ci.F.Delete(ci.F{}, namespace, context, volumes)
	logging.Log([]string{"", vars.EmojiRemote, vars.EmojiSuccess}, "Deleted "+resType+".\n", 0)
	return "success"
}