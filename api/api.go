// Package api provides the http endpoints
package api

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"stravaplot/config"
	"stravaplot/db"
)

type renderData struct {
	Coordinates template.JS
	Location    string
	Poly        []string
	Zoom        string
}

// Start fires off the api
func Start(ctx context.Context) {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("/v1/render", func(w http.ResponseWriter, _ *http.Request) {
		go renderPage(ctx)
		w.Write([]byte("accepted"))
	})
	http.ListenAndServe(config.C.Bind, logMiddleware(mux))
}

// logMiddleware captures http requests and prints them to stdout
func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		h.ServeHTTP(w, r)
	})
}

func renderPage(ctx context.Context) {
	// update db with latest info or render what we got
	if err := db.UpdateDB(ctx); err != nil {
		log.Print("error: ", err)
	}

	// get polylines from db
	pl, err := db.GetPolylines()
	if err != nil {
		log.Print("error querying db: ", err)
		return
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
	tpl.Execute(file, renderData{
		Coordinates: config.C.Coordinates,
		Location:    config.C.Location,
		Poly:        pl,
		Zoom:        config.C.Zoom,
	})
}
