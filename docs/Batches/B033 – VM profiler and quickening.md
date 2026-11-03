# B033 – VM profiler and quickening (M10 Performance, language side)

Repo: `../lizard-script` · Lane: Go (compiler, vm, main) · Needs: B029

> "Continue building the milestones; if it would be better for the lang to change some things and be more optimized, do it." — the developer, 2026-11-03

## 1. Profiler (first: measure before optimizing)
- `lizardscript --profile file.liz` (VM) prints, after the run, the top 25 lizard functions by **self time** and **total time**, with call counts and the source line of each function (`name (file:line)`; anonymous functions as `function@line`). Also the top builtins.
- For hosts: `vm.StartProfile()` / `vm.StopProfile() Profile` (a plain struct with the rows, and `Profile.Write(w io.Writer)`), so an app can profile one frame.
- Sampling or counting — choose the cheapest that is accurate enough; profiling off must cost ~nothing (one branch per call).
- Docs: a "Profiling" section in DESIGN.md and a short part in docs/learn (en + es) chapter 7 or a new short chapter.

## 2. Make the VM faster (use the profiler + `go test -bench` to prove each step)
- **Global and builtin lookups**: cache name reads of the host env and builtins (`push`, `len`, `math`, top-level functions) with correct invalidation (a `let`/assignment of that name, or the host defining it later). Measure.
- **Builtin calls**: avoid copying the argument slice when safe (e.g. a VM-owned scratch slice per call depth that builtins may not retain — document the contract in object.Builtin; check every builtin in evaluator and terrarium does not keep `args`).
- **Quickening / specialized instructions** where the profiler says it matters: e.g. integer-only `ADD/SUB/LT` fast paths already exist — add float fast paths, `GET_FIELD` with a constant string key (dot access) with a per-instruction cache, superinstructions for the hottest pairs (`GET_LOCAL; CONSTANT; ADD`), small-int and small-float caching where it avoids allocations.
- **Points and colors**: geometry code creates many POINT objects; reduce allocations in point arithmetic on the VM (e.g. avoid intermediate copies) without changing semantics.
- Keep behaviour identical: the differential tests (`vm/differential_test.go`) and both-engine program tests must stay green. Add benchmarks for dot access, float math, point math, builtin-heavy loops.

## Report
Before/after benchmark table (ns/op, allocs/op), a profile of `examples/terrarium/clock.liz` and of a geometry-heavy program, and what each optimization bought. Target: another ≥ 1.5× on float/point/dot-access-heavy code.

## Rules
lizard-script CLAUDE.md (both engines identical; tests green; clean code: one concept per file, documented). Don't commit. Another builder works on the lexer/parser/evaluator builtins at the same time (B034): don't edit `lexer/`, `parser/`, `evaluator/builtins.go`, `evaluator/shapes.go`, `repl/`; coordinate through your report if you must.
