# R002 revision (pen settings as a pop-up + test pad, pen picker with a gear per row)
# and R003 (select, lasso, color wheel + harmonies). Builds on r002.py / gen.py.
# Usage: python3 gen_r002b.py <outdir> [keys...]
import os, sys, math, colorsys
GEN = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, GEN)
import gen, r002
from gen import (R, G, T, I, C, tw, button, field, frame_lines, svg_doc, W, H, S, ROW_H, MAIN_Y, STATUS_Y, TOP_H, TAB_H,
                 sel_box, cursor, catmull, variable_outline, nm)
from r002 import (Lay, chrome, sfield, slider, seg, pixel_poly, pixel_circle, vec, code_panel, props2, canvas_area2, drawing_g,
                  layer_bodies, rows_layers, rail, status2, tabbar2, window2, curve_eval, CURVE, sw, hexrgb, rgbhex, mix, row_heights,
                  toolbar_artist, tool_column, inspector, right_column, art_layers, art_lay, var_grid, VARS, RECENT12, TX, PROPS_Y,
                  hsv_of, tail_px)

OUT = sys.argv[1]
os.makedirs(OUT, exist_ok=True)
PAPER = "#fbf8ef"
INK = "#3b2a6b"

gen.ICONS.update({
 "settings": "...##...|.#.##.#.|..####..|###..###|###..###|..####..|.#.##.#.|...##...",
 "trash": "..####..|########|.#....#.|.#.#.##.|.#.#.##.|.#.#.##.|.#....#.|..####..",
 "select": "#.......|##......|###.....|####....|#####...|###.....|#.##....|...##...",
 "lasso": ".#####..|#.....#.|#......#|#......#|.#....#.|..####..|..#.....|.##.....",
 "lasso-fill": ".#####..|#######.|########|########|.######.|..####..|..#.....|.##.....",
 "grip": "........|.#..#...|........|.#..#...|........|.#..#...|........|........",
 "move": "...##...|..####..|...##...|##.##.##|##.##.##|...##...|..####..|...##...",
 "rotate": "..####..|.#....#.|#......#|#.......|#.....#.|.#...###|..####..|........",
})

# ================= pen settings pop-up =================
PP = {
 "en": dict(title="PEN SETTINGS", pen="ink", pressure="PRESSURE", presets=["soft", "linear", "firm"], width="width", opacity="opacity",
            taper="taper", smoothing="smoothing", speed="speed", inputs=["Trackpad", "Tablet", "Mouse"], floor="floor",
            floor_hint="light taps still draw", reset="Reset", name="name", save="Save as new pen", newname="soft ink", raw="raw", eff="out",
            pr="pressure", inp="trackpad", test="test only", hint="draw here to try", hint2="pad strokes stay here", clear="clear",
            code="CODE", codefor="pens.soft_ink",
            stab="STABILIZATION", modes=["Off", "String", "Average"], length="length", catch="catch up at the end",
            readhint="the sky line on the curve"),
 "es": dict(title="AJUSTES DE PLUMA", pen="tinta", pressure="PRESIÓN", presets=["suave", "lineal", "firme"], width="grosor", opacity="opacidad",
            taper="punta", smoothing="suavizado", speed="velocidad", inputs=["Trackpad", "Tableta", "Ratón"], floor="mínimo",
            floor_hint="un toque leve pinta", reset="Restaurar", name="nombre", save="Guardar pluma nueva", newname="tinta suave", raw="entrada",
            eff="salida", pr="presión", inp="trackpad", test="de prueba", hint="dibuja aquí para probar", hint2="no van al dibujo", clear="borrar",
            code="CÓDIGO", codefor="pens.soft_ink",
            stab="ESTABILIZACIÓN", modes=["No", "Cuerda", "Promedio"], length="largo", catch="alcanzar al final",
            readhint="la línea azul de la curva"),
}
LIVE = 0.62
POP_W = 480

