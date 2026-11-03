# B030 – the Terrarium: lizard-script's standard library

Milestone: std (new, target for 2026-11-03) · Repo: `../lizard-script` · Lane: Go + lizard · Needs: nothing (runs in parallel with B029, see "Don't collide")

> "Set the build target for today for lizardscript standard library." Naming: "quirky umbrella, plain modules" — the developer, 2026-11-03

## Goal

Any lizard program can open a window, read input (with pen pressure), draw pixels, read and write files, load and save images, parse JSON and tell the time — **with no host of its own**:

```
lizardscript paint.liz
```

The library is called **the Terrarium** (where lizards live). Its modules have plain names. This is how a whole app (an editor, a game) can be 100% `.liz`: the Go lives in the language, not in the app.

The developer relaxed the "additive only" rule of lizard-script's CLAUDE.md for this work (2026-11-03), but the existing tests still must stay green and nothing that works today may break.

## Shape

- Go package `lizardscript/terrarium` with `terrarium.Install(env *object.Environment)` that defines the module hashes in `env` (each module a hash of builtins, like `math`/`text`/`list` are; use dot access: `fs.read(path)`).
- Each module in its own file (`terrarium/fs.go`, …) or subpackage when it has OS-specific code (`terrarium/pen/pen_darwin.go`, `pen_windows.go`, `pen_other.go`).
- Modules written **in lizard** where speed isn't critical (`json` first): `.liz` files embedded with `//go:embed`, run into the env by `Install`.
- `window` depends on Ebitengine (used **only** to open a window, get events and show a finished pixel buffer — lizard draws every pixel). Put it in `terrarium/window` so a program/host that doesn't want a window doesn't link it: `terrarium.Install` installs everything; `terrarium.InstallCore` everything but `window`/`pen`.
- Opaque values (a window, a pixel buffer, an image): add `object.Native` (a new type in a new file `object/native.go`: `Kind string`, `Value any`; `Type()` "NATIVE", `Inspect()` like `<image 640x480>`). Equality: identity.
- Every builtin: checks its arguments with the repo's message style (`fs.read: want a STRING path, got INTEGER`), returns `*object.Error` on failure (catchable by `try`).

## Modules (v0)

| Module | Functions | Source to port |
|---|---|---|
| `fs` | `read(path)`, `write(path, text)` (atomic), `append(path, text)`, `exists(path)`, `list(dir)` → names, `mkdir(path)`, `remove(path)` (refuses a directory with content), `home()`, `join(a, b, …)` | `../linkz/host/osprims.go` (fsHash, WriteFileAtomic) |
| `os` | `args()`, `env(name)`, `exit(code)`, `platform()` ("macos"/"windows"/"linux"), `locale()`, `config_dir()` | osprims.go (osHash, Locale) |
| `time` | `ms()` (monotonic), `now()` → hash {year, month, day, hour, minute, second}, `sleep(ms)` | osprims.go (timeHash) |
| `json` (lizard) | `parse(text)`, `write(value, indent = 0)` | `../linkz/lib/mcp/json.liz` |
| `img` | `load(path, max_side = 0)` → image, `save(image, path)` (PNG), `new(w, h, color = #00000000)`, `size(image)`, `get(image, x, y)`, `set(image, x, y, color)`, `pixels(image)` → flat list of colors, `from_pixels(size, pixels)` | `../linkz/host/img.go` |
| `clip` | `get()`, `set(text)` | osprims.go (clipHash) |
| `window` | `open(w, h, title, options = {})` (options: `scale`, `resizable`, `frameless`, `fps`), `run(w, tick)` (calls the lizard `tick(events)` every frame until it returns `false` or the window closes), events as hashes: `{"type": "down"/"up"/"move"/"key"/"type"/"scroll"/"resize"/"close", "pos", "button", "key", "mods", "text", "dy", "pen": {"pressure", "source"}}`; drawing on the window's pixel buffer: `size(w)`, `fill(w, rect, color)`, `span(w, y, x0, x1, color)`, `spans(w, list, color)`, `blit(w, image, pos)`, `clear(w, color)`, `title(w, text)` | `../linkz/host/window/*`, `../linkz/host/fb.go` (the general parts only) |
| `pen` | pressure comes in the window's events (`e.pen.pressure`, `e.pen.source` "mouse"/"trackpad"/"tablet"); `pen.available()` | `../linkz/host/pressure/*` (macOS NSEvent via purego, Windows Ink) |

Copy and generalize; do **not** change anything in `../linkz` (it migrates to the Terrarium in a later batch).

## CLI

`lizardscript program.liz` installs the Terrarium before running (so `fs`, `img`, `window` … just exist). The REPL too, minus `window`. A user `let fs = …` shadows a module like any global.

## Examples (in `examples/terrarium/`, each runs)

- `files.liz`: write, read, list, JSON round-trip.
- `image.liz`: make an image, draw a gradient pixel by pixel, save a PNG.
- `paint.liz`: a 600×400 window; drag to paint circles whose size follows pen pressure; `c` clears; Esc quits. (~60 lines, the "an app in pure lizard" demo.)
- `clock.liz`: a window with an analog clock drawn with `span`s.

## Docs

`docs/TERRARIUM.md`: what it is, one section per module (every function: signature, what it does, errors, an example), how to add a module. Link it from README. A row in DESIGN.md §7 is **not** needed (rule relaxed) but add a short "The Terrarium" section to DESIGN.md's architecture part — only after B029 has finished editing DESIGN.md (the orchestrator will tell you, or add it at the end and report the diff).

## Tests

Table-driven tests per module (files in `t.TempDir()`, images in memory, JSON round-trips incl. escapes and nesting, argument errors). `window`: test the event conversion and the buffer functions headless (no window opens in `go test`). `go test ./...` and `go vet ./...` green; build for darwin, windows and linux (`GOOS=… go build ./...`).

## Don't collide (B029 runs at the same time in the same repo)

B029 (the bytecode VM) is editing `evaluator/`, `object/object.go`, `repl/`, `main.go`, `docs/DESIGN.md`, `docs/SEMANTICS.md`, `tests/programs_test.go` and adding `code/`, `compiler/`, `vm/`. **Don't edit those files.** You may add `object/native.go`, `terrarium/**`, `examples/terrarium/**`, `docs/TERRARIUM.md`, and the README link. For the CLI wiring (`main.go`, `repl`), write the exact change in your report instead of applying it; the orchestrator applies it after B029.

## Done when

All modules work with tests, the four examples run (paint.liz opened by hand by the developer — report how to run it), cross-builds pass, docs written.
