# B044 – `lizardscript fmt` and `lizardscript check`

Repo: `../lizard-script` · Lane: Go (new packages `format/`, `check/`, `main.go` subcommands) · Parallel with B043 (VM), B045 (parser errors), B031/B032 (lily2d/3d)

## fmt
- `lizardscript fmt [-w] [--check] files/dirs…`: one canonical layout, like gofmt. Prints the formatted code (or rewrites with `-w`; `--check` lists unformatted files and exits 1).
- Rules (document them in docs/STYLE.md, with the guide's style): 2-space indentation, one space around binary operators and after commas, `{` on the same line, `} else {`/`} elif (…) {`, blank lines kept but collapsed to one, trailing `-~>` comments aligned within a block, long lines not reflowed (v1), LEP 0001 respected (operators at the end of a line when a statement continues), hash literals one pair per line when they span lines, no trailing whitespace, final newline.
- Built on `lexer` tokens (`KeepComments`) and `ast` spans; never changes meaning: a test formats every `.liz` in the repo (tests/, examples/, terrarium/lib/, terrarium/std/ when present, bench/, docs/learn code blocks) and checks the program parses to the same AST and prints the same output; formatting twice is a no-op (idempotence); round-trip on a fuzz corpus.
- Don't reformat the repo's own files in this batch; report how many would change (a later commit can run `fmt -w`).

## check
- `lizardscript check files/dirs…` reports, without running, with `[line L, column C]` and a severity:
  - errors: syntax errors (all of them), a name that's never defined anywhere it could come from (scopes, modules via `use`, builtins, Terrarium modules, standard modules), `break`/`continue`/`return` misuse the parser doesn't catch, calling a known function with the wrong number of arguments, assigning to an undeclared name;
  - warnings: unused `let`/parameters (`_name` exempt), a `let` that shadows an outer name, unreachable code after `return`/`break`, `==` where both sides are literals, a `use` whose names are never used, self-assignment;
  - hints: "did you mean `len`?" (edit distance over the names in scope) for unknown names and unknown hash/module members of known modules (`math.sine` → `math.sin`).
- Also as a Go API (`check.File(path)`, `check.Source(name, text)` → `[]Diagnostic{Severity, Line, Column, EndLine, EndColumn, Message, Hint}`) and as a Terrarium function `lizard.check(code, name)` so editors written in lizard (Linkz) can show diagnostics.
- Must respect lizard's late-bound scoping (a name may be defined later in an enclosing scope — read compiler/scopes.go's explanation) and modules; no false positives on any program in the repo (tests/valid, examples, terrarium/lib, docs/learn blocks): a test runs check over all of them and allows only intended warnings, listed.
- `--json` output for tools (B046 builds the language server on it).

## Done when
Tests (table-driven per rule, the repo-wide no-false-positive test, fmt idempotence + meaning-preservation over the repo and the guide), docs (docs/STYLE.md, TERRARIUM.md for `lizard.check`, guide chapter 18 "tools" section en + es), main.go subcommands. Clean code. Don't commit. Ownership: B043 owns `code/`, `compiler/`, `vm/`, `object/`; B045 owns `parser/` error recovery; lily builders own `terrarium/std/` and `evaluator/modules.go`.
