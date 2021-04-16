package routines

import (
	"strconv"

	"../../../tools-go/config"
	"../../../tools-go/logging"
	"../../../tools-go/timing"
	"../../../tools-go/vars"
	"../../infrastructure"
)

type F struct{}

var clusterSelfDeletionActive bool

func (f F) Router() {
	if vars.InfraProviderActive != vars.InfraProviderSelfHosted {
		if config.Setting("get", "dev", "Spec.API.Address", "") != "localhost" {
			if config.Setting("get", "infrastructure", "Spec.Cluster.SelfDeletion.Active", "") == "true" {
				if !clusterSelfDeletionActive {
					go f.clusterSelfDeletion()
				}
			}
		}
	}
}

func (f F) clusterSelfDeletion() {
	rName := "clusterSelfDeletion"
	timeDuration, _ := strconv.Atoi(config.Setting("get", "infrastructure", "Spec.Cluster.SelfDeletion.TimeDurationHours", ""))
	clusterSelfDeletionActive = true
	logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Activated "+rName, 0)
	var success bool = false
	for ok := true; ok; ok = !success {
		if timing.TimeDurationPassed(logging.LogTimeLast, timing.GetCurrentTime(), timeDuration, "h") { // TODO: Also save LogTimeLast to config (restart bug)
			logging.Log([]string{"\n", vars.EmojiAPI, vars.EmojiWarning}, "Triggered "+rName+"!", 0)
			logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiProcess}, rName+": Deleting setup..", 0)
			infrastructure.F.Router(infrastructure.F{}, []string{"infrastructure", "delete"}, false)
		}
		timing.TimeOut(1, "m")
	}
}
