package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RedirectShortened(eventChannle chan []byte) func (response http.ResponseWriter, request *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		path := chi.URLParam(request, "id") 
		eventChannle <- fmt.Appendf(nil, "Someone clicked %s", path)
		response.Write([]byte("Ok"))
	}
}
