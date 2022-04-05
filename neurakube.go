package main

import (
	"github.com/neurafuse/neurakube/router"
	"github.com/neurafuse/tools-go/env"
	"github.com/neurafuse/tools-go/vars"
)

func main() {
	env.F.SetFramework(env.F{}, vars.NeuraKubeNameID)
	router.F.Router(router.F{})
}
