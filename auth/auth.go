// Package auth provides oauth2
// contains copypasta from github.com/srabraham/strava-oauth-helper
package auth

import (
	"context"
	"encoding/gob"
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

// Auth authenticates to strava with oauth2
func Auth(ctx context.Context, oauth2ContextType fmt.Stringer, id string, secret string) (context.Context, error) {
	c := &oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.strava.com/api/v3/oauth/authorize",
			TokenURL: "https://www.strava.com/api/v3/oauth/token",
		},
		Scopes: []string{"read,activity:read_all"},
	}

	tokSource := c.TokenSource(ctx, getOAuthToken(ctx, &oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.strava.com/api/v3/oauth/authorize",
			TokenURL: "https://www.strava.com/api/v3/oauth/token",
		},
		Scopes: []string{"read,activity:read_all"},
	}))
	return context.WithValue(ctx, oauth2ContextType, tokSource), nil
}

// osUserCacheDir creates directory to store oauth token data
func osUserCacheDir() string {
	var (
		cacheDir string
		err      error
	)
	if cacheDir, err = os.UserCacheDir(); err != nil {
		log.Fatalf("error getting token cache: %v", err)
	}
	subDir := filepath.Join(cacheDir, "OAuthTokens")
	if err = os.MkdirAll(subDir, 0770); err != nil {
		log.Fatalf("failed i/o on cache directory: %v", err)
	}
	return subDir
}

func tokenCacheFile(config *oauth2.Config) string {
	hash := fnv.New32a()
	hash.Write([]byte(config.ClientID))
	hash.Write([]byte(config.ClientSecret))
	hash.Write([]byte(strings.Join(config.Scopes, " ")))
	return filepath.Join(osUserCacheDir(), url.QueryEscape(fmt.Sprintf("%s%v", "strava-auth-tok", hash.Sum32())))
}

// tokenFromFile reads oauth2 token from file
func tokenFromFile(file string) (*oauth2.Token, error) {
	var (
		err error
		f   *os.File
		t   = new(oauth2.Token)
	)
	if f, err = os.Open(file); err != nil {
		return nil, err
	} else if err = gob.NewDecoder(f).Decode(t); err != nil {
		return nil, err
	}
	return t, nil
}

// saveToken saves oauth2 token to file
func saveToken(file string, token *oauth2.Token) {
	f, err := os.Create(file)
	if err != nil {
		log.Printf("failed to cache oauth token: %v", err)
		return
	}
	defer f.Close()
	gob.NewEncoder(f).Encode(token)
	log.Printf("saved token to %q", file)
}

// getOAuthToken fetches an oauth2 token from cache
func getOAuthToken(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	var (
		cacheFile = tokenCacheFile(config)
		err       error
		token     *oauth2.Token
	)
	if token, err = tokenFromFile(cacheFile); err != nil {
		token = tokenFromWeb(ctx, config)
		saveToken(cacheFile, token)
	} else {
		log.Printf("using cached token from %q", cacheFile)
	}
	return token
}

// tokenFromWeb fetches a new token
func tokenFromWeb(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	var (
		ch        = make(chan string)
		code      string
		randState = fmt.Sprintf("st%d", time.Now().UnixNano())
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/favicon.ico" {
			http.Error(w, "", 404)
			return
		} else if r.FormValue("state") != randState {
			log.Printf("State doesn't match: req = %#v", r)
			http.Error(w, "", 500)
			return
		} else if code = r.FormValue("code"); code == "" {
			log.Printf("no code")
			http.Error(w, "", 500)
			return
		}

		fmt.Fprintf(w, "<h1>Success</h1>Authorized.")
		w.(http.Flusher).Flush()
		ch <- code
	}))
	defer ts.Close()

	config.RedirectURL = ts.URL
	authURL := config.AuthCodeURL(randState)
	go openURL(authURL)
	log.Printf("Authorize this app at: %s", authURL)
	code = <-ch
	log.Printf("Got code: %s", code)

	token, err := config.Exchange(ctx, code)
	if err != nil {
		log.Fatalf("Token exchange error: %v", err)
	}
	return token
}

// openURL uses xdg-utils to spawn a browser window for the user to approve oauth2
func openURL(url string) {
	try := []string{"xdg-open", "firefox", "open"}
	for _, bin := range try {
		err := exec.Command(bin, url).Run()
		if err == nil {
			return
		}
	}
	log.Printf("Error opening URL in browser.")
}

func valueOrFileContents(value string, filename string) string {
	if value != "" {
		return value
	}
	slurp, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading %q: %v", filename, err)
	}
	return strings.TrimSpace(string(slurp))
}
