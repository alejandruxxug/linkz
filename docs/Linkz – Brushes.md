# 🖌️ Linkz: Brushes (note for later)

> "Keep a note on how we could import brushes and export brushes. Brush creator would allow to parametrize brushes Procreate style." — the developer, 2026-11-03

Status: **idea, not scheduled** (after v1.1). Builds on the pens in `lib/draw/pens.liz` and the tunable pressure from B009.

## 🧩 What a brush is in Linkz

A brush is a **lizard hash**, the same thing a pen is today, grown to Procreate-style parameters. Because it's lizard, a brush is also code: a user can write one by hand, and the brush creator just writes it for them.

| Group | Parameters (Procreate-like) |
|---|---|
| Stroke path | spacing, streamline (smoothing), jitter (position), fall-off |
| Taper | start/end length, pressure taper, size/opacity at the tips |
| Shape | a stamp shape (procedural: round, square, flat nib, star…, or a small bitmap), rotation (fixed / follow stroke / random), scatter, count per stamp, roundness |
| Grain | a texture tiled inside the stroke (procedural noise or a bitmap), moving vs fixed grain, scale, depth (contrast) |
| Dynamics | speed → size/opacity, jitter of size/opacity/color per stamp |
| Pressure | the B009 block: curve, size/opacity min–max, per-input (trackpad / tablet / mouse) overrides |
| Color | hue/saturation/brightness jitter per stamp or per stroke, blend mode of the brush |
| Rendering | "vector" (today's outline brushes) or "stamped" (a stroke is many shape stamps) |

Stamped brushes need a renderer that draws many small stamps along the stroke. In lizard that's many small fills; cache per stroke, and it's a strong reason for the bytecode VM.

## 📦 Export / import format: `.lzbrush`

- A `.lzbrush` file is **plain lizard-script**: `let brush = {...}` with a header `-~> lzbrush 1`, plus optional bitmap data (shape / grain) as run-length lizard data, the same way `splash_art.liz` stores pixels. One file = one brush; a `.lzbrushes` file = a list (a brush set).
- Export: from the brush picker's ⋯ → "Export brush…", or the whole set.
- Import: drop a `.lzbrush` onto the window or the picker → it's validated in a step-limited, sandboxed run (a brush file is code; it may only produce data, no file access), then added to `~/.linkz/brushes/`.
- Because it's text, brushes are diffable, shareable in a gist, and an agent can write one.

## 🔁 Importing other apps' brushes (best effort)

| Source | What's inside | What we could take |
|---|---|---|
| Procreate `.brush` / `.brushset` | a zip: `Brush.archive` (a binary plist of settings), `Shape.png`, `Grain.png` | Shape and grain bitmaps reliably; the common settings (spacing, jitter, taper, size/opacity ranges) by mapping the plist keys we recognize. Undocumented format → keep a mapping table and warn about what was skipped. |
| Photoshop `.abr` | binary, versioned; sampled tips + dynamics | The brush tip bitmaps first; dynamics later if worth it. |
| Clip Studio `.sut` | SQLite database | Tip/texture images and basic size/opacity settings. |

Decoding archives and images is OS/format work, so it lives in the host as small converters that output a `.lzbrush` (lizard data); everything that *draws* stays lizard.

## 🛠️ Brush creator (UI)

A panel (artist and coder mode) with the groups above as tabs, a **scribble pad** that redraws a test stroke live as you tune, the pressure curve editor (shared with the B009/R002 pen settings), shape/grain pickers (procedural presets + import an image), and "Save brush" / "Export". In coder mode it shows the brush hash it writes. Designed in Figma first, like everything else.

## 🧭 Order when we get there

1. Stamped renderer + the extra parameters in `pens.liz` (headless, tested).
2. `.lzbrush` export/import (sandboxed).
3. Brush creator UI (Figma → build).
4. Procreate importer, then ABR / SUT if wanted.
