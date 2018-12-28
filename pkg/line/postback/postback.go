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
	Cursor    int
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
func NewSelectAudioLanguagePostbackAction(mid string, lang language.Language) *linebot.PostbackAction {
	data := &Data{
		Action:    SelectAudioLanguageAction,
		MessageID: mid,
		Language:  lang.Code,
	}
	bytes, _ := json.Marshal(data)
	return linebot.NewPostbackAction(lang.Label, string(bytes), "", lang.Label)
}

// NewSelectAudioSpeechLanguagePostbackAction ...
func NewSelectAudioSpeechLanguagePostbackAction(mid string, code language.SpeechCode, cursor int) *linebot.PostbackAction {
	data := &Data{
		Action:    SelectAudioLanguageAction,
		MessageID: mid,
		Language:  code.Code,
		Cursor:    cursor,
	}
	bytes, _ := json.Marshal(data)
	return linebot.NewPostbackAction(code.Label, string(bytes), "", code.Label)
}
