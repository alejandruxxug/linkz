# Home screen frames (R001 8-10): gallery EN, empty state, gallery ES.
import os, sys, re, math
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import gen
from gen import R, G, T, I, C, tw, button, field, frame_lines, svg_doc, W, H, S, bezier_d, catmull, variable_outline, pressure_widths

gen.ICONS.update({
 "search": ".###....|#...#...|#...#...|#...#...|.###....|....#...|.....#..|......##",
 "settings": "...##...|.#.##.#.|..####..|###..###|###..###|..####..|.#.##.#.|...##...",
 "more": "........|........|........|##.##.##|##.##.##|........|........|........",
 "folder": "........|###.....|#..####.|#......#|#......#|#......#|########|........",
 "trash": "..####..|########|.#....#.|.#.#.##.|.#.#.##.|.#.#.##.|.#....#.|..####..",
 "copy": "####....|#..#....|#.#####.|#.#...#.|###...#.|..#...#.|..#####.|........",
 "rename": "........|######..|.....#..|.....#..|.##..#..|.....#..|######..|........",
 "open": "........|.#####..|.#...#..|.#....##|.#...#.#|.#####..|........|........",
 "sort": "...#....|..###...|.#.#.#..|...#....|...#....|...#....|...#....|........",
})

T_ = {
 "en": {"search": "Search drawings", "working": "WORKING ON", "all": "ALL DRAWINGS", "new": "New drawing", "sort": "Last edited",
        "count": "12 drawings · ~/Linkz", "menu": ["Open", "Rename", "Duplicate", "Show in Finder", "Delete"], "lang": "EN",
        "big": {"2h": "edited 2 h ago", "y": "edited yesterday", "mon": "edited Monday"},
        "ago": {"2h": "2 h ago", "y": "yesterday", "mon": "Monday", "1w": "last week", "2w": "2 wk ago", "4d": "4 days ago", "5d": "5 days ago", "3w": "3 wk ago", "1m": "1 mo ago", "2m": "2 mo ago"},
        "empty_title": "No drawings yet", "empty_sub": "Every drawing is a lizard-script file you can read.", "example": "Open an example",
        "empty_hint": "or drop a .lzdraw file here", "first": "first launch"},
 "es": {"search": "Buscar dibujos", "working": "TRABAJANDO EN", "all": "TODOS LOS DIBUJOS", "new": "Nuevo dibujo", "sort": "Última edición",
        "count": "12 dibujos · ~/Linkz", "menu": ["Abrir", "Renombrar", "Duplicar", "Mostrar en el Finder", "Eliminar"], "lang": "ES",
        "big": {"2h": "editado hace 2 h", "y": "editado ayer", "mon": "editado el lunes"},
        "ago": {"2h": "hace 2 h", "y": "ayer", "mon": "el lunes", "1w": "hace 1 sem.", "2w": "hace 2 sem.", "4d": "hace 4 días", "5d": "hace 5 días", "3w": "hace 3 sem.", "1m": "hace 1 mes", "2m": "hace 2 meses"},
        "empty_title": "Aún no hay dibujos", "empty_sub": "Cada dibujo es un archivo de lizard-script.", "example": "Abrir un ejemplo",
        "empty_hint": "o suelta aquí un archivo .lzdraw", "first": "primer inicio"},
}

# ---------- thumbnails: (size, svg body in drawing coordinates) ----------
def th_sunset():
    return (480, 320), gen.drawing().split('scale(0.375)">', 1)[1].rsplit("</g>", 1)[0]

def th_hello():
    b = ['<rect width="400" height="200" fill="#f4f1e8"/>', '<line x1="20" y1="100" x2="380" y2="100" stroke="#1d1d1f" stroke-width="2"/>']
    for i in range(6):
        b.append(f'<circle cx="{40+i*64}" cy="100" r="{6+i*4}" fill="{"#e0a526" if i % 2 else "#3fa34d"}"/>')
    return (400, 200), "".join(b)

