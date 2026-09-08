package main

import (
	"encoding/json"
	"net/http"
)

func writeJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status) // shouldn't this come after encoding like if any error occured while encoding then?
	return json.NewEncoder(w).Encode(data)
}

func readJson(w http.ResponseWriter, r *http.Request, data any) error {
	max_bytes := 1_048_576                                    // 1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(max_bytes)) // what is this?
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data) // What is this data?
}

func writeErrorJson(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}

	return writeJson(w, status, &envelope{Error: message})
}
