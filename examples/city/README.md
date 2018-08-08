# city

What is it for (i.e. goals):

- fetch all buildings in a city
- normalize their coordinates
- triangulate them
- calculate errors
- find out correctness %

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

The data is later exported to `*.geojson` format and
saved in `assets/cracow_tmp`.
