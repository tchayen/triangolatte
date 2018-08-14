# city

What it does:

- loads all buildings in a city from file
- normalizes their coordinates
- triangulates them
- calculates errors
- saves triangles data to `assets/json_tmp`
- finds out correctness %

## Running

> **NOTE:** _get data from Overpass-turbo using the instructions below_

```bash
go run city.go
```

## Fetching data

Data comes from **OpenStreetMap**. You can use the [Overpass-turbo](https://overpass-turbo.eu/) query listed below, for fetching all buildings in
Cracow:

```
[out:json];
(
    way[building](50.0,19.85,50.105,20.13);
    relation[building](50.0,19.85,50.105,20.13);
);
out body;
>;
out skel qt;
```

_(shortened)_

```
[out:json];(way[building](50.0,19.85,50.105,20.13);relation[building](50.0,19.85,50.105,20.13););out body;>;out skel qt;
```

Then look for an option to export data in `*.geojson` format and
save it in `assets/cracow_tmp`.

> **TODO:** _Include data fetching in the program (requires OSM to GeoJSON conversion)._
