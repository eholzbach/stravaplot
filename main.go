package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	strava "github.com/eholzbach/strava"
	"github.com/eholzbach/stravaplot/auth"
	"github.com/eholzbach/stravaplot/config"
	"github.com/eholzbach/stravaplot/data"
)

func main() {
	log.Println("starting")

	// read configuration
	cpath := flag.String("config", "stravaplot.json", "configuration file")
	flag.Parse()

	log.Println("reading configuration")
	config, err := config.GetConfig(*cpath)

	if err != nil {
		log.Print("error reading configuration: ", err)
		os.Exit(1)
	}

	// connect to db
	log.Println("connecting to database")
	db, err := data.ConnectDB(config.Database)

	if err != nil {
		log.Print("error connecting to db: ", err)
		os.Exit(1)
	}

	// authenticate
	log.Println("starting oauth with Strava, check your browser for issues")

	oauth, err := auth.Auth(context.Background(), strava.ContextOAuth2, config.ClientID, config.ClientSecret)
	if err != nil {
		log.Print("error authenticating: ", err)
		os.Exit(1)
	}

	// regiester handlers and start server
	log.Printf("starting server on %s", config.Bind)
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("/v1/render", func(w http.ResponseWriter, r *http.Request) {
		renderPage(oauth, w, r, config, db)
	})

	http.ListenAndServe(config.Bind, logRequest(mux))
}
