package handlers

import (
	"net/http"
)

func (h *Handler) Test(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello world!"))
}
