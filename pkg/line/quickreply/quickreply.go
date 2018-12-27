package quickreply

import (
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/sh0e1/translation-konjac/pkg/language"
	"github.com/sh0e1/translation-konjac/pkg/line/postback"
)

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
		action := postback.NewSelectAudioLanguagePostbackAction(lang, mid)
		button := linebot.NewQuickReplyButton("", linebot.QuickReplyAction(action))
		buttons = append(buttons, button)
	}
	return linebot.NewQuickReplyItems(buttons...)
}
