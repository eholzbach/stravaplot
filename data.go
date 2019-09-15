package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/strava/go.strava"
	"time"
)

func connectDB(dbpath string) (*sql.DB, error) {
	// open db
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return nil, err
	}

	// create table if it doesn't exist
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS sp (ID INTEGER PRIMARY KEY, Name VARCHAR, Distance FLOAT, MovingTime INT, ElapsedTime INT, TotalElevationGain FLOAT, Type VARCHAR, StravaID INT UNIQUE, StartDate DATETIME, StartDateLocal DATETIME, TimeZone VARCHAR, City VARCHAR, State VARCHAR, Country VARCHAR, MapId VARCHAR, MapPolyline VARCHAR, MapSummaryPolyline VARCHAR, AverageSpeed FLOAT, MaximunSpeed FLOAT, AveragePower FLOAT, Kilojoules FLOAT, GearId VARCHAR )")

	if err != nil {
		return nil, err
	}

	statement.Exec()

	return db, err
}

func getPolylines(config Config, db *sql.DB) ([]string, error) {
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

func updateDB(config Config, db *sql.DB) (err error) {
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
	client := strava.NewClient(config.Accesstoken)
	athlete := strava.NewCurrentAthleteService(client)
	service := strava.NewActivitiesService(client)

	// get new activities from strava
	activities, err := athlete.ListActivities().PerPage(200).After(timestamp).Do()

	if err != nil {
		fmt.Println(err)
		return err
	}

	// write new activities to db
	for _, v := range activities {
		// only bicycles, only if gps data
		if v.Type == "Ride" && len(v.Map.SummaryPolyline) > 0 {
			// get full activity
			a, err := service.Get(v.Id).IncludeAllEfforts().Do()
			if err != nil {
				return err
			}

			statement, err := db.Prepare("INSERT OR IGNORE INTO sp (Name, Distance, MovingTime, ElapsedTime, TotalElevationGain, Type, StravaID, StartDate, StartDateLocal, TimeZone, City, State, Country, MapId, MapPolyline, MapSummaryPolyline, AverageSpeed, MaximunSpeed, AveragePower, Kilojoules, GearId) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")

			if err != nil {
				return err
			}

			_, err = statement.Exec(a.Name, a.Distance, a.MovingTime, a.ElapsedTime, a.TotalElevationGain, a.Type, a.Id, a.StartDate, a.StartDateLocal, a.TimeZone, a.City, a.State, a.Country, a.Map.Id, a.Map.Polyline, a.Map.SummaryPolyline, a.AverageSpeed, a.MaximunSpeed, a.AveragePower, a.Kilojoules, a.GearId)

			if err != nil {
				return err
			}
		}
	}

	return nil
}
