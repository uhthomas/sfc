package sfc

import (
	"encoding/json"
	"net/http"
)

// Handler handles HTTP requests.
type Handler struct {
	Client     *Client
	FileServer http.Handler
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	allow := "GET, HEAD, OPTIONS"
	switch r.Method {
	case http.MethodHead, http.MethodGet:
	case http.MethodOptions:
		w.Header().Set("Access-Control-Allow-Methods", allow)
		return
	default:
		w.Header().Set("Allow", allow)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if q := r.URL.Query().Get("q"); q != "" {
		h.HandleTrack(w, r, q)
		return
	}
	h.FileServer.ServeHTTP(w, r)
}

func (h Handler) HandleTrack(w http.ResponseWriter, r *http.Request, q string) {
	res, err := h.Client.Track(r.Context(), q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	e.Encode(struct {
		Data interface{}     `json:"data"`
		Raw  json.RawMessage `json:"raw"`
	}{res, res.Body()})
}
