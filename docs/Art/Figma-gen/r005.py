# R005 frames: brush creator (coder + artist, two layouts), brush picker ⋯ menu, dropping and
# importing brush files (.lzbrush / Procreate .brush), and tool / pen tooltips with a mini video.
# Builds on gen.py, r002.py, gen_r002b.py (pen settings pop-up: curve editor, test pad) and r003.py.
# Usage: python3 r005.py <outdir> [frame keys]
import os, sys, math
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import gen, r002, gen_r002b, r003
from gen import R, G, T, I, C, tw, button, frame_lines, svg_doc, W, H, S, MAIN_Y, STATUS_Y, TOP_H, cursor, catmull, nm
from r002 import (Lay, chrome, sfield, slider, seg, pixel_poly, vec, props2, canvas_area2, drawing_g, layer_bodies, rail,
                  status2, tabbar2, window2, code_panel, rows_layers, curve_eval, CURVE, toolbar_artist, inspector,
                  right_column, art_layers, art_lay, menu_box, dither_dim)
from gen_r002b import title_bar, float_chrome, checkbox, PAPER, INK, pen_tip
from r003 import toolbar, tool_column_r3

OUT = sys.argv[1] if len(sys.argv) > 1 else os.path.join(os.path.dirname(os.path.abspath(__file__)), "out")
os.makedirs(OUT, exist_ok=True)

gen.ICONS.update({
 "more": "........|........|........|##.##.##|##.##.##|........|........|........",
 "taper": "........|#.......|###.....|######..|########|######..|###.....|#.......",
 "shape": "...#....|..###...|#######.|.#####..|..###...|.##.##..|.#...#..|........",
 "grain": "#..#..#.|..#...#.|#...#..#|.#.#..#.|#...#...|..#..#.#|.#..#...|#...#..#",
 "dynamics": "........|..##....|.#..#..#|#....##.|........|..##....|.#..#..#|#....##.",
 "render": "........|.##.....|####.##.|########|.##.####|.....##.|........|........",
 "import": "...##...|...##...|...##...|.######.|..####..|...##...|#......#|########",
 "file": ".####...|.#..##..|.#...##.|.#....#.|.#....#.|.#....#.|.#....#.|.######.",
 "play": "........|.##.....|.####...|.######.|.######.|.####...|.##.....|........",
 "still": "........|.##..##.|.##..##.|.##..##.|.##..##.|.##..##.|.##..##.|........",
 "duplicate": "#####...|#...#...|#..#####|#..#...#|####...#|...#...#|...#####|........",
 "image": "########|#......#|#.##...#|#.##...#|#....#.#|#..####.|#.######|########",
})
TAB_ICONS = ["stroke", "taper", "shape", "grain", "dynamics", "graph", "fill", "render"]

# ---------------------------------------------------------------- strings
BC = {
 "en": dict(title="BRUSH CREATOR", tabs=["Stroke", "Taper", "Shape", "Grain", "Dynamics", "Pressure", "Color", "Rendering"],
            reset="Reset", shape="SHAPE", stamp="STAMP", rotation="rotation", rot=["Fixed", "Follow", "Random"], angle="angle",
            scatter="scatter", count="count", roundness="roundness", import_img="Import image…", png_hint="PNG · white paints",
            nib_note="flat nib · follows the stroke", grain="GRAIN", preview="PREVIEW", texture="texture", tex=["Moving", "Fixed"],
            scale="scale", depth="depth", grain_hint="fixed: stays on the paper", with_size="scale with brush size",
            name="name", save="Save brush", export="Export…", sample="sample", pad_note="test only", based="from ink",
            pressure="PRESSURE", presets=["soft", "linear", "firm"], size="size", opacity="opacity", smoothing="smoothing",
            pr="pressure", inp="tablet", inputs=["Trackpad", "Tablet", "Mouse"], floor="floor", floor_hint="light taps still draw",
            raw="raw", eff="out", changed="changed from the base pen"),
 "es": dict(title="CREADOR DE PINCELES", tabs=["Trazo", "Punta", "Forma", "Grano", "Dinámica", "Presión", "Color", "Render"],
            reset="Restaurar", shape="FORMA", stamp="SELLO", rotation="giro", rot=["Fijo", "Sigue", "Azar"], angle="ángulo",
            scatter="dispersión", count="cantidad", roundness="redondez", import_img="Importar imagen…", png_hint="PNG · blanco pinta",
            nib_note="plumilla · sigue el trazo", grain="GRANO", preview="VISTA", texture="textura", tex=["Se mueve", "Fija"],
            scale="escala", depth="fuerza", grain_hint="fija: queda en el papel", with_size="escala con el pincel",
            name="nombre", save="Guardar pincel", export="Exportar…", sample="muestra", pad_note="solo prueba", based="de tinta",
            pressure="PRESIÓN", presets=["suave", "lineal", "firme"], size="grosor", opacity="opacidad", smoothing="suavizado",
            pr="presión", inp="tableta", inputs=["Trackpad", "Tableta", "Ratón"], floor="mínimo", floor_hint="un toque leve pinta",
            raw="entrada", eff="salida", changed="cambiado desde la pluma base"),
}

# ---------------------------------------------------------------- stamped brush renderer (vector, color content)
def hsh(i, j=0, k=0):
    """Deterministic 0..1 noise."""
    v = math.sin(i * 12.9898 + j * 78.233 + k * 37.719) * 43758.5453
    return v - math.floor(v)

def dense(pts, steps=10):
    out = catmull(pts, steps) if len(pts) > 2 else pts
    return out

def stamped(pts, press, color, smin=2, smax=11, spacing=0.1, round_=0.45, angle=30, shape="nib", op=(0.3, 0.55),
            grain=True, taper=0.18, scale=1.0, seed=1):
    """A stamped stroke: many shape stamps along the path + a fixed grain knocked out of it. UI units."""
    d = dense(pts)
    seg_len = [math.dist(d[i], d[i + 1]) for i in range(len(d) - 1)]
    total = sum(seg_len) or 1
    out, cover = [], []
    pos, i, acc, n = 0.0, 0, 0.0, 0
    while i < len(seg_len):
        t = pos / total
        p = press[min(len(press) - 1, int(t * (len(press) - 1)))]
        tp = min(1, max(0.12, min(t, 1 - t) / taper)) if taper else 1
        size = (smin + (smax - smin) * p) * tp * scale
        # point at distance pos
        while i < len(seg_len) and acc + seg_len[i] < pos:
            acc += seg_len[i]; i += 1
        if i >= len(seg_len): break
        f = (pos - acc) / (seg_len[i] or 1)
        x = d[i][0] + (d[i + 1][0] - d[i][0]) * f
        y = d[i][1] + (d[i + 1][1] - d[i][1]) * f
        dirn = math.degrees(math.atan2(d[i + 1][1] - d[i][1], d[i + 1][0] - d[i][0]))
        o = op[0] + (op[1] - op[0]) * p * (0.7 + 0.3 * hsh(n, seed))
        jx, jy = (hsh(n, 2, seed) - 0.5) * 0.5, (hsh(n, 3, seed) - 0.5) * 0.5
        if shape == "nib":
            out.append(f'<ellipse cx="{x+jx:.2f}" cy="{y+jy:.2f}" rx="{size/2:.2f}" ry="{size*round_/2:.2f}" '
                       f'transform="rotate({dirn+angle:.1f} {x+jx:.2f} {y+jy:.2f})" fill-opacity="{o:.2f}"/>')
        else:
            out.append(f'<circle cx="{x+jx:.2f}" cy="{y+jy:.2f}" r="{size/2:.2f}" fill-opacity="{o:.2f}"/>')
        cover.append((x, y, size / 2))
        pos += max(0.35, spacing * size)
        n += 1
    body = f'<g fill="{color}">{"".join(out)}</g>'
    if grain:
        holes = []
        xs = [c[0] for c in cover]; ys = [c[1] for c in cover]
        g = 0.9
        gx0, gy0 = math.floor(min(xs) - 8), math.floor(min(ys) - 8)
        nx, ny = int((max(xs) - gx0 + 8) / g), int((max(ys) - gy0 + 8) / g)
        for a in range(nx):
            for b in range(ny):
                px_, py_ = gx0 + a * g, gy0 + b * g
                # fixed grain: noise of the paper position, not of the stroke
                if hsh(round(px_ / g), round(py_ / g), 7) < 0.55: continue
                near = any((px_ - cx) ** 2 + (py_ - cy) ** 2 < (r + 0.3) ** 2 for cx, cy, r in cover[::2])
                if near:
                    holes.append(f"M{px_:.2f} {py_:.2f}h.55v.55h-.55z")
        body += f'<path d="{"".join(holes)}" fill="{PAPER}" fill-opacity="0.85"/>'
    return body

