package main

import (
	"context"
	"html/template"
	"log"
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

var index = `<!DOCTYPE html>
<html>
  <head>
    <style>
      body {
        background-color: rgb({{ .red }}, {{ .green }}, {{ .blue }});
      }
    </style>
  </head>
  <body></body>
</html>`

var indexTemplate *template.Template

func init() {
	indexTemplate = template.Must(template.New("index").Parse(index))
}

type Handler struct {
	client cloudevents.Client
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	color := map[string]int{
		"red":   0,
		"green": 0,
		"blue":  255,
	}
	_ = indexTemplate.Execute(w, color)

	if h.client != nil {
		if result := h.client.Send(context.Background(), newEvent(color)); cloudevents.IsUndelivered(result) {
			log.Printf("failed to send cloudevent: %v\n", result.Error())
		}
	}
}

func newEvent(data interface{}) cloudevents.Event {
	event := cloudevents.NewEvent() // Sets version
	event.SetType("com.n3wscott.atlanta.colors")
	event.SetSource("github.com/n3wscott/k8s-meetup-atlanta/cmd/colors")
	if err := event.SetData(cloudevents.ApplicationJSON, data); err != nil {
		log.Printf("failed to cloudevents event: %v\n", err)
	}
	return event
}
