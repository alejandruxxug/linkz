# Generates the Linkz R002 mockups as SVG (same 8-bit rules as gen.py / home.py):
# shape menu, fill editor (linear, stepped, radial), palette sidebar, artist mode
# (+ code drawer), blend modes, tutorial, new-drawing dialog, object budget,
# pen settings, autocomplete. Usage: python3 r002.py <outdir>
import os, sys, re, math, colorsys
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import gen, home
from gen import (R, G, T, I, C, tw, button, field, frame_lines, svg_doc, W, H, S, ROW_H, LAYER_GAP, MAIN_Y,
                 STATUS_Y, TOP_H, TAB_H, checker, sel_box, cursor, bezier_d, catmull, variable_outline, nm)

OUT = sys.argv[1] if len(sys.argv) > 1 else os.path.join(os.path.dirname(os.path.abspath(__file__)), "out")
os.makedirs(OUT, exist_ok=True)
PROPS_Y = 282

def mirror(s): return "|".join(r[::-1] for r in s.split("|"))
UP = "...##...|..####..|.######.|##.##.##|...##...|...##...|...##...|........"
gen.ICONS.update({
 "up": UP, "down": "|".join(reversed(UP.split("|"))),
 "fill": "...#....|..###...|.#####..|#######.|#######.|#######.|.#####..|..###...",
 "code": "........|..#.##..|.#..#.#.|#..#...#|.#.#..#.|..##.#..|........|........",
 "solid": "|".join(["########"] * 8),
 "gradient": "########|##.#.#.#|#.#.#..#|##.#...#|#.#.#..#|##.#...#|#.#.#..#|########",
 "radial": "..####..|.#....#.|#..##..#|#.####.#|#.####.#|#..##..#|.#....#.|..####..",
 "eye-off": "......#.|..####..|.#..#.#.|#..##..#|#.##...#|.#....#.|.#####..|#.......",
 "steps": "........|......##|......##|....####|....####|..######|..######|########",
 "blend": ".###....|#...#...|#..###..|#.#.#.#.|.###...#|..#....#|...####.|........",
 "opacity": "########|#...####|#...####|#...####|#...####|#...####|#...####|########",
 "info": "..####..|.#.##.#.|#......#|#..##..#|#..##..#|#..##..#|.#....#.|..####..",
 "chevron-left": mirror(gen.ICONS["chevron-right"]),
 "sparkle": "...#....|...#....|..###...|#######.|..###...|...#....|...#....|........",
 "dropper": "......##|.....###|....###.|...#.#..|..#.#...|.#.#....|#.#.....|.#......",
 "objects": "###..###|###..###|###..###|........|........|###..###|###..###|###..###",
 "minus": "........|........|........|.######.|.######.|........|........|........",
 "keyword": "........|.#......|.#..#...|.#.#....|.##.....|.#.#....|.#..#...|........",
 "pen": "......##|.....###|....###.|...###..|..###...|.###....|.##.....|#.......",
 "graph": "#.......|#.....#.|#....#..|#...#...|#..#....|#.#.....|##......|########",
 "warn": "...##...|...##...|..#..#..|..#..#..|.#....#.|.#.##..#|#......#|########",
})

