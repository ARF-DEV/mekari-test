package healthhandlr

import "net/http"

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (handler *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Hello World!"))
}
