package handler

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/sh0e1/translation-konjac/pkg/service/storage"
	"github.com/sh0e1/translation-konjac/pkg/service/translate"
	"github.com/sh0e1/translation-konjac/pkg/service/vision"
)

// BaseEventHandler ...
type BaseEventHandler struct {
	Bot            *linebot.Client
	Event          *linebot.Event
	Translator     translate.Translator
	ImageAnnotator vision.ImageAnnotator
	Storager       storage.Storager
	Topic          *pubsub.Topic
}

// Handle ...
func (b *BaseEventHandler) Handle(ctx context.Context) error {
	return nil
}
