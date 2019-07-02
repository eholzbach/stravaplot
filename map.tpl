<!DOCTYPE html>
<html>
<head>
	
	<title>{{ .Location }}</title>

	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	
	<link rel="shortcut icon" type="image/x-icon" href="docs/images/favicon.ico" />

	<link rel="stylesheet" href="https://unpkg.com/leaflet@1.0.3/dist/leaflet.css" integrity="sha512-07I2e+7D8p6he1SIM+1twR5TIrhUQn9+I6yjqD53JQjFiMf8EtC93ty0/5vJTZGF8aAocvHYNEDJajGdNx1IsQ==" crossorigin=""/>
	<script src="https://unpkg.com/leaflet@1.0.3/dist/leaflet.js" integrity="sha512-A7vV8IFfih/D732iSSKi20u/ooOfj/AGehOKq0f4vLT1Zr2Y+RX7C+w8A1gaSasGtRUZpF/NZgzSAu4/Gc41Lg==" crossorigin=""></script>
	<script src="https://raw.githubusercontent.com/jieter/Leaflet.encoded/68e781908783fb12236fc7d8d90d71d785f892e9/Polyline.encoded.js" integrity="sha512-IxcXX9OwJ72ucNMR833ngaxl3HIXfrm1ZdnHJFpXOhJeLNLfkM/q0iL6lGVt8Xt4yl124ybQn+F/6L+ZmH57kg==" crossorigin=""></script>
        <script type="text/javascript" src="./{{ .Location }}.js"></script>
	<style>
	body {
		padding: 0;
		margin: 0;
	}
	html, body, #mapid {
		height: 100%;
		width: 100%;
	}
	</style>
</head>
<body>

<div id="mapid" ></div>
<script>
	var mymap = L.map('mapid', {zoomControl: false}).setView([{{ .Coordinates }}], 9);
	L.tileLayer('https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token=pk.eyJ1IjoibWFwYm94IiwiYSI6ImNpejY4NXVycTA2emYycXBndHRqcmZ3N3gifQ.rJcFIG214AriISLbB6B5aw', {
		maxZoom: 18,
		attribution: 'Map data &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors, ' +
			'<a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, ' +
			'Imagery © <a href="http://mapbox.com">Mapbox</a>',
		id: 'mapbox.streets'
	}).addTo(mymap);

	L.control.zoom({position:'topright'}).addTo(mymap);
</script>
</body>
</html>
