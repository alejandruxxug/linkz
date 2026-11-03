# B015 – Home screen

Milestone: M3 · Repo: `linkz` · Lanes: lizard (`lib/ui/home_*.liz`) + host (thumbnail cache primitives) · Needs: B007b (the editor view), approved Figma frames R001 #8–#10

## Goal
`linkz` / `linkz edit` with no file opens a **home screen**: a gallery of the user's drawings, working ones first, each with a small preview rendered by lizard.

## What it shows
- **Working on**: the most recently edited drawings (from `~/.linkz/recent.liz`, max 8), bigger cards.
- **All drawings**: every `.lzdraw` in the Linkz folder (`~/Documents/Linkz/` by default; `settings.drawings_dir` overrides), newest first, smaller cards. Search box filters by name.
- **＋ New drawing** card first: creates `Untitled N.lzdraw` with the `lzdraw 1` header and a canvas line, opens it in the editor.
- Card: thumbnail, file name, relative edit time via `t()` ("edited 2 h ago" / "editado hace 2 h"), ⋯ menu: open, rename, duplicate, show in Finder, delete (delete asks to confirm in-app; never deletes without it).
- Empty state (no drawings): the lizard, "＋ New drawing", "Open an example" (copies `examples/layers.lzdraw` into the folder).
- Opening a drawing switches to the editor; the editor's toolbar has a 🏠 button back to the home screen (unsaved changes → confirm).

## Thumbnails
- Rendered by lizard: `src.run(code)` → `raster.paint(paint_order(shapes))` into an offscreen buffer of the card size (canvas scaled to fit, centered, background drawn), saved with `fb.save_png` into `~/.linkz/thumbs/<hash>.png`.
- Cache key: a hash of the file's text (`host` primitive `fs.hash(path)` or `text.hash(s)` — add the smallest one needed) so an unchanged file is never re-rendered; also refreshed when the editor saves.
- Loading cached thumbs needs a host primitive to read a PNG back into a lizard-drawable form: `fb.load_png(path)` → a named kept buffer that `fb.blit(name, pos)` draws (Go copies pixels lizard already rendered; no drawing decisions).
- Rendering happens lazily, one card per idle paint, newest first, with a placeholder (checker + file name) until ready. A drawing that errors gets an error thumbnail (dim canvas + ⚠).

## Backdoor
`Call Home`, `Call OpenCard n`, `Call NewDrawing`, `Dump` includes the cards (name, path, thumb state, rect).

## Tests
Lizard tests for listing/sorting/searching/recent bookkeeping/relative times (en + es); a Go test that a thumbnail is rendered once and reused while the file is unchanged; the tester opens the home screen with 6 drawings via the backdoor, screenshots it, opens one, goes home.

## Done when
The home screen matches the approved frames, thumbnails appear and are cached, new/open/rename/duplicate/delete work, `go test ./...` green.
