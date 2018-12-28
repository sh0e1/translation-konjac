package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
	ps "github.com/sh0e1/translation-konjac/pkg/service/pubsub"
)

func main() {
	var (
		googleCloudProject      = os.Getenv("GOOGLE_CLOUD_PROJECT")
		pubsubTopic             = os.Getenv("PUBSUB_TOPIC")
		pubsubSubscription      = os.Getenv("PUBSUB_SUBSCRIPTION")
		pubsubVerificationToken = os.Getenv("PUBSUB_VERIFICATION_TOKEN")
	)

	ctx := context.Background()

	pubsubClient, err := ps.NewClient(ctx, googleCloudProject)
	if err != nil {
		log.Fatal(err)
	}
	defer pubsubClient.Close()

	topic, err := pubsubClient.CreateTopicIfNotExists(ctx, pubsubTopic)
	if err != nil {
		log.Fatal(err)
	}

	endpoint := fmt.Sprintf("https://%s.appspot.com/audios/subscribe?token=%s",
		googleCloudProject, pubsubVerificationToken)
	cfg := pubsub.SubscriptionConfig{
		Topic: topic,
		PushConfig: pubsub.PushConfig{
			Endpoint: endpoint,
		},
	}
	if _, err := pubsubClient.CreateSubscriptionIfNotExists(ctx, pubsubSubscription, cfg); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/audios/subscribe", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello Audio Converter")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
