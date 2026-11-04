# 🔍 2026-11-04 – Linkz design check 2 (v1.2: B053 fixes, gestures, R004, trace dialog)

Second, smaller check. It looks again at every frame the [first check](2026-11-03%20v1.1%20design%20check.md) marked 🔴 or 🟡, at the gesture frames it couldn't capture before (drive now has `Press` / `Move` / `Release` / `Path`), at the seven R004 layer-hover frames and at the R003-6 trace dialog. No code was changed. Only the items still open are listed.

**How the shots were taken.** main at `4fecc1b`, `lizardscript run --vm . edit` with the stable `lizardscript` from the scratchpad, `LIZARD_AUTOMATION=1`, and an isolated `HOME`, `XDG_CONFIG_HOME` and `TMPDIR` set in the same shell for every command, drive included (`~/.linkz` and `~/Documents/Linkz` were never touched). Driven with `lizardscript drive --pid N`. Window 1024×720 (512×360 UI px). Home and dialog shots come back at 512×360 and are scaled ×2 nearest-neighbor. The drawings are the first check's (`docs/Art/Checks/v1.1/drawings/`) plus [`r004.lzdraw`](../Art/Checks/v1.2/drawings/r004.lzdraw): the sunset with the sun's faint halo, a nested `layer rays` (`blend("screen")`, `opacity(0.7)`, 8 lines from a `for`) and `layer boat -~> [hidden]`. The machine's budget in the isolated HOME was 3,900 objects. All shots are in [`docs/Art/Checks/v1.2/`](../Art/Checks/v1.2/).

Legend: ✅ matches (or the fix holds) · 🟡 small differences left · 🔴 wrong. Coordinates are UI px.

## Summary

| Verdict | Count | Frames |
|---|---|---|
| ✅ | 20 | R002-1, 5, 6, 9, 10, 11, 12, 13, 14, 15, 16, 18b, 19b, 21, 23 · R003-2, 3 · R004-2, 5, 6 |
| 🟡 | 17 | R002-2, 3, 4, 7, 8, 17, 20, 22, 24, 25 · R003-1, 4, 5, 6 · R004-1, 3, 7 |
| 🔴 | 1 | R004-4 |

Every B053 fix from the first check holds. What is left is mostly items the first check listed without a ✓ (properties of gradients and `let colors`, the opacity label, names), plus what the new gestures and R004 frames show.

### Remaining items, ranked

