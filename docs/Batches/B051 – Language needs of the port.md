# B051 – What the all-lizard port needs from the language (M11)

Repo: `../lizard-script` · Needs: B041 findings (`../linkz/docs/Linkz – Language Needs.md`, newest rows) · Parallel with B043 (VM: owns code/, compiler/, vm/, object/, evaluator/evaluator.go) and B046 (lsp/, editors/)

The developer allows language changes that make it better (2026-11-03). Each item: tests on both engines, docs (SEMANTICS/TERRARIUM/DESIGN §7, the LEP when it's `use`), guide where users see it (en + es, real outputs).

1. **`use` cycles allowed** (Python's model; it blocks the editor port today): while a module is loading, a `use` of it gets its live, partly loaded scope instead of an error, so names defined so far are visible and the rest appear as the module finishes (live names make this work). Amend LEP 0002 (Amendment 2: why, the rule, the pitfall "a name used at load time before its module defines it", how `check` should warn about it). The old "use cycle" error tests change on purpose; list them. The one-line prototype is in `/private/tmp/claude-501/-Users-urrexx-Documents-Uni-4th-Semester-Langs-Compilers-lizard-script/7cd1ec77-d373-4271-93ca-7500f2cde90e/scratchpad/use-cycles.patch` (evaluator/modules.go: `if m.loading && node.Name != nil`), but make it clean and complete (also `use … as` during a cycle, error messages that point at the cycle when a missing name is really "not loaded yet").
2. **Faster repeated runs**: cache parsed + compiled modules by file and content (both engines), so `lizard.run` with the same libraries doesn't reparse/recompile them each run; `lizard.run`'s `on_statement` gets an option to fire only for the run's own source (`"hook_source": name`), not for library statements.
3. **`lizard.run` with live modules**: an option `"use": ["path", …]` so a sandboxed run sees modules as live names (writes visible to the library), not copies through `globals`.
4. **Time**: `time.now()` adds `epoch` (seconds since 1970, float), `utc_offset` (seconds) and `weekday`; `time.utc()`.
5. **Content hash**: `text.digest(s, "sha256"|"fnv")` and `fs.digest(path, …)` (hex strings).
6. **Scaled blit**: `blit(target, source, pos, {"scale": k, "filter": "nearest"|"linear"})` on every target; `window.screenshot(w, {"full": true})` returns the full window pixels even without a surface.
7. **Errors raised in library functions at the caller's line** (lily2d/3d found this): a way for lizard-written library functions to report an error at the position of the call into the module (e.g. `error(message, {"at": "caller"})`, or automatically for functions of standard modules) — design it, document it, use it in lily2d/lily3d's argument checks.

Clean code. Don't commit. Never mention Linkz inside lizard-script. Report per item: semantics, files, tests, any changed .expected/test with reasons.
