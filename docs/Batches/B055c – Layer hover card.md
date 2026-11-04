# B055c – Layer hover preview card (R004)

Model: Sonnet. Rules: `_v11-rules.md`. Read only: `docs/Art/Reviews/R004/*.png`, the Style Guide's R004 tables, `lib/ui/view_popups.liz` (tooltip timing), `layer_thumbs.liz`, `artist_layers.liz`.

Approved: the card on hover (isolation is a separate task). After 400 ms on a layer row (code rows and the layers panel, both modes): the card beside the row with only that layer's shapes on the canvas background (checker if none), title (`sun/rays` for nested), object count (shapes it draws), blend badge when not Normal/100%, "click: select layer"; hidden layers show a `hidden` chip. Switching rows while open is instant; hides 150 ms after leaving, at once on click/drag/scroll/Esc; never while dragging or drawing. Cache previews like thumbnails. en + es. Rig tests + shots vs R004 frames in `docs/Art/Checks/v1.1/`.
