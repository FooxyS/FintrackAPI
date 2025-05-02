package router

import (
	"net/http"

	handler "github.com/FooxyS/FintrackAPI/handler"
)

func SetupRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/register", handler.RegisterHandler)
	mux.HandleFunc("/login", handler.UserLogin)
	return mux
}
