# B002 – Floats, math, units

Milestone: M1 Language basics · Repo: `../lizard-script` · Lane: Language · Needs: B001 (dot access)

## Goal
lizard-script has decimal numbers, a `math` library reached with dot access, and angle units, so curves and rotations can be written.

## Rules (additive only: every existing program behaves the same; `1 / 2` is still `0`)

### FLOAT
- Literal: digits `.` digits (`1.5`, `0.25`, `10.0`). Not `.5`, not `1.` (`1.` followed by a name is dot access on an integer, which is just an error at runtime).
- Mixed INTEGER and FLOAT → FLOAT. INTEGER with INTEGER → INTEGER exactly as now.
- `/` with a float on either side is real division. `//` with a float → floor, as FLOAT. `%` with a float → the remainder that pairs with `//` (`a == (a // b) * b + a % b`).
- Division or `%` by `0` or `0.0` is the existing division-by-zero error. Never produce Inf or NaN: any result that would be Inf/NaN is a runtime error (`the result is not a number`).
- `**`: if either side is FLOAT → `math.Pow`; `2 ** -1` (both integers) stays the existing error, `2.0 ** -1` is `0.5`.
- Comparisons `< > <= >=` work across INTEGER and FLOAT by value.
- `==` (loose): `1 == 1.0` is true; numeric strings convert to numbers as today (`"1.5" == 1.5`). `===`: `1 === 1.0` is false (different types).
- Prefix `-` works on floats. Truthiness: `0.0` is false.
- Printing (`puts`, `str`): shortest form that reads back as the same float, always with a `.` or exponent so it never looks like an integer: `1.0`, `0.5`, `0.1`, `1e+21`.
- Builtins: `int(1.9)` → `1` (toward zero), `int("1.5")` stays an error; new `float(x)` from integer, float or numeric string; `type(1.5)` → `"FLOAT"`.
- FLOAT can't be a hash key (same error as other unsupported keys).

### `math` (a builtin value, used with dot access)
- `math` is a builtin HASH (or equivalent) whose values are builtins: `math.sin`, `cos`, `tan`, `atan2(y, x)`, `sqrt`, `abs`, `floor`, `ceil`, `round` (half away from zero; `floor`/`ceil`/`round` return INTEGER), `min(a, b, ...)`, `max(a, b, ...)`, `random()` (FLOAT in [0, 1)), `random(a, b)` (INTEGER in [a, b]), `seed(n)`; plus `math.pi`, `math.tau`.
- Angles are radians.
- `random` is **deterministic**: a fixed default seed, reset at the start of every run (same program → same picture), changed with `math.seed(n)`. Expose an exported reset for hosts alongside the step reset from B001.
- Errors: `math.sqrt(-1)` → runtime error; wrong argument types/counts → the same style as other builtins.
- A user variable named `math` hides it, like other builtins.

### Units (lexer)
- A number immediately followed by a unit: `45deg`, `1.5rad`, `0.25turn`, `10px`.
- `deg` → radians (`45deg` is `0.785…` FLOAT), `rad` → unchanged (FLOAT), `turn` → `x * 2π` (FLOAT), `px` → the number unchanged (keeps INTEGER if it was one).
- A number directly followed by any other letters is a lexical error: `unknown unit "kg"`. (Today `45deg` is a syntax error, so this is additive.)

## Tests
Table tests in lexer/parser/evaluator for every rule above, including the edge cases (mixed types, `//` and `%` with negatives and floats, Inf/NaN guards, printing, `math.seed` determinism, units). Programs: `tests/valid/floats.liz`, `tests/valid/math_and_units.liz`, `tests/errors/runtime_not_a_number.liz`, `tests/errors/lexical_unknown_unit.liz` (+ `.expected`). Add `examples/10_floats_and_math.liz`.

## Docs
GRAMMAR.md (FLOAT, units), SEMANTICS.md (numbers, `math`, units), DESIGN.md §3 "Integer arithmetic only" stays true for integers but needs a sentence on floats, §4 remove "Integers only", §7 rows.

## Done when
`go test ./...` green, pre-existing `.expected` files untouched, `gofmt -l .` empty. Never mention Linkz in that repo.
