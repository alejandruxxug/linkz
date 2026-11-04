# B021 – New drawing dialog + object counter (v1.1, wave 2)

Design: R002 New drawing dialog + `status counter/{normal,near,over}`; Style Guide design log row "object budget". Rules: `docs/Batches/_v11-rules.md`. The calibration already exists (`linkz calibrate`, `lib/app/machine.liz`, `lib/gfx/calibrate.liz`, `os.memory()`).

## Build
- **New drawing dialog** (from Home "+ New" and Cmd+N): name, size presets (square 512, 1080×1080, A4 at 150 dpi, 1920×1080, phone 1170×2532, custom W×H), background color, "Recommended: up to N objects on this computer" (from the machine calibration, measured once per machine/app version, refreshed with a "measure again" link) → writes `canvas = {"size": …, "background": …, "max_objects": N}`.
- **Status bar counter** `objects n / max` (normal; `amber` ≥ 80%; `ember` over), with a tooltip explaining the number; opening a drawing made on a faster machine warns once without editing the file.
- Revisit the calibration formula: tiles/patch runs mean edits no longer repaint every object (B035 noted it) — base the recommendation on full-repaint and open time at the 2k/10k fixtures' per-object costs, documented.
