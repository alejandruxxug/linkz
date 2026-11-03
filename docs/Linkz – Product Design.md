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
| 2026-11-02 | The language may diverge from the course design | "If I make it unique it isn't a bad thing" — every divergence is logged in lizard-script's DESIGN.md |
| 2026-11-02 | Named **Linkz** (lizard inkz); files stay `.lzdraw` | The developer's pick |
| 2026-11-02 | lizard-script stays its own repo; Linkz imports it | Language work counts for the course; Linkz stays private until the reveal |
