package handler

import (
	"context"
	"fmt"

	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/sh0e1/translation-konjac/internal/message"
	"github.com/sh0e1/translation-konjac/pkg/datastore/resources"
	"github.com/sh0e1/translation-konjac/pkg/language"
	"github.com/sh0e1/translation-konjac/pkg/line/quickreply"
)

// MessageHandler ...
type MessageHandler struct {
	*BaseEventHandler
}

// Handle ...
func (h *MessageHandler) Handle(ctx context.Context) error {
	var fn messageHandleFunc
	switch h.Event.Message.(type) {
	case *linebot.TextMessage:
		fn = h.handleTextMessage
	case *linebot.ImageMessage:
		fn = h.handleImageMessage
	case *linebot.AudioMessage:
		fn = h.handleAudioMessage
	}
	return fn(ctx, h.Event.Message)
}

type messageHandleFunc func(ctx context.Context, message linebot.Message) error

func (h *MessageHandler) handleTextMessage(ctx context.Context, m linebot.Message) error {
	msg := m.(*linebot.TextMessage)
	u := &resources.User{
		ID: h.Event.Source.UserID,
	}
	if err := u.Load(ctx); err != nil {
		return err
	}

	if message.ChangeLanguage.Equal(msg.Text) {
		reply := linebot.NewTextMessage(message.ApplyChangeLanguage.String())
		quickReplies := quickreply.SelectLanguageItem()
		_, err := h.Bot.ReplyMessage(h.Event.ReplyToken, reply.WithQuickReplies(quickReplies)).WithContext(ctx).Do()
		return err
	}

	detection, err := h.Translator.DetectLanguage(ctx, msg.Text)
	if err != nil {
		return err
	}
	var (
		source = detection
		target = u.SelectLanguage
	)
	if source == target {
		source = target
		target = language.JapaneseLanguageCode
	}
	translated, err := h.Translator.Translate(ctx, msg.Text, source, target)
	if err != nil {
		return err
	}

	reply := linebot.NewTextMessage(translated)
	_, err = h.Bot.ReplyMessage(h.Event.ReplyToken, reply).WithContext(ctx).Do()
	return err
}

func (h *MessageHandler) handleImageMessage(ctx context.Context, m linebot.Message) error {
	msg := m.(*linebot.ImageMessage)
	u := &resources.User{
		ID: h.Event.Source.UserID,
	}
	if err := u.Load(ctx); err != nil {
		return err
	}

	messageContent, err := h.Bot.GetMessageContent(msg.ID).Do()
	if err != nil {
		return err
	}
	defer messageContent.Content.Close()

	detection, err := h.ImageAnnotator.DetectTexts(ctx, messageContent.Content)
	if err != nil {
		return err
	}
	var (
		source = detection.Locale
		target = u.SelectLanguage
	)
	if source == target {
		source = target
		target = language.JapaneseLanguageCode
	}
	translated, err := h.Translator.Translate(ctx, detection.Description, source, target)
	if err != nil {
		return err
	}

	reply := linebot.NewTextMessage(translated)
	_, err = h.Bot.ReplyMessage(h.Event.ReplyToken, reply).WithContext(ctx).Do()
	return err
}

func (h *MessageHandler) handleAudioMessage(ctx context.Context, m linebot.Message) error {
	msg := m.(*linebot.AudioMessage)
	u := &resources.User{
		ID: h.Event.Source.UserID,
	}
	if err := u.Load(ctx); err != nil {
		return err
	}

	messageContent, err := h.Bot.GetMessageContent(msg.ID).Do()
	if err != nil {
		return err
	}
	defer messageContent.Content.Close()

	name := fmt.Sprintf("audios/%s/%s.mp4", h.Event.Source.UserID, msg.ID)
	if err := h.Storager.Upload(ctx, name, messageContent.ContentType, messageContent.Content); err != nil {
		return err
	}

	audio := resources.NewAudio(msg.ID, h.Event.Source.UserID, name)
	if err := audio.Save(ctx); err != nil {
		return err
	}

	reply := linebot.NewTextMessage(message.AudioLanguage.String())
	quickReplies := quickreply.SelectAudioLanguageItem(msg.ID)
	_, err = h.Bot.ReplyMessage(h.Event.ReplyToken, reply.WithQuickReplies(quickReplies)).WithContext(ctx).Do()
	return err
}
