# B019 – First-run tutorial (v1.1, wave 3, after B018)

Design: R002-10..13 tutorial steps 1–4, R002-14 step 1 (ES); Style Guide "Spotlight" row and design log row "first-run tutorial". Rules: `docs/Batches/_v11-rules.md`.

## Build
- Shown the first time the editor opens (settings `let tutorial_done = true` after finishing or skipping; "Show tutorial again" in Home/help).
- 62% `void` spotlight with `amber` rings around the target, a coach-mark with step dots, Skip, Next, Back; the steps from the frames: (1) the tool strip → make a first shape, (2) variables + the fill field, (3) the brush + canvas with a dotted ghost stroke to trace → a first brushstroke, (4) the new code rows (the lizard celebrates: the code is the drawing). Each step waits for the real action (or Next).
- Works in coder and artist mode (targets come from the live layout rects), both languages, any window size (coach-mark repositions).
- Tests: the rig walks every step in both modes and both languages; skip and resume.