def wave_pts(x0, y0, w, amp, n=7, phase=0):
    return [(x0 + w * k / (n - 1), y0 + amp * math.sin(k / (n - 1) * math.pi * 2 + phase)) for k in range(n)]

def ramp(n, a=0.15, b=1.0, bell=True):
    return [a + (b - a) * (math.sin(k / (n - 1) * math.pi) if bell else k / (n - 1)) for k in range(n)]

# ---------------------------------------------------------------- stamp + grain tiles
def stamp_svg(kind, cx, cy, s, col="#f2f5ec"):
    """A stamp shape drawn smooth (content), s = box size."""
    r = s / 2
    if kind == "round":
        return f'<circle cx="{cx}" cy="{cy}" r="{r*0.8:.2f}" fill="{col}"/>'
    if kind == "square":
        return f'<rect x="{cx-r*0.7:.2f}" y="{cy-r*0.7:.2f}" width="{r*1.4:.2f}" height="{r*1.4:.2f}" fill="{col}"/>'
    if kind == "nib":
        return f'<ellipse cx="{cx}" cy="{cy}" rx="{r*0.9:.2f}" ry="{r*0.32:.2f}" transform="rotate(-30 {cx} {cy})" fill="{col}"/>'
    if kind == "star":
        p = []
        for k in range(10):
            a = -math.pi / 2 + k * math.pi / 5
            rr = r * (0.85 if k % 2 == 0 else 0.38)
            p.append(f"{cx + rr*math.cos(a):.2f},{cy + rr*math.sin(a):.2f}")
        return f'<polygon points="{" ".join(p)}" fill="{col}"/>'
    if kind == "spray":
        return "".join(f'<circle cx="{cx + (hsh(k,1)-0.5)*r*1.6:.2f}" cy="{cy + (hsh(k,2)-0.5)*r*1.6:.2f}" r="{0.25 + hsh(k,3)*r*0.12:.2f}" fill="{col}"/>'
                       for k in range(26))
    if kind == "leaf":
        return (f'<path d="M{cx-r*0.85:.2f} {cy+r*0.6:.2f} Q{cx-r*0.6:.2f} {cy-r*0.8:.2f} {cx+r*0.85:.2f} {cy-r*0.6:.2f} '
                f'Q{cx+r*0.6:.2f} {cy+r*0.8:.2f} {cx-r*0.85:.2f} {cy+r*0.6:.2f}Z" fill="{col}"/>')
    if kind == "bitmap":  # an imported grainy chalk blob (Shape.png)
        dots = []
        for a in range(-10, 11):
            for b in range(-10, 11):
                x, y = a / 10 * r, b / 10 * r
                if (x / (r * 0.9)) ** 2 + (y / (r * 0.6)) ** 2 < 1 and hsh(a, b, 4) > 0.3:
                    dots.append(f"M{cx + x:.2f} {cy + y:.2f}h{r/10:.2f}v{r/10:.2f}h-{r/10:.2f}z")
        return f'<path d="{"".join(dots)}" fill="{col}" fill-opacity="0.9"/>'
    return ""

