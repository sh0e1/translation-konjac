package handler

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/sh0e1/translation-konjac/internal/message"
	"github.com/sh0e1/translation-konjac/pkg/datastore/resources"
	"github.com/sh0e1/translation-konjac/pkg/language"
	"github.com/sh0e1/translation-konjac/pkg/line/postback"
	ps "github.com/sh0e1/translation-konjac/pkg/service/pubsub"
)

// PostbackHandler ...
type PostbackHandler struct {
	*BaseEventHandler
}

// Handle ...
func (h *PostbackHandler) Handle(ctx context.Context) error {
	var data postback.Data
	if err := json.Unmarshal([]byte(h.Event.Postback.Data), &data); err != nil {
		return err
	}

	var fn postbackHandleFunc
	switch data.Action {
	case postback.SelectLanguageAction:
		fn = h.handleSelectLanguagePostBack
	case postback.SelectAudioLanguageAction:
		fn = h.handleSelectAudioLanguagePostBack
	}
	return fn(ctx, &data)
}

type postbackHandleFunc func(ctx context.Context, data *postback.Data) error

func (h *PostbackHandler) handleSelectLanguagePostBack(ctx context.Context, data *postback.Data) error {
	u := &resources.User{ID: h.Event.Source.UserID}
	if err := u.Load(ctx); err != nil {
		return err
	}

	u.SelectLanguage = data.Language
	if err := u.Save(ctx); err != nil {
		return err
	}

	reply := linebot.NewTextMessage(message.SelectedLanguage.Format(u.SelectLanguage))
	_, err := h.Bot.ReplyMessage(h.Event.ReplyToken, reply).WithContext(ctx).Do()
	return err
}

func (h *PostbackHandler) handleSelectAudioLanguagePostBack(ctx context.Context, data *postback.Data) error {
	if language.IsMultipleSpeechCode(data.Language) {
	}

	audio := &resources.Audio{ID: data.MessageID}
	if err := audio.Load(ctx); err != nil {
		return err
	}
	audio.SourceLanguage = language.GetSpeechCode(data.Language)
	if err := audio.Save(ctx); err != nil {
		return err
	}

	msgData := &ps.Data{
		MessageID:      data.MessageID,
		ReplayToken:    h.Event.ReplyToken,
		AudioFilePath:  audio.Path,
		SourceLanguage: audio.SourceLanguage,
	}
	msg := &pubsub.Message{
		Data: msgData.Marshal(),
	}
	if _, err := h.Topic.Publish(ctx, msg).Get(ctx); err != nil {
		return err
	}

	return nil
}
