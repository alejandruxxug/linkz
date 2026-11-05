# Linkz – Tracing

Turns a picture (PNG or JPEG) into **editable lizard-script**: a `.lzdraw` with one layer per color and one `path(...)` per region. The result is code, not an SVG blob, so it can be recolored, moved and edited like anything else in Linkz. The algorithm is our own (potrace is GPL); it is plain lizard-script in `lib/trace/trace.liz`.

## Use

```bash
linkz trace photo.png out.lzdraw [--colors 8] [--detail 1.5] [--min-area 8] [--max-side 256] [--scale 1] [--corner 55]
linkz trace sketch.png out.lzdraw --mode lineart [--line-threshold 0.5] [--line-min 6] [--continuation 35] [--fills true|false] [--solid 0] [--focus x,y,w,h[,f]]
linkz render out.lzdraw            # writes out.svg
```

From lizard (needs the `img` primitive and `time`; the CLI and `host.Trace` provide them):

```
let r = trace("photo.png", {"colors": 6, "detail": 2})
r.code    -- the .lzdraw source
r.stats   -- {colors, regions, holes, points, merged, ms, size, stages}
```

`trace_pixels(size, pixels, options, name)` does the same on pixels already in memory (`size` a point, `pixels` a flat list of colors, row by row), which is how the tests build tiny pictures.

## Options (`trace_options`, with tooltips in the source)

| Option | Default | Meaning |
|---|---|---|
| `colors` | 8 | Palette size, 2 to 64. More colors, more shapes. |
| `detail` | 1.5 | Simplification tolerance in pixels. Small keeps wiggles; large gives few, smooth shapes. |
| `min_area` | 8 | Regions under this many pixels are merged into the neighbour they touch most. 0 keeps every speck. |
| `corner` | 55 | Turn in degrees from which an outline point stays a sharp corner. |
| `smooth` | true | false writes straight-sided shapes. |
| `background` | true | The biggest color becomes one full-canvas `rect` behind everything. |
| `refine` | 2 | k-means rounds after the median cut. |
| `scale` | 1 | Multiplies every coordinate. |
| `max_side` | 256 | The picture is shrunk until its longer side is at most this. |

Unknown option names are an error, so a typo doesn't pass silently.

## How it works

1. **Load** (`host/img.go`, Go, format and size work only): decode, box-filter down to `max_side` (never enlarge), flatten transparency onto white. `img.load(path, max_side)` gives `{"size": (w, h), "pixels": [COLOR, ...]}`. The pixels are one flat row-major array (`pixels[y * w + x]`): one index per read, no inner array to fetch.
2. **Reduce colors** (`trace_quantize`): pixels are grouped by their top 5 bits per channel (so the palette work is done on at most a few thousand distinct colors, not every pixel), cut into `colors` boxes by **median cut** (widest channel, split at the pixel-count median), then polished with `refine` **k-means** rounds. Nothing is random: the same picture gives the same palette.
3. **Clean** (`trace_label`, `trace_clean`): connected pieces of one color (sides touching, not corners) smaller than `min_area` take the color of the neighbour they touch most.
4. **Outlines** (`trace_rings`): every pixel side that faces another color is an edge, directed so the pixel is on its right; edges are chained corner to corner. This is marching squares done on the pixel-corner lattice, so it needs one pass over the picture and only touches the boundary. Outer outlines run clockwise on screen, holes counter-clockwise. A corner where two same-color pixels only touch diagonally always turns right, so they stay separate regions (the same rule as the 4-connected pieces).
5. **Simplify** (`trace_ring_simplify`): Ramer-Douglas-Peucker on each closed outline (`simplify` from `geometry.liz`), tolerance `detail`.
6. **Write code** (`trace_code`): the `let colors = {...}` palette, `canvas`, then `layer traced` with one `layer color_N` per color.

## What the code looks like

```
let colors = {
  "c1": #f0f0e6,
  "c2": #c82828
}
canvas = {"size": (100, 100), "background": colors.c1}

layer traced
  layer color_2
    path([(45, 20), (55, 20), ...], colors.c2)
  layer color_1
    rect((0, 0), (100, 100), colors.c1)
```

- `c1` is the biggest color, `c2` the next, and so on. Every shape uses `colors.cN`, so recoloring is one edit.
- In Linkz a layer lower in the file is *behind*, so the layers are written smallest area first and the biggest last: the biggest color sits at the back.
- A region with holes is one `path(points, fill, true, true, [hole1, hole2])`.
- The smoothing is done by `path()` (new in `lib/draw/shapes.liz`): a Catmull-Rom curve through the points as cubic Béziers, except that a point where the line turns more than `corner` degrees stays sharp, and the run between two sharp points is straight. Edit a point or the `corner` argument and the shape follows.

