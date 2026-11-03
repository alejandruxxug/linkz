# Generates the Linkz R001 mockups as SVG (UI pixels x2, bitmap font as paths).
import json, re, math, os, sys
from xml.sax.saxutils import escape

ROOT = "/Users/urrexx/Documents/Uni/4th Semester/Langs&Compilers/linkz"
OUT = os.path.join(os.path.dirname(os.path.abspath(__file__)), "out")
os.makedirs(OUT, exist_ok=True)
S = 2  # one UI pixel = 2 screen pixels
E = 1  # emitted units per UI pixel (the svg root scales x2)

PAL = [
    ("void",   "#0b0c0b", "Outlines, bevel shadow, canvas frame"),
    ("coal",   "#131512", "Window background, tab bar, inset fields"),
    ("panel",  "#1a1d19", "Docked panel background, toolbar, status bar"),
    ("row",    "#232721", "Hovered row, tooltip and menu background"),
    ("raised", "#2e332c", "Button face, header strips"),
    ("bevel",  "#434a40", "Bevel highlight (top/left), splitters"),
    ("dim",    "#5f685b", "Disabled icons, guides, blank-line marks (never body text)"),
    ("mute",   "#a0ab9b", "Secondary text, comments, panel titles"),
    ("text",   "#d4dacd", "Body text, code identifiers, icons"),
    ("bright", "#f2f5ec", "Selected-row text, headings, tooltip titles"),
    ("moss",   "#163a20", "Selected row, active tool face"),
    ("lizard", "#4fc46a", "Accent: selection, anchors, keywords... active tool icon, saved"),
    ("amber",  "#f2b33d", "Focus outline, caret, hovered anchor, numbers in code"),
    ("ember",  "#f36c5b", "Errors"),
    ("sky",    "#72b8e6", "Keywords in code, links, info"),
    ("peach",  "#eb9a6f", "Strings in code"),
]
C = {n: h for n, h, _ in PAL}

# ---------- font ----------
FONT = {}
for line in open(os.path.join(ROOT, "lib/gfx/font_data.liz"), encoding="utf-8"):
    m = re.match(r'^\s*"((?:\\.|[^"\\])*)": (\[.*\]),?\s*$', line)
    if m:
        FONT[json.loads('"' + m.group(1) + '"')] = json.loads(m.group(2))
FALLBACK = FONT.get("fallback", [])
GW = 7

def tw(s):
    return len(s) * GW

GLYPH_RECTS = {}
def glyph_rects(ch):
    if ch in GLYPH_RECTS: return GLYPH_RECTS[ch]
    runs = sorted(FONT.get(ch, FALLBACK))
    rects = []  # [x0, y, w, h]
    open_ = {}
    for (gy, x0, x1) in runs:
        key = (x0, x1)
        if key in open_ and open_[key][1] + open_[key][3] == gy:
            open_[key][3] += 1
        else:
            r = [x0, gy, x1 - x0, 1]; rects.append(r); open_[key] = r
    GLYPH_RECTS[ch] = rects
    return rects

def rects_d(rects):
    """rects in absolute UI units -> compact relative path."""
    d, cx, cy = [], 0, 0
    for (x, y, w, h) in rects:
        d.append(f"{'M' if not d else 'm'}{x - cx if d else x} {y - cy if d else y}h{w}v{h}h{-w}z")
        cx, cy = x, y
    return "".join(d)

def text_d(s, x, y, k=1):
    rects = []
    for i, ch in enumerate(s):
        for (gx, gy, w, h) in glyph_rects(ch):
            rects.append(((x + i*GW + gx)*k, (y + gy)*k, w*k, h*k))
    return rects_d(rects)

def nm(s):
    return escape(s, {'"': "&quot;"})

def T(s, x, y, color="text", name=None):
    if not s.strip():
        return ""
    return f'<path id="{nm(name or ("text: " + s))}" fill="{C.get(color, color)}" d="{text_d(s, x, y)}"/>'

# ---------- primitives (UI units) ----------
def R(x, y, w, h, color, name=None):
    nid = f' id="{nm(name)}"' if name else ""
    return f'<rect{nid} x="{x}" y="{y}" width="{w}" height="{h}" fill="{C.get(color, color)}"/>'

def G(name, *items):
    return f'<g id="{nm(name)}">' + "".join(i for i in items if i) + "</g>"

def frame_lines(x, y, w, h, color, name="outline"):
    d = f"M{x} {y}h{w}v{h}h{-w}z M{x+1} {y+1}v{h-2}h{w-2}v{-(h-2)}z"
    return f'<path id="{nm(name)}" fill-rule="evenodd" fill="{C.get(color,color)}" d="{d}"/>'

def bevel(x, y, w, h, hi, lo, name="bevel"):
    # top/left in hi, bottom/right in lo, 1u inside (x,y,w,h)
    return G(name,
             R(x, y, w-1, 1, hi, "bevel top"), R(x, y, 1, h-1, hi, "bevel left"),
             R(x+1, y+h-1, w-1, 1, lo, "bevel bottom"), R(x+w-1, y+1, 1, h-1, lo, "bevel right"))

ICONS = {
 "circle": "..####..|.#....#.|#......#|#......#|#......#|#......#|.#....#.|..####..",
 "rect": "########|#......#|#......#|#......#|#......#|#......#|#......#|########",
 "line": ".......#|......#.|.....#..|....#...|...#....|..#.....|.#......|#.......",
 "polygon": "...##...|..#..#..|.#....#.|#......#|#......#|.#....#.|.#....#.|.######.",
 "stroke": "........|.##.....|#..#....|#...#..#|.....##.|........|........|........",
 "pencil": "......##|.....#.#|....#.#.|...#.#..|..#.#...|.#.#....|##......|#.......",
 "brush": "......##|.....##.|....##..|...##...|..#.....|.###....|####....|##......",
 "text": "########|#..##..#|...##...|...##...|...##...|...##...|...##...|..####..",
 "function": "....###.|...#....|...#....|.#####..|...#....|...#....|...#....|###.....",
 "layer": "...##...|.##..##.|#......#|.##..##.|#..##..#|.##..##.|...##...|........",
 "eye": "........|..####..|.#....#.|#..##..#|#..##..#|.#....#.|..####..|........",
 "lock": "..####..|.#....#.|.#....#.|########|###..###|###..###|########|########",
 "chevron-down": "........|........|##....##|.##..##.|..####..|...##...|........|........",
 "chevron-right": "..#.....|..##....|...##...|....##..|...##...|..##....|..#.....|........",
 "save": "#######.|#.###..#|#.###..#|#......#|#.####.#|#.#..#.#|#.####.#|########",
 "export": "...##...|..####..|.#.##.#.|...##...|...##...|#..##..#|#......#|########",
 "language": "..####..|.#.#..#.|#..#...#|########|#..#...#|#..#...#|.#.#..#.|..####..",
 "block": "..#..#..|.#....#.|.#....#.|#......#|.#....#.|.#....#.|..#..#..|........",
 "canvas": "########|#......#|#.####.#|#.####.#|#.####.#|#.####.#|#......#|########",
 "plus": "........|...##...|...##...|.######.|.######.|...##...|...##...|........",
 "error": "...##...|..####..|..####..|..####..|...##...|........|...##...|...##...",
 "check": "........|.......#|......##|#....##.|##..##..|.####...|..##....|........",
 "close": "........|.#....#.|..#..#..|...##...|...##...|..#..#..|.#....#.|........",
 "dot": "........|........|..####..|.######.|.######.|..####..|........|........",
 "comment": "........|........|........|##.##...|........|..##.##.|........|........",
}
ICON_ORDER = ["circle", "rect", "line", "polygon", "stroke", "pencil", "brush", "text", "function",
              "layer", "eye", "lock", "chevron-down", "chevron-right", "save", "export", "language",
              "block", "canvas", "plus", "error", "check", "close", "dot", "comment"]

