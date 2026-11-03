# B039 – lizard-script: modules, reflection, io, OS, test runner, build (M11)

Repo: `../lizard-script` · Plan: [[Linkz – All Lizard Plan]] · Parallel with B040 (see ownership)

## Build (general features, no app names)
1. **Modules: `use`.** `use "draw/shapes"` loads `draw/shapes.liz` relative to the file that says it (or a module search path), runs it once per program in its own top scope, and makes its top-level names available. Decide and document the exact semantics (all names vs `use "x" as x` giving a hash; cycles are an error; a module's errors carry its file name via source names). Works on both engines and in `lizardscript build` bundles. A LEP (`proposals/0002-modules.md`, Final) plus GRAMMAR/SEMANTICS/DESIGN/guide (a chapter, en + es).
2. **Syntax trees with end positions**: every AST node gets its end line/column (additive). 
3. **`lizard` Terrarium module (reflection)**: `lizard.tokens(text)` (kind, text, line, column, end), `lizard.parse(text, name = "")` → the syntax tree as plain hashes/arrays with start/end positions and parse errors (like the existing host `src.parse` shape: read `../linkz/host/src.go`, `scan.go` for what an editor needs), `lizard.run(code, options)` → runs in a fresh environment seeded from `options.globals` (a hash), with `steps` limit, a source `name`, optional statement/layer hooks as lizard functions, returning `{value, error: {message, line, column, source}}`; `lizard.docs(code)` → the `-~>` doc comments of functions (name, params, defaults, summary, example) like `src.library_docs`. Both engines.
4. **Terrarium `io`**: `io.read_line()` (null at EOF), `io.write(text)`, `io.error(text)` (stderr). Port from `../linkz/host/stdio.go`.
5. **`fs`/`os` additions**: `fs.info(path)` → {size, mtime, dir}, `fs.list` with info option, `fs.rename`, `fs.copy`, `fs.trash(path)` (macOS Trash, Windows Recycle Bin, Linux freedesktop Trash; port `../linkz/host/thumbs.go`, `trashinfo.go`, `sys_*.go`), `os.reveal(path)` (Finder/Explorer/xdg-open), `os.memory()` → {total, free} (port `mem*.go`), `os.open_url`.
6. **`lizardscript test [dir]`**: runs every `*.liz` with a sibling `.expected` (output compared, like tests/programs_test.go) and every `*_bench.liz` when `--bench`; exit code; works with `--vm`.
7. **`lizardscript build <dir> -o <exe>`**: bundles a program dir (`main.liz` + everything it `use`s + data files it lists) into a single executable for the current or a given GOOS/GOARCH (e.g. append an archive to a copy of the lizardscript binary, or generate a tiny Go main that embeds the files and build it); `lizardscript run <dir>`.

## Ownership (B040 runs at the same time)
B039 owns: `ast/`, `lexer/`, `parser/`, `token/`, `evaluator/` (module loading), `compiler/`, `vm/`, `repl/`, `main.go`, `terrarium/lizard*`, `terrarium/io.go`, `terrarium/fs.go`, `terrarium/os.go`, `tools/`, docs for these. B040 owns `terrarium/window/**`, `terrarium/img.go`, `terrarium/internal/**`, `terrarium/pen/**`, `terrarium/gfx*`. `terrarium/terrarium.go` (registration): coordinate by adding only your own lines.

## Done when
Tests on both engines, guide + docs, cross-builds, `lizardscript build` produces a working single file on macOS (and cross-builds for Windows/Linux). Clean code. Don't commit.
