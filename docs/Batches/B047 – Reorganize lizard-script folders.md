# B047 – Reorganize lizard-script's folders (grouped by role)

Repo: `../lizard-script` · Needs: every running lizard-script batch committed (B043, B044, B045, B031, B032; B046 after) · One atomic commit

> Chosen by the developer, 2026-11-03: "Grouped by role", plus graphics libraries as the **flower** family.

```
lizard-script/
├── main.go            the lizardscript command (stays at the root: `go build` keeps working)
├── lang/              token/ lexer/ ast/ parser/
├── engine/            object/ evaluator/ code/ compiler/ vm/
├── terrarium/         one folder per module: fs/ os/ io/ img/ clip/ time/ json/ window/ pen/ lizard/ (internal/ stays)
│   └── std/flower/    lily2d, lily3d (use "flower/lily2d")
├── tools/             repl/ test/ build/ bundle/ fmt/ check/ lsp/ drive/ checknewlines/
├── editors/           vscode/ zed/ tree-sitter-lizard/
├── docs/              DESIGN, GRAMMAR, SEMANTICS, STYLE, TERRARIUM, flower/, learn/, proposals/
└── tests/ examples/ bench/
```

- `git mv` only (history kept), import paths updated (`lizardscript/lang/lexer`, `lizardscript/engine/vm`, …), gofmt, no behaviour change: `go test ./...`, both engines, `lizardscript test tests`, `bench/compare.py`, cross-builds.
- Every doc/guide path and the README project tree updated; CLAUDE.md "Read first" paths; CONTRIBUTING; links in the LEPs.
- The Linkz Go host imports lizardscript packages: list every import that changes for the Linkz port/B042 (Linkz stays on its pin until then).
