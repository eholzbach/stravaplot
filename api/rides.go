// Package api provides the http router and data funcs
package api

import (
	"context"
	"log"
	"stravaplot/db"
	"time"

	"github.com/antihax/optional"
	strava "github.com/eholzbach/strava"
)

// only bicycles with gps data
func filterRides(ctx context.Context, client *strava.APIClient, activities []strava.SummaryActivity) ([]strava.DetailedActivity, error) {
	var filtered []strava.DetailedActivity
	for _, v := range activities {
		if *v.Type_ == "Ride" && len(v.Map_.SummaryPolyline) > 0 {
			a, _, err := client.ActivitiesApi.GetActivityById(ctx, v.Id, &strava.ActivitiesApiGetActivityByIdOpts{
				IncludeAllEfforts: optional.NewBool(true),
			})
			if err != nil {
				return filtered, err
			}
			filtered = append(filtered, a)
		}
	}
	return filtered, nil
}

func updateActivities(ctx context.Context, ts time.Time) error {
	var (
		client = strava.NewAPIClient(strava.NewConfiguration())
		err    error
		page   []strava.SummaryActivity
		rides  []strava.DetailedActivity
	)

	for i := int32(1); ; i++ {
		if page, _, err = client.ActivitiesApi.GetLoggedInAthleteActivities(ctx, &strava.ActivitiesApiGetLoggedInAthleteActivitiesOpts{
			Page:    optional.NewInt32(i),
			PerPage: optional.NewInt32(100),
			After:   optional.NewInt32(int32(ts.Unix())),
		}); err != nil {
			return err
		}

		if rides, err = filterRides(ctx, client, page); err != nil {
			log.Print("error: ", err)
		} else if len(rides) > 0 {
			if err = db.InsertRides(rides); err != nil {
				return err
			}
		}

		if len(page) <= 0 {
			break
		}
	}
	return nil
}

// updateDB checks for new strava data and writes it into the database
func updateDB(ctx context.Context) {
	var (
		err error
		ts  time.Time
	)

	if ts, err = db.GetLatestFromDB(); err != nil {
		log.Println(err)
		return
	} else if err := updateActivities(ctx, ts); err != nil {
		log.Println(err)
		return
	}
	log.Println("db updated")
	return
}
