package web

import (
	"encoding/json"
	"net/http"
)

// WriteJSON sets Content-Type to JSON and encodes v with the given status code.
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, err = w.Write(b)
	return err
}
