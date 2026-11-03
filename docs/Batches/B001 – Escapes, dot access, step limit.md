# B001 – Escapes, dot access, step limit

Milestone: M1 Language basics · Repo: `../lizard-script` · Lane: Language

## Goal
Strings can hold `"` and line breaks, `a.b` reads and writes hash keys, and a host can stop a program that runs too long.

## Rules (lizard-script CLAUDE.md applies: additive only, existing tests and `.expected` files unchanged)

### 1. String escapes (lexer)
- Inside a string: `\"` → `"`, `\\` → `\`, `\n` → line break, `\t` → tab.
- Any other `\x` is a **lexical error** (an ILLEGAL token, reported like the existing ones), message names the escape: `unknown escape \q in string`.
- A `\"` doesn't end the string, so `"say \"hi\""` is one string.
- An unterminated string ending in `\` is still the existing unterminated-string error.
- Note for the extensions log: this changes what a backslash means inside a string. No existing program or test contains one (checked 2026-11-02).

### 2. Dot access (parser only, desugars to index)
- `a.b` is exactly `a["b"]`: same AST node (IndexExpression with a string literal key), same evaluator path. Precedence and associativity of `[` (so `a.b.c`, `a.b[0]`, `a.b(1)`, `f().x` work).
- `a.b = x` is exactly `a["b"] = x`.
- After `.` must come an identifier; a keyword or anything else is a syntax error with a clear message (`expected a name after "."`).
- The token: add `DOT`. A `.` that isn't followed by an identifier is a syntax error, not a lexical one. (B002 adds floats like `1.5`; don't add them here, but make sure the lexer change won't block it.)
- `String()` of the AST may print `a["b"]` — fine.

### 3. Step limit (evaluator)
- An exported knob, e.g. `evaluator.MaxSteps` (int, 0 = no limit, the default so the REPL and existing tests are unchanged) and the counter reset per run (an exported `evaluator.ResetSteps()` or a `Run`-level reset; pick the cleanest in the repo's style).
- One step per loop iteration and per function call (builtins included).
- Over the limit → runtime error at the loop or call: `[line L, column C] runtime error: the program ran too long (more than N steps)`.
- The host (an editor that re-runs on every keystroke) sets it; it must make `while (true) {}` stop.

## Tests (the repo's style)
- Lexer table tests for every escape and the errors.
- Parser tests: `a.b`, `a.b.c`, `a.b[0]`, `a.b(1)`, `f().x`, `a.b = 1`, `a.if` / `a.` / `a.5` errors.
- Evaluator tests: dot reads/writes on hashes, missing key gives `null` like `a["b"]`, dot on a non-hash is the same error as indexing it; step limit on `while (true) {}`, on infinite recursion below the stack limit, and unlimited when 0.
- Programs: `tests/valid/escapes_and_dots.liz` (+ `.expected`), `tests/errors/lexical_unknown_escape.liz`, `tests/errors/syntax_dot_without_name.liz` (+ `.expected`).

## Docs (same change)
- GRAMMAR.md: escapes, the `.` rule and its precedence.
- SEMANTICS.md: escapes, dot access, the step limit (for hosts).
- DESIGN.md §4: remove "Strings have no escape sequences". §7: one row per extension (date 2026-11-02). §6 Use of AI: add a line that the extensions were written with Claude Code.

## Done when
`go test ./...` green in lizard-script with every pre-existing `.expected` file untouched; `gofmt -l .` empty. Never mention Linkz in that repo.