def th_strokes():
    b = ['<rect width="480" height="400" fill="#fbf8ef"/>']
    cols = ["#1d1d1f", "#4a4a4a", "#2d6a8f", "#c0392b", "#3b2a6b"]
    pr = [0.15, 0.8, 0.5, 1.0, 0.3, 0.7]
    for i in range(5):
        y = 50 + i * 54
        wave = [(40, y + 14), (110, y - 14), (190, y + 10), (270, y - 14), (350, y + 12), (430, y - 8)]
        dense = catmull(wave, 5)
        ws = pressure_widths(len(dense), pr, [1, 1, 10, 2, 1.5][i], [10, 2.5, 10, 28, 14][i], [0.3, 0.1, 0, 0.4, 0][i])
        op = [1, 0.7, 0.55, 0.9, 1][i]
        b.append(f'<path d="{variable_outline(dense, ws)}" fill="{cols[i]}" fill-opacity="{op}"/>')
    pts = [(40, 330), (120, 290), (200, 360), (290, 300), (380, 350), (440, 310)]
    b.append(f'<path d="{bezier_d(pts)}" fill="none" stroke="#3fa34d" stroke-width="5" stroke-linecap="round"/>')
    for p in pts:
        b.append(f'<circle cx="{p[0]}" cy="{p[1]}" r="7" fill="#1d1d1f"/>')
    return (480, 400), "".join(b)

def th_mountains():
    return (480, 320), ('<rect width="480" height="320" fill="#cfe3ea"/><circle cx="380" cy="70" r="26" fill="#fff6d8"/>'
        '<polygon points="0,320 0,200 90,110 170,190 250,90 360,210 480,140 480,320" fill="#6b7f8f"/>'
        '<polygon points="250,90 280,124 262,120 240,132 226,118" fill="#f4f6f8"/><polygon points="90,110 112,134 96,130 80,140 72,128" fill="#f4f6f8"/>'
        '<polygon points="0,320 0,250 120,210 230,250 330,220 480,260 480,320" fill="#3f5a4a"/>')

def th_night():
    b = ['<rect width="480" height="320" fill="#1d2142"/>', '<circle cx="390" cy="60" r="22" fill="#f3efc9"/>', '<circle cx="400" cy="54" r="20" fill="#1d2142"/>']
    xs = [(10, 70, 160), (85, 50, 210), (140, 80, 130), (225, 60, 230), (290, 90, 170), (385, 85, 200)]
    for x, w, h in xs:
        b.append(f'<rect x="{x}" y="{320-h}" width="{w}" height="{h}" fill="#2f355e"/>')
        for wy in range(320 - h + 14, 310, 22):
            for wx in range(x + 8, x + w - 8, 16):
                if (wx * 7 + wy * 3) % 7 < 2:
                    b.append(f'<rect x="{wx}" y="{wy}" width="7" height="9" fill="#f2c45a"/>')
    return (480, 320), "".join(b)

def th_lily():
    b = ['<rect width="320" height="320" fill="#f8eef2"/>']
    for k in range(6):
        a = k * 60
        b.append(f'<ellipse cx="160" cy="105" rx="26" ry="62" fill="#f4a6c8" transform="rotate({a} 160 160)"/>')
        b.append(f'<line x1="160" y1="160" x2="160" y2="70" stroke="#c0267a" stroke-width="3" transform="rotate({a} 160 160)"/>')
    b.append('<circle cx="160" cy="160" r="18" fill="#c8d84a"/>')
    for k in range(6):
        a = math.radians(k * 60 + 30)
        b.append(f'<circle cx="{160+44*math.cos(a):.1f}" cy="{160+44*math.sin(a):.1f}" r="5" fill="#6e2a12"/>')
    return (320, 320), "".join(b)

def th_badge():
    return (320, 320), ('<rect width="320" height="320" fill="#1a1d19"/><circle cx="160" cy="160" r="120" fill="#4fc46a"/>'
        '<circle cx="160" cy="160" r="96" fill="#163a20"/><polygon points="160,84 182,140 242,140 194,174 212,232 160,198 108,232 126,174 78,140 138,140" fill="#f2b33d"/>')

def th_waves():
    b = ['<rect width="480" height="320" fill="#0f3b5a"/>']
    for i in range(7):
        y = 40 + i * 40
        pts = [(0, y), (80, y - 16), (160, y + 10), (240, y - 14), (320, y + 12), (400, y - 10), (480, y + 6)]
        b.append(f'<path d="{bezier_d(pts)}" fill="none" stroke="{["#72b8e6", "#a9d3e4", "#2d6a8f"][i % 3]}" stroke-width="{4 + i}" stroke-linecap="round"/>')
    return (480, 320), "".join(b)

def th_poster():
    return (320, 480), ('<rect width="320" height="480" fill="#f2ecd2"/><rect x="30" y="30" width="260" height="260" fill="#d94a3a"/>'
        '<circle cx="160" cy="160" r="90" fill="#f2ecd2"/><rect x="30" y="320" width="200" height="22" fill="#1b1a17"/><rect x="30" y="354" width="140" height="14" fill="#1b1a17"/>'
        '<rect x="30" y="380" width="170" height="14" fill="#1b1a17"/>')

