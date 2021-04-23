package main

import (
	"html/template"
	"log"
	"net/http"
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

func handler(w http.ResponseWriter, r *http.Request) {
	_ = indexTemplate.Execute(w, map[string]string{
		"red":   "0",
		"green": "0",
		"blue":  "255",
	})
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
