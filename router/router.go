package router

import (
	"net/http"

	handler "github.com/FooxyS/FintrackAPI/handler"
)

func SetupRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/register", handler.RegisterHandler)
	return mux
}
