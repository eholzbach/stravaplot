package auth

import (
	"context"
	"encoding/gob"
	"fmt"
	"hash/fnv"
	"io/ioutil"
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

// Mostly copypasta from github.com/srabraham/strava-oauth-helper

// Auth authenticates to strava with oauth2
func Auth(parentCtx context.Context, oauth2ContextType fmt.Stringer, id string, secret string) (context.Context, error) {
	c := &oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.strava.com/oauth/authorize",
			TokenURL: "https://www.strava.com/oauth/token",
		},
		Scopes: []string{"read,activity:read_all"},
	}

	tok := getOAuthToken(parentCtx, c)
	tokSource := c.TokenSource(parentCtx, tok)
	oauthCtx := context.WithValue(parentCtx, oauth2ContextType, tokSource)

	return oauthCtx, nil
}

// osUserCacheDir creates directory to store oauth token data
func osUserCacheDir() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatalf("Error getting UserCacheDir: %v", err)
	}

	subDir := filepath.Join(cacheDir, "OAuthTokens")
	if err := os.MkdirAll(subDir, 0770); err != nil {
		log.Fatalf("Failed getting or making cache dir: %v", err)
	}

	return subDir
}

func tokenCacheFile(config *oauth2.Config) string {
	hash := fnv.New32a()

	hash.Write([]byte(config.ClientID))
	hash.Write([]byte(config.ClientSecret))
	hash.Write([]byte(strings.Join(config.Scopes, " ")))

	fn := fmt.Sprintf("%s%v", "strava-auth-tok", hash.Sum32())

	return filepath.Join(osUserCacheDir(), url.QueryEscape(fn))
}

// tokenFromFile reads oauth2 token from file
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	t := new(oauth2.Token)
	err = gob.NewDecoder(f).Decode(t)

	return t, err
}

// saveToken saves oauth2 token to file
func saveToken(file string, token *oauth2.Token) {
	f, err := os.Create(file)
	if err != nil {
		log.Printf("Warning: failed to cache oauth token: %v", err)
		return
	}

	defer f.Close()

	gob.NewEncoder(f).Encode(token)
}

// getOAuthToken fetches an oauth2 token from cache
func getOAuthToken(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	cacheFile := tokenCacheFile(config)

	token, err := tokenFromFile(cacheFile)
	if err != nil {
		token = tokenFromWeb(ctx, config)
		saveToken(cacheFile, token)
		log.Printf("Saved new token %#v to %q", token, cacheFile)
	} else {
		log.Printf("Using cached token %#v from %q", token, cacheFile)
	}

	return token
}

// tokenFromWeb fetches a new token
func tokenFromWeb(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	ch := make(chan string)
	randState := fmt.Sprintf("st%d", time.Now().UnixNano())

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/favicon.ico" {
			http.Error(rw, "", 404)
			return
		}
		if req.FormValue("state") != randState {
			log.Printf("State doesn't match: req = %#v", req)
			http.Error(rw, "", 500)
			return
		}
		if code := req.FormValue("code"); code != "" {
			fmt.Fprintf(rw, "<h1>Success</h1>Authorized.")
			rw.(http.Flusher).Flush()
			ch <- code
			return
		}
		log.Printf("no code")
		http.Error(rw, "", 500)
	}))
	defer ts.Close()

	config.RedirectURL = ts.URL
	authURL := config.AuthCodeURL(randState)

	go openURL(authURL)

	log.Printf("Authorize this app at: %s", authURL)

	code := <-ch
	log.Printf("Got code: %s", code)

	token, err := config.Exchange(ctx, code)
	if err != nil {
		log.Fatalf("Token exchange error: %v", err)
	}

	return token
}

// openURL uses xdg-utils to spawn a browser window for the user to approve oauth2
func openURL(url string) {
	try := []string{"xdg-open", "google-chrome", "open"}

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

	slurp, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading %q: %v", filename, err)
	}

	return strings.TrimSpace(string(slurp))
}
