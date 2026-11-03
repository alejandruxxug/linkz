# 🏛️ Linkz: Architecture

## 🧭 Principles

1. **Linkz is a lizard-script program.** `main.liz` and `lib/` are the whole app; `lizardscript` runs it (`lizardscript run . edit`) and `lizardscript build` ships it as one executable. There is no Go in this repo (the Go host was removed in B042, M11).
2. **Logic lives in `lib/*.liz`.** Shapes, SVG, rasterizing, tracing, the MCP protocol, JSON, the editor, the app loop: all lizard-script.
3. **Missing language feature → fix the language.** OS access, pixels, windows, speed: that work goes into lizard-script's Terrarium or VM (logged in [[Linkz – Language Needs]]), never into a workaround here.
4. **Slow lizard is data.** Measure it (`tests/editor/editor_bench.liz`, `lizardscript --profile`), log it, and make the language faster.
5. **Per-pixel color math is the Terrarium's.** Gradients and layer blend modes are general Terrarium primitives (`img.span_gradient`, `img.push_layer`, `img.composite`, ...). Lizard decides everything else: it builds the gradient (`lib/draw/gradient.liz`), says which layers need a buffer and in which order they nest (`lib/draw/shapes.liz`, `lib/gfx/raster.liz`), and writes the SVG.

## 🪟 The native app

The full contract is in [[Linkz – Editor Design]] → "How it's built".

```
main.liz      the command line: edit, render, paint, trace, calibrate, mcp
lib/app       the host: the window loop, splash, settings, machine measurement, the
              commands, drawing runs (lizard.run with the drawing libraries alive),
              the scanner, and the names the libraries get (fb, fs, src, linkz)
lib/ui        the editor: rows, properties, tools, text editing, canvas interaction
lib/gfx       rasterizer (paths → spans) and the bitmap font
lib/draw      shapes, geometry, SVG writer
lib/mcp       the MCP server and JSON
lib/trace     picture → drawing
```

The window, its surfaces, pen samples and the agent backdoor (`lizardscript drive`) are the Terrarium's `window`, `img`, `pen` and automation. Nothing runs while idle: the window is on demand, lizard handles an event, repaints dirty regions, the window presents them.

## 🗂️ Folder structure

| Path | Holds |
|---|---|
| `main.liz` | The `linkz` command |
| `lib/app/` | The host in lizard (see above) |
| `lib/` | The libraries; `lib/lizard-check.liz` lists the names the host gives them |
| `bundle.txt` | The data files a built Linkz carries (font, sprites, fixtures) |
| `tests/` | `lizardscript test` programs with `.expected` output, the benchmark |
| `tools/` | Dev-time lizard programs (`pngdiff`, `mcpdiff`, `uses`) |
| `scripts/` | Tests on both engines, parity, smoke, cross build and install |
| `examples/` | `.lzdraw` drawings, also used by tests |
| `docs/` | This knowledge base |
| `../lizard-script` | The language and `lizardscript` (separate repo) |

## ❓ Where does new code go?

| You need | Put it |
|---|---|
| A new shape, color helper, transform | `lib/draw`, `lib/gfx` |
| A file format, a protocol (JSON, MCP) | `lib/*.liz` |
| A command, a setting, how a drawing runs | `lib/app` |
| To touch the OS (files, windows, pixels, clipboard) | A Terrarium module in lizard-script, general for any program |
| A language feature (syntax, type, builtin, speed) | lizard-script repo + its `docs/DESIGN.md` extensions log |

## 🧪 Testing

- `scripts/test.sh`: every `tests/` program on the evaluator and on the VM, with HOME in a folder of its own.
- `scripts/parity.sh`: the outputs against the last Go host's (`scripts/parity/golden`).
- `scripts/smoke.sh`: the commands and a real editor driven with `lizardscript drive`.
- `lizardscript fmt --check` and `lizardscript check` on `main.liz lib tests tools`.
- `go test ./...` in lizard-script for any language change (its `tests/valid` + `tests/errors` programs).
