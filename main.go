// stravaplot
package main

import (
	"context"
	"flag"
	"log"

	"stravaplot/api"
	"stravaplot/auth"
	"stravaplot/config"
	"stravaplot/db"

	strava "github.com/eholzbach/strava"
)

func main() {
	var (
		cpath = flag.String("config", "stravaplot.json", "configuration file")
		ctx   context.Context
		err   error
	)

	flag.Parse()
	if err = config.GetConfig(*cpath); err != nil {
		log.Fatalf("error reading configuration: %s", err)
	} else if err = db.Connect(config.C.Database); err != nil {
		log.Fatalf("error connecting to db: %s", err)
	} else if ctx, err = auth.Auth(context.Background(), strava.ContextOAuth2, config.C.ClientID, config.C.ClientSecret); err != nil {
		log.Fatalf("error authenticating: %s", err)
	}
	api.Start(ctx)
}