# ---------- named colors: `colors.sky` gets a swatch like a literal ----------
NAMED = {"colors.sky": "#f2a65a", "colors.sea": "#2d6a8f", "colors.sun": "#fff3c4", "colors.boat": "#3b2a6b", "colors.dusk": "#e0607a", "colors.glow": "#fff3c4", "colors.dune": "#d8a070"}
_code = gen.code
def code(s, x, y, base="text", sel=False, maxw=None, name="code"):
    items, cx, groups, clipped = [], x, {}, False
    for m in gen.TOK.finditer(s):
        com, st, col, num, ident, ws, other = m.groups()
        tok = m.group(0)
        if com: color = "mute"
        elif st: color = "peach"
        elif col: color = base if not sel else "bright"
        elif num: color = "amber"
        elif ident: color = "sky" if ident in gen.KW else (base if not sel else "bright")
        else: color = "mute" if not sel else "text"
        swc = col if col else (NAMED.get(ident) if ident else None)
        if swc:
            if maxw and cx + 8 - x > maxw: clipped = True; break
            items.append(R(cx, y + 4, 6, 6, "void", "swatch outline"))
            items.append(R(cx + 1, y + 5, 4, 4, swc[:7], f"swatch {swc}"))
            cx += 8
        if maxw and cx + tw(tok) - x > maxw:
            tok = tok[:max((maxw - (cx - x)) // gen.GW - 1, 0)]; clipped = True
        if tok.strip():
            groups.setdefault(color, []).append(gen.text_d(tok, cx, y))
        cx += tw(tok)
        if clipped: break
    for color, ds in groups.items():
        items.insert(0, f'<path id="{nm("code " + color)}" fill="{C.get(color, color)}" d="{"".join(ds)}"/>')
    if clipped:
        items.append(T("»", x + maxw - 7, y, "mute", "clip marker »"))
    return G(f'{name}: {s}', *items)
gen.code = code
gen.KIND_ICON.update({"blend": "blend", "opacity": "opacity", "ellipse": "circle", "let": "function"})

# ---------- strings ----------
TX = {
 "en": dict(layers="LAYERS", props="PROPERTIES", line="line", saved="saved", unsaved="unsaved", objects="objects", lang="EN",
            color="COLOR", recent="RECENT", palette="VARIABLES", gradients="GRADIENTS", add_pal="Variable", used="used by",
            menu=[("rename", "Rename", "F2"), ("copy", "Duplicate", "Cmd+D"), ("up", "Bring forward", "Cmd+]"), ("down", "Send backward", "Cmd+["),
                  None, ("eye-off", "Hide", "H"), ("lock", "Lock", "L"), None, ("fill", "Fill", ">"), None,
                  ("code", "Copy code", ""), ("open", "Show in code", ""), None, ("trash", "Remove", "Del")],
            fill_sub=[("solid", "Solid color", ""), ("gradient", "Linear gradient", ""), ("radial", "Radial gradient", ""), None, ("fill", "Edit fill…", "F")],
            fill="FILL", tabs=["Solid", "Linear", "Radial"], at="at", soft="soft", hard="hard", steps="steps", off="off",
            hint1="click bar: add a stop", hint2="drag a stop off: remove", pick="stop 2 · pick in palette",
            code_btn="</> Code", export="Export PNG", save="Save"),
 "es": dict(layers="CAPAS", props="PROPIEDADES", line="línea", saved="guardado", unsaved="sin guardar", objects="objetos", lang="ES",
            color="COLOR", recent="RECIENTES", palette="VARIABLES", gradients="DEGRADADOS", add_pal="Variable", used="usado en",
            menu=[("rename", "Renombrar", "F2"), ("copy", "Duplicar", "Cmd+D"), ("up", "Traer adelante", "Cmd+]"), ("down", "Enviar atrás", "Cmd+["),
                  None, ("eye-off", "Ocultar", "H"), ("lock", "Bloquear", "L"), None, ("fill", "Relleno", ">"), None,
                  ("code", "Copiar código", ""), ("open", "Mostrar en el código", ""), None, ("trash", "Eliminar", "Supr")],
            fill_sub=[("solid", "Color sólido", ""), ("gradient", "Degradado lineal", ""), ("radial", "Degradado radial", ""), None, ("fill", "Editar relleno…", "F")],
            fill="RELLENO", tabs=["Sólido", "Lineal", "Radial"], at="en", soft="suave", hard="dura", steps="pasos", off="no",
            hint1="clic: añadir parada", hint2="arrastra fuera: quitar", pick="parada 2 · elige en paleta",
            code_btn="</> Código", export="Exportar PNG", save="Guardar"),
}
def fmt(n, lang): return f"{n:,}".replace(",", "." if lang == "es" else ",")

# ---------- layout ----------
class Lay:
    """Coder mode: code panel | canvas | palette sidebar (open 106, rail 14 or none)."""
    def __init__(s, sidebar="open", cw=256, cy0=MAIN_Y, cy1=STATUS_Y, cx0=None, cx1=None, zoom=None):
        s.cw, s.sidebar = cw, sidebar
        s.sb_w = {"open": 106, "rail": 14, None: 0}[sidebar]
        s.sb_x = W - s.sb_w
        s.cx0 = cw + 2 if cx0 is None else cx0
        s.cx1 = (s.sb_x - (2 if sidebar else 0)) if cx1 is None else cx1
        s.cy0, s.cy1 = cy0, cy1
        aw, ah = s.cx1 - s.cx0, s.cy1 - s.cy0
        s.z = zoom or next(z for z in (0.625, 0.45, 0.375, 0.275) if 480 * z <= aw - 10 and 320 * z <= ah - 10)
        s.dw, s.dh = int(480 * s.z), int(320 * s.z)
        s.dx = s.cx0 + (aw - s.dw) // 2
        s.dy = s.cy0 + (ah - s.dh) // 2
        s.pct = f"{int(round(s.z * 200))}%"
    def cv(s, p): return (s.dx + p[0] * s.z, s.dy + p[1] * s.z)
    def cvs(s, p): x, y = s.cv(p); return (x * S, y * S)

# ---------- small parts ----------
def chrome(x, y, w, h, bg="row", name="popup"):
    return [R(x + 2, y + 2, w, h, "void", "drop shadow"), frame_lines(x, y, w, h, "void"), R(x + 1, y + 1, w - 2, h - 2, bg, name + " bg"),
            R(x + 1, y + 1, w - 2, 1, "bevel", "top highlight")]

def header(x, y, w, title, close=True, icon=None):
    it = [R(x + 1, y + 1, w - 2, 13, "raised", "header bg"), R(x + 1, y + 14, w - 2, 1, "void", "header line")]
    tx = x + 5
    if icon:
        it.append(I(icon, x + 5, y + 4, "mute")); tx += 10
    it.append(T(title, tx, y + 1, "mute", "title " + title))
    if close:
        it.append(I("close", x + w - 12, y + 4, "mute", "close"))
    return G("header " + title, *it)

def sw(x, y, s, col, name=None):
    return G(name or f"swatch {col}", R(x, y, s, s, "void", "outline"), R(x + 1, y + 1, s - 2, s - 2, col[:7], col))

def sfield(x, y, w, text, name, color="text", focused=False, chev=False, swatch=None, h=16):
    items = []
    tx = x + 4
    if swatch:
        items += [R(x + 3, y + 3, 10, 10, "void", "swatch outline"), R(x + 4, y + 4, 8, 8, swatch[:7], "swatch " + swatch)]
        tx = x + 16
    items.append(T(text, tx, y + (h - 13) // 2 + 1, color))
    if chev:
        items.append(I("chevron-down", x + w - 11, y + (h - 8) // 2, "mute", "dropdown chevron"))
    return field(x, y, w, h, items, name, focused)

def slider(x, y, w, a, b=None, name="slider", col="lizard"):
    """Track 6 high centred in a 16 row, filled from 0 (or a) to a (or b); 5x12 thumbs."""
    lo, hi = (0, a) if b is None else (a, b)
    it = [frame_lines(x, y + 5, w, 6, "void", "track outline"), R(x + 1, y + 6, w - 2, 4, "coal", "track"),
          R(x + 1 + int((w - 2) * lo), y + 6, max(1, int((w - 2) * (hi - lo))), 4, col, "track fill")]
    for f in ([a] if b is None else [a, b]):
        px = x + 1 + int((w - 2) * f)
        it += [frame_lines(px - 2, y + 2, 5, 12, "void", "thumb outline"), R(px - 1, y + 3, 3, 10, "bright", "thumb")]
    return G(name, *it)

def seg(x, y, labels, active, name, h=14, icons=None, w=None):
    it, cx = [], x
    for i, lb in enumerate(labels):
        bw = w or (tw(lb) + 8 + (10 if icons else 0))
        it.append(button(cx, y, bw, h, f"{name} {lb}" + (" (active)" if i == active else ""), icon=icons[i] if icons else None,
                         label=lb, state="active" if i == active else "normal"))
        cx += bw - 1
    return G(name, *it), cx

def dither_dim(x, y, w, h, holes, name="spotlight dim", op=0.62):
    """void at op over everything except the holes (rect subtraction into bands)."""
    rects = [(x, y, x + w, y + h)]
    for (hx0, hy0, hx1, hy1) in holes:
        out = []
        for (a, b, c, d) in rects:
            if hx1 <= a or hx0 >= c or hy1 <= b or hy0 >= d:
                out.append((a, b, c, d)); continue
            if hy0 > b: out.append((a, b, c, hy0))
            if hy1 < d: out.append((a, hy1, c, d))
            mb, md = max(b, hy0), min(d, hy1)
            if hx0 > a: out.append((a, mb, hx0, md))
            if hx1 < c: out.append((hx1, mb, c, md))
        rects = out
    body = "".join(f'<rect x="{a}" y="{b}" width="{c-a}" height="{d-b}" fill="{C["void"]}" fill-opacity="{op}"/>' for (a, b, c, d) in rects)
    ring = [frame_lines(hx0 - 1, hy0 - 1, hx1 - hx0 + 2, hy1 - hy0 + 2, "amber", "spotlight ring (amber)") for (hx0, hy0, hx1, hy1) in holes]
    return G(name, f'<g id="dim (void 62%)">{body}</g>', *ring)

def pixel_circle(cx, cy, r, color, name):
    pts = set()
    x, y, d = r, 0, 1 - r
    while x >= y:
        for a, b in ((x, y), (y, x), (-x, y), (-y, x), (x, -y), (y, -x), (-x, -y), (-y, -x)):
            pts.add((cx + a, cy + b))
        y += 1
        if d < 0: d += 2 * y + 1
        else: x -= 1; d += 2 * (y - x) + 1
    return f'<path id="{nm(name)}" fill="{C.get(color, color)}" d="{"".join(f"M{a} {b}h1v1h-1z" for a, b in sorted(pts))}"/>'

def pixel_poly(pts, color, name):
    """1 px line through UI points (for the pressure curve)."""
    px = set()
    for (x0, y0), (x1, y1) in zip(pts, pts[1:]):
        n = int(max(abs(x1 - x0), abs(y1 - y0))) + 1
        for i in range(n + 1):
            t = i / n
            px.add((int(round(x0 + (x1 - x0) * t)), int(round(y0 + (y1 - y0) * t))))
    return f'<path id="{nm(name)}" fill="{C.get(color, color)}" d="{"".join(f"M{a} {b}h1v1h-1z" for a, b in sorted(px))}"/>'

GID = [0]
def gid(p):
    GID[0] += 1; return f"{p}{GID[0]}"

def hexrgb(h): return tuple(int(h[i:i + 2], 16) for i in (1, 3, 5))
def rgbhex(c): return "#" + "".join(f"{max(0, min(255, int(round(v)))):02x}" for v in c)
def mix(a, b, t): return rgbhex([x + (y - x) * t for x, y in zip(hexrgb(a), hexrgb(b))])

def grad_at(stops, t):
    """stops: [(pos, color, hard)]. A hard stop has no blend into it: the segment before it keeps the previous color."""
    if t <= stops[0][0]: return stops[0][1][:7]
    for (p0, c0, _), (p1, c1, h1) in zip(stops, stops[1:]):
        if t < p1:
            return c0[:7] if h1 else mix(c0[:7], c1[:7], (t - p0) / (p1 - p0))
    return stops[-1][1][:7]

def svg_stops(stops):
    out = []
    for i, (p, c, h) in enumerate(stops):
        if h and i > 0:
            out.append((p, stops[i - 1][1]))
        out.append((p, c))
    def one(p, c):
        op = int(c[7:9], 16) / 255 if len(c) > 7 else 1
        return f'<stop offset="{p:.3f}" stop-color="{c[:7]}" stop-opacity="{op:.3f}"/>'
    return "".join(one(p, c) for p, c in out)

def even(cols):
    n = len(cols)
    return [(i / (n - 1), c.lstrip("!"), c.startswith("!")) for i, c in enumerate(cols)]

# ---------- the drawing (drawing coordinates, vector) ----------
def layer_bodies(sky=("solid", "#f2a65a"), glow=("solid", "#f2a65a", 0.4), core="#fff3c4", sea="#2d6a8f", extra_sun="", wave="#a9d3e4"):
    hz = 200
    b = {}
    if sky[0] == "solid":
        b["sky"] = f'<rect id="rect line 26" x="0" y="0" width="480" height="{hz}" fill="{sky[1]}"/>'
    else:
        _, (x1, y1), (x2, y2), stops, steps = sky
        if steps:
            parts = [f'<rect x="0" y="0" width="480" height="{y1}" fill="{grad_at(stops, 0)}"/>'] if y1 > 0 else []
            bh = (y2 - y1) / steps
            for i in range(steps):
                parts.append(f'<rect id="band {i+1}" x="0" y="{y1 + i*bh:.2f}" width="480" height="{bh + 0.5:.2f}" fill="{grad_at(stops, (i + 0.5) / steps)}"/>')
            b["sky"] = f'<g id="rect line 26 (stepped linear, {steps} bands)">{"".join(parts)}</g>'
        else:
            g = gid("sky")
            b["sky"] = (f'<g id="rect line 26 (linear gradient)"><defs><linearGradient id="{g}" gradientUnits="userSpaceOnUse" x1="{x1}" y1="{y1}" x2="{x2}" y2="{y2}">'
                        f'{svg_stops(stops)}</linearGradient></defs><rect x="0" y="0" width="480" height="{hz}" fill="url(#{g})"/></g>')
    waves = []
    for i in range(5):
        pts = [(40 + i*20, hz + 30 + i*18), (140 + i*30, hz + 24 + i*18), (240 + i*20, hz + 32 + i*18), (440 - i*10, hz + 26 + i*18)]
        waves.append(f'<path d="{bezier_d(pts)}" fill="none" stroke="{wave}" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>')
    b["sea"] = f'<rect x="0" y="{hz}" width="480" height="120" fill="{sea}"/><g id="layer waves">{"".join(waves)}</g>'
    b["boat"] = ('<polygon points="120,214 180,214 165,232 135,232" fill="#3b2a6b"/>'
                 '<line x1="150" y1="214" x2="150" y2="170" stroke="#3b2a6b" stroke-width="3" stroke-linecap="round"/>'
                 '<polygon points="152,174 152,208 182,208" fill="#fbf8ef"/>')
    if glow[0] == "solid":
        gl = f'<circle id="glow line 11" cx="360" cy="90" r="56" fill="{glow[1]}" fill-opacity="{glow[2]}"/>'
    else:
        _, (cx, cy), r, stops = glow
        g = gid("glow")
        gl = (f'<g id="glow line 11 (radial gradient)"><defs><radialGradient id="{g}" gradientUnits="userSpaceOnUse" cx="{cx}" cy="{cy}" r="{r}">{svg_stops(stops)}</radialGradient></defs>'
              f'<circle cx="{cx}" cy="{cy}" r="{r}" fill="url(#{g})"/></g>')
    b["sun"] = gl + (f'<circle id="core line 10" cx="360" cy="90" r="34" fill="{core}"/>' if core else "") + extra_sun
    b["bg"] = '<rect x="0" y="0" width="480" height="320" fill="#f6d9a8"/>'
    return b

SHADOWS = ('<polygon id="polygon line 13 (dusk shade)" points="0,120 150,96 230,320 0,320" fill="#8a6fb0"/>'
           '<ellipse id="ellipse line 14 (boat shadow)" cx="152" cy="240" rx="44" ry="7" fill="#3b2a6b"/>')

def drawing_g(lay, bodies, order=("bg", "sky", "sea", "boat", "sun"), name="canvas drawing (vector, full-res)", extra=""):
    body = "".join(bodies.get(k, "") for k in order) + extra
    return f'<g id="{nm(name + " " + lay.pct)}" data-size="480 320" transform="translate({lay.dx} {lay.dy}) scale({lay.z})">{body}</g>'

def vec(name, x, y, w, h, body, scale=1.0):
    return f'<g id="{nm(name)}" data-size="{w} {h}" transform="translate({x} {y}) scale({scale})">{body}</g>'

def px_grad(name, x, y, w, h, at, two_d=False, check=None):
    """Smooth color content as exact screen pixels (0.5 UI px rects): at(u, v) -> (hex, alpha).
    Alpha composites over a checker (UI px size `check`) of #d4dacd / #a0ab9b."""
    out = []
    nx, ny = int(w * S), int(h * S)
    def comp(c, a, sx, sy):
        if a >= 1 or not check: return c[:7]
        k = ((int(sx / S) // check) + (int(sy / S) // check)) % 2
        return mix("#a0ab9b" if k else "#d4dacd", c[:7], a)
    if not two_d and not check:
        for i in range(nx):
            c, a = at((i + 0.5) / nx, 0.5)
            out.append(f'<rect x="{x + i / S:g}" y="{y}" width="0.5" height="{h}" fill="{c[:7]}"/>')
    else:
        for j in range(ny):
            run_c, run_s = None, 0
            for i in range(nx + 1):
                c = None
                if i < nx:
                    cc, a = at((i + 0.5) / nx, (j + 0.5) / ny)
                    c = comp(cc, a, i, j)
                if c != run_c:
                    if run_c is not None:
                        out.append(f'<rect x="{x + run_s / S:g}" y="{y + j / S:g}" width="{(i - run_s) / S:g}" height="0.5" fill="{run_c}"/>')
                    run_c, run_s = c, i
    return f'<g id="{nm(name)}">{"".join(out)}</g>'

def stops_at(stops):
    def f(u, v):
        c = grad_at_a(stops, u)
        return c
    return f

def grad_at_a(stops, t):
    """Like grad_at but keeps alpha: returns (hex, alpha)."""
    def ca(c): return (c[:7], int(c[7:9], 16) / 255 if len(c) > 7 else 1.0)
    if t <= stops[0][0]: return ca(stops[0][1])
    for (p0, c0, _), (p1, c1, h1) in zip(stops, stops[1:]):
        if t < p1:
            if h1: return ca(c0)
            f = (t - p0) / (p1 - p0)
            (a0, al0), (a1, al1) = ca(c0), ca(c1)
            return (mix(a0, a1, f), al0 + (al1 - al0) * f)
    return ca(stops[-1][1])

# ---------- window pieces ----------
def tabbar2(name="layers.lzdraw", unsaved=False):
    w = tw(name) + 26
    items = [R(0, TOP_H, W, TAB_H, "coal", "tab bar bg"), R(0, MAIN_Y - 1, W, 1, "void", "tab bar bottom line")]
    tab = [R(2, TOP_H + 1, w, TAB_H, "panel", "tab face"), R(2, TOP_H + 1, w, 1, "lizard", "tab active top line"),
           R(1, TOP_H + 1, 1, TAB_H - 1, "void", "tab left"), R(2 + w, TOP_H + 1, 1, TAB_H - 1, "void", "tab right"),
           T(name, 7, TOP_H + 1), I("dot" if unsaved else "close", 7 + tw(name) + 4, TOP_H + 4, "amber" if unsaved else "mute",
                                    "unsaved dot" if unsaved else "close tab")]
    items.append(G(f"tab {name} (active)", *tab))
    items.append(G("tab new", T("+", 2 + w + 6, TOP_H + 1, "mute")))
    return G("tab bar", *items)

def statusbar_y(y, parts_left, parts_right, w=W, name="status bar"):
    items = [R(0, y, w, 16, "panel", "status bg"), R(0, y, w, 1, "void", "status top line")]
    x = 6
    for i, (ic, s, col) in enumerate(parts_left):
        if ic:
            items.append(I(ic, x, y + 4, col)); x += 10
        items.append(T(s, x, y + 2, col, f"status: {s}")); x += tw(s) + 6
        if i < len(parts_left) - 1:
            items.append(R(x, y + 4, 1, 8, "bevel", "separator"))
        x += 7
    x = w - 6
    for i, (ic, s, col) in enumerate(reversed(parts_right)):
        x -= tw(s)
        items.append(T(s, x, y + 2, col, f"status: {s}"))
        if ic:
            x -= 10; items.append(I(ic, x, y + 4, col))
        x -= 7
        if i < len(parts_right) - 1:
            items.append(R(x, y + 4, 1, 8, "bevel", "separator"))
        x -= 6
    return G(name, *items)

def counter_part(lang, n, mx):
    st = "over" if n > mx else ("near" if n >= 0.8 * mx else "normal")
    col = {"normal": "mute", "near": "amber", "over": "ember"}[st]
    ic = {"normal": "objects", "near": "objects", "over": "warn"}[st]
    return (ic, f"{TX[lang]['objects']} {fmt(n, lang)} / {fmt(mx, lang)}", col), st

def status2(lang, lay=None, saved=True, objects=(14, 1800), file="layers.lzdraw", pos="x 352 y 88", zoom=None):
    t = TX[lang]
    left = [("canvas", file, "text"), ("check", t["saved"], "lizard") if saved else ("dot", t["unsaved"], "amber"),
            counter_part(lang, *objects)[0]]
    right = ([(None, pos, "mute")] if pos else []) + [(None, zoom or (lay.pct if lay else "75%"), "mute"), ("language", t["lang"], "text")]
    return statusbar_y(STATUS_Y, left, right)

def row_heights(rows):
    y, out = MAIN_Y + 15, {}
    for r in rows:
        if r["kind"] == "layer": y += LAYER_GAP
        out[r.get("line")] = y
        y += 4 if r["kind"] == "blank" else (10 if r["kind"] == "add" else ROW_H)
    return out

CANVAS_CODE = 'canvas = {"size": (480, 320), "background": #f6d9a8}'

def code_panel(lay, lang, rows, sel_line=None, hover_lines=(), props=None, file="layers.lzdraw", canvas_code=CANVAS_CODE, extra_after=None):
    t = TX[lang]
    w = lay.cw
    items = [R(0, MAIN_Y, w, STATUS_Y - MAIN_Y, "panel", "panel bg"),
             G("panel header", R(0, MAIN_Y, w, 13, "raised", "header bg"), R(0, MAIN_Y + 13, w, 1, "void", "header line"),
               T(t["layers"], 6, MAIN_Y, "mute"), T(file, w - tw(file) - 6, MAIN_Y, "dim"))]
    y = MAIN_Y + 15
    rows_g, hidden = [], False
    for r in rows:
        if r["kind"] == "layer": y += LAYER_GAP
        if y + (ROW_H if r["kind"] not in ("blank", "add") else 4) > PROPS_Y - ROW_H - 2:
            hidden = True; break
        el, h = gen.row_el(r, y, w=w, sel=(sel_line is not None and r.get("line") == sel_line),
                           hover=(r.get("line") in hover_lines), caret=r.get("caret"))
        rows_g.append(el); y += h
        if extra_after and r.get("line") == extra_after[0]:
            el2, h2 = extra_after[1](y); rows_g.append(el2); y += h2
    items.append(G("rows", *rows_g))
    if hidden:
        top, bot = MAIN_Y + 15, PROPS_Y - ROW_H - 3
        items.append(G("scrollbar (more rows below)", R(w - 5, top, 4, bot - top, "coal", "track"),
                       R(w - 5, top, 4, (bot - top) * 2 // 3, "bevel", "thumb")))
    py = PROPS_Y - ROW_H - 1
    el, _ = gen.row_el(dict(kind="canvas", d=0, code=canvas_code, line=5), py, w=w)
    items.append(G("pinned canvas row (line 5)", R(0, py - 1, w, 1, "void", "pin line"), R(0, py, w, ROW_H, "coal", "pinned bg"), el))
    items.append(props or "")
    return G("layers panel (the code)", *items)

def props2(lay, lang, title, fields_fn, hint=None):
    t = TX[lang]
    w, y = lay.cw, PROPS_Y
    items = [R(0, y, w, STATUS_Y - y, "panel", "properties bg"), R(0, y, w, 1, "void", "top line"),
             R(0, y + 1, w, 15, "raised", "header bg"), R(0, y + 16, w, 1, "void", "header line"),
             T(t["props"] + " · " + title, 6, y + 2, "mute")]
    items.append(fields_fn(y + 22))
    if hint:
        items.append(T(hint, 6, y + 45, "mute", "hint: " + hint))
    return G("properties panel", *items)

def canvas_area2(lay, drawing_svg, overlay="", after=(), splitter=True):
    items = []
    if splitter:
        items += [R(lay.cw, MAIN_Y, 2, STATUS_Y - MAIN_Y, "void", "splitter"), R(lay.cw + 1, MAIN_Y, 1, STATUS_Y - MAIN_Y, "bevel", "splitter highlight")]
    items.append(checker(lay.cx0, lay.cy0, lay.cx1 - lay.cx0, lay.cy1 - lay.cy0, "coal", "#161815", 16, "canvas surround (checker)"))
    items.append(frame_lines(lay.dx - 1, lay.dy - 1, lay.dw + 2, lay.dh + 2, "void", "canvas frame"))
    items.append(drawing_svg)
    items += list(after)
    items.append(overlay)
    return G("canvas area", *items)

# ---------- file rows ----------
def rows_layers(top_cut=True, sky_code="rect((0, 0), (480, horizon), #f2a65a)", glow_code="circle((360, 90), 56, #f2a65a66)",
                sun_open=True, sea_open=False, extra_top=()):
    rows = []
    if not top_cut:
        rows += [dict(kind="comment", d=0, code="-~> lzdraw 1", line=1), dict(kind="blank", line=4)]
    rows += list(extra_top)
    rows += [dict(kind="value", d=0, code="let horizon = 200", line=7), dict(kind="blank", line=8),
             dict(kind="layer", d=0, code="layer sun", line=9, open=sun_open)]
    if sun_open:
        rows += [dict(kind="circle", d=1, code="circle((360, 90), 34, #fff3c4)", line=10), dict(kind="circle", d=1, code=glow_code, line=11)]
    else:
        rows[-1]["count"] = 2
    rows += [dict(kind="blank", line=12), dict(kind="layer", d=0, code="layer boat", line=13, open=False, count=3), dict(kind="blank", line=17),
             dict(kind="layer", d=0, code="layer sea", line=18, open=sea_open, **({} if sea_open else {"count": 6})), dict(kind="blank", line=24),
             dict(kind="layer", d=0, code="layer sky", line=25, open=True), dict(kind="rect", d=1, code=sky_code, line=26)]
    return rows

# ---------- palette sidebar ----------
def hsv_of(hexc):
    r, g, b = [v / 255 for v in hexrgb(hexc)]
    return colorsys.rgb_to_hsv(r, g, b)

def sv_square(x, y, w, h, hexc, name="saturation/value square (vector)"):
    hh, s, v = hsv_of(hexc)
    base = rgbhex([c * 255 for c in colorsys.hsv_to_rgb(hh, 1, 1)])
    def at(u, v):
        return (rgbhex([c * 255 for c in colorsys.hsv_to_rgb(hh, u, 1 - v)]), 1)
    mx, my = x + int(s * (w - 1)), y + int((1 - v) * (h - 1))
    marker = G("picker marker", frame_lines(mx - 3, my - 3, 7, 7, "void", "marker outer"), frame_lines(mx - 2, my - 2, 5, 5, "bright", "marker ring"))
    return [frame_lines(x - 1, y - 1, w + 2, h + 2, "void", "square frame"), px_grad(name.replace("(vector)", "(pixels)"), x, y, w, h, at, two_d=True), marker]

def hue_strip(x, y, w, h, hexc, vertical=False):
    hh = hsv_of(hexc)[0]
    def at(u, v):
        t = v if vertical else u
        return (rgbhex([c * 255 for c in colorsys.hsv_to_rgb(t % 1, 1, 1)]), 1)
    if vertical:
        ty = y + int(hh * (h - 1))
        th = G("hue thumb", frame_lines(x - 2, ty - 2, w + 4, 5, "void", "thumb outline"), frame_lines(x - 1, ty - 1, w + 2, 3, "bright", "thumb"))
    else:
        tx = x + int(hh * (w - 1))
        th = G("hue thumb", frame_lines(tx - 2, y - 2, 5, h + 4, "void", "thumb outline"), frame_lines(tx - 1, y - 1, 3, h + 2, "bright", "thumb"))
    return [frame_lines(x - 1, y - 1, w + 2, h + 2, "void", "hue frame"), px_grad("hue strip", x, y, w, h, at, two_d=vertical), th]

def alpha_strip(x, y, w, h, hexc, a=1.0, vertical=False):
    def at(u, v):
        return (hexc[:7], (1 - v) if vertical else u)
    if vertical:
        ty = y + int((1 - a) * (h - 1))
        th = G("alpha thumb", frame_lines(x - 2, ty - 2, w + 4, 5, "void", "thumb outline"), frame_lines(x - 1, ty - 1, w + 2, 3, "bright", "thumb"))
    else:
        tx = x + int(a * (w - 1))
        th = G("alpha thumb", frame_lines(tx - 2, y - 2, 5, h + 4, "void", "thumb outline"), frame_lines(tx - 1, y - 1, 3, h + 2, "bright", "thumb"))
    return [frame_lines(x - 1, y - 1, w + 2, h + 2, "void", "alpha frame"), px_grad("alpha strip (over checker)", x, y, w, h, at, check=2), th]

PALETTE = [("sky", "#f2a65a"), ("sea", "#2d6a8f"), ("sun", "#fff3c4"), ("boat", "#3b2a6b")]
RECENT = ["#e0607a", "#3b2a6b", "#f2a65a", "#a9d3e4", "#fbf8ef", "#4fc46a", "#1d1d1f", "#8a6fb0"]
SAVED_GRADS = [("dusk", ["#3b2a6b", "#e0607a", "#f2a65a"]), ("sea glass", ["#a9d3e4", "#2d6a8f", "#12304a"]), ("ember", ["#fff3c4", "#f36c5b", "#3b2a6b"])]

def grad_bar_vec(name, x, y, w, h, cols):
    st = even(cols)
    return px_grad(name.replace("(vector)", "(pixels)"), x, y, w, h, lambda u, v: grad_at_a(st, u))

# color variables: solid or gradient, interchangeable (`colors.dusk` works wherever a color does)
VARS = [("sky", "s", "#f2a65a"), ("sea", "s", "#2d6a8f"), ("sun", "s", "#fff3c4"), ("boat", "s", "#3b2a6b"),
        ("dusk", "l", ["#3b2a6b", "#e0607a", "#f2a65a"]), ("glow", "r", ["#fff3c4", "#f2a65a00"]), ("dune", "st", ["#f6d9a8", "#c0703a"])]

def var_chip(x, y, w, h, var, name_col="text"):
    n, k, v = var
    it = [frame_lines(x, y, w, h, "void", "chip outline")]
    ix, iy, iw, ih = x + 1, y + 1, w - 2, h - 2
    if k == "s":
        it.append(R(ix, iy, iw, ih, v[:7], v))
    elif k == "l":
        st = even(v)
        it.append(px_grad(f"linear chip", ix, iy, iw, ih, lambda u, vv: grad_at_a(st, u)))
    elif k == "r":
        st = even(v)
        it.append(px_grad(f"radial chip", ix, iy, iw, ih, lambda u, vv: grad_at_a(st, min(1, math.hypot((u - 0.5) * iw, (vv - 0.5) * ih) / (iw / 2))), check=2))
    else:
        n_ = 4
        for i in range(n_):
            x0 = ix + iw * i // n_; x1 = ix + iw * (i + 1) // n_
            it.append(R(x0, iy, x1 - x0, ih, mix(v[0], v[1], i / (n_ - 1)), f"band {i+1}"))
    return G(f"variable chip {n} ({ {'s': 'solid', 'l': 'linear', 'r': 'radial', 'st': 'stepped'}[k]})", *it)

def var_grid(x, y, cols, cw, vars_, sel=None, focus=None, chip_h=12, name_col="text"):
    """Chips with the name under each. Cell cw wide, chip_h + 14 high."""
    it = []
    ch = chip_h + 15
    for i, var in enumerate(vars_):
        cx, cy = x + (i % cols) * cw, y + (i // cols) * ch
        g = []
        if var[0] == sel:
            g.append(R(cx - 1, cy - 1, cw - 1, ch - 1, "moss", "cell bg (selected)"))
        g.append(var_chip(cx, cy, cw - 3, chip_h, var))
        g.append(T(var[0], cx, cy + chip_h, "bright" if var[0] == sel else name_col, "name " + var[0]))
        if var[0] == focus:
            g.append(frame_lines(cx - 2, cy - 2, cw + 1, ch + 1, "amber", "keyboard focus ring (amber)"))
        it.append(G(f"variable colors.{var[0]}" + (" (selected)" if var[0] == sel else "") + (" (keyboard focus)" if var[0] == focus else ""), *g))
    rows = (len(vars_) + cols - 1) // cols
    return G("variables (solid and gradient)", *it), rows * ch

def sidebar(lay, lang, color, title=None, vars_=VARS, sel_name=None, used=None, alpha="100%", hover_add=False, hex_focus=False, var_kind=None):
    t = TX[lang]
    x0, w = lay.sb_x, lay.sb_w
    ix, iw = x0 + 5, w - 10
    items = [R(x0 - 2, MAIN_Y, 2, STATUS_Y - MAIN_Y, "void", "splitter"), R(x0 - 1, MAIN_Y, 1, STATUS_Y - MAIN_Y, "bevel", "splitter highlight"),
             R(x0, MAIN_Y, w, STATUS_Y - MAIN_Y, "panel", "sidebar bg"),
             G("sidebar header", R(x0, MAIN_Y, w, 13, "raised", "header bg"), R(x0, MAIN_Y + 13, w, 1, "void", "header line"),
               T(title or t["color"], x0 + 5, MAIN_Y, "mute", "title"), I("chevron-right", x0 + w - 11, MAIN_Y + 3, "mute", "collapse to rail"))]
    y = MAIN_Y + 17
    if var_kind is not None:
        sg, _ = seg(ix, y, ["", "", ""], var_kind, "variable kind (solid / linear / radial)", icons=["solid", "gradient", "radial"], w=(iw + 2) // 3)
        items.append(sg)
        y += 18
    sq = 44 if var_kind is not None else 52
    items += sv_square(ix, y, iw, sq, color); y += sq + 5
    items += hue_strip(ix, y, iw, 6, color); y += 11
    items += alpha_strip(ix, y, iw, 6, color); y += 11
    hw = tw("#e0607a") + 8
    items.append(sfield(ix, y, hw, color[:7], "hex", "text", focused=hex_focus))
    items.append(sfield(ix + hw + 2, y, iw - hw - 2, alpha, "alpha", "amber"))
    y += 20
    items.append(T(t["recent"], ix, y - 1, "mute", "label recent"))
    y += 13
    items.append(G("recent colors", *[sw(ix + i * 12, y, 10, c) for i, c in enumerate(RECENT)]))
    y += 15
    items.append(G("label variables", T(t["palette"], ix, y - 1, "mute"), I("plus", ix + iw - 8, y + 2, "lizard", "add variable")))
    y += 13
    g, gh = var_grid(ix + 1, y, 3, 32, vars_, sel=sel_name)
    items.append(g)
    y += gh
    if used:
        items.append(T(used, ix, y, "lizard", "used-by hint"))
        y += 13
    by = STATUS_Y - 20
    items.append(button(ix, by, iw, 16, "button + save as variable" + (" (hover)" if hover_add else ""), icon="plus", label=t["add_pal"],
                        state="hover" if hover_add else "normal"))
    assert y < by, ("sidebar overflow", y, by)
    return G("palette sidebar (open)", *items)

def rail(lay):
    x0 = lay.sb_x
    items = [R(x0 - 2, MAIN_Y, 2, STATUS_Y - MAIN_Y, "void", "splitter"), R(x0 - 1, MAIN_Y, 1, STATUS_Y - MAIN_Y, "bevel", "splitter highlight"),
             R(x0, MAIN_Y, 14, STATUS_Y - MAIN_Y, "panel", "rail bg"), button(x0 + 1, MAIN_Y + 2, 12, 12, "button expand palette", icon="chevron-left")]
    y = MAIN_Y + 20
    for var in VARS:
        items.append(var_chip(x0 + 2, y, 10, 10, var)); y += 12
    items.append(R(x0 + 3, y + 1, 8, 1, "bevel", "separator")); y += 5
    for c in RECENT[:4]:
        items.append(sw(x0 + 2, y, 10, c, f"recent {c}")); y += 12
    return G("palette sidebar (collapsed rail)", *items)

# ---------- menus ----------
def menu_box(x, y, entries, name, title=None, hover=None, disabled=(), danger=("trash",), checks=(), w=None, swatches=None):
    """entries: (icon, label, shortcut) or None for a separator. Rows 14."""
    sc_w = max([tw(e[2]) for e in entries if e and e[2] not in ("", ">")] + [0])
    w = w or (max(tw(e[1]) for e in entries if e) + 18 + (sc_w + 12 if sc_w else 0) + 14)
    h = 4 + (14 if title else 0) + sum(14 if e else 4 for e in entries) + (16 if swatches else 0)
    items = chrome(x, y, w, h, name="menu")
    yy = y + 2
    if title:
        items.append(G("menu title", R(x + 1, y + 1, w - 2, 13, "raised", "title bg"), R(x + 1, y + 14, w - 2, 1, "void", "title line"),
                       T(title, x + 6, y + 1, "mute")))
        yy = y + 16
    for e in entries:
        if not e:
            items.append(R(x + 4, yy + 1, w - 8, 1, "bevel", "separator")); yy += 4; continue
        ic, lb, sc = e
        st = "hover" if lb == hover else ("disabled" if lb in disabled else "normal")
        dang = ic in danger
        it = []
        if st == "hover":
            it.append(R(x + 1, yy, w - 2, 14, "moss", "item bg (hover)"))
            it.append(R(x + 1, yy, 2, 14, "lizard", "hover bar"))
        tcol = {"hover": "bright", "disabled": "dim", "normal": "ember" if dang else "text"}[st]
        icol = {"hover": "lizard", "disabled": "dim", "normal": "ember" if dang else "mute"}[st]
        it += [I(ic, x + 6, yy + 3, icol), T(lb, x + 18, yy, tcol)]
        if sc == ">":
            it.append(I("chevron-right", x + w - 12, yy + 3, "bright" if st == "hover" else "mute", "submenu chevron"))
        elif sc:
            it.append(T(sc, x + w - 6 - tw(sc), yy, "dim" if st == "disabled" else ("text" if st == "hover" else "mute"), "shortcut " + sc))
        if lb in checks:
            it.append(I("check", x + w - 12, yy + 3, "lizard", "current"))
        items.append(G(f"menu item {lb}" + {"hover": " (hover)", "disabled": " (disabled)", "normal": ""}[st] + (" (danger)" if dang else ""), *it))
        yy += 14
    if swatches:
        items.append(G("palette swatches", *[sw(x + 6 + i * 12, yy + 3, 10, c, f"swatch {n}") for i, (n, c) in enumerate(swatches)]))
    return G(name, *items), w, h

# ---------- fill editor ----------
MARK = ["....o....", "...ooo...", "..ooooo..", "ooooooooo", "occccccco", "occccccco", "occccccco", "occccccco", "ooooooooo"]
MARK_HARD = ["....o....", "....o....", "....o....", "ooooooooo", "occccccco", "occccccco", "occccccco", "occccccco", "ooooooooo"]
def stop_marker(cx, top, col, sel=False, hard=False, ghost=False, name="stop"):
    pat = MARK_HARD if hard else MARK
    oc = "bright" if sel else ("dim" if ghost else "void")
    do, dc = [], []
    for r, row in enumerate(pat):
        for c, ch in enumerate(row):
            if ch == "o": do.append(f"M{cx - 4 + c} {top + r}h1v1h-1z")
            elif ch == "c": dc.append(f"M{cx - 4 + c} {top + r}h1v1h-1z")
    it = []
    if sel:
        it.append(frame_lines(cx - 5, top + 2, 11, 8, "void", "selected outer outline"))
    it.append(f'<path id="marker outline" fill="{C[oc]}" d="{"".join(do)}"/>')
    it.append(f'<path id="marker color" fill="{C["row"] if ghost else col[:7]}" d="{"".join(dc)}"/>')
    if ghost:
        it.append(T("+", cx - 3, top + 9, "dim", "plus"))
    return G(name + (" (selected)" if sel else "") + (" (hard)" if hard else "") + (" (ghost: click to add)" if ghost else ""), *it)

FE_W = 186
def fill_editor(x, y, lang, kind, stops, sel=1, steps=0, code_lines=(), title="", tail=None, ghost_at=None, at=None):
    """kind 0 solid, 1 linear, 2 radial. stops: [(pos, color, hard)]."""
    t = TX[lang]
    w = FE_W
    h = 186
    items = chrome(x, y, w, h, name="fill editor")
    items.append(header(x, y, w, f'{t["fill"]} · {title}'))
    tabs, _ = seg(x + 6, y + 18, t["tabs"], kind, "tabs fill kind", icons=["solid", "gradient", "radial"], w=(w - 12 + 2) // 3)
    items.append(tabs)
    bx, by, bw, bh = x + 6, y + 38, w - 12, 14
    bar = [frame_lines(bx, by, bw, bh, "void", "bar outline")]
    if steps:
        n = steps
        for i in range(n):
            x0 = bx + 1 + (bw - 2) * i // n; x1 = bx + 1 + (bw - 2) * (i + 1) // n
            bar.append(R(x0, by + 1, x1 - x0, bh - 2, grad_at(stops, (i + 0.5) / n), f"band {i+1}"))
    else:
        alpha = any(len(c) > 7 for _, c, _ in stops)
        bar.append(px_grad("gradient bar fill", bx + 1, by + 1, bw - 2, bh - 2, lambda u, v: grad_at_a(stops, u), check=3 if alpha else None))
    for i, (p, c, hd) in enumerate(stops):
        if hd:
            bar.append(R(bx + 1 + int((bw - 3) * p), by + 1, 1, bh - 2, "bright", f"hard stop line {i+1}"))
    items.append(G("gradient bar" + (f" (stepped, {steps} bands)" if steps else ""), *bar))
    for i, (p, c, hd) in enumerate(stops):
        items.append(stop_marker(bx + 1 + int((bw - 3) * p), by + bh + 1, c, sel=(i == sel), hard=hd, name=f"color stop {i+1} {c}"))
    if ghost_at is not None:
        items.append(stop_marker(bx + 1 + int((bw - 3) * ghost_at), by + bh + 1, "#000000", ghost=True, name="color stop"))
    y1 = y + 66
    sp, sc, shd = stops[sel]
    at = at or f"{int(round(sp * 100))}%"
    items.append(sfield(x + 6, y1, tw("#e0607a") + 22, sc[:7] if len(sc) <= 7 or sc.endswith("ff") else sc, "stop color", "text", swatch=sc))
    ax = x + 6 + tw("#e0607a") + 22 + 6
    items.append(G("field at", T(t["at"], ax, y1 + 1, "mute"), sfield(ax + tw(t["at"]) + 4, y1, tw(at) + 8, at, "stop position", "amber")))
    items.append(button(x + w - 20, y1 + 1, 14, 14, "button remove stop", icon="trash"))
    y2 = y1 + 20
    sg, ex = seg(x + 6, y2, [t["soft"], t["hard"]], 1 if shd else 0, "stop blend soft/hard")
    items.append(sg)
    stx = ex + 7
    sv = str(steps) if steps else t["off"]
    items.append(G("field steps", I("steps", stx, y2 + 3, "mute"), T(t["steps"], stx + 10, y2, "mute"),
                   sfield(stx + 12 + tw(t["steps"]), y2 - 1, tw(sv) + 8, sv, "steps", "amber" if steps else "mute"),
                   button(stx + 12 + tw(t["steps"]) + tw(sv) + 9, y2, 12, 14, "button steps minus", icon="minus"),
                   button(stx + 12 + tw(t["steps"]) + tw(sv) + 20, y2, 12, 14, "button steps plus", icon="plus")))
    y3 = y2 + 20
    ch = len(code_lines) * 13 + 6
    items.append(G("code preview (what it writes)", frame_lines(x + 6, y3, w - 12, ch, "void", "inset outline"), R(x + 7, y3 + 1, w - 14, ch - 2, "coal", "inset"),
                   *[code(ln, x + 10, y3 + 3 + i * 13, name=f"code line {i+1}") for i, ln in enumerate(code_lines)]))
    y4 = y3 + ch + 3
    items.append(T(t["hint1"], x + 6, y4, "mute", "hint add"))
    items.append(T(t["hint2"], x + 6, y4 + 12, "mute", "hint remove"))
    items.append(G("button + save as variable", I("plus", x + 6, y4 + 29, "lizard"), T(CP[lang]["save"], x + 16, y4 + 26, "lizard")))
    hh = y4 + 26 + 16 - y
    items = chrome(x, y, w, hh, name="fill editor") + items[4:]
    if tail:
        side, ty = tail
        tl = []
        for k in range(6):
            span = 11 - 2 * k
            if span <= 0: break
            if side == "right":
                cx = x + w - 1 + k
            else:
                cx = x - k
            top = ty - span // 2
            tl.append(R(cx, top, 1, span, "void", "tail edge"))
            if span > 2:
                tl.append(R(cx, top + 1, 1, span - 2, "row", "tail fill"))
        items.append(G("tail (points at the shape)", *tl))
    return G(f"fill editor ({['solid', 'linear', 'radial'][kind]}" + (f", {steps} steps" if steps else "") + ")", *items), hh

def linear_handles(lay, a, b, stops, name="gradient handles (linear)"):
    ax, ay = [int(v) for v in lay.cv(a)]; bx, by = [int(v) for v in lay.cv(b)]
    it = [R(ax - 1, ay, 3, by - ay, "void", "line outline"), R(ax, ay, 1, by - ay, "bright", "gradient line")]
    for i, (p, c, hd) in enumerate(stops):
        yy = int(ay + (by - ay) * p)
        it.append(G(f"stop dot {i+1}", R(yy and ax - 3, yy - 3, 7, 7, "void", "outline"), R(ax - 2, yy - 2, 5, 5, c[:7], c)))
    it.append(gen.anchor(ax * S, ay * S, "handle start (from)", "alt"))
    it.append(gen.anchor(bx * S, by * S, "handle end (to)", "normal"))
    return G(name, *it)

# ---------- shared frame bits ----------
def window2(parts): return svg_doc(W * S, H * S, "".join(parts), bg="coal", scale=S)

def circle_fields(lang, x0=6, fill="#fff3c4"):
    def f(y):
        g1, x = gen.point_field(x0, y, "center", "360", "90", "center")
        g2 = gen.num_field(x + 4, y, "r", "34", "r")
        x2 = x + 4 + tw("r") + 4 + tw("34") + 8 + 8
        g3, _ = gen.color_field(x2, y, "fill", fill, "fill")
        return G("fields (from circle(center, r, fill))", g1, g2, g3)
    return f

def circ_cr(r, cx="360", cy="90"):
    def f(y):
        g1, x = gen.point_field(6, y, "center", cx, cy, "center")
        return G("fields (center, r; fill wraps)", g1, gen.num_field(x + 4, y, "r", r, "r"))
    return f

# ======================= frames =======================
def f_menu():
    lang = "en"; t = TX[lang]
    lay = Lay("rail")
    rows = rows_layers()
    props = props2(lay, lang, f"circle · {t['line']} 10", circ_cr("34"), "line 10 · right-click: more")
    sx, sy = [int(v) for v in lay.cv((360, 90))]
    x0, y0 = lay.cvs((326, 56)); x1, y1 = lay.cvs((394, 124))
    ov = [sel_box(x0, y0, x1, y1), gen.anchor(*lay.cvs((360, 90)), "anchor center", "normal"), gen.anchor(*lay.cvs((394, 90)), "anchor radius", "alt")]
    m, mw, mh = menu_box(0, 0, t["menu"], "x", title="circle · line 10", hover="Fill", disabled=("Bring forward",))
    mx, my = sx - mw - 2, sy + 4
    m, mw, mh = menu_box(mx, my, t["menu"], "context menu on shape (canvas right-click)", title="circle · line 10", hover="Fill", disabled=("Bring forward",))
    fy = my + 16 + 4 * 14 + 4 + 2 * 14 + 4 + 2 - 2
    sub, sw_, sh_ = menu_box(0, 0, t["fill_sub"], "x", swatches=PALETTE + [("recent", c) for c in RECENT[:3]])
    sub, sw_, sh_ = menu_box(mx - sw_ + 2, fy, t["fill_sub"], "submenu Fill ▸", checks=("Solid color",), hover="Linear gradient", swatches=PALETTE + [("recent", c) for c in RECENT[:3]])
    return window2([gen.toolbar(None, lang), tabbar2(), code_panel(lay, lang, rows, sel_line=10, props=props),
                    canvas_area2(lay, drawing_g(lay, layer_bodies()), G("selection (circle line 10)", *ov)), rail(lay),
                    m, sub, cursor(mx - 50, fy + 20), status2(lang, lay, pos="x 360 y 90")])

SKY_STOPS = even(["#3b2a6b", "#e0607a", "#f2a65a"])
def f_fill_linear():
    lang = "en"; t = TX[lang]
    lay = Lay("open")
    sky_code = "rect((0, 0), (480, horizon), gradient([#3b2a6b, #e0607a, #f2a65a], 90deg))"
    rows = rows_layers(sky_code=sky_code, sun_open=False)
    props = props2(lay, lang, f"rect · {t['line']} 26", lambda y: G("fields", gen.point_field(6, y, "pos", "0", "0", "pos")[0],
                   G("field fill (gradient)", T("fill", 104, y + 2, "mute"), field(104 + 32, y, 116, 18, [grad_bar_vec("fill swatch gradient (vector)", 104 + 35, y + 4, 24, 10, ["#3b2a6b", "#e0607a", "#f2a65a"]),
                     T("gradient(…)", 104 + 63, y + 3, "text")], "input fill", focused=True))), "line 26 · drag the handles on the canvas")
    bodies = layer_bodies(sky=("linear", (240, 0), (240, 200), SKY_STOPS, 0))
    hnd = linear_handles(lay, (240, 0), (240, 200), SKY_STOPS)
    ax, ay = lay.cv((0, 0))
    ov = G("selection (rect line 26)", sel_box(*lay.cvs((0, 0)), *lay.cvs((480, 200))), hnd)
    lines = ["gradient([#3b2a6b,", "  #e0607a, #f2a65a],", "  90deg)"]
    fe, fh = fill_editor(lay.cx0 - FE_W - 4, MAIN_Y + 22, lang, 1, SKY_STOPS, sel=1, code_lines=lines, title="rect · line 26",
                         tail=("right", int(lay.cv((0, 60))[1])), ghost_at=0.75)
    gx = lay.cx0 - FE_W - 4 + 6 + 1 + int((FE_W - 15) * 0.75)
    return window2([gen.toolbar(None, lang), tabbar2(unsaved=True), code_panel(lay, lang, rows, sel_line=26, props=props),
                    canvas_area2(lay, drawing_g(lay, bodies), ov), sidebar(lay, lang, "#e0607a", title="STOP 2"),
                    fe, cursor(gx + 1, MAIN_Y + 22 + 40), status2(lang, lay, saved=False, pos=None)])

def f_fill_stepped():
    lang = "en"; t = TX[lang]
    lay = Lay("open")
    sky_code = "rect((0, 0), (480, horizon), gradient([#3b2a6b, #e0607a, #f2a65a], 90deg, 5))"
    rows = rows_layers(sky_code=sky_code, sun_open=False)
    props = props2(lay, lang, f"rect · {t['line']} 26", lambda y: G("fields", gen.point_field(6, y, "pos", "0", "0", "pos")[0],
                   gen.num_field(104, y, "steps", "5", "steps", focused=True)), "line 26 · 5 flat bands, no blending")
    bodies = layer_bodies(sky=("linear", (240, 0), (240, 200), SKY_STOPS, 5))
    hnd = linear_handles(lay, (240, 0), (240, 200), SKY_STOPS)
    ov = G("selection (rect line 26)", sel_box(*lay.cvs((0, 0)), *lay.cvs((480, 200))), hnd)
    lines = ["gradient([#3b2a6b,", "  #e0607a, #f2a65a],", "  90deg, 5)"]
    fe, fh = fill_editor(lay.cx0 - FE_W - 4, MAIN_Y + 22, lang, 1, SKY_STOPS, sel=1, steps=5, code_lines=lines, title="rect · line 26",
                         tail=("right", int(lay.cv((0, 60))[1])))
    return window2([gen.toolbar(None, lang), tabbar2(unsaved=True), code_panel(lay, lang, rows, sel_line=26, props=props),
                    canvas_area2(lay, drawing_g(lay, bodies), ov), sidebar(lay, lang, "#e0607a", title="STOP 2"),
                    fe, status2(lang, lay, saved=False, pos=None)])

def f_palette():
    lang = "en"; t = TX[lang]
    lay = Lay("open")
    new_sky = "#e9806e"
    pal = [("sky", "s", new_sky)] + VARS[1:]
    NAMED["colors.sky"] = new_sky
    top = [dict(kind="value", d=0, code=f'let colors = {{"sky": {new_sky}, "sea": #2d6a8f, "sun": #fff3c4, "boat": #3b2a6b, "dusk": gradient([#3b2a6b, #e0607a, #f2a65a], 90deg)}}', line=6)]
    rows = rows_layers(sky_code="rect((0, 0), (480, horizon), colors.sky)", extra_top=top)
    for r in rows:
        if r.get("line") == 7: r["line"] = 7
    def f(y):
        g1, x = gen.color_field(6, y, "sky", new_sky, "sky")
        g2, x = gen.color_field(x, y, "sea", "#2d6a8f", "sea")
        return G("fields (hash literal keys)", g1, g2, T("+2", x, y + 2, "mute"))
    props = props2(lay, lang, f"colors · {t['line']} 6", f, "changing sky repaints its 2 shapes")
    bodies = layer_bodies(sky=("solid", new_sky), wave=new_sky)
    out = window2([gen.toolbar(None, lang), tabbar2(unsaved=True), code_panel(lay, lang, rows, sel_line=6, hover_lines=(18, 26), props=props),
                   canvas_area2(lay, drawing_g(lay, bodies)), sidebar(lay, lang, new_sky, title="colors.sky", vars_=pal, sel_name="sky", var_kind=0,
                   used="in 2 shapes", hover_add=False, hex_focus=True),
                   status2(lang, lay, saved=False, pos=None)])
    NAMED["colors.sky"] = "#f2a65a"
    return out

# ---------- artist mode ----------
ART_X1 = 360
RC_X = 362
def toolbar_artist(lang, code_on=False):
    t = TX[lang]
    items = [R(0, 0, W, TOP_H, "panel", "toolbar bg"), R(0, TOP_H - 1, W, 1, "void", "toolbar bottom line")]
    cw = tw(t["code_btn"]) + 8
    items.append(button(3, 1, cw, 14, "button </> Code (toggle code drawer)" + (" (active)" if code_on else ""), label=t["code_btn"],
                        state="active" if code_on else "normal"))
    items.append(R(3 + cw + 4, 3, 1, 10, "bevel", "separator"))
    items.append(T("artist mode" if lang == "en" else "modo artista", 3 + cw + 10, 1, "mute", "mode label"))
    ex_w, sv_w, lg_w = tw(t["export"]) + 18, tw(t["save"]) + 18, tw(t["lang"]) + 18
    rx = W - 3 - ex_w
    items.append(button(rx, 1, ex_w, 14, "button export png", icon="export", label=t["export"]))
    rx -= sv_w + 2
    items.append(button(rx, 1, sv_w, 14, "button save", icon="save", label=t["save"]))
    rx -= lg_w + 6
    items.append(button(rx, 1, lg_w, 14, "button language", icon="language", label=t["lang"]))
    return G("toolbar (artist)", *items)

def tool_column(active="circle", fill="#fff3c4", chips_y=None):
    items = [R(0, 50, 20, STATUS_Y - 50, "panel", "tool column bg"), R(19, 50, 1, STATUS_Y - 50, "void", "tool column line")]
    y = 53
    for tl in ["circle", "rect", "line", "polygon", "stroke", None, "pencil", "brush", None, "text", "layer"]:
        if tl is None:
            items.append(R(5, y + 1, 10, 1, "bevel", "separator")); y += 4; continue
        items.append(button(3, y, 14, 14, f"tool {tl}" + (" (active)" if tl == active else ""), icon=tl, state="active" if tl == active else "normal"))
        y += 15
    cy = chips_y or (STATUS_Y - 27)
    items.append(G("fill / stroke chips", sw(7, cy + 5, 10, "#3b2a6b", "stroke chip"), sw(3, cy, 10, fill, "fill chip (front)")))
    return G("tool strip (vertical)", *items)

def inspector(lang, parts_fn, title_icon, title):
    items = [R(0, MAIN_Y, W, 20, "panel", "inspector bg"), R(0, MAIN_Y + 19, W, 1, "void", "inspector bottom line"),
             I(title_icon, 6, MAIN_Y + 6, "lizard"), T(title, 18, MAIN_Y + 3, "bright", "selection name")]
    x = 18 + tw(title) + 6
    items.append(R(x, MAIN_Y + 5, 1, 10, "bevel", "separator"))
    items.append(parts_fn(x + 7, MAIN_Y + 2))
    return G("inspector (context bar)", *items)

def ins_circle(lang):
    def f(x, y):
        it = [T("center", x, y + 1, "mute"), T("(", x + 46, y + 1, "mute")]
        fx = x + 54
        it.append(sfield(fx, y, 29, "360", "center.x", "amber")); fx += 31
        it.append(T(",", fx - 1, y + 1, "mute")); fx += 6
        it.append(sfield(fx, y, 22, "90", "center.y", "amber")); fx += 23
        it.append(T(")", fx, y + 1, "mute")); fx += 12
        it.append(T("r", fx, y + 1, "mute")); fx += 11
        it.append(sfield(fx, y, 22, "34", "r", "amber")); fx += 30
        it.append(T("fill", fx, y + 1, "mute")); fx += 32
        it.append(sfield(fx, y, tw("#fff3c4") + 22, "#fff3c4", "fill", swatch="#fff3c4")); fx += tw("#fff3c4") + 26
        it.append(button(fx, y + 1, tw("Fill") + 22, 14, "button Fill ▸ (opens fill editor)", icon="fill", label="Fill"))
        return G("fields (from circle(center, r, fill))", *it)
    return f

BLEND_MODES = [("Normal", "", "normal"), ("Multiply", "MUL", "multiply"), ("Darken", "DRK", "darken"), ("Screen", "SCR", "screen"),
               ("Lighten", "LGT", "lighten"), ("Add", "ADD", "add"), ("Overlay", "OVR", "overlay")]
BLEND_ES = {"Normal": "Normal", "Multiply": "Multiplicar", "Screen": "Trama", "Overlay": "Superponer", "Darken": "Oscurecer", "Lighten": "Aclarar", "Add": "Sumar"}

def blend_px(mode, a, b):
    out = []
    for x, y in zip(hexrgb(a), hexrgb(b)):
        x, y = x / 255, y / 255
        r = {"normal": y, "multiply": x * y, "screen": 1 - (1 - x) * (1 - y), "overlay": 2 * x * y if x < 0.5 else 1 - 2 * (1 - x) * (1 - y),
             "darken": min(x, y), "lighten": max(x, y), "add": min(1, x + y)}[mode]
        out.append(r * 255)
    return rgbhex(out)

def blend_list(x, y, cur="Multiply", lang="en", w=None):
    labels = [BLEND_ES[m] if lang == "es" else m for m, _, _ in BLEND_MODES]
    w = w or max(tw(l) for l in labels) + 18 + 12 + 30
    h = len(BLEND_MODES) * 14 + 4 + 8
    items = chrome(x, y, w, h, name="blend list")
    yy = y + 2
    for i, (m, ab, key) in enumerate(BLEND_MODES):
        if i == 1:
            items.append(R(x + 4, yy + 1, w - 8, 1, "bevel", "separator")); yy += 4
        if i in (2, 4, 6) and False:
            pass
        lb = labels[i]
        it = []
        is_cur = m == cur
        if is_cur:
            it += [R(x + 1, yy, w - 2, 14, "moss", "item bg (current)"), R(x + 1, yy, 2, 14, "lizard", "bar")]
        res = blend_px(key, "#f2a65a", "#8a6fb0")
        it += [G("preview: sky under shade", R(x + 6, yy + 2, 10, 10, "void", "outline"), R(x + 7, yy + 3, 4, 8, "#f2a65a", "under"), R(x + 11, yy + 3, 4, 8, res, "result"))]
        it.append(T(lb, x + 20, yy, "bright" if is_cur else "text"))
        if ab:
            it.append(T(ab, x + w - 6 - tw(ab), yy, "text" if is_cur else "mute", "badge " + ab))
        items.append(G(f"blend mode {m}" + (" (current)" if is_cur else ""), *it))
        yy += 14
        if i in (2, 5):
            items.append(R(x + 4, yy + 1, w - 8, 1, "bevel", "separator")); yy += 4
    hh = yy + 2 - y
    items = chrome(x, y, w, hh, name="blend list") + items[4:]
    return G("blend mode dropdown (open)", *items), w, hh

THUMBS = []
def thumb_g(name, x, y, body, bg_check=True):
    """Frame + checker stay in the row's pixels; the rendered drawing is queued as a vector drawn on top."""
    it = [frame_lines(x, y, 29, 20, "void", "thumb outline")]
    if bg_check:
        it.append(checker(x + 1, y + 1, 27, 18, "#d4dacd", "#a0ab9b", 3, "transparency checker"))
    THUMBS.append(vec(f"thumbnail {name} (vector)", x + 1, y + 1, 480, 320, body, 27 / 480))
    return G(f"thumbnail {name}", *it)

def a_layer_row(x0, y, w, name, body, sel=False, open_=None, badge=None, kind="layer", bg_check=True):
    it = []
    if sel:
        it += [R(x0, y, w, 22, "moss", "row bg (selected)"), R(x0, y, 2, 22, "lizard", "selection bar")]
    it += [I("eye", x0 + 4, y + 7, "text" if sel else "mute", "toggle visible (eye)"), I("lock", x0 + 14, y + 7, "dim", "toggle lock")]
    if open_ is not None:
        it.append(I("chevron-down" if open_ else "chevron-right", x0 + 24, y + 7, "mute", "expand chevron"))
    it.append(thumb_g(name, x0 + 34, y + 1, body, bg_check))
    it.append(T(name, x0 + 67, y + 5, "bright" if sel else "text", "layer name " + name))
    if badge:
        bx = x0 + w - tw(badge) - 8
        it.append(G("blend badge", R(bx, y + 6, tw(badge) + 4, 10, "raised", "badge bg"), T(badge, bx + 2, y + 5, "text" if sel else "mute")))
    return G(f"layer row {name}" + (" (selected)" if sel else ""), *it), 22

def a_obj_row(x0, y, w, kind, name, sel=False):
    it = []
    if sel:
        it += [R(x0, y, w, ROW_H, "moss", "row bg (selected)"), R(x0, y, 2, ROW_H, "lizard", "selection bar")]
    it += [I("eye", x0 + 4, y + 5, "text" if sel else "mute"), I("lock", x0 + 14, y + 5, "dim"), R(x0 + 27, y, 1, ROW_H, "bevel", "indent guide"),
           I(kind, x0 + 44, y + 5, "lizard" if sel else "text", "kind icon " + kind), T(name, x0 + 56, y + 3, "bright" if sel else "text")]
    return G(f"object row {name}" + (" (selected)" if sel else ""), *it), ROW_H

def right_column(lang, color, layers, mode=("Normal", "100%"), dropdown=False, title=None):
    t = TX[lang]
    THUMBS.clear()
    x0, w = RC_X, W - RC_X
    items = [R(ART_X1, 50, 2, STATUS_Y - 50, "void", "splitter"), R(ART_X1 + 1, 50, 1, STATUS_Y - 50, "bevel", "splitter highlight"),
             R(x0, 50, w, STATUS_Y - 50, "panel", "right column bg")]
    pal = [G("color header", R(x0, 50, w, 13, "raised", "header bg"), R(x0, 63, w, 1, "void", "header line"), T(title or t["color"], x0 + 5, 50, "mute"),
             I("chevron-down", x0 + w - 11, 53, "mute", "collapse"))]
    pal += sv_square(x0 + 5, 68, 116, 40, color)
    pal += hue_strip(x0 + 125, 68, 8, 40, color, vertical=True)
    pal += alpha_strip(x0 + 137, 68, 8, 40, color, vertical=True)
    pal.append(sfield(x0 + 5, 112, 57, color[:7], "hex"))
    pal.append(sfield(x0 + 64, 112, 36, "100%", "alpha", "amber"))
    pal.append(button(x0 + 102, 113, 14, 14, "button eyedropper", icon="dropper"))
    pal.append(G("variables (solid and gradient chips; names in tooltip)", *[var_chip(x0 + 5 + i * 12, 131, 10, 10, v) for i, v in enumerate(VARS)]))
    pal.append(R(x0 + 5 + 7 * 12 + 1, 132, 1, 8, "bevel", "separator"))
    pal.append(G("recent colors", *[sw(x0 + 5 + 7 * 12 + 5 + i * 12, 131, 10, c, f"recent {c}") for i, c in enumerate(RECENT[:4])]))
    items.append(G("palette (compact)", *pal))
    y = 145
    lay_items = [R(x0, y, w, 13, "raised", "header bg"), R(x0, y + 13, w, 1, "void", "header line"), T(t["layers"], x0 + 5, y, "mute"),
                 I("plus", x0 + w - 22, y + 3, "lizard", "add layer"), I("trash", x0 + w - 11, y + 3, "mute", "delete layer")]
    y += 16
    md = BLEND_ES[mode[0]] if lang == "es" else mode[0]
    lay_items.append(sfield(x0 + 5, y, 82, md, "blend mode", "text", chev=True, focused=dropdown))
    lay_items.append(G("opacity", I("opacity", x0 + 92, y + 4, "mute"), sfield(x0 + 103, y, 42, mode[1], "opacity", "amber", chev=False)))
    y += 19
    rows = []
    for L in layers:
        if L["kind"] == "layer":
            el, h = a_layer_row(x0, y, w, L["name"], L["body"], sel=L.get("sel"), open_=L.get("open"), badge=L.get("badge"))
        else:
            el, h = a_obj_row(x0, y, w, L["icon"], L["name"], sel=L.get("sel"))
        rows.append(el); y += h
    lay_items.append(G("rows", *rows))
    py = STATUS_Y - 23
    assert y <= py, ("layers overflow", y, py)
    bgb = '<rect x="0" y="0" width="480" height="320" fill="#f6d9a8"/>'
    el, _ = a_layer_row(x0, py + 1, w, "canvas" if lang == "en" else "lienzo", bgb, bg_check=False)
    lay_items.append(G("pinned canvas row (background)", R(x0, py, w, 1, "void", "pin line"), R(x0, py + 1, w, 22, "coal", "pinned bg"), el))
    items.append(G("layers panel (artist)", *lay_items))
    out = G("right column (palette + layers)", *items) + G("layer thumbnails (rendered, vector)", *THUMBS)
    THUMBS.clear()
    if dropdown:
        dl, dw, dh = blend_list(x0 + 5, 161 + 17, mode[0], lang)
        out += dl
    return out

def art_layers(b, sel="circle", shadows=False, lang="en"):
    L = []
    if shadows:
        L.append(dict(kind="layer", name="shadows", body=SHADOWS, sel=(sel == "shadows"), open=False, badge="MUL"))
    L.append(dict(kind="layer", name="sun", body=b["sun"], open=not shadows, sel=(sel == "sun")))
    if not shadows:
        L += [dict(kind="obj", icon="circle", name="circle", sel=(sel == "circle")), dict(kind="obj", icon="circle", name="circle")]
    L += [dict(kind="layer", name="boat", body=b["boat"], open=False), dict(kind="layer", name="sea", body=b["sea"], open=False),
          dict(kind="layer", name="sky", body=b["sky"], open=False)]
    return L

def art_lay(drawer=False):
    if drawer:
        return Lay(None, cw=0, cx0=20, cx1=ART_X1, cy0=50, cy1=196)
    return Lay(None, cw=0, cx0=20, cx1=ART_X1, cy0=50, cy1=STATUS_Y)

def f_artist(drawer=False):
    lang = "en"; t = TX[lang]
    lay = art_lay(drawer)
    b = layer_bodies()
    x0, y0 = lay.cvs((326, 56)); x1, y1 = lay.cvs((394, 124))
    ov = G("selection (circle line 10)", sel_box(x0, y0, x1, y1), gen.anchor(*lay.cvs((360, 90)), "anchor center", "normal"),
           gen.anchor(*lay.cvs((394, 90)), "anchor radius", "alt"))
    parts = [toolbar_artist(lang, drawer), tabbar2(), inspector(lang, ins_circle(lang), "circle", "circle · sun"), tool_column("circle"),
             canvas_area2(lay, drawing_g(lay, b), ov, splitter=False)]
    if drawer:
        dy = 196
        dw = ART_X1
        it = [R(0, dy, dw, STATUS_Y - dy, "panel", "drawer bg"), R(0, dy, dw, 1, "void", "drawer top line"),
              R(0, dy + 1, dw, 13, "raised", "drawer header bg"), R(0, dy + 14, dw, 1, "void", "header line"),
              T("</> CODE · layers.lzdraw", 6, dy + 1, "mute"), I("chevron-down", dw - 22, dy + 4, "mute", "lower drawer"),
              I("close", dw - 11, dy + 4, "mute", "close drawer")]
        rows = [r for r in rows_layers() if r.get("line") not in (7, 8)]
        y = dy + 16
        rg = []
        for r in rows:
            if r["kind"] == "layer": y += LAYER_GAP
            if y + ROW_H > STATUS_Y - 1: break
            el, h = gen.row_el(r, y, w=dw, sel=r.get("line") == 10)
            rg.append(el); y += h
        parts.append(G("code drawer (same rows as coder mode)", *it, G("rows", *rg)))
    parts += [right_column(lang, "#fff3c4", art_layers(b)), status2(lang, lay, pos="x 360 y 90")]
    return window2(parts)

def f_blend_coder():
    lang = "en"; t = TX[lang]
    lay = Lay("rail")
    rows = [dict(kind="value", d=0, code="let horizon = 200", line=7), dict(kind="blank", line=8),
            dict(kind="layer", d=0, code="layer shadows", line=9, open=True, count="MUL 80%"),
            dict(kind="blend", d=1, code='blend("multiply")', line=10), dict(kind="opacity", d=1, code="opacity(0.8)", line=11),
            dict(kind="polygon", d=1, code="polygon([(0, 120), (150, 96), (230, 320), (0, 320)], #8a6fb0)", line=12),
            dict(kind="ellipse", d=1, code="ellipse((152, 240), (44, 7), #3b2a6b)", line=13), dict(kind="blank", line=14),
            dict(kind="layer", d=0, code="layer sun", line=15, open=False, count=2), dict(kind="blank", line=18),
            dict(kind="layer", d=0, code="layer boat", line=19, open=False, count=3), dict(kind="blank", line=23),
            dict(kind="layer", d=0, code="layer sea", line=24, open=False, count=6)]
    def f(y):
        lx = 6 + tw("opacity") + 4
        md = sfield(lx, y - 2, 90, "Multiply", "blend mode", chev=True, focused=True, h=18)
        op = G("field opacity (slider)", T("opacity", 6, y + 20, "mute"), slider(lx, y + 19, 120, 0.8, name="opacity slider"),
               sfield(lx + 126, y + 18, 36, "80%", "opacity value", "amber", h=18))
        return G("fields (layer: blend, opacity)", T("blend", 6, y, "mute"), md, op)
    props = props2(lay, lang, f"layer shadows · {t['line']} 9", f, None)
    b = layer_bodies()
    shadows = vec("layer shadows (vector) [blend MULTIPLY, opacity 80%]", lay.dx, lay.dy, 480, 320, SHADOWS, lay.z)
    x0, y0 = lay.cvs((0, 96)); x1, y1 = lay.cvs((230, 320))
    ov = G("selection (layer shadows)", sel_box(x0, y0, x1, y1))
    dl, dw, dh = blend_list(0, 0)
    dl, dw, dh = blend_list(6 + tw("opacity") + 4, PROPS_Y + 20 - dh - 1)
    return window2([gen.toolbar(None, lang), tabbar2(unsaved=True), code_panel(lay, lang, rows, sel_line=9, props=props),
                    canvas_area2(lay, drawing_g(lay, b), ov, after=[shadows]), rail(lay), dl,
                    status2(lang, lay, saved=False, pos=None, objects=(16, 1800))])

def f_blend_artist():
    lang = "en"; t = TX[lang]
    lay = art_lay()
    b = layer_bodies()
    shadows = vec("layer shadows (vector) [blend MULTIPLY, opacity 80%]", lay.dx, lay.dy, 480, 320, SHADOWS, lay.z)
    x0, y0 = lay.cvs((0, 96)); x1, y1 = lay.cvs((230, 320))
    ov = G("selection (layer shadows)", sel_box(x0, y0, x1, y1))
    def ins(x, y):
        return G("fields (layer)", T("blend", x, y + 1, "mute"), sfield(x + 40, y, 82, "Multiply", "blend", chev=True),
                 T("opacity", x + 130, y + 1, "mute"), slider(x + 183, y, 64, 0.8, name="opacity slider"), sfield(x + 252, y, 36, "80%", "opacity", "amber"))
    return window2([toolbar_artist(lang), tabbar2(unsaved=True), inspector(lang, ins, "layer", "layer · shadows"), tool_column("polygon", "#8a6fb0"),
                    canvas_area2(lay, drawing_g(lay, b), ov, after=[shadows], splitter=False),
                    right_column(lang, "#8a6fb0", art_layers(b, sel="shadows", shadows=True), mode=("Multiply", "80%"), dropdown=True),
                    status2(lang, lay, saved=False, pos=None, objects=(16, 1800))])

GLOW_STOPS = [(0, "#fff3c4", False), (0.6, "#f2a65a", True), (1, "#f2a65a00", False)]
GLOW_CODE = "circle((360, 90), 56, radial_gradient([#fff3c4, hard(#f2a65a, 0.6), #f2a65a00]))"

def f_spanish_menu():
    lang = "es"; t = TX[lang]
    lay = Lay("rail")
    rows = rows_layers(glow_code=GLOW_CODE)
    props = props2(lay, lang, f"circle · {t['line']} 11", circ_cr("56"), "línea 11 · clic derecho: más")
    b = layer_bodies(glow=("radial", (360, 90), 56, GLOW_STOPS), core=None)
    ys = row_heights(rows)
    m, mw, mh = menu_box(36, ys[11] + 10, t["menu"], "context menu on row (right-click, two-finger click)", title="circle · línea 11", hover="Relleno")
    fy = ys[11] + 10 + 16 + 4 * 14 + 4 + 2 * 14 + 4 + 2 - 2
    sub, sw_, sh_ = menu_box(36 + mw - 2, fy, t["fill_sub"], "submenu Relleno ▸", checks=("Degradado radial",), hover="Editar relleno…",
                             swatches=PALETTE + [("r", c) for c in RECENT[:3]])
    return window2([gen.toolbar(None, lang), tabbar2(unsaved=True), code_panel(lay, lang, rows, sel_line=11, props=props),
                    canvas_area2(lay, drawing_g(lay, b)), rail(lay), m, sub, cursor(36 + mw + 40, fy + 4 * 14 + 8),
                    status2(lang, lay, saved=False, pos=None)])

def f_spanish_fill():
    lang = "es"; t = TX[lang]
    lay = Lay("rail")
    rows = rows_layers(glow_code=GLOW_CODE)
    props = props2(lay, lang, f"circle · {t['line']} 11", circ_cr("56"), "línea 11 · arrastra el centro o el radio")
    b = layer_bodies(glow=("radial", (360, 90), 56, GLOW_STOPS), core=None)
    cx, cy = [int(v) for v in lay.cv((360, 90))]
    rr = int(56 * lay.z)
    hnd = G("gradient handles (radial)", pixel_circle(cx, cy, rr + 1, "void", "radius ring outline"), pixel_circle(cx, cy, rr, "bright", "radius ring"),
            gen.anchor(cx * S, cy * S, "handle center", "alt"), gen.anchor((cx + rr) * S, cy * S, "handle radius", "normal"))
    lines = ["radial_gradient(", "  [#fff3c4,", "   hard(#f2a65a, 0.6),", "   #f2a65a00])"]
    fe, fh = fill_editor(lay.cx0 - FE_W - 4, MAIN_Y + 22, lang, 2, GLOW_STOPS, sel=1, code_lines=lines, title="línea 11", at="60%",
                         tail=("right", MAIN_Y + 60))
    return window2([gen.toolbar(None, lang), tabbar2(unsaved=True), code_panel(lay, lang, rows, sel_line=11, props=props),
                    canvas_area2(lay, drawing_g(lay, b), hnd), rail(lay), fe,
                    status2(lang, lay, saved=False, pos=None)])

# ---------- tutorial ----------
def coach(x, y, w, title, body, step, total, lang, tail=None, done=False, lizard=False, next_label=None):
    t = {"en": ("Skip", "Next", "Let's draw!"), "es": ("Saltar", "Siguiente", "¡A dibujar!")}[lang]
    h = 8 + 15 + len(body) * 14 + 8 + 18 + 6
    items = chrome(x, y, w, h, name="coach-mark")
    items.append(R(x + 1, y + 1, 2, h - 2, "lizard", "accent bar"))
    items.append(T(title, x + 9, y + 6, "bright", "title"))
    for i, ln in enumerate(body):
        items.append(T(ln, x + 9, y + 22 + i * 14, "text", "body " + ln))
    by = y + h - 22
    dots = [R(x + 9 + i * 7, by + 7, 4, 4, "lizard" if i == step else "bevel", f"step dot {i+1}" + (" (current)" if i == step else "")) for i in range(total)]
    items.append(G("step dots", *dots))
    nl = next_label or (t[2] if done else t[1])
    nw = tw(nl) + 14
    nx = x + w - 8 - nw
    items.append(G("button primary " + nl, frame_lines(nx, by, nw, 16, "void"), R(nx + 1, by + 1, nw - 2, 14, "moss", "face"),
                   gen.bevel(nx + 1, by + 1, nw - 2, 14, "lizard", "void"), T(nl, nx + 7, by + 1, "bright")))
    if not done:
        items.append(T(t[0], nx - 10 - tw(t[0]), by + 1, "mute", "skip"))
    if tail:
        side, p = tail
        tl = []
        for k in range(6):
            span = 11 - 2 * k
            if span <= 0: break
            top = p - span // 2
            if side == "up":
                tl += [R(top, y - k, span, 1, "void", "tail edge")] + ([R(top + 1, y - k, span - 2, 1, "row", "tail fill")] if span > 2 else [])
            elif side == "down":
                tl += [R(top, y + h - 1 + k, span, 1, "void", "tail edge")] + ([R(top + 1, y + h - 1 + k, span - 2, 1, "row", "tail fill")] if span > 2 else [])
            elif side == "left":
                tl += [R(x - k, top, 1, span, "void", "tail edge")] + ([R(x - k, top + 1, 1, span - 2, "row", "tail fill")] if span > 2 else [])
            elif side == "right":
                tl += [R(x + w - 1 + k, top, 1, span, "void", "tail edge")] + ([R(x + w - 1 + k, top + 1, 1, span - 2, "row", "tail fill")] if span > 2 else [])
        items.append(G("tail", *tl))
    return G(f"coach-mark step {step+1} of {total}", *items), h

def tool_span():
    x = 3; out = []
    for tool in gen.TOOLS:
        if tool is None: x += 4; continue
        out.append((tool[0], x)); x += 15
    return dict(out)

TUT_CANVAS = 'canvas = {"size": (480, 320), "background": #ffffff}'
def tut_body(step):
    b = {"bg": '<rect x="0" y="0" width="480" height="320" fill="#ffffff"/>'}
    if step >= 1:
        b["c"] = f'<circle cx="240" cy="160" r="60" fill="{"#4fc46a" if step == 1 else "#e0607a"}"/>'
    if step >= 3:
        pts = [(70, 250), (150, 200), (250, 240), (340, 170), (420, 210)]
        b["s"] = gen.brush_shape(pts, [0.2, 0.9, 0.5, 1.0, 0.3], 2, 14, 0.3, "#3b2a6b", 1, "brush stroke")
    return b

def lizard_runs(i=3):
    src = open(os.path.join(gen.ROOT, "lib/ui/splash_art.liz"), encoding="utf-8").read()
    sps = re.findall(r'\{"size": \((\d+), (\d+)\), "runs": \[(.*?)\]\}', src)
    lw, lh = int(sps[i][0]), int(sps[i][1])
    return lw, lh, [(int(a), int(b), int(c), d) for a, b, c, d in re.findall(r'\[(\d+), (\d+), (\d+), (#[0-9a-f]{6})\]', sps[i][2])]

TUT = {
 "en": [("Make your first shape", ["Pick a shape tool up here, then", "click the canvas. Every shape is", "one line of code."]),
        ("Change its color", ["Click a color in the palette, or", "the fill swatch below. The code", "changes as you pick."]),
        ("Draw your first brush stroke", ["Pick the brush and drag on the", "canvas. Press harder for a", "thicker line."]),
        ("Peek at the code you made", ["Every shape you made is a line", "here. Change the code and the", "drawing follows."])],
 "es": [("Crea tu primera forma", ["Elige una forma aquí arriba y", "haz clic en el lienzo. Cada", "forma es una línea de código."])],
}

def f_tutorial(step, lang="en"):
    t = TX[lang]
    lay = Lay("open")
    base = [dict(kind="comment", d=0, code="-~> lzdraw 1", line=1), dict(kind="blank", line=2)]
    if step >= 1:
        base.append(dict(kind="circle", d=0, code=f'circle((240, 160), 60, {"#4fc46a" if step == 1 else "#e0607a"})', line=4))
    if step >= 3:
        base.append(dict(kind="brush", d=0, code="brush([(70, 250), (150, 200), (250, 240), (340, 170), (420, 210)], [0.2, 0.9, 0.5, 1, 0.3], #3b2a6b, pens.ink)", line=3))
        base[-1], base[-2] = base[-2], base[-1]
        base[-2]["line"], base[-1]["line"] = 3, 4
    if step == 0:
        base.append(dict(kind="add", d=0, line=None))
    sel = {1: 4, 2: None, 3: None}.get(step)
    if step == 1:
        def pf(y):
            g1, x = gen.point_field(6, y, "center", "240", "160", "center")
            g2 = gen.num_field(x + 4, y, "r", "60", "r")
            g3, _ = gen.color_field(6, y + 20, "fill", "#4fc46a", "fill")
            return G("fields (fill wraps to a second row at 256)", g1, g2, g3)
        props = props2(lay, lang, f"circle · {t['line']} 4", pf, None)
    else:
        props = props2(lay, lang, "", lambda y: T("nothing selected" if lang == "en" else "nada seleccionado", 6, y + 2, "mute"), None)
    b = tut_body(step)
    ov = ""
    if step == 1:
        x0, y0 = lay.cvs((180, 100)); x1, y1 = lay.cvs((300, 220))
        ov = G("selection (circle line 4)", sel_box(x0, y0, x1, y1), gen.anchor(*lay.cvs((240, 160)), "anchor center"), gen.anchor(*lay.cvs((300, 160)), "anchor radius", "alt"))
    if step == 2:
        pts = [lay.cv(p) for p in [(70, 250), (150, 200), (250, 240), (340, 170), (420, 210)]]
        dots = []
        for (ax, ay), (bx, by) in zip(pts, pts[1:]):
            for k in range(0, 10, 2):
                f = k / 10
                dots.append(R(int(ax + (bx - ax) * f), int(ay + (by - ay) * f), 1, 1, "lizard", "dot"))
        ov = G("ghost stroke hint (dotted lizard)", *dots)
    active = "brush" if step == 2 else None
    tb = gen.toolbar(active, lang)
    file = "untitled.lzdraw"
    parts = [tb, tabbar2(file, unsaved=step > 0), code_panel(lay, lang, base, sel_line=sel, props=props, file=file, canvas_code=TUT_CANVAS),
             canvas_area2(lay, drawing_g(lay, b, order=("bg", "c", "s")), ov),
             sidebar(lay, lang, "#e0607a" if step >= 1 else "#4fc46a", hover_add=False),
             status2(lang, lay, saved=step == 0, objects=(min(step, 2) if step < 3 else 2, 1800), file=file, pos=None)]
    ts = tool_span()
    title, body = TUT[lang][step if lang == "en" else 0]
    if step == 0:
        hole = [(ts["circle"] - 1, 0, ts["layer"] + 15, 16)]
        bub, h = coach(ts["circle"] + 2, 26, 250, title, body, 0, 4, lang, tail=("up", ts["circle"] + 7))
        parts += [dither_dim(0, 0, W, H, hole), bub]
    elif step == 1:
        holes = [(lay.sb_x, MAIN_Y + 186 - 1, W, MAIN_Y + 186 + 14 + 4 * ROW_H + 2), (154, PROPS_Y + 20, 252, PROPS_Y + 42)]
        fx1 = 6 + tw("fill") + 4 + tw("#4fc46a") + 22
        holes = [(lay.sb_x + 1, 177, W - 1, 265), (4, PROPS_Y + 41, fx1 + 2, PROPS_Y + 61)]
        bub, h = coach(6, PROPS_Y + 41 - 1 - 97 - 7, 250, title, body, 1, 4, lang, tail=("down", 60))
        parts += [dither_dim(0, 0, W, H, holes), bub]
    elif step == 2:
        holes = [(ts["brush"] - 1, 0, ts["brush"] + 15, 16), (lay.cx0, lay.dy - 4, lay.cx1, lay.dy + lay.dh + 4)]
        bub, h = coach(ts["brush"] - 8, 26, 230, title, body, 2, 4, lang, tail=("up", ts["brush"] + 7))
        cx, cy = [int(v) for v in lay.cv((150, 200))]
        parts += [dither_dim(0, 0, W, H, holes), bub, cursor(cx, cy)]
    else:
        ys = row_heights(base)
        holes = [(0, ys[3] - 1, lay.cw, ys[4] + ROW_H + 1)]
        bx = 70
        bub, h = coach(bx, ys[4] + ROW_H + 30, 240, title, body, 3, 4, lang, tail=("up", 120), done=True)
        lw, lh, runs = lizard_runs(3)
        k = 2
        lx, ly = bx + 240 - lw * k - 10, ys[4] + ROW_H + 30 - lh * k + 4
        liz = gen.sprite_svg(runs, lx, ly, k, "lizard celebrating (frame 3, x2)")
        sparks = G("sparkles", I("sparkle", lx - 12, ly - 4, "amber"), I("sparkle", lx + lw * k + 2, ly - 8, "lizard"), I("sparkle", lx + lw * k + 14, ly + 12, "sky"),
                   I("sparkle", lx - 24, ly + 8, "peach"))
        parts += [dither_dim(0, 0, W, H, holes), bub, liz, sparks]
    return window2(parts)

# ---------- new drawing dialog ----------
ND = {
 "en": dict(title="NEW DRAWING", name="Name", size="Canvas size", custom="Custom", rec="Recommended: up to ~1,800 objects", measured="measured on this Mac",
            cancel="Cancel", create="Create", tip=["Measured when Linkz starts, from this Mac's", "memory and drawing speed. Above it, repaints", "may get slow. You can raise it in Settings."],
            file="sunset-2"),
 "es": dict(title="NUEVO DIBUJO", name="Nombre", size="Tamaño del lienzo", custom="Otro", rec="Recomendado: hasta ~1.800 objetos", measured="medido en este Mac",
            cancel="Cancelar", create="Crear", tip=["Se mide al abrir Linkz, según la memoria", "y la velocidad de dibujo de este Mac. Por", "encima, repintar puede ir lento. Puedes", "subirlo en Ajustes."],
            file="atardecer-2"),
}

def f_new_drawing(lang="en"):
    import xml.etree.ElementTree as ET
    d = ND[lang]
    bg = home.frame_gallery(lang)
    NS = "http://www.w3.org/2000/svg"
    ET.register_namespace("", NS)
    root = ET.fromstring(bg)
    ui = next(e for e in root if (e.get("id") or "").startswith("ui ("))
    for e in list(ui):
        if e.get("id") in ("card menu ⋯ (open)", "mouse cursor"):
            ui.remove(e)
    for e in ui.iter():
        if (e.get("id") or "").startswith("card ") and "(hover" in (e.get("id") or ""):
            pass
    w, h = 300, 214
    x, y = (W - w) // 2, 64
    it = [f'<rect id="modal dim (void 62%)" x="0" y="0" width="{W}" height="{H}" fill="{C["void"]}" fill-opacity="0.62"/>']
    it += chrome(x, y, w, h, bg="panel", name="dialog")
    it.append(header(x, y, w, d["title"]))
    yy = y + 22
    it.append(T(d["name"], x + 12, yy, "mute", "label name"))
    it.append(sfield(x + 12, yy + 14, w - 24, d["file"], "name", "text", focused=True, h=18))
    it.append(R(x + 12 + 4 + tw(d["file"]), yy + 17, 1, 12, "amber", "caret"))
    it.append(T(".lzdraw", x + 12 + 6 + tw(d["file"]), yy + 17, "dim", "extension"))
    yy += 40
    it.append(T(d["size"], x + 12, yy, "mute", "label size"))
    presets = ["480×320", "512×512", "1080×1080", "1920×1080", d["custom"]]
    bx, by = x + 12, yy + 14
    row = []
    for i, p in enumerate(presets):
        pw = tw(p) + 8
        if bx + pw > x + w - 12:
            bx, by = x + 12, by + 17
        row.append(button(bx, by, pw, 15, f"preset {p}" + (" (selected)" if i == 0 else ""), label=p, state="active" if i == 0 else "normal"))
        bx += pw + 3
    it.append(G("size presets", *row))
    yy = by + 22
    it.append(G("custom size", T("W", x + 12, yy + 2, "mute"), sfield(x + 24, yy, 40, "480", "width", "amber", h=18), T("×", x + 68, yy + 2, "mute"),
                T("H", x + 80, yy + 2, "mute"), sfield(x + 92, yy, 40, "320", "height", "amber", h=18), T("px", x + 136, yy + 2, "mute")))
    yy += 26
    it.append(R(x + 12, yy, w - 24, 1, "bevel", "rule"))
    yy += 6
    it.append(G("recommended objects", I("objects", x + 12, yy + 3, "lizard"), T(d["rec"], x + 24, yy, "text"),
                I("info", x + 24 + tw(d["rec"]) + 4, yy + 3, "sky", "info button (hovered)"),
                T(d["measured"], x + 24, yy + 13, "mute")))
    by2 = y + h - 24
    cw_, crw = tw(d["cancel"]) + 12, tw(d["create"]) + 22
    cx = x + w - 12 - crw
    it.append(G("button primary " + d["create"], frame_lines(cx, by2, crw, 18, "void"), R(cx + 1, by2 + 1, crw - 2, 16, "moss", "face"),
                gen.bevel(cx + 1, by2 + 1, crw - 2, 16, "lizard", "void"), I("plus", cx + 6, by2 + 5, "lizard"), T(d["create"], cx + 16, by2 + 2, "bright")))
    it.append(button(cx - 4 - cw_, by2, cw_, 18, "button " + d["cancel"], label=d["cancel"]))
    tipw = max(tw(s) for s in d["tip"]) + 14
    ix = x + 24 + tw(d["rec"]) + 8
    tip = gen.tooltip_box(min(ix - 20, W - tipw - 6), yy - len(d["tip"]) * 15 - 6 - 6, [lambda a, b, s=s: T(s, a, b, "text") for s in d["tip"]], "tooltip recommended objects", w=tipw)
    it.append(G("new drawing dialog", *[i for i in it[1:]]))
    frag = f'<svg xmlns="{NS}"><g>{it[0]}{it[-1]}{tip}{cursor(ix + 3, yy + 7)}</g></svg>'
    for e in ET.fromstring(frag)[0]:
        ui.append(e)
    return ET.tostring(root, encoding="unicode")

# ---------- object counter strip ----------
def f_counter():
    lang = "en"
    h = 164
    parts = []
    labels = [("normal", 840), ("near", 1520), ("over", 2130)]
    ys = [6, 40, 128]
    for (st, n), y in zip(labels, ys):
        left = [("canvas", "layers.lzdraw", "text"), ("check", "saved", "lizard"), counter_part(lang, n, 1800)[0]]
        right = [(None, "x 352 y 88", "mute"), (None, "75%", "mute"), ("language", "EN", "text")]
        parts.append(statusbar_y(y, left, right, name=f"status bar · objects {st}"))
        parts.append(T({"normal": "normal (< 80%)", "near": "near the limit (>= 80%)", "over": "over the limit"}[st], 6, y + 16 + 1, "dim", "label " + st))
    x = 6 + 10 + tw("layers.lzdraw") + 6 + 7 + 10 + tw("saved") + 6 + 7
    tip = gen.tooltip_box(x, 74, [lambda a, b: G("title", T("2,130 objects", a, b, "ember"), T(" · recommended 1,800", a + tw("2,130 objects"), b, "mute")),
                                  lambda a, b: T("Repaints may get slow above this on this Mac.", a, b, "text"),
                                  lambda a, b: T("Measured from memory and speed · Settings", a, b, "mute")], "tooltip objects over",
                          w=tw("Repaints may get slow above this on this Mac.") + 14)
    parts += [tip, cursor(x + 60, 133)]
    return svg_doc(W * S, h * S, "".join(parts), bg="coal", scale=S).replace(f'height="{H*S}"', f'height="{h*S}"')

# ---------- pen settings ----------
PS = {
 "en": dict(title="PEN · ink", back="pens", pressure="PRESSURE", presets=["soft", "linear", "firm"], width="width", opacity="opacity", taper="taper",
            smoothing="smoothing", speed="speed", inputs=["Trackpad", "Tablet", "Mouse"], floor="floor", floor_hint="light touches still draw",
            preview="preview · redraws as you tune", name="name", save="Save as new pen", reset="Reset", newname="soft ink", raw="raw", eff="out"),
 "es": dict(title="PLUMA · tinta", back="plumas", pressure="PRESIÓN", presets=["suave", "lineal", "firme"], width="grosor", opacity="opacidad", taper="punta",
            smoothing="suavizado", speed="velocidad", inputs=["Trackpad", "Tableta", "Ratón"], floor="mínimo", floor_hint="un toque leve también pinta",
            preview="vista previa · se redibuja", name="nombre", save="Guardar pluma nueva", reset="Restaurar", newname="tinta suave", raw="entrada", eff="salida"),
}
CURVE = [(0, 0), (0.3, 0.6), (1, 1)]
def curve_eval(pts, x):
    dense = catmull(pts, 24)
    for (a, b), (c, d) in zip(dense, dense[1:]):
        if a <= x <= c:
            return b + (d - b) * ((x - a) / (c - a) if c > a else 0)
    return dense[-1][1]

def pen_settings(x, y, lang, coder=True):
    p = PS[lang]
    w = 320
    items = []
    rst = p["reset"]
    items.append(G("header", R(x + 1, y + 1, w - 2, 13, "raised", "header bg"), R(x + 1, y + 14, w - 2, 1, "void", "header line"),
                   I("chevron-left", x + 5, y + 4, "mute", "back to pen list"), T(p["back"], x + 15, y + 1, "mute"), R(x + 15 + tw(p["back"]) + 4, y + 4, 1, 8, "bevel", "sep"),
                   T(p["title"], x + 15 + tw(p["back"]) + 10, y + 1, "bright"), T(rst, x + w - 22 - tw(rst), y + 1, "mute", "reset (text button)"),
                   I("close", x + w - 12, y + 4, "mute", "close")))
    gx, gy, gs = x + 8, y + 31, 84
    items.append(T(p["pressure"], x + 6, y + 16, "mute", "label pressure"))
    g = [frame_lines(gx - 1, gy - 1, gs + 2, gs + 2, "void", "graph outline"), R(gx, gy, gs, gs, "coal", "graph bg")]
    for k in (1, 2, 3):
        g.append(G(f"grid {k}", *[R(gx + i, gy + gs * k // 4, 1, 1, "bevel", "dot") for i in range(0, gs, 3)],
                   *[R(gx + gs * k // 4, gy + i, 1, 1, "bevel", "dot") for i in range(0, gs, 3)]))
    g.append(pixel_poly([(gx, gy + gs - 1), (gx + gs - 1, gy)], "dim", "linear reference"))
    dense = catmull(CURVE, 24)
    g.append(pixel_poly([(gx + px * (gs - 1), gy + (1 - py) * (gs - 1)) for px, py in dense], "lizard", "pressure curve"))
    for i, (px, py) in enumerate(CURVE):
        ax, ay = gx + int(px * (gs - 1)), gy + int((1 - py) * (gs - 1))
        hov = i == 1
        s_ = 7 if hov else 5
        g.append(G(f"curve point {i+1}" + (" (hovered, dragging)" if hov else ""), R(ax - s_ // 2 - 1, ay - s_ // 2 - 1, s_ + 2, s_ + 2, "void", "outline"),
                   R(ax - s_ // 2, ay - s_ // 2, s_, s_, "amber" if hov else "lizard", "fill")))
    lx = gx + int(0.55 * (gs - 1))
    g.append(G("live pressure (pen down)", *[R(lx, gy + j, 1, 2, "sky", "dash") for j in range(0, gs, 4)]))
    g.append(T(p["raw"], gx + gs - tw(p["raw"]) - 2, gy + gs - 14, "mute", "x axis label"))
    g.append(T(p["eff"], gx + 2, gy + 1, "mute", "y axis label"))
    items.append(G("pressure curve editor", *g))
    items.append(sfield(gx - 1, gy + gs + 4, gs + 2, p["presets"][0], "curve preset (soft / linear / firm)", chev=True))
    rx = x + 104
    lw = max(tw(p[k]) for k in ("width", "opacity", "taper", "smoothing", "speed")) + 6
    vw = tw("80-100%") + 4
    sl_w = w - (rx - x) - lw - vw - 6
    rows = [("width", 1 / 28, 12 / 28, "1-12"), ("opacity", 0.8, 1.0, "80-100%"), ("taper", 0.6, None, "0.3"), ("smoothing", 0.4, None, "40%"), ("speed", 0.2, None, "20%")]
    yy = y + 17
    for key, a, b, val in rows:
        items.append(G(f"setting {key}", T(p[key], rx, yy, "mute"), slider(rx + lw, yy - 1, sl_w, a, b, name=f"{key} slider"),
                       T(val, x + w - 6 - tw(val), yy, "amber")))
        yy += 18
    iw_ = (w - (rx - x) - 6 + 2) // 3
    sg, _ = seg(rx, yy, p["inputs"], 0, "input switch", w=iw_)
    items.append(sg)
    yy += 17
    items.append(G("trackpad overrides", R(rx, yy, w - (rx - x) - 6, 31, "coal", "override bg"), R(rx, yy, 1, 31, "lizard", "override bar"),
                   T(p["floor"], rx + 5, yy + 1, "mute"), slider(rx + lw, yy, sl_w, 0.25, name="floor slider"), T("25%", x + w - 6 - tw("25%"), yy + 1, "amber"),
                   T(p["floor_hint"], rx + 5, yy + 16, "mute")))
    py = y + 156
    pv = [frame_lines(x + 6, py, w - 12, 34, "void", "preview outline"), R(x + 7, py + 1, w - 14, 32, "#fbf8ef", "preview paper")]
    n = 60
    pts = [(14 + i * 4.8, 17 + 8 * math.sin(i / n * math.pi * 2.2)) for i in range(n)]
    raw = [0.1 + 0.85 * math.sin(i / (n - 1) * math.pi) ** 1.5 for i in range(n)]
    eff = [0.25 + 0.75 * curve_eval(CURVE, r) for r in raw]
    ws = [(1 + 11 * e) * min(1, max(0.15, min(i / (n - 1), 1 - i / (n - 1)) / 0.3)) for i, e in enumerate(eff)]
    pv.append(vec("live stroke preview (vector, redraws as you tune)", x + 7, py + 1, w - 14, 32,
                  f'<path d="{variable_outline(pts, ws)}" fill="#3b2a6b" fill-opacity="0.92"/>'))
    items.append(G("live preview", *pv))
    yy = py + 38
    if coder:
        lines = ['pens.soft_ink = {"name": "soft ink",', '  "min": 1, "max": 12, "taper": 0.3,', '  "opacity": [0.8, 1], "pressure": {',
                 '   "curve": [(0, 0), (0.3, 0.6), (1, 1)],', '   "smoothing": 0.4, "speed": 0.2,', '   "trackpad": {"floor": 0.25}}}']
        ch = len(lines) * 13 + 6
        items.append(G("code preview (the pen it writes)", frame_lines(x + 6, yy, w - 12, ch, "void", "inset outline"), R(x + 7, yy + 1, w - 14, ch - 2, "coal", "inset"),
                       *[code(ln, x + 10, yy + 3 + i * 13, name=f"code line {i+1}") for i, ln in enumerate(lines)]))
        yy += ch + 4
    sv_w = tw(p["save"]) + 22
    sx = x + w - 6 - sv_w
    nx = x + 10 + tw(p["name"])
    items.append(G("field name", T(p["name"], x + 6, yy + 2, "mute"), sfield(nx, yy, sx - 4 - nx, p["newname"], "pen name", h=18)))
    items.append(G("button primary " + p["save"], frame_lines(sx, yy, sv_w, 18, "void"), R(sx + 1, yy + 1, sv_w - 2, 16, "moss", "face"),
                   gen.bevel(sx + 1, yy + 1, sv_w - 2, 16, "lizard", "void"), I("plus", sx + 6, yy + 5, "lizard"), T(p["save"], sx + 16, yy + 2, "bright")))
    h = yy + 24 - y
    return G("pen settings panel", *chrome(x, y, w, h, name="pen settings"), *items), h

def f_pen_coder():
    lang = "en"
    lay = Lay("rail")
    rows = rows_layers()
    bx = 3 + 15 * 5 + 4 + 15
    ps, h = pen_settings(bx - 10, TOP_H + 1, lang, True)
    props = props2(lay, lang, "brush · pens.ink", lambda y: T("pen settings write a pen hash: pens.soft_ink = {…}", 6, y + 2, "mute"), None)
    return window2([gen.toolbar("brush", lang), tabbar2(unsaved=True), code_panel(lay, lang, rows, props=props),
                    canvas_area2(lay, drawing_g(lay, layer_bodies())), rail(lay), ps, cursor(bx - 10 + 8 + 25 + 4, TOP_H + 1 + 31 + 34 + 4),
                    status2(lang, lay, saved=False, pos=None)])

def f_pen_artist_es():
    lang = "es"
    lay = art_lay()
    b = layer_bodies()
    stroke = gen.brush_shape([(70, 160), (160, 100), (250, 140), (340, 70), (430, 120)], [0.2, 0.9, 0.5, 1.0, 0.3], 1, 12, 0.3, "#3b2a6b", 1, "brush stroke (tinta suave)")
    ps, h = pen_settings(24, 52, lang, False)
    def ins(x, y):
        return G("fields (brush)", T("pluma", x, y + 1, "mute"), sfield(x + 42, y, 100, "tinta suave", "pen", chev=True),
                 T("color", x + 150, y + 1, "mute"), sfield(x + 188, y, tw("#3b2a6b") + 22, "#3b2a6b", "color", swatch="#3b2a6b"))
    return window2([toolbar_artist(lang), tabbar2(unsaved=True), inspector(lang, ins, "brush", "pincel · sun"), tool_column("brush", "#3b2a6b"),
                    canvas_area2(lay, drawing_g(lay, b, extra=stroke), splitter=False),
                    right_column(lang, "#3b2a6b", art_layers(b, sel="sun")), ps,
                    status2(lang, lay, saved=False, pos=None)])

# ---------- autocomplete ----------
def ac_popup(x, y, items, sel, doc, name):
    w = 214
    h = 4 + len(items) * 14 + (6 + len(doc) * 13 + 4 if doc else 0) + 14
    it = chrome(x, y, w, h, name="autocomplete")
    yy = y + 2
    for i, (ic, label, match, kind) in enumerate(items):
        s = i == sel
        g = []
        if s:
            g += [R(x + 1, yy, w - 2, 14, "moss", "item bg (selected)"), R(x + 1, yy, 2, 14, "lizard", "bar")]
        g.append(I(ic, x + 6, yy + 3, "lizard" if s else "mute"))
        cx = x + 18
        for j, ch in enumerate(label):
            hit = j in match
            g.append(T(ch, cx, yy, "lizard" if hit else ("bright" if s else "text")))
            cx += 7
        g.append(T(kind, x + w - 6 - tw(kind), yy, "text" if s else "mute", "kind " + kind))
        it.append(G(f"completion {label}" + (" (selected)" if s else ""), *g))
        yy += 14
    if doc:
        it.append(R(x + 1, yy + 2, w - 2, 1, "void", "doc line"))
        yy += 5
        for d in doc:
            it.append(d(x + 6, yy)); yy += 13
        yy += 2
    it.append(G("footer", R(x + 1, yy, w - 2, 13, "panel", "footer bg"), T("Tab insert · Esc close", x + 6, yy - 1, "mute")))
    hh = yy + 14 - y
    it = chrome(x, y, w, hh, name="autocomplete") + it[4:]
    return G(name, *it), hh

def sig(x, y, fn, params, cur, extra=""):
    out = [T(fn + "(", x, y, "bright")]
    cx = x + tw(fn + "(")
    for i, pn in enumerate(params):
        out.append(T(pn, cx, y, "lizard" if i == cur else "mute", "param " + pn + (" (current)" if i == cur else "")))
        if i == cur:
            out.append(R(cx, y + 12, tw(pn), 1, "lizard", "current param underline"))
        cx += tw(pn)
        if i < len(params) - 1:
            out.append(T(", ", cx, y, "mute")); cx += 14
    out.append(T(")" + extra, cx, y, "bright"))
    return G("signature " + fn, *out)

def ac_frame(variant):
    lang = "en"
    lay = Lay(None, cw=320)
    rows = rows_layers()
    out = []
    for r in rows:
        out.append(r)
        if r.get("line") == 11:
            break
    typed = {"ci": "ci", "sig": "circle((360, 90), ", "math": "let a = math.s"}[variant]
    new = dict(kind="circle" if variant != "math" else "value", d=1, code=typed, line=12, caret=tw(typed))
    rest = rows[len(out):]
    rows2 = out + [new] + [dict(r, line=r["line"] + 1) if r.get("line") else r for r in rest]
    props = gen.props_panel("en", "line 12", lambda y: T("typing…", 6, y + 2, "mute"), None) if False else None
    props = props2(lay, lang, "line 12", lambda y: T("typing… Tab completes, Enter keeps the line", 6, y + 2, "mute"), None)
    ys = row_heights(rows2)
    cx = 56 + tw(typed)
    parts = [gen.toolbar(None, lang), tabbar2(unsaved=True), code_panel(lay, lang, rows2, sel_line=12, props=props),
             canvas_area2(lay, drawing_g(lay, layer_bodies()))]
    if variant == "ci":
        items = [("circle", "circle", [0, 1], "function"), ("function", "circle_ring", [0, 1], "line 30"), ("function", "city_glow", [0, 1], "value"),
                 ("layer", "city", [0, 1], "layer"), ("keyword", "continue", [0, 4], "keyword"), ("pen", "pens.calligraphy", [5, 10], "pen")]
        doc = [lambda a, b: sig(a, b, "circle", ["center", "r", "fill"], -1), lambda a, b: T("Draws a filled circle.", a, b, "text")]
        pop, h = ac_popup(56 - 18 + 2, ys[12] + ROW_H, items, 0, doc, "autocomplete popup (typing ci)")
        parts.append(pop)
    elif variant == "sig":
        lines = [lambda a, b: sig(a, b, "circle", ["center", "r", "fill"], 1), lambda a, b: G("param doc", T("r", a, b, "lizard"), T("the radius in pixels", a + 14, b, "text"))]
        sw_ = tw("circle(center, r, fill)") + 14
        tip = gen.tooltip_box(cx - tw("circle(center, ") - 6, ys[12] - 38, lines, "signature hint (above the row)", w=sw_)
        parts.append(tip)
    else:
        items = [("function", "sin", [0], "math"), ("function", "sqrt", [0], "math"), ("function", "seed", [0], "math")]
        doc = [lambda a, b: sig(a, b, "math.sin", ["x"], -1), lambda a, b: T("Sine of x in radians.", a, b, "text")]
        pop, h = ac_popup(cx - 18 - 7 + 2, ys[12] + ROW_H, items, 0, doc, "member completion (math.s)")
        parts.append(pop)
    parts.append(status2(lang, lay, saved=False, pos=None))
    return window2(parts)

# ---------- style page components ----------
def comps():
    out = {}
    def wrap(w, h, body): return svg_doc(w * S, h * S, body, scale=S)
    m, mw, mh = menu_box(0, 0, TX["en"]["menu"], "menu (shape)", title="circle · line 10", hover="Fill", disabled=("Bring forward",))
    out["menu/shape"] = wrap(mw + 2, mh + 2, m)
    for st, kw in [("normal", {}), ("hover", {"hover": "Duplicate"}), ("disabled", {"disabled": ("Duplicate",)})]:
        mi, w_, h_ = menu_box(0, 0, [("copy", "Duplicate", "Cmd+D")], f"menu item {st}", **kw)
        out[f"menu item/{st}"] = wrap(w_ + 2, h_ + 2, mi)
    mi, w_, h_ = menu_box(0, 0, [("trash", "Remove", "Del")], "menu item danger")
    out["menu item/danger"] = wrap(w_ + 2, h_ + 2, mi)
    mi, w_, h_ = menu_box(0, 0, [("fill", "Fill", ">")], "menu item submenu")
    out["menu item/submenu"] = wrap(w_ + 2, h_ + 2, mi)
    mi, w_, h_ = menu_box(0, 0, [("solid", "Solid color", "")], "menu item current", checks=("Solid color",))
    out["menu item/current"] = wrap(w_ + 2, h_ + 2, mi)
    b = [frame_lines(0, 0, 120, 14, "void"), px_grad("gradient", 1, 1, 118, 12, lambda u, v: grad_at_a(SKY_STOPS, u))]
    out["gradient bar/smooth"] = wrap(120, 14, G("gradient bar", *b))
    out["gradient bar/stepped"] = wrap(120, 14, G("gradient bar stepped", frame_lines(0, 0, 120, 14, "void"),
                                         *[R(1 + 118 * i // 5, 1, 118 * (i + 1) // 5 - 118 * i // 5, 12, grad_at(SKY_STOPS, (i + .5) / 5), "band") for i in range(5)]))
    out["color stop/normal"] = wrap(11, 10, stop_marker(5, 0, "#3b2a6b"))
    out["color stop/selected"] = wrap(11, 10, stop_marker(5, 0, "#e0607a", sel=True))
    out["color stop/hard"] = wrap(11, 10, stop_marker(5, 0, "#f2a65a", hard=True))
    out["color stop/ghost"] = wrap(11, 10, stop_marker(5, 0, "#000000", ghost=True))
    hp = [R(0, 0, 106, 100, "panel", "bg")] + sv_square(5, 5, 96, 56, "#e0607a") + hue_strip(5, 66, 96, 6, "#e0607a") + alpha_strip(5, 77, 96, 6, "#e0607a") \
         + [sfield(5, 88 - 0, 57, "#e0607a", "hex"), sfield(64, 88, 37, "100%", "alpha", "amber")]
    out["hsv picker"] = wrap(106, 106, G("hsv picker", *hp))
    bb = layer_bodies()
    THUMBS.clear()
    el, _ = a_layer_row(0, 0, 150, "sun", bb["sun"])
    out["layer row/normal"] = wrap(150, 22, R(0, 0, 150, 22, "panel", "bg") + el + "".join(THUMBS))
    THUMBS.clear()
    el, _ = a_layer_row(0, 0, 150, "shadows", SHADOWS, sel=True, open_=False, badge="MUL")
    out["layer row/selected + blend badge"] = wrap(150, 22, R(0, 0, 150, 22, "panel", "bg") + el + "".join(THUMBS))
    THUMBS.clear()
    cm, ch = coach(0, 6, 250, "Make your first shape", TUT["en"][0][1], 0, 4, "en", tail=("up", 20))
    out["coach-mark"] = wrap(254, ch + 10, cm)
    dl, dw, dh = blend_list(0, 0)
    out["blend dropdown (open)"] = wrap(dw + 2, dh + 2, dl)
    out["slider/opacity"] = wrap(80, 16, slider(2, 0, 76, 0.8))
    out["slider/range"] = wrap(80, 16, slider(2, 0, 76, 0.1, 0.45))
    sg, ex = seg(0, 0, ["Solid", "Linear", "Radial"], 1, "tabs", icons=["solid", "gradient", "radial"], w=58)
    out["segmented tabs"] = wrap(ex + 1, 14, sg)
    st = statusbar_y(0, [counter_part("en", 840, 1800)[0]], [], w=180)
    out["status counter/normal"] = wrap(180, 16, st)
    out["status counter/near"] = wrap(180, 16, statusbar_y(0, [counter_part("en", 1520, 1800)[0]], [], w=180))
    out["status counter/over"] = wrap(180, 16, statusbar_y(0, [counter_part("en", 2130, 1800)[0]], [], w=180))
    out["row badge/blend"] = wrap(280, ROW_H, R(0, 0, 280, ROW_H, "panel", "bg") + gen.row_el(dict(kind="layer", d=0, code="layer shadows", line=9, open=True, count="MUL 80%"), 0, w=280)[0])
    pop, ph = ac_popup(0, 0, [("circle", "circle", [0, 1], "function"), ("function", "circle_ring", [0, 1], "line 30")], 0,
                       [lambda a, b: sig(a, b, "circle", ["center", "r", "fill"], -1)], "autocomplete")
    out["autocomplete"] = wrap(216, ph + 2, pop)
    for v in VARS[3:]:
        out[f"variable chip/{ {'s': 'solid', 'l': 'linear', 'r': 'radial', 'st': 'stepped'}[v[1]]}"] = wrap(32, 27, R(0, 0, 32, 27, "panel", "bg") + var_grid(1, 1, 1, 32, [v])[0])
    pop, ph = color_popover(0, 7, "en", "#fff3c4", tail=("up", 20))
    out["color popover"] = wrap(154, ph + 10, pop)
    fe, fh = fill_editor(0, 0, "en", 1, SKY_STOPS, sel=1, code_lines=["gradient([#3b2a6b,", "  #e0607a, #f2a65a],", "  90deg)"], title="rect · line 26")
    out["fill editor"] = wrap(FE_W + 2, fh + 2, fe)
    ps, ph2 = pen_settings(0, 0, "en", True)
    out["pen settings"] = wrap(322, ph2 + 2, ps)
    out["sidebar rail"] = wrap(14, 200, rail(Lay("rail")).replace(f'y="{MAIN_Y}"', 'y="0"')) if False else None
    return {k: v for k, v in out.items() if v}

# ---------- color popover ----------
CP = {"en": dict(vars="VARIABLES", recent="RECENT", save="Save as variable…", writes="writes", name="name", ok="Save"),
      "es": dict(vars="VARIABLES", recent="RECIENTES", save="Guardar variable…", writes="escribe", name="nombre", ok="Guardar")}
RECENT12 = RECENT + ["#f36c5b", "#72b8e6", "#c8f08a", "#6e4a2a", "#a9d3e4", "#fbf8ef"]

def tail_px(side, x, y, w, h, p):
    tl = []
    for k in range(6):
        span = 11 - 2 * k
        if span <= 0: break
        top = p - span // 2
        if side == "up":
            tl += [R(top, y - k, span, 1, "void", "tail edge")] + ([R(top + 1, y - k, span - 2, 1, "row", "tail fill")] if span > 2 else [])
        elif side == "down":
            tl += [R(top, y + h - 1 + k, span, 1, "void", "tail edge")] + ([R(top + 1, y + h - 1 + k, span - 2, 1, "row", "tail fill")] if span > 2 else [])
        elif side == "left":
            tl += [R(x - k, top, 1, span, "void", "tail edge")] + ([R(x - k, top + 1, 1, span - 2, "row", "tail fill")] if span > 2 else [])
    return G("tail", *tl)

def color_popover(x, y, lang, cur, focus="sun", current="sun", save_open=False, tail=None, target="fill"):
    c = CP[lang]
    w = 150
    ix, iw = x + 5, w - 10
    it = []
    yy = y + 3
    it.append(T(c["vars"], ix, yy, "mute", "label variables")); yy += 14
    g, gh = var_grid(ix + 1, yy, 4, 35, VARS, sel=current, focus=focus, chip_h=10)
    it.append(g)
    yy += gh + 2
    it.append(T(c["recent"], ix, yy, "mute", "label recent")); yy += 14
    it.append(G("recent colors (14)", *[sw(ix + i * 10, yy, 9, col) for i, col in enumerate(RECENT12)])); yy += 14
    it += sv_square(ix, yy, iw, 36, cur, "saturation/value square (pixels)"); yy += 40
    it += hue_strip(ix, yy, iw, 5, cur); yy += 9
    it.append(sfield(ix, yy, 62, cur[:7], "hex"))
    it.append(sfield(ix + 64, yy, 38, "100%", "alpha", "amber"))
    it.append(button(ix + iw - 14, yy + 1, 14, 14, "button eyedropper", icon="dropper"))
    yy += 20
    if save_open:
        bw = tw(c["ok"]) + 8
        it.append(G("save as variable (open)", sfield(ix, yy, iw - bw - 2, "glow", "variable name", focused=True),
                    R(ix + 4 + tw("glow"), yy + 2, 1, 12, "amber", "caret"),
                    button(ix + iw - bw, yy + 1, bw, 14, "button save variable", label=c["ok"])))
        yy += 18
        it.append(T("colors.glow = " + cur[:7], ix, yy, "mute", "writes hint"))
        yy += 14
    else:
        it.append(G("button + save as variable", I("plus", ix, yy + 3, "lizard"), T(c["save"], ix + 10, yy, "lizard")))
        yy += 14
        it.append(G("writes hint", T(c["writes"], ix, yy, "mute"), code("colors." + focus, ix + tw(c["writes"]) + 7, yy, name="ref")))
        yy += 14
    h = yy + 3 - y
    out = chrome(x, y, w, h, name="color popover") + it
    if tail:
        out.append(tail_px(tail[0], x, y, w, h, tail[1]))
    return G(f"color popover ({target})", *out), h

COLORS_ROW = dict(kind="value", d=0, code='let colors = {"sky": #f2a65a, "sea": #2d6a8f, "sun": #fff3c4, "boat": #3b2a6b}', line=6)

def f_pop_row():
    lang = "en"; t = TX[lang]
    lay = Lay("rail")
    rows = rows_layers(extra_top=[COLORS_ROW])
    ys = row_heights(rows)
    sx = 46 + 10 + tw("circle((360, 90), 34, ")
    ring = G("swatch focus (clicked)", frame_lines(sx - 2, ys[10] + 2, 10, 10, "amber", "focus ring"))
    pop, h = color_popover(sx - 20, ys[10] + ROW_H + 6, lang, "#fff3c4", tail=("up", sx + 3), target="code-row swatch, line 10")
    props = props2(lay, lang, f"circle · {t['line']} 10", circ_cr("34"), "click a swatch: color popover")
    return window2([gen.toolbar(None, lang), tabbar2(), code_panel(lay, lang, rows, sel_line=10, props=props), ring,
                    canvas_area2(lay, drawing_g(lay, layer_bodies())), rail(lay), pop, status2(lang, lay, pos=None)])

def f_pop_props():
    lang = "en"; t = TX[lang]
    lay = Lay("rail")
    rows = rows_layers(extra_top=[COLORS_ROW], glow_code="circle((360, 90), 56, #f2a65a66)")
    def pf(y):
        g1, x = gen.point_field(6, y, "center", "360", "90", "center")
        g2 = gen.num_field(x + 4, y, "r", "56", "r")
        g3 = G("field fill (color, focused)", T("fill", 6, y + 22, "mute"),
               field(38, y + 20, tw("#f2a65a66") + 22, 18, [R(41, y + 24, 10, 10, "void", "swatch outline"), R(42, y + 25, 8, 8, "#f2a65a", "swatch"),
                     T("#f2a65a66", 54, y + 23, "text")], "input fill", focused=True))
        return G("fields", g1, g2, g3)
    props = props2(lay, lang, f"circle · {t['line']} 11", pf, None)
    pop, h = color_popover(30, 0, lang, "#f2a65a", save_open=True)
    py = PROPS_Y + 42 - 6 - h
    pop, h = color_popover(30, py, lang, "#f2a65a", focus=None, current=None, save_open=True, tail=("down", 48), target="properties fill field")
    return window2([gen.toolbar(None, lang), tabbar2(unsaved=True), code_panel(lay, lang, rows, sel_line=11, props=props),
                    canvas_area2(lay, drawing_g(lay, layer_bodies())), rail(lay), pop, status2(lang, lay, saved=False, pos=None)])

def f_pop_artist_es():
    lang = "es"
    lay = art_lay()
    b = layer_bodies()
    x0, y0 = lay.cvs((326, 56)); x1, y1 = lay.cvs((394, 124))
    ov = G("selection (circle line 10)", sel_box(x0, y0, x1, y1), gen.anchor(*lay.cvs((360, 90)), "anchor center", "normal"),
           gen.anchor(*lay.cvs((394, 90)), "anchor radius", "alt"))
    def ins(x, y):
        return G("fields", T("relleno", x, y + 1, "mute"), sfield(x + 54, y, tw("colors.sun") + 22, "colors.sun", "fill", swatch="#fff3c4", focused=True),
                 T("centro (360, 90)  r 34", x + 54 + tw("colors.sun") + 30, y + 1, "mute"))
    fx = 18 + tw("círculo · sun") + 6 + 7 + 54
    pop, h = color_popover(fx - 10, MAIN_Y + 25, lang, "#fff3c4", focus="sky", current="sun", tail=("up", fx + 8), target="artist inspector fill")
    return window2([toolbar_artist(lang), tabbar2(), inspector(lang, ins, "circle", "círculo · sun"), tool_column("circle"),
                    canvas_area2(lay, drawing_g(lay, b), ov, splitter=False), right_column(lang, "#fff3c4", art_layers(b)), pop,
                    status2(lang, lay, pos=None)])

FRAMES = [
 ("r1", "R002-1 Shape menu – right-click on canvas, Fill ▸ open", f_menu),
 ("r2", "R002-2 Fill editor – linear gradient on the sky", f_fill_linear),
 ("r3", "R002-3 Fill editor – stepped gradient, 5 bands", f_fill_stepped),
 ("r4", "R002-4 Palette sidebar – editing named color sky", f_palette),
 ("r5", "R002-5 Artist mode – layers.lzdraw", f_artist),
 ("r6", "R002-6 Artist mode – code drawer open", lambda: f_artist(True)),
 ("r7", "R002-7 Blend modes – coder, shadows in Multiply", f_blend_coder),
 ("r8", "R002-8 Blend modes – artist, dropdown open", f_blend_artist),
 ("r9a", "Spanish A – row menu, Relleno ▸", f_spanish_menu),
 ("r9b", "Spanish B – fill editor, radial with hard stop", f_spanish_fill),
 ("r10", "R002-10 Tutorial 1 – Make your first shape", lambda: f_tutorial(0)),
 ("r11", "R002-11 Tutorial 2 – Change its color", lambda: f_tutorial(1)),
 ("r12", "R002-12 Tutorial 3 – First brush stroke", lambda: f_tutorial(2)),
 ("r13", "R002-13 Tutorial 4 – Peek at the code (done)", lambda: f_tutorial(3)),
 ("r14", "R002-14 Tutorial 1 – Spanish", lambda: f_tutorial(0, "es")),
 ("r15", "R002-15 New drawing dialog – object budget", lambda: f_new_drawing("en")),
 ("r16", "R002-16 New drawing dialog – Spanish", lambda: f_new_drawing("es")),
 ("r17", "R002-17 Status bar – object counter (normal, near, over)", f_counter),
 ("r18", "R002-18 Pen settings – coder (pressure curve + code)", f_pen_coder),
 ("r19", "R002-19 Pen settings – artist mode, Spanish", f_pen_artist_es),
 ("r20", "R002-20 Autocomplete – typing ci", lambda: ac_frame("ci")),
 ("r21", "R002-21 Autocomplete – signature hint", lambda: ac_frame("sig")),
 ("r22", "R002-22 Autocomplete – member math.s", lambda: ac_frame("math")),
 ("r23", "R002-23 Color popover – from a code-row swatch", f_pop_row),
 ("r24", "R002-24 Color popover – from the fill field, save as variable", f_pop_props),
 ("r25", "R002-25 Color popover – artist mode, Spanish", f_pop_artist_es),
]

if __name__ == "__main__":
    only = sys.argv[2:] or None
    for k, name, fn in FRAMES:
        if only and k not in only: continue
        s = fn()
        open(os.path.join(OUT, k + ".svg"), "w").write(s)
        print(k, len(s))
    if not only or "comps" in only:
        os.makedirs(os.path.join(OUT, "comps"), exist_ok=True)
        for k, v in comps().items():
            open(os.path.join(OUT, "comps", k.replace("/", "_").replace(" ", "-") + ".svg"), "w").write(v)
        print("comps ok")
