package routines

import (
	"strconv"

	"github.com/neurafuse/neurakube/infrastructure"
	"github.com/neurafuse/tools-go/config"
	infraConfig "github.com/neurafuse/tools-go/config/infrastructure"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/timing"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

var clusterSelfDeletionActive bool

func (f F) Router() {
	if infraConfig.F.ProviderIDIsActive(infraConfig.F{}, "gcloud") {
		if config.APILocationCluster() {
			if config.Setting("get", "infrastructure", "Spec.Cluster.SelfDeletion.Active", "") == "true" {
				if !clusterSelfDeletionActive {
					go f.clusterSelfDeletion()
				}
			}
		}
	}
}

func (f F) clusterSelfDeletion() {
	var rName string = "clusterSelfDeletion"
	timeDuration, _ := strconv.Atoi(config.Setting("get", "infrastructure", "Spec.Cluster.SelfDeletion.TimeDurationHours", ""))
	clusterSelfDeletionActive = true
	logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Activated "+rName, 0)
	var success bool
	for ok := true; ok; ok = !success {
		if timing.TimeDurationPassed(logging.LogTimeLast, timing.GetCurrentTime(), timeDuration, "h") { // TODO: Also save LogTimeLast to config (restart bug)
			logging.Log([]string{"\n", vars.EmojiAPI, vars.EmojiWarning}, "Triggered "+rName+"!", 0)
			logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiProcess}, rName+": Deleting setup..", 0)
			infrastructure.F.Router(infrastructure.F{}, []string{"infrastructure", "delete"}, false)
		}
		timing.Sleep(1, "m")
	}
}
