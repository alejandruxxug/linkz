# B003 – Colors and points

Milestone: M1 Language basics · Repo: `../lizard-script` · Lane: Language · Needs: B002 (floats, dot access)

## Goal
Colors and points are values the language understands, so drawing code reads like a drawing: `circle((120, 80), 30, #3fa34d)`.

## Rules (additive only: every construct here is an error today)

### COLOR
- Literals: `#rgb`, `#rrggbb`, `#rrggbbaa` (hex, any case). Other lengths or non-hex digits → lexical error (`bad color #12`). `#` is illegal today, so this is additive.
- Value: r, g, b, a as integers 0–255 (a = 255 when not written).
- Printing: `#rrggbb` lowercase, or `#rrggbbaa` when a ≠ 255. `type` → `"COLOR"`.
- Dot access reads channels: `c.r`, `c.g`, `c.b`, `c.a` (INTEGER). Colors can't be changed (assigning `c.r = 1` is an error, like strings).
- `==`/`===` compare all four channels. Truthy. Can be a hash key? No (same unsupported-key error).
- Builtins: `rgb(r, g, b)`, `rgba(r, g, b, a)` (each 0–255; integers, or floats rounded; out of range → error), `hsl(h, s, l)` (h in degrees 0–360, s and l 0–100) → COLOR.

### POINT
- Literal: `(x, y)` with two expressions. `(x)` stays a grouped expression; `()` and `(a, b, c)` are syntax errors (`a point has two coordinates`). Inside call arguments, `f((1, 2))` passes one point.
- Coordinates are numbers (INTEGER or FLOAT) and keep their type: `(1, 2) + (3, 4)` is `(4, 6)` with integers.
- Arithmetic: `p + q`, `p - q`, `p * n`, `n * p`, `p / n` (numeric rules from B002 per coordinate, so `(3, 4) / 2` is `(1, 2)`), `-p`.
- Dot access: `p.x`, `p.y` (read-only, like colors). `==`/`===` per coordinate with the usual numeric rules.
- Printing: `(1, 2)`, `(0.5, 2)`. `type` → `"POINT"`.
- Builtins: `point(x, y)`; a `math.dist(p, q)` and `math.angle(p, q)` (radians) for convenience.

## Tests
Lexer, parser and evaluator tables for every rule, including errors and printing. Programs: `tests/valid/colors.liz`, `tests/valid/points.liz`, `tests/errors/lexical_bad_color.liz`, `tests/errors/syntax_three_coordinates.liz` (+ `.expected`). Add `examples/11_colors_and_points.liz`.

## Docs
GRAMMAR.md, SEMANTICS.md, DESIGN.md §7 rows.

## Done when
`go test ./...` green, pre-existing `.expected` files untouched, `gofmt -l .` empty. Never mention Linkz in that repo.
