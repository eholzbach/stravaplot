// load configuration
package main

import (
	"encoding/json"
	"html/template"
	"os"
)

type Config struct {
	Accesstoken string
	Athleteid   string
	Coordinates template.JS
	Location    string
	Zoom        string
}

func getConfig() (Config, error) {
	c := Config{}

	// get user home
	home, err := os.UserHomeDir()
	if err != nil {
		return c, err
	}

	// open config file
	file, err := os.Open(home + "/" + ".stravaplot.json")
	if err != nil {
		return c, err
	}

	defer file.Close()

	// decode json
	d := json.NewDecoder(file)
	err = d.Decode(&c)
	if err != nil {
		return c, err
	}

	file.Close()

	return c, err
}
