package twitter

import (
	"../../../../tools-go/vars"
	"../../../../tools-go/logging"
)

type F struct{}

func (f F) Router(module, dataPath string) {
	logging.Log([]string{"", vars.EmojiGlobe, vars.EmojiInspect}, "Starting twitter API..", 0)
}