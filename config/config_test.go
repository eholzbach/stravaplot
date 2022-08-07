package main

import "testing"

func TestGetConfig(t *testing.T) {
	a := Config{}

	a, err := getConfig("stravaplot.json")

	if err != nil {
		t.Errorf("Could not find configuration file: %s", err)
	}

	if a.ClientSecret != "1234567890123456789012345678901234567890" {
		t.Errorf("ClientSecret incorrect, got: %s, want: %s", a.ClientSecret, "1234567890123456789012345678901234567890")
	}

	if a.ClientID != "123123" {
		t.Errorf("Clientid incorrect, got: %s, want: %s", a.ClientID, "123123")
	}

	if a.Bind != "127.0.0.1" {
		t.Errorf("Bind address incorrect, got %s, want %s", a.Bind, "127.0.0.1")
	}

	if a.Coordinates != "47.5800, -122.3000" {
		t.Errorf("Coordinates incorrect, got: %s, want: %s", a.Coordinates, "47.5800, -122.3000")
	}

	if a.Database != "/var/db/stravaplot/stravaplot.db" {
		t.Errorf("Database location incorrect, got %s, want %s", a.Database, "/var/db/stravaplot/stravaplot.db")
	}

	if a.Location != "Seattle" {
		t.Errorf("Location incorrect, got: %s, want: %s", a.Location, "Seattle")
	}

	if a.Zoom != "11" {
		t.Errorf("Zoom incorrect, got: %s, want: %s", a.Zoom, "11")
	}
}
