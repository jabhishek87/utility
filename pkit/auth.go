package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService struct {
	config *Config
}

func NewAuthService(config *Config) *AuthService {
	return &AuthService{config: config}
}

func (a *AuthService) GetClient(ctx context.Context) (*http.Client, error) {
	tok, err := a.tokenFromFile()
	if err != nil || !tok.Valid() {
		if tok != nil && tok.RefreshToken != "" {
			config, err := a.getOAuthConfig()
			if err != nil {
				return nil, err
			}
			tok, err = a.refreshToken(config, tok)
			if err == nil {
				// a.saveToken(tok)
				a.saveCredentialsWithToken(tok, a.config.Google.Auth.TokenFile)
			}
		}
		tok, err = a.runLocalServer()
		if err != nil {
			return nil, err
		}
		a.saveCredentialsWithToken(tok, a.config.Google.Auth.TokenFile)
	}

	config, err := a.getOAuthConfig()
	if err != nil {
		return nil, err
	}
	return config.Client(ctx, tok), nil
}

func (a *AuthService) getOAuthConfig() (*oauth2.Config, error) {
	b, err := os.ReadFile(a.config.Google.Auth.CredentialsFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}
	fmt.Printf("Using scopes: %v\n", a.config.Google.Auth.Scopes)
	return google.ConfigFromJSON(b, a.config.Google.Auth.Scopes...)
}

func (a *AuthService) runLocalServer() (*oauth2.Token, error) {
	config, err := a.getOAuthConfig()
	if err != nil {
		return nil, err
	}

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, err
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	config.RedirectURL = fmt.Sprintf("http://localhost:%d", port)

	codeCh := make(chan string)
	errorCh := make(chan error)

	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			if code == "" {
				errorCh <- fmt.Errorf("no code in callback")
				return
			}
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, "<html><body><h1>Authentication successful!</h1><p>You can close this window.</p></body></html>")
			codeCh <- code
		}),
	}

	go server.Serve(listener)
	defer server.Close()

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the following URL to authorize the application:\n%s\n", authURL)

	select {
	case code := <-codeCh:
		return config.Exchange(context.Background(), code)
	case err := <-errorCh:
		return nil, err
	case <-time.After(2 * time.Minute):
		return nil, fmt.Errorf("authentication timeout")
	}
}

func (a *AuthService) tokenFromFile() (*oauth2.Token, error) {
	f, err := os.Open(a.config.Google.Auth.TokenFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func (a *AuthService) refreshToken(config *oauth2.Config, token *oauth2.Token) (*oauth2.Token, error) {
	return config.TokenSource(context.Background(), token).Token()
}

func (a *AuthService) saveCredentialsWithToken(token *oauth2.Token, filename string) error {
	credsData, err := os.ReadFile(a.config.Google.Auth.CredentialsFile)
	if err != nil {
		return err
	}

	var creds map[string]interface{}
	err = json.Unmarshal(credsData, &creds)
	if err != nil {
		return err
	}

	installed := creds["installed"].(map[string]interface{})
	combined := map[string]interface{}{
		"client_id":     installed["client_id"],
		"client_secret": installed["client_secret"],
		"access_token":  token.AccessToken,
		"token_type":    token.TokenType,
		"refresh_token": token.RefreshToken,
		"expiry":        token.Expiry,
		"scopes":        a.config.Google.Auth.Scopes,
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Printf("Unable to cache oauth token: %v", err)
		return
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(combined)
}
