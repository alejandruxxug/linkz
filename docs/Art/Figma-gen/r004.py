# R004 frames: layer hover preview. Hovering a layer row (coder code row or artist layer row) shows a
# hover card with only that layer rendered; the "isolation" option dims everything else on the canvas.
# Usage: python3 r004.py <outdir> [frame keys]
import os, sys, math
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import gen, r002, r003
from gen import R, G, T, I, C, tw, frame_lines, W, H, S, ROW_H, MAIN_Y, STATUS_Y, cursor, nm, checker
from r002 import (Lay, props2, canvas_area2, drawing_g, layer_bodies, status2, tabbar2, window2, TX, vec,
                  toolbar_artist, inspector, right_column, art_lay, a_layer_row)
from r003 import toolbar, code_panel_multi, tool_column_r3

OUT = sys.argv[1] if len(sys.argv) > 1 else os.path.join(os.path.dirname(os.path.abspath(__file__)), "out")
os.makedirs(OUT, exist_ok=True)

BG = "#f6d9a8"           # the drawing's canvas background
HALO = "circle((360, 90), 80, #fff3c433)"
CARD_W = 152             # 6 padding + 140 preview + 6
PV_W, PV_H = 140, 93     # preview at the drawing's 3:2 aspect ratio
DIM_OP = 0.70            # isolation: void at 70% over everything outside the layer

S4 = {
 "en": dict(objects=lambda n: f"{n} object" + ("" if n == 1 else "s"), hint="click: select layer", hidden="hidden",
            hidden_hint="only shown here", nothing="nothing selected", nothing_hint="hover a layer to preview it",
            line="line", iso="isolated"),
 "es": dict(objects=lambda n: f"{n} objeto" + ("" if n == 1 else "s"), hint="clic: elegir la capa", hidden="oculta",
            hidden_hint="solo se ve aquí", nothing="nada seleccionado", nothing_hint="pasa sobre una capa para verla",
            line="línea", iso="aislada"),
}

# the compact palette's SV square as two vector gradients (the pixel version is 680 KB of SVG)
import colorsys
def sv_square_vec(x, y, w, h, hexc, name="saturation/value square (vector)"):
    hh, sat, val = r002.hsv_of(hexc)
    base = r002.rgbhex([c * 255 for c in colorsys.hsv_to_rgb(hh, 1, 1)])
    g1, g2 = r002.gid("svh"), r002.gid("svv")
    body = (f'<g id="{nm(name)}"><defs><linearGradient id="{g1}" x1="0" y1="0" x2="1" y2="0"><stop offset="0" stop-color="#ffffff"/>'
            f'<stop offset="1" stop-color="{base}"/></linearGradient><linearGradient id="{g2}" x1="0" y1="0" x2="0" y2="1">'
            f'<stop offset="0" stop-color="#000000" stop-opacity="0"/><stop offset="1" stop-color="#000000"/></linearGradient></defs>'
            f'<rect x="{x}" y="{y}" width="{w}" height="{h}" fill="url(#{g1})"/><rect x="{x}" y="{y}" width="{w}" height="{h}" fill="url(#{g2})"/></g>')
    mx, my = x + int(sat * (w - 1)), y + int((1 - val) * (h - 1))
    marker = G("picker marker", frame_lines(mx - 3, my - 3, 7, 7, "void", "marker outer"), frame_lines(mx - 2, my - 2, 5, 5, "bright", "marker ring"))
    return [frame_lines(x - 1, y - 1, w + 2, h + 2, "void", "square frame"), body, marker]
# r002.sv_square = sv_square_vec  (kept for reference; raster.py needs plain fills)

# ---------- the drawing, with the sun's third object and the optional nested rays layer ----------
RAYS_COL = "#e0607a"
def rays_svg():
    out = []
    for k in range(8):
        a = k * math.pi / 4 + math.pi / 8
        x0, y0 = 360 + 42 * math.cos(a), 90 + 42 * math.sin(a)
        x1, y1 = 360 + 66 * math.cos(a), 90 + 66 * math.sin(a)
        out.append(f'<line x1="{x0:.1f}" y1="{y0:.1f}" x2="{x1:.1f}" y2="{y1:.1f}" stroke="{RAYS_COL}" stroke-width="4" stroke-linecap="round"/>')
    return f'<g id="layer rays (screen 70%)" opacity="0.7">{"".join(out)}</g>'

def bodies(rays=False, hide_boat=False):
    halo = '<circle id="halo line 12" cx="360" cy="90" r="80" fill="#fff3c4" fill-opacity="0.2"/>'
    b = layer_bodies()
    b["sun"] = halo + b["sun"] + (rays_svg() if rays else "")
    if hide_boat:
        b["boat_hidden"], b["boat"] = b["boat"], ""
    return b

