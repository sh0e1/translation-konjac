package postback

import (
	"encoding/json"

	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/sh0e1/translation-konjac/pkg/language"
)

// Action ...
type Action int

// ...
const (
	SelectLanguageAction = iota + 1
	SelectAudioLanguageAction
)

// Data ...
type Data struct {
	Action    Action
	Language  string
	MessageID string
}

// NewSelectLanguagePostbackAction ...
func NewSelectLanguagePostbackAction(lang language.Language) *linebot.PostbackAction {
	data := &Data{
		Action:   SelectLanguageAction,
		Language: lang.Code,
	}
	bytes, _ := json.Marshal(data)
	return linebot.NewPostbackAction(lang.Label, string(bytes), "", lang.Label)
}

// NewSelectAudioLanguagePostbackAction ...
func NewSelectAudioLanguagePostbackAction(lang language.Language, mid string) *linebot.PostbackAction {
	data := &Data{
		Action:    SelectAudioLanguageAction,
		Language:  lang.Code,
		MessageID: mid,
	}
	bytes, _ := json.Marshal(data)
	return linebot.NewPostbackAction(lang.Label, string(bytes), "", lang.Label)
}
