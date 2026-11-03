# B046 – `lizardscript lsp`, a VS Code extension and a Zed extension

Repo: `../lizard-script` · Needs: B044 (fmt + check) committed · Folder: `editors/`

> "VS Code support (LSP) … also an extension here for Zed!" — the developer, 2026-11-03

## Language server (`lizardscript lsp`, Go, package `lsp/`)
Speaks LSP over stdio (JSON-RPC; standard library only, or a small vetted dependency). Features:
- diagnostics on open/change (syntax errors + `check`'s errors/warnings/hints, with ranges from AST spans), debounced;
- formatting (`fmt`), on save if the editor asks;
- completion: names in scope (with lizard's late-bound rules), keywords, snippets for `function`/`if`/`for`/`layer`/`try`, members after `.` for builtin libraries (`math.`, `text.`, `list.`, `hash.`), Terrarium modules, standard modules (`lily2d.`, `lily3d.`) and hashes whose keys are known, with signatures and the `-~>` doc summaries (`lizard.docs`);
- hover: signature + doc comment for functions and builtins, the value of literal constants, the type of a literal (COLOR shows a swatch hex, POINT its coordinates, units converted: `90deg = 1.5708 rad`);
- signature help inside calls (current parameter highlighted);
- go to definition and find references (lets, functions, parameters, across `use`d modules);
- document symbols (functions, top-level lets, layers as a tree) and workspace symbols;
- rename (scope-aware), folding ranges (blocks, layers, hashes), semantic tokens (colors and points highlighted as such);
- color provider: COLOR literals (`#f2a65a`) show the editor's color picker and edits write back a lizard literal.
Tests: an in-process client driving every request on fixture files; no editor needed.

## Syntax grammar
- `editors/tree-sitter-lizard/`: a Tree-sitter grammar (grammar.js, generated parser, queries: highlights, indents, folds, locals, brackets) covering the whole language incl. LEP 0001 lines, units, colors, points, layers, `use`/`as`, `-~>` comments; corpus tests (`tree-sitter test`) from tests/valid and examples. (Generating needs the tree-sitter CLI via npm/npx; it's a build tool only.)
- `editors/vscode/syntaxes/lizard.tmLanguage.json`: a TextMate grammar for VS Code (it doesn't use Tree-sitter).

## VS Code extension (`editors/vscode/`, TypeScript)
Language id `lizard`, `.liz` (and `.lzdraw` as an associated extension), the TextMate grammar, language configuration (comments `-~>`, brackets, auto-closing, indentation rules), snippets, the LSP client starting `lizardscript lsp` (path setting `lizard.path`), commands "Run file" (terminal `lizardscript file`), "Run on VM", "Format", "Check workspace", an icon. Package a `.vsix` locally (`npx @vscode/vsce package`) — **don't publish** to the Marketplace (the developer decides).

## Zed extension (`editors/zed/`)
`extension.toml`, the language config (`languages/lizard/config.toml`: name, `.liz`/`.lzdraw` suffixes, comment `-~> `, brackets, tab size 2), Tree-sitter queries (from the grammar), and the Rust extension (`src/lib.rs`, zed_extension_api) that starts `lizardscript lsp` (finding it on PATH or a setting). Document "Install dev extension" in Zed. **Don't publish** to the Zed registry.

## Docs
`editors/README.md` (install for each editor, building from source), guide chapter 18 "tools" section (en + es), README link. Clean code. Don't commit. Never mention Linkz in lizard-script.

## Done when
LSP tests green; tree-sitter corpus passes; the `.vsix` builds; the Zed extension builds (`cargo build --target wasm32-wasip1` or Zed's dev-install flow) and its queries parse; a manual check list for the developer (open a .liz in VS Code and in Zed: highlighting, diagnostics, completion, hover, format, go to definition, color picker).
