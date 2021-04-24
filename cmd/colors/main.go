package main

import (
	"log"
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
)

type envConfig struct {
	Port int    `envconfig:"PORT" default:"8080" required:"true"`
	Sink string `envconfig:"K_SINK"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("failed to process env var: %s", err)
	}

	handler := new(Handler)

	if env.Sink != "" {
		if client, err := cloudevents.NewClientHTTP(cloudevents.WithTarget(env.Sink)); err != nil {
			log.Printf("failed to make cloudevents client: %v\n", err)
		} else {
			handler.client = client
		}
	}

	log.Printf("Server starting on port :%d\n", env.Port)
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("failed to start server, %s", err.Error())
	}
}
