package main

import (
	"fmt"
	"net/http"

	"github.com/dzemildupljak/risc_monolith/server/infrastructure"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Hello from backend")

	logger := infrastructure.NewLogger()

	infrastructure.Load(logger)

	db := infrastructure.SetupDatabaseConnection()

	defer infrastructure.CloseDatabaseConnection(db)

	api := ApiImplementation(*db, logger)

	r := mux.NewRouter()

	// Auth rounting
	authR := r.PathPrefix("/auth").Subrouter()
	authR.HandleFunc("/signup", api.authController.SignUp).Methods("POST")
	authR.HandleFunc("/login", api.authController.Login).Methods("POST")
	authR.Use(api.authController.MiddlewareValidateUser)

	// Refresh token
	refR := r.PathPrefix("/refresh-token").Subrouter()
	refR.HandleFunc("", api.authController.RefreshToken)
	refR.Use(api.authController.MiddlewareValidateRefreshToken)

	// 	// User routing
	usrR := r.PathPrefix("/user").Subrouter()
	usrR.HandleFunc("", api.authController.Index).Methods("GET")
	usrR.Use(api.authController.MiddlewareValidateAccessToken)

	http.ListenAndServe(":8080", r)
}
