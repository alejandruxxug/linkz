# B047 – Reorganize lizard-script's folders (grouped by role)

Repo: `../lizard-script` · Needs: every running lizard-script batch committed (B043, B044, B045, B031, B032; B046 after) · One atomic commit

> Chosen by the developer, 2026-11-03: "Grouped by role", plus graphics libraries as the **flower** family.

```
lizard-script/
├── main.go            the lizardscript command (stays at the root: `go build` keeps working)
├── lang/              token/ lexer/ ast/ parser/
├── engine/            object/ evaluator/ code/ compiler/ vm/
├── terrarium/         unchanged inside for now (flat module files, window/, pen/, internal/, std/flower/)
├── tools/             repl/ test/ build/ bundle/ fmt/ check/ lsp/ drive/ checknewlines/
├── editors/           vscode/ zed/ tree-sitter-lizard/
├── docs/              DESIGN, GRAMMAR, SEMANTICS, STYLE, TERRARIUM, flower/, learn/, proposals/
└── tests/ examples/ bench/
```

- `git mv` only (history kept), import paths updated (`lizardscript/lang/lexer`, `lizardscript/engine/vm`, …), gofmt, no behaviour change: `go test ./...`, both engines, `lizardscript test tests`, `bench/compare.py`, cross-builds.
- Every doc/guide path and the README project tree updated; CLAUDE.md "Read first" paths; CONTRIBUTING; links in the LEPs.
- The Linkz Go host imports lizardscript packages: list every import that changes for the Linkz port/B042 (Linkz stays on its pin until then).

## Scope decided (2026-11-03)
- Parts 1–3 only: move the packages into `lang/`, `engine/`, `tools/`, `editors/` (already there), `docs/` (`proposals/` → `docs/proposals/`); **leave the Terrarium's inside as it is** (splitting it into one package per module is a later, separate task).
- **Rename the Go module** from `lizardscript` to `github.com/alejandruxxug/lizard-script` in the same pass, so `go install github.com/alejandruxxug/lizard-script@latest` works (check `go.mod`'s `replace` for `third_party/ebiten` still resolves, and `tools/ebitencopy`'s paths).
- Also update: `bench/compare.py`, `fuzz/`, `editors/` (the VS Code/Zed extensions only call the binary — check nothing imports packages), `lsp` tests that use repo paths, `.github` templates, every doc and guide path, the README tree, CLAUDE.md, CONTRIBUTING.
- Model: Sonnet. One commit. Verify: `go vet`, `go test ./...`, `lizardscript test tests` on both engines, `bench/compare.py --runs 1`, cross-builds, `go install` from a clean clone of the commit into a temp GOBIN.
