module github.com/dzemildupljak/risc_monolith

// +heroku goVersion go1.17
go 1.17

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/lib/pq v1.10.4
	golang.org/x/crypto v0.0.0-20220112180741-5e0467b6c7ce
)

require github.com/felixge/httpsnoop v1.0.1 // indirect
