package auth_rest

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dzemildupljak/risc_monolith/server/domain"
	"github.com/dzemildupljak/risc_monolith/server/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUser struct {
	Id             string `json:"id"`
	Email          string `json:"email"`
	Verified_email bool   `json:"verified_email"`
	Picture        string `json:"picture"`
	Hd             string `json:"hd"`
}

// Scopes: OAuth 2.0 scopes provide
// a way to limit the amount of access that is granted to an access token.
var googleOauthConfig = &oauth2.Config{
	RedirectURL:  fmt.Sprintf("%s/v1/oauth/google/callback", os.Getenv("HOST_ADRESS")),
	ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func (ac *AuthController) OauthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Create oauthState cookie
	oauthState := generateStateOauthCookie(w)

	googleOauthConfig.RedirectURL = fmt.Sprintf("%s/v1/oauth/google/callback", os.Getenv("HOST_ADRESS"))
	googleOauthConfig.ClientID = os.Getenv("OAUTH_CLIENT_ID")
	googleOauthConfig.ClientSecret = os.Getenv("OAUTH_CLIENT_SECRET")

	/*
		AuthCodeURL receive state that is a token to protect the user from CSRF attacks.
		You must always provide a non-empty string and
		validate that it matches the the state query parameter on your redirect callback.
	*/
	u := googleOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(2 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

// OauthGoogleCallback
func (ac *AuthController) OauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Read oauthState from Cookie
	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data, err := getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// GetOrCreate User in your db.
	// Redirect or response with a token.
	// More code .....

	var obj GoogleUser
	if err := json.Unmarshal(data, &obj); err != nil {
		panic(err)
	}

	usr, err := ac.ai.UserByEmail(r.Context(), obj.Email)

	if err != nil {
		tokenhash := utils.GenerateRandomString(15)

		usr, err = ac.ai.RegisterOauthUser(context.Background(),
			domain.CreateOauthUserParams{
				Email:      obj.Email,
				Tokenhash:  tokenhash,
				Isverified: obj.Verified_email,
				OauthID:    []string{obj.Id},
			})

		if err != nil {
			ac.logger.LogError("err", err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(
				&utils.GenericResponse{
					Status:  false,
					Message: "Unable to create user.Please try again later",
				})
			return
		}
	}

	if !usr.Isverified {
		ac.logger.LogError("unverified user")

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Please verify your mail address before login",
			})

		return
	}

	accessToken, err := ac.ai.GenerateAccessToken(&usr)
	if err != nil {
		ac.logger.LogError("unable to generate access token", "error", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to login the user. Please try again later",
			})

		return
	}

	refreshToken, err := ac.ai.GenerateRefreshToken(&usr)
	if err != nil {
		ac.logger.LogError("unable to generate refresh token", "error", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to login the user. Please try again later",
			})

		return
	}

	cookie := &http.Cookie{
		Name:  "access_token",
		Value: accessToken,
	}

	http.SetCookie(w, cookie)

	http.Redirect(w, r,
		"http://localhost:3000/home?token="+refreshToken,
		http.StatusTemporaryRedirect)
}

func getUserDataFromGoogle(code string) ([]byte, error) {
	// Use code to get token and get user info from Google.

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}
