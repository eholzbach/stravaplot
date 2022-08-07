// load configuration
package config

import (
	"encoding/json"
	"html/template"
	"os"
)

// Config is user configuration read from file
type Config struct {
	Bind         string
	ClientID     string
	ClientSecret string
	Coordinates  template.JS
	Database     string
	Location     string
	Zoom         string
}

// getConfig reads a json configuration file and returns type Config
func getConfig(cpath string) (Config, error) {
	c := Config{}

	// open config file
	file, err := os.Open(cpath)
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
