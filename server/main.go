package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/dzemildupljak/risc_monolith/docs"
	"github.com/dzemildupljak/risc_monolith/server/infrastructure"
	"github.com/go-openapi/runtime/middleware"
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
	r.Handle("/swagger.json", http.FileServer(http.Dir("./")))

	// opts := middleware.SwaggerUIOpts{SpecURL: "swagger.json"}
	// sh := middleware.SwaggerUI(opts, nil)
	// r.Handle("/docs", sh)

	opts1 := middleware.RedocOpts{SpecURL: "/swagger.json"}
	sh1 := middleware.Redoc(opts1, nil)
	r.Handle("/docs", sh1)

	rv1 := r.PathPrefix("/v1").Subrouter()

	// Auth rounting
	authR := rv1.PathPrefix("/auth").Subrouter()
	authR.HandleFunc("/signup", api.authController.SignUp).Methods("POST")
	authR.HandleFunc("/login", api.authController.Login).Methods("POST") // VALIDATE
	authR.HandleFunc("/forgot-password-code", api.authController.ForgotPasswordCode).Methods("POST")
	authR.HandleFunc("/forgot-password", api.authController.ForgotPassword).Methods("POST")
	// authR.Use(api.authController.MiddlewareValidateUser)

	verfyR := rv1.PathPrefix("/verify").Subrouter()
	verfyR.HandleFunc("/mail", api.authController.VerifyMail).Methods("GET")
	verfyR.Use(api.authController.MiddlewareValidateMailVerificationData)

	// Refresh token
	refR := rv1.PathPrefix("/refresh-token").Subrouter()
	refR.HandleFunc("", api.authController.RefreshToken).Methods("GET")
	refR.Use(api.authController.MiddlewareValidateRefreshToken)

	// Reset password
	getR := rv1.PathPrefix("").Subrouter()
	getR.HandleFunc("/get-password-reset-code", api.authController.GeneratePassResetCode).Methods("GET")
	getR.HandleFunc("/password-reset", api.authController.PasswordResetCode).Methods("POST")
	getR.Use(api.authController.MiddlewareValidateAccessToken)

	// 	// User routing
	usrR := rv1.PathPrefix("/user").Subrouter()
	usrR.HandleFunc("", api.authController.Index).Methods("GET")
	usrR.Use(api.authController.MiddlewareValidateAccessToken)

	headers := handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE", "PUT"})
	origins := handlers.AllowedOrigins([]string{"*"})

	http.ListenAndServe(app_port, handlers.CORS(headers, methods, origins)(r))
}
