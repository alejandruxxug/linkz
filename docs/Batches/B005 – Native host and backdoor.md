# B005 – Native host and backdoor

Milestone: M3 Editor · Repo: `linkz` · Lane: Host (Go) · Needs: lizard-script B001–B004 (escapes, dot, step limit, floats, math, colors, points, layers, hooks, text/list helpers)

Read first: `docs/Linkz – Editor Design.md` (the contract), `docs/Linkz – Architecture.md`, `../lizard-script/docs/SEMANTICS.md`.

## Goal
`linkz edit drawing.lzdraw` opens a native window whose every pixel is decided by lizard-script, and `linkz debug <Command>` lets an agent drive and inspect it. Plus a headless `linkz paint` for tests.

## The rule
Go does OS work and memset-style pixel filling only. **No drawing decisions in Go**: no shapes, curves, text, layout, colors chosen by Go. If a primitive starts deciding *what* to draw, it belongs in lizard.

## Pieces

### 1. Window (`host/window`)
- Ebitengine v2.10.5 used **only** as a window: `WritePixels` of our RGBA buffer, input polling, resize, clipboard is not in Ebitengine (use `golang.design/x/clipboard` only if pure-Go on macOS; otherwise keep the clipboard inside the app and note it as a gap). Never call Ebitengine's drawing, text or vector APIs.
- Behind a small interface (`Open(title, w, h)`, `Events()`, `Present(buf)`), so a hand-written Cocoa layer can replace it later.
- `SetScreenClearedEveryFrame(false)`, low TPS (30), present only when lizard repainted. Nearest-neighbor scaling (crisp pixel look). Logical pixels = buffer pixels.
- Events → lizard hashes: `{"type": "down"|"up"|"move"|"key"|"text"|"wheel"|"resize", "pos": (x, y), "button": 0, "key": "Enter", "text": "a", "dy": 1, "mods": {"shift": bool, "alt": bool, "cmd": bool, "ctrl": bool}, "pen": {"pressure": 0.0–1.0, "tilt": (x, y), "rotation": deg, "source": "mouse"|"trackpad"|"tablet", "eraser": bool}}`. `pen` comes from `host/pressure` (B010, built in parallel: Force Touch trackpad + drawing tablets through an NSEvent monitor; call its `Start()` from the first Update); without it `pressure` is 1.0 and `source` "mouse". The backdoor's `Drag` takes an optional pressure ramp (`--pressure 0.2,1.0`) so agents can test brush strokes. `move` only while a button is down or when lizard asked for hover (`app.wants_hover`). **Coalesce**: if lizard is slower than input, consecutive moves collapse to the latest. Key repeat for held keys.

