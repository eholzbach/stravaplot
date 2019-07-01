package main

import (
	"fmt"
	"github.com/strava/go.strava"
	"html/template"
	"log"
	"net/http"
	"time"
)

type renderData struct {
	Coordinates string
	Poly        []strava.Polyline
}

func renderPage(w http.ResponseWriter, r *http.Request, config Config, location string) {
	data := renderData{}
	var date string

	// set the start date
	switch location {
	case "seattle":
		date = "2016-03-01"
		data.Coordinates = "47.6332745,-122.3235537"
	default:
		current := time.Now().Local()
		date = current.Format("2006-01-02")
		data.Coordinates = "47.6332745,-122.3235537"
	}

	sd, err := time.Parse("2006-1-2", date)
	if err != nil {
		log.Println("error parsing timestamp: %s", err)
		return
	}

	t := sd.Unix()

	// build strava api client
	client := strava.NewClient(config.Accesstoken)
	athlete := strava.NewCurrentAthleteService(client)
	service := strava.NewActivitiesService(client)

	// get list of activities
	activities, err := athlete.ListActivities().After(int(t)).Do()

	for _, v := range activities {
		if len(v.Map.SummaryPolyline) > 0 {
			//fmt.Printf("\"%s\",\n", v.Id)
			activity, _ := service.Get(v.Id).IncludeAllEfforts().Do()
			data.Poly = append(data.Poly, activity.Map.Polyline)
		}
	}

	if err != nil {
		log.Print("error getting list of activities: %s", err)
		return
	}

	// render template to file
	tpl, _ := template.ParseFiles("seattle.tpl")

	fmt.Println(tpl)
}
