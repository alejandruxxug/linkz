# R003 frames: select / move, lasso (select + fill) with the string stabilizer, brush with the
# pulled-string visual, color wheel + harmonies in the color popover, trace image dialog.
import math, colorsys
import gen, r002
from gen import R, G, T, I, C, tw, button, field, frame_lines, svg_doc, W, H, S, ROW_H, MAIN_Y, STATUS_Y, TOP_H, sel_box, cursor, nm, variable_outline, catmull
from r002 import (Lay, chrome, sfield, slider, seg, pixel_poly, pixel_circle, vec, props2, canvas_area2, drawing_g, layer_bodies,
                  rail, status2, tabbar2, window2, sw, hexrgb, rgbhex, mix, TX, PROPS_Y, var_grid, VARS, RECENT12, tail_px, CP,
                  toolbar_artist, inspector, right_column, art_layers, art_lay, hsv_of, COLORS_ROW)
from gen_r002b import checkbox, float_chrome, title_bar, INK, PAPER

TOOLS0 = list(gen.TOOLS)
def toolbar(active, lang="en"):
    gen.TOOLS[:] = [("select",), ("lasso",), None] + TOOLS0
    try:
        return gen.toolbar(active, lang)
    finally:
        gen.TOOLS[:] = TOOLS0

def tool_column_r3(active, fill="#fff3c4"):
    items = [R(0, 50, 20, STATUS_Y - 50, "panel", "tool column bg"), R(19, 50, 1, STATUS_Y - 50, "void", "tool column line")]
    y = 53
    for tl in ["select", "lasso", None, "circle", "rect", "line", "polygon", "stroke", None, "pencil", "brush", None, "text", "layer"]:
        if tl is None:
            items.append(R(5, y + 1, 10, 1, "bevel", "separator")); y += 4; continue
        items.append(button(3, y, 14, 14, f"tool {tl}" + (" (active)" if tl == active else ""), icon=tl, state="active" if tl == active else "normal"))
        y += 15
    cy = STATUS_Y - 27
    items.append(G("fill / stroke chips", sw(7, cy + 5, 10, INK, "stroke chip"), sw(3, cy, 10, fill, "fill chip (front)")))
    return G("tool strip (vertical, select + lasso first)", *items)

def code_panel_multi(lay, lang, rows, sel_lines=(), hover_lines=(), props=None, new_lines=()):
    t = TX[lang]
    w = lay.cw
    items = [R(0, MAIN_Y, w, STATUS_Y - MAIN_Y, "panel", "panel bg"),
             G("panel header", R(0, MAIN_Y, w, 13, "raised", "header bg"), R(0, MAIN_Y + 13, w, 1, "void", "header line"),
               T(t["layers"], 6, MAIN_Y, "mute"), T("layers.lzdraw", w - tw("layers.lzdraw") - 6, MAIN_Y, "dim"))]
    y = MAIN_Y + 15
    rows_g = []
    for r in rows:
        if r["kind"] == "layer": y += gen.LAYER_GAP
        if y + ROW_H > PROPS_Y - ROW_H - 2: break
        el, h = gen.row_el(r, y, w=w, sel=r.get("line") in sel_lines, hover=r.get("line") in hover_lines)
        rows_g.append(el)
        if r.get("line") in new_lines:
            bw = tw("new") + 4
            rows_g.append(G("badge new row", R(w - bw - 6, y + 4, bw, 10, "lizard", "badge"), T("new", w - bw - 4, y + 3, "void")))
        y += h
    items.append(G("rows", *rows_g))
    py = PROPS_Y - ROW_H - 1
    el, _ = gen.row_el(dict(kind="canvas", d=0, code=r002.CANVAS_CODE, line=5), py, w=w)
    items.append(G("pinned canvas row (line 5)", R(0, py - 1, w, 1, "void", "pin line"), R(0, py, w, ROW_H, "coal", "pinned bg"), el))
    items.append(props or "")
    return G("layers panel (the code)", *items)

def dotted_circle(cx, cy, r, color, name, every=3):
    pts = set()
    for k in range(int(2 * math.pi * r * 2)):
        a = k / (2 * r * 2)
        pts.add((int(round(cx + r * math.cos(a))), int(round(cy + r * math.sin(a)))))
    pts = sorted(pts, key=lambda p: math.atan2(p[1] - cy, p[0] - cx))
    keep = [p for i, p in enumerate(pts) if i % every == 0]
    return f'<path id="{nm(name)}" fill="{C.get(color, color)}" d="{"".join(f"M{a} {b}h1v1h-1z" for a, b in keep)}"/>'

