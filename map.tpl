<!DOCTYPE html>
<html>
<head>
	
	<title>{{ .Location }}</title>

	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	
	<link rel="shortcut icon" type="image/x-icon" href="docs/images/favicon.ico" />

	<link rel="stylesheet" href="https://unpkg.com/leaflet@1.0.3/dist/leaflet.css" integrity="sha512-07I2e+7D8p6he1SIM+1twR5TIrhUQn9+I6yjqD53JQjFiMf8EtC93ty0/5vJTZGF8aAocvHYNEDJajGdNx1IsQ==" crossorigin=""/>
	<script src="https://unpkg.com/leaflet@1.0.3/dist/leaflet.js" integrity="sha512-A7vV8IFfih/D732iSSKi20u/ooOfj/AGehOKq0f4vLT1Zr2Y+RX7C+w8A1gaSasGtRUZpF/NZgzSAu4/Gc41Lg==" crossorigin=""></script>
	<script src="./polyline.encoded.js" integrity="sha512-IxcXX9OwJ72ucNMR833ngaxl3HIXfrm1ZdnHJFpXOhJeLNLfkM/q0iL6lGVt8Xt4yl124ybQn+F/6L+ZmH57kg==" crossorigin=""></script>
	<style>
	body {
		padding: 0;
		margin: 0;
	}
	html, body, #map {
		height: 100%;
		width: 100%;
	}
	</style>
</head>
<body>

<div id="map" ></div>
<script>
	var map = L.map('map').setView([{{ .Coordinates }}], {{ .Zoom }});
		L.tileLayer('http://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}.png', {
			attribution:'&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>, &copy; <a href="https://carto.com/attributions">CARTO</a>',
			subdomains: 'abcd',
			maxZoom: 20,
			minZoom: 0
	}).addTo(map);

	var encodedRoutes = [
	{{ range .Poly }}
		{{ . }},
	{{ end }}
	]

	for (let encoded of encodedRoutes) {
		var coordinates = L.Polyline.fromEncoded(encoded).getLatLngs();

		L.polyline(
			coordinates,
			{
				color: 'white',
				weight: 1,
				opacity: .7
			}
		).addTo(map);
	}
</script>
</body>
</html>
