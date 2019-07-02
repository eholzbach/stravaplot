package main

import (
	"github.com/strava/go.strava"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type renderData struct {
	Coordinates string
	Location    string
	Poly        []strava.Polyline
}

func renderPage(w http.ResponseWriter, r *http.Request, config Config, location string) {
	data := renderData{}
	var date string

	// select location and date
	switch location {
	case "seattle":
		date = "2016-03-01"
		data.Coordinates = "47.6332745,-122.3235537"
	default:
		current := time.Now().Local()
		date = current.Format("2006-01-02")
		data.Coordinates = "47.6332745,-122.3235537"
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
	activities, err := athlete.ListActivities().After(int(timestamp)).Do()

	for _, v := range activities {
		if len(v.Map.SummaryPolyline) > 0 {
			//fmt.Printf("\"%s\",\n", v.Id)
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
		Poly:        data.Poly,
	}

	// render map template
	tpl, err := template.ParseFiles("map.tpl")
	if err != nil {
		log.Print("error parsing template: %s\n", err)
		return
	}

	// open file
	file, err := os.Create("seattle.html")
	if err != nil {
		log.Print("error writing to file: %s\n", err)
		os.Exit(1)
	}

	// execute template and write to file
	tpl.Execute(file, a)

	// render data template
	tpl, err = template.ParseFiles("js.tpl")
	if err != nil {
		log.Print("error parsing template: %s\n", err)
		return
	}

	// open file
	file, err = os.Create("seattle.js")
	if err != nil {
		log.Print("error writing to file: %s\n", err)
		os.Exit(1)
	}

	// execute template and write to file
	tpl.Execute(file, a)

}
