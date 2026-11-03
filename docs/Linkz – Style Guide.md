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
| **Linkz** (Figma desktop, edited through the Figwright plugin) | 🖼️ Inputs (direction + references) · 🎨 Style (palette, font 1× and 2×, components, 8×8 icons) · 🧪 Reviews (one section per review; R001 – Editor views is kept there as history) · ✅ Approved (the frames to build from) |

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

Style page: palette section `4:501` (swatches bound to the **Linkz Dark** variables `color/<name>`, paint styles `Linkz Dark/<name>`), font `4:502`, components `4:503` (`button/*`, `tool/*`, `field/*`, `row/*`, `tooltip`), icons `4:504` (components `icon/<name>`).

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
