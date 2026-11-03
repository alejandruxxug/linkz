# B029 – Bytecode compiler + VM (M6)

Milestone: M6 Speed · Repo: `../lizard-script` (language) then `linkz` (switch the host) · Lane: Go, in lizard-script · Needs: nothing

> "continue building today bytecode VM lands" — the developer, 2026-11-03

## Goal

lizard-script gets a second way to run a program: **compile the AST to bytecode and run it on a stack VM**, several times faster than the tree-walker, with **exactly the same behaviour**. The tree-walker stays, untouched in meaning (CLAUDE.md: the language grows by adding; the Go foundation stays). Linkz then runs on the VM.

## Rules (from lizard-script `CLAUDE.md`, read it first)

- Purely additive. `evaluator.Eval` keeps working and stays the default for `repl.Run` and the REPL. No existing test or `.expected` file changes.
- New packages, in the style of the repo (plain-word comments, small functions, errors with `[line L, column C]`): `code` (opcodes + encoding), `compiler` (AST → bytecode), `vm` (runs it). The course follows *Writing an Interpreter in Go*; this follows its sequel's shape (*Writing a Compiler in Go*) so the team recognizes it, but lizard's semantics win wherever they differ from Monkey's.
- Never mention Linkz in lizard-script (commits, docs, comments). Describe it as a language feature: "a faster engine for hosts that run a lot of code".
- DESIGN.md: a row in §7 (*Extensions beyond the course design*) and a short section on the second engine (architecture, what is resolved at compile time, the interop points). SEMANTICS.md: a note that both engines give the same results. GRAMMAR.md: no change.
- Don't commit; the orchestrator commits.

## Entry points

```go
vm.Run(program *ast.Program, env *object.Environment) object.Object
```

Same contract as `evaluator.Eval(program, env)`: the program's value, or an `*object.Error`. `env` is the same `*object.Environment` the hosts already fill with their builtins and globals, and it must see the program's top-level `let`s afterwards (hosts read them with `env.Get`, and run later programs in the same env). It honours `evaluator.MaxSteps` / `Steps` / `ResetSteps` / `SetSteps`, `evaluator.Output`, `evaluator.OnStatement`, `evaluator.OnLayer` (one counter, one set of hooks, shared by both engines).

Also: `lizardscript --vm program.liz` in `main.go`, and `repl.RunVM` (or an option) for the tests.

## Behaviour that must be identical (the hard part)

Read `evaluator/evaluator.go`, `object/object.go`, `docs/SEMANTICS.md`. In particular:

1. **Scopes are dynamic environments today.** Each block, each loop round, a for-loop's init, each call, each catch and each layer body gets a new scope; `let` defines in the current scope at run time; a name is looked up innermost-out *at the moment it is read*. Consequences the VM must keep:
   - `let y = x; let x = 2` in a scope reads the **outer** `x` first, then the new one.
   - A closure created before a later `let x` in an enclosing scope sees that `x` once it exists.
   - Closures capture variables **by reference** (assignment through `=` is seen by every closure sharing the scope); each loop round has a fresh scope, so closures made in a loop capture that round's variables.
   - `let` twice in one scope redefines; `=` on an undeclared name is the existing error.
   - Names not found in scopes fall back, in this order, to the builtins, `null`, `math`, `text`, `list` (a user `let len = …` shadows the builtin). Host globals live in the `env` map and may be added **after** compilation (between runs) — globals must be looked up by name (cached is fine if invalidation is correct).
   - Suggested approach: resolve each scope's `let` names at compile time into slots (cells for captured ones), and where a read might happen before its slot is defined, emit a read that falls back to the next candidate when the slot is still unset. Prove it with tests for every case above.
2. **Values of statements**: a block's value is its last statement's value (`if` is an expression; a function without `return` gives its last statement's value; loops give `null`; `let` gives `null`); a `return` at top level ends the program with that value.
3. **Errors**: same messages and positions (the same tokens), the same `Uncatchable` rules, `try/catch` hash (`message`, `line`, `column`), builtin errors wrapped as `name: message` except `error()`, errors from functions a builtin called keep `Uncatchable`.
4. **Steps**: one step per loop round and per call (builtins included), counted at the same moments, so the step-limit error happens on the same token.
5. **Limits**: `maxCallDepth` 10 000 with the same message; the VM has no Go recursion for lizard calls, so "nest too deeply" may not occur on the VM — if any test program depends on it, report it (don't change the test).
6. **Hooks**: `OnStatement` gets the same `ast.Statement` before every statement (program, blocks, layer bodies, function bodies) — keep the statement nodes in the constant pool and skip the call when the hook is nil. `OnLayer` enter/exit exactly as now, including exits by error, return and break.
7. **Functions are values**: `type(f)`, `puts(f)` (prints the source), `==`/`===` (identity), wrong-argument-count messages, default parameters evaluated at each call in the call's scope, recursion, closures, functions stored in hashes and called with dot access.
8. **Interop between engines**, both ways (additive extension points):
   - The VM's closures report `FUNCTION_OBJ` and inspect like `*object.Function`.
   - Builtins that call functions (`list.sort` with a comparison, any other `applyFunction` caller) must work with VM closures: add an additive hook in `evaluator` (e.g. an `object.Callable` interface or `evaluator.CallFunction` variable used before "not a function").
   - The VM must call a tree-walker `*object.Function` found in `env` (a host may have defined it with `Eval`): export an additive `evaluator.ApplyFunction`.
   - Programs run with `Eval` after `vm.Run` in the same env (and the reverse) keep working.

## Tests

- Every program in `tests/valid`, `tests/errors` and `examples/` runs on **both** engines with byte-identical output; add `TestValidProgramsVM` / `TestErrorProgramsVM` next to the existing ones (new tests, old ones untouched).
- Table-driven unit tests for `code`, `compiler`, `vm`, including every scope case in rule 1, closures in loops, try/catch, layers + hooks (record the hook calls on both engines and compare), step-limit position, interop both ways.
- A **differential test**: many generated or hand-written small programs run on both engines, outputs compared (include deliberately failing ones).
- `go test ./...` and `go vet ./...` green.

## Speed (report numbers)

Benchmarks (`go test -bench`) on both engines: fib(25) recursive, a nested loop with arithmetic, string building, array push/index, hash get/set, closure calls. Target: **≥ 3× faster** than the tree-walker on the loop/call benchmarks, and fewer allocations (small integers cached, no map lookup for locals). Report the table.

## Phase 2 (linkz repo, after phase 1 is green)

- The host runs programs with `vm.Run` (one switch in `host/`, e.g. `host.Engine`, default VM, `LINKZ_ENGINE=tree` to fall back), all `evaluator.Eval` call sites through it.
- Update the pinned `../.lizard-stable` only when the orchestrator says so; until then test with a local `go.work` pointing at `../lizard-script`.
- `go test ./...` green in linkz on the VM; report before/after: `linkz calibrate` (recommended objects), repaint ms with 2,000 shapes (QA BUG-7 numbers: 40–80 ms), `linkz trace docs/Art/References/lilies-stargazer.png` at `--max-side 256` (1.25 s today), splash load time.

## Done when

Phase 1: both engines give identical output on every test program and the differential set; benchmarks ≥ 3× on loops/calls; docs updated. Phase 2: Linkz green on the VM with the measured speed-ups.
