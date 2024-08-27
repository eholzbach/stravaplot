package db

import (
	"database/sql"
	"log"
	"time"

	strava "github.com/eholzbach/strava"

	_ "github.com/mattn/go-sqlite3"
)

// DB represents a db connection
var DB *sql.DB

type ctxKey string

const (
	limitHeader ctxKey = "ratelimit"
	useHeader   ctxKey = "ratelimituse"
)

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

// GetLatestFromDB returns the timestamp of latest ride stored in db and error
func GetLatestFromDB() (time.Time, error) {
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

// InsertRides accepts a slice of type []strava.DetailedActivity and inserts it to the db
func InsertRides(rides []strava.DetailedActivity) error {
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
