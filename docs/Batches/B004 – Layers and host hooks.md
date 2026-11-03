# B004 – Layers and host hooks

Milestone: M1 Language basics · Repo: `../lizard-script` · Lane: Language · Needs: B003

## Goal
`layer name` opens an indented block (Python style), and a host can observe which statement and which layer is running, so it can tag what a program makes with where in the source it came from.

## Rules (additive only: `layer name` on one line is a syntax error today; a variable named `layer` keeps working)

### Layer statement
- A statement that starts with the identifier `layer` followed **on the same line** by a name (an identifier, or a string for names with spaces) and then **nothing else on that line** (a trailing `-~>` comment is fine) is a layer statement.
- Its body is every following statement whose **first token's column is greater** than the column of `layer`, until the first statement at that column or less (or the end of the enclosing block / file). Use the token positions the parser already has; no INDENT/DEDENT tokens.
- An empty body is allowed (a layer with nothing yet).
- Layers nest: a `layer` inside a layer body is a child layer.
- `layer` anywhere else is the ordinary identifier: `let layer = 1`, `layer + 1`, `layer = 2`, `f(layer)` all behave exactly as now.
- Tabs: a tab counts as one column like today's lexer; mixing is the user's problem (document it). The editor writes 2 spaces.
- Inside a `{ }` block (function body, loop), layer statements work the same way, measured against the `layer` keyword's column.
- AST: a new `LayerStatement{Token, Name string, Body *BlockStatement}`; `String()` prints it back as `layer name` + indented body.
- Evaluation: runs the body in a new enclosed environment (like other blocks); value is `null`. `return`/`break`/errors pass through like any block.

### Host hooks (evaluator, exported, nil by default, zero cost when nil)
- `evaluator.OnStatement func(node ast.Statement)`: called right before every statement is evaluated (program, blocks, layer bodies, function bodies).
- `evaluator.OnLayer func(name string, enter bool, tok token.Token)`: called when a layer body starts (`enter = true`) and when it ends, also when it ends by error/return/break.
- The REPL and `Run` leave both nil, so nothing changes for them.

### Text and list helpers (builtins, additive)
An editor written in lizard needs to cut and join text. Two builtin hashes like `math`:
- `text.slice(s, start, end)` (character indexes, clamped; negative counts from the end), `text.find(s, sub, [from])` (index or -1), `text.split(s, sep)`, `text.join(list, sep)`, `text.replace(s, old, new)`, `text.trim(s)`, `text.lines(s)` (keeps empty lines; `\r\n` → lines), `text.starts(s, prefix)`, `text.ends(s, suffix)`, `text.upper`, `text.lower`, `text.code(s, i)` (character code) and `text.char(code)`, `text.number(s)` (INTEGER/FLOAT or `null` if not a number, never an error).
- `list.slice(a, start, end)` (new array), `list.insert(a, i, x)` and `list.remove(a, i)` (change `a` in place and return it, like `push`), `list.index(a, x)` (-1 if missing, uses `===`), `list.reverse(a)` (new array), `list.sort(a, [less])` (new array; numbers or strings, or a lizard function `less(x, y)`), `list.copy(a)`.
- Wrong types/counts → the usual builtin errors. A user variable named `text` or `list` hides them.

## Tests
- Text and list helpers: a table test per function, including unicode characters (indexes count characters), empty strings, out-of-range indexes, `sort` with a lizard comparator.
- Parser: layer with body, empty layer, nested layers, layer inside a function and a loop, the body ending at a dedent, a string name, a trailing comment, `layer` used as a variable in every position above (unchanged trees), `layer sun foo` (syntax error with a clear message).
- Evaluator: variables inside a layer are scoped to it; outer variables are visible and assignable; hooks fire in the right order with the right names (enter/exit, including on error); hooks nil by default.
- Programs: `tests/valid/layers.liz`, `tests/errors/syntax_layer_extra_tokens.liz` (+ `.expected`). Example `examples/12_layers.liz`.

## Docs
GRAMMAR.md (the layer rule and how indentation is measured), SEMANTICS.md (layers and the two hooks, for hosts), DESIGN.md §3 a sentence on mixing indentation for layers with braces for code, §7 rows.

## Done when
`go test ./...` fully green with every pre-existing test and `.expected` file unchanged; `gofmt -l .` empty. Never mention Linkz in that repo.
