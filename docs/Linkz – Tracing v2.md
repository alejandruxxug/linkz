# Linkz – Tracing v2: line art and smooth edges

> "I would like a method for vectorialization to detect lines and to adjust rects like that. For accurate vectorializations of anime art. Detect line flow for hair so we can have spiky edges and continuous ones too. This is the main focus: to detect edges and for smooth edges to be formed." — the developer, 2026-11-03

Status: **design, not scheduled** (build later, after the bytecode VM). Extends [[Linkz – Tracing]] (v1: posterize → pixel-edge outlines → RDP → `path`). v1 stays as the "photo" mode.

## 🎯 Goal

Anime/cel art traced into code an artist would have written:

- **Line art becomes strokes**, not filled slivers: one `brush(points, pressures, color, pens.ink)` per line, the width read from the image, so lines thicken and taper like real inking.
- **Flat colors become fills** that sit *under* the lines with no seams.
- **Edges are smooth curves**, except where the art has a real corner: hair spikes, chin, collar points stay sharp.
- **Simple shapes are recognized**: a rectangle, circle, ellipse or straight line becomes `rect`/`circle`/`ellipse`/`line`, not a 40-point path.
- **Hair follows its flow**: strands are continuous curves along the hair direction, tips end in sharp spikes.

## 🧭 Pipeline (mode `"lineart"`)

