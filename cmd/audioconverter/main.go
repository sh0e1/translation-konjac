package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
	"google.golang.org/appengine"
)

func main() {
	var (
		project = os.Getenv("GOOGLE_CLOUD_PROJECT")
		token   = os.Getenv("PUBSUB_VERIFICATION_TOKEN")
	)

	client, err := pubsub.NewClient(context.Background(), project)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	sub := client.Subscription();
	endpoint := fmt.Sprintf(
		"https://%s/subscribe?token=%s",
		appengine.DefaultVersionHostname(context.Background()),
		token,
	)
	cfg := pubsub.SubscriptionConfig{
		PushConfig: pubsub.PushConfig{
			Endpoint: endpoint,
		},
	}
	if _, err := client.CreateSubscription(context.Background(), project, cfg)

	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello Audio Converter")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
