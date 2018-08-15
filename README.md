<h1 align="center">
  <img src="assets/logo.png"><br />
  triangolatte
</h1>

<p align="center">
  <a href="https://travis-ci.org/tchayen/triangolatte">
    <img src="https://travis-ci.org/tchayen/triangolatte.svg?branch=master" alt="Travis status">
  </a>
  <a href="https://coveralls.io/github/tchayen/triangolatte?branch=master">
    <img src="https://coveralls.io/repos/github/tchayen/triangolatte/badge.svg?branch=master" alt="Coveralls status">
  </a>
  <a href="https://godoc.org/github.com/tchayen/triangolatte">
    <img src="https://godoc.org/github.com/tchayen/triangolatte?status.svg" alt="Godoc reference">
  </a>
  <a href="https://gitter.im/triangolatte/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge">
    <img src="https://badges.gitter.im/triangolatte/Lobby.svg" alt="Join the chat at https://gitter.im/triangolatte/Lobby">
  </a>
</p>

2D triangulation library. Allows translating lines and polygons (both based on
points) to the language of GPUs.

Features normal and miter joint line triangulation. Handles polygons using ear
clipping algorithm with hole elimination included.

> **For reference:** _triangulates 99.76% of 75 thousand buildings in Cracow under 3.43s on
average programmer notebook (single threaded)._

## Table of contents

- [Installation](#installation)
- [Usage](#usage)
- [Examples](#examples)
- [Features](#features)
  - [API](#api)
  - [Types](#types)
- [Tests](#tests)
- [Benchmarks](#benchmarks)
  - [Flame Graphs](#flame-graphs)
- [Future plans](#future-plans)
- [Contributing](#contributing)
- [License](#license)

## Installation

Nothing surprising
```bash
go get github.com/tchayen/triangolatte
```

## Usage

#### Basic example
```go
vertices := []Point{{10, 20}, {30, 40}, {50, 60}}
t, err = triangolatte.Polygon(vertices)
```

## Examples

In `/examples` you can find:
- **buildings** – full-blown WebGL previewer of buildings triangulated in _city_ example
- **city** – triangulation of whole city downloaded from Open Street Map
- **gpx** – GPX format parsing and triangulation of its data
- **wireframe** – desktop OpenGL wireframe previewer for triangulated shapes

You will find instructions for running the code there.

## Features

### API

#### `Polygon(points []Point, holes [][]Point) ([]float64, error)`

Takes array of points and produces array of triangle coordinates.

Based on the following [paper](https://www.geometrictools.com/Documentation/TriangulationByEarClipping.pdf) and inspired by [EarCut](https://github.com/mapbox/earcut).

#### `JoinHoles(points [][]Point) ([]Point, error)`

Removes holes, joining them with the rest of the polygon. Provides preprocessing
for `Polygon`. First element of the points array is the outer polygon, the rest
of them are considered as holes to be removed.

#### `Line(joint Joint, points []Point, width int) ([]float64, error)`

Takes array of points and triangulates them to resemble a line of given
width. Returns array of two-coordinate CCW triangles one after another.

### Types

To select method of joining line segments.
```go
type Joint int

const (
	// No joint correction.
	Normal Joint = 0
	// Producing miter joints, i.e. extending the lines until they meet at some point.
	Miter Joint = 1
)
```

For calculations using points.
```go
type Point struct {
  X, Y float64
}
```

A wrapper for Point used in cyclic list.
```go
type Element struct {
	Prev, Next *Element
	Point      Point
}
```

## Tests

Code is (more or less) covered in tests. You can run them like this:

```bash
go test -v
```

You can also run benchmarks for selected functions (refer to the `*_test.go` files for availability). For example:

```bash
go test -run NONE -bench IsInsideTriangle
```

## Benchmarks

> **NOTE:** _This section contains work in progress. Numbers below are better reference point than nothing, but still far from perfect._

`Polygon()` on shape with 10 vertices takes `754ns` on average.

Triangulation of 75 thousand buildings runs for around `3.43s`.

_Using average programmer's notebook. Expect speed up on faster CPUs or while splitting execution into separate threads._

### Flame Graphs

CPU time % usage snaphost using Flame Graphs:

![assets/torch.svg](assets/torch.svg)

Want to learn what is it or maybe you are willing to generate one yourself? Check [FlameGraphs](flame_graphs.md) document in this repository.

## Future plans

### Optimizations

> **NOTE**: _this library is developed mostly with map data triangulation in
mind and it will be its main performance target._

- explore possibilities for optimizations in `JoinHoles(...)`
- maybe allow reusing point array for massive allocation reduction

### Content

- provide more examples (e.g. desktop OpenGL usage, mobile app, live rendering pipeline, other unusual use cases...)
- add benchmarks with comparison to libraries in
other languages

### WebAssembly

One of the core plans for this library's development is creating, as soon as it
becomes possible, some kind of [WebAssembly](https://webassembly.org/) module
for use in JS.

## Contributing

You are welcome to create an issue or pull request if you've got an idea what to
do. It is usually a good idea to visit [Gitter](https://gitter.im/triangolatte/Lobby)
and discuss your thoughts.

Don't have one, but still want to contribute? Get in touch with us and we can
brainstorm some ideas.

## License

MIT License – refer to the [LICENSE](LICENSE) file.
