package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dzemildupljak/risc_monolith/server/infrastructure"
	"github.com/gorilla/handlers"
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

	app_port := ":" + os.Getenv("PORT")

	r := mux.NewRouter()

	// Auth rounting
	authR := r.PathPrefix("/auth").Subrouter()
	authR.HandleFunc("/signup", api.authController.SignUp).Methods("POST")
	authR.HandleFunc("/login", api.authController.Login).Methods("POST")
	authR.Use(api.authController.MiddlewareValidateUser)

	verfyR := r.PathPrefix("/verify").Subrouter()
	verfyR.HandleFunc("/mail", api.authController.VerifyMail).Methods("GET")
	verfyR.Use(api.authController.MiddlewareValidateVerificationData)

	// Refresh token
	refR := r.PathPrefix("/refresh-token").Subrouter()
	refR.HandleFunc("", api.authController.RefreshToken)
	refR.Use(api.authController.MiddlewareValidateRefreshToken)

	// Reset password
	getR := r.PathPrefix("").Subrouter()
	getR.HandleFunc("/get-password-reset-code", api.authController.GeneratePassResetCode)
	getR.HandleFunc("/password-reset", api.authController.PasswordResetCode)
	getR.Use(api.authController.MiddlewareValidateAccessToken)

	// 	// User routing
	usrR := r.PathPrefix("/user").Subrouter()
	usrR.HandleFunc("", api.authController.Index).Methods("GET")
	usrR.Use(api.authController.MiddlewareValidateAccessToken)

	headers := handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE", "PUT"})
	origins := handlers.AllowedOrigins([]string{"*"})

	http.ListenAndServe(app_port, handlers.CORS(headers, methods, origins)(r))
}
