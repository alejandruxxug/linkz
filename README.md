# Linkz 🦎

*lizard inkz*

Vector art you write in [lizard-script](https://github.com/alejandruxxug/lizard-script). Linkz itself is built on lizard-script: the shapes, the SVG writer and everything after them are lizard code, on top of a thin Go host.

```
-~> lzdraw 1
canvas = {"width": 400, "height": 200, "background": "#f4f1e8"}
line(20, 100, 380, 100)
for (let i = 0; i < 6; i = i + 1) {
  circle(40 + i * 64, 100, 6 + i * 4)
}
```

## Try it

Linkz needs `lizard-script` checked out next to it.

```bash
go build -o linkz .
./linkz render examples/hello.lzdraw   # writes examples/hello.svg
```

The plan, the architecture and the design log are in [docs/](docs/Linkz.md).
