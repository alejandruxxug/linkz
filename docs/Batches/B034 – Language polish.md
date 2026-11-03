# B034 – Language polish (bugs and learner surprises found by the guide)

Repo: `../lizard-script` · Lane: Go (lexer, parser, evaluator builtins, repl) + compiler/vm where the AST changes · Needs: nothing

Found while writing docs/learn (2026-11-03). The developer allowed language changes that make it better; keep every existing program working unless listed under "changes behaviour".

| # | What | Kind |
|---|---|---|
| 1 | `rgb(300, 0, 0)` says `got 300.0`; show integers as integers (check every color/point builtin's messages) | bug |
| 2 | The REPL shows strings inside arrays/hashes unescaped (`"[1, "a"]"`): quote and escape them as code so they can be pasted back (`["1", "a"]` style per the language's own literal syntax). Check `.expected` files: if `puts`/`Inspect` output changes for existing tests, keep `puts` of a string unquoted (as today) and only change nested strings — report any `.expected` change | bug (changes output) |
| 3 | Exponent float literals `1e3`, `2.5e-4`, `1E+6` (today: "unknown unit e3"); units stay whole words | addition |
| 4 | `puts(5kg)` / `puts(1e3)` today say `expected ")"`: report `unknown unit "kg"` inside calls too | better error |
| 5 | Accept `else if` as a synonym of `elif` (today a syntax error) | addition |
| 6 | `continue` outside a loop: a syntax error like `break` (today `identifier not found: continue` at run time) | better error |
| 7 | Document (DESIGN §4, SEMANTICS, guide's "common mistakes") the surprises we keep: negative indexes give null on `[]` but count from the end in `text.slice`/`list.*`; `"5" * 2` is `"55"`; byte-wise string comparison | docs |
| 8 | A line starting with `-` continues the previous expression: **don't change it** (it would break programs); write it up as a LEP draft `proposals/0001-statement-ends-at-newline.md` (status Draft) with the options | LEP |

Every change: both engines (if the AST changes, the compiler must handle it), tests (lexer/parser/evaluator tables + programs in tests/valid or tests/errors with .expected), GRAMMAR/SEMANTICS/DESIGN §7, and the guide (`docs/learn/en` + `es`: chapters 4, 5, 6, 11, cheatsheet, exercises if affected; every output really run).

Rules: lizard-script CLAUDE.md, clean code. Don't commit. Another builder (B033) is editing `compiler/` and `vm/` for speed at the same time: keep your compiler/vm edits minimal (only what a new AST shape needs) and list them in the report.
