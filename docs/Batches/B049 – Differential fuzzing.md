# B049 – Differential fuzzing of the two engines and the parser

Repo: `../lizard-script` · New folder only: `fuzz/` · Read-only everywhere else

- **Go native fuzz tests** (`go test -fuzz`) in `fuzz/`:
  1. *Parser*: random bytes and token soups never panic, never hang (time-boxed), and every error has a valid `[line L, column C]` inside the input.
  2. *Engines agree*: a **program generator** (well-typed-ish random lizard programs: lets, arithmetic incl. floats/units/points/colors, strings, arrays, hashes, ifs, loops with break/continue, functions and closures incl. defaults and recursion, try/catch, layers) runs each program on the tree-walker and the VM with a step limit and compares output, value, error message/position/uncatchable, step count and hook calls — the same contract as vm/differential_test.go. Shrink failures to a minimal program automatically.
  3. *fmt/round-trip* (if B044's `format` package exists when you run): formatting never changes behaviour.
- A seed corpus from tests/, examples/ and the guide's code blocks.
- Run it for a meaningful time (e.g. 10–20 minutes per target) and **report every bug found** with its minimal program (don't fix engine bugs yourself — other builders own those files; write each as a failing test in `fuzz/testdata/` named by what it shows, so the owner can fix it).
- `fuzz/README.md`: how to run, how to add a generator feature. Don't commit. Don't edit files outside `fuzz/`.
