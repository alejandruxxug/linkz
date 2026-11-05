# 🎨 Linkz: Style Guide

## 🧭 Direction in one line

**An Aseprite-like pixel editor in near-black, around smooth vector art.** The UI is built like 8-bit pixel art (half resolution, scaled ×2, hard pixels, 1 px lines, bevels, no gradients) but colored dark: near-black panels, one lizard-green accent, amber for focus. The drawing on the canvas stays full-resolution and anti-aliased.

## 🕹️ 8-bit UI rules

| Rule | Detail |
|---|---|
| Resolution | UI drawn at half the window resolution, scaled ×2 nearest-neighbor by the host. No anti-aliasing anywhere in the UI. |
| Palette | Only the 16 Linkz Dark colors below. No gradients, blur or soft shadows; shade with the darker/lighter neighbour or a 2×2 checker dither. |
| Font | The 7×13 bitmap font (`lib/gfx/font_data.liz`), scale 1 in the UI buffer (it's already ×2 on screen). |
| Outlines | 1 px `void` outline on panels, buttons, fields, tooltips. |
| Buttons | `raised` face, 1 px `void` outline, 1 px bevel inside: `bevel` top/left, `coal` bottom/right; hover = `bevel` face; pressed/active = `moss` face, bevel swapped, content moved 1 px down-right, active tool icon `lizard`. |
| Selection | Selected row: `moss` background, 2 px `lizard` bar on the left, `bright` text. Focused field: `amber` outline. On the canvas: dashed 1 px `lizard` box (3 on, 2 off) and 5×5 anchors (`lizard`, `bright` for size handles, 7×7 `amber` when hovered), each with a `void` outline. |
| Icons | 8×8 pixel icons as lizard data (one color + transparent, tinted when drawn): ○ □ ╱ ⬠ ～ ✎ 🖌 T ƒ ▤ layer, eye, lock, chevrons. |
| Caret | Solid 1 px `amber` bar (nothing animates while idle, so no blink). |
| Tooltips | `row` box, `void` outline, `bevel` top line, 2 px `void` drop shadow below-right. |
| Canvas | Full-res drawing inside a 1 px `void` frame on a `coal` / `#161815` 16 px checker surround (the one off-palette tint, kept because it is only a texture). |
| Color content ✅ (R002) | The one exception to "only the 16 colors, no gradients": UI parts that *show a user color* — the HSV square, hue and alpha strips, gradient bars, variable chips, swatches, layer thumbnails, the pen preview — show the exact colors, smooth, at screen resolution. Everything around them (1 px `void` frame, labels, thumbs, markers) stays 8-bit. Alpha shows over a `#d4dacd`/`#a0ab9b` checker (2–3 UI px). |
| Spotlight ✅ (R002) | Tutorial focus: everything else under `void` at 62%, the target ringed 1 px `amber` (focus color). Same 62% `void` dim behind modal dialogs. |

## 🎨 Palette: Linkz Dark

Replaces Lizard 16 (developer: "similar in UI to Aseprite, but black and dark, not the 8-bit palette"). In Figma as the variable collection **Linkz Dark** and paint styles `Linkz Dark/<name>`.

| Name | Hex | Role |
|---|---|---|
| void | `#0b0c0b` | 1 px outlines, bevel shadow, canvas frame, drop shadows |
| coal | `#131512` | Window background, tab bar, inset fields, canvas surround |
| panel | `#1a1d19` | Docked panels, toolbar, status bar |
| row | `#232721` | Hovered row, tooltips, menus |
| raised | `#2e332c` | Button face, panel header strips, badges |
| bevel | `#434a40` | Bevel highlight (top/left), splitters, indent guides |
| dim | `#5f685b` | Disabled icons, unlocked lock, "add" row (decorative: never body text) |
| mute | `#a0ab9b` | Secondary text, panel titles, comments, labels |
| text | `#d4dacd` | Body text, identifiers, icons |
| bright | `#f2f5ec` | Selected-row text, tooltip titles |
| moss | `#163a20` | Selected row, active tool face, current menu item |
| lizard | `#4fc46a` | The accent: selection bar/outline, anchors, active tool icon, "saved", parameter names |
| amber | `#f2b33d` | Focus outline, caret, hovered anchor, numbers in code, unsaved dot |
| ember | `#f36c5b` | Errors |
| sky | `#72b8e6` | Keywords in code, links, info |
| peach | `#eb9a6f` | Strings in code |

Contrast (WCAG, body text needs ≥ 4.5:1):

| Text \ on | coal | panel | row | raised | moss |
|---|---|---|---|---|---|
| bright | 16.7 | 15.4 | 13.8 | 11.7 | 11.5 |
| text | 12.9 | 11.9 | 10.6 | 9.0 | 8.9 |
| mute | 7.7 | 7.1 | 6.4 | 5.4 | 5.3 |
| lizard | 8.3 | 7.6 | 6.8 | 5.8 | 5.7 |
| amber | 9.9 | 9.2 | 8.2 | 6.9 | 6.8 |
| sky | 8.5 | 7.9 | 7.0 | 6.0 | 5.9 |
| peach | 8.2 | 7.6 | 6.8 | 5.8 | 5.7 |
| ember | 6.2 | 5.8 | 5.1 | 4.4 | 4.3 |
| dim | 3.2 | 2.9 | 2.6 | 2.2 | 2.2 |

Every text color passes 4.5:1 on every background it is used on. `ember` stays off `raised` and `moss` (4.4 / 4.3): error rows keep the `panel` background with an `ember` bar. `dim` is only for non-text marks and disabled states.

Code colors: `sky` keywords (`let`, `layer`, `for`), `amber` numbers, `peach` strings, `mute` comments and punctuation, `text` identifiers (`bright` on the selected row), colors drawn as a 6×6 swatch (1 px `void` outline) before the literal, errors `ember`.

Old → new names (for the rule table above): ink → `void`, night → `coal`, panel-dark → `panel`, panel → `row`, shadow → `void`/`coal` (bevel shadow), stone → `mute`, light → `bevel`, paper → `bright`, leaf → `lizard`, sun → `amber`, butter → `amber`, rust → `peach`, sea → `sky`. Selection is now `moss` background + 2 px `lizard` bar + `bright` text (not `leaf` + `ink`). Tooltips are `row` with a `void` outline, a `bevel` top line and a 1 px `void` drop shadow.

## 🧩 Figma

| File | Pages |
|---|---|
| **Linkz** (Figma desktop, edited through the Figwright plugin) | 🖼️ Inputs (direction + references) · 🎨 Style (palette, font 1× and 2×, components, 8×8 icons) · 🧪 Reviews (one section per review while it is open) · ✅ Approved (the frames to build from) |

⚠️ 2026-11-03: the 🧪 Reviews page holds only **R004 – Layer hover preview** (`25:4344`, open). 2026-11-04: **R005 – Brush creator, brush files, video tooltips** (`29:4808`, approved 2026-11-04) sits under it at y 2720. The R001, R002 and R003 review sections (`3:2`, `6:1257`, `9:2302`) were moved, not copied, to ✅ Approved (they now have new `21:*` ids), so the old review node ids in the tables below are history only. Their pictures are still in `docs/Art/Reviews/R002/` and `R003/`. If the team wants the review copies back on 🧪 Reviews, clone the three ✅ Approved sections there.

Mockups are drawn at window size (1024×720) with every UI pixel as a 2×2 block. Each UI part is its own layer, rendered pixel-exact from the real 7×13 font and 8×8 icons. The canvas drawing is vector and smooth, and home thumbnails are rendered drawings. Layer names say what each part is (`row line 10 · circle (selected)`, `properties panel`, `tool strip`, `card menu ⋯ (open)`), so the builder reads the frame tree. The generator that produced them works in UI pixels with the numbers below.

Approved frames (page ✅ Approved, node ids):

| Frame | Node | Review copy |
|---|---|---|
| Editor – main | `4:829` | `3:3` |
| Editor – new row, inline tool strip | `4:873` | `3:98` |
| Editor – pencil drawing (anchors, one hovered) | `4:918` | `3:99` |
| Editor – brush + pen picker | `4:964` | `3:100` |
| Editor – tooltip + error row | `4:1028` | `3:101` |
| Editor – Spanish UI | `4:1075` | `4:418` |
| Splash banner (200×80 at 4×) | `4:1119` | `3:103` |
| Home – gallery | `4:1122` | `3:104` |
| Home – empty (first launch) | `4:1184` | `3:105` |
| Home – gallery, Spanish | `4:1194` | `4:748` |

The ✅ Approved page also has section **R001 – Editor views** `21:2692` (the review frames, moved from 🧪 Reviews; same pictures as the table above).

Approved sections for R002 and R003 (page ✅ Approved). Each PNG in `docs/Art/Approved/` is named after its node id, exported at 1× (1024×720 unless noted). The app's v1.1 check against them is in `docs/Reviews/2026-11-03 v1.1 design check.md`.

**R002 – Menu, fill, palette, artist mode, tutorial** (section `21:3120`):

| # | Frame | Node | PNG | Review node (history) |
|---|---|---|---|---|
| 1 | Shape menu – right-click on canvas, Fill ▸ open | `21:3121` | `21-3121.png` | `6:1258` |
| 2 | Fill editor – linear gradient on the sky | `21:3151` | `21-3151.png` | `6:1259` |
| 3 | Fill editor – stepped gradient, 5 bands | `21:3181` | `21-3181.png` | `6:1260` |
| 4 | Palette sidebar – editing variable `colors.sky` | `21:3214` | `21-3214.png` | `6:1261` |
| 5 | Artist mode – layers.lzdraw | `21:3240` | `21-3240.png` | `6:1262` |
| 6 | Artist mode – code drawer open | `21:3285` | `21-3285.png` | `6:1263` |
| 7 | Blend modes – coder, `shadows` in Multiply | `21:3331` | `21-3331.png` | `6:1264` |
| 8 | Blend modes – artist, dropdown open | `21:3362` | `21-3362.png` | `6:1265` |
| 9 | Spanish – shape menu + fill editor (radial, hard stop), two windows | `21:3414` | `21-3414.png` (2128×720) | `6:1266` |
| 10 | Tutorial 1 – Make your first shape | `21:3472` | `21-3472.png` | `6:1267` |
| 11 | Tutorial 2 – Change its color | `21:3487` | `21-3487.png` | `6:1268` |
| 12 | Tutorial 3 – First brush stroke | `21:3504` | `21-3504.png` | `6:1269` |
| 13 | Tutorial 4 – Peek at the code (done) | `21:3522` | `21-3522.png` | `6:1270` |
| 14 | Tutorial 1 – Spanish | `21:3541` | `21-3541.png` | `6:1271` |
| 15 | New drawing dialog – object budget | `21:3826` | `21-3826.png` | `7:1963` |
| 16 | New drawing dialog – Spanish | `21:3890` | `21-3890.png` | `7:2025` |
| 17 | Status bar – object counter (normal, near, over) | `21:3556` | `21-3556.png` (1024×328) | `6:1272` |
| 18 | Pen settings – coder (old panel, superseded by 18b) | `21:3566` | `21-3566.png` | `6:1273` |
| 19 | Pen settings – artist ES (old panel, superseded by 19b) | `21:3598` | `21-3598.png` | `6:1274` |
| 18b | Pen settings pop-up – coder (test pad, stabilization, code) | `21:4003` | `21-4003.png` | `9:2304` |
| 18b-2 | Pen settings pop-up – empty test pad | `21:4091` | `21-4091.png` (982×464) | `9:2306` |
| 18c | Pen picker – ⚙ per row | `21:3954` | `21-3954.png` | `9:2303` |
| 19b | Pen settings pop-up – artist mode, Spanish | `21:4038` | `21-4038.png` | `9:2305` |
| 20 | Autocomplete – typing `ci` | `21:3647` | `21-3647.png` | `6:1275` |
| 21 | Autocomplete – signature hint | `21:3673` | `21-3673.png` | `6:1276` |
| 22 | Autocomplete – member `math.s` | `21:3699` | `21-3699.png` | `6:1277` |
| 23 | Color popover – from a code-row swatch (square picker; the wheel replaces it, see R003-5) | `21:3725` | `21-3725.png` | `6:1278` |
| 24 | Color popover – fill field, save as variable | `21:3753` | `21-3753.png` | `6:1279` |
| 25 | Color popover – artist mode, Spanish | `21:3780` | `21-3780.png` | `6:1280` |

**R003 – Select, lasso, stabilizer, color wheel, trace** (section `21:4094`):

| # | Frame | Node | PNG | Review node (history) |
|---|---|---|---|---|
| 1 | Select tool – moving 2 objects | `21:4095` | `21-4095.png` | `9:2307` |
| 2 | Lasso select – loop in progress, string stabilizer | `21:4123` | `21-4123.png` | `9:2308` |
| 3 | Lasso fill – loop became a polygon row | `21:4149` | `21-4149.png` | `9:2309` |
| 4 | Brush – pulled-string stabilizer (stroke in progress) | `21:4176` | `21-4176.png` | `9:2310` |
| 5 | Color popover – wheel + harmonies | `21:4223` | `21-4223.png` | `9:2311` |
| 6 | Trace image dialog – line art, before/after (design only, not built) | `21:4251` | `21-4251.png` | `9:2312` |

Style page: palette section `4:501` (swatches bound to the **Linkz Dark** variables `color/<name>`, paint styles `Linkz Dark/<name>`), font `4:502`, components `4:503` (`button/*`, `tool/*`, `field/*`, `row/*`, `tooltip`), icons `4:504` (components `icon/<name>`). R002 components ✅ in section `7:2125`: `menu/shape`, `menu item/{normal,hover,disabled,danger,submenu ▸,current}`, `gradient bar/{smooth,stepped}`, `color stop/{normal,selected,hard,ghost}`, `segmented tabs`, `hsv picker`, `variable chip/{solid,linear,radial,stepped}`, `color popover`, `fill editor`, `layer row (thumbnail)/{normal,selected + blend badge}`, `row/layer with blend badge`, `blend dropdown (open)`, `slider/{single,range}`, `coach-mark`, `status counter/{normal,near,over}`, `autocomplete popup`, `pen settings panel`, and 30 more icons (`icon/up … icon/settings`, incl. the home icons).

R002 – Menu, fill, palette, artist mode, tutorial ✅ (review ids, history: 🧪 Reviews section `6:1257`, now ✅ Approved `21:3120`, see above; generator `docs/Art/Figma-gen/r002.py` + `r002_split.py`; PNGs in `docs/Art/Reviews/R002/`):

| # | Frame | Node |
|---|---|---|
| 1 | Shape menu – right-click on canvas, Fill ▸ open | `6:1258` |
| 2 | Fill editor – linear gradient on the sky | `6:1259` |
| 3 | Fill editor – stepped gradient, 5 bands | `6:1260` |
| 4 | Palette sidebar – editing variable `colors.sky` | `6:1261` |
| 5 | Artist mode – layers.lzdraw | `6:1262` |
| 6 | Artist mode – code drawer open | `6:1263` |
| 7 | Blend modes – coder, `shadows` in Multiply | `6:1264` |
| 8 | Blend modes – artist, dropdown open | `6:1265` |
| 9 | Spanish – row menu + fill editor (radial, hard stop), two windows | `6:1266` (`6:1281`, `6:1282`) |
| 10–13 | Tutorial 1–4 (shape, color, brush, code + lizard) | `6:1267` `6:1268` `6:1269` `6:1270` |
| 14 | Tutorial 1 – Spanish | `6:1271` |
| 15 / 16 | New drawing dialog (object budget) EN / ES | `7:1963` / `7:2025` |
| 17 | Status bar – object counter, 3 states | `6:1272` |
| 18 / 19 | Pen settings – coder EN / artist ES (old panel, kept for history; replaced by 18b / 19b if approved) | `6:1273` / `6:1274` |
| 18b / 19b ✅ | Pen settings **pop-up** – coder EN (test pad with 3 strokes, live pressure 0.62, stabilization, code scrolled to the changed line) / artist ES (opened from the ⚙ in the inspector) · row 2 of the section | `9:2304` / `9:2305` |
| 18b-2 ✅ | Pen settings pop-up alone – empty test pad ("draw here to try", readout —) | `9:2306` |
| 18c ✅ | Pen picker – ⚙ on every row (hover = button + tooltip), new "soft ink" row, footer "double-click brush: its settings" | `9:2303` |
| 20–22 | Autocomplete – `ci`, signature hint, `math.s` | `6:1275` `6:1276` `6:1277` |
| 23–25 | Color popover – code-row swatch, fill field (save as variable), artist ES | `6:1278` `6:1279` `6:1280` |

### 📏 UI metrics (UI pixels: the half-res buffer, shown ×2)

| Part | Metric |
|---|---|
| Window | 512×360 (1024×720 on screen) |
| Font | 7×13 cell, 7 px advance, ascent 11. In an 18 px row the text top is at +3 |
| Toolbar | height 16, bottom 1 px `void` line. Tool buttons 14×14 at y 1, x from 3, step 15. Groups separated by 4 px with a 1×10 `bevel` mark. Right side: `Export PNG`, `Save`, language buttons, 14 high, sized from text (icon at +3, label at +13, width = text + 18), 2 px apart (6 px before language) |
| Tab bar | height 14 (y 16–30). Active tab: `panel` face, 1 px `lizard` top line, 1 px `void` sides, width = name + 26 |
| Status bar | height 16 (y 344–360), 1 px `void` top line. Text at +2, icons 8×8 at +4. Parts are 6 px + a 1×8 `bevel` separator + 7 px apart. Cursor position, zoom and language are right-aligned |
| Layers panel | x 0–320. Header 13 `raised` + 1 `void` line. Rows start at y 45 |
| Row | **pitch 18**. Icons at +5, text at +3. Extra **4 px gap above every `layer` row**. Blank file lines are 4 px. `⊕ add` row is 10 px. Selected: `moss` + 2 px `lizard` bar on the left. Error: 2 px `ember` bar + `error` icon + message row of 16 below |
| Row columns | eye x 4, lock x 14, tree x 24. **Indent step 10 per depth**: chevron at 24 + 10d, kind icon at 34 + 10d, code at 46 + 10d. Indent guides are 1 px `bevel` at 27 + 10k. Count badge is right-aligned 8 from the edge. Clipped code ends in `»` |
| Code | color literal = 6×6 swatch (1 px `void` + 4×4 color) and 2 px gap before the text |
| Pinned canvas row | y 263 (props top − 19), `coal` background, 1 px `void` line above |
| Scrollbar | 4 wide at x 315, `coal` track, `bevel` thumb. Shown when rows are hidden |
| Splitter | 2 px at x 320 (`void` + 1 px `bevel`). Canvas area x 322–512, 16 px `coal`/`#161815` checker |
| Properties panel | y 282–344. Header 15 `raised` + 1 `void` line, title at +2. Fields row at +22, hint at +45 |
| Field | height 18, 1 px `void` outline (`amber` when focused), `coal` inset with 1 px `void` inner top line, text at x+4 / y+3. Point = `(` [x] `,` [y] `)`. Color = 10×10 swatch at +3 + hex |
| Button | 1 px `void` outline + 1 px bevel. Tool icon centered. Pressed/active content moves 1 px down-right |
| Inline tool strip (new row) | row 22 high, 8 tool buttons 14×14 at y + 4, step 15, starting at the code column − 4, then a `mute` hint |
| Canvas overlay | dashed 1 px `lizard` box (3 on, 2 off). Anchors 5×5 + 1 px `void` (size handles `bright`), hovered 7×7 `amber`. Drawing at 75%, 1 px `void` frame |
| Tooltip | `row` fill, `void` outline, `bevel` top line, 2 px `void` shadow, padding 6, line pitch 15 |
| Pen picker | 230 wide under the brush button, header 13, rows 20, preview 50×16 inset in a `void` frame (full-res strokes) |
| Home | top bar 22 (search field 212×16 at x 150). Section title + 1 px `bevel` rule. Working-on row: `+ New drawing` card 80 wide (dashed `bevel` outline, label wraps when too long), 3 recent cards 128 wide, 8 apart, thumb 76 high + 32 text. Grid: 5 cards 91 wide (thumb 50 + 28 text), 8 apart, 6 between rows. Hover card: `lizard` outline + `⋯` button. Card menu rows 14 (Delete in `ember` after a separator) |

R002 metrics ✅ (approved 2026-11-03, same units):

| Part | Metric |
|---|---|
| Coder layout with palette | code panel 0–256, splitter 2, canvas 258–404, splitter 2, sidebar 406–512 (106, inner 96 at +5). Sidebar collapsed = 14 px rail at 498 (variables 10×10 step 12, a `bevel` mark, recent), canvas grows to 258–496. Canvas zoom = largest of 62.5/45/37.5/27.5% (UI) that leaves 10 px around: 55% open, 90% rail |
| Properties at 256 | fields that don't fit wrap to a second 20 px row (the hint goes) |
| Palette sidebar | header 13 (title = what is being edited: `COLOR`, `STOP 2`, `colors.sky`) + collapse chevron. Variable kind tabs (icons only, 3 × 32) when a variable is selected. SV square 96×52 (44 with tabs), hue 96×6, alpha 96×6 (5 px apart), hex field 57×16 + alpha field, `RECENT` 8 swatches 10×10 step 12, `VARIABLES` grid 3 columns × 32 (chip 29×12 + name under it, cell 27 high; selected cell `moss`), "in N shapes" hint `lizard`, `+ Variable` button 96×16 at the bottom |
| Variable chip | solid fill, linear (smooth), radial (over checker), stepped (4 bands); 1 px `void` outline. Gradients and colors are interchangeable |
| Context menu | rows 14, header 13 `raised` (`circle · line 10`), icon at +6, label at +18, shortcut right-aligned 6 from the edge, `chevron-right` for submenus, separators 4 (1 px `bevel`). Hover: `moss` + 2 px `lizard` bar, `bright` text, `lizard` icon. Disabled: `dim` icon/label. Danger (Remove): `ember` label/icon, on hover keeps `row` (ember fails 4.5:1 on `moss`). Current: `lizard` check. Flips left/up at the window edge; submenu opens on the side with room. Submenu can end with a swatch row (10×10 step 12) |
| Fill editor | floating, 186 wide, beside the shape (over the code panel when the canvas is narrow), 5-step pixel tail toward the shape. Header 13, kind tabs (3 × 59, 14 high), gradient bar 174×14 (`void` frame, checker under alpha), stops below it, stop row (color field, `at` field, remove), soft/hard segmented + `steps` field with −/+, code preview inset (13 px lines, wraps at 166), two `mute` hint lines, `+ Save as variable…` |
| Color stop | 9×9 marker: 3-row pointer + 7×5 color body, `void` outline; selected = `bright` outline + `void` ring; hard = 1-px stem pointer + 1 px `bright` line through the bar; ghost (add) = `dim` outline + `+` |
| Gradient handles | linear: 1 px `bright` line (with `void` sides) from start (5×5 `bright`) to end (5×5 `lizard`), stops as 5×5 color dots on it. Radial: 5×5 center, 1 px `bright` pixel circle (+ `void` ring) at the radius, radius handle 5×5 `lizard` |
| Color popover | 150 wide, beside the clicked color with a tail: `VARIABLES` 4 columns × 35 (chip 32×10 + name), current `moss`, keyboard focus 1 px `amber` ring; `RECENT` 14 swatches 9×9 step 10; SV 140×36, hue 140×5, hex + alpha + eyedropper; `+ Save as variable…` (opens inline name field + Save); `writes colors.sun` preview |
| Artist mode | toolbar (`</> Code` toggle, mode label), tab bar, inspector context bar 20 high (y 30–50, fields 16 high), tool column 20 wide (tools 14×14 step 15 from y 53, fill/stroke chips at the bottom), canvas 20–360, right column 362–512: compact palette (header, SV 116×40 + vertical hue/alpha 8×40, hex/alpha/eyedropper, one row of variable chips + recent), then `LAYERS` header (+ add, trash), blend dropdown 82×16 + opacity field, layer rows, `canvas` row pinned at the bottom |
| Layer row (artist) | 22 high: eye +4, lock +14, chevron +24, thumbnail 29×20 at +34 (rendered layer on a 3 px checker), name at +67, blend badge right. Object rows under an open layer are 18 (kind icon +44, name +56) |
| Code drawer | opens under the canvas at y 196 (canvas area shrinks, zoom drops to 75%), full width 0–360, header 13 `</> CODE · file` + lower/close; same rows as coder mode |
| Blend | layer row badge in the count slot (`MUL 80%`); modes Normal │ Multiply, Darken │ Screen, Lighten, Add │ Overlay, each with a 10×10 preview (sky under the shade, result) and a 3-letter code; opacity slider (track 6, thumb 5×12 `bright`) + `%` field |
| Coach-mark | tooltip box, 250 wide, 2 px `lizard` left bar, title `bright`, body `text` at 14 pitch, footer: step dots 4×4 (current `lizard`, others `bevel`), `Skip` (`mute`), primary `Next` (`moss` + `lizard` bevel); 5-step tail toward the target |
| New drawing dialog | 300×214 centered at y 64 over a 62% `void` dim; header 15; name field (`.lzdraw` in `dim`), size presets as buttons (selected = active), custom W × H fields, rule, `objects` icon + "Recommended: up to ~1,800 objects" + `info` (`sky`), `measured on this Mac` (`mute`), Cancel + primary Create |
| Object counter | status bar part after `saved`: `objects N / MAX`, `mute` < 80%, `amber` ≥ 80%, `ember` + `warn` icon over; tooltip on hover. It replaces "steps · ms" in the status bar |
| Pen settings pop-up ✅ (R002 rev, 18b) | floating window 480 wide (coder 285 high, artist 220), over the editor at the left: 1 px `void` outline, **3 px hard `void` drop shadow**, `panel` face, title bar 14 (`raised`, 1 px `lizard` top line when focused, `dim` 2×3 grip dots = drag, `PEN SETTINGS · ink`, `Reset` text button, ✕). Opened from the ⚙ on a pen row, the ⚙ in the brush inspector, or double-clicking the brush tool; Esc / ✕ closes, changes apply live. Three columns split by 1 px `void` lines: **A** pressure graph 84×84 + preset dropdown; **B** five sliders (label 69, track 70, `amber` value) + rule + live readout `pressure 0.62 · trackpad` (value `sky`, matches the dashed `sky` line and a 3×3 `sky` dot on the curve; `—` in `dim` when idle); **C** `STABILIZATION` (Off / String / Average segmented, `length` slider + px, ☑ catch up at the end), input switch Trackpad / Tablet / Mouse and that input's overrides box (`coal` + `lizard` bar). Row 2: **test pad** 300×66 paper (`#fbf8ef`, 1 px `void` frame, 1 px shade line), strokes in the pen's real look, 13×13 trash button in its top-right corner, faint `dim` "draw here to try" centred when empty; strokes there never reach the drawing (`dim` note "pad strokes stay here"); right of it the name field + primary `Save as new pen`. Coder adds the code block: 4 visible lines of the pen hash, scrolled to what you just changed (`row` line + 2 px `amber` bar), 3 px scrollbar |
| Pen settings (old panel, R002 18) | the pen picker turns into this panel (`◂ pens` back), 320 wide: pressure graph 84×84 (`coal`, `bevel` dotted quarters, `dim` diagonal, `lizard` curve, 5×5 points, hovered 7×7 `amber`, `sky` dashed live-pressure line), preset dropdown; sliders (label, track, `amber` value): width (range), opacity (range), taper, smoothing, speed; input switch Trackpad / Tablet / Mouse with that input's overrides in a `coal` box with a `lizard` bar (trackpad: floor); live preview 308×34 on paper; coder mode adds the code it writes (6 lines); name field + primary `Save as new pen` |
| Pen picker ⚙ ✅ (18c) | every row ends in an 8×8 `settings` icon 18 from the right edge: `dim` at rest, `text` on the current pen, a 14×14 hover button on the hovered row + tooltip "Pen settings…". Footer row 14 under a `bevel` rule: "double-click brush: its settings" |
| Autocomplete | 214 wide under the caret, rows 14: kind icon, label with matched letters `lizard`, kind right `mute`; selected `moss` + bar; doc block (signature, one line); footer `Tab insert · Esc close`. Signature hint: tooltip above the row, current parameter `lizard` + 1 px underline, its doc line below |

R003 – Select, lasso, stabilizer, color wheel, trace ✅ (🧪 Reviews, section `9:2302` at y 3060; generator in the designer's scratchpad: `gen_r002b.py` + `r003.py`, built on `r002.py`; PNGs in `docs/Art/Reviews/R003/`):

| # | Frame | Node |
|---|---|---|
| 1 | Select tool – 2 objects (sun core + glow) being moved, rows 10–11 selected, ghost of the old place, handles + rotation handle | `9:2307` |
| 2 | Lasso select – marching-ants loop in progress with the pulled string, objects inside pre-selected (rows hovered), options strip | `9:2308` |
| 3 | Lasso fill – the loop became `polygon([…], #8a6fb0)` line 17 (`new` badge), selected | `9:2309` |
| 4 | Brush – pulled-string stabilizer while drawing (artist): raw cursor path vs smooth ink, string + radius | `9:2310` |
| 5 | Color popover – hue ring + triangle, harmony markers on the ring, 6 harmony modes, swatches + add all | `9:2311` |
| 6 | Trace image dialog – line art mode, before/after split, sliders, toggles, stats, Insert as layer / Open as drawing | `9:2312` |

R003 metrics ✅ (UI pixels):

| Part | Metric |
|---|---|
| Tool order | `select` (V) and `lasso` (L) first, then a separator, then the shape tools (toolbar and the artist tool column). Icons: `select` arrow, `lasso`, `lasso-fill`, `grip`, `move`, `rotate` (8×8) |
| Select | click = select (row `moss` + bar), Shift-click adds, drag on empty canvas = marquee (same dashed `lizard` box), drag a selected object = move all. Bounding box: dashed `lizard`, 8 size handles 5×5 `bright` + `void`, rotation handle = 6 px `lizard` stalk + r 2 ring above the top edge (optional). While moving: the old place as a dotted `mute` ghost, tooltip `move -30, +40` (`amber`), rows update live, properties `2 objects · move ( dx , dy )` |
| Lasso | loop = marching ants (1 px, `bright` 3 / `void` 2, the phase moves while drawing), `bevel` 1 px closing line back to the start, objects fully inside get the dashed box and their rows the hover `row` background. Options strip (properties, 2 rows): Select / Fill segmented (icons) + ☑ catch up; stabilizer Off / String / Average + length slider + `24 px`. Fill mode: the header shows the current color swatch; on release the loop becomes a `polygon` row in the current layer (`new` badge, `lizard`, 3 s), selected |
| String stabilizer | live while drawing (brush, pencil, lasso): the ink / loop is laid only at the **pen tip** (3×3 `bright`, `void` outline); 1 px `sky` string from the tip to the cursor; dotted `sky` circle of radius = string length around the cursor; the raw cursor path is not drawn in the app (shown `dim` in the mockup only). Release with catch up = the tip runs to the cursor. Average = moving average of the last N points (no string, slider = strength) |
| Color wheel | in the color popover (172 wide): ring 88 Ø, 7 thick, triangle inside (hue / white / black corners, rotates with the hue), both smooth color content in a 1 px `void` ring; markers: current hue 5×5 `bright` ring, harmony hues 5×5 filled with their color, triangle point 5×5 `bright` ring. Right of it: hex, alpha + eyedropper, old / new swatch |
| Harmony | `HARMONY` + name right, 6 icon segments (complementary, split, analogous, triadic, square, mono; 27 wide), then the swatches 24×14 (click = use, drag = palette) and `+ add all` (`lizard`) → `colors.sky_2`, `sky_3` |
| Trace dialog | 452×270 at y 34 over the 62% `void` dim; header 15 (`TRACE IMAGE · file`); preview 240×180 with before (left) / after (right) and a vertical 1 px `bright` split (+ `void` side) with an 11×17 grip; chips `original` / `traced`; source row (file · size + `Change…`), drop hint; right: MODE Photo / Line art, sliders line threshold, colors, detail, snap to shapes, flow; ☑ show lines, ☑ show fills; stats (`amber` numbers: strokes, fills, primitives, ms); footer Insert as layer + primary Open as drawing |

R004 – Layer hover preview ✅ (🧪 Reviews, section `25:4344`; generator `docs/Art/Figma-gen/r004.py` on top of `r002.py` / `r003.py`, split into layers by `raster.py`; PNGs in `docs/Art/Reviews/R004/`). Developer: "when we hover on a layer it should show a preview of only the items in the layer!" The sun layer has a third object for this review (`circle((360, 90), 80, #fff3c433)`, a faint halo) so it reads "3 objects".

| # | Frame | Node | PNG |
|---|---|---|---|
| 1 | Coder – hover card on `layer sun` (option **card**) | `25:4345` | `R004-01-coder-hover-card.png` |
| 2 | Coder – canvas isolation only, no card (option **isolation**) | `25:4396` | `R004-02-coder-isolation.png` |
| 3 | Artist – hover card beside the layers panel | `25:4348` | `R004-03-artist-hover-card.png` |
| 4 | Artist – card + isolation (option **both**) | `25:4731` | `R004-04-artist-card-and-isolation.png` |
| 5 | Coder – Spanish ("3 objetos", "clic: elegir la capa") | `25:4441` | `R004-05-coder-spanish.png` |
| 6 | Coder – nested `layer rays` inside `layer sun`: only the rays in the preview, path `sun/rays`, `8 objects` (made by a `for`), `SCR 70%` badge | `25:4486` | `R004-06-nested-layer.png` |
| 7 | Coder – hidden `layer boat` (`amber` eye-off on the row): the preview still shows it with a `hidden` chip, hint "only shown here"; the canvas does not | `25:4531` | `R004-07-hidden-layer.png` |

R004 metrics (UI pixels, proposed):

| Part | Metric |
|---|---|
| Hover card | tooltip construction (`row` face, 1 px `void` outline, `bevel` top line, 2 px `void` drop shadow), 152 wide, padding 6. Preview 140×93 (the drawing's 3:2 ratio, whole canvas so you see where the layer sits) in a 1 px `void` frame, the layer alone drawn full-res over the canvas background (`#f6d9a8` here); a drawing with no background shows the 3 px `#d4dacd`/`#a0ab9b` checker instead (color-content rule). Under it, pitch 15: `layer` icon (`lizard`) + name `bright` (parent path `sun/` in `mute` for a nested layer), shape count right-aligned `mute` (`3 objects` / `3 objetos`, the shapes the layer draws, loops included, not its rows); a blend/opacity badge line (`raised`, like the layer-row badge: `SCR 70%`) only when not Normal 100%; a `mute` hint line. Hidden layer: 13 px `raised` chip with `void` outline, `amber` eye-off + `hidden` in the preview's top-left. Height 128 (148 with the badge line) |
| Card placement | coder: right of the hovered row's text (+22), top 6 above the row, 5-step pixel tail (`void` edge, `row` fill) on the left pointing at the row; artist: left of the layers panel (7 px gap), tail on the right. Flips up near the bottom edge like menus. Proposed 2026-11-04 (design check 2): the card never covers the hovered layer's box on the canvas: it slides along the panel edge until it clears the box (the tail stays on the row), else it goes to the other side of the box; otherwise Alt isolation hides the one layer it shows |
| Card preview blend (2026-11-04) | the preview draws the layer's shapes in Normal at full opacity, nested layers too, so a `screen` or `add` layer stays readable on the light background; the badge line (`SCR 70%`) names the real mode, and Alt isolation shows the real blend |
| Timing | appears after the tooltip delay, 400 ms of rest on a layer row; while a card is open, moving to another layer row swaps it at once (no new delay); hides 150 ms after the cursor leaves the row, at once on click, drag, scroll or Esc. Never while dragging rows or drawing |
| Canvas isolation | over the drawing area only: `void` at 70% over the whole drawing, then the hovered layer drawn again above the dim (its real colors and blend), plus a chip under the canvas (`row` + `void` outline + `bevel` top line, 16 high: `layer` icon `lizard`, `sun · 3 objects`). The surround, panels and selection overlays are not dimmed. Hidden layers: no isolation (the card is the only place they show) |
| Hidden row | the row's eye becomes `eye-off` in `amber` |

R005 – Brush creator, brush files, video tooltips ✅ **approved 2026-11-04** ("yeah do the brush creator"), with layout A and these answers: (1) the creator replaces the R002 pen pop-up, ⚙ opens it on Pressure, stabilization moves to the Stroke tab; (2) Save brush writes the library `~/.linkz/brushes/<key>.lzbrush`, and a drawing gets a copy of the hash the first time it uses one; (3) Export set = all custom brushes; (4) the UI says "brush", the code keeps `pens.`; (5) the extra stamps and grains in the frames (spray, leaf, bitmap; paper, canvas, chalk, hatch, bitmap) are built with the creator; (6) tooltip videos are drawn live by the app, reduced motion follows the OS; (7) import result: the dialog for Procreate/ABR/SUT, quiet with a toast for `.lzbrush`; (8) a file dragged over the brush tool opens the picker after 500 ms (🧪 Reviews, section `29:4808` at y 2720; generator `docs/Art/Figma-gen/r005.py` on top of `r002.py` / `gen_r002b.py` / `r003.py`; PNGs in `docs/Art/Reviews/R005/`). From `docs/Linkz – Brushes.md`: a brush is a pen hash grown to Procreate-style groups, the creator writes it, `.lzbrush` files carry it. Each frame is three to fifteen layers: `editor behind (unchanged app, flattened)` plus one exact PNG per R005 part, named after it (`brush creator (floating window, layout A)`, `pen picker (open)`, `import result dialog`…). The smooth content (stamps, grain, test strokes, video frames) is rendered at screen resolution; every sub-part is named in `r005.py`.

| # | Frame | Node | PNG |
|---|---|---|---|
| 1 | Brush creator – coder, Shape tab, **layout A** (side tab rail, wide pad under the tab) | `29:4809` | `R005-01-creator-coder-shape-layoutA.png` |
| 2 | Brush creator – coder, Shape tab, **layout B** (top tab strip, tall pad on the right) | `29:4810` | `R005-02-creator-coder-shape-layoutB.png` |
| 3 | Brush creator – artist, Grain tab (opened from the ⚙ in the brush inspector) | `29:4811` | `R005-03-creator-artist-grain.png` |
| 4 | Brush creator – artist, Spanish, Pressure tab (the R002 curve editor, shared) | `29:4812` | `R005-04-creator-artist-es-pressure.png` |
| 5 | Pen picker – ⋯ menu open (New brush…, Edit in creator…, Duplicate · Export brush…, Export set…, Import…), hover Export brush + tooltip | `29:4813` | `R005-05-picker-menu.png` |
| 6 | Drop **option picker**: a `.brush` dragged onto the brush tool springs the picker open, drop zone row at the bottom | `29:4814` | `R005-06-drop-on-picker.png` |
| 7 | Drop **option window**: the file anywhere over the window, 62% dim + dashed drop card naming the file and the types | `29:4815` | `R005-07-drop-on-window.png` |
| 8 | Import result **option dialog**: preview stroke, shape + grain from the zip, SKIPPED 2 / CHANGED 1 / TAKEN 9, Cancel · Open in creator… · Add brush | `29:4816` | `R005-08-import-dialog.png` |
| 9 | Import result **option quiet**: the brush is added at once (`new` badge + `amber` warn on its row, hover = what was skipped) and a toast with Details / Undo | `29:4817` | `R005-09-import-toast.png` |
| 10 | Tool tooltip with a mini video – lasso (layout A, video on top), playing | `29:4818` | `R005-10-tool-tooltip-video.png` |
| 11 | Pen tooltip with a mini video – calligraphy row (layout B under the picker: no room beside it at 1024 wide) | `29:4819` | `R005-11-pen-tooltip-video.png` |
| 12 | Video tooltip states sheet (1800×620): layout A / B × loading, playing, reduced motion | `29:4820` | `R005-12-video-tooltip-states.png` |

Options to choose between (designer's pick first):

| Choice | Options | Pick |
|---|---|---|
| Creator layout | **A** side tab rail (8 × 18 px rows, icon + label + `amber` changed dot), tab content beside it, pad 488×70 under everything · **B** tab strip on top (labels only), content 290 wide, pad 188×148 on the right | **A**: every tab fits without scrolling (Shape shows all 5 settings; B fits 3 and drops scatter/count), the rail reads like the layer rows, and a wide pad is where long test strokes go. B keeps the pad in view while you scroll a tab, which matters only if tabs grow |
| Drop target | **picker** drop zone (spring-loaded open over the brush tool) · **window** card (anywhere) | **both**: the window card catches any drop (it names what will happen to that file type, so `.lzdraw` and images keep their old behaviour), the picker zone is the precise target once it is open |
| Import result | **dialog** before adding · **quiet** add + toast + warn on the row | **dialog** for Procreate / ABR / SUT (best effort: you should see what was skipped before it lands), **quiet** for `.lzbrush` (lossless: nothing to warn about) |
| Video tooltip | **A** video on top (180 wide) · **B** video beside (282 wide, 96×54 video) | **A** for the toolbar and the artist tool column, **B** for picker rows and anywhere narrow (frame 11) |

R005 metrics (UI pixels, proposed):

| Part | Metric |
|---|---|
| Brush creator window | the pen pop-up construction (R002 18b): floating, 500 wide, coder 324 high at y 18, artist 266 at y 52; 1 px `void` outline, 3 px hard `void` shadow, `panel` face, title bar 14 `BRUSH CREATOR · chalk` + `from ink` (`dim`, the pen it started from), `Reset`, ✕. Opened from the picker's ⋯ → New brush… / Edit in creator…, the ⚙ on a pen row, the ⚙ in the brush inspector, or double-clicking the brush tool. **It replaces the pen settings pop-up** (see questions): the ⚙ opens it on the Pressure tab |
| Tab rail (A) | 92 wide, `coal`, 1 px `void` line on its right, 8 rows of 18: icon at +6 (8×8, `mute`; `lizard` when selected), label at +18 (`text`; `bright` selected), selected = `moss` + 2 px `lizard` bar. A 3×3 `amber` dot 5 from the right edge marks a group changed from the base pen. Icons (new 8×8): `stroke`, `taper`, `shape`, `grain`, `dynamics`, `graph` (pressure), `fill` (color), `render` |
| Tab strip (B) | 15 high, `coal`; tab = label + 12; selected like a file tab (`panel` face, 1 px `lizard` top line, `void` sides); changed dot 2×2 `amber` top-right |
| Tab content | from rail + 8, 3 columns: presets (tiles), preview, settings. Settings rows pitch 18: label (`mute`, column = longest label + 6), slider, value right-aligned (`amber`, 32 wide). Segmented controls 14 high. A rule + one `mute` note line under the rows |
| Preset tiles | 24×24 at step 27, 4 per row: `coal` well, 1 px `void` outline; selected = 1 px `lizard` outline + 5×5 `lizard` corner. The stamp / grain inside is content (smooth, `bright` = paints). Last tile = import: `panel` face, dotted `bevel` edge, `plus`. `Import image…` button 16 high under the tiles, hint `PNG · white paints` |
| Stamp / grain preview | 64×64 in a 1 px `void` frame. Stamp: the shape large on `coal` + a 1 px `lizard` line at the angle. Grain: a short stroke of the brush on paper so you see the texture inside the stroke |
| Shape tab | stamp: round, square, flat nib, star, spray, leaf, imported bitmap, import; rotation Fixed / Follow / Random, angle, scatter, count, roundness |
| Grain tab | none, noise, paper, canvas, chalk, hatch, imported, import; texture Moving / Fixed, scale, depth, ☐ scale with brush size; note "fixed: stays on the paper" |
| Pressure tab | exactly the R002 pop-up's curve editor (84×84, presets, `sky` live line + readout) and its input switch + overrides box; sliders size (range), opacity (range), smoothing |
| Scribble pad | the R002 test pad grown: paper `#fbf8ef`, 1 px `void` frame. A **sample stroke** (a wave with a pressure ramp, left) is redrawn on every change; your own test strokes stay until cleared; `sample` label `mute` bottom-left, `test only` `dim` bottom-right. Three 13×13 buttons top-right: sample on/off (`stroke`, active), paper light / dark (`fill`, for light brushes), clear (`trash`). Strokes never reach the drawing. A: 488×70; B: 188×148 |
| Footer | `name` + field 120×18, then right-aligned `Export…` (button, `export` icon, writes `.lzbrush`) and primary `Save brush` (`moss` + `lizard` bevel, `save` icon), 4 apart |
| Coder code block | as the pen pop-up: `coal` inset, 4 lines of 13, the pen hash it writes, scrolled to what changed (`row` line + 2 px `amber` bar), label `pens.chalk · 1-4 / 8` `dim` top-right, 3 px scrollbar |
| Picker ⋯ | 8×8 `more` icon 16 from the header's right edge, `text`; open = 14×12 active button. Menu (context-menu construction, rows 14) under it, title = the current pen: New brush…, Edit in creator…, Duplicate, separator, Export brush…, Export set…, Import…. Right-click on a row opens the same menu for that row |
| Drop zone (picker) | a 24-high row above the footer: `moss` fill, 1 px dashed `lizard` (2 on 2 off), `import` icon + `Drop to add` `bright`, type right-aligned `mute`. The picker springs open after 500 ms of dragging a brush file over the brush tool |
| Drop card (window) | 62% `void` dim over the window, card 300×138 at y 100 (tooltip construction), 2 px dashed `lizard` border 6 in (4 on 2 off), `import` icon, `Drop to import a brush` `bright`, file + source `text`, rule, accepted types `mute`, "images still open the trace dialog" `dim` |
| Drag ghost | file chip 16 high (`row`, `void` outline, 2 px shadow, `file` icon + name) 10 right / 12 below the cursor, 9×9 `lizard` copy badge with a `void` plus at its top-left |
| Import dialog | 380×276 at y 34 over the 62% dim; header 15 `IMPORT BRUSH · file`. Left 150: preview 150×56 (sample stroke with the imported brush), shape and grain tiles with file name + size, `stamped` chip. Right: name `bright`, source `mute` + "best effort", rule, then **warnings first**: `SKIPPED n` (`warn`, `amber`) with one `text` row each, `CHANGED n` (`info`, `sky`) with what it became in `mute`, `TAKEN n` (`check`, `lizard`) as a wrapped `mute` list. Footer: where it will be saved (`dim`), Cancel · Open in creator… · primary Add brush |
| Import quiet | the new row gets a `new` badge (`raised` + `lizard`, 3 s like the lasso fill row) and an 8×8 `amber` `warn` 12 left of the ⚙ while notes exist; hovering it shows the notes tooltip. Toast: tooltip construction, 236 wide, 2 px `lizard` left bar, bottom-right 6 above the status bar, `check` + "Added …" `bright`, `warn` + "3 notes" + `Details` (`sky`) + `Undo` (`mute`, right); 6 s, stays while hovered |
| Video tooltip A | tooltip construction, 180 wide: video 168×94 (336×188 on screen, 16:9) at +6 in a 1 px `void` frame, progress 1 px under it (`bevel` track, `lizard` played), then pitch 15: title `bright` + shortcut badge right (`raised`, `mute`), 2 description lines `text` (24 characters a line) |
| Video tooltip B | 282 wide: video 96×54 left, text column 168 to its right (title + badge, 2 lines), progress under the video |
| Video states | **loading**: the poster frame under 62% `void` + three 3×3 dots (the lit `lizard` one walks, the others `mute`); progress track empty. **playing**: loops, progress fills. **reduced motion** (OS setting or a Linkz preference): the last frame still, no progress, chip `motion off` (13 high, `raised`, `void` outline, `still` icon) top-left. No sound ever |
| Video timing | appears with the normal tooltip delay (400 ms); plays from the start each time the tooltip opens; loops 3-4 s clips up to 3 times, then rests on the last frame; stops when the tooltip closes |
| New 8×8 icons | `more`, `taper`, `shape`, `grain`, `dynamics`, `render`, `import`, `file`, `play`, `still`, `duplicate`, `image` (in `r005.py`) |

R005 API (the hash the creator writes). It follows **B060 – Stamped brushes** (`lib/draw/stamps.liz`, being built in parallel): a brush is a pen with **flat** keys, so it works with `brush(points, pressures, color, pens.chalk)` unchanged. The creator's tabs group the keys in the UI only; the code stays flat, which is easier to write by hand and matches the `pens.liz` header. The creator writes only the keys that differ from the pen it started from (those are the tabs with the `amber` dot).

```
pens.chalk = {"name": "chalk", "render": "stamped", "min": 2, "max": 14,
  "spacing": 0.1, "jitter": 0.05, "shape": "flat", "rotation": "follow", "angle": 30,
  "scatter": 0, "count": 1, "roundness": 0.35,
  "grain": {"kind": "chalk", "scale": 1.2, "depth": 0.7, "moving": false},
  "size_jitter": 0.1, "opacity_jitter": 0.15, "speed": 0.2,
  "pressure": {"curve": [[0, 0], [0.3, 0.6], [1, 1]], "opacity": [0.3, 0.9]}}
```

What the frames show beyond B060 (proposed, needs a language/library decision): stamp shapes `"spray"`, `"leaf"` and an imported bitmap (`"shape": {"image": <run-length data>}`), grain kinds `"paper"`, `"canvas"`, `"chalk"`, `"hatch"` and an imported bitmap, `streamline` / `falloff` in the Stroke tab, and per-stroke vs per-stamp color jitter. `grain` as a number keeps today's meaning (edge roughness), as B060 already does.

## 🖼️ Inputs

| File | For | What to take from it |
|---|---|---|
| [Aseprite docs: workspace](https://www.aseprite.org/docs/workspace/), [tool bar](https://www.aseprite.org/docs/tool-bar/), [context bar](https://www.aseprite.org/docs/context-bar/), [timeline](https://www.aseprite.org/docs/timeline/), [status bar](https://www.aseprite.org/docs/status-bar/) | R001 layout and feel | Docked compact panels, beveled 1 px pixel buttons, a tool strip, file tabs, a status bar that changes with the mouse, dense rows. Not its art, icons or colors. |
| — | | Drop sketches or references in `docs/Images/` |
| `docs/Art/References/lilies-stargazer.png` | Splash bloom lilies | Oriental "Stargazer" lily: six wide pointed petals curving back at the tips; white rim grading to strong pink toward the middle; a deeper magenta stripe down each petal's center; scattered crimson speckles on the inner petals; yellow-green star throat; green stamens with rust/brown anthers; long closed purple-pink buds; narrow lance-shaped leaves |

## 📓 Art log

| Date | What | Verdict |
|---|---|---|
| 2026-11-02 | 8-bit UI direction (developer: "an 8 bit style UI since it's easier for you to draw it that way") | ✅ |
| 2026-11-02 | Splash logo "Linkz" (B012): 104×31 custom chunky letters, 6 px strokes, `leaf` fill with a 1 px `lime` top highlight, `moss` shadow 1 px down-right, `ink` 8-neighbour outline. Source `docs/Art/splash.aseprite`, previews `docs/Art/Previews/logo-3x.png`, `banner-3x.png`, `banner-dark-3x.png` | ✅ crisp at 1×, reads on transparent and near-black |
| 2026-11-02 | Splash lizard walk (B012): 48×16 side view, 6 frames, diagonal-pair gait (4 stance + 2 swing frames, feet move 2 px/frame so it advances `lizard_step` = 2 px/frame without sliding), tail wave, `sun` eye with `ink` pupil, far legs `moss` with their own outline behind the body. First pass rejected (legs merged into the body, duck-bill head); second pass fixed both. Previews `lizard-walk-4x.png`, `crawl-3x.gif` | ✅ reads as a lizard walking |
| 2026-11-02 | Splash layout: lizard's back hides behind the logo's lower edge, belly and legs show below it (walking behind a sign); tagline UI font `#a39e88` with a 1 px `#1b1a17` drop shadow so it holds up on light desktops too | ✅ |
| 2026-11-02 | `splash_art.liz` stores hex COLOR literals, not palette names (palette moving to "Linkz Dark"), so a palette swap can't break it; `tools/spritegen` accepts any opaque color, rejects half-transparent pixels | ✅ |
| 2026-11-02 | Logo "ink" (developer: "highlight the ink in Linkz"): i, n, k filled wet ink `#3b40a8`, `#a8b0ff` top sheen, `#6870e0` left edge, `#eef0ff` glints, `#1c1f5c` shadow; L and z stay green; shared `ink` outline keeps it one word. A drip inside the n's arch was rejected (read as "m"); the drop now runs off the bottom of the k stem. Logo is 104×36 with the drop. Previews `logo-3x.png`, `logo-dark-3x.png` | ✅ "ink" pops on dark and on the checker |
| 2026-11-02 | Splash vines: 24 growth stages, from a 4 px sprout at the base of the L to vines wrapping the L, i, n, k, z and two side runners. Each stage only adds pixels (Catmull-Rom paths drawn up to a length that grows each stage); growing tips are pale `#c8f08a`, stems are `#2f9470` with a `#0f3a2f` edge 1 px down-right, leaves are `#3fc28e`/`#b4f2d2` and go from bud to half to full. The bluish greens stay visible on the green and on the ink letters. One closed purple-pink lily bud on the n tendril from stage 21. Previews `vines-sheet-dark-2x.png`, `banner-vines-{0,50,100}-dark-3x.png`, `banner-vines-100-3x.png` | ✅ organic, letters still read at 100% |
| 2026-11-02 | Splash bloom: the first pass used 5 small mixed flowers (coral, yellow, lilac, white) and was replaced at the developer's request ("lirios, big white, pinkish-white lilies", Stargazer reference). Now two Stargazer lilies (R 14, ~28 px) on the side runners, staggered 2 frames, over 10 frames: long bud → swollen pink bud → split → open → stamens → sparkle. White rim, pink body, 1 px magenta line per petal, crimson speckles, yellow-green throat, green stamens with rust anthers, `#2a1424` outline, back petals pale lilac, overlaps shaded `#b77aa6`. Previews `bloom-sheet-dark-2x.png`, `banner-final-dark-3x.png`, `banner-final-3x.png` | ✅ |
| 2026-11-02 | Linkz Dark palette (16 colors, contrast table above) replaces Lizard 16, per "similar in UI to Aseprite, but black and dark" | ✅ approved 2026-11-02 |
| 2026-11-02 | R001 – Editor views, Figma "Linkz" → 🧪 Reviews (section `3:2`): main, new row with inline tool strip, pencil (anchors, one hovered), brush + pen picker, tooltip + error row, Spanish, splash, plus home gallery, home empty and home Spanish. The six decisions are approved too (no name column, 4 px blank lines, English parameter labels in Spanish, full-res pen previews, last good drawing on error, splash keeps its art colors) | ✅ approved 2026-11-02 ("design is approved! I'm very happy with it"); copies on ✅ Approved |
| 2026-11-02 | Row spacing (developer: "maybe the lines are too close together"): row pitch 14 → 18, 4 px gap above layer rows, properties rows to 18 high. Frame 1 now scrolls past `layer sea` (scrollbar shown) | ✅ approved 2026-11-02 |
| 2026-11-02 | Home screen (developer request): gallery with "Working on" (+ New drawing first, 3 recent big cards) and "All drawings" grid, rendered thumbnails, hover card with ⋯ menu (Open, Rename, Duplicate, Show in Finder, Delete); empty first-launch state with the lizard, "+ New drawing" and "Open an example"; Spanish version (labels fit; "Nuevo dibujo" wraps) | ✅ approved 2026-11-02 |
| 2026-11-03 | R002 – shape menu (canvas or row right-click / two-finger click): Rename, Duplicate, Bring forward / Send backward (= move the line; disabled when already first in its layer), Hide, Lock, Fill ▸ (Solid / Linear / Radial / Edit fill… + variable swatches), Copy code, Show in code, Remove (`ember`). Frames `6:1258`, ES `6:1281` | ✅ approved 2026-11-03 |
| 2026-11-03 | R002 – fill editor: floating beside the shape, kind tabs, gradient bar with draggable stops (click to add, drag off to remove), stop color from the palette sidebar, soft/hard per stop, steps (off or 2–16 flat bands), live code preview, Save as variable; canvas handles (linear start/end, radial center/radius). Linear on the sky `6:1259`, 5-band stepped `6:1260`, radial with a hard stop (ES) `6:1282`. Figma's SVG import drops duplicate-offset (hard) stops, so the hard-stop sun is drawn as a disk + ring | ✅ approved 2026-11-03 |
| 2026-11-03 | R002 – palette sidebar (106 wide, collapsible to a 14 px rail): HSV picker, recent, VARIABLES grid (solid and gradient chips with names, kind tabs Solid/Linear/Radial when one is selected, "in N shapes"), + Variable. Editing `colors.sky` repaints the sky and the waves `6:1261`. Proposes the "color content" exception to the no-gradient rule | ✅ approved 2026-11-03 |
| 2026-11-03 | R002 – artist mode: vertical tools left, canvas center, right column = compact palette + layers with thumbnails, blend dropdown + opacity in the layers header, properties as an Aseprite-style context bar under the tabs, `</> Code` toggles a code drawer under the canvas (same rows). `6:1262`, drawer `6:1263` | ✅ approved 2026-11-03 |
| 2026-11-03 | R002 – blend modes + opacity: `blend("multiply")` / `opacity(0.8)` rows inside the layer, `MUL 80%` badge on the layer row, properties blend dropdown + opacity slider; artist header dropdown. `shadows` layer really in Multiply 80% in Figma (visibly darker over sky, sea and boat). Coder `6:1264`, artist `6:1265` | ✅ approved 2026-11-03 |
| 2026-11-03 | R002 – first-run tutorial: 62% `void` spotlight with `amber` rings, coach-mark with step dots, Skip, Next; (1) tool strip, (2) variables + fill field, (3) brush + canvas with a dotted ghost stroke, (4) the new code rows, lizard celebrating with sparkles, "Let's draw!". `6:1267`–`6:1270`, Spanish step 1 `6:1271` | ✅ approved 2026-11-03 |
| 2026-11-03 | R002 – object budget: New drawing dialog (name, size presets, custom W×H, "Recommended: up to ~1,800 objects · measured on this Mac" + ⓘ tooltip) EN `7:1963` / ES `7:2025` (cloned from the approved home galleries); status bar counter normal / near (≥ 80%, `amber`) / over (`ember`, tooltip) `6:1272`; the counter replaces "steps · ms" | ✅ approved 2026-11-03 |
| 2026-11-03 | R002 – pen settings: pressure curve editor (soft/linear/firm presets), width and opacity ranges, taper, smoothing, speed, Trackpad/Tablet/Mouse overrides (trackpad floor), live preview, Save as new pen; coder shows the `pens.soft_ink = {…, "pressure": {…}}` it writes `6:1273`; artist (ES) from the brush tool `6:1274` | ✅ approved 2026-11-03 |
| 2026-11-03 | R002 – autocomplete polish frames (approved R001 layout, 320 code panel): `ci` popup with kinds and matched letters `6:1275`, signature hint with the current parameter `6:1276`, `math.s` members `6:1277` | ✅ approved 2026-11-03 |
| 2026-11-03 | R002 – color popover beside the clicked color (code-row swatch `6:1278`, properties fill field with inline Save as variable `6:1279`, artist inspector fill, Spanish `6:1280`): Variables (solid + gradient chips, picking writes `colors.<name>`), Recent 14, compact HSV + hex + alpha + eyedropper, keyboard focus `amber` | ✅ approved 2026-11-03 |
| 2026-11-03 | R002 rev – pen settings as a **pop-up** (developer: "make the pressure adjustment be a pop up and have a tiny space to test the brush"): floating window opened from ⚙ on each pen row / the brush inspector / double-click brush, canvas and code behind it; everything the panel had, plus a 300×66 test pad (3 example strokes, trash, faint hint when empty, strokes never saved) and a live readout `pressure 0.62 · trackpad` matching the `sky` line; + STABILIZATION block (Off / String / Average, length, catch up). Coder `9:2304`, artist ES `9:2305`, empty pad `9:2306`, picker with ⚙ `9:2303`. Proposes a 3 px hard drop shadow for floating windows (was 2 px for tooltips/popovers) | ✅ approved 2026-11-03 |
| 2026-11-03 | R003 – select / move tool (V, first tool): click, Shift-click, marquee, drag to move, handles (+ optional rotation handle), ghost + delta while moving, rows highlighted `9:2307` | ✅ approved 2026-11-03 |
| 2026-11-03 | R003 – lasso (L): select mode (marching ants, objects inside get selected, then move like Select) `9:2308`; fill mode (Clip Studio style: the loop becomes a `polygon` row with the current color) `9:2309`; stabilizer control in its options strip | ✅ approved 2026-11-03 |
| 2026-11-03 | R003 – stabilization, live while drawing: pulled string (`sky` string from cursor to tip, radius circle, ink only at the tip, catch up at the end) for brush, pencil and lasso; per pen (in the pen hash) and in the brush inspector `9:2310` | ✅ approved 2026-11-03 |
| 2026-11-03 | R003 – color wheel (hue ring + triangle, Krita/Painter style) + harmonies (complementary, split, analogous, triadic, square, mono) with ring markers, swatches and "+ add all as variables" in the color popover `9:2311` | ✅ approved 2026-11-03 |
| 2026-11-03 | R003 – trace image dialog (Tracing v2 UI): Photo / Line art, line threshold, colors, detail, snap to shapes, flow, before/after split on an original placeholder head (no copyrighted art), show lines / fills, stats, Insert as layer / Open as drawing `9:2312` | ✅ approved 2026-11-03 (design only, not scheduled) |
| 2026-11-03 | R002 + R003 moved to ✅ Approved (developer: "everything on figma is approved, move it and check it"): sections `21:3120` (R002, 29 frames) and `21:4094` (R003, 6 frames), plus R001 `21:2692`. They were moved, not copied, before this pass, so 🧪 Reviews is empty. The 35 frames are exported at 1× to `docs/Art/Approved/21-*.png` and listed in the approved tables | ✅ moved · ⚠️ review copies gone from 🧪 Reviews (pictures kept in `docs/Art/Reviews/`) |
| 2026-11-03 | v1.1 design check, real app vs approved frames (`docs/Reviews/2026-11-03 v1.1 design check.md`, shots in `docs/Art/Checks/v1.1/`) | 3 ✅ · 23 🟡 · 5 🔴 (R002-6 artist selection not revealed, R002-9 submenu off-window, R002-20 / 22 autocomplete layout + error while typing, R002-25 named colors open no popover) · R003-3 not verifiable · R003-6 not built · R002-18/19 superseded |
| 2026-11-03 | R004 – layer hover preview (developer: "when we hover on a layer it should show a preview of only the items in the layer!"): hover card with the layer alone over the canvas background at 140×93, name, shape count, blend badge, hint; canvas isolation option (70% `void`, layer drawn again on top, chip under the canvas); coder `25:4345`, isolation `25:4396`, artist `25:4348`, both `25:4731`, Spanish `25:4441`, nested `sun/rays` `25:4486`, hidden `boat` `25:4531` (section `25:4344`). Designer's pick: card on hover, isolation while Alt is held | ✅ approved 2026-11-03: card on hover, isolation while Alt is held |
| 2026-11-03 | Proposal from the check: R002-9's code preview `hard(#f2a65a, 0.6)` should read `hard(#f2a65a)` (the library's one-argument `hard` is the better API). Additions the app made (Arrange + Blend ▸ in the shape menu, fill-editor `angle` / `alcance`, dialog Background + `measure again`, tutorial `Back`) should be drawn into the frames next pass | ⏳ proposed |
| 2026-11-04 | Design check 2 (`docs/Reviews/2026-11-04 design check 2.md`, shots in `docs/Art/Checks/v1.2/`): the first check's 🔴/🟡 frames, the gesture frames now that drive has Press/Move/Release/Path, the R004 frames and the R003-6 trace dialog | 20 ✅ · 17 🟡 · 1 🔴. Every B053 fix holds. 🔴 R004-4: the artist card covers the hovered layer, so Alt isolation shows nothing. Top 🟡: live brush seams, hidden rows (`dim` eye-off, `-~> [hidden]` shown), the eye click scrolls the code to the top, palette stuck on `colors.sky`, no "nothing selected" state |
| 2026-11-04 | R004 spec notes from check 2: the card clears the hovered layer's canvas box; the card preview draws in Normal (the badge names the blend) | ⏳ proposed |
| 2026-11-04 | R005 – brush creator (Brushes doc): floating window like the pen pop-up, 8 group tabs, shape / grain preset tiles + import image, stamp / grain preview, scribble pad with a live sample stroke, Save brush / Export…, coder code block. Layout A (side rail) `29:4809`, layout B (top tabs) `29:4810`, artist Grain `29:4811`, artist ES Pressure with the shared curve editor `29:4812` (section `29:4808`) | ⏳ waiting for approval · designer's pick: layout A |
| 2026-11-04 | R005 – brush files: picker ⋯ menu (New, Edit in creator, Duplicate, Export brush, Export set, Import) `29:4813`; drop on the picker `29:4814` or the window `29:4815`; import result as a dialog with skipped / changed / taken `29:4816` or quietly with a row warning + toast `29:4817` | ⏳ waiting for approval · pick: both drop targets; dialog for foreign formats, quiet for `.lzbrush` |
| 2026-11-04 | R005 – mini videos in tooltips: lasso tool (A, video on top) `29:4818`, calligraphy pen row (B, video beside, under the picker) `29:4819`, states sheet loading / playing / reduced motion × A / B `29:4820` | ⏳ waiting for approval · pick: A for tools, B for picker rows |
| 2026-11-04 | R005 process: first render check found and fixed faint chalk strokes, a caption over a hint, a label touching its changed dot, layout B sliders running into the footer, a grain note off the window, a pen tooltip off the window (moved under the picker as layout B), an import list hitting the save line, and the ⋯ glyph missing from the font (hints now say "menu"). Figma exports checked against the generator renders. The coder code block first used nested groups; changed to B060's flat keys (`"render": "stamped"`, `"shape": "flat"`, `"grain": {"kind"…}`) so the creator writes what `stamps.liz` reads | ✅ fixed before review |
| 2026-11-04 | R006 – app icon. **Round 1** (A nib with an eye, B lizard whose tail is a Bézier path, C lizard on an ink drop): developer: "I imagined it more like a lizard face with a paintbrush on the side" ❌ all three rejected as a direction (B had the best silhouette). **Round 2** (same Figma section, second row; generator `docs/Art/Figma-gen/r006b.py`; PNGs and `overview-2.png` in `docs/Art/Reviews/R006/`): a big lizard face with a paintbrush (peach handle, `mute` ferrule, ink-blue bristles with `#a8b0ff` sheen). 52×52 cell art, 832 px squircle, same palette, `void` outline, no soft shadows; the 16 px file is a hand-simplified 14×14 sprite, not a downscale. **D** friendly: round front face, big eyes, small crest, blush, smile, brush tucked behind the right ear ✅ most readable at 32, but a little frog-like (the crest helps). **E** cool: sharp three-quarter face facing right, spiky crest, slanted brow, fangs, brush held in the teeth ✅ most lizard, best attitude; most detail, so the 16 px sprite is the weakest. **F** bold: front face with a spiked crest and toothy grin, a clawed hand holds a tall brush with a fat ink drop ✅ brush story is clearest, but the brush squeezes the face to the right and nearly vanishes at 16.  **Round 3**: developer picked **E**, with one change ("on the second one just remove the left eye, no changes elsewhere"): the small far eye is gone and nothing else moved, so the pixels it covered keep the face's own green and shading; the brow above it stays (it still sits over the big eye). The hand-drawn 16 px sprite never had that eye, so E2-16 equals E-16. Files `E2-cool-brush-in-teeth-one-eye-*` next to `E-*`, `overview-3.png` (E left, E2 right), generator `r006c.py`; Figma: E2 added to the right of F in the round 2 row | ✅ 2026-11-04 developer: "Icon E2 is the best". E2 is the app icon |

## Approval of R002 and R003 (2026-11-03)

> "The requisites are approved they look incredible" — the developer, 2026-11-03

Every R002 and R003 frame is approved as drawn. The designer's open questions are settled with the designer's recommendations: the pen settings pop-up keeps its full size (no tabs); the coder code block scrolls (4 of 6 lines); floating windows use a 3 px hard shadow; `dim` placeholder text on paper is allowed (empty states only); the harmony icons get hover tooltips; the color wheel replaces the square picker in the color popover, and the palette sidebar and fill editor get a ring/square toggle (ring by default); the selection box keeps the rotation handle; the "color content" exception holds (UI that shows the user's colors may show gradients); saved pens live in `~/.linkz/pens/`.
