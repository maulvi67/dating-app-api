package base

import (
	"context"

	"encoding/json"
	"net/http"
)

func EncodeResponseHTTP(ctx context.Context, w http.ResponseWriter, resp interface{}) error {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	result := GetHttpResponse(resp)
	switch result.Meta.Code {
	case 404:
		w.WriteHeader(http.StatusNotFound)
	case 401:
		w.WriteHeader(http.StatusUnauthorized)
	case 400:
		w.WriteHeader(http.StatusBadRequest)
	case 200:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	return json.NewEncoder(w).Encode(resp)
}
