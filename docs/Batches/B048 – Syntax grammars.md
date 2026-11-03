# B048 – Syntax grammars: Tree-sitter and TextMate (split out of B046)

Repo: `../lizard-script` · New folders only: `editors/tree-sitter-lizard/`, `editors/vscode/syntaxes/` · Read-only everywhere else

The editor extensions (B046) need these, and they don't depend on the language server, so they're built now.

- **Tree-sitter grammar** (`editors/tree-sitter-lizard/`): grammar.js for the whole language as it is on main (read docs/GRAMMAR.md, proposals/0001 and 0002, parser/): statements ending at line ends with the LEP 0001 continuation rules (use an external scanner in C if needed for newline significance), `let`, functions (defaults), `if/elif/else if/else`, `while`, `for`, `break/continue/return`, `try/catch`, `layer` blocks with Python-style indentation (external scanner for indent/dedent), `use "…" as name`, numbers with exponents and units (`90deg`, `2px`), colors `#rrggbb[aa]`, points `(x, y[, z])`, strings with escapes, hashes, arrays, dot access, `-~>` comments. Queries: `highlights.scm`, `indents.scm`, `folds.scm`, `locals.scm`, `brackets.scm`, `injections.scm` (none needed?). Corpus tests (`test/corpus/*.txt`) generated from tests/valid and examples, plus error cases; `tree-sitter test` passes; parse every `.liz` in the repo without ERROR nodes (report any).
- **TextMate grammar** (`editors/vscode/syntaxes/lizard.tmLanguage.json`) for VS Code, same highlighting categories (keywords, functions, builtins, Terrarium modules, colors, points, units, strings, numbers, comments, operators).
- Needs the tree-sitter CLI (npx/npm, build tool only) and a C compiler; nothing is published.
- `editors/README.md` section on the grammars. Clean, commented grammar. Don't commit. Don't edit files outside `editors/`.
