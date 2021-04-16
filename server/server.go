package server

import (
	"./http"
)

type F struct{}

func (f F) Router(cliArgs []string, routeAssistant bool) {
	http.F.StartServing(http.F{})
}