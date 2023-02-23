package service

import (
	"context"
	"fmt"
	"googleauth/helper"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauthConfig = &oauth2.Config{
		ClientID:     helper.OAUTH_CLIENT_ID,
		ClientSecret: helper.OAUTH_CLIENT_SECRET,
		RedirectURL:  helper.OAUTH_REDIRECT_URL,
		Scopes:       helper.SCOPES,
		Endpoint:     google.Endpoint,
	}
	oauthState = helper.OAUTH_STATE
)

/*
HandleGoogleLogin Function
*/
func Login(w http.ResponseWriter, r *http.Request) {
	HandleLogin(w, r, oauthConfig, oauthState)
}

/*
CallBackFromGoogle Function
*/
func CallBack(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	// Exchange the authorization code for an access token and refresh token
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to exchange token: %v", err), http.StatusBadRequest)
		return
	}

	// Set the access token as a cookie
	accessToken := &http.Cookie{
		Name:     "access_token",
		Value:    token.AccessToken,
		HttpOnly: true,
	}
	if r.FormValue("remember_me") == "on" {
		accessToken.Expires = time.Now().AddDate(1, 0, 0) // Expires in 1 year
	}
	http.SetCookie(w, accessToken)

	// Set the refresh token as a cookie
	refreshToken := &http.Cookie{
		Name:     "refresh_token",
		Value:    token.RefreshToken,
		HttpOnly: true,
	}
	if r.FormValue("remember_me") == "on" {
		fmt.Println("inside remember me functionality")
		refreshToken.Expires = time.Now().AddDate(1, 0, 0) // Expires in 1 year
	}
	http.SetCookie(w, refreshToken)

	// Redirect the user to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)

	OAuthGmailService()
	status, err := SendEmailOAUTH2(helper.EMAIL_TO)
	if err != nil {
		log.Println(err)
	}
	if status {
		log.Println("Email sent successfully using oauth and gmail api")
	}
}
