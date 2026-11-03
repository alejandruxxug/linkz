# 🎨 Linkz: Style Guide

## 🧭 Direction in one line

**An 8-bit editor around smooth vector art.** The UI is chunky pixel art (half resolution, scaled ×2, hard pixels); the drawing on the canvas stays full-resolution and anti-aliased.

## 🕹️ 8-bit UI rules

| Rule | Detail |
|---|---|
| Resolution | UI drawn at half the window resolution, scaled ×2 nearest-neighbor by the host. No anti-aliasing anywhere in the UI. |
| Palette | Only the 16 colors below. No gradients; shade with the darker/lighter neighbour or a 2×2 checker dither. |
| Font | The 7×13 bitmap font (`lib/gfx/font_data.liz`), scale 1 in the UI buffer (it's already ×2 on screen). |
| Outlines | 1 px `ink` outline on panels, buttons, fields, tooltips. |
| Buttons | 1 px bevel: `light` top/left, `shadow` bottom/right; pressed = bevel swapped and content moved 1 px down-right. |
| Selection | Selected row: `leaf` background, `ink` text. Focused field: `sun` outline. |
| Icons | 8×8 pixel icons as lizard data (one color + transparent, tinted when drawn): ○ □ ╱ ⬠ ～ ✎ 🖌 T ƒ ▤ layer, eye, lock, chevrons. |
| Caret | Solid 1 px `sun` bar (nothing animates while idle, so no blink). |
| Tooltips | `paper` box, `ink` outline, 1 px `shadow` drop shadow below-right. |
| Canvas | Full-res drawing inside a 1 px `ink` frame on a `panel-dark` checker surround. |

## 🎨 Palette: Lizard 16

| Name | Hex | Role |
|---|---|---|
| ink | `#1b1a17` | Outlines, text on light |
| night | `#2b2a24` | Window background |
| panel-dark | `#3a3830` | Panel background |
| panel | `#46443a` | Rows |
| shadow | `#6e6a5a` | Bevel shadow, muted text |
| stone | `#a39e88` | Secondary text |
| light | `#d6cfb2` | Bevel highlight |
| paper | `#f2ecd2` | Text on dark, tooltips, default canvas |
| moss | `#2f6b3a` | Dark green |
| leaf | `#3fa34d` | Lizard green: selection, default fill |
| lime | `#8fd16a` | Highlights, success |
| sun | `#e0a526` | Focus, caret, accent |
| butter | `#f2d06b` | Numbers in code |
| rust | `#c8642c` | Strings in code |
| ember | `#d94a3a` | Errors |
| sea | `#2d6a8f` | Keywords in code, links |

Code colors on the dark rows: `lime` keywords, `butter` numbers, `rust` strings, `stone` comments, colors drawn as a 6×6 swatch before the literal, errors `ember`.

## 🖼️ Inputs

| File | For | What to take from it |
|---|---|---|
| — | | Drop sketches or references in `docs/Images/` |

## 📓 Art log

| Date | What | Verdict |
|---|---|---|
| 2026-11-02 | 8-bit UI direction (developer: "an 8 bit style UI since it's easier for you to draw it that way") | ✅ |
| 2026-11-02 | Splash logo "Linkz" (B012): 104×31 custom chunky letters, 6 px strokes, `leaf` fill with a 1 px `lime` top highlight, `moss` shadow 1 px down-right, `ink` 8-neighbour outline. Source `docs/Art/splash.aseprite`, previews `docs/Art/Previews/logo-3x.png`, `banner-3x.png`, `banner-dark-3x.png` | ✅ crisp at 1×, reads on transparent and near-black |
| 2026-11-02 | Splash lizard walk (B012): 48×16 side view, 6 frames, diagonal-pair gait (4 stance + 2 swing frames, feet move 2 px/frame so it advances `lizard_step` = 2 px/frame without sliding), tail wave, `sun` eye with `ink` pupil, far legs `moss` with their own outline behind the body. First pass rejected (legs merged into the body, duck-bill head); second pass fixed both. Previews `lizard-walk-4x.png`, `crawl-3x.gif` | ✅ reads as a lizard walking |
| 2026-11-02 | Splash layout: lizard's back hides behind the logo's lower edge, belly and legs show below it (walking behind a sign); tagline UI font `#a39e88` with a 1 px `#1b1a17` drop shadow so it holds up on light desktops too | ✅ |
| 2026-11-02 | `splash_art.liz` stores hex COLOR literals, not palette names (palette moving to "Linkz Dark"), so a palette swap can't break it; `tools/spritegen` accepts any opaque color, rejects half-transparent pixels | ✅ |