1. 🔴 **The artist hover card covers the layer it previews, so Alt isolation shows nothing** (R004-3, R004-4). With the top layer (`sun`) at y 111 in LAYERS (compact palette collapsed), the card sits at x 204–354, y 95–228 and covers the sun on the canvas (x 243–286, y 130–175). Holding Alt dims everything else, and the one layer left bright is under the card. Expected: the card never covers the hovered layer's canvas box. Slide it along the panel edge (down or up) until it clears the box; the tail stays on the row. When nothing clears it, put the card on the other side of the box.
2. 🟡 **Live brush ink has light seams** (R003-4). While the stroke is drawn, every segment join shows a 1-screen-px vertical hairline in a lighter color through the black ink (`R003-04-brush-string.png`, at each sample). The stroke you get after release is clean. Expected: the live stroke draws as one piece, or with joins that overlap, so it looks like the final stroke.
3. 🟡 **Hidden layer rows** (R004-7). (a) The `eye-off` on a hidden row is `dim` `#5f685b`, in the code rows and in LAYERS; the spec is `amber`. (b) The row shows the flag comment: `layer boat -~> [hidden]`, also after the editor's own eye click writes it. The frame shows `layer boat`, and the eye already tells you it's hidden. Expected: `amber` eye-off, and the `[hidden]` / `[locked]` flag comments left out of the row text (still in the file).
4. 🟡 **The eye click scrolls the code to the top** (seen while checking R004-7). Code scrolled to `layer sea`, clicking its eye hides the layer and the panel jumps back to line 1. Expected: the scroll and the caret stay where they were.
5. 🟡 **The palette stays on `colors.sky`** (R002-8). After a double-click on the `sky` chip, selecting a layer (`shadows`), or even opening another file, still leaves the compact palette's header on `colors.sky` with its hex. Expected: as with shapes (R002-6 fix), selecting a layer or opening a file ends the variable edit (`COLOR` + the current color).
6. 🟡 **No "nothing selected" state** (R004-1, also R004-3). A file opens with its line-1 comment row selected. The properties read `PROPERTIES · line 1 / Nothing to edit on this row` and the artist bar reads `Nothing to edit on this row`. Escape and a click on empty canvas don't change that. Expected (R004 frames): `PROPERTIES · nothing selected`, `canvas (480, 320)`, hint `hover a layer to preview it` (ES `nada seleccionado`, `pasa sobre una capa para verla`); artist bar `nothing selected | hover a layer to preview it`. The string isn't in `lib/` yet.
7. 🟡 **Gradient fills in the properties panel are raw text** (R002-2, R002-3, also R002-9). `fill gradient([#f2a65a, #fadbbd], 90` is cut at the panel edge. While the fill editor is open, nothing in the properties marks it. With steps there's no `steps [5]` field and no hint. Expected: a gradient chip + `gradient(…)` field, `amber` while the editor is open, and with steps on, `steps [5]` and the hint `line 26 · 5 flat bands, no blending`.
8. 🟡 **`let colors` has no properties** (R002-4, R003-5). Selecting the row, or double-clicking a variable chip to edit it, shows `PROPERTIES · line 6 / Nothing to edit on this row`. The chip's double-click also leaves the previous row selected. Expected: row 6 selected, `colors · line 6`, `sky [■ #f2a65a] sea [■ #2d6a8f] +5`, hint `changing sky repaints its 1 shape` (or `click a swatch: color popover`).
9. 🟡 **Stabilizer feedback is only partly built** (R003-4, R003-2). The string, the radius circle and the ink at the tip are all there. Not there: the tooltip `ink follows the string / release: catches up to the cursor`, `pressure 0.71` in the status bar during a tablet stroke (drive sent pressure 0.6–0.75), the dotted `mute` line of the raw cursor path, and the 5×5 `sky` cursor square at the string's end (lasso too).
10. 🟡 **Trace preview: the two halves differ** (R003-6). The `original` half shows the PNG's transparent pixels on the dark preview face, and the `traced` half on white paper. The original is also blocky: it's drawn at UI resolution, while the color-content rule says user pictures show smooth at screen resolution. Expected: both halves on the same ground (the 3 px checker for alpha), and the original at full screen resolution.
11. 🟡 **The cursor position goes stale** (R003-1, R002-17). During a select drag the status bar keeps the press point (`x 360 y 91`) instead of following the pointer (frame `x 338 y 140`). On `many-3400.lzdraw` (near state, `amber`), hovering the canvas shows no position at all, though about 130 px of bar is free. On `layers.lzdraw` hovering shows it.
12. 🟡 **The completion popup covers the PROPERTIES header** (R002-20, R002-22). It ends at y 289 and the strip starts at 282, so `PROPERTIES · line 12` shows as `PROPER… ine 12`. Expected: it ends above the properties strip (shorter list, or it flips above the row).
13. 🟡 **Layer properties** (R002-7). The opacity slider still has only the ◧ icon and no `opacity` label. The header reads `PROPERTIES · Layer · line 9`, capitalized, where other kinds read `rect`, `circle`. Expected: an `opacity` label, and `layer shadows · line 9`.
14. 🟡 **Names and wraps in the color popover** (R002-24, R002-25). The suggested name for a literal is still generic: `as colors.color_2` and the name field `color`; it should come from the hue (`cream_2`). `writes ■colors.color …` ends with a stray `…`. In Spanish, `era ahora / anillo: / tono` still wraps onto three lines. Also, the artist compact palette opens collapsed when the coder sidebar was collapsed, while the frames show it open. The two should keep separate states.

---

## R002 – the first check's 🔴 / 🟡 frames

