# B055a – The layers panel in coder mode

Model: Sonnet. Rules: `_v11-rules.md`. Read only: `lib/ui/artist_layers.liz`, `layer_thumbs.liz`, `view_layout.liz`, `palette_view.liz`, `mode_state.liz`.

> "I like the layer system and I think it should be for everyone, not only artist mode, so keep it." — the developer

Show the artist layers panel in coder mode too: right column, under the palette sidebar (they share its width; the rail state collapses both). Same panel code (no copy): eye, lock, fold, blend badge, thumbnails, drag to reorder, right-click menu, blend + opacity header. A header chevron hides it; remembered as `let layers_panel = true` in settings. Rig test in coder mode; M10 benchmarks within 5% in coder mode.
