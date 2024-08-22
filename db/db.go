package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/antihax/optional"
	strava "github.com/eholzbach/strava"

	_ "github.com/mattn/go-sqlite3"
)

// DB represents a db connection
var DB *sql.DB

// Connect creates a connection to sqlite db, creates table if it doesn't exist
func Connect(dbpath string) error {
	var (
		err  error
		stmt *sql.Stmt
		q    = "CREATE TABLE IF NOT EXISTS sp (ID INTEGER PRIMARY KEY, Name VARCHAR, Distance FLOAT," +
			"MovingTime INT, ElapsedTime INT, TotalElevationGain FLOAT, Type VARCHAR, StravaID INT UNIQUE, StartDate DATETIME," +
			"StartDateLocal DATETIME, Timezone VARCHAR, MapId VARCHAR, MapPolyline VARCHAR, MapSummaryPolyline VARCHAR," +
			"AverageSpeed FLOAT, MaxSpeed FLOAT, AveragePower FLOAT, Kilojoules FLOAT, GearId VARCHAR )"
	)

	if DB, err = sql.Open("sqlite3", dbpath); err != nil {
		return err
	} else if stmt, err = DB.Prepare(q); err != nil {
		return err
	}
	stmt.Exec()
	return nil
}

// GetPolylines queries the db for polyline data and returns a slice of strings and error
func GetPolylines() ([]string, error) {
	var (
		data []string
		err  error
		rows *sql.Rows
	)

	if rows, err = DB.Query("SELECT MapPolyline FROM sp;"); err != nil {
		return data, err
	}

	defer rows.Close()
	for rows.Next() {
		var a string
		if err = rows.Scan(&a); err != nil {
			return data, err
		}
		data = append(data, a)
	}
	return data, nil
}

// UpdateDB checks for new strava data and writes it into the database
func UpdateDB(ctx context.Context) error {
	var (
		err error
		ts  time.Time
	)

	if ts, err = getLatestFromDB(); err != nil {
		return err
	} else if err := updateActivities(ctx, ts); err != nil {
		return err
	}
	return nil
}

func getLatestFromDB() (time.Time, error) {
	var (
		err       error
		rows      *sql.Rows
		timestamp time.Time
	)

	if rows, err = DB.Query("SELECT StartDate FROM sp ORDER BY ID DESC LIMIT 1;"); err != nil {
		return timestamp, err
	}

	for rows.Next() {
		if err = rows.Scan(&timestamp); err != nil {
			return timestamp, err
		}
	}

	// if empty set to epoch
	a, _ := time.Parse("2006-01-02 15:04:05", "0001-01-01 00:00:00")
	if timestamp == a {
		timestamp, _ = time.Parse("1/2/2006 15:04:05", "1/1/1970 12:00:00")
		log.Printf("db empty, populating")
	}
	return timestamp, nil
}

/*
func getStravaActivities(ctx context.Context, ts time.Time) ([]strava.SummaryActivity, error) {
	var (
		all    []strava.SummaryActivity
		client = strava.NewAPIClient(strava.NewConfiguration())
	)

	for i := int32(1); ; i++ {
		page, _, err := client.ActivitiesApi.GetLoggedInAthleteActivities(ctx, &strava.ActivitiesApiGetLoggedInAthleteActivitiesOpts{
			Page:    optional.NewInt32(i),
			PerPage: optional.NewInt32(100),
			After:   optional.NewInt32(int32(ts.Unix())),
		})
		if err != nil {
			return all, err
		}

		all = append(all, page...)
		if len(page) <= 0 {
			break
		}
	}
	return all, nil
}
*/
// only bicycles with gps data
func filterRides(ctx context.Context, activities []strava.SummaryActivity) ([]strava.DetailedActivity, error) {
	var (
		client   = strava.NewAPIClient(strava.NewConfiguration())
		filtered []strava.DetailedActivity
	)
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

func insertRides(rides []strava.DetailedActivity) error {
	insert := "INSERT OR IGNORE INTO sp (Name, Distance, MovingTime, ElapsedTime, TotalElevationGain, Type, StravaID, " +
		"StartDate, StartDateLocal, Timezone, MapId, MapPolyline, MapSummaryPolyline, AverageSpeed, MaxSpeed, Kilojoules," +
		"GearId) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

	for _, v := range rides {
		log.Printf("adding activity: %s - %s", v.Name, v.StartDateLocal.Format(time.DateTime))
		stmt, err := DB.Prepare(insert)
		if err != nil {
			return err
		}

		if _, err := stmt.Exec(v.Name, v.Distance, v.MovingTime, v.ElapsedTime, v.TotalElevationGain, v.Type_, v.ExternalId,
			v.StartDate, v.StartDateLocal, v.Timezone, v.Map_.Id, v.Map_.Polyline, v.Map_.SummaryPolyline, v.AverageSpeed,
			v.MaxSpeed, v.Kilojoules, v.GearId); err != nil {
			return err
		}
	}
	return nil
}

func updateActivities(ctx context.Context, ts time.Time) error {
	var (
		client = strava.NewAPIClient(strava.NewConfiguration())
	)

	for i := int32(1); ; i++ {
		page, _, err := client.ActivitiesApi.GetLoggedInAthleteActivities(ctx, &strava.ActivitiesApiGetLoggedInAthleteActivitiesOpts{
			Page:    optional.NewInt32(i),
			PerPage: optional.NewInt32(100),
			After:   optional.NewInt32(int32(ts.Unix())),
		})
		if err != nil {
			return err
		}

		fmt.Println(i)
		fmt.Println(len(page))
		rides, err := filterRides(ctx, page)
		if err != nil {
			return err
		}

		if err := insertRides(rides); err != nil {
			return err
		}

		if len(page) <= 0 {
			break
		}
	}
	return nil
}

/*
	for _, v := range activities {
		// only bicycles with gps data
		if *v.Type_ == "Ride" && len(v.Map_.SummaryPolyline) > 0 {
			var a strava.DetailedActivity

			if a, _, err = client.ActivitiesApi.GetActivityById(ctx, v.Id, &strava.ActivitiesApiGetActivityByIdOpts{
				IncludeAllEfforts: optional.NewBool(true),
			}); err != nil {
				return err
			}

			log.Printf("adding activity: %s", a.Name)
			if stmt, err = DB.Prepare(insert); err != nil {
				return err
			}

			if _, err = stmt.Exec(a.Name, a.Distance, a.MovingTime, a.ElapsedTime, a.TotalElevationGain, a.Type_, a.ExternalId, a.StartDate, a.StartDateLocal,
				a.Timezone, a.Map_.Id, a.Map_.Polyline, a.Map_.SummaryPolyline, a.AverageSpeed, a.MaxSpeed, a.Kilojoules, a.GearId); err != nil {
				return err
			}
		}
	}
*/
