package quickreply

import (
	"encoding/json"

	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/sh0e1/translation-konjac/pkg/language"
	"github.com/sh0e1/translation-konjac/pkg/line/postback"
)

const maxItems = 13

// SelectLanguageItem ...
func SelectLanguageItem() *linebot.QuickReplyItems {
	buttons := make([]*linebot.QuickReplyButton, 0, len(language.Languages)-1)
	for _, lang := range language.Languages[1:] {
		action := postback.NewSelectLanguagePostbackAction(lang)
		button := linebot.NewQuickReplyButton("", linebot.QuickReplyAction(action))
		buttons = append(buttons, button)
	}
	return linebot.NewQuickReplyItems(buttons...)
}

// SelectAudioLanguageItem ...
func SelectAudioLanguageItem(mid string) *linebot.QuickReplyItems {
	buttons := make([]*linebot.QuickReplyButton, 0, len(language.Languages))
	for _, lang := range language.Languages {
		action := postback.NewSelectAudioLanguagePostbackAction(mid, lang)
		button := linebot.NewQuickReplyButton("", linebot.QuickReplyAction(action))
		buttons = append(buttons, button)
	}
	return linebot.NewQuickReplyItems(buttons...)
}

// SelectAudioSpeechCodeItem ...
func SelectAudioSpeechCodeItem(mid, code string, cursor int) *linebot.QuickReplyItems {
	codes := language.GetSpeechCodes(code)[cursor:]
	if len(codes) > maxItems {
		codes = codes[:maxItems-2]
		cursor = len(code) - 2
	}

	buttons := make([]*linebot.QuickReplyButton, 0, len(codes))
	for _, code := range codes {
		action := postback.NewSelectAudioSpeechLanguagePostbackAction(mid, code, cursor)
		button := linebot.NewQuickReplyButton("", linebot.QuickReplyAction(action))
		buttons = append(buttons, button)
	}
	if cursor != 0 {
		data := &postback.Data{
			Action:    postback.SelectAudioLanguageAction,
			MessageID: mid,
			Language:  code,
			Cursor:    cursor,
		}
		bytes, _ := json.Marshal(data)
		action := linebot.NewPostbackAction("next", string(bytes), "", "next")
		button := linebot.NewQuickReplyButton("", linebot.QuickReplyAction(action))
		buttons = append(buttons, button)
	}

	return linebot.NewQuickReplyItems(buttons...)
}
