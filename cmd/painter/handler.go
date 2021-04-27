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
        background-image: url("canvas.png");
		background-size: cover;
      }
    </style>
  </head>
  <body></body>
</html>`

var indexTemplate *template.Template

func init() {
	indexTemplate = template.Must(template.New("index").Parse(index))
}

func displayHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		_ = indexTemplate.Execute(w, nil)
	} else {
		w.Header().Set("Content-Type", "image/png") // <-- set the content-type header
		w.Write(lastImage)
	}
}

/*

for i in {1..2500}; do
  curl http://colors.default.d2k.n3wscott.com &
done

*/
