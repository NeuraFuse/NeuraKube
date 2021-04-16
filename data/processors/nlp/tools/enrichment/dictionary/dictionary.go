package dictionary

import (
	"../../../../../../../tools-go/errors"
	"../../../../../../../tools-go/logging"
	"../../../../../../../tools-go/objects/strings"
	"../../../../../../../tools-go/runtime"
	"../../../../../../../tools-go/vars"
	"./api"
	prose "github.com/jdkato/prose"
)

type F struct{}

func (f F) AddDefinitions(input string, wordTypes []string) string {
	logging.Log([]string{"", vars.EmojiInspect, vars.EmojiProcess}, "Adding definitions ("+strings.Join(wordTypes, "")+")..", 0)
	input = strings.Join(f.detect(input, wordTypes), " ")
	return input
}

func (f F) detect(input string, wordTypes []string) []string {
	var output []string
	doc, err := prose.NewDocument(input)
	errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to process input!", false, false, true)
	var iWords int
	var wordsWithAddedDef []string
	for i, tok := range doc.Tokens() {
		if i == 50 {
			break
		}
		word := tok.Text
		output = append(output, word)
		for _, t := range wordTypes {
			if strings.HasPrefix(tok.Tag, t) {
				if !strings.ArrayContains(wordsWithAddedDef, word) {
					output = append(output, api.F.GetDefinition(api.F{}, word)...)
					wordsWithAddedDef = append(wordsWithAddedDef, word)
				}
			}
		}
		if iWords == 16 {
			output = append(output, "\n")
			iWords = 0
		}
		iWords++
	}
	return output
}
