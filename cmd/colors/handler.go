package main

import (
	"html/template"
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

type Handler struct {
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet || r.URL.Path != "/" {
		return
	}

	color := map[string]int{
		"red":   255,
		"green": 0,
		"blue":  0,
	}
	_ = indexTemplate.Execute(w, color)
}
