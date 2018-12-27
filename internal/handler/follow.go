package handler

import (
	"context"

	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/sh0e1/translation-konjac/internal/message"
	"github.com/sh0e1/translation-konjac/pkg/datastore/resources"
	"github.com/sh0e1/translation-konjac/pkg/language"
	"github.com/sh0e1/translation-konjac/pkg/line/quickreply"
)

// FollowHandler ...
type FollowHandler struct {
	*BaseEventHandler
}

// Handle ...
func (h *FollowHandler) Handle(ctx context.Context) error {
	u := resources.NewUser(h.Event.Source.UserID, language.DefaultLanguageCode)
	if err := u.Save(ctx); err != nil {
		return err
	}
	reply := linebot.NewTextMessage(message.Follow.String())
	quickReplies := quickreply.SelectLanguageItem()
	_, err := h.Bot.ReplyMessage(h.Event.ReplyToken, reply.WithQuickReplies(quickReplies)).WithContext(ctx).Do()
	return err
}
