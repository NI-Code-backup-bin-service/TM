package routers

import (
	"github.com/gorilla/mux"
	"net/http"
	"nextgen-tms-website/handlers"
)

func CreateAutomationAPIRouter() *mux.Router {
	r := mux.NewRouter()

	handlers.ValidatorInit()

	r.Use(handlers.LoggingMiddleware)
	r.Use(handlers.AuthMiddleware)

	v1 := r.PathPrefix("/api/v1").Subrouter()
	v1.HandleFunc("/tid", handlers.CreateTID).Methods(http.MethodPost)
	v1.HandleFunc("/tid", handlers.UpdateTID).Methods(http.MethodPut)
	v1.HandleFunc("/otp", handlers.GenerateOTP).Methods(http.MethodPost)

	return r
}
