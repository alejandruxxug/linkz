# B055d – Alt isolation of the hovered layer (R004)

Model: Sonnet. Rules: `_v11-rules.md`. Read only: R004 frames 2 and 4, `lib/ui/view_canvas.liz`, the card file from B055c.

While the layer card is open and Alt is held, the canvas dims to 70% `void` and the hovered layer is drawn on top of it (hidden layers: no isolation). Releasing Alt or leaving the row restores the canvas. Partial repaints only (no full repaint per hover). Rig test.
