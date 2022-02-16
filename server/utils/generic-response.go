package utils

// This text will appear as description of your response body.
// swagger:response genericResponse
type GenericResponseWrapper struct {
	// in:body
	Body GenericResponse
}
type GenericResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// AuthResponse
// swagger:response authResponse
type AuthResponseWrapper struct {
	// in:body
	Body AuthResponse
}
type AuthResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	Username     string `json:"username"`
}

type TokenResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}
