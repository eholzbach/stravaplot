package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	strava "github.com/eholzbach/strava"
	"github.com/eholzbach/strava/auth"
	"github.com/eholzbach/strava/config"
)

func main() {
	log.Println("starting service")

	// read configuration
	log.Println("reading configuration file")
	cpath := flag.String("config", "stravaplot.json", "configuration file")
	flag.Parse()

	settings, err := config.getConfig(*cpath)
	if err != nil {
		log.Printf("error reading configuration: %s", err)
		os.Exit(1)
	}

	// connect to db
	log.Println("connecting to db")
	db, err := connectDB(settings.Database)
	if err != nil {
		log.Printf("error connecting to db: %s", err)
		os.Exit(1)
	}

	// authenticate
	log.Println("authenticating to Strava")
	oauth, err := auth.Auth(context.Background(), strava.ContextOAuth2, settings.ClientID, settings.ClientSecret)
	if err != nil {
		log.Printf("error authenticating: %s", err)
		os.Exit(1)
	}

	// register handlers and serve
	log.Println("starting api")

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/render", func(w http.ResponseWriter, r *http.Request) {
		renderPage(oauth, w, r, settings, db)
	})

	mux.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(config.Bind, logRequest(mux))
}
