package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// GmailService : Gmail client for sending email
var GmailService *gmail.Service

func OAuthGmailService() {

	token := oauth2.Token{
		AccessToken:  "ya29.a0AVvZVsr_fT7VZzkxkMcyx8H4iluxBi4YVHTAZtpisdepH3Y2QjBpZ8OHlRg-wZ54XDviYsMJDi6rYREQbjJ9MTAIW193bwPODhOu0U4mfP251mnIWa6__Q5O_0NapR8r7RYrWMrle9oCpICm04Os2jrEaaGUaCgYKAdcSARESFQGbdwaIM0awz8qZAovOYURBRK8AMg0163",
		RefreshToken: "1//04PB6ZpGaqfEACgYIARAAGAQSNwF-L9IryNGqeBJ1gZ8j25-UYDl0xmwiYzINILxbNDA_exNwNIrLb5H-mbgJcRaqw52A2iH5cYo",
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	var tokenSource = oauthConfig.TokenSource(context.Background(), &token)

	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		fmt.Println("Unable to retrieve Gmail client:", err)
	}

	GmailService = srv
	if GmailService != nil {
		fmt.Println("Email service is initialized")
	}
}

func SendEmailOAUTH2(to string) (bool, error) {
	var message gmail.Message
	emailTo := "To: " + to + "\r\n"
	subject := "Subject: " + "Task2 Test Email form Gmail API using OAuth" + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	msg := []byte(emailTo + subject + mime + "\n" + "Task2: Implementing google sign-in using OAuth.\n\nThis email is sent by gmail API using OAuth.")

	message.Raw = base64.URLEncoding.EncodeToString(msg)

	_, err := GmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		return false, err
	}
	return true, nil
}
