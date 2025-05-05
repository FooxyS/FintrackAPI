package router

import (
	"net/http"

	handler "github.com/FooxyS/FintrackAPI/handler/auth"
)

func SetupRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/register", handler.RegisterHandler)
	mux.HandleFunc("/auth/login", handler.LoginHandler)
	mux.HandleFunc("/auth/refresh", handler.RefreshHandler)
	mux.HandleFunc("/auth/logout", handler.LogoutHandler)
	mux.HandleFunc("/auth/me", handler.MeHandler)
	return mux
}
