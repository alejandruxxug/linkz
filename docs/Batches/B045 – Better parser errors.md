# B045 – Better parser errors: recovery and "did you mean"

Repo: `../lizard-script` · Lane: Go (`parser/`, `lexer/` messages, `tests/errors`) · Parallel with B044 (fmt/check)

## 1. Error recovery (pending since the first app version)
Today one bad token can make the parser report fake follow-on errors on later, valid lines (Linkz dims them as "may be caused by the error above"). Implement proper panic-mode recovery: after an error, skip to a synchronization point (the end of the current line under LEP 0001, a `}` that closes the current block, or the next statement keyword at line start), so each real mistake gives one error and valid lines after it give none. Keep reporting several independent errors in one run.
- This intentionally changes some `.expected` files in `tests/errors` (e.g. `syntax_several_errors.expected`): list each change with before/after and why. The developer allowed language changes that make it better; the team gets a heads-up in DESIGN.md §7.
- A test corpus: take every valid program in tests/valid + examples, inject one random error (delete/duplicate/replace a token) at many positions, and check (a) exactly one error is reported near the injected spot in most cases (report the rate), (b) never more than 2.

## 2. Messages that teach
- "did you mean" for keywords and common slips: `fucntion` → `function`, `elseif`/`elsif` → `elif`, `=` in a condition (`if (x = 1)` → "did you mean =="?), missing `)`/`}` with where it was opened ("the ( on line 3 is never closed"), unterminated strings with their start, `let` without a name, `:` vs `=` confusions in hashes.
- Every message keeps `[line L, column C]` and is short and plain; update GRAMMAR/SEMANTICS examples and the guide's error examples (en + es, real outputs).

## Done when
Tests (parser tables, the injected-error corpus with its rates, tests/errors programs), the `.expected` changes listed, docs + guide. Clean code. Don't commit. Ownership: don't edit `code/`, `compiler/`, `vm/`, `object/` (B043), `evaluator/modules.go`/`terrarium/std/` (lily), `format/`, `check/`, `main.go` (B044).
