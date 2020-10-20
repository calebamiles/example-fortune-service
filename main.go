package main

import (
	"net/http"

	"github.com/calebamiles/example-fortune-service/service"
)

func main() {
	http.HandleFunc("/fortune", service.HandleGetFortune)
	http.HandleFunc("/healthz", service.HandleGetHealthz)

	http.ListenAndServe(":8090", nil)
}