def grain_svg(kind, x, y, w, h, col="#f2f5ec"):
    """A grain texture tile (content), light = paints."""
    cell = {"noise": 1.0, "paper": 1.5, "canvas": 2.0, "chalk": 1.2, "hatch": 1.0, "bitmap": 1.3}.get(kind, 1)
    d = []
    for a in range(int(w / cell)):
        for b in range(int(h / cell)):
            px_, py_ = x + a * cell, y + b * cell
            v = hsh(a, b, {"noise": 1, "paper": 2, "chalk": 3, "bitmap": 5}.get(kind, 9))
            if kind == "canvas": v = 0.9 if (a % 3 == 0) != (b % 3 == 0) else 0.2
            if kind == "hatch": v = 0.9 if (a + b) % 4 == 0 else 0.1
            if kind == "paper": v = v * 0.5 + 0.5 * hsh(a // 3, b // 3, 2)
            if kind == "chalk": v = 1 - (v * 0.6 + 0.4 * hsh(a // 2, b // 4, 8))
            if v > 0.45:
                d.append(f"M{px_:.2f} {py_:.2f}h{cell:.2f}v{cell:.2f}h-{cell:.2f}z")
    return f'<rect x="{x}" y="{y}" width="{w}" height="{h}" fill="#5f685b"/><path d="{"".join(d)}" fill="{col}"/>'

def tile(x, y, kind, sel=False, name=None, grain=False, s=24):
    """A shape / grain preset tile: coal well, content drawn smooth, selected = lizard outline."""
    it = [frame_lines(x, y, s, s, "lizard" if sel else "void", "tile outline" + (" (selected)" if sel else "")),
          R(x + 1, y + 1, s - 2, s - 2, "coal", "tile well")]
    if kind == "import":
        it = [frame_lines(x, y, s, s, "void", "tile outline"), R(x + 1, y + 1, s - 2, s - 2, "panel", "tile face")]
        it += [R(x + 1 + k, y + 1, 1, 1, "bevel", "dash") for k in range(0, s - 2, 3)]
        it += [R(x + 1 + k, y + s - 2, 1, 1, "bevel", "dash") for k in range(0, s - 2, 3)]
        it += [R(x + 1, y + 1 + k, 1, 1, "bevel", "dash") for k in range(0, s - 2, 3)]
        it += [R(x + s - 2, y + 1 + k, 1, 1, "bevel", "dash") for k in range(0, s - 2, 3)]
        it.append(I("plus", x + (s - 8) // 2, y + (s - 8) // 2, "mute"))
    elif kind == "none":
        it.append(pixel_poly([(x + 4, y + s - 5), (x + s - 5, y + 4)], "dim", "none slash"))
    elif grain:
        it.append(vec(f"grain {kind} (content)", x + 2, y + 2, s - 4, s - 4, grain_svg(kind, 0, 0, s - 4, s - 4)))
    else:
        it.append(vec(f"stamp {kind} (content)", x + 1, y + 1, s - 2, s - 2, stamp_svg(kind, (s - 2) / 2, (s - 2) / 2, s - 6)))
    if sel:
        it.append(R(x + s - 6, y + s - 6, 5, 5, "lizard", "selected corner"))
        it.append(I("check", x + s - 8, y + s - 8, "void") if False else "")
    return G(name or f"tile {kind}" + (" (selected)" if sel else ""), *it)

def big_preview(x, y, s, body_fn, name):
    return G(name, frame_lines(x - 1, y - 1, s + 2, s + 2, "void", "preview frame"), R(x, y, s, s, "coal", "preview well"),
             vec(name + " (content)", x, y, s, s, body_fn(s)))

# ---------------------------------------------------------------- shared bits
def label_row(x, y, w, label, ctrl_fn, val, lw, name, val_col="amber"):
    vw = tw("100%") + 4
    return G(name, T(label, x, y, "mute"), ctrl_fn(x + lw, y - 1, w - lw - vw - 2), T(val, x + w - tw(val), y, val_col))

def sl_row(x, y, w, label, a, b, val, lw, name):
    return label_row(x, y, w, label, lambda sx, sy, sw_: slider(sx, sy, sw_, a, b, name=name + " slider"), val, lw, "setting " + name)

def primary(x, y, label, icon="plus"):
    w = tw(label) + 22
    return G("button primary " + label, frame_lines(x, y, w, 18, "void"), R(x + 1, y + 1, w - 2, 16, "moss", "face"),
             gen.bevel(x + 1, y + 1, w - 2, 16, "lizard", "void"), I(icon, x + 6, y + 5, "lizard"), T(label, x + 16, y + 2, "bright")), w

def plain_btn(x, y, label, icon, h=18, name=None):
    w = tw(label) + 22
    return button(x, y, w, h, name or ("button " + label), icon=icon, label=label), w

def curve_editor(x, y, gs=84, live=0.62, curve=CURVE, labels=("raw", "out"), name="pressure curve editor (shared with pen settings)"):
    """The R002 pen pop-up curve editor, same construction."""
    gx, gy = x + 1, y + 1
    g = [frame_lines(x, y, gs + 2, gs + 2, "void", "graph outline"), R(gx, gy, gs, gs, "coal", "graph bg")]
    for k in (1, 2, 3):
        g.append(G(f"grid {k}", *[R(gx + i, gy + gs * k // 4, 1, 1, "bevel", "dot") for i in range(0, gs, 3)],
                   *[R(gx + gs * k // 4, gy + i, 1, 1, "bevel", "dot") for i in range(0, gs, 3)]))
    g.append(pixel_poly([(gx, gy + gs - 1), (gx + gs - 1, gy)], "dim", "linear reference"))
    dn = catmull(curve, 24)
    g.append(pixel_poly([(gx + px * (gs - 1), gy + (1 - py) * (gs - 1)) for px, py in dn], "lizard", "pressure curve"))
    if live is not None:
        lx = gx + int(live * (gs - 1)); ly = gy + int((1 - curve_eval(curve, live)) * (gs - 1))
        g.append(G(f"live pressure {live}", *[R(lx, gy + j, 1, 2, "sky", "dash") for j in range(0, gs, 4)],
                   R(lx - 2, ly - 2, 5, 5, "void", "live dot outline"), R(lx - 1, ly - 1, 3, 3, "sky", "live dot")))
    for i, (px, py) in enumerate(curve):
        ax, ay = gx + int(px * (gs - 1)), gy + int((1 - py) * (gs - 1))
        g.append(G(f"curve point {i+1}", R(ax - 3, ay - 3, 7, 7, "void", "outline"), R(ax - 2, ay - 2, 5, 5, "lizard", "fill")))
    g.append(T(labels[0], gx + gs - tw(labels[0]) - 2, gy + gs - 14, "mute", "x axis label"))
    g.append(T(labels[1], gx + 2, gy + 1, "mute", "y axis label"))
    return G(name, *g)

# ---------------------------------------------------------------- tab contents
def shape_tab(x0, y0, w, lang, compact=False):
    t = BC[lang]
    it = [T(t["shape"], x0, y0, "mute", "label shape")]
    kinds = [("round", False), ("square", False), ("nib", True), ("star", False), ("spray", False), ("leaf", False), ("bitmap", False), ("import", False)]
    for k, (kd, sel) in enumerate(kinds):
        it.append(tile(x0 + (k % 4) * 27, y0 + 15 + (k // 4) * 27, kd, sel,
                       name={"bitmap": "tile chalk.png (imported image)", "import": "tile + import image"}.get(kd)))
    by = y0 + 15 + 54 + 3
    b, bw = plain_btn(x0, by, t["import_img"], "image", h=16, name="button import image")
    it.append(b)
    if not compact:
        it.append(T(t["png_hint"], x0, by + 18, "mute", "hint png"))
    # stamp preview
    px = x0 + 117
    pv = 48 if compact else 64
    it.append(T(t["stamp"], px, y0, "mute", "label stamp"))
    it.append(big_preview(px + 1, y0 + 16, pv, lambda s: stamp_svg("nib", s / 2, s / 2, s * 0.8) +
                          f'<line x1="{s/2}" y1="{s/2}" x2="{s/2 + s*0.42*math.cos(math.radians(-30)):.1f}" y2="{s/2 + s*0.42*math.sin(math.radians(-30)):.1f}" stroke="#4fc46a" stroke-width="0.6"/>',
                          "stamp preview (flat nib, 30°)"))
    if compact:
        sx, sy, sw_ = x0, y0 + 94, w
    else:
        sx, sy, sw_ = x0 + 192, y0, w - 192
    lw = max(tw(t[k]) for k in ("rotation", "angle", "scatter", "count", "roundness")) + 6
    sg, _ = seg(sx + lw, sy - 1, t["rot"], 1, "rotation (fixed / follow stroke / random)")
    it.append(G("setting rotation", T(t["rotation"], sx, sy, "mute"), sg))
    rows = [("angle", 30 / 180, None, "30°"), ("scatter", 0.0, None, "0%"), ("count", 1 / 8, None, "1"), ("roundness", 0.35, None, "35%")]
    if compact:
        rows = [rows[0], rows[3]]
    yy = sy + 18
    for key, a, b_, val in rows:
        it.append(sl_row(sx, yy, sw_, t[key], a, b_, val, lw, key))
        yy += 18
    if not compact:
        it.append(R(sx, yy + 2, sw_, 1, "bevel", "rule"))
        it.append(T(t["nib_note"], sx, yy + 6, "mute", "shape note"))
    return G("tab content: Shape", *it)

def grain_tab(x0, y0, w, lang):
    t = BC[lang]
    it = [T(t["grain"], x0, y0, "mute", "label grain")]
    kinds = [("none", False), ("noise", False), ("paper", False), ("canvas", False), ("chalk", True), ("hatch", False), ("bitmap", False), ("import", False)]
    for k, (kd, sel) in enumerate(kinds):
        it.append(tile(x0 + (k % 4) * 27, y0 + 15 + (k // 4) * 27, kd, sel, grain=True,
                       name={"bitmap": "tile paper.png (imported image)", "import": "tile + import image", "none": "tile no grain"}.get(kd)))
    by = y0 + 15 + 54 + 3
    b, bw = plain_btn(x0, by, t["import_img"], "image", h=16, name="button import image")
    it.append(b)
    it.append(T(t["png_hint"], x0, by + 18, "mute", "hint png"))
    px = x0 + 117
    it.append(T(t["preview"], px, y0, "mute", "label preview"))
    def body(s):
        return (f'<rect width="{s}" height="{s}" fill="{PAPER}"/>' +
                stamped(wave_pts(6, s / 2, s - 12, 10, 6), ramp(6, 0.5, 1), INK, 6, 22, 0.1, 0.6, 20, "nib", (0.3, 0.5), True, 0.2))
    it.append(big_preview(px + 1, y0 + 16, 64, body, "grain preview (chalk in a stroke, on paper)"))
    sx, sy, sw_ = x0 + 192, y0, w - 192
    lw = max(tw(t[k]) for k in ("texture", "scale", "depth")) + 6
    sg, _ = seg(sx + lw, sy - 1, t["tex"], 1, "texture (moving / fixed)")
    it.append(G("setting texture", T(t["texture"], sx, sy, "mute"), sg))
    it.append(sl_row(sx, sy + 18, sw_, t["scale"], 0.4, None, "120%", lw, "scale"))
    it.append(sl_row(sx, sy + 36, sw_, t["depth"], 0.7, None, "70%", lw, "depth"))
    it.append(checkbox(sx, sy + 54, False, t["with_size"], "toggle grain scales with brush size"))
    it.append(R(sx, sy + 76, sw_, 1, "bevel", "rule"))
    it.append(T(t["grain_hint"], sx, sy + 80, "mute", "grain note"))
    return G("tab content: Grain", *it)

def pressure_tab(x0, y0, w, lang):
    t = BC[lang]
    it = [T(t["pressure"], x0, y0, "mute", "label pressure"), curve_editor(x0, y0 + 15, labels=(t["raw"], t["eff"]))]
    it.append(sfield(x0, y0 + 103, 86, t["presets"][0], "curve preset", chev=True))
    sx, sw_ = x0 + 100, w - 100
    lw = max(tw(t[k]) for k in ("size", "opacity", "smoothing")) + 6
    yy = y0
    it.append(sl_row(sx, yy, sw_, t["size"], 2 / 24, 14 / 24, "2-14", lw, "size (range)")); yy += 18
    it.append(sl_row(sx, yy, sw_, t["opacity"], 0.3, 0.9, "30-90%", lw, "opacity (range)")); yy += 18
    it.append(sl_row(sx, yy, sw_, t["smoothing"], 0.3, None, "30%", lw, "smoothing")); yy += 20
    it.append(R(sx, yy, sw_, 1, "bevel", "rule")); yy += 4
    it.append(G("live pressure readout", T(t["pr"], sx, yy, "mute"), T("0.62", sx + tw(t["pr"]) + 7, yy, "sky", "value"),
                T("· " + t["inp"], sx + tw(t["pr"]) + 7 + tw("0.62") + 7, yy, "mute"))); yy += 18
    sg, _ = seg(sx, yy, t["inputs"], 1, "input switch (overrides for)")
    it.append(sg); yy += 17
    fl_x = sx + 5 + tw(t["floor"]) + 6
    fl_w = sw_ - (fl_x - sx) - tw("10%") - 6
    it.append(G("tablet overrides", R(sx, yy, sw_, 31, "coal", "override bg"), R(sx, yy, 1, 31, "lizard", "override bar"),
                T(t["floor"], sx + 5, yy + 1, "mute"), slider(fl_x, yy, fl_w, 0.1, name="floor slider"),
                T("10%", sx + sw_ - 2 - tw("10%"), yy + 1, "amber"), T(t["floor_hint"], sx + 5, yy + 15, "mute")))
    return G("tab content: Pressure (shared curve editor)", *it)

TAB_FN = {"Shape": shape_tab, "Grain": grain_tab, "Pressure": pressure_tab}
CHANGED = {"Shape", "Grain", "Rendering", "Pressure"}

# ---------------------------------------------------------------- scribble pad
def scribble_pad(x, y, w, h, lang, dark=False, sample=True, tall=False):
    t = BC[lang]
    paper = "#1a1d19" if dark else PAPER
    it = [frame_lines(x, y, w, h, "void", "pad outline"), R(x + 1, y + 1, w - 2, h - 2, paper, "pad paper"),
          R(x + 1, y + 1, w - 2, 1, "#e4dfd0", "pad inner top shade")]
    iw, ih = w - 2, h - 2
    body = ""
    if tall:
        if sample:
            body += stamped(wave_pts(12, 26, iw - 40, 9, 7), ramp(7, 0.1, 1.0, False), INK, 2, 11, seed=1)
        body += stamped([(14, 70), (50, 56), (80, 78), (112, 58), (150, 76)], ramp(5, 0.3, 1), INK, 2, 11, seed=2)
        body += stamped([(20, 108), (60, 96), (40, 116), (100, 100), (82, 118), (160, 104)], ramp(6, 0.2, 0.7), INK, 2, 11, seed=3)
    else:
        if sample:
            body += stamped(wave_pts(12, ih / 2, 170, 16, 7), ramp(7, 0.1, 1.0, False), INK, 2, 11, seed=1)
        body += stamped([(210, 48), (240, 22), (262, 50), (292, 20), (318, 46)], ramp(5, 0.2, 1), INK, 2, 11, seed=2)
        s3 = [(350 + i * 1.1 + 9 * math.cos(i / 5), 34 + 12 * math.sin(i / 5)) for i in range(0, 110, 6)]
        body += stamped(s3, ramp(len(s3), 0.15, 0.6), INK, 2, 11, seed=3)
    it.append(vec("test strokes, stamped chalk (content, never saved to the drawing)", x + 1, y + 1, iw, ih, body))
    if sample:
        it.append(T(t["sample"], x + 5, y + h - 15, "mute", "label sample stroke (redrawn on every change)"))
    it.append(T(t["pad_note"], x + w - 5 - tw(t["pad_note"]), y + h - 15, "dim", "pad note: strokes here never reach the drawing"))
    # pad tool buttons, top-right
    bx = x + w - 15
    it.append(button(bx, y + 2, 13, 13, "button clear pad (trash)", icon="trash"))
    it.append(button(bx - 14, y + 2, 13, 13, "button paper light / dark", icon="fill"))
    it.append(button(bx - 28, y + 2, 13, 13, "button sample stroke (on)", icon="stroke", state="active"))
    return G("scribble pad (live)", *it)

# ---------------------------------------------------------------- the creator window
def code_block(x, y, w, lines, changed, label):
    ch = len(lines) * 13 + 6
    cb = [frame_lines(x, y, w, ch, "void", "inset outline"), R(x + 1, y + 1, w - 2, ch - 2, "coal", "inset"),
          R(x + 1, y + 3 + changed * 13, w - 8, 13, "row", "changed line"), R(x + 1, y + 3 + changed * 13, 2, 13, "amber", "changed bar")]
    cb += [r002.code(ln, x + 5, y + 3 + i * 13, name=f"code line {i+1}") for i, ln in enumerate(lines)]
    cb += [R(x + w - 5, y + 2, 3, ch - 4, "coal", "scroll track"), R(x + w - 5, y + 2, 3, (ch - 4) // 2, "bevel", "scroll thumb (lines 1-4 of 8)")]
    cb.append(T(label, x + w - 10 - tw(label), y + 3, "dim", "code label"))
    return G("code block (the pen hash it writes, scrolled to what changed)", *cb), ch

# flat keys, as lib/draw/stamps.liz (B060) reads them; the tabs group them only in the UI
CODE_SHAPE = ['pens.chalk = {"render": "stamped",',
              '  "spacing": 0.1, "shape": "flat", "rotation": "follow",',
              '  "angle": 30, "scatter": 0, "count": 1, "roundness": 0.35,',
              '  "grain": {"kind": "chalk", "depth": 0.7, "moving": false},']

def creator(x, y, lang, tab="Shape", layout="A", coder=True, w=500):
    t = BC[lang]
    tabs_en = BC["en"]["tabs"]
    ti = tabs_en.index(tab)
    items = []
    def tl(tx, ty):
        return G("title", T(t["title"], tx, ty, "mute"), T("·", tx + tw(t["title"]) + 5, ty, "mute"),
                 T("chalk", tx + tw(t["title"]) + 17, ty, "bright", "title brush name"),
                 T(t["based"], tx + tw(t["title"]) + 17 + tw("chalk") + 7, ty, "dim", "based on"))
    items.append(title_bar(x, y, w, tl, t["reset"]))
    top = y + 16
    if layout == "A":
        rw = 92
        rail_ = [R(x + 1, top, rw, 147, "coal", "tab rail bg"), R(x + 1 + rw, top, 1, 147, "void", "tab rail line")]
        for i, nm_ in enumerate(t["tabs"]):
            ry = top + 1 + i * 18
            sel = i == ti
            r_ = []
            if sel:
                r_ += [R(x + 1, ry, rw, 18, "moss", "tab bg (selected)"), R(x + 1, ry, 2, 18, "lizard", "selected bar")]
            r_ += [I(TAB_ICONS[i], x + 7, ry + 5, "lizard" if sel else "mute"), T(nm_, x + 19, ry + 3, "bright" if sel else "text")]
            if tabs_en[i] in CHANGED:
                r_.append(R(x + rw - 5, ry + 8, 3, 3, "amber", "changed dot"))
            rail_.append(G(f"tab {tabs_en[i]}" + (" (selected)" if sel else "") + (" · changed" if tabs_en[i] in CHANGED else ""), *r_))
        items.append(G("group tabs (vertical rail)", *rail_))
        cx0, cw = x + rw + 8, w - rw - 14
        items.append(TAB_FN[tab](cx0, top + 3, cw, lang))
        py, ph = top + 151, 70
        items.append(scribble_pad(x + 6, py, w - 12, ph, lang))
        fy = py + ph + 4
    else:
        strip = [R(x + 1, top, w - 2, 15, "coal", "tab strip bg"), R(x + 1, top + 14, w - 2, 1, "void", "tab strip line")]
        tx = x + 4
        for i, nm_ in enumerate(t["tabs"]):
            tw_ = tw(nm_) + 12
            sel = i == ti
            if sel:
                strip.append(G(f"tab {tabs_en[i]} (selected)", R(tx, top + 1, tw_, 14, "panel", "tab face"), R(tx, top + 1, tw_, 1, "lizard", "tab top line"),
                               R(tx - 1, top + 1, 1, 13, "void", "tab left"), R(tx + tw_, top + 1, 1, 13, "void", "tab right"),
                               T(nm_, tx + 6, top + 1, "bright")))
            else:
                strip.append(G(f"tab {tabs_en[i]}" + (" · changed" if tabs_en[i] in CHANGED else ""), T(nm_, tx + 6, top + 1, "mute"),
                               R(tx + tw_ - 3, top + 3, 2, 2, "amber", "changed dot") if tabs_en[i] in CHANGED else ""))
            tx += tw_ + 1
        items.append(G("group tabs (top strip)", *strip))
        cx0, cw = x + 6, 290
        items.append(shape_tab(cx0, top + 19, cw, lang, compact=True) if tab == "Shape" else TAB_FN[tab](cx0, top + 19, cw, lang))
        pw = w - cw - 18
        items.append(R(x + cw + 9, top + 18, 1, 148, "void", "column divider"))
        items.append(scribble_pad(x + w - 6 - pw, top + 18, pw, 148, lang, tall=True))
        fy = top + 170
    # footer: name, export, save
    items.append(G("field name", T(t["name"], x + 6, fy + 2, "mute"), sfield(x + 11 + tw(t["name"]), fy, 120, "chalk", "brush name", h=18)))
    pb, pw_ = primary(0, 0, t["save"], "save")
    sx = x + w - 6 - pw_
    pb, _ = primary(sx, fy, t["save"], "save")
    eb, ew = plain_btn(0, 0, t["export"], "export")
    eb, _ = plain_btn(sx - 4 - ew, fy, t["export"], "export", name="button export (.lzbrush)")
    items += [eb, pb]
    yy = fy + 22
    if coder:
        cb, ch = code_block(x + 6, yy, w - 12, CODE_SHAPE, 2, "pens.chalk · 1-4 / 9")
        items.append(cb); yy += ch + 4
    h = yy + 1 - y
    return G(f"brush creator (floating window, layout {layout})", *float_chrome(x, y, w, h, "brush creator"), *items), h

# ================================================================ frames: brush creator
def bg_coder(lang="en", props_title="brush · pens.chalk"):
    lay = Lay("rail")
    props = props2(lay, lang, props_title, lambda y: T("pen picker menu: New brush…", 6, y + 2, "mute"), None)
    return lay, [toolbar("brush", lang), tabbar2(unsaved=True), code_panel(lay, lang, rows_layers(), props=props),
                 canvas_area2(lay, drawing_g(lay, layer_bodies())), rail(lay)]

def bg_artist(lang):
    lay = art_lay()
    b = layer_bodies()
    pen = "chalk"
    def ins(x, y):
        return G("fields (brush)", T("pen" if lang == "en" else "pluma", x, y + 1, "mute"), sfield(x + 42, y, 100, pen, "pen", chev=True),
                 button(x + 145, y + 1, 14, 14, "button brush creator (open)", icon="settings", state="active"),
                 T("color", x + 166, y + 1, "mute"), sfield(x + 204, y, tw(INK) + 22, INK, "color", swatch=INK))
    return lay, [toolbar_artist(lang), tabbar2(unsaved=True), inspector(lang, ins, "brush", ("brush" if lang == "en" else "pincel") + " · sun"),
                 tool_column_r3("brush", INK), canvas_area2(lay, drawing_g(lay, b), splitter=False),
                 right_column(lang, INK, art_layers(b, sel="sun"))]

def f_creator(lang="en", tab="Shape", layout="A", coder=True):
    if coder:
        lay, parts = bg_coder(lang)
        x, y = 6, TOP_H + 2
    else:
        lay, parts = bg_artist(lang)
        x, y = 6, MAIN_Y + 22
    cr, h = creator(x, y, lang, tab, layout, coder)
    parts.append(cr)
    # pen tip in the pad (drawing a test stroke)
    if layout == "A":
        parts.append(pen_tip(x + 6 + 1 + 318 + 4, y + 16 + 151 + 1 + 46))
    else:
        parts.append(pen_tip(x + 500 - 6 - 192 + 1 + 150, y + 16 + 18 + 1 + 76))
    parts.append(status2(lang, lay, saved=False, pos=None))
    return window2(parts)


# ================================================================ brush picker
BX = 131                     # brush tool button x (toolbar with select + lasso first)
PK_W, PK_RH = 230, 20
PENS = [("ink", "crisp, tapers"), ("pencil", "light, grainy"), ("marker", "see-through"), ("brush", "soft, wide"),
        ("calligraphy", "45° nib"), ("soft ink", "yours"), ("chalk", "stamped")]

def chalk_preview(x, y, w=50, h=16, seed=4, col="#d4dacd"):
    body = stamped(wave_pts(4, h / 2, w - 8, 3.5, 5), ramp(5, 0.3, 1), col, 1, 5, seed=seed, grain=False, op=(0.45, 0.8))
    return vec("stroke preview stamped (content)", x, y, w, h, body)

def picker(px, py, pens=PENS, current="chalk", hover=None, more="normal", new=None, drop=None, warn=None):
    pw, rh = PK_W, PK_RH
    extra = 30 if drop else 0
    hh = len(pens) * rh + 18 + 14 + extra
    items = [R(px + 2, py + 2, pw, hh, "void", "drop shadow"), frame_lines(px, py, pw, hh, "void"),
             R(px + 1, py + 1, pw - 2, hh - 2, "row", "menu bg"), R(px + 1, py + 1, pw - 2, 13, "raised", "menu header"),
             T("PEN · pens.list", px + 5, py, "mute")]
    mx = px + pw - 16
    if more == "normal":
        items.append(I("more", mx + 2, py + 3, "text", "button ⋯ (more)"))
    else:
        items.append(button(mx - 1, py + 1, 14, 12, "button ⋯ (open)", icon="more", state="active"))
    for i, (pn, desc) in enumerate(pens):
        ry = py + 16 + i * rh
        it = []
        cur, hov = pn == current, pn == hover
        if cur:
            it += [R(px + 1, ry, pw - 2, rh, "moss", "item bg (current)"), R(px + 1, ry, 2, rh, "lizard", "current bar")]
        if hov:
            it += [R(px + 1, ry, pw - 2, rh, "raised", "item bg (hover)")]
        it.append(G("preview frame", R(px + 6, ry + 2, 50, 16, "coal", "preview bg"), frame_lines(px + 5, ry + 1, 52, 18, "void")))
        if pn in ("chalk", "soft charcoal"):
            it.append(chalk_preview(px + 6, ry + 2, seed=5 if pn == "chalk" else 9))
        else:
            kind = "ink" if pn == "soft ink" else pn
            it.append(f'<g id="stroke preview (full-res)" transform="scale(0.5)">{gen.pen_preview(kind, (px + 9) * S, (ry + 5) * S, 88, 16, C["text"])}</g>')
        it.append(T(pn, px + 62, ry + 1, "bright" if cur else "text"))
        dx = px + 62 + tw(pn) + 7
        if pn == new:
            it.append(G("new badge", R(dx, ry + 4, tw("new") + 4, 11, "raised", "badge bg"), T("new", dx + 2, ry + 1, "lizard")))
            dx += tw("new") + 8
        elif tw(pn) + tw(desc) + 7 < pw - 66 - 30:
            it.append(T(desc, dx, ry + 1, "mute", f"desc {desc}"))
        gx = px + pw - 18
        if pn == warn:
            it.append(I("warn", gx - 12, ry + 6, "amber", "import warnings (hover: details)"))
        it.append(I("settings", gx + 2, ry + 6, "text" if cur else "dim", "pen settings ⚙"))
        items.append(G(f"pen {pn}" + (" (current)" if cur else "") + (" (hover)" if hov else ""), *it))
    fy = py + 16 + len(pens) * rh
    if drop:
        dz = [R(px + 4, fy + 2, pw - 8, 24, "moss", "drop zone bg")]
        for k in range(0, pw - 8, 4):
            dz += [R(px + 4 + k, fy + 2, 2, 1, "lizard", "dash"), R(px + 4 + k, fy + 25, 2, 1, "lizard", "dash")]
        for k in range(0, 24, 4):
            dz += [R(px + 4, fy + 2 + k, 1, 2, "lizard", "dash"), R(px + pw - 5, fy + 2 + k, 1, 2, "lizard", "dash")]
        dz += [I("import", px + 10, fy + 10, "lizard"), T(drop[0], px + 22, fy + 2, "bright"), T(drop[1], px + 22, fy + 12, "mute") if False else ""]
        dz.append(T(drop[1], px + pw - 8 - tw(drop[1]), fy + 2, "mute"))
        items.append(G("drop zone (file over the picker)", *dz))
        fy += extra
    items.append(G("footer hint", R(px + 1, fy, pw - 2, 1, "bevel", "separator"), T("double-click brush: its settings", px + 5, fy + 1, "mute")))
    return G("pen picker (open)", *items), hh

def coder_bg(lang="en", title="brush · pens.chalk", hint="pick a pen · menu for brush files"):
    lay = Lay("rail")
    props = props2(lay, lang, title, lambda y: T(hint, 6, y + 2, "mute"), None)
    return lay, [toolbar("brush", lang), tabbar2(unsaved=True), code_panel(lay, lang, rows_layers(), props=props),
                 canvas_area2(lay, drawing_g(lay, layer_bodies())), rail(lay)]

def f_picker_menu():
    lay, parts = coder_bg()
    px, py = BX, TOP_H + 1
    pk, hh = picker(px, py, more="open")
    entries = [("plus", "New brush…", ""), ("settings", "Edit in creator…", ""), ("duplicate", "Duplicate", ""), None,
               ("export", "Export brush…", ""), ("export", "Export set…", ""), ("import", "Import…", "")]
    mx = px + PK_W - 17
    mn, mw, mh = menu_box(mx, py + 14, entries, "menu ⋯ (brush files)", title="chalk", hover="Export brush…")
    tip = gen.tooltip_box(mx + mw - 170, py + 14 + mh + 4, [lambda a, b: T("Saves chalk.lzbrush", a, b, "bright"),
                                                               lambda a, b: T("plain lizard, shareable", a, b, "mute")],
                          "tooltip export brush", w=tw("plain lizard, shareable") + 14)
    parts += [pk, mn, tip, cursor(mx + 60, py + 14 + 62 + 5), status2("en", lay, saved=False, pos=None)]
    return window2(parts)

def drag_ghost(x, y, name):
    w = tw(name) + 22
    return G("dragged file (ghost chip)", R(x + 2, y + 2, w, 16, "void", "shadow"), frame_lines(x, y, w, 16, "void"),
             R(x + 1, y + 1, w - 2, 14, "row", "chip bg"), I("file", x + 5, y + 4, "text"), T(name, x + 16, y + 2, "text"),
             G("copy badge", R(x - 6, y - 6, 9, 9, "void", "badge outline"), R(x - 5, y - 5, 7, 7, "lizard", "badge"),
               R(x - 4, y - 2, 5, 1, "void", "plus h"), R(x - 2, y - 4, 1, 5, "void", "plus v")))

def f_drop_picker():
    lay, parts = coder_bg()
    px, py = BX, TOP_H + 1
    pk, hh = picker(px, py, pens=PENS, drop=("Drop to add", "Procreate .brush"))
    fy = py + 16 + len(PENS) * PK_RH
    cx, cy = px + 120, fy + 14
    parts += [button(BX, 1, 14, 14, "tool brush (spring-loaded open)", icon="brush", state="active"), pk,
              cursor(cx, cy), drag_ghost(cx + 10, cy + 12, "Soft Charcoal.brush"), status2("en", lay, saved=False, pos=None)]
    return window2(parts)

def dashed_box(x, y, w, h, col, name):
    it = []
    for k in range(0, w, 6):
        it += [R(x + k, y, min(4, w - k), 2, col, "dash"), R(x + k, y + h - 2, min(4, w - k), 2, col, "dash")]
    for k in range(0, h, 6):
        it += [R(x, y + k, 2, min(4, h - k), col, "dash"), R(x + w - 2, y + k, 2, min(4, h - k), col, "dash")]
    return G(name, *it)

def f_drop_window():
    lay, parts = coder_bg()
    parts.append(f'<rect id="dim (void 62%)" x="0" y="0" width="{W}" height="{H}" fill="{C["void"]}" fill-opacity="0.62"/>')
    w, h = 300, 138
    x, y = (W - w) // 2, 100
    it = list(chrome(x, y, w, h, name="drop card"))
    it.append(dashed_box(x + 6, y + 6, w - 12, h - 12, "lizard", "drop target (dashed lizard, 2 px)"))
    def ctr(s_, yy, col, nm_=None): return T(s_, x + (w - tw(s_)) // 2, yy, col, nm_)
    it += [I("import", x + (w - 8) // 2, y + 18, "lizard", "import icon"),
           ctr("Drop to import a brush", y + 31, "bright"), ctr("Soft Charcoal.brush · Procreate", y + 47, "text"),
           R(x + 40, y + 68, w - 80, 1, "bevel", "rule"),
           ctr(".lzbrush .lzbrushes .brush .brushset", y + 75, "mute"), ctr("best effort for Procreate files", y + 90, "mute"),
           ctr("images still open the trace dialog", y + 105, "dim")]
    parts.append(G("window drop overlay (brush file)", *it))
    cx, cy = 300, 250
    parts += [cursor(cx, cy), drag_ghost(cx + 10, cy + 12, "Soft Charcoal.brush")]
    return window2(parts)

# ---------- import result: dialog option
IM = dict(skipped=["wet mix (dilution, charge)", "tilt -> size"], changed=("grain blend: color burn", "drawn as multiply"),
          taken=["spacing, streamline, jitter,", "taper, size 3-40, opacity,", "shape + rotation, grain,", "pressure curve"])

def soft_charcoal_body(w, h):
    return (f'<rect width="{w}" height="{h}" fill="{PAPER}"/>' +
            stamped(wave_pts(8, h / 2, w - 16, h * 0.22, 7), ramp(7, 0.2, 1, False), "#2b2b2b", 2, 13, 0.08, 0.9, 0, "round",
                    (0.18, 0.4), True, 0.25, seed=11))

def f_import_dialog():
    lay, parts = coder_bg()
    parts.append(f'<rect id="dim (void 62%)" x="0" y="0" width="{W}" height="{H}" fill="{C["void"]}" fill-opacity="0.62"/>')
    w, h = 380, 276
    x, y = (W - w) // 2, 34
    it = [R(x + 3, y + 3, w, h, "void", "drop shadow 3px"), frame_lines(x, y, w, h, "void"), R(x + 1, y + 1, w - 2, h - 2, "panel", "dialog bg"),
          R(x + 1, y + 1, w - 2, 14, "raised", "header bg"), R(x + 1, y + 15, w - 2, 1, "void", "header line"),
          T("IMPORT BRUSH · Soft Charcoal.brush", x + 6, y + 2, "mute"), I("close", x + w - 12, y + 5, "mute")]
    # left: preview pad + shape / grain
    lx, ly = x + 8, y + 22
    it.append(G("preview (sample stroke with the imported brush)", frame_lines(lx - 1, ly - 1, 152, 58, "void", "pad outline"),
                vec("sample stroke (content)", lx, ly, 150, 56, soft_charcoal_body(150, 56))))
    ty = ly + 64
    it.append(tile(lx, ty, "bitmap", name="shape tile (Shape.png)"))
    it.append(G("shape label", T("shape", lx + 30, ty - 1, "text"), T("Shape.png 256", lx + 30, ty + 11, "mute")))
    it.append(tile(lx, ty + 30, "paper", grain=True, name="grain tile (Grain.png)"))
    it.append(G("grain label", T("grain", lx + 30, ty + 29, "text"), T("Grain.png 512", lx + 30, ty + 41, "mute")))
    it.append(G("render chip", R(lx, ty + 64, tw("stamped") + 6, 11, "raised", "chip bg"), T("stamped", lx + 3, ty + 61, "text")))
    # right: name, warnings first, taken summary
    rx, ry = x + 172, y + 20
    r_ = [T("Soft Charcoal", rx, ry, "bright", "brush name"), T("from Procreate · best effort", rx, ry + 14, "mute"),
          R(rx, ry + 31, w - 180, 1, "bevel", "rule")]
    yy = ry + 36
    r_.append(G("skipped", I("warn", rx, yy + 3, "amber"), T(f"SKIPPED {len(IM['skipped'])}", rx + 12, yy, "amber")))
    yy += 15
    for sname in IM["skipped"]:
        r_.append(T(sname, rx + 12, yy, "text", "skipped: " + sname)); yy += 14
    yy += 3
    r_.append(G("changed", I("info", rx, yy + 3, "sky"), T("CHANGED 1", rx + 12, yy, "sky")))
    yy += 15
    r_ += [T(IM["changed"][0], rx + 12, yy, "text"), T(IM["changed"][1], rx + 12, yy + 14, "mute")]
    yy += 31
    r_.append(G("taken", I("check", rx, yy + 3, "lizard"), T("TAKEN 9", rx + 12, yy, "lizard")))
    yy += 15
    for ln in IM["taken"]:
        r_.append(T(ln, rx + 12, yy, "mute")); yy += 14
    it.append(G("import report", *r_))
    # footer
    fy = y + h - 24
    it.append(R(x + 1, fy - 5, w - 2, 1, "void", "footer line"))
    it.append(T("-> brushes/soft_charcoal.lzbrush", x + 8, fy - 20, "dim", "where it is saved"))
    pb, pw_ = primary(0, 0, "Add brush")
    pb, _ = primary(x + w - 8 - pw_, fy, "Add brush")
    ob, ow = plain_btn(0, 0, "Open in creator…", "settings")
    ob, _ = plain_btn(x + w - 8 - pw_ - 4 - ow, fy, "Open in creator…", "settings")
    cw_ = tw("Cancel") + 12
    cb = button(x + 8, fy, cw_, 18, "button cancel", label="Cancel")
    it += [cb, ob, pb]
    parts.append(G("import result dialog", *it))
    parts.append(cursor(x + w - 40, fy + 10))
    return window2(parts)

# ---------- import result: quiet option (row + toast)
def f_import_toast():
    lay, parts = coder_bg(hint="soft charcoal added")
    px, py = BX, TOP_H + 1
    pens = PENS + [("soft charcoal", "")]
    pk, hh = picker(px, py, pens=pens, current="soft charcoal", new="soft charcoal", warn="soft charcoal")
    ry = py + 16 + 7 * PK_RH
    tip = gen.tooltip_box(px + PK_W - 30 - 10, ry + PK_RH + 2, [lambda a, b: T("2 settings skipped", a, b, "bright"),
                                                               lambda a, b: T("wet mix, tilt -> size", a, b, "text"),
                                                               lambda a, b: T("1 changed: color burn", a, b, "text"),
                                                               lambda a, b: T("drawn as multiply", a, b, "mute")],
                          "tooltip import warnings", w=tw("1 changed: color burn") + 14)
    tw_ = 236
    tx, ty = W - 14 - tw_ - 8, STATUS_Y - 6 - 36
    toast = gen.tooltip_box(tx, ty, [lambda a, b: G("l1", I("check", a, b + 3, "lizard"), T("Added soft charcoal", a + 12, b, "bright")),
                                     lambda a, b: G("l2", I("warn", a, b + 3, "amber"), T("3 notes", a + 12, b, "text"),
                                                    T("Details", a + 12 + tw("3 notes") + 10, b, "sky"), T("Undo", a + tw_ - 12 - tw("Undo"), b, "mute"))],
                            "toast imported (bottom right, 6 s)", w=tw_)
    toast = toast.replace('<g id="toast imported (bottom right, 6 s)">', '<g id="toast imported (bottom right, 6 s)">' + R(tx + 1, ty + 1, 2, 34, "lizard", "toast bar"), 1)
    parts += [pk, tip, cursor(px + PK_W - 40, ry + 8), toast, status2("en", lay, saved=False, pos=None)]
    return window2(parts)

# ================================================================ tooltips with a mini video
VW, VH = 168, 94          # the video, UI px (336×188 on screen)
def video_lasso(w, h, t=0.62):
    body = f'<rect width="{w}" height="{h}" fill="#f6d9a8"/>'
    body += f'<circle cx="{w*0.3:.1f}" cy="{h*0.45:.1f}" r="{h*0.17:.1f}" fill="#f2a65a"/>'
    body += f'<rect x="{w*0.5:.1f}" y="{h*0.32:.1f}" width="{h*0.3:.1f}" height="{h*0.3:.1f}" fill="#2d6a8f"/>'
    body += f'<polygon points="{w*0.78:.1f},{h*0.75:.1f} {w*0.9:.1f},{h*0.5:.1f} {w*0.98:.1f},{h*0.78:.1f}" fill="#8a6fb0"/>'
    loop = [(w * 0.15 + w * 0.47 * (0.5 + 0.5 * math.cos(a)) + 2 * math.sin(a * 5), h * 0.45 + h * 0.36 * math.sin(a))
            for a in [math.pi + k * 2 * math.pi / 40 for k in range(41)]]
    n = int(len(loop) * t)
    d = "M" + " L".join(f"{x:.1f} {y:.1f}" for x, y in loop[:n])
    body += f'<path d="{d}" fill="none" stroke="#0b0c0b" stroke-width="1.2"/><path d="{d}" fill="none" stroke="#f2f5ec" stroke-width="1.2" stroke-dasharray="3 2"/>'
    if t >= 1:
        for (bx, by, bw, bh) in ((w * 0.3 - h * 0.17, h * 0.28, h * 0.34, h * 0.34), (w * 0.5, h * 0.32, h * 0.3, h * 0.3)):
            body += f'<rect x="{bx-2:.1f}" y="{by-2:.1f}" width="{bw+4:.1f}" height="{bh+4:.1f}" fill="none" stroke="#4fc46a" stroke-width="1" stroke-dasharray="3 2"/>'
    else:
        ex, ey = loop[n - 1]
        body += f'<path d="M{ex:.1f} {ey:.1f} l0 9 l2.5 -2.5 l3.5 0 z" fill="#f2f5ec" stroke="#0b0c0b" stroke-width="0.8"/>'
    return body

def video_calli(w, h, t=0.7):
    body = f'<rect width="{w}" height="{h}" fill="{PAPER}"/>'
    pts = [(14 + (w - 28) * k / 60, h * 0.55 - h * 0.25 * math.sin(k / 60 * math.pi * 2.2) + 6 * math.sin(k / 4)) for k in range(61)]
    n = max(2, int(len(pts) * t))
    ws = []
    for i in range(n):
        a = math.atan2(pts[min(i + 1, 60)][1] - pts[i][1], pts[min(i + 1, 60)][0] - pts[i][0])
        ws.append(1 + 7 * abs(math.sin(a - math.radians(45))))
    body += f'<path d="{gen.variable_outline(pts[:n], ws)}" fill="{INK}"/>'
    if t < 1:
        ex, ey = pts[n - 1]
        body += f'<circle cx="{ex:.1f}" cy="{ey:.1f}" r="4" fill="none" stroke="#0b0c0b" stroke-width="1"/><circle cx="{ex:.1f}" cy="{ey:.1f}" r="0.8" fill="#0b0c0b"/>'
    return body

def video_box(x, y, w, h, body_fn, state, t, name):
    it = [frame_lines(x - 1, y - 1, w + 2, h + 2, "void", "video frame")]
    tt = 1.0 if state == "still" else t
    it.append(vec(f"video frame (content, {state})", x, y, w, h, body_fn(w, h, tt if state != "loading" else 0.02)))
    if state == "loading":
        it.append(f'<rect id="poster dim (void 62%)" x="{x}" y="{y}" width="{w}" height="{h}" fill="{C["void"]}" fill-opacity="0.62"/>')
        cx, cy = x + w // 2, y + h // 2
        it.append(G("loading dots (3, the lit one walks)", R(cx - 7, cy - 1, 3, 3, "lizard"), R(cx - 1, cy - 1, 3, 3, "mute"), R(cx + 5, cy - 1, 3, 3, "mute")))
    if state == "still":
        cw_ = tw("motion off") + 16
        it.append(G("still chip", R(x + 3, y + 3, cw_, 13, "void", "chip outline"), R(x + 4, y + 4, cw_ - 2, 11, "raised", "chip bg"),
                    I("still", x + 5, y + 5, "mute"), T("motion off", x + 15, y + 2, "text")))
    # progress line under the video
    if state != "still":
        it.append(R(x, y + h + 2, w, 1, "bevel", "progress track"))
        if state == "playing":
            it.append(R(x, y + h + 2, int(w * t), 1, "lizard", f"progress {int(t*100)}%"))
    return G(name, *it)

def video_tip(x, y, kind, state="playing", layout="A", lang="en"):
    if kind == "lasso":
        title, key, desc, body = "Lasso", "L", ["Draw a loop to select", "what's inside it."], video_lasso
    else:
        title, key, desc, body = "calligraphy", "45° nib", ["Width follows the angle", "of the stroke."], video_calli
    if layout == "A":
        w = VW + 12
        h = 6 + VH + 5 + 3 * 15 + 3
        it = [R(x + 2, y + 2, w, h, "void", "drop shadow"), frame_lines(x, y, w, h, "void"), R(x + 1, y + 1, w - 2, h - 2, "row", "tooltip bg"),
              R(x + 1, y + 1, w - 2, 1, "bevel", "top highlight")]
        it.append(video_box(x + 6, y + 6, VW, VH, body, state, 0.62, "mini video (loops)"))
        ty = y + 6 + VH + 5
        tx = x + 6
    else:
        vw, vh = 96, 54
        w = 6 + vw + 8 + 168 + 6
        h = 6 + vh + 3 + 6
        it = [R(x + 2, y + 2, w, h, "void", "drop shadow"), frame_lines(x, y, w, h, "void"), R(x + 1, y + 1, w - 2, h - 2, "row", "tooltip bg"),
              R(x + 1, y + 1, w - 2, 1, "bevel", "top highlight")]
        it.append(video_box(x + 6, y + 6, vw, vh, body, state, 0.62, "mini video (loops)"))
        ty = y + 4
        tx = x + 6 + vw + 8
        w_txt = 168
    it.append(T(title, tx, ty, "bright", "title"))
    kx = (x + w - 6 - tw(key) - 4) if layout == "A" else (tx + w_txt - tw(key) - 4)
    it.append(G("shortcut / note badge", R(kx, ty + 1, tw(key) + 4, 11, "raised", "badge bg"), T(key, kx + 2, ty - 2, "mute")))
    for i, ln in enumerate(desc):
        it.append(T(ln, tx, ty + 15 * (i + 1), "text", "description"))
    return G(f"tooltip with mini video · {title} ({state}, layout {layout})", *it), w, h

def f_tool_video():
    lay, parts = coder_bg(title="lasso", hint="L: lasso · loop to select")
    el, w, h = video_tip(18, TOP_H + 3, "lasso", "playing")
    parts += [button(18, 1, 14, 14, "tool lasso (hover)", icon="lasso", state="hover"), el, cursor(18 + 7, 8), status2("en", lay, saved=True, pos=None)]
    return window2(parts)

def f_pen_video():
    lay, parts = coder_bg()
    px, py = BX, TOP_H + 1
    pk, hh = picker(px, py, hover="calligraphy", current="ink")
    ry = py + 16 + 4 * PK_RH
    el, w, h = video_tip(px, py + hh + 4, "calli", "playing", "B")
    parts += [pk, el, cursor(px + 90, ry + 10), status2("en", lay, saved=False, pos=None)]
    return window2(parts)

def f_video_states():
    """Sheet: layouts A (video on top) and B (video beside) × loading / playing / reduced motion."""
    cols = [("loading", "loading"), ("playing", "playing"), ("still", "reduced motion")]
    items = [R(0, 0, 900, 320, "coal", "sheet bg")]
    items.append(T("A · VIDEO ON TOP", 20, 6, "lizard", "row label A"))
    items.append(T("B · VIDEO BESIDE", 20, 206, "lizard", "row label B"))
    for ci, (st, lab) in enumerate(cols):
        x = 20 + ci * 290
        items.append(T(lab, x + 140, 6, "mute", "column " + lab))
        el, w, h = video_tip(x, 24, "lasso", st, "A")
        items.append(el)
        items.append(T(lab, x + 140, 206, "mute", "column " + lab))
        el, w, h = video_tip(x, 224, "lasso", st, "B")
        items.append(el)
    return svg_doc(900 * S, 310 * S, G("ui (video tooltip states)", *items), bg="coal", scale=S)

FRAMES = {
 "c1": ("R005-1 Brush creator – coder, Shape tab (layout A: side tabs, wide pad)", lambda: f_creator("en", "Shape", "A")),
 "c2": ("R005-2 Brush creator – coder, Shape tab (layout B: top tabs, tall pad)", lambda: f_creator("en", "Shape", "B")),
 "c3": ("R005-3 Brush creator – artist, Grain tab", lambda: f_creator("en", "Grain", "A", coder=False)),
 "c4": ("R005-4 Brush creator – artist, Spanish, Pressure tab (shared curve editor)", lambda: f_creator("es", "Pressure", "A", coder=False)),
 "p5": ("R005-5 Pen picker – ⋯ menu (New, Edit, Duplicate, Export brush, Export set, Import)", f_picker_menu),
 "p6": ("R005-6 Drop a brush file on the picker (option: picker drop zone)", f_drop_picker),
 "p7": ("R005-7 Drop a brush file on the window (option: whole-window target)", f_drop_window),
 "p8": ("R005-8 Import result – dialog with skipped / changed / taken (option: dialog)", f_import_dialog),
 "p9": ("R005-9 Import result – new row + toast (option: quiet)", f_import_toast),
 "v10": ("R005-10 Tool tooltip with mini video – lasso, playing", f_tool_video),
 "v11": ("R005-11 Pen tooltip with mini video – calligraphy, playing", f_pen_video),
 "v12": ("R005-12 Video tooltip states – layout A / B × loading, playing, reduced motion", f_video_states),
}

if __name__ == "__main__":
    only = sys.argv[2:] or list(FRAMES)
    for k in only:
        s = FRAMES[k][1]()
        open(os.path.join(OUT, k + ".svg"), "w").write(s)
        print(k, len(s))
