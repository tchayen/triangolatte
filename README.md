# triangolate
Triangulation library for the Go ecosystem.

Features triangulation of **lines** and **polygons** (handles those with holes too).

## Installation
```bash
go get github.com/Tchayen/triangolate
```

## Usage
```go
vertices := []Point{{10, 20}, {30, 40}, {50, 60}}
triangolate.earCut(vertices)
```

## Tests
```bash
go test
```

## Contributing

You are welcome to create an issue or pull request. I will be glad to help.
