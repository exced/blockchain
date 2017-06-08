package api

import (
	"net/http"

	"github.com/exced/simple-blockchain/backend/model"
)

type CRUDHandlerFunc func(w http.ResponseWriter, r *http.Request, s model.Storage) http.Handler

type CRUDHandler interface {
	serveHTTP(w http.ResponseWriter, r *http.Request)
}