| Frame | Verdict | Shot | Notes |
|---|---|---|---|
| R002-1 shape menu | ✅ | [menu](../Art/Checks/v1.2/R002-01-shape-menu.png), [Fill ▸](../Art/Checks/v1.2/R002-01-fill-sub.png) | Ends at y 338, above the status bar. Fill ▸ flips left inside the window. The extra items stay a design proposal (draw them into the frame). |
| R002-2 fill editor linear | 🟡 | [shot](../Art/Checks/v1.2/R002-02-fill-linear.png) | Only the dashed box + gradient handles ✅. Item 7. |
| R002-3 stepped | 🟡 | [shot](../Art/Checks/v1.2/R002-03-fill-stepped.png) | 5 bands on bar and canvas ✅. Item 7. |
| R002-4 palette sidebar | 🟡 | [shot](../Art/Checks/v1.2/R002-04-palette-sidebar.png) | Variable edit, kind tabs, `in 1 shape` ✅, 55% ✅. Item 8. |
| R002-5 artist mode | ✅ | [layer](../Art/Checks/v1.2/R002-05-artist-mode.png), [sky](../Art/Checks/v1.2/R002-05-artist-sky.png) | Nothing clips, Fill button is an icon, the fill chip follows the selection, 125%. Small: the bar says `rectangle` where the menu and the coder properties say `rect`. |
| R002-6 artist + drawer | ✅ | [shot](../Art/Checks/v1.2/R002-06-artist-code-drawer.png) | The click unfolds `sun`, the `circle` row gets `moss` + bar, the drawer (header at y 196) selects line 10, 75%. |
| R002-7 blend coder | 🟡 | [shot](../Art/Checks/v1.2/R002-07-blend-coder.png) | Dashed layer box ✅, 90% ✅. Item 13. |
| R002-8 blend artist | 🟡 | [shot](../Art/Checks/v1.2/R002-08-blend-artist.png) | One `amber` ring, `MUL` badge after the name ✅. Item 5. |
| R002-9 Spanish menu + radial | ✅ | [menu](../Art/Checks/v1.2/R002-09-spanish-menu-fill.png), [editor](../Art/Checks/v1.2/R002-09-spanish-fill-editor.png) | `Relleno ▸` flips left and stays inside the window. The editor sits over the code with its tail to the shape, and the radial handles show. The fill field is item 7. |
| R002-10 tutorial 1 | ✅ | [shot](../Art/Checks/v1.2/R002-10-tutorial-1.png) | Sidebar opens, tool strip ringed, 480×320. |
| R002-11 tutorial 2 | ✅ | [shot](../Art/Checks/v1.2/R002-11-tutorial-2.png) | Shot with no shape (step 1 skipped), so the ring falls back to the whole sidebar; the first check's shape case holds. |
| R002-12 tutorial 3 | ✅ | [shot](../Art/Checks/v1.2/R002-12-tutorial-3.png) | Brush and canvas ringed, dotted ghost. |
| R002-13 tutorial 4 | ✅ | [step 4](../Art/Checks/v1.2/R002-13-tutorial-4.png), [done](../Art/Checks/v1.2/R002-13-tutorial-done.png) | A `Path` stroke moves on to step 4 by itself. The brush inspector takes two rows, the canvas row stays at y 263, and the counter stays. |
| R002-14 tutorial ES | ✅ | [shot](../Art/Checks/v1.2/R002-14-tutorial-1-es.png) | |
| R002-15 new drawing | ✅ | [shot](../Art/Checks/v1.2/R002-15-new-drawing.png) | 480×320 first and default. |
| R002-16 new drawing ES | ✅ | [shot](../Art/Checks/v1.2/R002-16-new-drawing-es.png) | `Sin título`. (The `Call NewDrawing` backdoor skips the dialog and makes `Untitled` 512×512 white even in Spanish: a drive shortcut, not the UI.) |
| R002-17 object counter | 🟡 | [near](../Art/Checks/v1.2/R002-17-near.png) | `amber` near state, never shortened ✅. Item 11. |
| R002-18b pen pop-up | ✅ | [live readout](../Art/Checks/v1.2/R002-18b-pen-pad-readout.png) | Now verifiable: mid-stroke `pressure 0.62 · tablet` in `sky` and the dashed `sky` line on the curve. The pressure tabs stay on `Trackpad` while the source is a tablet; fine. |
| R002-19b pen pop-up ES | ✅ | [shot](../Art/Checks/v1.2/R002-19b-pen-settings-artist-es.png) | `AJUSTES DE PLUMA · tinta`, `mi tinta`, the bar fits. |
| R002-20 autocomplete | 🟡 | [shot](../Art/Checks/v1.2/R002-20-autocomplete.png) | Layout, kinds, footer, held errors ✅. Item 12. |
| R002-21 signature hint | ✅ | [shot](../Art/Checks/v1.2/R002-21-signature-hint.png) | |
| R002-22 `math.s` | 🟡 | [shot](../Art/Checks/v1.2/R002-22-member-completion.png) | `math.sin(x) / Sine of x in radians.` ✅. Item 12. |
| R002-23 popover from a row | ✅ | [shot](../Art/Checks/v1.2/R002-23-color-popover-row.png) | `amber` swatch ring, no hint. It opens right of the literal, over the canvas, as the first check accepted. |
| R002-24 popover + save | 🟡 | [before](../Art/Checks/v1.2/R002-24-before.png), [after](../Art/Checks/v1.2/R002-24-after.png) | Keeps its size and place ✅. Item 14. |
| R002-25 popover artist ES | 🟡 | [shot](../Art/Checks/v1.2/R002-25-color-popover-artist-es.png) | `relleno [■ colors.sun]` opens it, `sun` current, tail up ✅, `como colors.sun_2` ✅. Item 14. |

## R003 – gestures (new) and the rest

