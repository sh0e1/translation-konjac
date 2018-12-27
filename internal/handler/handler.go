package handler

import (
	"context"

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
}

// Handle ...
func (b *BaseEventHandler) Handle(ctx context.Context) error {
	return nil
}
