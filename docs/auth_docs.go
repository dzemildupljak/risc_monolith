package docs

// A ValidationError is an error that is used when the required input fails validation.
// swagger:response validationError
type validationError struct {
	// The error message
	// in: body
	Body struct {
		Message   string
		FieldName string
	}
}

// A ValidationError is an error that is used when the required input fails validation.
// swagger:response validationErrors
type validationErrors struct {
	// The error message
	// in: body
	Body []struct {
		Message   string
		FieldName string
	}
}

/////////////////////////////////////////////////
// swagger:parameters Login
type loginRequest struct {
	// in:body
	Body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
}

// swagger:route POST /auth/login Auth Login
// Get acces to other parts of aplication
// Get tokens for next requests
// responses:
//   200: genericResponse
//   401: genericResponse

/////////////////////////////////////////////////
// swagger:parameters Signup
type signupRequest struct {
	// in:body
	Body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
		Name     string `json:"name"`
	}
}

// swagger:route POST /auth/signup Auth Signup
// Create new user with adequate data
// responses:
//   200: genericResponse
//   401: genericResponse

/////////////////////////////////////////////////

// swagger:route GET /get-password-reset-code Auth GeneratePasswordResetCode
// You receive an email with a code to change the password
// responses:
//   200: genericResponse
//   500: genericResponse

/////////////////////////////////////////////////

// swagger:route GET /verify/mail  Auth VerifyMail
// Recieve mail with a:href to verify email addres
// responses:
//   200: genericResponse
//   401: genericResponse

/////////////////////////////////////////////////

// swagger:parameters ResetPassword
type resetPasswordRequest struct {
	// in:body
	Body struct {
		Code                string `json:"code"`
		Old_password        string `json:"old_password"`
		New_password        string `json:"new_password"`
		New_password_second string `json:"new_password_second"`
	}
}

// swagger:route POST /password-reset Auth ResetPassword
// Need to send old pass/code/new_password
// responses:
//   200: genericResponse
//   500: genericResponse

/////////////////////////////////////////////////

// swagger:parameters ForgotPassword
type forgotPasswordRequest struct {
	// in:body
	Body struct {
		Email               string `json:"email"`
		Code                string `json:"code"`
		New_password        string `json:"new_password"`
		New_password_second string `json:"new_password_second"`
	}
}

// swagger:route POST /auth/forgot-password Auth ForgotPassword
// Need to send new password
// responses:
//   200: genericResponse
//   500: genericResponse

/////////////////////////////////////////////////

// swagger:parameters ForgotPasswordCode
type forgotPasswordCodeRequest struct {
	// in:body
	Body struct {
		Email string `json:"email"`
	}
}

// swagger:route POST /auth/forgot-password-code Auth ForgotPasswordCode
// Need to send new password
// responses:
//   200: genericResponse
//   500: genericResponse

// swagger:route GET /oauth/google/login OAuth OauthGoogleLogin
// Need to send request for google account check
// responses:
//   200: genericResponse
//   401: genericResponse


/////////////////////////////////////////////////

// swagger:parameters SetNewPassword
type setPasswordRequest struct {
	// in:body
	Body struct {
		New_password        string `json:"new_password"`
		New_password_second string `json:"new_password_second"`
	}
}

// swagger:route POST /set-password Auth SetNewPassword
// Need to send old pass/code/new_password
// responses:
//   200: genericResponse
//   500: genericResponse
