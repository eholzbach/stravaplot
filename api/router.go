// Package api provides the http router and data funcs
package api

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"

	"stravaplot/config"
	"stravaplot/db"
)

// Start starts the api
func Start(ctx context.Context) {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("/v1/rides", rides(ctx))
	http.ListenAndServe(config.C.Bind, logMiddleware(mux))
}

// logMiddleware captures http requests and prints them to stdout
func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		h.ServeHTTP(w, r)
	})
}

func rides(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("action") == "render" {
			if err := render(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(http.StatusText(http.StatusOK)))
			return
		} else if r.URL.Query().Get("action") == "update" {
			go updateDB(ctx)
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte(http.StatusText(http.StatusAccepted)))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
	}
}

func render() error {
	var (
		err  error
		file *os.File
		pl   []string
		tpl  *template.Template
	)

	// get polylines from db, parse template, write to disk
	if pl, err = db.GetPolylines(); err != nil {
		return err
	} else if tpl, err = template.ParseFiles("map.tpl"); err != nil {
		return err
	} else if file, err = os.Create("static/index.html"); err != nil {
		return err
	} else if err = tpl.Execute(file, struct {
		Coordinates template.JS
		Location    string
		Poly        []string
		Zoom        string
	}{
		Coordinates: config.C.Coordinates,
		Location:    config.C.Location,
		Poly:        pl,
		Zoom:        config.C.Zoom,
	}); err != nil {
		return err
	}
	file.Close()
	return nil
}
