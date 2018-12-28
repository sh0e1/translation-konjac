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

func main() {
	var (
		channelID          = os.Getenv("CHANNEL_ID")
		channelSecret      = os.Getenv("CHANNEL_SECRET")
		pubsubTopic        = os.Getenv("PUBSUB_TOPIC")
		googleCloudProject = os.Getenv("GOOGLE_CLOUD_PROJECT")
		languagesFilePath  = os.Getenv("LANGUAGES_FILE_PATH")
	)

	if err := language.LoadLanguages(languagesFilePath); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	translator, err := translate.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer translator.Close()

	imageAnnotator, err := vision.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer imageAnnotator.Close()

	storager, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer storager.Close()

	pubsubClient, err := pubsub.NewClient(ctx, googleCloudProject)
	if err != nil {
		log.Fatal(err)
	}
	defer pubsubClient.Close()

	topic, err := pubsubClient.CreateTopicIfNotExists(ctx, pubsubTopic)
	if err != nil {
		log.Fatal(err)
	}

	handler := &handler.Handler{
		Translator:     translator,
		ImageAnnotator: imageAnnotator,
		Storager:       storager,
		Topic:          topic,
		ChannelSecret:  channelSecret,
	}

	mw := middleware.Chain(
		middleware.AppEngineContext,
		middleware.Auth(middleware.NewAuthConfigure(channelID, channelSecret)),
	)
	http.HandleFunc("/hook", mw.Then(handler.WebHook))

	appengine.Main()
}
