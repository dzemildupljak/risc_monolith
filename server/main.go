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
	// /v1/auth/signup
	authR.HandleFunc("/signup", api.authController.SignUp).Methods("POST")
	// /v1/auth/login
	authR.HandleFunc("/login", api.authController.Login).Methods("POST")
	// /v1/auth/forgot-password-code
	authR.HandleFunc("/forgot-password-code", api.authController.ForgotPasswordCode).Methods("POST")
	// /v1/auth/forgot-password
	authR.HandleFunc("/forgot-password", api.authController.ForgotPassword).Methods("POST")

	// Google 0Auth2
	googleR := rv1.PathPrefix("/oauth").Subrouter()
	// /v1/oauth/google/login
	googleR.HandleFunc("/google/login", api.authController.OauthGoogleLogin).Methods("GET")
	// /v1/oauth/google/callback
	googleR.HandleFunc("/google/callback", api.authController.OauthGoogleCallback)

	// Varify mail
	verfyR := rv1.PathPrefix("/verify").Subrouter()
	// /v1/verify/mail
	verfyR.HandleFunc("/mail", api.authController.VerifyMail).Methods("GET")
	verfyR.Use(api.authController.MiddlewareValidateMailVerificationData)

	// Refresh token
	refR := rv1.PathPrefix("/refresh-token").Subrouter()
	// /v1/refresh-token
	refR.HandleFunc("", api.authController.RefreshToken).Methods("GET")
	refR.Use(api.authController.MiddlewareValidateRefreshToken)

	// Reset password
	getR := rv1.PathPrefix("").Subrouter()
	// /v1/get-password-reset-code
	getR.HandleFunc("/get-password-reset-code", api.authController.GeneratePassResetCode).Methods("GET")
	// /v1/password-reset
	getR.HandleFunc("/password-reset", api.authController.PasswordResetCode).Methods("POST")
	// /v1/set-password
	getR.HandleFunc("/set-password", api.authController.SetNewPassword).Methods("POST")
	getR.Use(api.authController.MiddlewareValidateAccessToken)

	// 	// User routing
	usrR := rv1.PathPrefix("/user").Subrouter()
	// /v1/user/
	usrR.HandleFunc("/users", api.userController.ListUsers).Methods("GET")
	// /v1/user/current
	usrR.HandleFunc("/current", api.userController.CurrentUser).Methods("GET")
	// /v1/user/{user_id}
	usrR.HandleFunc("/{user_id}", api.userController.UserById).Methods("GET")
	// /v1/user/{user_id}
	usrR.HandleFunc("/{user_id}", api.userController.UpdateUserById).Methods("PUT")
	// // /v1/user/{user_id}
	// usrR.HandleFunc("/{user_email}", api.userController.UserById).Methods("GET")
	usrR.Use(api.authController.MiddlewareValidateAccessToken)

	// Donor Routing
	donorR := rv1.PathPrefix("/donor").Subrouter()
	donorR.HandleFunc("/donors", api.donorController.ListDonors).Methods("GET")
	donorR.Use(api.authController.MiddlewareValidateAccessToken)

	// Admin routing
	adminR := rv1.PathPrefix("/admin").Subrouter()
	// /v1/user/
	adminR.HandleFunc("/users", api.authController.Index).Methods("GET")
	adminR.Use(api.authController.MiddlewareValidateAdminAccessToken)

	headers := handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE", "PUT"})
	origins := handlers.AllowedOrigins([]string{"*"})

	http.ListenAndServe(app_port, handlers.CORS(headers, methods, origins)(r))
}
