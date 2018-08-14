# gpx

GPX (_GPS eXchange Format_) parsing example.

- load `assets/gpx_tmp` file from assets
- convert it from `*.gpx` to an array of point segments
- triangulate those segments using `triangolatte.Line(Miter, ...)`

## Running

```bash
go run gpx.go
```
