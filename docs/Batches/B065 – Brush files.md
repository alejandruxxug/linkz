# B065 – `.lzbrush` files and the brush library

Model: Sonnet. Rules: `_v11-rules.md`. Read only: `docs/Linkz – Brushes.md` (the `.lzbrush` section), the R005 approval paragraph in `docs/Linkz – Style Guide.md` (answers 2, 3, 7), `lib/ui/pen_library.liz`, `lib/ui/pen_code.liz`, and how the `lizard` Terrarium module runs code sandboxed (`lizard.run` with `steps` and no `use`/file access; see how `lib/app` already uses it).

Engine and file layer only; the picker menu, drop and import dialog come in a later batch, the creator UI in B066.

- **Format:** `-~> lzbrush 1` header, then `let brush = {...}` (flat keys as `lib/draw/pens.liz` documents, including the stamped ones), bitmap shape/grain data as run-length lizard data (like `lib/ui/splash_art.liz`). `.lzbrushes` = header `-~> lzbrushes 1` + `let brushes = [...]`.
- **API** (new `lib/app/brush_files.liz`, documented): `brush_write(pen) -> text`, `brush_read(text) -> {"brush", "warnings"}` or an error with a reason (a step-limited sandboxed run: no `use`, no fs/os/window, only data; unknown keys are dropped with a warning; values out of range are clamped with a warning), `brushes_write(list)`, `brushes_read(text)`, `brush_library_dir()` = `~/.linkz/brushes`, `brush_library_save(pen) -> path`, `brush_library_load() -> list`, `brush_library_export_all(path)`, `brush_import_file(path) -> {"added": [...], "warnings": [...]}`.
- **Migration:** pens saved by the old pop-up in `~/.linkz/pens` are read as before and written to `~/.linkz/brushes` the first time; nothing is deleted. `pen_library.liz` loads from the brush library.
- **Copy on first use:** when a drawing first uses a library brush, its hash is written into the drawing's code (`pens.<key> = {...}`, after the `let colors` block, one undo step with the stroke), so a `.lzdraw` renders the same anywhere. Built-in pens are never copied.
- Round-trip tests (every built-in and the four stamped example pens: write → read is identical), hostile files (a loop, file access, a huge list, a syntax error: each refused with its reason, nothing written), migration, copy-on-first-use in the rig. Both engines.
- Others work in `lib/ui/pen_window*.liz` (B066) and `lib/draw/stamps*.liz` (B067) at the same time: stay out; add only your lines to `pen_library.liz` and the shared files, and list them.
