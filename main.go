package main

import (
	"log"
	"net/http"

	router "github.com/FooxyS/FintrackAPI/router"
)

func main() {
	router := router.SetupRouter()

	err := http.ListenAndServe("localhost:8080", router)
	if err != nil {
		log.Fatalf("problem with starting the server: %v\n", err)
	}
}
