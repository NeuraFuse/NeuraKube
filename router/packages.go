package router

import (
	"../container"
	"../data"
	"../infrastructure/ci"
	"../server"
)

type Packages struct {
	Server server.F
	Ci     ci.F
	Data   data.F
	Container container.F
}