# B036 – Incremental editor (M10 Performance, Linkz side)

Repo: `linkz` · Lane: lizard (`lib/ui/editor_*`, `view_rows`/`view_panels`, `complete_*`) + `host/src.go`/`host/scan.go` only if the profile proves it · Needs: B035

## Why
After B035 the canvas costs 3.6 ms per edit at 2k shapes, but a whole tick is still 84–114 ms (B035's profile): every tick rebuilds the document from scratch, several times.
- `view_rows`: `src.parse` of the whole document (25 ms, Go `ScanSource` + converting to lizard) + `doc_rows` (18 ms).
- `doc_row(m, line)` rebuilds every row and is called 3–4 times per tick (`props_of`, `move_shape`, the props panel): ~12 ms each.
- `cmp_sync` (autocomplete) 15–36 ms per tick, 205 ms at 10k.
- Go's GC is ~45% of CPU (boxed numbers; the VM work in lizard-script B033 helps separately).

## Build
1. **Rows cached per document version**: `doc_rows` computed once per text change (keyed by the model's version), `doc_row(m, line)` an O(1) lookup into it. First cheap win (~35 ms at 2k).
2. **Incremental scan/parse**: when an edit touches a few lines, rescan only the statements those lines belong to and splice them into the cached row list (line numbers after it shift). Keep a full rescan as the fallback and as the check in tests (incremental result `===` full result after random edits).
3. **Props panel and rows panel** draw from the cache; only visible rows are laid out (virtualized list).
4. **Autocomplete**: `cmp_sync` keeps its symbol index across ticks and updates it only for changed statements; it does nothing when the caret and text didn't change.
5. **Per-tick budget report**: the backdoor `Stats` gets a breakdown per phase (rows, props, run, canvas, complete) so regressions are visible.

## Done when
Before/after table on the B035 fixtures (`docs/Batches/fixtures/shapes-2k.lzdraw`, `shapes-10k.lzdraw`, `go test ./host -bench Editor`): small edit, drag frame, keystroke — target **< 16 ms at 2k** and < 50 ms at 10k. A randomized test that incremental == full after many edits. All tests and smoke green on the VM and `LINKZ_ENGINE=tree`.

## Rules
Clean code (one concept per file, documented, no hacks). New lizard code must already follow lizard-script's upcoming newline rule (LEP 0001, accepted): **put binary operators at the end of a line, never start a line with an operator, `(` or `[`** — so it needs no migration later. linkz builds against `../.lizard-stable` (don't change go.work). Testers: `HOME=$(mktemp -d)`, `--no-splash`. Don't commit.
