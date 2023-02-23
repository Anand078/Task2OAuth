package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

/*
HandleMain Function renders the index page when the application index route is called
*/
func Main(w http.ResponseWriter, r *http.Request) {
	accessToken, err := r.Cookie("access_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Use the tokens to obtain an authenticated http.Client
	token := oauth2.Token{
		AccessToken:  accessToken.Value,
		RefreshToken: refreshToken.Value,
	}
	client := oauthConfig.Client(context.Background(), &token)

	// Fetch the user's profile information from the Google API
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user info: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Parse the user's profile information from the response body
	var profileInfo struct {
		ID    string `json:"sub"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &profileInfo); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse user info: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Hello, %s (%s)!", profileInfo.Name, profileInfo.Email)
}

/*
HandleLogin Function
*/
func HandleLogin(w http.ResponseWriter, r *http.Request, oauthConf *oauth2.Config, oauthStateString string) {
	URL, err := url.Parse(oauthConf.Endpoint.AuthURL)
	if err != nil {
		fmt.Println("Parse: " + err.Error())
	}
	parameters := url.Values{}
	parameters.Add("client_id", oauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)
	URL.RawQuery = parameters.Encode()
	url := URL.String()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
