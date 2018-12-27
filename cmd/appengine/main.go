package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"google.golang.org/appengine"

	"github.com/sh0e1/translation-konjac/pkg/handler"
	"github.com/sh0e1/translation-konjac/pkg/language"
	"github.com/sh0e1/translation-konjac/pkg/middleware"
	"github.com/sh0e1/translation-konjac/pkg/service/pubsub"
	"github.com/sh0e1/translation-konjac/pkg/service/storage"
	"github.com/sh0e1/translation-konjac/pkg/service/translate"
	"github.com/sh0e1/translation-konjac/pkg/service/vision"
)

const languagesFilePath = "./languages.json"

func main() {
	var (
		channelID     = os.Getenv("CHANNEL_ID")
		channelSecret = os.Getenv("CHANNEL_SECRET")
		pubsubTopic   = os.Getenv("PUBSUB_TOPIC")
	)

	if err := language.LoadLanguages(languagesFilePath); err != nil {
		log.Fatal(err)
	}

	translator, err := translate.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	defer translator.Close()

	imageAnnotator, err := vision.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	defer imageAnnotator.Close()

	storager, err := storage.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	defer storager.Close()

	pubsubClient, err := pubsub.NewClient(appengine.AppID(context.Background()), pubsubTopic)
	if err != nil {
		log.Fatal(err)
	}
	defer pubsubClient.Close()

	handler := &handler.Handler{
		Translator:     translator,
		ImageAnnotator: imageAnnotator,
		Storager:       storager,
		Topic:          pubsubClient.Topic(pubsubTopic),
		ChannelSecret:  channelSecret,
	}

	mw := middleware.Chain(
		middleware.AppEngineContext,
		middleware.Auth(middleware.NewAuthConfigure(channelID, channelSecret)),
	)
	http.HandleFunc("/hook", mw.Then(handler.WebHook))

	appengine.Main()
}