def th_grid():
    b = ['<rect width="400" height="400" fill="#fbf8ef"/>']
    for i in range(5):
        for j in range(5):
            col = ["#3fa34d", "#e0a526", "#2d6a8f", "#d94a3a"][(i + j * 2) % 4]
            if (i + j) % 2:
                b.append(f'<circle cx="{40+i*80}" cy="{40+j*80}" r="28" fill="{col}"/>')
            else:
                b.append(f'<rect x="{14+i*80}" y="{14+j*80}" width="52" height="52" fill="{col}"/>')
    return (400, 400), "".join(b)

def th_bird():
    pts = [(60, 150), (150, 92), (240, 128), (330, 58), (430, 104)]
    return (480, 320), ('<rect width="480" height="320" fill="#eef3f6"/>' + f'<path d="{bezier_d(pts)}" fill="none" stroke="#1d1d1f" stroke-width="4" stroke-linecap="round"/>'
        + f'<path d="{bezier_d([(80, 240), (180, 200), (260, 236), (400, 190)])}" fill="none" stroke="#3b2a6b" stroke-width="3" stroke-linecap="round"/>')

def th_tree():
    return (320, 400), ('<rect width="320" height="400" fill="#e8f0de"/><rect x="146" y="250" width="28" height="120" fill="#6e4a2a"/>'
        '<circle cx="160" cy="180" r="90" fill="#3fa34d"/><circle cx="110" cy="210" r="56" fill="#2f6b3a"/><circle cx="215" cy="205" r="60" fill="#2f8a44"/>'
        '<rect x="0" y="360" width="320" height="40" fill="#8fbf5a"/>')

def thumb(x, y, w, h, fn, name):
    (dw, dh), body = fn()
    body = re.sub(r'(\d+\.\d)\d+', r'\1', body)
    body = re.sub(r' id="[^"]*"', '', body)
    sc = min(w / dw, h / dh)
    ox, oy = x + (w - dw * sc) / 2, y + (h - dh * sc) / 2
    return G(f"thumbnail {name} (rendered drawing)", R(x, y, w, h, "coal", "thumb bg"),
             f'<g id="drawing {name} (vector)" data-size="{dw} {dh}" transform="translate({ox:.2f} {oy:.2f}) scale({sc:.5f})">{body}</g>')

# ---------- pieces ----------
def wordmark(x, y):
    return G("logo Linkz (small)", T("L", x, y, "lizard"), T("ink", x + 7, y, "sky"), T("z", x + 28, y, "lizard"))

def topbar(lang):
    t = T_[lang]
    items = [R(0, 0, W, 22, "panel", "top bar bg"), R(0, 21, W, 1, "void", "top bar line"),
             I("layer", 8, 7, "lizard", "logo mark (lizard layers)"), wordmark(20, 4)]
    sx, sw = 150, 212
    items.append(field(sx, 3, sw, 16, [I("search", sx + 5, 7, "mute"), T(t["search"], sx + 17, 4, "dim" if False else "mute", "placeholder")], "search field"))
    lw = tw(t["lang"]) + 18
    items.append(button(W - 8 - 16 - 4 - lw, 4, lw, 14, "button language", icon="language", label=t["lang"]))
    items.append(button(W - 8 - 16, 4, 16, 14, "button settings", icon="settings"))
    return G("top bar", *items)

def section_title(s, y, right=None):
    items = [T(s, 12, y, "mute", f"section title {s}"), R(12 + tw(s) + 6, y + 7, W - 24 - tw(s) - 6 - (tw(right) + 26 if right else 0), 1, "bevel", "rule")]
    if right:
        rx = W - 12 - tw(right) - 20
        items.append(G("sort dropdown", I("sort", rx, y + 3, "mute"), T(right, rx + 10, y, "text"), I("chevron-down", rx + 12 + tw(right), y + 3, "mute")))
    return G(f"section {s}", *items)

