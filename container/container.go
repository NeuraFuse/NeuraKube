package container

import (
	"../../tools-go/api/client"
	"../../tools-go/apps/tensorflow/tensorboard"
	"../../tools-go/crypto"
	"../../tools-go/env"
	"../../tools-go/errors"
	"../../tools-go/logging"
	"../../tools-go/projects"
	"../../tools-go/runtime"
	"../../tools-go/terminal"
	"../../tools-go/users"
	"../../tools-go/vars"
	"../data/providers/crawler"

	//"../data/providers/twitter"
	"./env/python"
)

type F struct{}

func (f F) Router(cliArgs []string, routeAssistant bool) {
	logging.Log([]string{"", vars.EmojiContainer, vars.EmojiProcess}, "Starting container..\n", 0)
	f.init()
	project, assistant := f.getProject(cliArgs)
	module := f.getModule(cliArgs, project, assistant)
	pathExec, dataPath := f.getPaths(project)
	f.dataAggregation(module, project, dataPath)
	if module != "modelserver" {
		tensorboard.F.Start(tensorboard.F{})
	}
	python.F.Router(python.F{}, project, module, pathExec, f.GetProjectSyncWaitMsg())
}

func (f F) init() {
	if env.F.Container(env.F{}) {
		f.connectAPI()
	}
}

func (f F) dataAggregation(module, project, dataPath string) {
	switch module {
	case "gpt":
		crawler.F.Router(crawler.F{}, project, dataPath)
	}
}

func (f F) connectAPI() {
	users.SetIDActive(f.getID())
	projects.F.CheckAuth(projects.F{})
	client.F.Connect(client.F{})
}

func (f F) getID() string {
	return "container-" + crypto.RandomString(8)
}

func (f F) getPaths(project string) (string, string) {
	pathExec := projects.F.GetExternalExecPath(projects.F{}, project)
	dataPath := projects.F.GetExternalDataPath(projects.F{}) + "data/"
	return pathExec, dataPath
}

func (f F) getProject(cliArgs []string) (string, bool) {
	var project string
	var assistant bool = true
	if len(cliArgs) > 1 {
		project = cliArgs[1]
		if project == "lightning" {
			assistant = false
		} else {
			errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "Unsupported project argument: "+project, true, false, true)
		}
	}
	if assistant {
		project = terminal.GetUserSelection("Please choose a project", []string{"lightning"}, false, false)
	}
	return project, assistant
}

func (f F) getModule(cliArgs []string, project string, assistant bool) string {
	var module string
	if len(cliArgs) == 3 || assistant {
		module = terminal.GetUserSelection("Please choose a module for the project "+project+"", []string{"gpt", "modelserver"}, false, false)
	} else if len(cliArgs) > 3 {
		module = cliArgs[3]
	}
	return module
}

func (f F) GetProjectSyncWaitMsg() string {
	return "Waiting for module " + projects.F.GetWorkingDir(projects.F{}) + " to get synced in.."
}
