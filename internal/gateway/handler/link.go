package handler

import (
	"net/http"

	"github.com/sxd0/go_url-shortener/internal/gateway/linkclient"
)

type LinkHandler struct {
	client *linkclient.LinkClient
}

func NewLinkHandler(deps Deps) *LinkHandler {
	return &LinkHandler{
		client: deps.LinkClient,
	}
}

func (h *LinkHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (h *LinkHandler) Create() http.HandlerFunc { return h.List() }
func (h *LinkHandler) Get() http.HandlerFunc    { return h.List() }
func (h *LinkHandler) Update() http.HandlerFunc { return h.List() }
func (h *LinkHandler) Delete() http.HandlerFunc { return h.List() }
