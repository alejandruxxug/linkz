# 🦎 Linkz: Editor Design

The motto: **the code is the drawing.** The editor isn't a text box next to a canvas. The code itself is the layers panel and the properties panel.

## 🪟 The window

```
┌─ LAYERS = THE .lzdraw ──────────────────────────────┬──────────────────────────┐
│  ○ □ ╱ ⬠ ～ T ƒ ▤                      ⤓ png  💾    │                          │
│  ƒ horizon      let horizon = 300                   │                          │
│ ▾◉ sun          layer sun                           │                          │
│    ○ core         circle((380, 90), 40, #e0a526)  ◀ │         CANVAS           │
│    ○ glow         circle((380, 90), 60, #e0a52640)  │   click → select line    │
│    ⊕                                                │   drag  → code changes   │
│ ▾◉ sea          layer sea                           │                          │
│   ▸⌇ waves        for (let i = 0; i < 5; ...) {     │                          │
│    □ water        rect((0, horizon), (512, 212), #2d6a8f)                      │
│ ▣ canvas        canvas = {"size": (512, 512), ...}  │                          │
├─ PROPERTIES · core ─────────────────────────────────┤                          │
│  center  ( 380 , 90 )   r  40   fill ■ #e0a526       │                          │
└─────────────────────────────────────────────────────┴──────────────────────────┘
```

## 🎯 v1 finish line (the developer's goal)

> "At least be able to draw with a mouse a line that could be edited like when I draw a vector line on an art program."

1. Pick the **pencil** ✎ and drag on the canvas.
2. Lizard records the mouse points, simplifies them (Ramer–Douglas–Peucker), and writes one line into the current layer: `stroke([(40, 300), (120, 210), (260, 240), (380, 120)], #1d1d1f, 4)`.
3. The stroke is drawn as a smooth curve through its points (Catmull–Rom → cubic Béziers, in lizard).
4. Select it → every point shows as an **anchor handle**. Drag an anchor → that one point literal in the array is rewritten. Double-click the curve → a point is inserted. Select an anchor + Delete → removed. The properties panel shows color and width.
5. Undo works for every step.
6. **Brush with pressure** (trackpad Force Touch / drawing tablet; speed with a mouse) and **default pen styles** (ink, pencil, marker, brush, calligraphy): `brush(points, pressures, color, pens.ink)`; switching the pen re-styles the stroke.

Everything else in v1 serves this.

## 📐 Rules

