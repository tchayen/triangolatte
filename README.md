# triangolatte

> **NOTE:** _The library is in its very early stage. Near future will bring real package class API, tests and optimizations._

Triangulation library. Allows translating lines and polygons to the language of GPUs.

Features normal and miter joint line triangulation. Handles polygons using ear clipping algorithm with hole elimination included.

As soon as possible, the library will be distrubuted additionally as some kind of [WebAssembly](https://webassembly.org/) module for use in JS for fast triangulation computing.

## Installation

```bash
go get github.com/Tchayen/triangolate
```

## Usage

```go
vertices := []Point{{10, 20}, {30, 40}, {50, 60}}
holes := [][]Point{}
triangolate.EarCut(vertices, holes)
```

## Features

> **NOTE**: _this library is developed mostly in map data triangulation and it will be its main performance target once it reaches benchmarking phase_

#### `line.Normal(points []Point, width int)`
Normal triangulation. Produces ugly in zoom, but fast to compute and sometimes acceptable joints.

#### `line.Miter(points []Point, width)`

Triangulates lines using miter joint. With little computational overhead, produces no more vertices than normal one.

Refer to this [forum post](https://forum.libcinder.org/topic/smooth-thick-lines-using-geometry-shader) for sketches, code examples and ideas.

#### `polygon.EarCut(points, holes [][]Point)`

Based on the following [paper](https://www.geometrictools.com/Documentation/TriangulationByEarClipping.pdf).

## Tests

```bash
go test
```

## Contributing

You are welcome to create an issue or pull request. I will be glad to help.