def rows_sun(sun_open=True, nested=False, boat_hidden=False):
    rows = [dict(kind="value", d=0, code="let horizon = 200", line=7), dict(kind="blank", line=8),
            dict(kind="layer", d=0, code="layer sun", line=9, open=sun_open)]
    if sun_open:
        rows += [dict(kind="circle", d=1, code="circle((360, 90), 34, #fff3c4)", line=10),
                 dict(kind="circle", d=1, code="circle((360, 90), 56, #f2a65a66)", line=11),
                 dict(kind="circle", d=1, code=HALO, line=12)]
        if nested:
            rows += [dict(kind="layer", d=1, code="layer rays", line=13, open=True),
                     dict(kind="blend", d=2, code='blend("screen")', line=14), dict(kind="opacity", d=2, code="opacity(0.7)", line=15),
                     dict(kind="block", d=2, code="for i in range(8) {", line=16), dict(kind="line", d=3, code="line(ray(i, 42), ray(i, 66), #e0607a, 4)", line=17)]
    else:
        rows[-1]["count"] = 3
    base = 13 if not nested else 19
    rows += [dict(kind="blank", line=base), dict(kind="layer", d=0, code="layer boat", line=base + 1, open=False, count=3, hidden=boat_hidden),
             dict(kind="blank", line=base + 5), dict(kind="layer", d=0, code="layer sea", line=base + 6, open=False, count=6),
             dict(kind="blank", line=base + 12), dict(kind="layer", d=0, code="layer sky", line=base + 13, open=False, count=1)]
    return rows

# rows marked hidden get the eye-off icon (gen.row_el always draws an eye)
_row_el = gen.row_el
def row_el(r, y, **kw):
    el, h = _row_el(r, y, **kw)
    if r.get("hidden"):
        bg = "row" if kw.get("hover") else "panel"
        el = G(f'row line {r["line"]} · layer (hidden)', el, R(4, y + 5, 8, 8, bg, "eye cover"), I("eye-off", 4, y + 5, "amber", "toggle visible (eye off = hidden)"))
    return el, h
gen.row_el = row_el

def row_y(rows, line):
    y = MAIN_Y + 15
    for r in rows:
        if r["kind"] == "layer": y += gen.LAYER_GAP
        if r.get("line") == line: return y
        y += 4 if r["kind"] == "blank" else ROW_H
    raise KeyError(line)

# ---------- the hover card ----------
def tail_side(side, x, y, h, p):
    """5-step pixel tail out of the card's left or right edge, centred on p (y)."""
    tl = []
    for k in range(1, 6):
        span = 11 - 2 * k
        top = p - span // 2
        xx = x - k if side == "left" else x + CARD_W - 1 + k
        tl.append(R(xx, top, 1, span, "void", "tail edge"))
        if span > 2: tl.append(R(xx, top + 1, 1, span - 2, "row", "tail fill"))
    return G("tail → hovered row", *tl)

def hover_card(x, y, lang, name, body, count, badge=None, hidden=False, path=None, hint=None, tail=None, transparent=False):
    t = S4[lang]
    lines = 2 + (1 if badge else 0)  # hidden is marked by the chip on the preview
    h = 6 + PV_H + 2 + 4 + lines * 15 + 2
    it = [R(x + 2, y + 2, CARD_W, h, "void", "drop shadow 2px"), frame_lines(x, y, CARD_W, h, "void"),
          R(x + 1, y + 1, CARD_W - 2, h - 2, "row", "card bg"), R(x + 1, y + 1, CARD_W - 2, 1, "bevel", "top highlight")]
    px, py = x + 6, y + 6
    it.append(frame_lines(px - 1, py - 1, PV_W + 2, PV_H + 2, "void", "preview frame"))
    if transparent:
        it.append(checker(px, py, PV_W, PV_H, "#d4dacd", "#a0ab9b", 3, "transparency checker (no canvas background)"))
        bg = ""
    else:
        bg = f'<rect x="0" y="0" width="480" height="320" fill="{BG}"/>'
    it.append(vec(f"preview: only layer {name} (vector, full-res)", px, py, 480, 320, bg + body, PV_W / 480))
    if hidden:
        bw = tw(t["hidden"]) + 16
        it.append(G("hidden chip (on the preview)", R(px + 3, py + 3, bw, 13, "void", "chip outline"), R(px + 4, py + 4, bw - 2, 11, "raised", "chip bg"),
                    I("eye-off", px + 6, py + 5, "amber"), T(t["hidden"], px + 16, py + 2, "amber")))
    ty = py + PV_H + 4
    tx = x + 6
    it.append(I("layer", tx, ty + 3, "lizard", "layer icon"))
    tx += 12
    if path:
        it.append(T(path + "/", tx, ty, "mute", "parent layer")); tx += tw(path + "/")
    it.append(T(name, tx, ty, "bright", "layer name"))
    cnt = t["objects"](count)
    it.append(T(cnt, x + CARD_W - 6 - tw(cnt), ty, "mute", "shape count"))
    ty += 15
    if badge:
        it.append(G("blend badge", R(x + 6, ty + 1, tw(badge) + 4, 10, "raised", "badge bg"), T(badge, x + 8, ty, "text")))
        ty += 15
    hint = hint or t["hint"]
    it.append(T(hint, x + 6, ty, "mute", "hint: " + hint))
    if tail:
        it.append(tail_side(tail[0], x, y, h, tail[1]))
    return G(f"layer hover card · {name}", *it), h

