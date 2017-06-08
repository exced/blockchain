package api

import (
	"net/http"

	"github.com/exced/simple-blockchain/backend/model"
)

func WithStorage(h http.Handler, s model.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		h.ServeHTTP(w, r, s)

	})
}
