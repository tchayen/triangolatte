# triangolatte

> **NOTE:** _The library is in its very early stage. Near future will bring
usable API, tests and optimizations._

Triangulation library. Allows translating lines and polygons (both based on
points) to the language of GPUs.

Features normal and miter joint line triangulation. Handles polygons using ear
clipping algorithm with hole elimination included.

#### WebAssembly

One of the core plans for this library's development is creating, as soon as it
becomes possible, some kind of [WebAssembly](https://webassembly.org/) module
for use in JS.

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

> **NOTE**: _this library is developed mostly with map data triangulation in
mind and it will be its main performance target once it reaches benchmarking
phase_

### API

#### `line.Normal(points []Point, width int) (triangles []float64)`

Normal triangulation. Produces joints that are ugly in zoom, but fast to compute
and sometimes acceptable.

#### `line.Miter(points []Point, width int) (triangles []float64)`

Triangulates lines using miter joint. With little computational overhead,
produces no more vertices than normal one. Comes with one drawback: very sharp
angles can potentially explode to infinity.

Refer to this [forum post](https://forum.libcinder.org/topic/smooth-thick-lines-using-geometry-shader)
for sketches, code examples and ideas.

#### `polygon.EarCut(points []Point) (triangles []float64)`

Based on the following [paper](https://www.geometrictools.com/Documentation/TriangulationByEarClipping.pdf).

Removes holes, joining them with the rest of the polygon.

Biggest place for impromevent. Ranging from more effective data structures for
vertex removal, spatial indices, z-ordering and so on...

### Helpers

Some helpers for polygon-related operations.

**`polygon.IsInsideTriangle(t Triangle, p Point) bool`** – checks whether
`point` lies in the `triangle` or not.

**`polygon.IsReflex(a, b, c Point) bool`** - checks if given point `b` is
reflex with respect to neighbor vertices `a` and `c`.

**`polygon.SameSide(p1, p2, a, b Point) bool`** – checks if `p1` lies on the
same side of segment `a<->b` as `p2`.

### Point

Api of the `point` package used for calculating maths.

**`(p Point) Add(r Point) Point`** – add some point to the current one i.e.
`p.Add(Point{0, 1})`.

**`(p Point) Sub(r Point) Point`** – subtract point.

**`(p Point) Scale(f float64) Point`** – scale point **to** given length.

**`(p Point) Normalize() Point`** – normalize point (make length equal to `1`).

**`(p Point) Dot(r Point) float64`** – calculate dot product of two vectors.

**`(p Point) Cross(r Point) float64`** – calculate 2D cross product (like in 3D
but taking only non-zero coordinate as result)

**`(p Point) Distance2(r Point) float64`** – squared distance to given point.
_Useful for comparisons._

**`(p Point) Pair() (x, y float64)`** – split point into pair of floats.

### Types

```
type Point struct {
  X, Y float64
}
```

```
type Set map[int]bool
```

```
type Triangle struct {
  A, B, C Point
}
```

## Tests

```bash
go test
```

## Contributing

You are welcome to create an issue or pull request. I will be glad to help.

## Roadmap

- use real life data for tests
- rename/move directories to match Go's conventions
- provide proper, tested installation and usage info
- complete usage example
- try to break something – more solid error handling
- research WebAssembly usage in Go
- setup meaningful benchmarks
- more agressive optimizations
- docs website on `github.io`
