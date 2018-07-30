# city

[Overpass-turbo](https://overpass-turbo.eu/) query for fetching all buildings in
Cracow.

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

```
[out:json];(way[building](50.0,19.85,50.105,20.13);relation[building](50.0,19.85,50.105,20.13););out body;>;out skel qt;
```
