# StravaPlot [![Build Status](https://travis-ci.org/eholzbach/stravaplot.svg?branch=master)](https://travis-ci.org/eholzbach/stravaplot)
#### Bicycle rides over open map tiles.

Ride bicycles. Use [Strava](https://www.strava.com). Collect polylines.

![Example](example/sea.jpg?raw=true "Seattle")

## Configuration
Create a [Strava API token](https://developers.strava.com/docs/getting-started/#account). Grab your `athleteid` with `curl -s -H 'Authorization: Bearer <yourapitoken>' 'https://www.strava.com/api/v3/athlete' | jq .id`

Stravaplot's configuration file is in json.
```
{
        "athleteid": "123123",
        "accesstoken": "1234567890123456789012345678901234567890",
        "bind": "127.0.0.1",
        "coordinates": "47.5800, -122.3000",
        "database": "/var/db/stravaplot/stravaplot.db",
        "location": "Seattle",
        "zoom": "11"
}
```

### Parameters

  **-config** configuration file (default "stravaplot.json")

## Endpoints
### /v1/render
 - Methods: GET
 - Response: 200
 - Function: Updates the database and renders a new map page

### /rides/strava.html
  - Methods: GET
  - Response: 200
  - Function: Serves the rendered map from a static file
