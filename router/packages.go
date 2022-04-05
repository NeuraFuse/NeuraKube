package router

import (
	"github.com/neurafuse/neurakube/server"
	"github.com/neurafuse/tools-go/ci"
	conRun "github.com/neurafuse/tools-go/container/runtime"
	"github.com/neurafuse/tools-go/data"
)

type Packages struct {
	Server    server.F
	Ci        ci.F
	Data      data.F
	Container conRun.F
}
