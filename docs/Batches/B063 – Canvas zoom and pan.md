# B063 – Zoom and pan the canvas

Model: Sonnet. Rules: `_v11-rules.md`. Read only: `lib/ui/view_layout.liz` (`view_zoom_steps`, the canvas layout, and the UI↔drawing point conversions near line 500–560), `lib/ui/view_canvas.liz`, the scroll and key handling in `lib/ui/view_input.liz`, `lib/gfx/tiles.liz`, the Style Guide's status-bar row, and `tests/editor/editor_bench.liz`.

> "Also allow people to zoom in into the canvas" — the developer, 2026-11-04

Today the zoom is always "fit" (the largest step that leaves room around the drawing). Add a user zoom and pan on top of it:

- **Zoom:** ⌘+ / ⌘− (Ctrl on Windows/Linux) step through the zoom steps, extended up to 3200% and down to 5%. ⌘0 = fit (today's behavior), ⌘1 = 100%. ⌘/Ctrl + scroll zooms smoothly **around the cursor**: the drawing point under the cursor stays put.
- **Pan:** two-finger scroll or wheel pans (Shift = horizontal); Space-drag and middle-drag pan with a grab cursor. The pan is clamped so the drawing can't be lost: at least 32 px stay visible.
- Fit comes back on its own only when the user hasn't zoomed. A window resize keeps the user's zoom, and the zoom/pan are per open document (not saved in the file).
- **Status bar:** the zoom % (already in the Style Guide) is live; clicking it opens a small menu with Fit, 50%, 100%, 200% and 400%. en/es strings.
- **Everything that maps points keeps working** at any zoom and pan: drawing (pencil, brush, lasso, select/move, anchors, hit tests), selection highlights, the layer card and Alt isolation, stroke stabilization. Go through every user of the conversions and the `cl.zoom` readers (`grep -n zoom lib/ui`).
- **Speed:** at high zoom, never make a canvas buffer bigger than the visible canvas area. Render only the visible part (the tiles and spatial index exist), and stroke widths scale with the zoom. The M10 bench must not regress more than 5% at fit; add bench cases at 400% and 1600% while panning, with a target under 16 ms a frame at 2k shapes.
- Rig tests: zoom around the cursor keeps the point, ⌘0/⌘1, pan clamping, drawing a stroke at 400% plus a pan lands at the right drawing coordinates, and a hit test at 1600%. Drive the real app (isolated HOME) and save shots to `docs/Art/Checks/v1.2/zoom/`.

Trackpad pinch is not in this batch: the Terrarium has no pinch event yet (that's a separate lizard-script task). Make zooming one function, e.g. `canvas_zoom_at(factor, ui_point)`, so the pinch can call it later.
