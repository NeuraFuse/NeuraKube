package main

import (
	"../tools-go/env"
	"../tools-go/vars"
	"./router"
)

func main() {
	env.F.SetFramework(env.F{}, vars.NeuraKubeNameRepo)
	router.F.Router(router.F{})
}