package request

import (
	"bytes"

	"../../../../../tools-go/api/client"
	"../../../../../tools-go/objects/strings"
)

type F struct{}

func (f F) Router(context string) string {
	var response string
	request := "{\"context\": \"" + context + "\"}"
	body := bytes.NewReader(strings.ToBytes(request))
	response = client.F.Router(client.F{}, "inference/gpt", "POST", "user/infrastructure", "", "", "", body)
	return response
}