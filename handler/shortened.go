package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zimlewis/shortened/repository"
)


func RedirectShortened(
	eventChannle chan []byte, 
	repo repository.Repository,
) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		path := chi.URLParam(request, "id") 

		full, err := repo.GetShortenedResult(ctx, path)
		if err != nil {
			response.Write([]byte("Error"))
			fmt.Printf("Cannot get shortened url: %s\n", err.Error())
			return
		}

		eventChannle <- []byte(path) 
		response.Write([]byte(full))
	}
}

func GetShortenedCount(
	repo repository.Repository,
) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		path := chi.URLParam(request, "id") 

		count, err := repo.GetClickedCount(ctx, path)
		if err != nil {
			response.Write([]byte("Error"))
			fmt.Printf("Cannot get the count of link %s: %s", path, err.Error())
			return
		}

		fmt.Fprintf(response, "%d", count)
	}
}

func AddShortened(
	repo repository.Repository,
) http.HandlerFunc {
	type Body struct {
		Shortened string `json:"shortened"`
		Full      string `json:"full"`
	}
	return func(response http.ResponseWriter, request *http.Request) {
		var body Body
		ctx := request.Context()

		err := json.NewDecoder(request.Body).Decode(&body)
		if err != nil {
			response.Write([]byte("Wrong"))
			fmt.Printf("Cannot get the request body: %s\n", err.Error())
			return
		}

		err = repo.AddShortenedLink(ctx, body.Shortened, body.Full)
		if err != nil {
			response.Write([]byte("Error"))
			fmt.Printf("Cannot add to database: %s\n", err.Error())
			return
		}

		response.Write([]byte("Ok"))
	}
}
