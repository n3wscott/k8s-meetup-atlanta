package main

import (
	"context"
	"log"
	"math/rand"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
)

type envConfig struct {
	Port int `envconfig:"PORT" default:"8080" required:"true"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("failed to process env var: %s", err)
	}

	client, err := cloudevents.NewClientHTTP(cloudevents.WithPort(env.Port))
	if err != nil {
		log.Printf("failed to make cloudevents client: %v\n", err)
	}

	log.Printf("Server starting on port :%d\n", env.Port)
	if err := client.StartReceiver(context.Background(), mix); err != nil {
		log.Fatalf("failed to start receiver, %s", err.Error())
	}
}

type color struct {
	Red   int `json:"red"`
	Green int `json:"green"`
	Blue  int `json:"blue"`
}

func mix(event cloudevents.Event) *cloudevents.Event {
	if event.Type() != "com.n3wscott.atlanta.colors" {
		return nil
	}

	c := new(color)
	if err := event.DataAs(c); err != nil {
		log.Println("failed to get colors from event,", err)
		return nil
	}

	c.Red = rand.Intn(c.Red+1) + rand.Intn(64)
	if c.Red > 255 {
		c.Red = 255
	}

	c.Green = rand.Intn(c.Green+1) + rand.Intn(64)
	if c.Green > 255 {
		c.Green = 255
	}

	c.Blue = rand.Intn(c.Blue+1) + rand.Intn(64)
	if c.Blue > 255 {
		c.Blue = 255
	}

	return newEvent(c)
}

func newEvent(data interface{}) *cloudevents.Event {
	event := cloudevents.NewEvent() // Sets version
	event.SetType("com.n3wscott.atlanta.mixed-colors")
	event.SetSource("github.com/n3wscott/k8s-meetup-atlanta/cmd/mixer")
	if err := event.SetData(cloudevents.ApplicationJSON, data); err != nil {
		log.Printf("failed to cloudevents event: %v\n", err)
	}
	return &event
}
