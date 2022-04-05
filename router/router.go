package router

import (
	"os"

	"github.com/neurafuse/tools-go/errors"
	"github.com/neurafuse/tools-go/logging"
	"github.com/neurafuse/tools-go/objects"
	"github.com/neurafuse/tools-go/objects/strings"
	"github.com/neurafuse/tools-go/runtime"
	"github.com/neurafuse/tools-go/terminal"
	"github.com/neurafuse/tools-go/updater"
	"github.com/neurafuse/tools-go/vars"
)

type F struct{}

func (f F) Router() {
	f.startup()
	var cliArgs []string = strings.ArrayRemoveString(os.Args, os.Args[0])
	var pack string
	if len(cliArgs) == 0 {
		pack = terminal.GetUserSelection("Which "+vars.NeuraKubeName+" module do you want to start?", []string{"server", "container"}, false, false)
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
	logging.Log([]string{"", vars.EmojiAstronaut, vars.EmojiSuccess}, "Ready to go.\n", 0)
}
