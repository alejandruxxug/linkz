# B007c – Autocomplete (mandatory for v1)

Milestone: M3 · Repo: `linkz` · Lane: lizard (`lib/ui/complete_*.liz` + small hooks in the editor view) · Needs: B007b

> "also mandatory autocomplete" — the developer, 2026-11-03

## Goal
While typing in a code row, a suggestion list completes names, and a signature hint shows the call's parameters, so nobody has to remember the API.

## Sources (all from data the editor already has)
| Kind | From |
|---|---|
| Drawing functions | `src.library_docs()` (names, params, doc summary, example) |
| The user's functions, `let` values, layer names | `src.parse` of the current text (with their `-~>` docs) |
| Members after a dot | `math.`, `text.`, `list.` (keys of the builtin hashes), `pens.` (`pens.list`), `colors.` / any user hash literal's keys when it's a known `let` |
| Keywords | `let function layer if elif else for while return break continue try catch true false null` |
| Snippets | the tool templates (`circle` → `circle((256, 256), 40, #3fa34d)` from the doc's `example:`) |

## Behaviour
- Opens after 1+ letters of an identifier (or right after `.`), closes on Esc, on moving the caret off the word, or when nothing matches.
- Ranking: exact prefix → word-start/fuzzy (`crc` → `circle`, `sw` → `stroke_width`) → recently accepted (per session) → shorter first. Max 8 visible, scroll for more.
- Each item: kind icon, name with the matched letters highlighted, short signature `circle(center, r, fill)`, doc summary on the selected item.
- Keys: ↑/↓ move, Tab or Enter accept, Esc dismiss. Enter only accepts when the list is open; otherwise it's the normal new-row Enter.
- Accepting a function inserts the call with its example arguments and selects the first argument (caret over it), so typing replaces it; Tab jumps to the next argument.
- **Signature help**: inside a call's parentheses, a small hint above the row shows `circle(center, r, fill)` with the current parameter highlighted and its doc line.
- Works in both languages (labels/hints via `t()`; names stay code).
- Styled like the approved menus / pen picker (Linkz Dark, 18 px rows); a polish frame comes in R002.

## Backdoor
`Dump` includes `complete: {open, items, selected}`; `Call Complete` forces it open at the caret.

## Tests
Lizard tests for ranking, member completion after `.`, user names from the document, snippet insertion + first-argument selection, signature hint's current parameter, Enter vs new-row behaviour. Tester: type `ci` → circle first; accept → template inserted and drawn; type `math.s` → sin, sqrt; inside `circle(` the hint highlights `center`, after the first comma `r`.

## Done when
All above pass, typing feels instant (report ms per keystroke with the list open on `layers.lzdraw`), `go test ./...` green.
