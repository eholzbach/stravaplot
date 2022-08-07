# StravaPlot
#### Bicycle rides over open map tiles.

Ride bicycles. Use [Strava](https://www.strava.com). Collect polylines.

![Example](example/sea.jpg?raw=true "Seattle")

## Usage
Start the service, then make a request to `curl localhost:8000/v1/render` and allow it time to query Strava. This doesn't respect pagination yet, so you may need to do this a few times until the db is fully populated. Fire up a browser and view the map at `http://localhost:8000/`

## Configuration
This authenticates against Strava using [oauth](https://developers.strava.com/docs/getting-started/#oauth). You can find your client secret and id in Strava's settings.

Stravaplot's configuration file is in json.
```
{
        "clientid": "123123",
        "clientsecret": "1234567890123456789012345678901234567890",
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

### /
  - Methods: GET
  - Response: 200
  - Function: Serves the rendered map from a static file
