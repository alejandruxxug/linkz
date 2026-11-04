"""R006 round 2: a lizard FACE with a paintbrush (developer: "a lizard face with a paintbrush on the side").
Three variations, each with a master and a hand-simplified 16 px sprite. Run: python3 r006b.py (needs Pillow)."""
import math, os
from PIL import Image, ImageDraw
import r006 as R
from r006 import *

PEACH, WOOD_HI = hx("eb9a6f"), hx("f7c9a2")
def ell(im, cx, cy, rx, ry, c):
    for y in range(R.G):
        for x in range(R.G):
            if ((x + .5 - cx) / rx) ** 2 + ((y + .5 - cy) / ry) ** 2 <= 1: px(im, x, y, c)
def stamp_line(im, a, b, w, c):
    n = int(max(abs(b[0]-a[0]), abs(b[1]-a[1])) * 2) + 1
    for i in range(n + 1):
        t = i / n; disc(im, a[0] + (b[0]-a[0]) * t, a[1] + (b[1]-a[1]) * t, w, c)
def eye(im, cx, cy, r, pupil=True):
    disc(im, cx, cy, r + 1.1, VOID); disc(im, cx, cy, r, AMBER)
    if pupil: rect(im, int(cx - .5), int(cy - r + .5), int(cx - .5) + 1, int(cy + r - .5), VOID)
    px(im, int(cx - r + 1), int(cy - r + 1), BRIGHT)
def shade(im, cx, cy, rx, ry, hi=.62, lo=.85):
    """bottom-right crescent in moss (dithered edge) and a lime rim top-left, only on lizard-green pixels"""
    for y in range(R.G):
        for x in range(R.G):
            if im.getpixel((x, y)) == LIZARD:
                u = ((x + .5 - cx) / rx + (y + .5 - cy) / ry) * .7
                if u > lo or (u > hi and (x + y) % 2 == 0): im.putpixel((x, y), MOSS)
                elif u < -lo: im.putpixel((x, y), LIME)
def brush(im, p0, p1, hr=2.2, br=3.3):
    L = math.hypot(p1[0]-p0[0], p1[1]-p0[1]); ux, uy = (p1[0]-p0[0]) / L, (p1[1]-p0[1]) / L; nx, ny = -uy, ux
    n = int(L * 3)
    for i in range(n + 1):
        t = i / n; x, y = p0[0] + ux * L * t, p0[1] + uy * L * t
        if t < .56: disc(im, x, y, hr, PEACH)
        elif t < .66: disc(im, x, y, hr + .9, MUTE)
        else:
            k = (t - .66) / .34; disc(im, x, y, br * (1 - k ** 2.4) + .9, INK)
    for i in range(int(L * .5)):                         # wood sheen
        t = i / (L * .5) * .53; px(im, p0[0] + ux*L*t - nx * hr * .55, p0[1] + uy*L*t - ny * hr * .55, WOOD_HI)
    for t in (.7, .76, .82):                             # bristle sheen
        px(im, p0[0] + ux*L*t - nx * 1.1, p0[1] + uy*L*t - ny * 1.1, INK_HI)
    px(im, p0[0] + ux*L*.72 - nx * 1.1, p0[1] + uy*L*.72 - ny * 1.1, INK_GL)
def inkdrop(im, x, y, big=True):
    s = layer()
    if big:
        poly(s, [(x, y - 3), (x + 2, y), (x + 2, y + 2), (x, y + 4), (x - 2, y + 2), (x - 2, y)], INK); disc(s, x, y + 2, 2.4, INK)
    else: rect(s, x - 1, y - 1, x + 1, y + 2, INK)
    px(s, x - 1, y + 1, INK_HI); px(s, x - 1, y, INK_GL)
    im.alpha_composite(outline(s))

# ---------- D: friendly front face, brush tucked behind the ear ----------
def option_d():
    bg = background("d"); b = layer()
    brush(b, (36, 47), (45, 9), hr=2.3, br=3.6)
    bg.alpha_composite(outline(b))
    h = layer()
    ell(h, 21, 30, 17, 14.5, LIZARD); ell(h, 11, 17, 6.4, 6, LIZARD); ell(h, 31, 17, 6.4, 6, LIZARD)
    for (x0, xa, x1, yt) in [(15, 18, 21, 10), (21, 24, 27, 9)]: poly(h, [(x0, 16), (xa, yt), (x1, 16)], LIZARD)   # little crest
    shade(h, 21, 30, 17, 14)
    for x in range(15, 28): px(h, x, 16, LIME)
    for (x, y) in [(15, 41), (20, 43), (25, 43), (29, 41)]: px(h, x, y, MOSS)
    h = outline(h)
    eye(h, 11, 17, 4.2); eye(h, 31, 17, 4.2)
    for (x, y) in [(18, 28), (18, 29), (24, 28), (24, 29)]: px(h, x, y, VOID)        # nostrils
    for x in range(10, 33):                                                         # smile
        px(h, x, 34 + round(((x - 21) / 11.5) ** 2 * 3.5), VOID)
    px(h, 9, 32, VOID); px(h, 33, 32, VOID)
    for (x, y) in [(6, 28), (8, 29), (7, 29), (34, 28), (36, 29), (35, 29)]: px(h, x, y, PEACH)   # blush
    bg.alpha_composite(h)
    d = layer(); inkdrop(d, 47, 4, False)
    return bg

