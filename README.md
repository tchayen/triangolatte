# triangolatte

[![GoDoc](https://godoc.org/github.com/Tchayen/triangolatte?status.svg)](https://godoc.org/github.com/Tchayen/triangolatte)

> **NOTE:** _The library is in its very early stage. Near future will bring more
tests and optimizations._

Triangulation library. Allows translating lines and polygons (both based on
points) to the language of GPUs.

Features normal and miter joint line triangulation. Handles polygons using ear
clipping algorithm with hole elimination included.

## Usage

#### Basic example
```go
vertices := []Point{{10, 20}, {30, 40}, {50, 60}}
t, err = triangolatte.EarCut(vertices, [][]Point{})
```

## Installation

Nothing surprising
```bash
go get github.com/Tchayen/triangolatte
```

## Examples

For rendering examples, go to `examples/webgl` directory.

You will find instructions for running the code there.

## Features

> **NOTE**: _this library is developed mostly with map data triangulation in
mind and it will be its main performance target once it reaches benchmarking
phase_

#### WebAssembly

One of the core plans for this library's development is creating, as soon as it
becomes possible, some kind of [WebAssembly](https://webassembly.org/) module
for use in JS.

### API

#### `Normal(points []Point, width int) (triangles []float64)`

Normal triangulation. Produces joints that are ugly in zoom, but fast to compute
and sometimes acceptable.

#### `Miter(points []Point, width int) (triangles []float64)`

Triangulates lines using miter joint. With little computational overhead,
produces no more vertices than normal one. Comes with one drawback: very sharp
angles can potentially explode to infinity.

Refer to this [forum post](https://forum.libcinder.org/topic/smooth-thick-lines-using-geometry-shader)
for sketches, code examples and ideas.

#### `EarCut(points []Point, holes [][]Point) ([]float64, error)`

Based on the following [paper](https://www.geometrictools.com/Documentation/TriangulationByEarClipping.pdf).

Removes holes, joining them with the rest of the polygon.

Currently uses `O(n^2)` algorithm which does not have much space remaining for
improvement. Most obvious option is speeding up vertex traversal using some kind
of z-ordering (checking close vertices first increases chance of proving that
given vertex cannot be an ear and therefore makes occasion for early return).

### Helpers

Some helpers for polygon-related operations.

**`IsInsideTriangle(t Triangle, p Point) bool`** – checks whether
`point` lies in the `triangle` or not.

**`IsReflex(a, b, c Point) bool`** - checks if given point `b` is
reflex with respect to neighbor vertices `a` and `c`.

### Types

For calculations using points.
```
type Point struct {
  X, Y float64
}
```

## Tests

Running tests during development happens in IDE and project currently lacks
support for running tests in the CLI. I am sorry for that.

**TODO:** _provide an instruction set for running tests (all at once)._

## Contributing

You are welcome to create an issue or pull request if you've got an idea what to do.

Don't have one, but still want to contribute? Get in touch with me via email (**TODO:** _put some email address here_) and we will brainstorm some ideas.

## Roadmap

- fix bug with incorrect hole elimination in AGH A0 example
- provide proper, tested installation and usage info
- complete usage example
- try to break something – more solid error handling
- research WebAssembly usage in Go
- optimize using z-order curve
