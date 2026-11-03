## Rules for every v1.1 batch
- **The approved design is the spec**: the frames and metrics in `docs/Linkz – Style Guide.md` (R002/R003 sections, approval note at the end) and the PNGs in `docs/Art/Reviews/R002/` and `R003/`. Match them in the real app; compare with `lizardscript drive Screenshot` side by side and report differences.
- Linkz is **all lizard** (no Go). Anything missing from the OS or the language goes to lizard-script as a need (report it), never a workaround.
- **The code is the drawing**: every visible change is a code change (rows/literals), one undo step per gesture, the code panel stays in sync; refused edits say why (as today for expression positions).
- English + Spanish strings for everything (`lib/ui/lang_en.liz`, `lang_es.liz`): **add only your own keys**, at the end of a section.
- Clean code: one concept per file (new files for new features), documented with `-~>`, LEP 0001 style; `lizardscript fmt --check` and `lizardscript check` clean.
- Tests: lizard tests in `tests/` (headless editor rig) for the feature, both engines (`scripts/test.sh`); `scripts/parity.sh` and `scripts/smoke.sh` still pass; the M10 benchmarks (`tests/editor/editor_bench.liz`) don't regress by more than 5%.
- Drive the real app to check it (`LIZARD_AUTOMATION=1`, isolated HOME and config dir, `lizardscript drive`); never touch the user's ~/Documents/Linkz, ~/.linkz or real config.
- Other builders work in `lib/ui` at the same time: own your files, add only your lines to shared ones (`view_main.liz`, `view_input.liz`, `view_popups.liz`, `app.liz`, lang files), and list every shared edit in the report. Don't commit. Don't edit ../lizard-script.