### 2. Framebuffer primitives (`host/fb.go`), injected as a `fb` hash of builtins
| Primitive | Does |
|---|---|
| `fb.size()` | `(w, h)` |
| `fb.span(y, x0, x1, color)` | Fills pixels `[x0, x1)` of row `y`; FLOAT ends get partial coverage at the two end pixels (that's how lizard gets soft edges); color alpha blends source-over |
| `fb.spans(list, offset, color)` | `list` = array of `[y, x0, x1]` (numbers), drawn at `offset` (a point) in `color`. One call per glyph / cached shape |
| `fb.fill(pos, size, color)` | Axis-aligned rectangle (identical spans) |
| `fb.clip(pos, size)` / `fb.unclip()` | Spans outside are cut |
| `fb.offscreen(size)` / `fb.screen()` | Switch target to a new offscreen buffer / back |
| `fb.save_png(path)` | Saves the current target |
| `fb.dirty(pos, size)` | Tells the host which region changed (present only that) |
Everything clips to the buffer; out-of-range never panics.

### 3. Source primitives (`host/src.go`), a `src` hash
- `src.parse(code)` → `{"lines": [...], "statements": [...], "errors": [...]}`:
  - lines: `{"line", "indent", "text", "comment", "name", "flags"}` (name = trailing comment text without `[flags]`; flags from `[hidden]`, `[locked]`).
  - statements (recursively, with `children` for layer bodies and blocks): `{"line", "end_line", "column", "depth", "kind": "layer"|"call"|"let"|"function"|"assign"|"block"|"other", "name", "callee", "args"}`; each arg `{"text", "line", "start", "end", "kind": "number"|"point"|"color"|"string"|"array"|"expr", "items"}` where `items` gives the same for the elements of an array literal (so one point inside `stroke([...])` can be rewritten). Spans come from token positions: an argument runs from its first token to just before the `,` / `)` that ends it at depth 0.
- **Docs for tooltips**: function statements also carry `"params": [names]` and `"doc"`: the consecutive `-~>` comment lines directly above the definition, parsed as `{"summary": first paragraph, "params": {name: text}, "example": text or null}` (param lines look like `-~> r: the radius in pixels`, an example line like `-~> example: circle((256, 256), 40, #3fa34d)`). `src.library_docs()` returns the same for every function the drawing libraries define (`{name: {"params", "doc"}}`), so the editor can show tooltips for built-in drawing functions.
- `src.run(code)` → `{"shapes", "errors": [{"line","column","message"}], "printed", "steps", "ms"}`:
  - fresh environment, the drawing libraries loaded (`lib/draw/*.liz`), `evaluator.MaxSteps` (e.g. 2,000,000), `ResetSteps`, `ResetMath`, `puts` captured.
  - Tags every shape: the libraries' `draw(shape)` calls host `src.here()`, which returns `{"line", "layer": [names…]}`, kept by the B004 hooks: the line is the innermost **user** statement that isn't inside a function body (collect those statement nodes from the user AST; ignore library statements), so a call to the user's own function is tagged with the calling line; the layer path from `OnLayer`.
  - Nested evaluation (lizard UI calling `src.run`) must restore the UI's own limits and output after.
- `fs.read`, `fs.write` (atomic: temp + rename), `fs.exists`; `time.ms()`; `clip.get`, `clip.set` (if a clipboard is available).

### 4. App runtime (`host/app.go`)
- Loads `lib/ui/*.liz` (order from `lib.go`) into the UI environment, calls `app.start(path)`, then per event batch: `app.event(e)` for each (coalesced), then `app.paint()` if `app.needs_paint` is true. Errors in UI code: show a red bar with the message using a built-in last-resort… **no**: print it to stderr and keep the last good frame; the backdoor's `Dump` reports it.
- This batch ships a **placeholder** `lib/ui/app.liz`: gray background, a panel, and a rectangle that follows the mouse while held, plus `app.debug_state()` and `app.debug_commands` with one command (`Ping`). The real editor is B007.

### 5. Backdoor (`host/debug.go`, `linkz debug`)
- HTTP on `127.0.0.1:0`; port + random token in `~/.linkz/debug.json` (mode 0600, removed on exit). Requests without the token header are rejected (stops other local processes and web pages).
- `linkz debug <Command> [args]` reads that file, posts, prints JSON.
- Commands (see the Editor Design table): `List`, `Dump`, `Screenshot [path]`, `Click x y [mods]`, `DoubleClick x y`, `Drag x1 y1 x2 y2 [steps]`, `Scroll dy`, `Key name [mods]`, `Type text`, `SetCode text`, `Open path`, `Save [path]`, `Export png path`, `Resize w h`, `Stats`, `Wait`, `Quit`, `Call Name args…` (lizard `app.debug_commands`). Input commands are injected into the same event queue as real input and return after the resulting paint (`Wait` semantics).
- `Stats`: last paint ms, last event ms, last `src.run` ms/steps, frames painted, events coalesced.

### 6. Commands
- `linkz edit [file.lzdraw]` (new file if missing), `linkz render file.lzdraw [out.svg]` (keep working), `linkz paint prog.liz out.png [WxH]` (headless: runs a lizard program with `fb` on an offscreen buffer; used by tests of the lizard graphics libraries), `linkz debug …`.

## Tests
- Go unit tests: span coverage and blending, clipping, `fs.write` atomicity, `src.parse` spans (calls, nested points in arrays, layers, comments/flags, multi-line statements), `src.run` tagging (top-level, inside a layer, inside a loop, inside a user function → calling line), step limit error, backdoor token rejection.
- A smoke script `scripts/smoke.sh`: start `linkz edit` in the background, `debug Wait`, `Ping`, `Drag`, `Screenshot`, `Dump`, `Quit`; exit non-zero on failure.

## Done when
`go test ./...` green; `scripts/smoke.sh` passes and its screenshot shows the placeholder; `go vet` and `gofmt -l .` clean; `linkz render examples/hello.lzdraw` still works.
