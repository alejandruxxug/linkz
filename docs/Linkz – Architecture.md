# 🏛️ Linkz: Architecture

## 🧭 Principles

1. **The host is thin.** Go code is either the lizard-script runtime (in the lizard-script repo) or a *primitive*: one small function lizard can't write itself (read/write a file, encode a PNG, serve HTTP, read stdin, the clock, randomness). A primitive does one OS-level thing and no drawing logic.
2. **Logic lives in `lib/*.liz`.** Shapes, SVG, rasterizing, tracing, the MCP protocol, JSON: all lizard-script.
3. **Missing language feature → fix the language.** If lizard needs floats, `use` or more speed, that work goes into lizard-script (logged in [[Linkz – Language Needs]]), not into a Go workaround here.
4. **Slow lizard is data.** Measure it, log it, and make the language faster.

## 🪟 The native app

The full contract is in [[Linkz – Editor Design]] → "How it's built".

```
Go host       Ebitengine as a window ONLY (pixels in, input out), host/pressure (trackpad + tablet),
              fb span filling, src.parse / src.run, files, the debug backdoor
lib/ui        the editor: rows, properties, tools, text editing, canvas interaction   (lizard)
lib/gfx       rasterizer (paths → spans) and the bitmap font                          (lizard)
lib/draw      shapes, geometry, SVG writer                                             (lizard)
```

Nothing runs while idle: lizard handles an event, repaints dirty regions, the host presents the buffer.

## 🗂️ Folder structure

| Path | Holds | Language |
|---|---|---|
| `main.go` | The `linkz` command | Go |
| `host/` | Runtime setup + primitives | Go |
| `lib/` | The libraries, embedded in the binary by `lib.go` | lizard-script |
| `tools/` | Dev-time programs (font generator, pressure probe) | Go |
| `examples/` | `.lzdraw` drawings, also used by tests | lizard-script |
| `docs/` | This knowledge base | |
| `../lizard-script` | The language (separate repo, imported with a `replace` in `go.mod`) | Go |

## ❓ Where does new code go?

| You need | Put it |
|---|---|
| A new shape, color helper, transform | `lib/*.liz` |
| A file format, a protocol (JSON, MCP) | `lib/*.liz` |
| To touch the OS (files, network, stdin, pixels in/out) | A primitive in `host/`, as small as possible |
| A language feature (syntax, type, builtin, speed) | lizard-script repo + its `docs/DESIGN.md` extensions log |

## 🧪 Testing

- `go test ./...` here: host tests render the `examples/` and check the SVG.
- `go test ./...` in lizard-script for any language change (its `tests/valid` + `tests/errors` programs).
