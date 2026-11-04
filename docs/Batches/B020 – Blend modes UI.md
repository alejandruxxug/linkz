# B020 – Layer blend modes and opacity in the UI (v1.1, wave 2)

Design: R002 blend frames (`row/layer with blend badge`, `blend dropdown (open)`), Style Guide design log row "blend modes + opacity". Rules: `docs/Batches/_v11-rules.md`. The engine already supports `blend("multiply")` / `opacity(0.8)` as a layer's first lines (lib/draw, composited per tile).

## Build
- Layer rows show a **blend badge** when the layer isn't Normal/100%; a dropdown on the layer row (and in the properties panel when a layer is selected) with Normal, Multiply, Screen, Overlay, Darken, Lighten, Add, plus an opacity slider/field; choosing writes/updates/removes the layer's `blend(...)` / `opacity(...)` line (one undo step; removing returns to Normal/100%).
- Wire B016's disabled **Blend ▸** entry in the shape menu: for a shape, its layer's blend (or, if the shape is alone in no layer, wrap it in a layer — ask nothing, just do the simplest documented thing and report it).
- Global shortcuts decided: Cmd+D duplicate, Cmd+] / Cmd+[ forward/backward work anywhere in the editor (they don't type text); single letters (H, L, F) stay menu-only. Implement the global Cmd ones (they belong to B016's actions in `editor_shape_ops.liz`; call them, don't duplicate).