def icon_d(name, x, y, s=E):
    d = []
    for r, row in enumerate(ICONS[name].split("|")):
        c = 0
        while c < 8:
            if row[c] == "#":
                e = c
                while e < 8 and row[e] == "#":
                    e += 1
                d.append(f"M{(x+c)*s} {(y+r)*s}h{(e-c)*s}v{s}h{-(e-c)*s}z")
                c = e
            else:
                c += 1
    return "".join(d)

def I(name, x, y, color="text", label=None):
    return f'<path id="{nm(label or ("icon: " + name))}" fill="{C.get(color,color)}" d="{icon_d(name, x, y)}"/>'

def button(x, y, w, h, name, icon=None, label=None, state="normal", icolor="text"):
    face = {"normal": "raised", "hover": "bevel", "pressed": "moss", "active": "moss"}[state]
    hi, lo = ("bevel", "coal") if state in ("normal",) else (("dim", "coal") if state == "hover" else ("coal", "bevel"))
    off = 1 if state in ("pressed", "active") else 0
    if state == "active":
        icolor = "lizard"
    items = [frame_lines(x, y, w, h, "void"), R(x+1, y+1, w-2, h-2, face, "face"), bevel(x+1, y+1, w-2, h-2, hi, lo)]
    cx = x + 3 + off
    if icon:
        items.append(I(icon, x + (w - 8)//2 + off if not label else cx, y + (h - 8)//2 + off, icolor))
        cx += 10
    if label:
        items.append(T(label, cx + (0 if icon else 1), y + (h - 13)//2 + off, "bright" if state in ("pressed", "active") else "text"))
    return G(name, *items)

def field(x, y, w, h, value_items, name, focused=False):
    items = [frame_lines(x, y, w, h, "amber" if focused else "void"), R(x+1, y+1, w-2, h-2, "coal", "inset"),
             R(x+1, y+1, w-2, 1, "void", "inset shadow top")]
    return G(name + (" (focused)" if focused else ""), *items, *value_items)

# ---------- code tokens ----------
KW = {"let", "layer", "for", "function", "if", "else", "return", "true", "false", "null"}
TOK = re.compile(r'(-~>.*$)|("(?:[^"\\]|\\.)*")|(#[0-9a-fA-F]{3,8})|(\d+(?:\.\d+)?)|([A-Za-z_][A-Za-z_0-9.]*)|(\s+)|(.)')

def code(s, x, y, base="text", sel=False, maxw=None, name="code"):
    """Syntax-coloured code text. Colour literals get a 6x6 swatch before them."""
    items, cx = [], x
    groups = {}
    clipped = False
    for m in TOK.finditer(s):
        com, st, col, num, ident, ws, other = m.groups()
        tok = m.group(0)
        if com: color = "mute"
        elif st: color = "peach"
        elif col: color = base if not sel else "bright"
        elif num: color = "amber"
        elif ident: color = "sky" if ident in KW else (base if not sel else "bright")
        else: color = "mute" if not sel else "text"
        if col:
            if maxw and cx + 8 - x > maxw: clipped = True; break
            items.append(R(cx, y + 4, 6, 6, "void", "swatch outline"))
            items.append(R(cx + 1, y + 5, 4, 4, "#" + col[1:7] if len(col) >= 7 else col, f"swatch {col}"))
            cx += 8
        if maxw and cx + tw(tok) - x > maxw:
            room = (maxw - (cx - x)) // GW - 1
            tok = tok[:max(room, 0)]
            clipped = True
        if tok.strip():
            groups.setdefault(color, []).append(text_d(tok, cx, y))
        cx += tw(tok)
        if clipped: break
    for color, ds in groups.items():
        items.insert(0, f'<path id="{nm("code " + color)}" fill="{C.get(color,color)}" d="{"".join(ds)}"/>')
    if clipped:
        items.append(T("»", x + maxw - 7, y, "mute", "clip marker »"))
    return G(f'{name}: {s}', *items)

# ---------- window ----------
W, H = 512, 360
PANEL_W = 320
TOP_H, TAB_H, STATUS_H = 16, 14, 16
MAIN_Y = TOP_H + TAB_H
STATUS_Y = H - STATUS_H
ROW_H = 18
LAYER_GAP = 4   # extra room above each layer row so layers read as groups
CANVAS_X = PANEL_W + 2

TOOLS = [("circle",), ("rect",), ("line",), ("polygon",), ("stroke",), None, ("pencil",), ("brush",), None, ("text",), ("function",), ("layer",)]

L = {
 "en": {"layers": "LAYERS", "props": "PROPERTIES", "line": "line", "save": "Save", "export": "Export PNG",
        "saved": "saved", "unsaved": "unsaved", "steps": "steps", "lang": "EN", "hint": "drag the shape or its anchors",
        "empty": "empty row · pick a tool or type code", "pick": "or type code", "error": "1 error",
        "lastgood": "last good drawing", "points": "points", "pen": "pen"},
 "es": {"layers": "CAPAS", "props": "PROPIEDADES", "line": "línea", "save": "Guardar", "export": "Exportar PNG",
        "saved": "guardado", "unsaved": "sin guardar", "steps": "pasos", "lang": "ES", "hint": "arrastra la forma o sus anclas",
        "empty": "fila vacía · elige una herramienta o escribe", "pick": "o escribe código", "error": "1 error",
        "lastgood": "último dibujo bueno", "points": "puntos", "pen": "pluma"},
}

def toolbar(active=None, lang="en", pressed_right=None):
    t = L[lang]
    items = [R(0, 0, W, TOP_H, "panel", "toolbar bg"), R(0, TOP_H-1, W, 1, "void", "toolbar bottom line")]
    x = 3
    for tool in TOOLS:
        if tool is None:
            items.append(R(x + 1, 3, 1, 10, "bevel", "separator")); x += 4; continue
        nmx = tool[0]
        items.append(button(x, 1, 14, 14, f"tool {nmx}" + (" (active)" if active == nmx else ""), icon=nmx,
                            state="active" if active == nmx else "normal"))
        x += 15
    # right side: language, save, export
    ex_w = tw(t["export"]) + 18
    sv_w = tw(t["save"]) + 18
    lg_w = tw(t["lang"]) + 18
    rx = W - 3 - ex_w
    items.append(button(rx, 1, ex_w, 14, "button export png", icon="export", label=t["export"]))
    rx -= sv_w + 2
    items.append(button(rx, 1, sv_w, 14, "button save", icon="save", label=t["save"]))
    rx -= lg_w + 6
    items.append(button(rx, 1, lg_w, 14, "button language", icon="language", label=t["lang"]))
    return G("toolbar", *items)

def tabbar(unsaved=False):
    name = "layers.lzdraw"
    w = tw(name) + 26
    items = [R(0, TOP_H, W, TAB_H, "coal", "tab bar bg"), R(0, MAIN_Y-1, W, 1, "void", "tab bar bottom line")]
    tab = [R(2, TOP_H+1, w, TAB_H, "panel", "tab face"), R(2, TOP_H+1, w, 1, "lizard", "tab active top line"),
           R(1, TOP_H+1, 1, TAB_H-1, "void", "tab left"), R(2+w, TOP_H+1, 1, TAB_H-1, "void", "tab right"),
           T(name, 7, TOP_H+1), I("dot" if unsaved else "close", 7 + tw(name) + 4, TOP_H+4, "amber" if unsaved else "mute",
                                    "unsaved dot" if unsaved else "close tab")]
    items.append(G("tab layers.lzdraw (active)", *tab))
    items.append(G("tab new", T("+", 2 + w + 6, TOP_H+1, "mute")))
    return G("tab bar", *items)

def statusbar(parts_left, parts_right):
    items = [R(0, STATUS_Y, W, STATUS_H, "panel", "status bg"), R(0, STATUS_Y, W, 1, "void", "status top line")]
    x = 6
    for p in parts_left:
        ic, s, col = p
        if ic:
            items.append(I(ic, x, STATUS_Y + 4, col)); x += 10
        items.append(T(s, x, STATUS_Y + 2, col)); x += tw(s) + 6
        if p is not parts_left[-1]:
            items.append(R(x, STATUS_Y + 4, 1, 8, "bevel", "separator"))
        x += 7
    x = W - 6
    for p in reversed(parts_right):
        ic, s, col = p
        x -= tw(s)
        items.append(T(s, x, STATUS_Y + 2, col))
        if ic:
            x -= 10; items.append(I(ic, x, STATUS_Y + 4, col))
        x -= 7
        if p is not parts_right[0]:
            items.append(R(x, STATUS_Y + 4, 1, 8, "bevel", "separator"))
        x -= 6
    return G("status bar", *items)

# ---------- rows ----------
def file_rows(variant="main"):
    """Rows of examples/layers.lzdraw. kind, depth, code, line, extra flags."""
    rows = [
        dict(kind="comment", d=0, code="-~> lzdraw 1", line=1),
        dict(kind="comment", d=0, code="-~> A sunset in layers. Layers higher in the file are in front, and so are", line=2),
        dict(kind="comment", d=0, code="-~> the objects higher inside a layer. The loop is one object.", line=3),
        dict(kind="blank", line=4),
        dict(kind="value", d=0, code="let horizon = 200", line=7),
        dict(kind="blank", line=8),
        dict(kind="layer", d=0, code="layer sun", line=9, open=True),
        dict(kind="circle", d=1, code="circle((360, 90), 34, #fff3c4)", line=10),
        dict(kind="circle", d=1, code="circle((360, 90), 56, #f2a65a66)", line=11),
        dict(kind="add", d=1, line=None),
        dict(kind="blank", line=12),
        dict(kind="layer", d=0, code="layer boat", line=13, open=False, count=3),
        dict(kind="blank", line=17),
        dict(kind="layer", d=0, code="layer sea", line=18, open=True),
        dict(kind="layer", d=1, code="layer waves", line=19, open=True),
        dict(kind="block", d=2, code="for (let i = 0; i < 5; i = i + 1) {", line=20, open=False, count=3),
        dict(kind="rect", d=1, code="rect((0, horizon), (480, 120), #2d6a8f)", line=23),
        dict(kind="blank", line=24),
        dict(kind="layer", d=0, code="layer sky", line=25, open=True),
        dict(kind="rect", d=1, code="rect((0, 0), (480, horizon), #f2a65a)", line=26),
    ]
    return rows

KIND_ICON = {"circle": "circle", "rect": "rect", "line": "line", "polygon": "polygon", "stroke": "stroke",
             "brush": "brush", "value": "function", "layer": "layer", "block": "block", "canvas": "canvas",
             "comment": "comment", "text": "text"}

def row_el(r, y, w=PANEL_W, sel=False, hover=False, error=False, caret=None):
    k = r["kind"]
    if k == "blank":
        return "", 4
    if k == "add":
        x = 44 + r["d"] * 10
        return G("row ⊕ add object to layer sun", I("plus", x - 10, y + 2, "dim"), T("add", x, y - 1, "dim")), 10
    items = []
    bg = "moss" if sel else ("row" if hover else None)
    if bg:
        items.append(R(0, y, w, ROW_H, bg, "row bg (selected)" if sel else "row bg (hover)"))
    if sel:
        items.append(R(0, y, 2, ROW_H, "lizard", "selection bar"))
    if error:
        items.append(R(0, y, 2, ROW_H, "ember", "error bar"))
    if k != "comment":
        items.append(I("eye", 4, y + 5, "mute" if not sel else "text", "toggle visible (eye)"))
        items.append(I("lock", 14, y + 5, "dim", "toggle lock (unlocked)"))
    x = 24 + r["d"] * 10
    for dd in range(r["d"]):
        items.append(R(24 + dd * 10 + 3, y, 1, ROW_H, "bevel", f"indent guide {dd+1}"))
    if r.get("open") is not None:
        items.append(I("chevron-down" if r["open"] else "chevron-right", x, y + 5, "mute", "expand chevron"))
    if error:
        items.append(I("error", x + 10, y + 5, "ember", "kind icon: error"))
    elif k != "comment":
        items.append(I(KIND_ICON[k], x + 10, y + 5, "lizard" if sel else "text", f"kind icon: {k}"))
    cx = x + 22 if k != "comment" else 24
    maxw = w - cx - 4
    if r.get("count"):
        maxw -= tw(str(r["count"])) + 8
    items.append(code(r["code"], cx, y + 3, sel=sel, maxw=maxw))
    if r.get("count"):
        bx = w - tw(str(r["count"])) - 8
        items.append(G("collapsed count", R(bx, y + 4, tw(str(r["count"])) + 4, 10, "raised", "count badge"),
                       T(str(r["count"]), bx + 2, y + 3, "mute")))
    if caret is not None:
        items.append(R(cx + caret, y + 3, 1, 12, "amber", "caret"))
    label = f'row line {r["line"]} · {k}' + (" (selected)" if sel else "") + (" (error)" if error else "") + (" (hover)" if hover else "")
    return G(label, *items), ROW_H

def layers_panel(lang="en", rows=None, sel_line=None, hover_line=None, error_line=None, extra_after=None, insert=None,
                 props=None, tooltip=None, caret_line=None):
    t = L[lang]
    items = [R(0, MAIN_Y, PANEL_W, STATUS_Y - MAIN_Y, "panel", "panel bg")]
    items.append(G("panel header", R(0, MAIN_Y, PANEL_W, 13, "raised", "header bg"), R(0, MAIN_Y + 13, PANEL_W, 1, "void", "header line"),
                   T(t["layers"], 6, MAIN_Y, "mute"),
                   T("layers.lzdraw", PANEL_W - tw("layers.lzdraw") - 6, MAIN_Y, "dim")))
    y = MAIN_Y + 15
    rows_g = []
    hidden = False
    for r in rows:
        if r["kind"] == "layer":
            y += LAYER_GAP
        if y + (ROW_H if r["kind"] not in ("blank", "add") else 4) > PROPS_Y - ROW_H - 2:
            hidden = True
            break
        el, h = row_el(r, y, sel=(r.get("line") == sel_line and r.get("line") is not None),
                       hover=(r.get("line") == hover_line and hover_line is not None),
                       error=(r.get("line") == error_line and error_line is not None),
                       caret=r.get("caret"))
        rows_g.append(el)
        y += h
        if extra_after and r.get("line") == extra_after[0]:
            el2, h2 = extra_after[1](y)
            rows_g.append(el2); y += h2
    items.append(G("rows", *rows_g))
    if hidden:
        top, bot = MAIN_Y + 15, PROPS_Y - ROW_H - 3
        items.append(G("scrollbar (more rows below)", R(PANEL_W - 5, top, 4, bot - top, "coal", "track"),
                       R(PANEL_W - 5, top, 4, (bot - top) * 2 // 3, "bevel", "thumb")))
    # pinned canvas row
    py = PROPS_Y - ROW_H - 1
    canvas_row = dict(kind="canvas", d=0, code='canvas = {"size": (480, 320), "background": #f6d9a8}', line=5)
    el, _ = row_el(canvas_row, py)
    items.append(G("pinned canvas row (line 5)", R(0, py - 1, PANEL_W, 1, "void", "pin line"), R(0, py, PANEL_W, ROW_H, "coal", "pinned bg"), el))
    items.append(props or "")
    return G("layers panel (the code)", *items)

PROPS_Y = 282

def props_panel(lang, title, fields_fn, hint):
    t = L[lang]
    y = PROPS_Y
    items = [R(0, y, PANEL_W, STATUS_Y - y, "panel", "properties bg"), R(0, y, PANEL_W, 1, "void", "top line"),
             R(0, y + 1, PANEL_W, 15, "raised", "header bg"), R(0, y + 16, PANEL_W, 1, "void", "header line"),
             T(t["props"] + " · " + title, 6, y + 2, "mute")]
    items.append(fields_fn(y + 22))
    if hint:
        items.append(T(hint, 6, y + 45, "mute", "hint: " + hint))
    return G("properties panel", *items)

def num_field(x, y, label, value, name, focused=False, w=None):
    w = w or tw(value) + 8
    return G(f"field {name}", T(label, x, y + 2, "mute", f"label {label}"),
             field(x + tw(label) + 4, y, w, 18, [T(value, x + tw(label) + 8, y + 3, "amber")], f"input {name}", focused))

def point_field(x, y, label, a, b, name, focused=False):
    lw = tw(label) + 4
    fx = x + lw
    items = [T(label, x, y + 2, "mute", f"label {label}"), T("(", fx, y + 2, "mute")]
    fx += 8
    items.append(field(fx, y, tw(a) + 8, 18, [T(a, fx + 4, y + 3, "amber")], f"input {name}.x", focused))
    fx += tw(a) + 10
    items.append(T(",", fx - 2, y + 2, "mute"))
    fx += 6
    items.append(field(fx, y, tw(b) + 8, 18, [T(b, fx + 4, y + 3, "amber")], f"input {name}.y"))
    fx += tw(b) + 9
    items.append(T(")", fx, y + 2, "mute"))
    return G(f"field {name} (point)", *items), fx + 8

def color_field(x, y, label, hexv, name):
    lw = tw(label) + 4
    w = tw(hexv) + 22
    fx = x + lw
    sw = hexv[:7]
    return G(f"field {name} (color)", T(label, x, y + 2, "mute", f"label {label}"),
             field(fx, y, w, 18, [R(fx + 3, y + 4, 10, 10, "void", "swatch outline"), R(fx + 4, y + 5, 8, 8, sw, f"swatch {hexv}"),
                                   T(hexv, fx + 16, y + 3, "text")], f"input {name}")), fx + w + 8

def circle_props(r="34", fill="#fff3c4", focused=False):
    def f(y):
        g1, x = point_field(6, y, "center", "360", "90", "center", focused)
        g2 = num_field(x + 4, y, "r", r, "r")
        x2 = x + 4 + tw("r") + 4 + tw(r) + 8 + 8
        g3, _ = color_field(x2, y, "fill", fill, "fill")
        return G("fields (from circle(center, r, fill))", g1, g2, g3)
    return f

# ---------- canvas ----------
CV_X, CV_Y, SC = 654, 254, 0.75   # screen px of drawing origin, zoom

def cv(p):
    return (CV_X + p[0] * SC, CV_Y + p[1] * SC)

def catmull(points, steps=12):
    pts = [points[0]] + points + [points[-1]]
    out = [points[0]]
    for i in range(1, len(pts) - 2):
        p0, p1, p2, p3 = pts[i-1], pts[i], pts[i+1], pts[i+2]
        for s in range(1, steps + 1):
            t = s / steps
            t2, t3 = t*t, t*t*t
            x = 0.5*((2*p1[0]) + (-p0[0]+p2[0])*t + (2*p0[0]-5*p1[0]+4*p2[0]-p3[0])*t2 + (-p0[0]+3*p1[0]-3*p2[0]+p3[0])*t3)
            y = 0.5*((2*p1[1]) + (-p0[1]+p2[1])*t + (2*p0[1]-5*p1[1]+4*p2[1]-p3[1])*t2 + (-p0[1]+3*p1[1]-3*p2[1]+p3[1])*t3)
            out.append((x, y))
    return out

def poly_d(pts, close=False):
    d = "M" + " L".join(f"{x:.2f} {y:.2f}" for x, y in pts)
    return d + ("Z" if close else "")

def bezier_d(points):
    """Catmull-Rom through points as cubic Beziers (what lizard draws)."""
    pts = [points[0]] + points + [points[-1]]
    d = f"M{points[0][0]:.2f} {points[0][1]:.2f}"
    for i in range(1, len(pts) - 2):
        p0, p1, p2, p3 = pts[i-1], pts[i], pts[i+1], pts[i+2]
        c1 = (p1[0] + (p2[0]-p0[0])/6, p1[1] + (p2[1]-p0[1])/6)
        c2 = (p2[0] - (p3[0]-p1[0])/6, p2[1] - (p3[1]-p1[1])/6)
        d += f" C{c1[0]:.2f} {c1[1]:.2f} {c2[0]:.2f} {c2[1]:.2f} {p2[0]:.2f} {p2[1]:.2f}"
    return d

def variable_outline(points, widths):
    """Closed polygon around a centre line with per-sample widths."""
    left, right = [], []
    n = len(points)
    for i, (x, y) in enumerate(points):
        a = points[max(i-1, 0)]; b = points[min(i+1, n-1)]
        dx, dy = b[0]-a[0], b[1]-a[1]
        ln = math.hypot(dx, dy) or 1
        nx, ny = -dy/ln, dx/ln
        w = widths[i] / 2
        left.append((x + nx*w, y + ny*w)); right.append((x - nx*w, y - ny*w))
    return poly_d(left + right[::-1], True)

def pressure_widths(n, pressures, mn, mx, taper):
    out = []
    for i in range(n):
        t = i / (n - 1)
        # interpolate pressure
        f = t * (len(pressures) - 1)
        j = min(int(f), len(pressures) - 2)
        p = pressures[j] + (pressures[j+1] - pressures[j]) * (f - j)
        w = mn + (mx - mn) * p
        if taper:
            edge = min(t, 1 - t) / taper
            w *= min(1, max(0.15, edge))
        out.append(w)
    return out

SUN_RAW = None

def drawing(variant="main", extra=""):
    """The sunset, full-res and smooth, in drawing coordinates (scaled by a transform)."""
    hz = 200
    parts = []
    parts.append(f'<rect id="canvas background #f6d9a8" x="0" y="0" width="480" height="320" fill="#f6d9a8"/>')
    parts.append(f'<g id="layer sky"><rect id="rect line 26" x="0" y="0" width="480" height="{hz}" fill="#f2a65a"/></g>')
    waves = []
    for i in range(5):
        pts = [(40 + i*20, hz + 30 + i*18), (140 + i*30, hz + 24 + i*18), (240 + i*20, hz + 32 + i*18), (440 - i*10, hz + 26 + i*18)]
        waves.append(f'<path id="stroke i={i}" d="{bezier_d(pts)}" fill="none" stroke="#a9d3e4" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>')
    parts.append(f'<g id="layer sea"><rect id="rect line 23" x="0" y="{hz}" width="480" height="120" fill="#2d6a8f"/>'
                 f'<g id="layer waves (for loop line 20)">{"".join(waves)}</g></g>')
    parts.append('<g id="layer boat">'
                 '<polygon id="hull line 14" points="120,214 180,214 165,232 135,232" fill="#3b2a6b"/>'
                 '<line id="mast line 15" x1="150" y1="214" x2="150" y2="170" stroke="#3b2a6b" stroke-width="3" stroke-linecap="round"/>'
                 '<polygon id="sail line 16" points="152,174 152,208 182,208" fill="#fbf8ef"/></g>')
    sun = []
    if variant != "error":
        sun.append('<circle id="glow line 11" cx="360" cy="90" r="56" fill="#f2a65a" fill-opacity="0.4"/>')
    else:
        sun.append('<circle id="glow line 11 (last good render)" cx="360" cy="90" r="56" fill="#f2a65a" fill-opacity="0.4"/>')
    sun.append('<circle id="core line 10" cx="360" cy="90" r="34" fill="#fff3c4"/>')
    sun.append(extra)
    parts.append(f'<g id="layer sun">{"".join(sun)}</g>')
    body = "".join(parts)
    return (f'<g id="canvas drawing (vector, full-res, 75%)" data-size="480 320" transform="translate({CV_X/S} {CV_Y/S}) scale({SC/S})">{body}</g>')

def checker(x, y, w, h, a, b, size, name):
    d = []
    for j in range(0, h, size):
        for i in range(0, w, size):
            if ((i // size) + (j // size)) % 2:
                d.append(f"M{x+i} {y+j}h{min(size, w-i)}v{min(size, h-j)}h{-min(size, w-i)}z")
    return G(name, R(x, y, w, h, a, "checker base"), f'<path id="checker squares" fill="{C.get(b,b)}" d="{"".join(d)}"/>')

def canvas_area(drawing_svg, overlay="", dim=False, badge=None):
    x0, y0 = CANVAS_X, MAIN_Y
    w, h = W - x0, STATUS_Y - y0
    items = [R(PANEL_W, MAIN_Y, 2, STATUS_Y - MAIN_Y, "void", "splitter"),
             R(PANEL_W + 1, MAIN_Y, 1, STATUS_Y - MAIN_Y, "bevel", "splitter highlight"),
             checker(x0, y0, w, h, "coal", "#161815", 16, "canvas surround (checker)")]
    # canvas frame 1u void around 360x240 screen px at (654,254)
    fx, fy = (CV_X // S) - 1, (CV_Y // S) - 1
    items.append(frame_lines(fx, fy, 180 + 2, 120 + 2, "void", "canvas frame"))
    items.append(drawing_svg)
    if dim:
        items.append(f'<rect id="last-good dim" x="{CV_X//S}" y="{CV_Y//S}" width="180" height="120" fill="{C["coal"]}" fill-opacity="0.45"/>')
    items.append(overlay)
    if badge:
        bw = tw(badge[1]) + 18
        bx, by = fx, fy + 122 + 4
        items.append(G("canvas badge", frame_lines(bx, by, bw, 16, "void"), R(bx + 1, by + 1, bw - 2, 14, "row", "badge bg"),
                       I(badge[0], bx + 4, by + 4, badge[2]), T(badge[1], bx + 14, by + 1, badge[2])))
    return G("canvas area", *items)

def anchor(px, py, name, kind="normal"):
    # px,py in screen px (drawing projected); anchors are UI pixels, snapped to the 2px grid
    ux, uy = int(round(px / S)), int(round(py / S))
    s = 7 if kind == "hover" else 5
    h = s // 2
    fill = {"normal": "lizard", "alt": "bright", "hover": "amber"}[kind]
    return G(name + (" (hovered)" if kind == "hover" else ""), R(ux - h - 1, uy - h - 1, s + 2, s + 2, "void", "anchor outline"),
             R(ux - h, uy - h, s, s, fill, "anchor fill"))

def sel_box(x0, y0, x1, y1, name="selection outline"):
    ux0, uy0, ux1, uy1 = int(x0 // S), int(y0 // S), int(math.ceil(x1 / S)), int(math.ceil(y1 / S))
    w, h = ux1 - ux0, uy1 - uy0
    # dashed 1u lizard box: 3 on, 2 off
    d = []
    for i in range(0, w, 5):
        l = min(3, w - i)
        d.append(f"M{ux0+i} {uy0}h{l}v1h{-l}z M{ux0+i} {uy1-1}h{l}v1h{-l}z")
    for j in range(0, h, 5):
        l = min(3, h - j)
        d.append(f"M{ux0} {uy0+j}h1v{l}h-1z M{ux1-1} {uy0+j}h1v{l}h-1z")
    return f'<path id="{nm(name)}" fill="{C["lizard"]}" d="{"".join(d)}"/>'

def svg_doc(w, h, body, bg=None, scale=1):
    bgr = f'<rect id="window bg" x="0" y="0" width="{w}" height="{h}" fill="{C[bg]}"/>' if bg else ""
    if scale != 1:
        body = f'<g id="ui (half-res, shown x{scale})" transform="scale({scale})">{body}</g>'
    return f'<svg xmlns="http://www.w3.org/2000/svg" width="{w}" height="{h}" viewBox="0 0 {w} {h}">{bgr}{body}</svg>'

def window(lang, body_parts):
    return svg_doc(W*S, H*S, "".join(body_parts), bg="coal", scale=S)

def cursor(x, y, name="mouse cursor"):
    rows = ["#.......", "##......", "#o#.....", "#oo#....", "#ooo#...", "#oooo#..", "#ooooo#.", "#oooooo#", "#ooo####", "#o#o#...", "##.#o#..", "#...##.."]
    dv, do = [], []
    for r, row in enumerate(rows):
        for c, ch in enumerate(row):
            if ch == "#": dv.append(f"M{x+c} {y+r}h1v1h-1z")
            elif ch == "o": do.append(f"M{x+c} {y+r}h1v1h-1z")
    return G(name, f'<path id="cursor outline" fill="{C["void"]}" d="{"".join(dv)}"/>', f'<path id="cursor fill" fill="{C["bright"]}" d="{"".join(do)}"/>')

def tooltip_box(x, y, lines, name, w=None):
    """lines: list of (items_fn(x, y)) each 13u tall. Box: row bg, void outline, 2u void drop shadow."""
    h = len(lines) * 15 + 6
    w = w or 200
    items = [R(x + 2, y + 2, w, h, "void", "drop shadow"), frame_lines(x, y, w, h, "void"), R(x + 1, y + 1, w - 2, h - 2, "row", "tooltip bg"),
             R(x + 1, y + 1, w - 2, 1, "bevel", "top highlight")]
    for i, fn in enumerate(lines):
        items.append(fn(x + 6, y + 4 + i * 15))
    return G(name, *items)

# ---------------- frames ----------------
def status_main(lang, saved=True, extra_left=None, err=False):
    t = L[lang]
    steps = "1,204" if lang == "en" else "1.204"
    left = [("canvas", "layers.lzdraw", "text")]
    if err:
        left.append(("error", t["error"] + " · " + t["line"] + " 26", "ember"))
    elif saved:
        left.append(("check", t["saved"], "lizard"))
    else:
        left.append(("dot", t["unsaved"], "amber"))
    left.append(extra_left or (None, f"{steps} {t['steps']} · 3 ms", "mute"))
    right = ([] if (err or extra_left) else [(None, "x 352 y 88", "mute")]) + [(None, "75%", "mute"), ("language", t["lang"], "text")]
    return statusbar(left, right)

def circle_overlay(kind="sel"):
    x0, y0 = cv((326, 56)); x1, y1 = cv((394, 124))
    return G("selection (circle line 10)", sel_box(x0, y0, x1, y1),
             anchor(*cv((360, 90)), "anchor center (move)", "normal"),
             anchor(*cv((394, 90)), "anchor radius", "alt"))

def frame_main(lang="en"):
    t = L[lang]
    props = props_panel(lang, f"circle · {t['line']} 10", circle_props(), f"{t['line']} 10 · {t['hint']}")
    return window(lang, [toolbar(None, lang), tabbar(), layers_panel(lang, file_rows(), sel_line=10, props=props),
                         canvas_area(drawing(), circle_overlay()), status_main(lang)])

def frame_new_row():
    lang = "en"; t = L[lang]
    rows = file_rows()
    def strip(y):
        x = 44 + 10
        items = [R(0, y, PANEL_W, ROW_H + 4, "moss", "row bg (selected)"), R(0, y, 2, ROW_H + 4, "lizard", "selection bar"),
                 R(27, y, 1, ROW_H + 4, "bevel", "indent guide 1"),
                 I("eye", 4, y + 7, "text", "toggle visible (eye)"), I("lock", 14, y + 7, "dim", "toggle lock (unlocked)"),
                 R(x - 10 + 1, y + 5, 1, 12, "amber", "caret")]
        bx = x - 4
        for tl in ["circle", "rect", "line", "polygon", "stroke", "text", "function", "layer"]:
            items.append(button(bx, y + 4, 14, 14, f"strip tool {tl}" + (" (hover)" if tl == "rect" else ""), icon=tl,
                                state="hover" if tl == "rect" else "normal"))
            bx += 15
        items.append(T(t["pick"], bx + 4, y + 5, "mute", "strip hint"))
        return G("row line 12 · new empty row (inline tool strip)", *items), ROW_H + 4
    props = props_panel(lang, f"{t['line']} 12", lambda y: T(t["empty"], 6, y + 2, "mute", "empty hint"),
                        "Enter below · Shift+Enter above · Tab indent")
    tip = tooltip_box(68, 205, [lambda x, y: G("tooltip text", T("rect", x, y, "bright"), T(" · rect(pos, size, fill)", x + tw("rect"), y, "mute"))],
                      "tooltip strip rect", w=tw("rect · rect(pos, size, fill)") + 14)
    # insert after line 11; drop the add row
    rows = [r for r in rows if r["kind"] != "add"]
    panel = layers_panel(lang, rows, extra_after=(11, strip), props=props)
    return window(lang, [toolbar(None, lang), tabbar(unsaved=True), panel, canvas_area(drawing()), tip,
                         cursor(73, 194), status_main(lang, saved=False)])

STROKE_PTS = [(60, 150), (150, 92), (240, 128), (330, 58), (430, 104)]

def frame_pencil():
    lang = "en"; t = L[lang]
    rows = file_rows()
    code_s = "stroke([(60, 150), (150, 92), (240, 128), (330, 58), (430, 104)], #1d1d1f, 4)"
    new = dict(kind="stroke", d=1, code=code_s, line=10)
    out = []
    for r in rows:
        if r.get("line") == 10:
            out.append(new)
            r = dict(r, line=11)
        elif r.get("line") == 11:
            r = dict(r, line=12)
        out.append(r)
    out = [r for r in out if r["kind"] != "add"]
    def f(y):
        g1 = G("field points (list)", T("points", 6, y + 2, "mute"),
               field(6 + tw("points") + 4, y, tw("5 points") + 20, 18, [T("5 points", 6 + tw("points") + 8, y + 3, "amber"),
                     I("chevron-down", 6 + tw("points") + 4 + tw("5 points") + 10, y + 5, "mute")], "input points"))
        x = 6 + tw("points") + 4 + tw("5 points") + 28
        g2, x = color_field(x, y, "color", "#1d1d1f", "color")
        g3 = num_field(x, y, "width", "4", "width")
        return G("fields (from stroke(points, color, width))", g1, g2, g3)
    props = props_panel(lang, f"stroke · {t['line']} 10", f, "anchor 4 (330, 58) · drag, Del removes")
    path = bezier_d(STROKE_PTS)
    extra = f'<path id="stroke line 10 (new, pencil)" d="{path}" fill="none" stroke="#1d1d1f" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>'
    xs = [cv(p) for p in STROKE_PTS]
    ov = [sel_box(min(p[0] for p in xs) - 6, min(p[1] for p in xs) - 6, max(p[0] for p in xs) + 6, max(p[1] for p in xs) + 6)]
    for i, p in enumerate(xs):
        ov.append(anchor(*p, f"anchor {i} {STROKE_PTS[i]}", "hover" if i == 3 else "normal"))
    hx, hy = int(xs[3][0] // S), int(xs[3][1] // S)
    ov.append(tooltip_box(hx - tw("(330, 58)") - 22, hy - 22, [lambda x, y: T("(330, 58)", x, y, "amber")], "anchor coords tooltip", w=tw("(330, 58)") + 14))
    ov.append(cursor(hx + 2, hy + 2))
    overlay = G("selection (stroke line 10)", *ov)
    return window(lang, [toolbar("pencil", lang), tabbar(unsaved=True), layers_panel(lang, out, sel_line=10, props=props),
                         canvas_area(drawing("main", extra), overlay), status_main(lang, saved=False, extra_left=("pencil", "5 points from 212 samples", "mute"))])

BRUSH_PTS = [(70, 160), (160, 100), (250, 140), (340, 70), (430, 120)]
BRUSH_P = [0.2, 0.9, 0.5, 1.0, 0.3]

def brush_shape(points, pressures, mn, mx, taper, color, op=1.0, name="brush"):
    dense = catmull(points, 16)
    ws = pressure_widths(len(dense), pressures, mn, mx, taper)
    return f'<path id="{nm(name)}" d="{variable_outline(dense, ws)}" fill="{color}" fill-opacity="{op}"/>'

def pen_preview(kind, x, y, w, h, color):
    pts = [(x, y + h*0.7), (x + w*0.33, y + h*0.15), (x + w*0.66, y + h*0.85), (x + w, y + h*0.3)]
    if kind == "ink":
        return brush_shape(pts, [0.2, 1, 0.6, 0.2], 0.6, 4.5, 0.3, color, 1, "preview ink")
    if kind == "pencil":
        d = bezier_d(pts)
        return (f'<g id="preview pencil"><path d="{d}" fill="none" stroke="{color}" stroke-width="1.5" stroke-opacity="0.75" stroke-linecap="round"/>'
                f'<path d="{d}" fill="none" stroke="{C["row"]}" stroke-width="1.6" stroke-dasharray="0.6 2.2" stroke-opacity="0.6"/></g>')
    if kind == "marker":
        return f'<path id="preview marker" d="{bezier_d(pts)}" fill="none" stroke="{color}" stroke-width="6" stroke-opacity="0.55" stroke-linecap="butt"/>'
    if kind == "brush":
        dense = catmull(pts, 16)
        ws = pressure_widths(len(dense), [0.1, 1, 0.7, 0.15], 1, 10, 0.4)
        return f'<path id="preview brush" d="{variable_outline(dense, ws)}" fill="{color}" fill-opacity="0.85"/>'
    if kind == "calligraphy":
        dense = catmull(pts, 24)
        nx, ny = math.cos(math.radians(45)) * 3.5, -math.sin(math.radians(45)) * 3.5
        quads = []
        for a, b in zip(dense, dense[1:]):
            q = [(a[0]+nx, a[1]+ny), (b[0]+nx, b[1]+ny), (b[0]-nx, b[1]-ny), (a[0]-nx, a[1]-ny)]
            quads.append(poly_d(q, True))
        return f'<path id="preview calligraphy" d="{" ".join(quads)}" fill="{color}"/>'

def frame_brush():
    lang = "en"; t = L[lang]
    rows = file_rows()
    code_s = "brush([(70, 160), (160, 100), (250, 140), (340, 70), (430, 120)], [0.2, 0.9, 0.5, 1, 0.3], #3b2a6b, pens.ink)"
    new = dict(kind="brush", d=1, code=code_s, line=10)
    out = []
    for r in rows:
        if r.get("line") == 10:
            out.append(new); r = dict(r, line=11)
        elif r.get("line") == 11:
            r = dict(r, line=12)
        out.append(r)
    out = [r for r in out if r["kind"] != "add"]
    def f(y):
        g1 = G("field points (list)", T("points", 6, y + 2, "mute"),
               field(6 + tw("points") + 4, y, tw("5") + 8, 18, [T("5", 6 + tw("points") + 8, y + 3, "amber")], "input points"))
        x = 6 + tw("points") + 4 + tw("5") + 16
        g2, x = color_field(x, y, "color", "#3b2a6b", "color")
        pw = tw("pens.ink") + 20
        g3 = G("field pen (picker)", T("pen", x, y + 2, "mute"),
               field(x + tw("pen") + 4, y, pw, 18, [T("pens.ink", x + tw("pen") + 8, y + 3, "text"),
                     I("chevron-down", x + tw("pen") + 4 + tw("pens.ink") + 10, y + 5, "mute")], "input pen", focused=True))
        return G("fields (from brush(points, pressures, color, pen))", g1, g2, g3)
    props = props_panel(lang, f"brush · {t['line']} 10", f, "switching the pen re-styles this stroke")
    extra = brush_shape(BRUSH_PTS, BRUSH_P, 1, 10, 0.3, "#3b2a6b", 1, "brush line 10 (pens.ink)")
    xs = [cv(p) for p in BRUSH_PTS]
    ov = [sel_box(min(p[0] for p in xs) - 8, min(p[1] for p in xs) - 8, max(p[0] for p in xs) + 8, max(p[1] for p in xs) + 8)]
    for i, p in enumerate(xs):
        ov.append(anchor(*p, f"anchor {i}"))
    # pen picker popup under the brush tool button (toolbar x: tools at 3 + 15*n + seps)
    bx = 3 + 15*5 + 4 + 15  # brush button x
    px, py = bx, TOP_H + 1
    pens = [("ink", "crisp, tapers"), ("pencil", "light, grainy"), ("marker", "see-through"), ("brush", "soft, wide"), ("calligraphy", "45° nib")]
    pw, rh = 230, 20
    items = [R(px + 2, py + 2, pw, len(pens)*rh + 18, "void", "drop shadow"), frame_lines(px, py, pw, len(pens)*rh + 18, "void"),
             R(px + 1, py + 1, pw - 2, len(pens)*rh + 16, "row", "menu bg"), R(px + 1, py + 1, pw - 2, 13, "raised", "menu header"),
             T(t["pen"].upper() + " · pens.list", px + 5, py, "mute")]
    for i, (pn, desc) in enumerate(pens):
        ry = py + 16 + i * rh
        it = []
        if pn == "ink":
            it += [R(px + 1, ry, pw - 2, rh, "moss", "item bg (current)"), R(px + 1, ry, 2, rh, "lizard", "current bar")]
        if pn == "calligraphy":
            it += [R(px + 1, ry, pw - 2, rh, "raised", "item bg (hover)")]
        it.append(G("preview frame", R(px + 6, ry + 2, 50, 16, "coal", "preview bg"), frame_lines(px + 5, ry + 1, 52, 18, "void")))
        it.append(f'<g id="stroke preview (full-res)" transform="scale(0.5)">{pen_preview(pn, (px + 9)*S, (ry + 5)*S, 88, 16, C["text"])}</g>')
        it.append(T(pn, px + 62, ry + 1, "bright" if pn == "ink" else "text"))
        it.append(T(desc, px + 62, ry + 1 + 0, "mute", "desc") if False else "")
        it.append(T(desc, px + 62 + tw(pn) + 7, ry + 1, "mute", f"desc {desc}") if tw(pn) + tw(desc) + 7 < pw - 66 else T("", 0, 0))
        out_items = G(f"pen {pn}" + (" (current)" if pn == "ink" else "") + (" (hover)" if pn == "calligraphy" else ""), *it)
        items.append(out_items)
    picker = G("pen picker (open)", *items)
    ov.append("")
    overlay = G("selection (brush line 10)", *ov)
    hy = py + 16 + 4*rh + 6
    return window(lang, [toolbar("brush", lang), tabbar(unsaved=True), layers_panel(lang, out, sel_line=10, props=props),
                         canvas_area(drawing("main", extra), overlay), picker, cursor(px + 120, hy),
                         status_main(lang, saved=False, extra_left=("brush", "pressure 0.62 (Force Touch)", "mute"))])

def frame_tooltip_error():
    lang = "en"; t = L[lang]
    rows = file_rows()
    rows = [r for r in rows if r.get("line") not in (1, 2, 3, 4, 19, 20, 23)]   # scrolled past the header comments
    for r in rows:
        if r.get("line") == 18:
            r["open"], r["count"] = False, 6
        if r.get("line") == 26:
            r["code"] = 'rect((0, 0), (480, horizn), #f2a65a)'
    def errmsg(y):
        msg = "identifier not found: horizn"
        return G("error message (line 26)", R(0, y, PANEL_W, 14, "#2a1714", "error message bg") if False else R(0, y, 2, 16, "ember", "error bar"),
                 R(47, y + 2, 1, 12, "ember", "message tick"), T(msg, 56, y + 1, "ember", "error text: " + msg)), 16
    props = props_panel(lang, f"circle · {t['line']} 10", circle_props(), f"{t['line']} 10 · {t['hint']}")
    panel = layers_panel(lang, rows, hover_line=10, error_line=26, extra_after=(26, errmsg), props=props)
    # underline "56" in row 11: compute row y
    # rows y: header MAIN_Y+15; compute by walking
    y = MAIN_Y + 15
    yl = {}
    for r in rows:
        if r["kind"] == "blank": h = 4
        elif r["kind"] == "add": h = 10
        else: h = ROW_H
        if r["kind"] == "layer": y += LAYER_GAP
        yl[r.get("line")] = y
        y += h
    col = 54 + tw('rect((0, 0), (480, ')
    ul = G("error underline", *[R(col + i, yl[26] + 15 + (i // 2) % 2, 1, 1, "ember", "squiggle px") for i in range(tw('horizn'))])
    ty = yl[10] + 19
    tx = 60
    swatch_ex = 'circle((256, 256), 40, #3fa34d)'
    lines = [
        lambda x, y: G("title", T("circle", x, y, "bright"), T("(center, r, fill)", x + tw("circle"), y, "mute")),
        lambda x, y: T("Draws a filled circle.", x, y, "text"),
        lambda x, y: G("param center", T("center", x, y, "lizard"), T("the middle of the circle, a point", x + 56, y, "text")),
        lambda x, y: G("param r", T("r", x, y, "lizard"), T("the radius in pixels", x + 56, y, "text")),
        lambda x, y: G("param fill", T("fill", x, y, "lizard"), T("the color inside", x + 56, y, "text")),
        lambda x, y: G("example", T("example", x, y, "mute"), code(swatch_ex, x + 56, y, name="example code")),
    ]
    tip = tooltip_box(tx, ty, lines, "tooltip circle (docs)", w=56 + tw("the middle of the circle, a point") + 14)
    return window(lang, [toolbar(None, lang), tabbar(unsaved=True), panel, ul,
                         canvas_area(drawing("error"), "", dim=True, badge=("error", t["lastgood"], "ember")),
                         tip, cursor(58, yl[10] + 8), status_main(lang, err=True)])

def frame_spanish():
    return frame_main("es")

# ---------- splash ----------
def load_sprite(key, idx=None):
    src = open(os.path.join(ROOT, "lib/ui/splash_art.liz"), encoding="utf-8").read()
    i = src.index(f'"{key}":')
    if idx is not None:
        i = src.index("[", i)
        for _ in range(idx + 1):
            i = src.index('{"size"', i + 1)
    else:
        i = src.index('{"size"', i)
    m = re.match(r'\{"size": \((\d+), (\d+)\), "runs": \[(.*?)\]\}', src[i:])
    w, h = int(m.group(1)), int(m.group(2))
    runs = re.findall(r'\[(\d+), (\d+), (\d+), (#[0-9a-f]{6})\]', m.group(3))
    return w, h, [(int(a), int(b), int(c), d) for a, b, c, d in runs]

REMAP = {"#1b1a17": "void", "#2f6b3a": "moss", "#3fa34d": "lizard", "#8fd16a": "text", "#e0a526": "amber", "#a39e88": "mute"}
REMAP_NAMES = {"void": "outline", "moss": "shadow", "lizard": "fill", "text": "highlight", "amber": "eye"}

def sprite_svg(runs, x, y, px, name, remap=REMAP):
    by = {}
    for (ry, x0, x1, col) in runs:
        by.setdefault(remap.get(col, col), []).append(f"M{x + x0*px} {y + ry*px}h{(x1-x0)*px}v{px}h{-(x1-x0)*px}z")
    return G(name, *[f'<path id="{nm(name + " " + REMAP_NAMES.get(c, c) + " (" + c + ")")}" fill="{C.get(c,c)}" d="{"".join(ds)}"/>' for c, ds in by.items()])

def frame_splash():
    px = 4
    bw, bh = 200, 80
    lw, lh, logo = load_sprite("logo")
    zw, zh, liz = load_sprite("lizard", 2)
    W2, H2 = bw * px, bh * px
    parts = [checker(0, 0, W2, H2, "#2b2b2b", "#3a3a3a", 16, "transparency checker (not part of the art)")]
    lx, ly = 48, 14
    lzx = 112  # mid-crawl, behind the k/z
    parts.append(sprite_svg(liz, lzx * px, 36 * px, px, "lizard (frame 3 of 6, mid-crawl)"))
    parts.append(sprite_svg(logo, lx * px, ly * px, px, "logo Linkz"))
    tag = "powered by lizard-script"
    tx = (bw - tw(tag)) // 2
    # tagline at 2x font scale (font px = 4 screen px here: S=2, so draw at banner px 4 -> use own path)
    def tagpath(color, ox, oy):
        d = []
        for i, ch in enumerate(tag):
            for (gy, x0, x1) in FONT.get(ch, FALLBACK):
                d.append(f"M{(tx + i*GW + x0)*px + ox} {(58 + gy)*px + oy}h{(x1-x0)*px}v{px}h{-(x1-x0)*px}z")
        return "".join(d)
    parts.append(f'<path id="tagline shadow (void)" fill="{C["void"]}" d="{tagpath(0, px, px)}"/>')
    parts.append(f'<path id="tagline: powered by lizard-script (mute)" fill="{C["mute"]}" d="{tagpath(0, 0, 0)}"/>')
    return svg_doc(W2, H2, "".join(parts))

# ---------- style page pieces ----------
def L_rel(hexv):
    c = [int(hexv[i:i+2], 16)/255 for i in (1, 3, 5)]
    c = [x/12.92 if x <= 0.03928 else ((x+0.055)/1.055)**2.4 for x in c]
    return 0.2126*c[0] + 0.7152*c[1] + 0.0722*c[2]

def contrast(a, b):
    a, b = L_rel(a), L_rel(b)
    return (max(a, b) + 0.05) / (min(a, b) + 0.05)

def font_specimen():
    lines = ["ABCDEFGHIJKLMNOPQRSTUVWXYZ", "abcdefghijklmnopqrstuvwxyz", "0123456789 ()[]{}#,.;:+-*/=<>!?\"'", "áéíóúñü ¿¡ ç · … » ° ×  línea guardado",
             "circle((360, 90), 34, #fff3c4)"]
    w1 = max(tw(s) for s in lines) + 16
    h1 = len(lines) * 15 + 12
    def block(scale):
        items = [f'<rect id="bg" x="0" y="0" width="{w1*scale}" height="{h1*scale}" fill="{C["panel"]}"/>']
        for i, s in enumerate(lines):
            d = [text_d(s, 8, 6 + i*15, scale)]
            items.append(f'<path id="{nm("line: " + s)}" fill="{C["text"]}" d="{"".join(d)}"/>')
        return svg_doc(w1*scale, h1*scale, "".join(items))
    return block(1), block(2)

def comp_svgs():
    out = {}
    def wrap(w, h, body):
        return svg_doc(w*S, h*S, body, scale=S)
    out["button/normal"] = wrap(40, 16, button(0, 0, 40, 16, "button normal", icon="save", label="Save"))
    out["button/hover"] = wrap(40, 16, button(0, 0, 40, 16, "button hover", icon="save", label="Save", state="hover"))
    out["button/pressed"] = wrap(40, 16, button(0, 0, 40, 16, "button pressed", icon="save", label="Save", state="pressed"))
    out["tool/normal"] = wrap(14, 14, button(0, 0, 14, 14, "tool normal", icon="circle"))
    out["tool/hover"] = wrap(14, 14, button(0, 0, 14, 14, "tool hover", icon="circle", state="hover"))
    out["tool/active"] = wrap(14, 14, button(0, 0, 14, 14, "tool active", icon="pencil", state="active"))
    out["field/number"] = wrap(40, 18, field(0, 0, 40, 18, [T("360", 4, 3, "amber")], "field number"))
    out["field/focused"] = wrap(40, 18, field(0, 0, 40, 18, [T("360", 4, 3, "amber"), R(4 + tw("360"), 3, 1, 12, "amber", "caret")], "field number", True))
    out["field/color"] = wrap(80, 18, field(0, 0, 80, 18, [R(3, 4, 10, 10, "void", "swatch outline"), R(4, 5, 8, 8, "#fff3c4", "swatch"), T("#fff3c4", 16, 3, "text")], "field color"))
    r = dict(kind="circle", d=1, code="circle((360, 90), 34, #fff3c4)", line=10)
    for st in ["normal", "hover", "selected", "error"]:
        el, _ = row_el(r, 0, w=280, sel=st == "selected", hover=st == "hover", error=st == "error")
        out[f"row/{st}"] = wrap(280, ROW_H, R(0, 0, 280, ROW_H, "panel", "panel bg") + el)
    lr = dict(kind="layer", d=0, code="layer sun", line=9, open=True)
    el, _ = row_el(lr, 0, w=280)
    out["row/layer"] = wrap(280, ROW_H, R(0, 0, 280, ROW_H, "panel", "panel bg") + el)
    tip = tooltip_box(0, 0, [lambda x, y: G("title", T("circle", x, y, "bright"), T("(center, r, fill)", x + tw("circle"), y, "mute")),
                             lambda x, y: T("Draws a filled circle.", x, y, "text")], "tooltip", w=tw("circle(center, r, fill)") + 14)
    out["tooltip"] = wrap(tw("circle(center, r, fill)") + 16, 40, tip)
    return out

if __name__ == "__main__":
    frames = {
        "f1": frame_main("en"), "f2": frame_new_row(), "f3": frame_pencil(), "f4": frame_brush(),
        "f5": frame_tooltip_error(), "f6": frame_spanish(), "f7": frame_splash(),
    }
    for k, v in frames.items():
        open(os.path.join(OUT, k + ".svg"), "w").write(v)
    a, b = font_specimen()
    open(os.path.join(OUT, "font1.svg"), "w").write(a)
    open(os.path.join(OUT, "font2.svg"), "w").write(b)
    for k, v in comp_svgs().items():
        open(os.path.join(OUT, "comp_" + k.replace("/", "_") + ".svg"), "w").write(v)
    for n in ICON_ORDER:
        open(os.path.join(OUT, f"icon_{n}.svg"), "w").write(svg_doc(16, 16, f'<path id="{n}" fill="{C["text"]}" d="{icon_d(n, 0, 0, 2)}"/>'))
    for fg in ["mute", "text", "bright", "lizard", "amber", "ember", "sky", "peach", "dim"]:
        print(fg, " ".join(f"{bg}:{contrast(C[fg], C[bg]):.1f}" for bg in ["coal", "panel", "row", "raised", "moss"]))
    for k in frames: print(k, len(frames[k]))