def title_bar(x, y, w, left_items_fn, right_text=None, focused=True):
    """Floating window title bar: raised strip, lizard top line when focused, grip dots, close."""
    it = [R(x + 1, y + 1, w - 2, 13, "raised", "title bar bg"), R(x + 1, y + 14, w - 2, 1, "void", "title bar line")]
    if focused:
        it.append(R(x + 1, y + 1, w - 2, 1, "lizard", "focused window top line"))
    it.append(I("grip", x + 4, y + 4, "dim", "drag grip (move window)"))
    it.append(left_items_fn(x + 14, y + 1))
    rx = x + w - 12
    it.append(I("close", rx, y + 4, "mute", "close ✕"))
    if right_text:
        rx -= 8 + tw(right_text)
        it.append(R(rx - 5, y + 4, 1, 8, "bevel", "sep"))
        it.append(T(right_text, rx, y + 1, "mute", "reset (text button)"))
    return G("title bar (drag to move)", *it)

def float_chrome(x, y, w, h, name):
    """Floating window: 1 px void outline, 3 px void drop shadow (hard, 8-bit), panel bg, bevel top-left highlight."""
    return [R(x + 3, y + 3, w, h, "void", "drop shadow (3 px)"), frame_lines(x, y, w, h, "void", "window outline"),
            R(x + 1, y + 1, w - 2, h - 2, "panel", name + " bg"), R(x + 1, y + 1, 1, h - 2, "bevel", "left highlight")]

def scribble(pts, raw, ws_max=12, mn=1, taper=0.3, floor=0.25):
    n = len(pts)
    eff = [floor + (1 - floor) * curve_eval(CURVE, r) for r in raw]
    ws = [(mn + (ws_max - mn) * e) * min(1, max(0.15, min(i / (n - 1), 1 - i / (n - 1)) / taper)) for i, e in enumerate(eff)]
    return variable_outline(pts, ws)

