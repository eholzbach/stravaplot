package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	strava "github.com/eholzbach/strava"
)

func main() {
	log.Print("starting...")

	// read configuration
	cpath := flag.String("config", "stravaplot.json", "configuration file")
	flag.Parse()

	config, err := getConfig(*cpath)
	if err != nil {
		log.Print("error reading configuration: ", err)
		os.Exit(1)
	}

	// connect to db
	db, err := connectDB(config.Database)
	if err != nil {
		log.Print("error connecting to db: ", err)
		os.Exit(1)
	}

	// authenticate
	oauth, err := auth(context.Background(), strava.ContextOAuth2, config.ClientID, config.ClientSecret)
	if err != nil {
		log.Print("error authenticating: ", err)
		os.Exit(1)
	}

	// set rendering handler
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/render", func(w http.ResponseWriter, r *http.Request) {
		renderPage(oauth, w, r, config, db)
	})

	// set display handler
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	http.ListenAndServe("0.0.0.0:8000", logRequest(mux))
}
