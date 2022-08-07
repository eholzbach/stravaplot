package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/antihax/optional"
	"github.com/eholzbach/strava"
	"github.com/eholzbach/stravaplot/config"

	_ "github.com/mattn/go-sqlite3"
)

// connectDB creates a connection to sqlite db, creates the table if it does not exist, and returns type DB
func ConnectDB(dbpath string) (*sql.DB, error) {
	// open db
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return nil, err
	}

	// create table if it doesn't exist
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS sp (ID INTEGER PRIMARY KEY, Name VARCHAR, Distance FLOAT, MovingTime INT, ElapsedTime INT, TotalElevationGain FLOAT, Type VARCHAR, StravaID INT UNIQUE, StartDate DATETIME, StartDateLocal DATETIME, Timezone VARCHAR, MapId VARCHAR, MapPolyline VARCHAR, MapSummaryPolyline VARCHAR, AverageSpeed FLOAT, MaxSpeed FLOAT, AveragePower FLOAT, Kilojoules FLOAT, GearId VARCHAR )")

	if err != nil {
		return nil, err
	}

	statement.Exec()

	return db, err
}

// getPolylines queries the db for polyline data and returns it in a slice of strings
func GetPolylines(config config.Config, db *sql.DB) ([]string, error) {
	var data []string

	// query db for all polylines
	rows, err := db.Query("SELECT MapPolyline FROM sp;")

	if err != nil {
		return data, err
	}

	defer rows.Close()

	for rows.Next() {
		var p string
		err = rows.Scan(&p)
		if err != nil {
			return data, err
		}
		data = append(data, p)
	}

	return data, err
}

// updateDB checks for new strava data and writes it into the database
func UpdateDB(oauth context.Context, config config.Config, db *sql.DB) (err error) {
	// get time of most recent activity in db
	row, err := db.Query("SELECT StartDate FROM sp ORDER BY ID DESC LIMIT 1;")
	if err != nil {
		return err
	}

	var ts time.Time
	for row.Next() {
		err = row.Scan(&ts)
		if err != nil {
			return err
		}
	}

	// if empty set to epoch
	t, _ := time.Parse("2006-01-02 15:04:05", "0001-01-01 00:00:00")
	if ts == t {
		ts, _ = time.Parse("1/2/2006 15:04:05", "1/1/1970 12:00:00")
	}

	// strava api after function expects int
	timestamp := int(ts.Unix())

	// build strava api client
	client := strava.NewAPIClient(strava.NewConfiguration())

	// get new activities from strava
	opts := &strava.GetLoggedInAthleteActivitiesOpts{
		PerPage: optional.NewInt32(200),
		After:   optional.NewInt32(int32(timestamp)),
	}

	activities, _, err := client.ActivitiesApi.GetLoggedInAthleteActivities(oauth, opts)

	if err != nil {
		fmt.Println(err)
		return err
	}

	// write new activities to db
	for _, v := range activities {
		// only bicycles, only if gps data
		if *v.Type_ == "Ride" && len(v.Map_.SummaryPolyline) > 0 {
			// get full activity
			opts := &strava.GetActivityByIdOpts{
				IncludeAllEfforts: optional.NewBool(true),
			}
			a, _, err := client.ActivitiesApi.GetActivityById(oauth, v.Id, opts)
			if err != nil {
				return err
			}

			statement, err := db.Prepare("INSERT OR IGNORE INTO sp (Name, Distance, MovingTime, ElapsedTime, TotalElevationGain, Type, StravaID, StartDate, StartDateLocal, Timezone, MapId, MapPolyline, MapSummaryPolyline, AverageSpeed, MaxSpeed, Kilojoules, GearId) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")

			if err != nil {
				return err
			}

			_, err = statement.Exec(a.Name, a.Distance, a.MovingTime, a.ElapsedTime, a.TotalElevationGain, a.Type_, a.Id, a.StartDate, a.StartDateLocal, a.Timezone, a.Map_.Id, a.Map_.Polyline, a.Map_.SummaryPolyline, a.AverageSpeed, a.MaxSpeed, a.Kilojoules, a.GearId)

			if err != nil {
				return err
			}
		}
	}

	return nil
}
