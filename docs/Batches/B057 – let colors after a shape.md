# B057 – Fix: `let colors` declared after a shape that uses it

Model: Sonnet. Rules: `_v11-rules.md`. Read only: `lib/ui/color_vars.liz`, `color_target.liz`, `editor_twoway.liz` (`twoway_after_decls`).

When `let colors` comes after a shape that uses `colors.x`, the run errors and the palette lists `size` and `background` (the canvas hash) as variables. Fix the variable lookup to find the real `let colors` hash only; decide and document what the palette shows while the drawing errors (the last good variables). Regression test.
