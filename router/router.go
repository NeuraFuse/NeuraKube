package router

import (
	"os"

	"../../tools-go/build"
	"../../tools-go/env"
	"../../tools-go/errors"
	"../../tools-go/logging"
	"../../tools-go/objects"
	"../../tools-go/objects/strings"
	"../../tools-go/runtime"
	"../../tools-go/terminal"
	"../../tools-go/updater"
	"../../tools-go/vars"
)

type F struct{}

func (f F) Router() {
	f.startup()
	cliArgs := strings.ArrayRemoveString(os.Args, os.Args[0])
	var pack string
	if len(cliArgs) == 0 {
		pack = terminal.GetUserSelection("Which "+vars.NeuraKubeName+" module do you want to start?", []string{"server", "Spec.container"}, false, false)
	} else {
		pack = cliArgs[0]
	}
	success, _ := objects.CallStructInterfaceFuncByName(Packages{}, strings.Title(pack), "Router", cliArgs, false)
	if !success {
		errors.Check(nil, runtime.F.GetCallerInfo(runtime.F{}, false), "The package "+pack+" does not exist!", true, true, true)
	}
	terminal.Exit(0, "")
}

func (f F) startup() {
	terminal.Init(false)
	updater.F.Check(updater.F{})
	build.F.CheckUpdates(build.F{}, env.F.GetActive(env.F{}, false), true)
	logging.Log([]string{"", vars.EmojiAstronaut, vars.EmojiSuccess}, "Ready to go.\n", 0)
}
