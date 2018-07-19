# triangolatte

> **NOTE:** _The library is in its very early stage. Near future will bring real package class API, tests and optimizations._

> **NOTE:** I have been writing Go for ~6 hours now, not everything will be perfect. _Also: proper usage examples will come later._

Triangulation library. Allows translating lines and polygons to the language of GPUs.

Features normal and miter joint line triangulation. Handles polygons using ear clipping algorithm with hole elimination included.

#### WebAssembly

One of the core plans for this library's development is creating, as soon as it becomes possible, some kind of [WebAssembly](https://webassembly.org/) module for use in JS.

## Installation

Nothing surprising
```bash
go get github.com/Tchayen/triangolatte
```

## Usage

```go
vertices := []Point{{10, 20}, {30, 40}, {50, 60}}
holes := [][]Point{}
triangolate.EarCut(vertices, holes)
```

## Features

> **NOTE**: _this library is developed mostly with map data triangulation in mind and it will be its main performance target once it reaches benchmarking phase_

#### `line.Normal(points []Point, width int)`
Normal triangulation. Produces joints that are ugly in zoom, but fast to compute and sometimes acceptable.

#### `line.Miter(points []Point, width)`

Triangulates lines using miter joint. With little computational overhead, produces no more vertices than normal one. Comes with one drawback: very sharp angles can potentially explode to infinity.

Refer to this [forum post](https://forum.libcinder.org/topic/smooth-thick-lines-using-geometry-shader) for sketches, code examples and ideas.

#### `polygon.EarCut(points, holes [][]Point)`

Based on the following [paper](https://www.geometrictools.com/Documentation/TriangulationByEarClipping.pdf).

Removes holes, joining them with the rest of the polygon.

Biggest place for impromevent. Ranging from more effective data structures for vertex removal, spatial indices, z-ordering and so on...

## Tests

```bash
go test
```

## Contributing

You are welcome to create an issue or pull request. I will be glad to help.
