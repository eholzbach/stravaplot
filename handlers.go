package main

import (
	"fmt"
	"github.com/strava/go.strava"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type renderData struct {
	Coordinates template.JS
	Location    string
	Poly        []strava.Polyline
	Zoom        string
}

func renderPage(w http.ResponseWriter, r *http.Request, config Config, location string) {
	data := renderData{}
	var date string

	// select location and date
	switch location {
	case "losangeles":
		date = "2011-06-13"
		data.Coordinates = "34.1000, -118.3800"
		data.Zoom = "12"
	case "sanfrancisco":
		date = "2012-04-07"
		data.Coordinates = "37.8150, -122.4400"
		data.Zoom = "12"
	case "seattle":
		date = "2016-04-09"
		data.Coordinates = "47.5800, -122.3000"
		data.Zoom = "11"
	default:
		current := time.Now().Local()
		date = current.Format("2006-01-02")
		data.Coordinates = "47.5800,-122.3000"
		data.Location = "seattle"
	}

	// set date
	sd, err := time.Parse("2006-1-2", date)
	if err != nil {
		log.Print("error parsing timestamp: %s\n", err)
		return
	}
	timestamp := sd.Unix()

	// build strava api client
	client := strava.NewClient(config.Accesstoken)
	athlete := strava.NewCurrentAthleteService(client)
	service := strava.NewActivitiesService(client)

	// get list of activities
	activities, err := athlete.ListActivities().PerPage(200).After(int(timestamp)).Do()

	for _, v := range activities {
		if len(v.Map.SummaryPolyline) > 0 {
			activity, _ := service.Get(v.Id).IncludeAllEfforts().Do()
			data.Poly = append(data.Poly, activity.Map.Polyline)
		}
	}

	if err != nil {
		log.Print("error getting list of activities: %s\n", err)
		return
	}

	// build config
	a := renderData{
		Coordinates: data.Coordinates,
		Location:    location,
		Poly:        data.Poly,
		Zoom:        data.Zoom,
	}

	// render map template
	tpl, err := template.ParseFiles("map.tpl")
	if err != nil {
		log.Print("error parsing template: %s\n", err)
		return
	}

	// open file
	filename := fmt.Sprintf("static/%s.html", location)
	file, err := os.Create(filename)
	if err != nil {
		log.Print("error writing to file: %s\n", err)
		return
	}

	// execute template and write to file
	tpl.Execute(file, a)
}