# ---------- isolation overlay ----------
def isolation(lay, layer_body, name, lang, count):
    t = S4[lang]
    dim = f'<rect id="dim (void 70%)" x="{lay.dx}" y="{lay.dy}" width="{lay.dw}" height="{lay.dh}" fill="{C["void"]}" fill-opacity="{DIM_OP}"/>'
    redraw = f'<g id="{nm("layer " + name + " (drawn again above the dim)")}" data-size="480 320" transform="translate({lay.dx} {lay.dy}) scale({lay.z})">{layer_body}</g>'
    lab = f"{name} · {t['objects'](count)}"
    cw = tw(lab) + 20
    cx, cy = lay.dx, lay.dy + lay.dh + 4
    chip = G("isolation chip (under the canvas)", R(cx, cy, cw, 16, "void", "chip outline"), R(cx + 1, cy + 1, cw - 2, 14, "row", "chip bg"),
             R(cx + 1, cy + 1, cw - 2, 1, "bevel", "top highlight"), I("layer", cx + 5, cy + 4, "lizard"), T(lab, cx + 15, cy + 2, "text"))
    return G(f"isolation: everything outside layer {name} dimmed", dim, redraw, chip)

def nothing_props(lay, lang):
    t = S4[lang]
    return props2(lay, lang, t["nothing"], lambda y: T("canvas (480, 320)", 6, y + 2, "mute"), t["nothing_hint"])

