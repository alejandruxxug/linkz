# B055b – Selecting a layer highlights its code block

Model: Sonnet. Rules: `_v11-rules.md`. Read only: `lib/ui/view_panels.liz` (row paint), `select_model.liz`, `artist_layers.liz`, `view_canvas.liz` (`view_show_line`).

> "If we select a layer then it should highlight the area in the code too." — the developer

Selecting a layer (layers panel, its `layer` row, or the canvas) highlights every code row of that layer — header through its last indented row — with a light band (row tone + 2 px `lizard` bar, lighter than the selected-row look so the caret row still stands out) and scrolls the code (or artist's open drawer) so the header shows. Nested layers highlight only their own block. Rig tests in both modes.