| # | Stage | What | Why |
|---|---|---|---|
| 1 | **Load + upscale** | Small images are upscaled 2× (bicubic, host) before analysis; work at up to 1024 px. | Line widths of 1–2 px need sub-pixel room. |
| 2 | **Edge-preserving smooth** | A small Kuwahara or bilateral pass. | Removes JPEG noise and screentone without blurring lines. |
| 3 | **Line mask** | "Lineness" per pixel: dark relative to its neighbourhood (difference of Gaussians, XDoG-style) and thin. Threshold with hysteresis (strong lines seed, weak ones only if connected) → a binary line mask. Colored lines (not only black) are kept as their own color. | Separates the inking from the colors. |
| 4 | **Orientation field** | Structure tensor on the smoothed image → a direction and a coherence at every pixel (how strongly the area "flows" one way). | The hair flow: tells strands apart and guides stages 5, 6 and 8. |
| 5 | **Centerlines** | Skeletonize the line mask (Zhang–Suen thinning), then a **distance transform** gives the half-width at every skeleton pixel. The skeleton becomes a graph: endpoints, junctions, paths between them. | Each path is one stroke, with its width profile. |
| 6 | **Stroke building** | At each junction, continue through the branch whose direction best matches the incoming one and the orientation field (good continuation), so a hair strand crossing another stays one stroke. Prune tiny spurs. Merge gaps < 2 px along the flow. | Continuous strands instead of fragments. |
| 7 | **Width → pressure** | The width profile along each stroke is smoothed and mapped to pressures 0–1 for a pen whose `min`/`max` are the thinnest and thickest widths found. Tips where the width falls to 0 are **spiky** (taper); tips that end at full width are **blunt** (no taper). | `brush()` reproduces the line weight and its taper. |
| 8 | **Fills under the lines** | Quantize the non-line pixels in Lab (fewer colors than photo mode, flat cel shading), label regions bounded by the line mask, then **grow each region into the line pixels** until they meet. Neighbouring regions share one boundary (built once, used by both). | No hairline seams: the fills meet under the ink, and the ink covers the join. |
| 9 | **Sub-pixel outlines** | Region outlines by marching squares on a slightly blurred region mask at the 0.5 level (not on pixel corners). | Real curves, not staircases. |
| 10 | **Corners** | Curvature along each outline at two scales; a point is a corner where the turn is sharp at both scales. Hair spikes: a narrow region whose two sides converge along the flow is forced to a sharp point. | Sharp where the art is sharp, smooth elsewhere. |
| 11 | **Curve fitting** | Between corners, fit cubic Béziers (Schneider's least-squares fitting, split where the error > `detail`). Strokes' centerlines are fitted the same way. | Few control points, smooth, editable. |
| 12 | **Primitive fitting** | Try each closed outline as an ellipse/circle (algebraic fit), a rectangle/rotated rect (min-area box), a polygon with few sides; each straight stroke as `line`. Accept if the max error < `snap` px. | Clean, readable code ("adjust rects like that"). |
| 13 | **Write code** | `let colors = {…}` (named palette), then `layer fills` (back) and `layer lines` (front), optional `layer shadows` for colors that are darker versions of a base color. | Recolor, re-ink or hide the lines with one edit. |

## ✍️ What the code looks like

```
let colors = {"skin": #fbe3d3, "hair": #3b2a5c, "hair_shade": #2a1d45, "ink": #1a1418}
pens.traced = {"name": "traced", "min": 0.6, "max": 3.2, "taper": 0, "opacity": [1, 1], "nib": null, "grain": 0}

layer fills
  path([(210, 88), …], colors.hair)
  ellipse((256, 230), (70, 88), colors.skin)
layer lines
  brush([(181, 60), (196, 92), (204, 140)], [0, 0.8, 0.3], colors.ink, pens.traced)
  line((140, 400), (372, 400), colors.ink, 2)
```

**Done.** `path` accepts Bézier items mixed in its points: `["c", c1, c2, p]` (cubic) and `["q", c, p]` (quadratic) draw a curve exactly from the previous point to `p` (flattened with `flatten_cubic`, no Catmull-Rom smoothing round them; SVG export writes real `C`/`Q`). Plain points keep their smoothing, and the points next to a curve are sharp ends. The tracer writes the fills' fitted curves this way (`path([(x, y), ["c", …], …], fill)`). Strokes stay points: `brush` takes one width and pressure per point, and following them along a curve would need widths at the curve ends and between them, which it has no way to say yet.

Measured on `wikipe-tan-casual.png` at 1024 (561×1024): fill outlines went from 6558 points to 3393 curves (about 10600 numbers' worth of points with their controls), the code from 134020 to 219676 bytes (fills 104 KB to 190 KB: three points per curve cost more than the simplified samples they replace), mean color error 5.48 to 5.46, ink IoU 0.832 to 0.832. The gain is exactness and editability (real curves, no resampling), not size.

## 🎛️ Options (new, with v1's)

| Option | Default | Meaning |
|---|---|---|
| `mode` | `"photo"` | `"lineart"` for anime/cel art (this pipeline); `"photo"` is v1. |
| `line_threshold` | 0.5 | How dark/contrasty a pixel must be to count as ink. |
| `line_min` | 6 | Strokes shorter than this (px) are dropped. |
| `continuation` | 35 | Max angle (deg) a stroke may turn at a junction to continue through it. |
| `flow` | 0.6 | How much the orientation field steers stroke continuation and spikes (0 = off). |
| `corner` | 55 | As v1, applied to curvature (two scales). |
| `detail` | 0.8 | Curve-fit tolerance in px (smaller than v1, since outlines are now sub-pixel). |
| `snap` | 1.5 | Max error (px) to replace an outline with `rect`/`circle`/`ellipse`/`line`; 0 = never. |
| `lines` | `"brush"` | `"brush"` (width as pressure), `"stroke"` (constant width), or `"fill"` (v1-style filled outlines). |

## 🧪 Test set and "done when"

- Synthetic tests: a circle, a rotated rectangle, a tapered stroke, a Y junction, two crossing strands, a spiky hair tuft (3 spikes), flat fills behind a black outline.
- Real tests: 3–4 anime frames/illustrations kept in `docs/Art/References/trace/` (with permission/licensing noted), scored by eye and by the **render difference**: render the traced code and compare with the original (mean color error, and line-mask overlap).
- Done when: the synthetic shapes come back as primitives / single strokes with spikes in the right place, the real images show continuous hair strands and no seams, and the code is under ~1 line per visible shape.

## ⚡ Performance

This is many per-pixel passes (smooth, DoG, structure tensor, thinning, distance transform) at up to 1024 px: ~1 M pixels × ~10 passes. v1 does about 0.2 s per 100 k pixel-pass-equivalents in the tree-walker, so lineart mode would take tens of seconds today. Following the project rule (slowness → fix the language), this waits for the **bytecode VM**; measure each stage, and if one pass is still too slow, add a general lizard feature (e.g. numeric arrays) rather than moving tracing into Go. The image *decode and resize* stay in the host, as in v1.

## 🖼️ UI

A trace dialog (Figma first): drop or open an image → mode Photo / Line art, the main sliders, a before/after preview with a split handle, toggles to show lines / fills, the stats (strokes, fills, primitives, ms), then "Open as drawing" or "Insert as layer".