| Frame | Verdict | Shot | Notes |
|---|---|---|---|
| R003-1 select, moving 2 | 🟡 | [mid-drag](../Art/Checks/v1.2/R003-01-select-move-ghost.png) | Now captured: the dotted `mute` ghost at the old place, the `move -30, +40` tooltip (`amber` numbers), the rows updating live, `2 objects · move (-30, +40)`. Matches the frame. Item 11. |
| R003-2 lasso in progress | ✅ | [ants](../Art/Checks/v1.2/R003-02-lasso-ants.png), [released](../Art/Checks/v1.2/R003-02-lasso-released.png) | Marching ants, `mute` closing line, `sky` string + dotted radius circle, tip square, the shape inside gets a dashed box, `1 inside`. On release, catch-up selects hull + mast. The cursor square is item 9. |
| R003-3 lasso fill | ✅ | [shot](../Art/Checks/v1.2/R003-03-lasso-fill.png) | `Path … --close` makes a `polygon` row in the current layer with the `new` badge, selected, `lasso fill · ■ #1d1d1f`. |
| R003-4 brush string | 🟡 | [mid-stroke](../Art/Checks/v1.2/R003-04-brush-string.png), [after](../Art/Checks/v1.2/R003-04-brush-after.png) | The `sky` string, the dotted radius circle, ink only up to the tip, and the bar `brush · boat \| pen [ink ▾] ⚙ \| stabilizer Off String Average 24 px` ✅. Items 2 and 9. |
| R003-5 wheel + harmonies | 🟡 | [shot](../Art/Checks/v1.2/R003-05-let-colors-props.png) | The popover (first check's fix) holds. Item 8. |
| R003-6 trace dialog | 🟡 | [empty](../Art/Checks/v1.2/R003-06-trace-empty.png), [line art](../Art/Checks/v1.2/R003-06-trace-dialog.png), [photo](../Art/Checks/v1.2/R003-06-trace-photo.png) | Built: 450×266 at y 34 over the dim, `TRACE IMAGE · girl.png`, before/after split with grip, `original` / `traced` chips, `girl.png · 255×426` + `Change…`, the drop hint, MODE Photo / Line art, the five sliders, show lines / fills, stats with `amber` numbers, Insert as layer (disabled from Home) + Open as drawing. Photo mode hides the options the engine doesn't have, as the batch says. Item 10. |

## R004 – layer hover preview (new)

| Frame | Verdict | Shot | Notes |
|---|---|---|---|
| R004-1 coder card | 🟡 | [shot](../Art/Checks/v1.2/R004-01-coder-hover-card.png) | 152 wide, 140×93 preview of the layer alone on `#f6d9a8`, name + `11 objects` (3 circles + 8 rays), `click: select layer`, 5-step tail at the row, row hover. Scroll hides it ✅. Item 6. |
| R004-2 Alt isolation | ✅ | [shot](../Art/Checks/v1.2/R004-02-coder-alt-isolation.png) | Drawing at 70% `void`, the layer drawn again on top in its real blend, chip `sun · 11 objects` under the canvas. In coder mode the card sits beside the row and leaves the canvas clear. |
| R004-3 artist card | 🟡 | [shot](../Art/Checks/v1.2/R004-03-artist-hover-card.png) | Card left of LAYERS, tail on the right. Items 1 and 6. |
| R004-4 artist card + isolation | 🔴 | [shot](../Art/Checks/v1.2/R004-04-artist-card-and-isolation.png) | The dim and the chip work, but the bright layer is under the card. Item 1. |
| R004-5 Spanish | ✅ | [shot](../Art/Checks/v1.2/R004-05-coder-spanish.png) | `11 objetos`, `clic: elegir la capa`, `CAPAS`, `PROPIEDADES`. |
| R004-6 nested `sun/rays` | ✅ | [shot](../Art/Checks/v1.2/R004-06-nested-layer.png) | `sun/rays`, `8 objects` from the `for`, `SCR 70%` badge line, only the rays, flipped above the status bar. The card draws the rays in Normal (pink) while the canvas and Alt isolation screen them (pale). Kept on purpose, see the Style Guide note: a screen layer would almost vanish on the light background, and the badge names the mode. |
| R004-7 hidden `boat` | 🟡 | [coder](../Art/Checks/v1.2/R004-07-hidden-layer.png), [artist + Alt](../Art/Checks/v1.2/R004-07-artist-hidden-alt.png), [rows](../Art/Checks/v1.2/R004-extra-hidden-rows.png) | The preview shows the boat with the `amber` `hidden` chip and `only shown here`, the canvas doesn't, and Alt doesn't isolate a hidden layer ✅. Items 3 and 4. |

## Notes on drive

- `Press` / `Move` / `Release` / `Path` were enough for every gesture frame. `Path … --close` draws a real lasso loop.
- Twice, a gesture command printed `Press failed` / `Path failed` and the same call worked on a retry. Once, the first `Press --pressure` stroke with the brush drew nothing (also fine on a retry). I couldn't reproduce either; worth a look if testers see it.
- `Move --steps 5` inside one drive command reaches the brush as one sample per command (a 3-command stroke gave 3 points), so a smooth test stroke needs many `Move`s or `--samples`.
