## Flame Graphs

What is a _flame graph_? Simply speaking, a human-readable insight into what
kept CPU busy (and it resembles fire 🔥).

It has intelligent zoom that lets you narrow down the range to a particular
function from the flame and go back quickly at will.

It also supports hovering while still being a regular, valid `*.svg` file.

![assets/torch.svg](assets/torch.svg)

You can view an example of `Polygon` benchmark flame graph in [assets/torch.svg](assets/torch.svg).

> **NOTE:** _you must display the image file directly to use cool features
described above._

## Generating flame graph

Install [go-torch](https://github.com/uber/go-torch) and [FlameGraph](https://github.com/brendangregg/FlameGraph)
if you haven't done it before.

```bash
go get github.com/uber/go-torch

# In any directory, $HOME for example.
git clone https://github.com/brendangregg/FlameGraph

# You might want to add this to your .bashrc or other equivalent.
# NOTE: replace the path with the one you chose for your FlameGraph installation.
export PATH=$PATH:$HOME/FlameGraph
```

From now on, every time you want to generate a flame graph, simply run the
commands below:

> You can replace `Polygon` with any function name from `*_test.go` file, with
> name starting with `Benchmark*`, stripping the prefix.
>
> For example, `func BenchmarkPolygon(...)` became `Polygon`.

```bash
go test -run NONE -bench Polygon -cpuprofile prof.cpu
go tool pprof triangolatte.test prof.cpu
go-torch triangolatte.test prof.cpu
```

Now you can open newly generated `torch.svg` in your web browser.