| Rule | Detail |
|---|---|
| **Rows are lines** | Every line of the file is a row. What you see is exactly the file. |
| **Layers are indentation** | `layer sun` opens a layer; the following lines indented deeper belong to it. Layers nest (a layer in a layer is a group). Code keeps its braces. |
| **Higher = in front, everywhere** | Layers over layers, objects over objects inside a layer. Code still *runs* top to bottom (so `let` goes above its users); only the *paint* order is bottom-up. |
| **Loose shapes are layers** | A shape line outside any layer is its own one-object layer. Simple drawings need no `layer` lines. |
| **One statement = one object** | A loop or `if` is one object that makes many shapes; it expands to show its lines. |
| **Names and flags are comments** | `layer sun -~> [hidden]`, `circle(...) -~> core [locked]`. A row's name is its trailing comment, else the function name. |
| **Row kinds** | layer ▾◉ · shape ○□╱⬠～T · block ⌇ · value ƒ (`let`) · definition ƒ (`function`) · canvas ▣ (pinned bottom, like Photoshop's Background) |

## ✍️ Creating and editing

- **Enter** at the end of a row → a new row **below** it at the same depth (behind it). **Shift+Enter** → above (in front). **Tab / Shift+Tab** → indent into / out of the layer above.
- An empty row shows the **tool strip inline** (○ □ ╱ ⬠ ～ T ƒ ▤ layer). Pick one → the row becomes a template, e.g. `circle((256, 256), 40, #3fa34d)`, the shape appears, and the properties panel focuses its first field. Typing code directly always works too.
- **⊕** at the end of a layer adds an object inside it.
- **Properties** are the call's arguments. Labels are the lizard function's own parameter names (so a user's own `function star(center, points, r)` gets fields too). Editing a field rewrites that argument's text. Literal → typed field (number, point, color swatch); expression → editable code text.
- **Undo/redo** (Cmd+Z / Shift+Cmd+Z) over text edits.
- **Drag rows** (or Alt+↑/↓) to reorder = move lines.

## 🖱️ Two-way canvas (protected for v1)

- **Click** a shape → its row is selected (and its layer expanded). Click a row → its shapes get handles.
- **Drag** a shape → if its position argument is a point literal, the literal is rewritten live (`(380, 90)` → `(402, 77)`). If it's an expression, the drag is refused and the row shows why ("position comes from `horizon + i * 20`").
- **Drag a handle** → radius / size literal changes the same way.
- Selection and hit-testing use the lizard geometry (point in path).

## 🕹️ Look

8-bit UI around a smooth canvas: the UI renders at half resolution and is scaled ×2 with hard pixels; the canvas is a separate full-resolution target composited into its rect. Rules and the Lizard 16 palette are in [[Linkz – Style Guide]].

## 🌍 Languages (built in from the start)

| Rule | Detail |
|---|---|
| No hard-coded UI text | Every label, button, status message and tooltip frame goes through `t("key")` (with `{name}` placeholders: `t("status.saved", {"file": name})`). |
| One file per language | `lib/ui/lang/en.liz`, `lib/ui/lang/es.liz`: a lizard hash of key → text. Adding a language = adding one file. `en` is the fallback for any missing key; a missing key everywhere shows the key itself (easy to spot). |
| Which language | `settings.language` in `~/.linkz/settings.liz` (itself lizard code), else `os.locale()`, else `en`. Switchable from the toolbar; the UI repaints, nothing restarts. |
| Tooltips | The library `-~>` doc comments are English; translations live in `lib/draw/docs/<lang>.liz` keyed by function and parameter name, falling back to English. |
| Layout | Nothing assumes English widths: buttons and fields size from `font.measure` (Spanish runs ~25% longer). |
| Font | The pixel font covers Latin-1 (á é í ó ú ñ ü ¿ ¡ ç …) plus a fallback box. |
| Code stays code | lizard-script keywords (`let`, `layer`, `function`) and the language's own error messages stay English for now; translating runtime errors needs error codes in lizard-script (later). |

## 🧱 How it's built

```
Go host (thin)          window + input via Ebitengine (window ONLY), a pixel buffer, span fill,
                        files, clipboard, the debug backdoor
lizard (everything)     ui.liz       layout, rows, fields, text editing, tool strip, events
                        raster.liz   shapes → paths → scanline spans with soft edges
                        font.liz     bitmap pixel font as lizard data
                        shapes.liz, color.liz, transform.liz, svg.liz
```

### Host ↔ lizard contract

| Lizard calls (host primitives) | Does |
|---|---|
| `fb.size()` | `(w, h)` of the window buffer |
| `fb.span(y, x0, x1, color)` | Fills one horizontal run (color alpha blends). The only per-pixel loop in Go |
| `fb.spans(list)` | Many spans in one call (speed) |
| `fb.clip(pos, size)` / `fb.unclip()` | Restricts spans (dirty regions, panels) |
| `fb.save_png(path, pos, size)` | Writes part of the buffer as PNG (export renders the canvas into an offscreen buffer first: `fb.offscreen(size)` / `fb.screen()`) |
| `src.parse(code)` | The language's own parser as data: statements with line, column, depth, kind, call name, argument text spans, trailing comment |
| `src.run(code)` | Runs the drawing in a fresh environment with a step limit; returns `{"shapes", "errors", "printed", "steps", "ms"}`; every shape tagged with its `line` and `layer` path |
| `fs.read`, `fs.write`, `fs.exists`, `clip.get`, `clip.set`, `time.ms` | OS access |

| Host calls (lizard entry points) | When |
|---|---|
| `app.start(path)` | Once |
| `app.event(e)` | On input only (mouse down/up/move while pressed, key, text, wheel, resize). `e` is a hash with `type`, `pos`, `key`, `text`, `mods`; `move` also has `samples` (below) |
| `app.paint()` | After events, only if lizard marked something dirty; lizard repaints only dirty regions |
| `app.debug_state()` | For the backdoor's `Dump` |
| `app.debug_commands` | A hash of named commands lizard exposes (e.g. `SelectRow`, `SetProp`, `AddShape`) |

- **Move samples:** a `move` event carries `samples`, a list of `{pos, canvas_pos, pressure, t}` (`pos` UI px, `canvas_pos` full-res or null, `t` seconds since start): every native drag event since the previous tick, in order, the last equal to the event's `pos`. The macOS event monitor records them (bounded queue of 4096, y flipped to top-left origin); elsewhere, and for backdoor input, it is the one point of the tick. Coalesced moves concatenate their samples. `down`/`up` carry their native pressure in `pen.pressure`.
- **No 60 fps:** nothing runs while idle. Dragging coalesces mouse moves (latest position only), so a slow frame drops moves instead of queuing them.
- **Caching:** each layer's spans are cached; only layers whose lines changed are re-rasterized.

## 🔧 Agent backdoor (tool calls, not scripts)

`linkz debug <Command> [args]` talks to the running window over `127.0.0.1` with a random token from `~/.linkz/debug.json` (0600). Every input command goes through the same event path as real input.

| Command | Does |
|---|---|
| `List` | Every command, host and lizard ones |
| `Dump` | JSON: file, dirty, rows (kind, depth, name, rect, flags), selection, properties, errors, last run stats, window size |
| `Screenshot [path]` | PNG of the window (agents Read it to look) |
| `Click x y [mods]`, `DoubleClick`, `Drag x1 y1 x2 y2 [steps]`, `Scroll dy` | Mouse |
| `Key name [mods]`, `Type "text"` | Keyboard |
| `SetCode "…"`, `Open path`, `Save [path]`, `Export png path` | File state |
| `Resize w h`, `Stats`, `Wait` (until idle), `Quit` | App |
| `Call Name args…` | A named lizard command from `app.debug_commands` (`SelectRow 3`, `SetProp r 50`, `AddShape circle`) |

## ⚡ Speed plan

Lizard drawing every pixel decision is the bottleneck. In order: event-driven painting → dirty regions → per-layer span caches → `fb.spans` batching → the bytecode VM (M6). Every `Stats` call reports lizard ms per paint, so we measure instead of guessing.
