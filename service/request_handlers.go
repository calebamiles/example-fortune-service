package service

import (
	"encoding/json"
	"log"

	"github.com/calebamiles/example-fortune-service/provider"

	"net/http"
)

const defaultFortune = `
Software engineering is what happens to programming when you add time and other programmers.

–Russ Cox
`

// HandleGetFortune returns a new fortune
func HandleGetFortune(w http.ResponseWriter, req *http.Request) {
	defaultFortune := []byte(defaultFortune)

	fortune := provider.NewFortune(defaultFortune)
	rawTxt, err := fortune.Get()
	if err != nil {
		log.Printf("error: getting fortune: %s", err)
	}

	n, err := w.Write(rawTxt)
	if err != nil {
		log.Printf("error: writing response: %s", err)
	}

	if n != len(rawTxt) {
		log.Printf("error: expected to write %d bytes, but only wrote %d", len(rawTxt), n)
	}
}

// HandleGetHealthz returns server health as ok
func HandleGetHealthz(w http.ResponseWriter, req *http.Request) {
	var response struct {
		Status string `json:"status"`
	}

	response.Status = "ok"
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("error: encoding status to JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	n, err := w.Write(responseJSON)
	if err != nil {
		log.Printf("error: writing response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if n != len(responseJSON) {
		log.Printf("error: expected to write %d bytes, but only wrote %d", len(responseJSON), n)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
