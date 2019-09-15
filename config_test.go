package main

import "testing"

func TestGetConfig(t *testing.T) {
	a := Config{}

	a, err := getConfig("stravaplot.json")

	if err != nil {
		t.Errorf("Could not find configuration file: %s", err)
	}

	if a.Accesstoken != "1234567890123456789012345678901234567890" {
		t.Errorf("Accesstoken incorrect, got: %s, want: %s", a.Accesstoken, "1234567890123456789012345678901234567890")
	}
	if a.Athleteid != "123123" {
		t.Errorf("Athleteid incorrect, got: %s, want: %s", a.Athleteid, "123123")
	}
	if a.Coordinates != "47.5800, -122.3000" {
		t.Errorf("Coordinates incorrect, got: %s, want: %s", a.Coordinates, "47.5800, -122.3000")
	}
	if a.Location != "Seattle" {
		t.Errorf("Location incorrect, got: %s, want: %s", a.Location, "Seattle")
	}
	if a.Zoom != "11" {
		t.Errorf("Zoom incorrect, got: %s, want: %s", a.Zoom, "11")
	}
}
