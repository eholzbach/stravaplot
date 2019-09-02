package main

import (
	"log"
	"log/syslog"
	"net/http"
	"os"
)

func logRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		h.ServeHTTP(w, r)
	})
}

func main() {
	// set up syslog
	logwriter, err := syslog.New(syslog.LOG_NOTICE, "stravaplot")
	if err == nil {
		log.SetOutput(logwriter)
	}

	log.Print("starting...")

	// read configuration
	config, err := getConfig()
	if err != nil {
		log.Print("error reading configuration: ", err)
		os.Exit(1)
	}

	// connect to db
	db, err := connectDB()
	if err != nil {
		log.Print("error connecting to db: ", err)
		os.Exit(1)
	}

	// set rendering handler
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/render", func(w http.ResponseWriter, r *http.Request) {
		renderPage(w, r, config, db)
	})

	// set display handler
	mux.Handle("/rides/", http.StripPrefix("/rides/", http.FileServer(http.Dir("./static"))))

	http.ListenAndServe("0.0.0.0:8000", logRequest(mux))
}