def card(x, y, w, th_h, fn, name, when, state="normal", big=False, lang="en"):
    h = th_h + (32 if big else 28)
    outline = "lizard" if state == "hover" else "void"
    items = [frame_lines(x, y, w, h, outline, "card outline" + (" (hover)" if state == "hover" else "")),
             R(x + 1, y + 1, w - 2, h - 2, "row" if state == "hover" else "panel", "card bg"),
             thumb(x + 1, y + 1, w - 2, th_h, fn, name), R(x + 1, y + 1 + th_h, w - 2, 1, "void", "thumb line")]
    maxc = (w - 6) // 7
    nm = name if len(name) <= maxc else name[:maxc - 1] + "…"
    items.append(T(nm, x + 3, y + th_h + 2, "bright" if state == "hover" else "text", f"file name {name}"))
    wh = when if len(when) <= maxc else when[:maxc - 1] + "…"
    items.append(T(wh, x + 3, y + th_h + (15 if big else 13), "mute", f"edited {when}"))
    if state == "hover":
        bx = x + w - 18
        items.append(button(bx, y + 4, 14, 12, "button ⋯ more (pressed)", icon="more", state="pressed"))
    return G(f"card {name}" + (" (hover, menu open)" if state == "hover" else ""), *items), h

def new_card(x, y, w, h, lang, big=True):
    t = T_[lang]
    items = [R(x, y, w, h, "coal", "card bg")]
    # dashed outline 1u (3 on 2 off)
    d = []
    for i in range(0, w, 5):
        l = min(3, w - i); d.append(f"M{x+i} {y}h{l}v1h{-l}z M{x+i} {y+h-1}h{l}v1h{-l}z")
    for j in range(0, h, 5):
        l = min(3, h - j); d.append(f"M{x} {y+j}h1v{l}h-1z M{x+w-1} {y+j}h1v{l}h-1z")
    items.append(f'<path id="dashed outline" fill="{C["bevel"]}" d="{"".join(d)}"/>')
    cx, cy = x + w // 2, y + h // 2 - 10
    items.append(G("plus badge", frame_lines(cx - 10, cy - 10, 20, 20, "void"), R(cx - 9, cy - 9, 18, 18, "moss", "badge face"), I("plus", cx - 4, cy - 4, "lizard")))
    lines = [t["new"]] if tw(t["new"]) <= w - 8 else t["new"].split(" ", 1)
    for i, ln in enumerate(lines):
        items.append(T(ln, cx - tw(ln) // 2, cy + 14 + i * 13, "lizard", "label " + ln))
    items.append(T("Cmd+N", cx - tw("Cmd+N") // 2, cy + 27 + (len(lines) - 1) * 13, "mute", "shortcut"))
    return G("card + New drawing", *items)

def menu(x, y, lang):
    t = T_[lang]
    entries = list(zip(["open", "rename", "copy", "folder", "trash"], t["menu"]))
    w = max(tw(s) for _, s in entries) + 28
    h = len(entries) * 14 + 6 + 4
    x = min(x, W - w - 6)
    items = [R(x + 2, y + 2, w, h, "void", "drop shadow"), frame_lines(x, y, w, h, "void"), R(x + 1, y + 1, w - 2, h - 2, "row", "menu bg"),
             R(x + 1, y + 1, w - 2, 1, "bevel", "top highlight")]
    yy = y + 3
    for i, (ic, s) in enumerate(entries):
        if i == 4:
            items.append(R(x + 4, yy + 1, w - 8, 1, "bevel", "separator")); yy += 4
        col = "ember" if ic == "trash" else "text"
        it = []
        if i == 0:
            it.append(R(x + 1, yy, w - 2, 14, "moss", "item bg (hover)"))
            col = "bright"
        it += [I(ic, x + 6, yy + 3, "ember" if ic == "trash" else ("lizard" if i == 0 else "mute")), T(s, x + 18, yy, col)]
        items.append(G(f"menu item {s}" + (" (hover)" if i == 0 else ""), *it))
        yy += 14
    return G("card menu ⋯ (open)", *items)

def statusbar_home(lang):
    t = T_[lang]
    return gen.statusbar([("folder", t["count"], "mute")], [(None, "Linkz 0.1 · lizard-script", "dim" if False else "mute")])

DRAWINGS = [
    ("layers.lzdraw", th_sunset, "2h"), ("strokes.lzdraw", th_strokes, "y"), ("night-city.lzdraw", th_night, "mon"),
    ("hello.lzdraw", th_hello, "1w"), ("mountains.lzdraw", th_mountains, "2w"), ("lily.lzdraw", th_lily, "4d"),
    ("badge.lzdraw", th_badge, "5d"), ("waves.lzdraw", th_waves, "3w"), ("poster.lzdraw", th_poster, "1m"),
    ("grid.lzdraw", th_grid, "1m"), ("bird.lzdraw", th_bird, "2m"), ("tree.lzdraw", th_tree, "2m"),
]

def frame_gallery(lang="en"):
    t = T_[lang]
    parts = [topbar(lang)]
    # Working on: + New and 3 most recent, big cards
    y0 = 30
    parts.append(section_title(t["working"], y0))
    cw, gap, x, nw = 128, 8, 12, 80
    th_h = 76
    big = [new_card(x, y0 + 16, nw, th_h + 32, lang)]
    for i, (n, fn, when) in enumerate(DRAWINGS[:3]):
        el, h = card(x + nw + gap + i * (cw + gap), y0 + 16, cw, th_h, fn, n, t["big"][when], big=True, lang=lang)
        big.append(el)
    parts.append(G("row working on", *big))
    # All drawings grid: 6 per row, 2 rows
    y1 = y0 + 16 + th_h + 32 + 10
    parts.append(section_title(t["all"], y1, t["sort"]))
    sw = 91
    sth = 50
    grid = []
    hover_idx = 4
    menu_el = ""
    for i, (n, fn, when) in enumerate(DRAWINGS[:10]):
        r, c = divmod(i, 5)
        cx = 12 + c * (sw + 8)
        cy = y1 + 16 + r * (sth + 28 + 6)
        el, h = card(cx, cy, sw, sth, fn, n.replace(".lzdraw", ""), t["ago"][when], state="hover" if i == hover_idx else "normal", lang=lang)
        grid.append(el)
        if i == hover_idx:
            menu_el = menu(cx + sw - 18, cy + 18, lang)
    parts.append(G("grid all drawings", *grid))
    parts.append(menu_el)
    parts.append(G("scrollbar", R(W - 7, y1 + 16, 3, 2 * (sth + 34) - 6, "panel", "track"), R(W - 7, y1 + 16, 3, 110, "bevel", "thumb")))
    parts.append(gen.cursor(12 + 4 * (sw + 8) + sw - 12, y1 + 16 + 24))
    parts.append(statusbar_home(lang))
    return svg_doc(W * S, H * S, "".join(parts), bg="coal", scale=S)

def frame_empty(lang="en"):
    t = T_[lang]
    parts = [topbar(lang)]
    # lizard sprite (frame 3) at 3 UI px per art px, centered
    src = open(os.path.join(gen.ROOT, "lib/ui/splash_art.liz"), encoding="utf-8").read()
    sps = re.findall(r'\{"size": \((\d+), (\d+)\), "runs": \[(.*?)\]\}', src)
    lw, lh, runs = int(sps[3][0]), int(sps[3][1]), [(int(a), int(b), int(c), d) for a, b, c, d in re.findall(r'\[(\d+), (\d+), (\d+), (#[0-9a-f]{6})\]', sps[3][2])]
    k = 3
    lx, ly = (W - lw * k) // 2, 72
    parts.append(G("lizard (frame 3, x3) on a ground line", R(lx - 20, ly + lh * k - 3, lw * k + 40, 1, "bevel", "ground line"),
                   gen.sprite_svg(runs, lx, ly, k, "lizard")))
    parts.append(T(t["empty_title"], (W - tw(t["empty_title"])) // 2, ly + lh * k + 12, "bright", "title"))
    parts.append(T(t["empty_sub"], (W - tw(t["empty_sub"])) // 2, ly + lh * k + 28, "mute", "subtitle"))
    # big new button + example button
    bw = tw(t["new"]) + 40
    bx, by = (W - bw) // 2, ly + lh * k + 52
    parts.append(G("button + New drawing (primary, big)", frame_lines(bx, by, bw, 26, "void"), R(bx + 1, by + 1, bw - 2, 24, "moss", "face"),
                   gen.bevel(bx + 1, by + 1, bw - 2, 24, "lizard", "void"), I("plus", bx + 10, by + 9, "lizard"), T(t["new"], bx + 24, by + 6, "bright"),
                   ))
    ew = tw(t["example"]) + 18
    ex = (W - ew) // 2
    parts.append(button(ex, by + 34, ew, 16, "button Open an example", icon="open", label=t["example"]))
    parts.append(T(t["empty_hint"], (W - tw(t["empty_hint"])) // 2, by + 58, "dim" if False else "mute", "drop hint"))
    parts.append(gen.statusbar([("folder", "~/Linkz", "mute")], [(None, "Linkz 0.1 · lizard-script", "mute")]))
    return svg_doc(W * S, H * S, "".join(parts), bg="coal", scale=S)

if __name__ == "__main__":
    for k, v in {"f8": frame_gallery("en"), "f9": frame_empty("en"), "f10": frame_gallery("es")}.items():
        open(os.path.join(gen.OUT, k + ".svg"), "w").write(v)
        print(k, len(v))
