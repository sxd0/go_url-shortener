package openapi

import (
  "embed"
  "net/http"

  "github.com/go-chi/chi/v5"
)

//go:embed openapi.yaml openapi.html
var contentFS embed.FS

func Mount(r chi.Router) {
  r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    data, err := contentFS.ReadFile("openapi.yaml")
    if err != nil {
      http.Error(w, "spec not found", http.StatusInternalServerError)
      return
    }
    w.Write(data)
  })

  r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    data, err := contentFS.ReadFile("openapi.html")
    if err != nil {
      http.Error(w, "docs not found", http.StatusInternalServerError)
      return
    }
    w.Write(data)
  })
}
