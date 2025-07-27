package openapi

import (
	"embed"
	"net/http"
)

var fs embed.FS

func OpenAPIServeYAML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	data, _ := fs.ReadFile("openapi.yaml")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func RedocHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`<!doctype html>
<html>
  <head>
    <title>API Docs</title>
    <meta charset="utf-8"/>
    <script src="https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js"></script>
    <style>html,body,#redoc{height:100%;margin:0;padding:0}</style>
  </head>
  <body>
    <div id="redoc"></div>
    <script>
      Redoc.init('/openapi.yaml', {}, document.getElementById('redoc'))
    </script>
  </body>
</html>`))
}
