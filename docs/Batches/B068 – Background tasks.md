# B068 – Background tasks (`task`), then a full-size trace preview

> "Why is the trace image different from the one we get on the actual canvas?" … "Do b now" — the developer, 2026-11-04

The trace dialog previews at 160 px but Insert/Open trace at up to 1024 px, so they differ. Fix it properly: let lizard run work off the main thread, then trace at full size in the background.

## Step 1 – lizard-script: the `task` Terrarium module (Opus)
Its CLAUDE.md rules (§7 row, TERRARIUM.md, GRAMMAR/SEMANTICS only if needed, guide en+es with run examples, both engines).
- `task.start(module, function_name, args = [])` → a task handle. It runs `function_name(args…)` from `module` (a path resolved like `use`, relative to the caller's file, and **from the bundle** in a `lizardscript build` app) in a **fresh interpreter on its own goroutine**, same engine as the caller (VM or tree-walker).
- **Nothing is shared.** Args and the result are deep-copied across (numbers, strings, booleans, null, lists, hashes, points, colors). A function or other non-data value in either direction is an error that names where it was found.
- `t.done()`, `t.result()` (the value, or raises the task's error as an error at the caller), `t.cancel()` (stops it soon: the VM/evaluator checks a cancel flag at the step-limit hook), `t.state()` (`"running" | "done" | "failed" | "cancelled"`). Optionally a window event `{"type": "task", "id"}` when one ends, if a window is open (same mailbox pattern as `choose_file`).
- Inside a task: `img`, `fs`, `os` (read-only parts are fine; say which), `math`, `text`, `json`, `hash`, `time`, `lizard` work; `window`, `pen`, `clip` raise "not available inside a task". `print` goes to stderr with a `[task]` prefix or is captured; decide and document.
- A task that never ends doesn't keep the program alive at exit; many tasks may run at once (limit to NumCPU running, queue the rest).
- Tests: results and errors on both engines, deep copy (mutating a passed hash in the task doesn't change the caller's), cancel of a busy loop within ~50 ms, forbidden modules, a function value refused, from a built bundle (`lizardscript build` then run it), `go test -race` on the task package. `go test ./...`, `go vet`, cross-builds. Never mention Linkz. Don't commit.
- Another builder (B064) is editing `tools/build/` and the window `icon` option right now: put the task code in new files (`terrarium/task/`), keep shared edits minimal and listed.

## Step 2 – Linkz: the dialog traces at full size in the background (Sonnet)
Rules `_v11-rules.md`. Read `lib/ui/trace_view.liz` (`tr_run`, `tr_code`, `tr_insert`, `tr_insert_scale`, `tr_open_drawing`) and `lib/ui/trace_model.liz`.
- Keep the instant 160 px preview while sliders move. When nothing has changed for ~400 ms, start a `task` that traces with exactly the options and size Insert/Open use; cancel it on any change. When it finishes, the preview shows it (a small "full size" badge; a "tracing full size…" badge while it runs).
- Insert / Open reuse the finished full-size result (no second trace) — the canvas is exactly the preview. If it isn't finished yet, they wait for it (showing the badge), not trace again.
- Rig tests (fake or real task, both engines): the result reused, cancel on slider change, Insert during a running task. Drive the real app (isolated HOME) on a reference image and save before/after shots of preview vs canvas into `docs/Art/Checks/v1.2/trace-full/`.
