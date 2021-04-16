package nodepools

import (
	"fmt"
	"strconv"

	"../../../../../tools-go/config"
	"../../../../../tools-go/errors"
	"../../../../../tools-go/logging"
	"../../../../../tools-go/objects/strings"
	"../../../../../tools-go/runtime"
	"../../../../../tools-go/vars"
	"../clients"
	gconfig "../config"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

type F struct{}

func (f F) Get(logResult bool) []string {
	var nodepools []string
	ctx, client := clients.F.GetContainer(clients.F{})
	request := &containerpb.ListNodePoolsRequest{
		ProjectId: config.Setting("get", "infrastructure", "Spec.Gcloud.ProjectID", ""),
		Zone:      config.Setting("get", "infrastructure", "Spec.Gcloud.Zone", ""),
		ClusterId: config.Setting("get", "infrastructure", "Spec.Cluster.Name", ""),
	}
	list, err := client.ListNodePools(ctx, request)
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to list nodepools!", false, true, true) {
		for i, d := range list.NodePools {
			if logResult {
				fmt.Printf(" ["+strconv.Itoa(i)+"] %s\n", d.Name)
			}
			nodepools = append(nodepools, d.Name)
		}
	}
	return nodepools
}

func (f F) Create(context, accType string) {
	logging.Log([]string{"", vars.EmojiInfra, vars.EmojiProcess}, "Creating nodepool "+context+"..", 0)
	if !f.Exists(context) {
		ctx, client := clients.F.GetContainer(clients.F{})
		accMachineType := config.Setting("get", "infrastructure", "Spec.Gcloud.Accelerator."+strings.ToUpper(accType)+".MachineType", "")
		request := &containerpb.CreateNodePoolRequest{
			ProjectId: config.Setting("get", "infrastructure", "Spec.Gcloud.ProjectID", ""),
			Zone:      config.Setting("get", "infrastructure", "Spec.Gcloud.Zone", ""),
			ClusterId: config.Setting("get", "infrastructure", "Spec.Cluster.Name", ""),
			NodePool:  gconfig.F.NodePoolConfigSingle(gconfig.F{}, context, accMachineType, accType),
		}
		_, err := client.CreateNodePool(ctx, request)
		if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to create nodepool "+context+"!", false, false, true) {
			logging.Log([]string{"", vars.EmojiInfra, vars.EmojiSuccess}, "Created nodepool "+context+".", 0)
		}
	} else {
		logging.Log([]string{"", vars.EmojiInfra, vars.EmojiWarning}, "Nodepool "+context+" already exists!", 0)
	}
}

func (f F) Delete(context string) {
	if f.Exists(context) {
		logging.Log([]string{"", vars.EmojiInfra, vars.EmojiProcess}, "Deleting nodepool "+context+"..", 0)
		ctx, client := clients.F.GetContainer(clients.F{})
		request := &containerpb.DeleteNodePoolRequest{
			ProjectId:  config.Setting("get", "infrastructure", "Spec.Gcloud.ProjectID", ""),
			Zone:       config.Setting("get", "infrastructure", "Spec.Gcloud.Zone", ""),
			ClusterId:  config.Setting("get", "infrastructure", "Spec.Cluster.Name", ""),
			Name:       context,
			NodePoolId: context,
		}
		_, err := client.DeleteNodePool(ctx, request)
		if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to delete nodepool "+context+".", false, false, true) {
			logging.Log([]string{"", vars.EmojiInfra, vars.EmojiSuccess}, "Deleted nodepool "+context+".", 0)
		}
	} else {
		logging.Log([]string{"", vars.EmojiInfra, vars.EmojiWarning}, "Nodepool "+context+" already deleted.", 0)
	}
}

func (f F) Exists(context string) bool {
	return strings.ArrayContains(f.Get(false), context)
}