# ---------- coder frames ----------
def f_coder(lang="en", iso=False, card=True, nested=False, hidden=False):
    t = S4[lang]
    lay = r003.lay_big()
    rows = rows_sun(nested=nested, boat_hidden=hidden)
    b = bodies(rays=nested, hide_boat=hidden)
    if nested:
        line, name, body, count, badge, path = 13, "rays", rays_svg(), 8, "SCR 70%", "sun"
    elif hidden:
        line, name, body, count, badge, path = 20 - 6, "boat", b["boat_hidden"], 3, None, None
    else:
        line, name, body, count, badge, path = 9, "sun", b["sun"], 3, None, None
    ry = row_y(rows, line)
    d = next(r["d"] for r in rows if r.get("line") == line)
    text_end = 46 + 10 * d + tw(f"layer {name}")
    parts = [toolbar(None, lang), tabbar2(), code_panel_multi(lay, lang, rows, hover_lines=(line,), props=nothing_props(lay, lang)),
             canvas_area2(lay, drawing_g(lay, b), isolation(lay, body, name, lang, count) if iso else "")]
    if nested:  # rows below line 17 are scrolled out of view
        top, bot = MAIN_Y + 15, r002.PROPS_Y - ROW_H - 3
        parts.append(G("scrollbar (more rows below)", R(lay.cw - 5, top, 4, bot - top, "coal", "track"), R(lay.cw - 5, top, 4, (bot - top) * 2 // 3, "bevel", "thumb")))
    if card:
        cx = text_end + 22
        el, h = hover_card(cx, ry - 6, lang, name, body, count, badge=badge, hidden=hidden, path=path,
                           hint=t["hidden_hint"] if hidden else t["hint"], tail=("left", ry + 9))
        parts.append(el)
    parts += [cursor(text_end + 3, ry + 6), status2(lang, lay, pos=None)]
    return window2(parts)

# ---------- artist frames ----------
def art_rows(b, hover="sun"):
    L = [dict(kind="layer", name="sun", body=b["sun"], open=False), dict(kind="layer", name="boat", body=b["boat"], open=False),
         dict(kind="layer", name="sea", body=b["sea"], open=False), dict(kind="layer", name="sky", body=b["sky"], open=False)]
    return L

def f_artist(lang="en", iso=False, card=True):
    t = S4[lang]
    lay = art_lay()
    b = bodies()
    rc = right_column(lang, "#fff3c4", art_rows(b))
    ry = 180  # first layer row (sun)
    hover_bg = G("layer row sun (hover)", R(r002.RC_X, ry, W - r002.RC_X, 22, "row", "row bg (hover)"))
    # redraw the hovered row on top of its hover background
    THUMBS = r002.THUMBS
    THUMBS.clear()
    row, _ = a_layer_row(r002.RC_X, ry, W - r002.RC_X, "sun", b["sun"], open_=False)
    row += G("thumb sun (vector)", *THUMBS); THUMBS.clear()
    ins = inspector(lang, lambda x, y: T(t["nothing_hint"], x, y + 1, "mute"), "layer", t["nothing"])
    parts = [toolbar_artist(lang), tabbar2(), ins, tool_column_r3("select"),
             canvas_area2(lay, drawing_g(lay, b), isolation(lay, b["sun"], "sun", lang, 3) if iso else "", splitter=False),
             rc, hover_bg, row]
    if card:
        cx = r002.RC_X - 7 - CARD_W
        el, h = hover_card(cx, ry - 4, lang, "sun", b["sun"], 3, hint=t["hint"], tail=("right", ry + 11))
        parts.append(el)
    parts += [cursor(r002.RC_X + 75, ry + 10), status2(lang, lay, pos=None)]
    return window2(parts)

FRAMES = {
 "v1": ("R004-1 Coder – hover card on layer sun", lambda: f_coder("en")),
 "v2": ("R004-2 Coder – canvas isolation only (no card)", lambda: f_coder("en", iso=True, card=False)),
 "v3": ("R004-3 Artist – hover card beside the layers panel", lambda: f_artist("en")),
 "v4": ("R004-4 Artist – card + isolation (both)", lambda: f_artist("en", iso=True)),
 "v5": ("R004-5 Coder – Spanish (3 objetos)", lambda: f_coder("es")),
 "v6": ("R004-6 Coder – nested layer sun/rays (only rays in the preview)", lambda: f_coder("en", nested=True)),
 "v7": ("R004-7 Coder – hidden layer boat (preview still shows it, marked hidden)", lambda: f_coder("en", hidden=True)),
}

# ---------- deltas: named groups cut out of a full frame, for cloning in Figma ----------
def extract(svg, gid):
    """The balanced <g id="gid">...</g> substring."""
    start = svg.index(f'<g id="{nm(gid)}">')
    depth, i = 0, start
    while True:
        o, c = svg.find("<g", i + 1 if depth else i), svg.find("</g>", i + 1)
        if o != -1 and o < c:
            depth += 1; i = o
        else:
            depth -= 1; i = c
            if depth == 0: return svg[start:c + 4]

def delta(svg, gids, name):
    body = "".join(extract(svg, g) for g in gids)
    return (f'<svg xmlns="http://www.w3.org/2000/svg" width="{W * S}" height="{H * S}" viewBox="0 0 {W * S} {H * S}">'
            f'<g id="{nm(name)}" transform="scale({S})">{body}</g></svg>')

# key: (source frame, groups to bring in, groups to delete in the clone of `base`)
DELTAS = {
 "v2": ("v1", ["isolation: everything outside layer sun dimmed"], ["layer hover card · sun"]),
 "v4": ("v3", ["isolation: everything outside layer sun dimmed", "layer hover card · sun"], ["layer hover card · sun"]),
 "v5": ("v1", ["button language", "panel header", "properties panel", "status bar", "layer hover card · sun"],
        ["button language", "panel header", "properties panel", "status bar", "layer hover card · sun"]),
 "v6": ("v1", ["canvas area", "rows", "layer hover card · rays", "mouse cursor"], ["canvas area", "rows", "layer hover card · sun", "mouse cursor"]),
 "v7": ("v1", ["canvas area", "row line 9 · layer", "row line 14 · layer (hidden)", "layer hover card · boat", "mouse cursor"],
        ["canvas area", "row line 9 · layer (hover)", "row line 14 · layer", "layer hover card · sun", "mouse cursor"]),
}

if __name__ == "__main__":
    only = sys.argv[2:] or list(FRAMES)
    for k in only:
        s = FRAMES[k][1]()
        open(os.path.join(OUT, k + ".svg"), "w").write(s)
        print(k, len(s))
    for k, (base, bring, _) in DELTAS.items():
        full = open(os.path.join(OUT, k + ".svg")).read()
        d = delta(full, bring, f"R004 delta {k}")
        open(os.path.join(OUT, "delta-" + k + ".svg"), "w").write(d)
        print("delta", k, len(d))
