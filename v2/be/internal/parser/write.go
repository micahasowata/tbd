package parser

import (
	"encoding/json"
	"net/http"
)

type Envelope map[string]any

func Write(w http.ResponseWriter, status int, v Envelope) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(data)
	if err != nil {
		return err
	}

	return nil
}
