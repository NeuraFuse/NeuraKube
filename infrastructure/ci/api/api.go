package api

import (
	ci ".."
	"../../../../tools-go/config"
	infraconfig "../../../../tools-go/config/infrastructure"
	"../../../../tools-go/container"
	"../../../../tools-go/env"
	"../../../../tools-go/errors"
	"../../../../tools-go/kubernetes/namespaces"
	"../../../../tools-go/logging"
	"../../../../tools-go/runtime"
	"../../../../tools-go/terminal"
	"../../../../tools-go/vars"
)

type F struct{}

var context string = env.F.GetContext(env.F{}, runtime.F.GetCallerInfo(runtime.F{}, true), false)

func (f F) Router(cliArgs []string, routeAssistant bool) {
	var action string
	if routeAssistant || len(cliArgs) < 2 {
		action = terminal.GetUserSelection("Which "+context+" setup action do you want to start?", []string{"create", "update", "restart", "recreate", "delete"}, false, false)
	} else {
		action = cliArgs[1]
	}
	if action == "create" || action == "cre" {
		f.Create()
	} else if action == "update" || action == "up" {
		f.update()
	} else if action == "restart" || action == "res" {
		f.restart()
	} else if action == "recreate" || action == "rec" {
		f.delete()
		f.Create()
	} else if action == "delete" || action == "del" {
		f.delete()
	}
}

func (f F) Create() {
	if f.EvalAction("create") {
		namespaces.F.Create(namespaces.F{}, namespaces.Default)
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiProcess}, "Creating "+vars.NeuraKubeName+"..\n", 0)
		if config.Setting("get", "dev", "Spec.Status", "") == "active" {
			container.F.CheckUpdates(container.F{}, f.GetContext(), true, false)
		}
		ci.F.Create(ci.F{}, f.GetNamespace(), f.GetContext(), container.F.GetImgAddrs(container.F{}, f.GetContext(), false, false), f.getResources(), ci.F.GetClusterIP(ci.F{}, 1, 9), f.getVolumes(), f.GetContainerPorts())
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Created "+vars.NeuraKubeName+".\n", 0)
	}
}

func (f F) update() {
	if f.EvalAction("update") {
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiProcess}, "Updating "+vars.NeuraKubeName+"..\n", 0)
		if config.Setting("get", "dev", "Spec.Status", "") == "active" {
			container.F.CheckUpdates(container.F{}, f.GetContext(), true, false)
		}
		ci.F.Update(ci.F{}, f.GetNamespace(), f.GetContext(), container.F.GetImgAddrs(container.F{}, f.GetContext(), false, false), f.getResources(), f.getVolumes(), f.GetContainerPorts())
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Updated "+vars.NeuraKubeName+".\n", 0)
	}
}

func (f F) restart() {
	if f.EvalAction("update") {
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiProcess}, "Restarting "+vars.NeuraKubeName+"..", 0)
		ci.F.Update(ci.F{}, f.GetNamespace(), f.GetContext(), container.F.GetImgAddrs(container.F{}, f.GetContext(), false, false), f.getResources(), f.getVolumes(), f.GetContainerPorts())
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Restarted "+vars.NeuraKubeName+".\n", 0)
	}
}

func (f F) delete() {
	if f.EvalAction("delete") {
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiProcess}, "Deleting "+vars.NeuraKubeName+"..", 0)
		ci.F.Delete(ci.F{}, f.GetNamespace(), f.GetContext(), f.getVolumes())
		logging.Log([]string{"", vars.EmojiAPI, vars.EmojiSuccess}, "Deleted "+vars.NeuraKubeName+".\n", 0)
	}
}

func (f F) EvalAction(action string) bool {
	if !config.ValidSettings("infrastructure", vars.NeuraKubeNameRepo, true) {
		infraconfig.F.SetNeuraKubeSpec(infraconfig.F{})
	}
	var actionValid bool
	if action == "create" {
		if ci.F.Exists(ci.F{}, f.GetNamespace(), f.GetContext()) {
			logging.Log([]string{"", vars.EmojiProcess, ""}, vars.NeuraKubeName+" is already created. You can update or delete it.", 0)
			action := terminal.GetUserSelection("Which action do you want to start?", []string{"update", "delete"}, false, false)
			if action == "update" {
				f.update()
			} else if action == "delete" || action == "del" {
				f.delete()
			}
		} else {
			actionValid = true
		}
	} else if action == "connect" || action == "update" {
		if config.Setting("get", "dev", "Spec.API.Address", "") != "localhost" {
			if !ci.F.Exists(ci.F{}, f.GetNamespace(), f.GetContext()) {
				actionPostfix := " "
				if action == "connect" {
					actionPostfix = " to "
				}
				logging.Log([]string{"", vars.EmojiAPI, vars.EmojiWarning}, "Unable to "+action+actionPostfix+vars.NeuraKubeName+" because it is not (fully) created yet.", 0)
				autoCreation := config.Setting("get", "infrastructure", "Spec.Neurakube.AutoCreation", "")
				if autoCreation != "true" {
					sel := terminal.GetUserSelection("Do you want to create it?", []string{}, false, true)
					if sel == "Yes" {
						sel = terminal.GetUserSelection("Do you want to enable the API auto creation for future situations like this?", []string{}, false, true)
						if sel == "Yes" {
							config.Setting("set", "infrastructure", "Spec.Neurakube.AutoCreation", "true")
						}
					} else {
						terminal.Exit(0, "")
					}
				} else {
					logging.Log([]string{"", vars.EmojiAPI, vars.EmojiInfra}, vars.NeuraKubeName+" auto creation is active.", 0)
				}
				f.Create()
			} else {
				actionValid = true
			}
		} else {
			actionValid = true
		}
	} else if action == "delete" {
		actionValid = true
	} else {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Unsupported action to evaluate: "+action, true, true, true)
	}
	return actionValid
}

func (f F) GetContext() string {
	return env.F.GetContext(env.F{}, vars.NeuraKubeNameRepo, false)
}

func (f F) GetNamespace() string {
	return namespaces.Default
}

func (f F) getVolumes() [][]string {
	return [][]string{
		{f.getContainerUserPath(), config.Setting("get", "infrastructure", "Spec.Neurakube.VolumeSizeGB", "") + "Gi"},
	}
}

func (f F) getContainerUserPath() string {
	return env.F.GetContainerWorkingDir(env.F{}) + "users"
}

func (f F) GetContainerPorts() [][]string {
	return [][]string{{"1111", "TCP"}}
}

func (f F) getResources() string {
	return ""
}
