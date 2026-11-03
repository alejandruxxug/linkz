# B026 – Pen settings pop-up, test pad and stabilizer (v1.1, wave 1)

Design: R002 18b / 18b-2 / 19b / 18c and R003-4 (brush with the pulled-string stabilizer); Style Guide rows "Pen settings pop-up", "Pen picker ⚙", "String stabilizer". Rules: `docs/Batches/_v11-rules.md`.

## Build
- **⚙ on every pen row** of the pen picker (and double-click the brush tool) opens the **floating pen settings window**: draggable title bar, Reset, ✕; pressure curve editor (drag points, add/remove, presets soft/linear/firm, live `sky` line at the current pressure), width/opacity ranges, taper, smoothing, speed, Trackpad/Tablet/Mouse overrides; coder mode shows the pen hash it writes (scrolling, edited line highlighted).
- **Test pad**: draws with the current settings live (the real brush code), clear button, "draw here to try" when empty, live readout `pressure 0.62 · trackpad`; pad strokes never reach the document.
- **Save as new pen** → `~/.linkz/pens/<name>.liz` (a `let pen = {…}` file), listed in the picker after the built-ins; editing a built-in saves a copy, never the built-in.
- **Stabilizer** per pen (pop-up and brush/pencil inspector): Off / String (pulled string: the tip follows the cursor at a string length in px; ink only at the tip; `sky` string line and dotted radius circle drawn live) / Average (moving average over N samples); "catch up at the end" finishes at the cursor on release. Applied live while drawing. Expose it as `stabilize(state, sample, settings)` in its own file so the lasso (B027) reuses it.
