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

// C represent parsed config data
var C Config

// GetConfig reads a json configuration file and returns type Config
func GetConfig(cpath string) error {
	var (
		err  error
		file *os.File
	)

	if file, err = os.Open(cpath); err != nil {
		return err
	}

	defer file.Close()

	if err = json.NewDecoder(file).Decode(&C); err != nil {
		return err
	}
	file.Close()
	return nil
}
