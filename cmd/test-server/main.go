package main

import (
	"net/http"

	"github.com/calebamiles/example-fortune-service/service"
)

func main() {
	// Don't use Cadence backend
	http.HandleFunc("/fortune", service.HandleGetFortuneDirect)
	http.HandleFunc("/healthz", service.HandleGetHealthz)

	http.ListenAndServe(":8080", nil)
}
