package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgraph-io/badger/v4"
	"github.com/go-chi/chi/v5"
	"github.com/zimlewis/shortened/repository"
)

func RedirectShortened(
	eventChannle chan []byte, 
	repo repository.Repository,
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		path := chi.URLParam(request, "id") 

		full, err := repo.GetShortenedResult(ctx, path)
		if err == badger.ErrKeyNotFound {
			writer.Header().Set("Content-Type", "application/json")
			resp := JSON {
				"error": fmt.Sprintf("Cannot find full link for %s", path),
			}

			respBytes, err := json.Marshal(resp)
			if err != nil { 
				fmt.Printf("Cannot marshal the response: %s", err.Error())
				writer.WriteHeader(http.StatusInternalServerError)
				return 
			}

			writer.WriteHeader(http.StatusNotFound)
			writer.Write(respBytes)
			return
		} else if err != nil {
			fmt.Printf("Cannot get shortened url: %s\n", err.Error())

			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusInternalServerError)

			resp := JSON {
				"error": fmt.Sprintf("Cannot find full link for %s", path),
			}

			respBytes, err := json.Marshal(resp)
			if err != nil { 
				fmt.Printf("Cannot marshal the response: %s", err.Error())
				writer.WriteHeader(http.StatusInternalServerError)
				return 
			}

			writer.Write(respBytes)
			return
		}

		eventChannle <- []byte(path) 

		http.Redirect(writer, request, full, http.StatusMovedPermanently)
	}
}

func GetShortenedCount(
	repo repository.Repository,
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		path := chi.URLParam(request, "id") 
		writer.Header().Set("Content-Type", "application/json")

		count, err := repo.GetClickedCount(ctx, path)
		if err != nil {
			fmt.Printf("Cannot get the count of link %s: %s", path, err.Error())
			var rsp = JSON {
				"error": "Cannot get count of the link",
			}

			b, err := json.Marshal(rsp)
			if err == nil {
				writer.Write(b)
			}

			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		rsp := JSON {
			"count": count,
			"link": path,
		}

		b, err := json.Marshal(rsp)
		if err == nil {
			writer.WriteHeader(http.StatusOK)
			writer.Write(b)
			return
		}

		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func AddShortened(
	repo repository.Repository,
) http.HandlerFunc {
	type Body struct {
		Shortened string `json:"shortened"`
		Full      string `json:"full"`
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		var body Body
		ctx := request.Context()
		writer.Header().Set("Content-Type", "application/json")
		

		err := json.NewDecoder(request.Body).Decode(&body)
		if err != nil {
			fmt.Printf("Cannot add to database: %s\n", err.Error())

			var rsp = JSON {
				"error": "Cannot decode the body",
			}
			b, err := json.Marshal(rsp)

			if err == nil {
				writer.Write(b)
			}
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		err = repo.AddShortenedLink(ctx, body.Shortened, body.Full)
		if err != nil {
			fmt.Printf("Cannot add to database: %s\n", err.Error())

			var rsp = JSON {
				"error": "Cannot add the link to database",
			}
			b, err := json.Marshal(rsp)

			if err == nil {
				writer.Write(b)
			}
			writer.WriteHeader(http.StatusInternalServerError)

			return
		}

		writer.WriteHeader(http.StatusCreated)
	}
}
