package main

import (
	"testing"
	"time"
)

func TestDB(t *testing.T) {
	// create db
	db, err := connectDB(":memory:")
	if err != nil {
		t.Errorf("Connection to db failed, got: %s", err)
	}

	c := Config{
		Accesstoken: "1234567890123456789012345678901234567890",
		Athleteid:   "123123",
		Coordinates: "47.5800, -122.3000",
		Location:    "Seattle",
		Zoom:        "11",
	}

	// write a row
	statement, err := db.Prepare("INSERT OR IGNORE INTO sp (Name, Distance, MovingTime, ElapsedTime, TotalElevationGain, Type, StravaID, StartDate, StartDateLocal, TimeZone, City, State, Country, MapId, MapPolyline, MapSummaryPolyline, AverageSpeed, MaximunSpeed, AveragePower, Kilojoules, GearId) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")

	if err != nil {
		t.Errorf("Preparing sql statement failed, got: %s", err)
	}

	_, err = statement.Exec("testing", "1234.0", "1234", "1234", "1234.0", "Ride", "1234", time.Now(), time.Now(), "PST", "Los Angeles", "CA", "USA", "123123", "123123123123123123123123", "123123123", "10.1", "12.1", "2.0", "3.0", "space horse")

	if err != nil {
		t.Errorf("Executing sql statement failed, got: %s", err)
	}

	// get polylines
	pl, err := getPolylines(c, db)
	if err != nil {
		t.Errorf("Error reading polylines from db, got: %s", err)
	}

	if pl[0] != "123123123123123123123123" {
		t.Errorf("Error reading polylines, got: %s, want: %s", pl[0], "123123123123123123123123")
	}
}
