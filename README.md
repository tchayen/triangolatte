# triangolatte

[![Build Status](https://travis-ci.org/Tchayen/triangolatte.svg?branch=master)](https://travis-ci.org/Tchayen/triangolatte)
[![Coverage Status](https://coveralls.io/repos/github/Tchayen/triangolatte/badge.svg?branch=master)](https://coveralls.io/github/Tchayen/triangolatte?branch=master)
[![GoDoc](https://godoc.org/github.com/Tchayen/triangolatte?status.svg)](https://godoc.org/github.com/Tchayen/triangolatte)

---

_Should I use it?_ **Not yet**

_The algorithm is generally working and should be quite optimal. The tests are
quite rich now and cover many aspects and things that can go wrong. However, I
am at a stage at which I keep finding bugs and not working cases. Come back in
two weeks and you should see vastly improved (and even faster) version,
hopefully ready for a real world usage._

---

2D triangulation library. Allows translating lines and polygons (both based on
points) to the language of GPUs.

Features normal and miter joint line triangulation. Handles polygons using ear
clipping algorithm ~~with hole elimination included~~ (**WIP**).

## Usage

#### Basic example
```go
vertices := []Point{{10, 20}, {30, 40}, {50, 60}}
t, err = triangolatte.EarCut(vertices)
```

## Installation

Nothing surprising
```bash
go get github.com/Tchayen/triangolatte
```

## Examples

For rendering examples, go to `examples/webgl` directory.

You will find instructions for running the code there.

## Docs
[![GoDoc](https://godoc.org/github.com/Tchayen/triangolatte?status.svg)](https://godoc.org/github.com/Tchayen/triangolatte)

Visit [godoc](https://godoc.org/github.com/Tchayen/triangolatte) for complete
documentation of the project.

## Features

> **NOTE**: _this library is developed mostly with map data triangulation in
mind and it will be its main performance target once it reaches benchmarking
phase_

### API

#### `EarCut(points []Point, holes [][]Point) ([]float64, error)`

Based on the following [paper](https://www.geometrictools.com/Documentation/TriangulationByEarClipping.pdf).

Removes holes, joining them with the rest of the polygon.

#### `Normal(points []Point, width int) (triangles []float64)`

Normal triangulation. Produces joints that are ugly in zoom, but fast to compute
and sometimes acceptable.

#### `Miter(points []Point, width int) (triangles []float64)`

Triangulates lines using miter joint. With little computational overhead,
produces no more vertices than normal one. Comes with one drawback: very sharp
angles can potentially explode to infinity.

Refer to this [forum post](https://forum.libcinder.org/topic/smooth-thick-lines-using-geometry-shader)
for sketches, code examples and ideas.

### Types

For calculations using points.
```
type Point struct {
  X, Y float64
}
```

## Tests

Code is (more or less) covered in tests. You can run them:

```bash
go test -v
```

You can for example benchmark speed of checking if point is inside triangle:

```bash
go test -run NONE -bench IsInsideTriangle
```

### Flame Graphs

What is a _flame graph_? Simply speaking, a human-readable insight into what
kept CPU busy (and it resembles fire 🔥).

It has intelligent zoom that lets you narrow down the range to a particular
function from the flame and go back quickly at will.

It also supports hovering while still being a regular, valid `*.svg` file.

![assets/torch.svg](assets/torch.svg)

You can view an example of `EarCut` benchmark flame graph in [assets/torch.svg](assets/torch.svg).

> **NOTE:** _you must display the image file directly to use cool features
described above._

#### Generating flame graph

Install [go-torch](https://github.com/uber/go-torch) and [FlameGraph](https://github.com/brendangregg/FlameGraph)
if you haven't done it before

```bash
go get github.com/uber/go-torch
git clone https://github.com/brendangregg/FlameGraph
export PATH=$PATH:$(pwd)/FlameGraph # You might want to add this to your .bashrc or other equivalent
```

From now on, every time you want to generate a flame graph, simply run the
commands below:

> You can replace `EarCut` with any function name from `*_test.go` file, with
> name starting with `Benchmark*`, stripping the prefix.
>
> For example, `func BenchmarkEarCut(...)` became `EarCut`.

```bash
go test -run NONE -bench EarCut -cpuprofile prof.cpu
go tool pprof triangolatte.test prof.cpu
go-torch triangolatte.test prof.cpu
```

Now you can open newly generated `torch.svg` in your web browser.

## Future plans

### Optimizations

**First:** rewrite ear list from list to a queue implemented on array. It should
be rather easy and will save tons of mallocs.

`EarCut` currently uses `O(n^2)` algorithm which does not have much space
remaining for improvement.

However, there are other obvious options for
speeding up vertex traversal using some kind of z-ordering (checking close
vertices first increases chance of proving that given vertex cannot be an ear
and therefore makes occasion for early return).

### Making the library production grade

By providing more examples, real benchmarks with comparison to other languages,
libraries. Testing behavior on huge, unusual real life examples.

### WebAssembly

One of the core plans for this library's development is creating, as soon as it
becomes possible, some kind of [WebAssembly](https://webassembly.org/) module
for use in JS.

## License

MIT License (refer to the [LICENSE](LICENSE) file).

## Contributing

You are welcome to create an issue or pull request if you've got an idea what to do.

Don't have one, but still want to contribute? Get in touch with me and we can
brainstorm some ideas.
