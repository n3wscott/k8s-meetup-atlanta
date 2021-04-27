package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"sync"

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

	client, err := cloudevents.NewClientHTTP(cloudevents.WithPort(env.Port), cloudevents.WithGetHandlerFunc(displayHandler))
	if err != nil {
		log.Printf("failed to make cloudevents client: %v\n", err)
	}

	log.Printf("Server starting on port :%d\n", env.Port)
	if err := client.StartReceiver(context.Background(), receive); err != nil {
		log.Fatalf("failed to start receiver, %s", err.Error())
	}
}

type mixedColor struct {
	Red   uint8 `json:"red"`
	Green uint8 `json:"green"`
	Blue  uint8 `json:"blue"`
}

func receive(event cloudevents.Event) *cloudevents.Event {
	if event.Type() != "com.n3wscott.atlanta.mixed-colors" {
		return nil
	}

	c := new(mixedColor)
	if err := event.DataAs(c); err != nil {
		log.Println("failed to get colors from event,", err)
		return nil
	}

	return paint(c)
}

var mux = sync.Mutex{}

const (
	width  = 50
	height = 50
)

var canvas *image.RGBA
var x, y int
var lastImage []byte

func paint(c *mixedColor) *cloudevents.Event {
	// Thread safe.
	mux.Lock()
	defer mux.Unlock()

	log.Printf("painting #%02x%02x%02x at [%02d,%02d]\n", c.Red, c.Green, c.Blue, x+1, y+1)

	if canvas == nil {
		canvas = image.NewRGBA(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: width, Y: height}},
		)
		x = 0
		y = 0
	}

	paint := color.RGBA{
		R: c.Red,
		G: c.Green,
		B: c.Blue,
		A: 0xff,
	}

	// Paint!
	canvas.Set(x, y, paint)

	x++
	if x >= width {
		x = 0
		y++
	}
	if y >= height {
		log.Println("#############################")
		log.Println("#                           #")
		log.Printf("# painted a picture (%dx%d) #\n", width, height)
		log.Println("#                           #")
		log.Println("#############################")

		// On the way out, reset.
		defer func() {
			x = 0
			y = 0
			canvas = nil
		}()

		var imageBuf bytes.Buffer
		if err := png.Encode(&imageBuf, canvas); err != nil {
			log.Printf("failed to encode image: %v\n", err)
			return nil
		}

		// Save the last image.
		lastImage = imageBuf.Bytes()

		// Return CloudEvent.
		event := cloudevents.NewEvent() // Sets version
		event.SetType("com.n3wscott.atlanta.painting")
		event.SetSource("github.com/n3wscott/k8s-meetup-atlanta/cmd/painter")
		if err := event.SetData(cloudevents.ApplicationJSON, map[string]string{"message": fmt.Sprintf("painted a picture (%dx%d)", width, height)}); err != nil {
			log.Printf("failed to cloudevents event: %v\n", err)
		}
		return &event
	}
	return nil
}