def test_pad(x, y, w, h, lang, empty=False):
    p = PP[lang]
    it = [frame_lines(x, y, w, h, "void", "pad outline"), R(x + 1, y + 1, w - 2, h - 2, PAPER, "pad paper"),
          R(x + 1, y + 1, w - 2, 1, "#e4dfd0", "pad inner top shade")]
    if empty:
        ht = p["hint"]
        it.append(T(ht, x + (w - tw(ht)) // 2, y + (h - 13) // 2, "dim", "hint: draw here to try (faint, empty pad)"))
    else:
        # three test scribbles in screen px inside a (w-2)x(h-2) UI box; vector, full-res
        iw, ih = w - 2, h - 2
        s1 = [(10 + i * 1.7, 16 + 9 * math.sin(i / 14)) for i in range(60)]
        r1 = [0.15 + 0.8 * math.sin(i / 59 * math.pi) ** 1.4 for i in range(60)]
        s2 = [(118 + i * 1.25, 46 - 20 * math.sin(i / 70 * math.pi) + 4 * math.sin(i / 6)) for i in range(70)]
        r2 = [0.2 + 0.75 * (i / 69) for i in range(70)]
        # loop-the-loop scribble, light pressure (shows the trackpad floor)
        s3 = [(222 + i * 0.85 + 7 * math.cos(i / 5), 32 + 9 * math.sin(i / 5) + i * 0.1) for i in range(72)]
        r3 = [0.05 + 0.15 * math.sin(i / 71 * math.pi) for i in range(72)]
        body = "".join(f'<path d="{scribble(s, r)}" fill="{INK}" fill-opacity="0.92"/>' for s, r in ((s1, r1), (s2, r2), (s3, r3)))
        it.append(vec("test strokes (vector, never saved to the drawing)", x + 1, y + 1, iw, ih, body))
        # live pen tip on the last stroke end
    # clear button, top-right corner inside the pad
    bx, by = x + w - 15, y + 2
    it.append(button(bx, by, 13, 13, "button clear test pad (trash)", icon="trash"))
    return G("test pad" + (" (empty)" if empty else " (3 test strokes)"), *it)

def checkbox(x, y, on, label, name, col="text"):
    it = [frame_lines(x, y + 2, 9, 9, "void", "box outline"), R(x + 1, y + 3, 7, 7, "moss" if on else "coal", "box face")]
    if on:
        it.append(pixel_poly([(x + 2, y + 6), (x + 3, y + 8), (x + 7, y + 4)], "lizard", "check"))
    it.append(T(label, x + 13, y, col))
    return G(name + (" (on)" if on else " (off)"), *it)

def pen_popup(x, y, lang, coder=True, empty=False):
    p = PP[lang]
    w = POP_W
    items = []
    def tl(tx, ty):
        return G("title", T(p["title"], tx, ty, "mute"), T("·", tx + tw(p["title"]) + 5, ty, "mute"),
                 T(p["pen"], tx + tw(p["title"]) + 17, ty, "bright", "title pen name"))
    items.append(title_bar(x, y, w, tl, p["reset"]))
    # --- column A: pressure curve editor
    gx, gy, gs = x + 8, y + 31, 84
    items.append(T(p["pressure"], x + 6, y + 16, "mute", "label pressure"))
    g = [frame_lines(gx - 1, gy - 1, gs + 2, gs + 2, "void", "graph outline"), R(gx, gy, gs, gs, "coal", "graph bg")]
    for k in (1, 2, 3):
        g.append(G(f"grid {k}", *[R(gx + i, gy + gs * k // 4, 1, 1, "bevel", "dot") for i in range(0, gs, 3)],
                   *[R(gx + gs * k // 4, gy + i, 1, 1, "bevel", "dot") for i in range(0, gs, 3)]))
    g.append(pixel_poly([(gx, gy + gs - 1), (gx + gs - 1, gy)], "dim", "linear reference"))
    dense = catmull(CURVE, 24)
    g.append(pixel_poly([(gx + px * (gs - 1), gy + (1 - py) * (gs - 1)) for px, py in dense], "lizard", "pressure curve"))
    if not empty:
        lx = gx + int(LIVE * (gs - 1))
        ly = gy + int((1 - curve_eval(CURVE, LIVE)) * (gs - 1))
        g.append(G(f"live pressure {LIVE} (pen down in the test pad)", *[R(lx, gy + j, 1, 2, "sky", "dash") for j in range(0, gs, 4)],
                   R(lx - 2, ly - 2, 5, 5, "void", "live dot outline"), R(lx - 1, ly - 1, 3, 3, "sky", "live dot on the curve")))
    for i, (px, py) in enumerate(CURVE):
        ax, ay = gx + int(px * (gs - 1)), gy + int((1 - py) * (gs - 1))
        g.append(G(f"curve point {i+1}", R(ax - 3, ay - 3, 7, 7, "void", "outline"), R(ax - 2, ay - 2, 5, 5, "lizard", "fill")))
    g.append(T(p["raw"], gx + gs - tw(p["raw"]) - 2, gy + gs - 14, "mute", "x axis label"))
    g.append(T(p["eff"], gx + 2, gy + 1, "mute", "y axis label"))
    items.append(G("pressure curve editor", *g))
    items.append(sfield(gx - 1, gy + gs + 4, gs + 2, p["presets"][0], "curve preset (soft / linear / firm)", chev=True))
    # --- column B: sliders + live readout
    bx, bw = x + 104, 192
    lw = max(tw(p[k]) for k in ("width", "opacity", "taper", "smoothing", "speed")) + 6
    vw = tw("80-100%") + 4
    sl_w = bw - lw - vw
    rows = [("width", 1 / 28, 12 / 28, "1-12"), ("opacity", 0.8, 1.0, "80-100%"), ("taper", 0.6, None, "0.3"),
            ("smoothing", 0.4, None, "40%"), ("speed", 0.2, None, "20%")]
    yy = y + 17
    for key, a, b, val in rows:
        items.append(G(f"setting {key}", T(p[key], bx, yy, "mute"), slider(bx + lw, yy - 1, sl_w, a, b, name=f"{key} slider"),
                       T(val, bx + bw - tw(val), yy, "amber")))
        yy += 18
    yy += 4
    items.append(R(bx, yy, bw, 1, "bevel", "rule"))
    yy += 5
    if empty:
        items.append(G("live pressure readout (idle)", T(p["pr"], bx, yy, "mute"), T("—", bx + tw(p["pr"]) + 7, yy, "dim"),
                       T("· " + p["inp"], bx + tw(p["pr"]) + 21, yy, "mute")))
    else:
        items.append(G("live pressure readout", T(p["pr"], bx, yy, "mute"), T(f"{LIVE:.2f}", bx + tw(p["pr"]) + 7, yy, "sky", "value (matches the sky line)"),
                       T("· " + p["inp"], bx + tw(p["pr"]) + 7 + tw("0.62") + 7, yy, "mute")))
    items.append(T(p["readhint"], bx, yy + 14, "dim", "readout hint"))
    # --- column C: stabilization + input overrides
    cx, cw = x + 304, w - 304 - 6
    yy = y + 16
    items.append(T(p["stab"], cx, yy, "mute", "label stabilization")); yy += 15
    sg, _ = seg(cx, yy, p["modes"], 1, "stabilizer mode (off / pulled string / average)")
    items.append(sg); yy += 17
    v = "24 px"
    sw_ = cw - tw(p["length"]) - 6 - tw(v) - 5
    items.append(G("setting string length", T(p["length"], cx, yy, "mute"), slider(cx + tw(p["length"]) + 6, yy - 1, sw_, 0.4, name="string length slider"),
                   T(v, cx + cw - tw(v), yy, "amber"))); yy += 18
    items.append(checkbox(cx, yy, True, p["catch"], "toggle catch up at the end of the stroke")); yy += 19
    sg, _ = seg(cx, yy, p["inputs"], 0, "input switch (overrides for)")
    items.append(sg); yy += 17
    fl_x = cx + 5 + tw(p["floor"]) + 6
    fl_w = cw - (fl_x - cx) - tw("25%") - 6
    items.append(G("trackpad overrides", R(cx, yy, cw, 31, "coal", "override bg"), R(cx, yy, 1, 31, "lizard", "override bar"),
                   T(p["floor"], cx + 5, yy + 1, "mute"), slider(fl_x, yy, fl_w, 0.25, name="floor slider"),
                   T("25%", cx + cw - 2 - tw("25%"), yy + 1, "amber"), T(p["floor_hint"], cx + 5, yy + 15, "mute", "floor hint")))
    for sx_ in (x + 98, x + 300):
        items.append(R(sx_, y + 18, 1, 128, "void", "column divider"))
    # --- row 2: test pad | name + save
    py = y + 150
    pw, ph = 300, 66
    items.append(test_pad(x + 6, py, pw, ph, lang, empty))
    rx2, rw2 = x + 6 + pw + 6, w - pw - 18
    items.append(G("field name", T(p["name"], rx2, py + 2, "mute"),
                   sfield(rx2 + tw(p["name"]) + 5, py, rw2 - tw(p["name"]) - 5, p["newname"], "pen name", h=18)))
    sv_w = tw(p["save"]) + 22
    sx = x + w - 6 - sv_w
    items.append(G("button primary " + p["save"], frame_lines(sx, py + 22, sv_w, 18, "void"), R(sx + 1, py + 23, sv_w - 2, 16, "moss", "face"),
                   gen.bevel(sx + 1, py + 23, sv_w - 2, 16, "lizard", "void"), I("plus", sx + 6, py + 27, "lizard"), T(p["save"], sx + 16, py + 24, "bright")))
    items.append(T(p["hint2"], rx2, py + ph - 13, "dim", "pad note: test strokes never go into the drawing"))
    yy = py + ph + 4
    if coder:
        lines = ['  "pressure": {"curve": [(0, 0), (0.3, 0.6), (1, 1)],',
                 '   "smoothing": 0.4, "speed": 0.2, "trackpad": {"floor": 0.25}},',
                 '  "stabilizer": {"mode": "string", "length": 24,', '   "catch_up": true}}']
        ch = len(lines) * 13 + 6
        cb = [frame_lines(x + 6, yy, w - 12, ch, "void", "inset outline"), R(x + 7, yy + 1, w - 14, ch - 2, "coal", "inset"),
              R(x + 7, yy + 3 + 2 * 13, w - 20, 13, "row", "changed line (just edited: mode)"), R(x + 7, yy + 3 + 2 * 13, 2, 13, "amber", "changed bar")]
        cb += [r002.code(ln, x + 11, yy + 3 + i * 13, name=f"code line {i+3}") for i, ln in enumerate(lines)]
        cb += [R(x + w - 11, yy + 2, 3, ch - 4, "coal", "scroll track"), R(x + w - 11, yy + 2 + (ch - 4) // 3, 3, (ch - 4) * 2 // 3, "bevel", "scroll thumb (lines 3-6 of 6)")]
        cb.append(T(p["codefor"] + " = {…", x + w - 16 - tw(p["codefor"] + " = {…"), yy + 3 + 3 * 13, "dim", "code label (scrolled: lines 3-6)"))
        items.append(G("code block (pens.soft_ink, scrolled to what changed)", *cb))
        yy += ch + 4
    h = yy + 1 - y
    return G("pen settings pop-up (floating window)", *float_chrome(x, y, w, h, "pen settings pop-up"), *items), h

def pen_tip(x, y, name="pen tip (drawing in the test pad)"):
    """Pen cursor: 1 px ring with a centre dot (brush size preview), void outlined."""
    return G(name, pixel_circle(x, y, 4, "void", "ring outline"), pixel_circle(x, y, 3, "bright", "ring"), R(x, y, 1, 1, "bright", "centre"))

def pen_row(lay, rows):
    return rows

def f_pen_popup_coder():
    lang = "en"
    lay = Lay("rail")
    rows = rows_layers()
    props = props2(lay, lang, "brush · pens.ink", lambda y: T("gear on a pen or double-click brush: settings", 6, y + 2, "mute"), None)
    px, py = 16, MAIN_Y + 3
    ps, h = pen_popup(px, py, lang, True)
    tip = pen_tip(px + 7 + 282, py + 151 + 48)
    return window2([gen.toolbar("brush", lang), tabbar2(unsaved=True), code_panel(lay, lang, rows, props=props),
                    canvas_area2(lay, drawing_g(lay, layer_bodies())), rail(lay), ps, tip, status2(lang, lay, saved=False, pos=None)])

def f_pen_popup_empty():
    """The pop-up alone (empty test pad: faint hint, no readout), on the coal window bg."""
    ps, h = pen_popup(4, 4, "en", False, empty=True)
    ww, hh = POP_W + 11, h + 11
    return svg_doc(ww * S, hh * S, ps, bg="coal", scale=S)

def f_pen_popup_artist_es():
    lang = "es"
    lay = art_lay()
    b = layer_bodies()
    stroke = gen.brush_shape([(70, 160), (160, 100), (250, 140), (340, 70), (430, 120)], [0.2, 0.9, 0.5, 1.0, 0.3], 1, 12, 0.3, INK, 1,
                             "brush stroke (tinta suave)")
    def ins(x, y):
        return G("fields (brush)", T("pluma", x, y + 1, "mute"), sfield(x + 42, y, 100, "tinta suave", "pen", chev=True),
                 button(x + 145, y + 1, 14, 14, "button pen settings ⚙ (open)", icon="settings", state="active"),
                 T("color", x + 166, y + 1, "mute"), sfield(x + 204, y, tw(INK) + 22, INK, "color", swatch=INK))
    ps, h = pen_popup(16, 62, lang, False)
    tip = pen_tip(16 + 7 + 282, 62 + 151 + 48)
    return window2([toolbar_artist(lang), tabbar2(unsaved=True), inspector(lang, ins, "brush", "pincel · sun"), tool_column("brush", INK),
                    canvas_area2(lay, drawing_g(lay, b, extra=stroke), splitter=False),
                    right_column(lang, INK, art_layers(b, sel="sun")), ps, tip, status2(lang, lay, saved=False, pos=None)])

# ================= pen picker with a gear per row =================
def f_picker_gear():
    lang = "en"
    lay = Lay("rail")
    rows = rows_layers()
    bx = 3 + 15 * 5 + 4 + 15
    px, py = bx, TOP_H + 1
    pens = [("ink", "crisp, tapers"), ("pencil", "light, grainy"), ("marker", "see-through"), ("brush", "soft, wide"), ("calligraphy", "45° nib"),
            ("soft ink", "yours")]
    pw, rh = 230, 20
    hh = len(pens) * rh + 18 + 14
    items = [R(px + 2, py + 2, pw, hh, "void", "drop shadow"), frame_lines(px, py, pw, hh, "void"),
             R(px + 1, py + 1, pw - 2, hh - 2, "row", "menu bg"), R(px + 1, py + 1, pw - 2, 13, "raised", "menu header"),
             T("PEN · pens.list", px + 5, py, "mute")]
    hover = "calligraphy"
    for i, (pn, desc) in enumerate(pens):
        ry = py + 16 + i * rh
        it = []
        cur, hov = pn == "ink", pn == hover
        if cur:
            it += [R(px + 1, ry, pw - 2, rh, "moss", "item bg (current)"), R(px + 1, ry, 2, rh, "lizard", "current bar")]
        if hov:
            it += [R(px + 1, ry, pw - 2, rh, "raised", "item bg (hover)")]
        it.append(G("preview frame", R(px + 6, ry + 2, 50, 16, "coal", "preview bg"), frame_lines(px + 5, ry + 1, 52, 18, "void")))
        kind = "ink" if pn == "soft ink" else pn
        it.append(f'<g id="stroke preview (full-res)" transform="scale(0.5)">{gen.pen_preview(kind, (px + 9) * S, (ry + 5) * S, 88, 16, C["text"])}</g>')
        it.append(T(pn, px + 62, ry + 1, "bright" if cur else "text"))
        if tw(pn) + tw(desc) + 7 < pw - 66 - 18:
            it.append(T(desc, px + 62 + tw(pn) + 7, ry + 1, "mute", f"desc {desc}"))
        # gear: every row; dim at rest, a real button on hover / current
        gx = px + pw - 18
        if hov:
            it.append(button(gx - 1, ry + 3, 14, 14, "button pen settings ⚙ (hover)", icon="settings", state="hover"))
        else:
            it.append(I("settings", gx + 2, ry + 6, "text" if cur else "dim", "pen settings ⚙"))
        it_name = f"pen {pn}" + (" (current)" if cur else "") + (" (hover)" if hov else "")
        items.append(G(it_name, *it))
    fy = py + 16 + len(pens) * rh
    items.append(G("footer hint", R(px + 1, fy, pw - 2, 1, "bevel", "separator"),
                   T("double-click 🖌 = settings of the current pen", px + 5, fy + 1, "mute") if False else
                   T("double-click brush: its settings", px + 5, fy + 1, "mute")))
    picker = G("pen picker (open, ⚙ per row)", *items)
    gx = px + pw - 18
    hy = py + 16 + 4 * rh
    tip = gen.tooltip_box(gx - 150, py + hh + 4, [lambda a, b: T("Pen settings…", a, b, "bright"),
                                                  lambda a, b: T("opens the pop-up for calligraphy", a, b, "mute")],
                          "tooltip ⚙", w=tw("opens the pop-up for calligraphy") + 14)
    props = props2(lay, lang, "brush · pens.ink", lambda y: T("pick a pen · the gear tunes it", 6, y + 2, "mute"), None)
    return window2([gen.toolbar("brush", lang), tabbar2(unsaved=True), code_panel(lay, lang, rows, props=props),
                    canvas_area2(lay, drawing_g(lay, layer_bodies())), rail(lay), picker, tip, cursor(gx + 6, hy + 10),
                    status2(lang, lay, saved=False, pos=None)])

FRAMES = {"p18b": f_pen_popup_coder, "p18b2": f_pen_popup_empty, "p19b": f_pen_popup_artist_es, "p18c": f_picker_gear}

if __name__ == "__main__":
    import r003
    FRAMES.update(r003.FRAMES)
    only = sys.argv[2:] or list(FRAMES)
    for k in only:
        s = FRAMES[k]()
        open(os.path.join(OUT, k + ".svg"), "w").write(s)
        print(k, len(s))
