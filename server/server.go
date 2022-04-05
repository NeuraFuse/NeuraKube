package server

import (
	"github.com/neurafuse/neurakube/server/http"
)

type F struct{}

func (f F) Router(cliArgs []string, routeAssistant bool) {
	http.F.StartServing(http.F{})
}
