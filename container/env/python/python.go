package python

import (
	"../../../../tools-go/exec"
	"../../../../tools-go/filesystem"
	"../../../../tools-go/kubernetes/resources"
	"../../../../tools-go/logging"
	"../../../../tools-go/timing"
	"../../../../tools-go/vars"
	"../../../../tools-go/env"
)

type F struct{}

func (f F) Router(project, module, pathExec, projectSyncWaitMsg string) {
	if env.F.Container(env.F{}) {
		resources.Check("container", "tpu")
	}
	pathExec = pathExec + ".py"
	firstRun := true
	for {
		if filesystem.Exists(pathExec) {
			if firstRun {
				timing.TimeOut(5, "s")
				firstRun = false
			}
			logging.PartingLine()
			logging.Log([]string{"", vars.EmojiProject, vars.EmojiInfo}, "Starting project "+project+"..\n", 0)
			argsExec := []string{pathExec}
			exec.WithLiveLogs("python", argsExec, true)
			logging.Log([]string{"", vars.EmojiLink, vars.EmojiInfo}, "Auto restart after 1s..", 0)
			logging.PartingLine()
		} else {
			logging.Log([]string{"", vars.EmojiLink, vars.EmojiInfo}, projectSyncWaitMsg, 0)
		}
		timing.TimeOut(1, "s")
	}
}