def handles(x0, y0, x1, y1, rot=True):
    """8 size handles (5x5 bright, void outline) + optional rotation handle above the top edge."""
    xs, ys = (x0, (x0 + x1) // 2, x1 - 1), (y0, (y0 + y1) // 2, y1 - 1)
    it = []
    for i, x in enumerate(xs):
        for j, y in enumerate(ys):
            if i == 1 and j == 1: continue
            it.append(G(f"size handle {i}{j}", R(x - 3, y - 3, 7, 7, "void", "outline"), R(x - 2, y - 2, 5, 5, "bright", "fill")))
    if rot:
        mx = (x0 + x1) // 2
        it.append(G("rotation handle (optional)", R(mx, y0 - 9, 1, 6, "lizard", "stalk"), pixel_circle(mx, y0 - 12, 3, "void", "ring outline"),
                    pixel_circle(mx, y0 - 12, 2, "lizard", "ring")))
    return G("transform handles", *it)

def ants(pts, name="marching ants (lasso)"):
    """1 px dashed path, bright 3 / void 2 alternating (marching ants)."""
    px = []
    for (x0, y0), (x1, y1) in zip(pts, pts[1:]):
        n = int(max(abs(x1 - x0), abs(y1 - y0))) + 1
        for i in range(n):
            t = i / n
            p = (int(round(x0 + (x1 - x0) * t)), int(round(y0 + (y1 - y0) * t)))
            if not px or px[-1] != p: px.append(p)
    on = [p for i, p in enumerate(px) if i % 5 < 3]
    off = [p for i, p in enumerate(px) if i % 5 >= 3]
    return G(name, f'<path id="ants void" fill="{C["void"]}" d="{"".join(f"M{a} {b}h1v1h-1z" for a, b in off)}"/>',
             f'<path id="ants bright" fill="{C["bright"]}" d="{"".join(f"M{a} {b}h1v1h-1z" for a, b in on)}"/>')

def string_viz(tip, cur, radius, name="stabilizer string (pulled string)"):
    """Cursor ahead, 1 px sky string back to the tip, dotted sky circle = string length at the cursor."""
    (tx, ty), (cx, cy) = tip, cur
    return G(name, dotted_circle(cx, cy, radius, "sky", "string radius (dotted)", 2),
             pixel_poly([tip, cur], "sky", "string (cursor → tip)"),
             R(tx - 2, ty - 2, 5, 5, "void", "tip outline"), R(tx - 1, ty - 1, 3, 3, "bright", "pen tip (ink is laid here)"),
             R(cx - 1, cy - 1, 3, 3, "sky", "cursor dot"))

def stab_strip(x, y, lang="en", label="stabilizer"):
    sg, ex = seg(x + tw(label) + 6, y, ["Off", "String", "Average"], 1, "stabilizer mode")
    return G("stabilizer control", T(label, x, y, "mute"), sg), ex

ROWS_BOAT = lambda sun=("circle((360, 90), 34, #fff3c4)", "circle((360, 90), 56, #f2a65a66)"), extra_boat=(): (
    [dict(kind="value", d=0, code="let horizon = 200", line=7), dict(kind="blank", line=8),
     dict(kind="layer", d=0, code="layer sun", line=9, open=True),
     dict(kind="circle", d=1, code=sun[0], line=10), dict(kind="circle", d=1, code=sun[1], line=11), dict(kind="blank", line=12),
     dict(kind="layer", d=0, code="layer boat", line=13, open=True),
     dict(kind="polygon", d=1, code="polygon([(120, 214), (180, 214), (165, 232), (135, 232)], #3b2a6b)", line=14),
     dict(kind="line", d=1, code="line((150, 214), (150, 170), #3b2a6b, 3)", line=15),
     dict(kind="polygon", d=1, code="polygon([(152, 174), (152, 208), (182, 208)], #fbf8ef)", line=16)] + list(extra_boat) +
    [dict(kind="blank", line=18), dict(kind="layer", d=0, code="layer sea", line=19, open=False, count=6)])

def lay_big():
    return Lay(None, cw=256, cx1=W, zoom=0.5)

def ui(lay, p):
    x, y = lay.cv(p); return (int(round(x)), int(round(y)))

# ---------- R003-1 select + move ----------
def f_select_move():
    lang = "en"; lay = lay_big()
    new = (330, 130)
    rows = ROWS_BOAT(sun=("circle((330, 130), 34, #fff3c4)", "circle((330, 130), 56, #f2a65a66)"))
    b = layer_bodies()
    b["sun"] = (f'<circle cx="{new[0]}" cy="{new[1]}" r="56" fill="#f2a65a" fill-opacity="0.4"/>'
                f'<circle cx="{new[0]}" cy="{new[1]}" r="34" fill="#fff3c4"/>')
    ox, oy = ui(lay, (360, 90))
    ghost = G("ghost (where they were)", dotted_circle(ox, oy, int(56 * lay.z), "mute", "ghost glow"), dotted_circle(ox, oy, int(34 * lay.z), "mute", "ghost core"))
    x0, y0 = ui(lay, (new[0] - 56, new[1] - 56)); x1, y1 = ui(lay, (new[0] + 56, new[1] + 56))
    box = G("selection box (2 objects)", sel_box(x0 * S, y0 * S, x1 * S, y1 * S), handles(x0, y0, x1, y1))
    cx, cy = ui(lay, (new[0] + 8, new[1] + 10))
    tip = gen.tooltip_box(cx - 40, cy + 16, [lambda a, bb: G("t", T("move", a, bb, "bright"), T("-30, +40", a + tw("move") + 7, bb, "amber"))],
                          "tooltip move delta", w=tw("move -30, +40") + 14)
    def pf(y):
        return G("fields (2 selected)", T("2 objects", 6, y + 2, "bright"), T("·", 6 + tw("2 objects") + 5, y + 2, "mute"),
                 T("move", 6 + tw("2 objects") + 17, y + 2, "mute"),
                 *gen.point_field(6 + tw("2 objects") + 17 + tw("move") + 6, y, "", "-30", "+40", "move delta")[:1])
    props = props2(lay, lang, "select · lines 10, 11", pf, "Shift-click adds · drag empty space: marquee")
    return window2([toolbar("select", lang), tabbar2(unsaved=True), code_panel_multi(lay, lang, rows, sel_lines=(10, 11), props=props),
                    canvas_area2(lay, drawing_g(lay, b), G("overlay", ghost, box)), cursor(cx, cy), tip,
                    status2(lang, lay, saved=False, pos="x 338 y 140")])

# ---------- lasso loop ----------
LOOP = [(104, 162), (124, 150), (152, 148), (178, 156), (196, 176), (200, 202), (196, 224), (182, 240), (158, 246), (130, 244),
        (108, 236), (98, 218), (96, 196), (100, 176)]

def wobble(pts, amp=2.2, seed=1):
    dense = catmull(pts, 6)
    return [(x + amp * math.sin(i * 1.7 + seed), y + amp * math.cos(i * 1.3 + seed)) for i, (x, y) in enumerate(dense)]

def lasso_props(lay, lang, mode, extra=None):
    def pf(y):
        sg, ex = seg(6, y, ["Select", "Fill"], 0 if mode == "select" else 1, "lasso mode", icons=["lasso", "lasso-fill"])
        cb = checkbox(ex + 10, y + 1, True, "catch up", "toggle catch up at the end")
        return G("lasso options", sg, cb)
    def hint_row():
        return ""
    p = props2(lay, lang, "lasso " + mode + (" ·   #8a6fb0" if mode == "fill" else ""), pf, None)
    y = PROPS_Y + 41
    sg, ex = seg(6, y, ["Off", "String", "Average"], 1, "stabilizer mode (off / pulled string / average)")
    row2 = [sg, slider(ex + 4, y, 256 - ex - 4 - tw("24 px") - 10, 0.4, name="string length slider"), T("24 px", 256 - 6 - tw("24 px"), y + 1, "amber")]
    if extra: row2.append(extra())
    return G("properties (lasso options strip)", p, *row2)

def f_lasso_select():
    lang = "en"; lay = lay_big()
    rows = ROWS_BOAT()
    b = layer_bodies()
    path = [ui(lay, p) for p in catmull(LOOP[:12], 5)]
    tipp = path[-1]
    cur = (tipp[0] - 10, tipp[1] - 9)
    raw = [ui(lay, p) for p in wobble(LOOP[9:12] + [(92, 200), (92, 186)], 3.0)]
    ov = G("overlay",
           G("objects inside the loop (will be selected)", *[sel_box(*[v * S for v in ui(lay, (120, 170))], *[v * S for v in ui(lay, (183, 233))], name="pre-selection box")]),
           pixel_poly(raw, "dim", "raw cursor path (wiggly, not used)"),
           ants(path), pixel_poly([tipp, path[0]], "bevel", "closing line (on release)"),
           string_viz(tipp, cur, 13))
    props = lasso_props(lay, lang, "select")
    return window2([toolbar("lasso", lang), tabbar2(), code_panel_multi(lay, lang, rows, hover_lines=(14, 15, 16), props=props),
                    canvas_area2(lay, drawing_g(lay, b), ov), status2(lang, lay, pos="3 inside")])

def f_lasso_fill():
    lang = "en"; lay = lay_big()
    loop = [(88, 238), (120, 232), (160, 230), (196, 234), (214, 242), (200, 250), (160, 254), (118, 252), (92, 246)]
    pts = catmull(loop + [loop[0]], 4)
    poly_code = "polygon([(88, 238), (120, 232), (160, 230), (196, 234), (214, 242), …], #8a6fb0)"
    rows = ROWS_BOAT(extra_boat=[dict(kind="polygon", d=1, code=poly_code, line=17)])
    b = layer_bodies()
    shade = f'<polygon points="{" ".join(f"{x:.0f},{y:.0f}" for x, y in pts)}" fill="#8a6fb0" fill-opacity="0.85"/>'
    b["boat"] = shade + b["boat"]
    x0, y0 = ui(lay, (86, 228)); x1, y1 = ui(lay, (216, 256))
    ov = G("overlay (new polygon selected)", sel_box(x0 * S, y0 * S, x1 * S, y1 * S), handles(x0, y0, x1, y1, rot=False))
    props = lasso_props(lay, lang, "fill", extra=lambda: sw(6 + tw("PROPERTIES · lasso fill · ") + 2, PROPS_Y + 4, 10, "#8a6fb0", "current color #8a6fb0 (header)"))
    return window2([toolbar("lasso", lang), tabbar2(unsaved=True), code_panel_multi(lay, lang, rows, sel_lines=(17,), new_lines=(17,), props=props),
                    canvas_area2(lay, drawing_g(lay, b), ov), status2(lang, lay, saved=False, pos=None)])

# ---------- brush with the pulled string (artist) ----------
def f_brush_string():
    lang = "en"
    lay = art_lay()
    b = layer_bodies()
    pts = [(20, 170), (70, 152), (130, 160), (190, 144), (250, 152), (290, 140)]
    stroke = gen.brush_shape(pts, [0.3, 0.7, 0.9, 0.8, 0.7, 0.6], 1, 10, 0.2, INK, 1, "brush stroke in progress (smooth)")
    tip = ui(lay, (290, 140))
    cur = (tip[0] + 13, tip[1] - 8)
    wig = [(20, 170), (44, 156), (60, 166), (84, 144), (106, 164), (130, 148), (156, 166), (174, 138), (200, 156), (222, 140), (246, 160),
           (264, 134), (290, 152), (312, 128)]
    raw = [ui(lay, p) for p in wig]
    ov = G("overlay", pixel_poly(raw, "dim", "raw cursor path (wiggles, ignored)"), string_viz(tip, cur, 15))
    def ins(x, y):
        st, ex = stab_strip(x + 136, y + 1)
        return G("fields (brush)", T("pen", x, y + 1, "mute"), sfield(x + 28, y, 84, "soft ink", "pen", chev=True),
                 button(x + 115, y + 1, 14, 14, "button pen settings", icon="settings"), st,
                 T("24 px", ex + 6, y + 1, "amber"))
    tipbox = gen.tooltip_box(cur[0] - 120, cur[1] - 58, [lambda a, bb: T("ink follows the string", a, bb, "bright"),
                                                        lambda a, bb: T("release: catches up to the cursor", a, bb, "mute")],
                             "tooltip string stabilizer", w=tw("release: catches up to the cursor") + 14)
    return window2([toolbar_artist(lang), tabbar2(unsaved=True), inspector(lang, ins, "brush", "brush · boat"), tool_column_r3("brush", INK),
                    canvas_area2(lay, drawing_g(lay, b, extra=stroke), ov, splitter=False), right_column(lang, INK, art_layers(b, sel="boat")),
                    tipbox, status2(lang, lay, saved=False, pos="pressure 0.71")])

# ---------- color wheel + harmonies ----------
def region(name, x, y, w, h, at):
    """Smooth color content at screen resolution, skipping None (outside the shape)."""
    out = []
    nx, ny = int(w * S), int(h * S)
    for j in range(ny):
        run_c, run_s = None, 0
        for i in range(nx + 1):
            c = at(x + (i + 0.5) / S, y + (j + 0.5) / S) if i < nx else None
            if c != run_c:
                if run_c is not None:
                    out.append(f'<rect x="{x + run_s / S:g}" y="{y + j / S:g}" width="{(i - run_s) / S:g}" height="0.5" fill="{run_c}"/>')
                run_c, run_s = c, i
    return f'<g id="{nm(name)}">{"".join(out)}</g>'

def hue_hex(hh, s=1, v=1): return rgbhex([c * 255 for c in colorsys.hsv_to_rgb(hh % 1, s, v)])

HARM = [("complementary", [0, 0.5]), ("split", [0, 5 / 12, 7 / 12]), ("analogous", [-1 / 12, 0, 1 / 12]), ("triadic", [0, 1 / 3, 2 / 3]),
        ("square", [0, 0.25, 0.5, 0.75]), ("mono", [0, 0, 0])]

def harm_icon(offs, mono=False):
    g = [["."] * 8 for _ in range(8)]
    pts = [(3.5 + 3 * math.cos(2 * math.pi * o - math.pi / 2), 3.5 + 3 * math.sin(2 * math.pi * o - math.pi / 2)) for o in offs]
    if mono:
        for k in range(3): g[1 + k * 2][3] = g[1 + k * 2][4] = "#"
    else:
        for (a, b) in pts:
            for dx in (0, 1):
                for dy in (0, 1):
                    xi, yi = min(7, max(0, int(round(a - 0.5)) + dx)), min(7, max(0, int(round(b - 0.5)) + dy))
                    g[yi][xi] = "#"
    return "|".join("".join(r) for r in g)

for nmh, offs in HARM:
    gen.ICONS["h-" + nmh] = harm_icon(offs, nmh == "mono")

def harmony_colors(hexc, kind):
    hh, s, v = hsv_of(hexc)
    offs = dict(HARM)[kind]
    if kind == "mono":
        return [hue_hex(hh, s, v), hue_hex(hh, s * 0.6, min(1, v * 1.0)), hue_hex(hh, s, v * 0.6)]
    return [hue_hex(hh + o, s, v) for o in offs]

def wheel(x, y, d, hexc, kind="triadic"):
    """Hue ring + inner triangle (hue / white / black corners)."""
    r_out, r_in = d / 2, d / 2 - 7
    cx, cy = x + d / 2, y + d / 2
    hh, s, v = hsv_of(hexc)
    ang = -2 * math.pi * hh  # hue 0 at the right, counter-clockwise
    rt = r_in - 2
    A = (cx + rt * math.cos(ang), cy + rt * math.sin(ang))                      # pure hue
    B = (cx + rt * math.cos(ang + 2 * math.pi / 3), cy + rt * math.sin(ang + 2 * math.pi / 3))   # white
    Cc = (cx + rt * math.cos(ang - 2 * math.pi / 3), cy + rt * math.sin(ang - 2 * math.pi / 3))  # black
    base = [c * 255 for c in colorsys.hsv_to_rgb(hh, 1, 1)]
    def bary(px, py):
        (x1, y1), (x2, y2), (x3, y3) = A, B, Cc
        den = (y2 - y3) * (x1 - x3) + (x3 - x2) * (y1 - y3)
        a = ((y2 - y3) * (px - x3) + (x3 - x2) * (py - y3)) / den
        b = ((y3 - y1) * (px - x3) + (x1 - x3) * (py - y3)) / den
        return a, b, 1 - a - b
    def at(px, py):
        dx, dy = px - cx, py - cy
        rr = math.hypot(dx, dy)
        if r_in <= rr <= r_out:
            return hue_hex(-math.atan2(dy, dx) / (2 * math.pi))
        a, b, c = bary(px, py)
        if min(a, b, c) >= 0:
            return rgbhex([a * base[k] + b * 255 for k in range(3)])
        return None
    it = [pixel_circle(int(cx), int(cy), int(r_out) + 1, "void", "ring outer outline"), pixel_circle(int(cx), int(cy), int(r_in) - 1, "void", "ring inner outline"),
          region("hue ring + triangle (smooth color content)", x, y, d, d, at)]
    # markers: harmony colors on the ring, current = bright ring
    rm = (r_out + r_in) / 2
    offs = dict(HARM)[kind]
    for k, o in enumerate(offs):
        a = -2 * math.pi * (hh + o)
        mx, my = int(cx + rm * math.cos(a)), int(cy + rm * math.sin(a))
        col = hue_hex(hh + o)
        if k == 0:
            it.append(G("current hue marker", frame_lines(mx - 3, my - 3, 7, 7, "void", "outer"), frame_lines(mx - 2, my - 2, 5, 5, "bright", "ring")))
        else:
            it.append(G(f"harmony marker {k} {col}", R(mx - 3, my - 3, 7, 7, "void", "outline"), R(mx - 2, my - 2, 5, 5, col, "fill")))
    # current color in the triangle: weights hue v*s, white v*(1-s), black 1-v
    wa, wb, wc = v * s, v * (1 - s), 1 - v
    px_ = int(wa * A[0] + wb * B[0] + wc * Cc[0]); py_ = int(wa * A[1] + wb * B[1] + wc * Cc[1])
    it.append(G("triangle marker (saturation / value)", frame_lines(px_ - 3, py_ - 3, 7, 7, "void", "outer"), frame_lines(px_ - 2, py_ - 2, 5, 5, "bright", "ring")))
    return G("color wheel (hue ring + triangle)", *it)

def wheel_popover(x, y, lang, cur, kind="triadic", tail=None, target="code-row swatch"):
    c = CP[lang]
    w = 172
    ix, iw = x + 6, w - 12
    it = []
    yy = y + 3
    it.append(T(c["vars"], ix, yy, "mute", "label variables")); yy += 14
    g, gh = var_grid(ix + 1, yy, 4, 40, VARS, sel="sky", focus=None, chip_h=10)
    it.append(g); yy += gh + 2
    it.append(T(c["recent"], ix, yy, "mute", "label recent")); yy += 14
    it.append(G("recent colors", *[sw(ix + i * 10, yy, 9, col) for i, col in enumerate(RECENT12[:16])])); yy += 14
    d = 88
    it.append(wheel(ix, yy, d, cur, kind))
    rx = ix + d + 6; rw = iw - d - 6
    it.append(sfield(rx, yy, rw, cur[:7], "hex"))
    it.append(sfield(rx, yy + 18, 36, "100%", "alpha", "amber"))
    it.append(button(rx + rw - 14, yy + 19, 14, 14, "button eyedropper", icon="dropper"))
    it.append(G("old / new", frame_lines(rx, yy + 38, rw, 14, "void", "outline"), R(rx + 1, yy + 39, rw // 2 - 1, 12, "#e0607a", "old #e0607a"),
                R(rx + rw // 2, yy + 39, rw - rw // 2 - 1, 12, cur[:7], "new")))
    it.append(T("old  new", rx, yy + 53, "dim", "old / new label"))
    it.append(T("ring: hue", rx, yy + 67, "dim", "help ring"))
    yy += d + 4
    it.append(R(ix, yy, iw, 1, "bevel", "rule")); yy += 4
    it.append(G("label harmony", T("HARMONY", ix, yy, "mute"), T(kind, ix + iw - tw(kind), yy, "text", "current harmony name"))); yy += 15
    sg, _ = seg(ix, yy, [""] * 6, [k for k, _ in HARM].index(kind), "harmony (complementary / split / analogous / triadic / square / mono)",
                icons=["h-" + k for k, _ in HARM], w=(iw + 5) // 6)
    it.append(sg); yy += 18
    cols = harmony_colors(cur, kind)
    sx = ix
    sws = []
    for k, col in enumerate(cols):
        sws.append(G(f"harmony swatch {k+1} {col} (click: use, drag: palette)", frame_lines(sx, yy, 24, 14, "void", "outline"), R(sx + 1, yy + 1, 22, 12, col, col)))
        sx += 26
    it.append(G("harmony swatches", *sws))
    it.append(G("button + add all as variables", I("plus", sx + 4, yy + 3, "lizard"), T("add all", sx + 14, yy, "lizard"))); yy += 16
    it.append(T("as colors.sky_2, sky_3", ix, yy, "dim", "add-all hint")); yy += 14
    it.append(G("writes hint", T(c["writes"], ix, yy, "mute"), r002.code("colors.sky", ix + tw(c["writes"]) + 7, yy, name="ref"))); yy += 14
    h = yy + 3 - y
    out = chrome(x, y, w, h, name="color popover (wheel)") + it
    if tail: out.append(tail_px(tail[0], x, y, w, h, tail[1]))
    return G(f"color popover with wheel ({target})", *out), h

def f_wheel_popover():
    lang = "en"
    lay = Lay("rail")
    rows = r002.rows_layers(extra_top=[COLORS_ROW])
    ys = r002.row_heights(rows)
    sx = 46 + tw('let colors = {"sky": ')
    ring = G("swatch focus (clicked)", frame_lines(sx - 2, ys[6] + 2, 10, 10, "amber", "focus ring"))
    pop, h = wheel_popover(sx - 20, ys[6] + ROW_H + 6, lang, "#f2a65a", tail=("up", sx + 3))
    props = props2(lay, lang, "let colors · line 6", lambda y: T("click a swatch: color popover", 6, y + 2, "mute"), None)
    return window2([toolbar(None, lang), tabbar2(), r002.code_panel(lay, lang, rows, sel_line=6, props=props), ring,
                    canvas_area2(lay, drawing_g(lay, layer_bodies())), rail(lay), pop, status2(lang, lay, pos=None)])

# ---------- trace image dialog ----------
SKIN, SKIN_SH, HAIR, HAIR_SH, LINE, IRIS = "#fbe3d3", "#efc3ad", "#3b2a5c", "#2a1d45", "#1a1418", "#4a7fc1"

def head_svg(traced):
    """Simple original cel-shaded head, 240x180 box. traced=False: soft 'bitmap' look (gradients); True: flat fills + ink strokes."""
    p = []
    if not traced:
        p.append('<defs><radialGradient id="tsk" cx="120" cy="96" r="70" gradientUnits="userSpaceOnUse"><stop offset="0" stop-color="#fdeee3"/>'
                 '<stop offset="1" stop-color="#eebda6"/></radialGradient><linearGradient id="thr" x1="0" y1="20" x2="0" y2="140" gradientUnits="userSpaceOnUse">'
                 '<stop offset="0" stop-color="#5a4682"/><stop offset="1" stop-color="#2a1d45"/></linearGradient>'
                 '<linearGradient id="tbg" x1="0" y1="0" x2="0" y2="180" gradientUnits="userSpaceOnUse"><stop offset="0" stop-color="#dfe9f2"/><stop offset="1" stop-color="#c7d6e6"/></linearGradient></defs>')
    bg = "url(#tbg)" if not traced else "#d6e2ee"
    sk = "url(#tsk)" if not traced else SKIN
    hr = "url(#thr)" if not traced else HAIR
    p.append(f'<rect x="0" y="0" width="240" height="180" fill="{bg}"/>')
    # back hair
    p.append(f'<path d="M62 70 C60 120 66 150 74 172 L96 150 L92 120 Z M178 70 C180 120 174 150 166 172 L144 150 L148 120 Z" fill="{hr}"/>')
    # neck + collar
    p.append(f'<path d="M104 132 L136 132 L138 160 L102 160 Z" fill="{SKIN_SH if traced else sk}"/>')
    p.append('<path d="M70 180 L84 156 L120 166 L156 156 L170 180 Z" fill="#f2f5ec"/><path d="M106 160 L120 176 L134 160 L120 166 Z" fill="#e0607a"/>')
    # face
    p.append(f'<path d="M78 70 C78 112 92 138 120 146 C148 138 162 112 162 70 C150 50 90 50 78 70 Z" fill="{sk}"/>')
    if traced:
        p.append(f'<path d="M80 92 C86 120 100 136 120 146 C108 132 100 112 98 92 Z" fill="{SKIN_SH}"/>')
    # eyes
    for ex, flip in ((100, 1), (140, -1)):
        p.append(f'<ellipse cx="{ex}" cy="98" rx="9" ry="12" fill="#ffffff"/><ellipse cx="{ex + flip}" cy="100" rx="7" ry="10" fill="{IRIS}"/>'
                 f'<ellipse cx="{ex + flip}" cy="103" rx="3.5" ry="5" fill="#1c2c4a"/><circle cx="{ex + flip * 3 - 1}" cy="95" r="2.6" fill="#ffffff"/>')
    p.append('<path d="M114 128 Q120 131 126 128" fill="none" stroke="#b0505a" stroke-width="1.6" stroke-linecap="round"/>')
    # blush
    p.append('<ellipse cx="94" cy="116" rx="7" ry="3" fill="#f4a3a3" fill-opacity="0.55"/><ellipse cx="146" cy="116" rx="7" ry="3" fill="#f4a3a3" fill-opacity="0.55"/>')
    # front hair with spikes (bangs)
    bangs = "M66 78 C64 40 92 22 120 22 C150 22 178 40 174 80 L164 66 L160 90 L148 64 L138 88 L128 60 L118 86 L108 60 L98 88 L90 64 L80 90 L76 70 Z"
    p.append(f'<path d="{bangs}" fill="{hr}"/>')
    if traced:
        p.append(f'<path d="M84 40 C100 30 116 30 126 34 C112 36 98 44 90 56 Z" fill="#5a4682"/>')  # hair highlight band
        p.append(f'<path d="M150 66 L160 90 L164 66 L170 76 C172 60 166 46 158 40 Z" fill="{HAIR_SH}"/>')
    # ink lines
    def ink(pts, ws, op=1.0):
        d = catmull(pts, 3)
        n = len(d)
        wl = [ws[0] + (ws[1] - ws[0]) * math.sin(math.pi * i / (n - 1)) for i in range(n)]
        import re as _re
        dd = _re.sub(r'(\d+\.\d)\d+', r'\1', variable_outline(d, wl)).replace(".0 ", " ")
        return f'<path d="{dd}" fill="{LINE}" fill-opacity="{op}"/>'
    op = 1.0 if traced else 0.75
    lines = [([(78, 72), (80, 104), (94, 130), (120, 146), (146, 130), (160, 104), (162, 72)], (0.3, 2.4)),
             ([(66, 78), (66, 50), (88, 28), (120, 22), (152, 28), (174, 50), (174, 80)], (0.4, 2.6)),
             ([(90, 86), (100, 82), (110, 86)], (0.4, 2.2)), ([(130, 86), (140, 82), (150, 86)], (0.4, 2.2)),
             ([(118, 112), (121, 118), (119, 120)], (0.2, 1.0)),
             ([(80, 90), (84, 76), (90, 64)], (0.2, 1.6)), ([(98, 88), (104, 72), (108, 60)], (0.2, 1.6)),
             ([(118, 86), (122, 72), (128, 60)], (0.2, 1.6)), ([(138, 88), (143, 76), (148, 64)], (0.2, 1.6)),
             ([(160, 90), (162, 78), (164, 66)], (0.2, 1.6))]
    for pts, ws in lines:
        p.append(ink(pts, ws, op))
    return "".join(p)

def f_trace_dialog():
    lang = "en"
    lay = Lay("rail")
    base = [toolbar(None, lang), tabbar2(), r002.code_panel(lay, lang, r002.rows_layers(), props=props2(lay, lang, "canvas", lambda y: "", None)),
            canvas_area2(lay, drawing_g(lay, layer_bodies())), rail(lay), status2(lang, lay, pos=None)]
    dim = G("dim (void 62%)", f'<rect x="0" y="0" width="{W}" height="{H}" fill="{C["void"]}" fill-opacity="0.62"/>')
    w, h = 452, 270
    x, y = (W - w) // 2, 34
    it = chrome(x, y, w, h, bg="panel", name="trace dialog")
    it.append(G("header", R(x + 1, y + 1, w - 2, 14, "raised", "header bg"), R(x + 1, y + 15, w - 2, 1, "void", "header line"),
                I("image" if "image" in gen.ICONS else "canvas", x + 6, y + 4, "mute"), T("TRACE IMAGE", x + 18, y + 2, "mute"),
                T("·", x + 18 + tw("TRACE IMAGE") + 5, y + 2, "mute"), T("girl.png", x + 18 + tw("TRACE IMAGE") + 17, y + 2, "bright"),
                I("close", x + w - 12, y + 5, "mute", "close")))
    # preview (left)
    pvx, pvy, pw, ph = x + 8, y + 22, 240, 180
    split = 112
    it.append(frame_lines(pvx - 1, pvy - 1, pw + 2, ph + 2, "void", "preview frame"))
    it.append(vec("before: original image (left of split)", pvx, pvy, split, ph, head_svg(False)))
    it.append(vec("after: traced vector (right of split)", pvx + split, pvy, pw - split, ph,
                  f'<g transform="translate({-split} 0)">{head_svg(True)}</g>'))
    sxl = pvx + split
    it.append(G("split handle (drag)", R(sxl - 1, pvy, 1, ph, "void", "split shadow"), R(sxl, pvy, 1, ph, "bright", "split line"),
                frame_lines(sxl - 5, pvy + ph // 2 - 8, 11, 17, "void", "grip outline"), R(sxl - 4, pvy + ph // 2 - 7, 9, 15, "raised", "grip face"),
                R(sxl - 2, pvy + ph // 2 - 3, 1, 7, "text", "grip ◂"), R(sxl + 2, pvy + ph // 2 - 3, 1, 7, "text", "grip ▸")))
    for lb, lx in (("original", pvx + 3), ("traced", pvx + pw - tw("traced") - 9)):
        it.append(G(f"chip {lb}", R(lx, pvy + 3, tw(lb) + 6, 13, "row", "chip bg"), frame_lines(lx, pvy + 3, tw(lb) + 6, 13, "void", "chip outline"),
                    T(lb, lx + 3, pvy + 3, "text")))
    # source row under the preview
    sy = pvy + ph + 4
    it.append(G("source", I("canvas", pvx, sy + 4, "mute"), T("girl.png · 640×480", pvx + 11, sy + 1, "text"),
                button(pvx + pw - tw("Change…") - 8, sy, tw("Change…") + 8, 14, "button change image (or drop one)", label="Change…")))
    # options (right)
    ox, ow = pvx + pw + 10, w - pw - 26
    yy = pvy
    it.append(T("MODE", ox, yy, "mute")); yy += 15
    sg, _ = seg(ox, yy, ["Photo", "Line art"], 1, "trace mode", w=(ow + 1) // 2)
    it.append(sg); yy += 20
    lw = tw("line threshold") + 6
    vw = tw("1.5 px")
    sl_w = ow - lw - vw - 6
    rows = [("line threshold", 0.5, "0.50"), ("colors", 6 / 16, "6"), ("detail", 0.4, "0.8 px"), ("snap to shapes", 0.5, "1.5 px"), ("flow", 0.6, "0.6")]
    for lb, f, val in rows:
        it.append(G(f"setting {lb}", T(lb, ox, yy, "mute"), slider(ox + lw, yy - 1, sl_w, f, name=lb + " slider"), T(val, ox + ow - tw(val), yy, "amber")))
        yy += 18
    yy += 2
    it.append(checkbox(ox, yy, True, "show lines", "toggle show lines")); yy += 16
    it.append(checkbox(ox, yy, True, "show fills", "toggle show fills")); yy += 18
    it.append(R(ox, yy, ow, 1, "bevel", "rule")); yy += 5
    st = [T("38", ox, yy, "amber"), T("strokes", ox + 18, yy, "mute"), T("12", ox + 82, yy, "amber"), T("fills", ox + 100, yy, "mute"),
          T("5", ox, yy + 14, "amber"), T("primitives", ox + 11, yy + 14, "mute"), T("1,240", ox + 96, yy + 14, "amber"), T("ms", ox + 135, yy + 14, "mute")]
    it.append(G("stats (strokes, fills, primitives, ms)", *st))
    # stats + buttons (footer)
    fy = y + h - 26
    it.append(R(x + 1, fy - 4, w - 2, 1, "void", "footer line"))
    pb = tw("Open as drawing") + 12
    ib = tw("Insert as layer") + 20
    bx = x + w - 8 - pb
    it.append(G("button primary Open as drawing", frame_lines(bx, fy, pb, 18, "void"), R(bx + 1, fy + 1, pb - 2, 16, "moss", "face"),
                gen.bevel(bx + 1, fy + 1, pb - 2, 16, "lizard", "void"), T("Open as drawing", bx + 6, fy + 2, "bright")))
    it.append(button(bx - 4 - ib, fy + 2, ib, 14, "button insert as layer", icon="layer", label="Insert as layer"))
    it.append(T("or drop an image on this dialog", pvx, sy + 17, "dim", "drop hint"))
    return window2(base + [dim, G("trace image dialog", *it), cursor(sxl + 2, pvy + ph // 2 + 4)])

FRAMES = {"t1": f_select_move, "t2": f_lasso_select, "t3": f_lasso_fill, "t4": f_brush_string, "t5": f_wheel_popover, "t6": f_trace_dialog}
