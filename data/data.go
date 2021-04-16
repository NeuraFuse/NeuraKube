package data

import (
	"../../tools-go/logging"
	"../../tools-go/vars"
	"../../tools-go/objects"
	"../../tools-go/objects/strings"
)

type F struct{}

func (f F) Router(cliArgs []string, routeAssistant bool) {
	pack := cliArgs[2]
	input := cliArgs[3]
	logging.Log([]string{"\n", vars.EmojiData, vars.EmojiInfo}, "Context "+input+" with package "+pack+"..", 0)
	switch pack {
		case "cc":
			pack = "commoncrawl"
	}
	objects.CallStructInterfaceFuncByName(Packages{}, strings.Title(pack), "Router", input)
}