## Speed (a tree-walker, measured on this machine)

Lilies reference photo (708 x 416 PNG), `max_side` 160 (160 x 94, 15 040 pixels), 8 colors, default options:

| Stage | ms |
|---|---|
| load and shrink (Go) | 15 |
| quantize (histogram, median cut, k-means) | 270 |
| clean (labeling, merge specks) | 190 |
| outlines | 85 |
| simplify | 75 |
| write code | 6 |
| **total** | **~640** |

Result: 216 regions, 1 727 points, 20.7 KB of code. At `max_side` 256 (8 colors): 403 regions, 3 717 points, about 1.25 s. Time grows roughly with the pixel count; the stage that dominates is the per-pixel color work, then the labeling.

## Limits

- The outline follows pixel edges. The simplified shapes pass through pixel corners, so curves are only as good as the picture's resolution: a traced circle of radius 30 comes out about 3% larger than the true circle, and a photo at 160 px looks faceted. A bigger `max_side` and a smaller `detail` give smoother shapes and many more points.
- Photos reduced to 8 colors look posterized, with many small shapes; this is a vectorizer for flat art, logos and sketches first, photos second. Gradients become bands.
- Neighbouring regions are separate shapes, so a faint hairline can show along shared edges (anti-aliasing). The biggest color is a plain rectangle (`background`) to avoid the most visible seams.
- Holes are not editable anchors in the editor (only the outer outline's points are anchors of a `path`).
- One palette distance: plain RGB. No dithering, no edge detection, no stroke (centerline) tracing.

## Line art mode

`--mode lineart` (`{"mode": "lineart"}` from lizard) traces inked drawings (pen sketches, comics, logos made of lines) as brush strokes plus fills for the areas between them, instead of color regions. The default mode is `photo`, the tracer described above. The design and the algorithm are in `Linkz – Tracing v2.md`.

| Flag | Option | Default | What it does |
|---|---|---|---|
| `--line-threshold x` | `line_threshold` | 0.5 | How much darker than its surroundings a pixel must be to start a line, 0 to 1. Lower finds fainter lines. |
| `--line-min n` | `line_min` | 6 | Strokes shorter than this many pixels are dropped. |
| `--continuation deg` | `continuation` | 35 | The largest turn a line may make at a junction and still go on as one stroke. |
| `--radius n` | `radius` | 0 | Half side of the neighbourhood a pixel is compared with; 0 picks it from the picture size. |
| `--fills true\|false` | `fills` | true | false traces only the lines. |
| `--solid n` | `solid` | 0 | Ink this far (px) from the paper all round is a solid fill, not a stroke; 0 picks it from the size. |
| `--focus x,y,w,h[,f]` | `focus` | none | A box of the picture (in its own pixels) with more detail (`f` above 1, 2 if left out) or less (`f` below 1); repeat the flag for more boxes. |

Focus boxes (Tracing v3 step 5b): from lizard, `{"focus": [{"box": [x, y, w, h], "detail": f}, ...]}`. `f` is kept between 0.25 and 4. Inside a box the smallest region a fill keeps is 1/f squared of what it would be, and the fit tolerance of its outlines 1/f (a later box wins where two meet). It works with `adaptive` off too. The Trace image dialog draws the boxes with its Detail tool, and Insert as layer writes them in a `-~> trace focus [...]` comment above the layer.

The merge pass (`lineart_merge`, Tracing v3 step 1) runs on the fills before their outlines are made: `snap` (palette colors under 2 apart in Lab become one), `join` (touching regions of one color become one shape), `stack` (a region with islands is drawn whole with the islands on top, not with holes), `absorb` (a region of the color of the shape under it is given to that shape) and `cover` (shapes hidden by later shapes are dropped). From lizard: `{"merge": false}` for none, or `{"merge": {"stack": false}}` to turn single steps off (the others stay on); no CLI flag yet. The stats have `merge: {snapped, joined, stacked, absorbed, covered}`.

`--colors` (default 16 here), `--min-area`, `--detail` (0.8), `--scale` and `--max-side` (1024) also apply, with their line-art defaults. The line-art flags are an error in photo mode. The summary line then reads `N strokes, N fills, N tips, N junctions, N ms`, followed by the time of each stage (the `stages` of the stats).

## Next: line art

The anime/line-art pipeline (strokes with width, fills under the lines, sub-pixel Bézier edges, corners and hair spikes, primitive fitting) is designed in [[Linkz – Tracing v2]].