# ---------- E: cool sharp three-quarter face, brush held in the teeth ----------
def option_e(far_eye=True):
    bg = background("e"); n = layer()
    poly(n, [(6, 32), (16, 36), (24, 44), (20, 51), (4, 51), (3, 40)], LIZARD)      # neck + dewlap
    for y in range(44, 51):
        for x in range(0, 30):
            if n.getpixel((x, y))[3] and (x + y) % 2 == 0: n.putpixel((x, y), MOSS)
    sp = layer()
    for (x0, x1, yt) in [(8, 13, 7), (14, 20, 3), (21, 27, 2), (28, 33, 4)]:        # crest spikes
        poly(sp, [(x0, 14), ((x0 + x1) // 2, yt), (x1, 14)], MOSS)
    bg.alpha_composite(outline(sp)); bg.alpha_composite(outline(n))
    h = layer()
    poly(h, [(7, 22), (11, 13), (22, 10), (34, 12), (43, 18), (49, 26), (47, 31), (37, 33), (31, 40), (21, 41), (11, 35)], LIZARD)
    shade(h, 27, 26, 21, 15, .55, .8)
    for x in range(14, 30): px(h, x, 11, LIME)
    stamp_line(h, (22, 15), (38, 21), 1.7, MOSS)                                     # brow ridge, slanting to the nose
    h = outline(h)
    eye(h, 30, 24, 4.2)
    if far_eye: eye(h, 16.5, 22, 2.6, pupil=False); rect(h, 16, 20, 17, 25, VOID)
    px(h, 46, 25, VOID); px(h, 46, 26, VOID)
    for x in range(26, 47): px(h, x, 30 + (x - 26) // 12, VOID)                      # mouth line
    poly(h, [(31, 31), (33, 31), (32, 35)], BRIGHT); poly(h, [(38, 31), (40, 32), (39, 35)], BRIGHT)   # fangs
    bg.alpha_composite(h)
    b = layer(); brush(b, (33, 32), (44, 44), hr=2.0, br=3.4)
    bg.alpha_composite(outline(b)); return bg

# ---------- F: front face, a clawed hand holds the brush up ----------
def option_f():
    bg = background("f"); h = layer()
    poly(h, [(20, 16), (42, 16), (47, 22), (47, 35), (41, 44), (21, 44), (16, 35), (16, 22)], LIZARD)
    for (x0, xa, x1, yt) in [(21, 24, 27, 9), (28, 31.5, 35, 6), (36, 39, 42, 9)]: poly(h, [(x0, 17), (xa, yt), (x1, 17)], LIZARD)   # crest
    shade(h, 31.5, 30, 16, 14, .6, .85)
    for x in range(21, 42): px(h, x, 17, LIME)
    h = outline(h)
    stamp_line(h, (19, 19), (27, 22), 1.0, VOID); stamp_line(h, (44, 19), (36, 22), 1.0, VOID)      # brows
    eye(h, 24, 26, 4.2); eye(h, 39, 26, 4.2)
    for (x, y) in [(29, 33), (34, 33)]: px(h, x, y, VOID); px(h, x, y + 1, VOID)
    for x in range(21, 42): px(h, x, 39, VOID)                                       # wide flat grin
    px(h, 20, 38, VOID); px(h, 42, 38, VOID)
    for x in range(24, 39, 3): rect(h, x, 40, x + 2, 42, BRIGHT)                     # teeth
    bg.alpha_composite(h)
    b = layer(); brush(b, (9, 46), (9, 14), hr=2.1, br=3.4)
    inkdrop(b, 9, 7, True)
    bg.alpha_composite(outline(b))
    c = layer()                                            # the claw: palm + three fingers wrapped on the handle
    disc(c, 10, 37, 3.6, LIZARD)
    for y in (34, 37, 40): stamp_line(c, (10, y), (14, y), 1.2, LIZARD); px(c, 15, y, LIME)
    px(c, 8, 35, LIME); px(c, 9, 35, LIME)
    bg.alpha_composite(outline(c)); return bg

# ---------- hand-simplified 16 px sprites (14x14 art cells, 1 px margin) ----------
def small(kind):
    R.G = 14; S = Image.new("RGBA", (14, 14), PANEL)
    if kind == "d":
        b = layer(); brush(b, (1, 13), (10.5, 2), hr=1.1, br=1.6); S.alpha_composite(outline(b))
        h = layer(); ell(h, 7, 8.5, 6, 4.8, LIZARD); ell(h, 3.8, 5, 2.4, 2.2, LIZARD); ell(h, 10.2, 5, 2.4, 2.2, LIZARD)
        h = outline(h)
        for x in (3, 10): rect(h, x, 4, x + 2, 6, AMBER); px(h, x + 1, 5, VOID)
        for x in range(4, 10): px(h, x, 11, MOSS)
        px(h, 4, 10, VOID); px(h, 9, 10, VOID)
        S.alpha_composite(h)
    if kind == "e":
        h = layer(); poly(h, [(1, 6), (3, 3), (8, 2), (12, 5), (13, 8), (9, 10), (7, 12), (3, 11)], LIZARD)
        for x in (3, 5, 7): px(h, x, 2, MOSS)
        h = outline(h); rect(h, 7, 5, 10, 8, AMBER); px(h, 8, 6, VOID); px(h, 7, 4, VOID); px(h, 8, 4, VOID); px(h, 9, 5, VOID)
        for x in range(8, 13): px(h, x, 9, VOID)
        S.alpha_composite(h)
        b = layer(); brush(b, (7, 9), (12, 13), hr=1.0, br=1.5); S.alpha_composite(outline(b))
    if kind == "f":
        h = layer(); poly(h, [(4, 3), (10, 3), (12, 5), (12, 10), (10, 12), (4, 12), (2, 10), (2, 5)], LIZARD)
        h = outline(h)
        for x in (3, 8): rect(h, x, 5, x + 2, 7, AMBER); px(h, x + 1, 6, VOID)
        for x in range(4, 10): px(h, x, 10, VOID)
        S.alpha_composite(h)
        b = layer(); brush(b, (1, 13), (1, 2), hr=0.9, br=1.4); px(b, 1, 1, INK); S.alpha_composite(outline(b))
    R.G = 52
    sm = Image.new("RGBA", (16, 16), (0, 0, 0, 0)); sm.paste(S, (1, 1))
    sm.putalpha(squircle_mask(16, 14, 5.0, 16)); return sm

def dock(m, s16, bgc):
    W, H = 560, 180; im = Image.new("RGBA", (W, H), bgc); light = sum(bgc[:3]) > 400
    bar = Image.new("RGBA", (W, H), (0, 0, 0, 0))
    ImageDraw.Draw(bar).rounded_rectangle([14, 30, W - 14, H - 14], 26, fill=(255, 255, 255, 120 if light else 36), outline=(255, 255, 255, 90 if light else 50))
    im.alpha_composite(bar); x = 34
    for s in (128, 64, 32, 16):
        ic = s16 if s == 16 else m.resize((s, s), Image.LANCZOS); im.alpha_composite(ic, (x, 30 + (H - 44 - s) // 2)); x += s + 26
    im.alpha_composite(s16.resize((64, 64), Image.NEAREST), (W - 92, 30 + (H - 44 - 64) // 2)); return im

def main():
    out = R.OUT; os.makedirs(out, exist_ok=True)
    opts = {"D-friendly-brush-behind-ear": (option_d, "d"), "E-cool-brush-in-teeth": (option_e, "e"), "F-claw-holds-brush": (option_f, "f")}
    sheet = []
    for name, (fn, k) in opts.items():
        m = R.master(fn()); m.save(f"{out}/{name}-1024.png")
        s16 = small(k); s16.save(f"{out}/{name}-16.png")
        for s in (32, 64, 128): m.resize((s, s), Image.LANCZOS).save(f"{out}/{name}-{s}.png")
        l, d = dock(m, s16, (226, 228, 232, 255)), dock(m, s16, (22, 24, 28, 255))
        l.save(f"{out}/{name}-dock-light.png"); d.save(f"{out}/{name}-dock-dark.png")
        S = Image.new("RGBA", (128 + 64 + 32 + 16 + 100, 168), (0, 0, 0, 0)); x = 20
        for s, ic in ((128, Image.open(f"{out}/{name}-128.png")), (64, Image.open(f"{out}/{name}-64.png")), (32, Image.open(f"{out}/{name}-32.png")), (16, s16)):
            S.alpha_composite(ic, (x, 20)); x += s + 20
        S.save(f"{out}/{name}-small-sizes.png"); sheet.append((m, l, d))
    T = Image.new("RGBA", (1064 + 560 + 40, 3 * 1044 + 20), (11, 12, 11, 255))
    for i, (m, l, d) in enumerate(sheet):
        y = 20 + i * 1044; T.alpha_composite(m, (20, y)); T.alpha_composite(l, (1064, y + 300)); T.alpha_composite(d, (1064, y + 520))
    T.resize((T.width // 2, T.height // 2), Image.LANCZOS).save(f"{out}/overview-2.png")
if __name__ == "__main__":
    main()
