// load configuration
package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Accesstoken string `json:"accesstoken"`
	Activityid  string `json:"activityid"`
	Athleteid   string `json:"athleteid"`
}

func getConfig() (Config, error) {

	c := Config{}

	// get user home
	home, err := os.UserHomeDir()
	if err != nil {
		log.Print(err)
		return c, err
	}

	// open config file
	file, err := os.Open(home + "/" + ".stravaplot.json")
	if err != nil {
		log.Print(err)
		log.Print("error reading configuration file, exiting...")
		return c, err
	}

	defer file.Close()

	// decode json
	d := json.NewDecoder(file)
	err = d.Decode(&c)
	if err != nil {
		log.Print(err)
		log.Print("error parsing configuration file, exiting...")
		return c, err
	}

	file.Close()

	return c, err
}
