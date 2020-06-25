package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type renderData struct {
	Coordinates template.JS
	Location    string
	Poly        []string
	Zoom        string
}

// renderPage polls strava for new data, updates the local database, and writes a freshly populated static html file for the display handler to server
func renderPage(oauth context.Context, w http.ResponseWriter, r *http.Request, config Config, db *sql.DB) {
	// update db with latest info
	err := updateDB(oauth, config, db)
	if err != nil {
		log.Print("error updating db: ", err)
		return
	}

	// get polylines from db
	pl, err := getPolylines(config, db)
	if err != nil {
		log.Print("error querying db: ", err)
		return
	}

	// render map template
	data := renderData{
		Coordinates: config.Coordinates,
		Location:    config.Location,
		Poly:        pl,
		Zoom:        config.Zoom,
	}

	tpl, err := template.ParseFiles("map.tpl")
	if err != nil {
		log.Print("error parsing template: ", err)
		return
	}

	// open file
	filename := fmt.Sprintf("static/index.html")
	file, err := os.Create(filename)
	if err != nil {
		log.Print("error writing to file: ", err)
		return
	}

	// execute template and write to file
	tpl.Execute(file, data)
}

// logRequest captures http requests and prints them to stdout
func logRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		h.ServeHTTP(w, r)
	})
}
