package handler

import (
	"context"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"

	"github.com/sh0e1/translation-konjac/internal/handler"
	"github.com/sh0e1/translation-konjac/pkg/middleware"
	"github.com/sh0e1/translation-konjac/pkg/service/storage"
	"github.com/sh0e1/translation-konjac/pkg/service/translate"
	"github.com/sh0e1/translation-konjac/pkg/service/vision"
)

// Handler ...
type Handler struct {
	Translator     translate.Translator
	ImageAnnotator vision.ImageAnnotator
	Storager       storage.Storager
	ChannelSecret  string
	HandleEvent    HandleEventFunc
}

// WebHook ...
func (h *Handler) WebHook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	token := middleware.GetChannelTokenFromContext(ctx)

	client := urlfetch.Client(ctx)
	bot, err := linebot.New(h.ChannelSecret, token, linebot.WithHTTPClient(client))
	if err != nil {
		log.Errorf(ctx, "%+v", err)
		return
	}

	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(http.StatusBadRequest)
		}
		log.Errorf(ctx, "%+v", err)
		return
	}

	if h.HandleEvent == nil {
		h.HandleEvent = h.defaultHandleEventFunc
	}

	for _, e := range events {
		eh := h.HandleEvent(bot, e)
		if err := eh.Handle(ctx); err != nil {
			log.Errorf(ctx, "%+v", err)
		}
	}
}

// HandleEventFunc ...
type HandleEventFunc func(bot *linebot.Client, e *linebot.Event) EventHandler

// EventHandler ...
type EventHandler interface {
	Handle(ctx context.Context) error
}

func (h *Handler) defaultHandleEventFunc(bot *linebot.Client, e *linebot.Event) (eh EventHandler) {
	base := &handler.BaseEventHandler{
		Bot:            bot,
		Event:          e,
		Translator:     h.Translator,
		ImageAnnotator: h.ImageAnnotator,
		Storager:       h.Storager,
	}

	switch e.Type {
	case linebot.EventTypeFollow:
		eh = &handler.FollowHandler{BaseEventHandler: base}
	case linebot.EventTypeUnfollow:
		eh = &handler.UnfollowHandler{BaseEventHandler: base}
	case linebot.EventTypeMessage:
		eh = &handler.MessageHandler{BaseEventHandler: base}
	case linebot.EventTypePostback:
		eh = &handler.PostbackHandler{BaseEventHandler: base}
	}
	return
}
