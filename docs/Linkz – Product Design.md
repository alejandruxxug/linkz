# 🎨 Linkz: Product Design

## 🔁 The core loop

Write lizard-script → see the drawing update → tweak (code or, later, drag on the canvas) → export.

## 🧱 Scope tiers

| Tier | What's in it |
|---|---|
| **Core** | Floats, dot access, `use` in lizard-script · `math`, `shapes`, `svg` libraries · web editor with live preview · `.lzdraw` open/save · PNG export |
| **Stretch** | `color`, `transform`, `3d` (projection to 2D paths) · MCP with `run` / `render` / `docs` / `save` · image tracing · bytecode VM for speed |
| **Later** | Two-way editing (drag a shape → its number in the code changes) · agents driving the open editor live · self-hosting compiler |
| **Out** | A real 3D engine (meshes, lighting, orbit viewport) · Photoshop-style raster editing (brushes, filters) |

## 🧩 Systems (one line each)

- **Shapes**: circle, rect, line, path, polygon, group; every shape is a hash, made and drawn in one call.
- **SVG writer**: drawing → SVG text.
- **Raster export**: a rasterizer in lizard turns shapes into pixels; the host encodes the PNG.
- **Editor**: code pane, live preview, errors pinned to their line.
- **Image tracing**: reduce colors → trace each region's outline → smooth → fit Bézier curves. Our own implementation (potrace is GPL).
- **MCP**: JSON-RPC over stdin/stdout, written in lizard; `render` returns a PNG so agents see what they drew.
- **3D**: rotate + perspective-project 3D points into 2D paths (wireframes, extrusions, isometric).

## 📓 Design log

| Date | Decision | Why |
|---|---|---|
| 2026-11-02 | A vector art program, not a code editor, game engine or pixel app | Fits lizard's arrays/hashes; runs once per edit, so interpreter speed matters less; the output is a file people keep |
| 2026-11-02 | Built entirely on lizard-script, Go only as a thin host | "I would prefer for the engine to be built entirely on lizard" — Linkz proves the language can carry real software |
| 2026-11-02 | `.lzdraw` = plain lizard-script text | Diffable, and agents can write it directly |
| 2026-11-02 | Raster export (PNG/JPEG/WebP) | Users need images, not just SVG |
| 2026-11-02 | Bundled MCP; agents write lizard-script, not SVG | Every agent becomes a lizard-script user |
| 2026-11-02 | 3D only as projection to 2D | Real 3D is a second project; projection fits vector art |
| 2026-11-02 | Image tracing is a headline feature | Picture → editable vectors |
| 2026-11-02 | The language grows beyond the course design, additively | "If I make it unique it isn't a bad thing"; "nothing diverts the go foundation so everything is purely additive" — every extension is logged in lizard-script's DESIGN.md |
| 2026-11-02 | **Native app, no browser** ("no chrome, I want a native lizard app") | The editor is a native window |
| 2026-11-02 | **Lizard draws every pixel**; Ebitengine only opens the window ("I would prefer lizard to draw, maybe keep the framerate not that high") | Go has no window API of its own; a hand-written Cocoa layer can replace Ebitengine later |
| 2026-11-02 | The code is the layers panel and the properties panel; each line is a row, **top row = front** | The developer's editor idea: "code is a drawing" |
| 2026-11-02 | **Layers are Python-style indentation** (`layer sun` + indented lines); code keeps braces | "for multiple objects in a same layer it would be messy… blocks would be indented" |
| 2026-11-02 | Layer names and hidden/locked flags live in trailing comments | The file stays plain lizard-script |
| 2026-11-02 | Shapes take points: `circle((380, 90), 40, #e0a526)` | Dragging rewrites one literal |
| 2026-11-02 | Language gets color literals, points and units in v1 | The identity: geometry and color are part of the syntax |
| 2026-11-02 | **v1 finish line: draw a line with the mouse and edit it like a vector line** | The developer's goal; two-way canvas editing is protected over panel polish |
| 2026-11-02 | Pressure-sensitive `brush` (Force Touch trackpad + drawing tablets via an NSEvent monitor; speed-based with a mouse) | Trackpad pressure confirmed working by the developer |
| 2026-11-02 | **Default pen styles** (ink, pencil, marker, brush, calligraphy); `brush` stores raw pressures and a pen style shapes them | "add some pen styles for default as to have some brush styles"; changing the pen re-styles a stroke without redrawing it |
| 2026-11-02 | **Tooltips** on functions, properties, tools and pens, from `-~>` doc comments above each function (the user's own functions get them too) | "when we hover over a function or a property make it give a small tooltip" |
| 2026-11-02 | Mini videos in tool/pen tooltips: **later**, once the structure is solid | The developer's call |
| 2026-11-02 | **8-bit UI** at half resolution ×2, Lizard 16 palette; the canvas stays full-res and smooth | "I also want a 8 bit style UI since it's easier for you to draw it that way"; 4× fewer pixels for lizard to fill |
| 2026-11-02 | **Multi-language UI from the start**: `t("key")`, one lizard file per language (en, es first), Latin-1 pixel font | "from the start could we maybe add multi lang support or for it to slot in easily" |
| 2026-11-02 | **Splash banner**: frameless, transparent window with the pixel Linkz logo, a lizard crawling behind it, "powered by lizard-script" | The developer's request; lizard pre-draws the frames, Go only flips them while the editor loads |
| 2026-11-02 | **UI like Aseprite, but black and dark** ("Linkz Dark" palette replaces Lizard 16); views designed in Figma via Figwright and approved there before B007 builds them | The developer's direction; UI is approved as Figma frames first |
| 2026-11-02 | **Home screen** gallery: working + past drawings with small previews rendered by lizard, cached by file hash | The developer's request |
| 2026-11-02 | Splash: lizard walks in place under "ink", bloom held 1 s, setting `skip_loading_anim` in `~/.linkz/settings.liz` (lizard code) | The developer's request |
| 2026-11-03 | **v1.1 features** (design R002 first, build after B009): right-click shape menu with a **Fill editor** (solid / linear / radial gradients as values: `linear(from, to, [stops])`), a **palette sidebar** on the right (HSV picker, recent, the drawing's own named colors `colors.sky`), an **artist mode** (Photoshop/CSP layout, code stashed, same document), a **first-run tutorial** (first shape → color → brush stroke → peek at the code) | The developer's requests |
| 2026-11-03 | Gradients: lizard computes the gradient's parameters, Go evaluates them per pixel along a span (like a shader) | Per-pixel color in lizard is too slow before the VM; logged as the one exception to "no color decisions in Go" |
| 2026-11-03 | **Stepped gradients** (`steps = n`, hard stops) and **layer blend modes** (Normal, Multiply, Screen, Overlay, Darken, Lighten, Add + opacity) written as `blend("multiply")` / `opacity(0.8)` inside the layer; blending math in the host like gradients | The developer's requests; no new syntax needed |
| 2026-11-02 | Named **Linkz** (lizard inkz); files stay `.lzdraw` | The developer's pick |
| 2026-11-02 | lizard-script stays its own repo; Linkz imports it | Language work counts for the course; Linkz stays private until the reveal |
