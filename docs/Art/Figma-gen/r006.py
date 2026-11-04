"""R006 App icon: three pixel-built options, master 1024 + small previews + Dock mocks.
Run:  python3 r006.py   (needs Pillow). Art is drawn on a 52x52 cell grid, 16 px per cell."""
import math, os
from PIL import Image, ImageDraw, ImageFilter

OUT = os.path.join(os.path.dirname(__file__), "..", "Reviews", "R006")
G, CELL, OFF = 52, 16, 96          # 832 px body (Big Sur grid ~824), centred on 1024
def hx(s): s = s.lstrip("#"); return tuple(int(s[i:i+2], 16) for i in (0, 2, 4)) + (255,)
VOID, COAL, PANEL, ROW, RAISED, BEVEL = map(hx, "0b0c0b 131512 1a1d19 232721 2e332c 434a40".split())
MUTE, TEXT, BRIGHT, MOSS, LIZARD, AMBER, SKY = map(hx, "a0ab9b d4dacd f2f5ec 163a20 4fc46a f2b33d 72b8e6".split())
INK, INK_HI, INK_LO, INK_GL = map(hx, "3b40a8 a8b0ff 1c1f5c eef0ff".split())   # splash logo ink, art colors
LIME = hx("a6e8a0")

def layer(): return Image.new("RGBA", (G, G), (0, 0, 0, 0))
def px(im, x, y, c):
    if 0 <= x < G and 0 <= y < G: im.putpixel((int(x), int(y)), c)
def disc(im, cx, cy, r, c):
    for y in range(G):
        for x in range(G):
            if (x + .5 - cx) ** 2 + (y + .5 - cy) ** 2 <= r * r: px(im, x, y, c)
def poly(im, pts, c): ImageDraw.Draw(im).polygon(pts, fill=c)
def rect(im, x0, y0, x1, y1, c): ImageDraw.Draw(im).rectangle([x0, y0, x1 - 1, y1 - 1], fill=c)
def outline(im, c=VOID):
    """1 cell 8-neighbour outline behind the sprite, like the splash logo."""
    a = im.getchannel("A").point(lambda v: 255 if v else 0).filter(ImageFilter.MaxFilter(3))
    o = Image.new("RGBA", (G, G), c); o.putalpha(a)
    o.alpha_composite(im); return o
def bez(p, t):
    a, b, c, d = p; u = 1 - t
    return tuple(u**3*a[i] + 3*u*u*t*b[i] + 3*u*t*t*c[i] + t**3*d[i] for i in (0, 1))

def background(kind):
    im = Image.new("RGBA", (G, G), PANEL)
    for y in range(G):
        for x in range(G):
            d = math.hypot(x - 25.5, y - 25.5)
            if d > 33 and (x + y) % 2 == 0: im.putpixel((x, y), COAL)      # 2x2-style checker dither toward the corners
            elif d > 29 and (x % 2 == 0 and y % 2 == 0): im.putpixel((x, y), COAL)
    return im
def bevel_frame(im):
    for i in range(G):
        px(im, i, 0, BEVEL); px(im, 0, i, BEVEL); px(im, i, G - 1, VOID); px(im, G - 1, i, VOID)
    px(im, 0, G - 1, VOID); px(im, G - 1, 0, BEVEL)

# ---------- A: the nib with a lizard eye ----------
def option_a():
    bg = background("a"); n = layer()
    poly(n, [(21, 5), (31, 5), (34, 10), (38, 21), (31, 32), (26, 42), (21, 32), (14, 21), (18, 10)], LIZARD)
    for y in range(G):                                   # right side in moss, dithered edge
        for x in range(G):
            if n.getpixel((x, y))[3]:
                if x >= 29 + (y - 5) // 9 and x > 28 and (x >= 31 or (x + y) % 2 == 0): n.putpixel((x, y), MOSS)
                elif x <= 16 and y > 8 and y < 26 and (x == 16 or (x == 17 and y % 2 == 0)): n.putpixel((x, y), LIME)
    for x in range(22, 31): px(n, x, 6, LIME)            # top sheen
    for y in range(22, 43): px(n, 26, y, VOID)           # the slit, from the eye to the point
    for y in range(10, 14): px(n, 17 + (y - 10) // 2, y, LIME)
    for x in range(21, 32):                              # collar line under the top edge
        px(n, x, 8, MOSS if x > 26 else LIME)
    n = outline(n)
    eye = layer(); disc(eye, 26, 17.5, 4.5, VOID); disc(eye, 26, 17.5, 3.5, AMBER)
    rect(eye, 26, 14, 27, 21, VOID)                      # vertical lizard pupil
    px(eye, 24, 16, BRIGHT); px(eye, 24, 17, BRIGHT)
    n.alpha_composite(eye)
    # tip band: nib point stripe
    pass
    drop = layer()                                       # ink drop under the tip
    poly(drop, [(26, 43), (28, 46), (28, 48), (26, 49), (24, 48), (24, 46)], INK)
    disc(drop, 26, 47, 2.6, INK)
    px(drop, 25, 46, INK_HI); px(drop, 25, 45, INK_GL); px(drop, 27, 48, INK_LO)
    drop = outline(drop)
    bg.alpha_composite(n); bg.alpha_composite(drop)
    # three tiny ink dots flying off the nib
    for (x, y) in [(10, 40), (42, 38), (8, 32)]: px(bg, x, y, INK); px(bg, x, y - 1, INK_HI)
    return bg

# ---------- B: the tail is a Bezier path ----------
def option_b():
    bg = background("b")
    for y in range(8, 46):                               # faint 4-cell grid of dots, an editor canvas
        for x in range(6, 46, 4):
            if y % 4 == 0: px(bg, x, y, ROW)
    seg1 = [(12, 15), (32, 4), (47, 22), (33, 31)]
    seg2 = [(33, 31), (23, 36), (14, 44), (24, 45)]      # tail curls in
    pts = [(bez(seg1, i / 60), 1 - i / 120) for i in range(61)] + [(bez(seg2, i / 60), .5 - i / 120) for i in range(61)]
    body = layer()
    for (x, y), t in pts:
        r = 1.0 + 3.6 * t
        disc(body, x, y, r, LIZARD)
    for i, ((x, y), t) in enumerate(pts):                # moss belly shade on lower/right side
        pass
    shade = layer()
    for (x, y), t in pts:
        disc(shade, x + 1.2, y + 1.4, 1.0 + 3.6 * t - 1.0, MOSS)
    sh = Image.new("RGBA", (G, G), (0, 0, 0, 0)); sh.alpha_composite(shade)
    mask = body.getchannel("A")
    sh.putalpha(Image.composite(sh.getchannel("A"), Image.new("L", (G, G), 0), mask))
    out = layer(); out.alpha_composite(body)
    for y in range(G):                                   # keep only moss where it's on the body edge zone
        for x in range(G):
            if sh.getpixel((x, y))[3] and body.getpixel((x, y))[3] and (x + y) % 2 == 0 and not body.getpixel((x - 1, y - 1))[3] == 0:
                pass
    # simple light/dark: top-left lime, bottom-right moss, by checking neighbours
    res = layer()
    for y in range(G):
        for x in range(G):
            if body.getpixel((x, y))[3]:
                up = body.getpixel((x, y - 1))[3] if y else 0
                lf = body.getpixel((x - 1, y))[3] if x else 0
                dn = body.getpixel((x, y + 1))[3] if y < G - 1 else 0
                rt = body.getpixel((x + 1, y))[3] if x < G - 1 else 0
                c = LIZARD
                if not up or not lf: c = LIME
                if not dn or not rt: c = MOSS
                res.putpixel((x, y), c)
    res = outline(res)
    # head at the start
    head = layer()
    disc(head, 10.5, 15.5, 5.6, LIZARD); disc(head, 6.5, 17, 3.6, LIZARD)
    for y in range(G):
        for x in range(G):
            if head.getpixel((x, y))[3] and (y >= 19 or x >= 14): head.putpixel((x, y), MOSS)
    px(head, 8, 12, LIME); px(head, 9, 12, LIME); px(head, 10, 12, LIME)
    head = outline(head)
    rect(head, 9, 13, 12, 16, VOID); rect(head, 9, 13, 11, 15, AMBER); px(head, 10, 14, VOID)
    px(head, 3, 19, hx("ef5d4c")); px(head, 2, 20, hx("ef5d4c"))   # tongue flick
    res.alpha_composite(head)
    # legs
    legs = layer()
    for (x0, y0, x1, y1) in [(15, 18, 12, 24), (21, 14, 24, 9)]:
        for i in range(8):
            px(legs, x0 + (x1 - x0) * i / 7, y0 + (y1 - y0) * i / 7, LIZARD); px(legs, x0 + (x1 - x0) * i / 7 + 1, y0 + (y1 - y0) * i / 7, LIZARD)
    bg.alpha_composite(outline(legs)); bg.alpha_composite(res)
    # Bezier handles: tail segment end, sky lines + circles, bright anchor
    ov = layer(); d = ImageDraw.Draw(ov)
    def handle(a, h):
        n = int(max(abs(h[0]-a[0]), abs(h[1]-a[1])))
        for i in range(n + 1):
            if i % 2 == 0: px(ov, a[0] + (h[0]-a[0]) * i / n, a[1] + (h[1]-a[1]) * i / n, SKY)
    handle((33, 31), (43, 38)); handle((33, 31), (24, 25))
    for h in [(43, 38), (24, 25)]:
        disc(ov, h[0] + .5, h[1] + .5, 2.5, VOID); disc(ov, h[0] + .5, h[1] + .5, 1.6, SKY)
    rect(ov, 30, 28, 37, 35, VOID); rect(ov, 31, 29, 36, 34, BRIGHT)    # the anchor, a 5x5 square
    rect(ov, 33, 31, 35, 33, AMBER)
    bg.alpha_composite(ov)
    return bg

# ---------- C: lizard on an ink drop ----------
def option_c():
    bg = background("c"); dr = layer()
    cx, cy, r = 26, 31, 15
    disc(dr, cx, cy, r, INK)
    poly(dr, [(26, 4), (26 - 9.5, 22), (26 + 9.5, 22)], INK)
    poly(dr, [(26, 4), (26 - 5, 14), (26 + 5, 14)], INK)
    res = layer()
    for y in range(G):
        for x in range(G):
            if dr.getpixel((x, y))[3]:
                c = INK
                if x + y * .35 >= 36 + ((y % 2) if False else 0) and x >= 31: c = INK_LO
                if x >= 35 or (x == 34 and (x + y) % 2 == 0): c = INK_LO
                res.putpixel((x, y), c)
    for (x, y) in [(17, 26), (17, 27), (17, 28), (18, 25), (19, 24), (18, 26)]: px(res, x, y, INK_HI)
    px(res, 18, 25, INK_GL); px(res, 17, 27, INK_GL); px(res, 24, 13, INK_HI); px(res, 24, 14, INK_HI); px(res, 24, 12, INK_GL)
    res = outline(res)
    bg.alpha_composite(res)
    # lizard, side view facing right, standing on the drop (diagonal-pair legs)
    z = layer()
    for i in range(0, 34):                               # tail: S curve ending in a curl
        t = i / 33; x = 20 - 12 * t; y = 31 + 4 * math.sin(t * 3.2) - 7 * t * t
        disc(z, x, y, 3.4 * (1 - t) + .8, LIZARD)
    disc(z, 26, 30.5, 6.2, LIZARD); disc(z, 31, 30, 5.0, LIZARD)    # body
    disc(z, 37, 29, 4.2, LIZARD); rect(z, 37, 27, 43, 31, LIZARD)    # head + snout
    for (x0, y0, x1, y1) in [(23, 33, 20, 39), (33, 33, 37, 39)]:    # legs
        for i in range(7):
            for dx in (0, 1, 2): px(z, x0 + (x1-x0)*i/6 + dx, y0 + (y1-y0)*i/6, LIZARD)
    res = layer()
    for y in range(G):
        for x in range(G):
            if z.getpixel((x, y))[3]:
                up = z.getpixel((x, y - 1))[3] if y else 0
                dn = z.getpixel((x, y + 1))[3] if y < G - 1 else 0
                res.putpixel((x, y), LIME if not up else (MOSS if not dn else LIZARD))
    for i in range(7):                                   # belly dither
        for x in range(21, 34):
            if res.getpixel((x, 33))[3] and res.getpixel((x, 33)) == LIZARD and (x % 2 == 0): res.putpixel((x, 33), MOSS)
    rect(res, 38, 26, 41, 29, AMBER); px(res, 40, 27, VOID); px(res, 40, 28, VOID)
    px(res, 44, 30, hx("ef5d4c")); px(res, 45, 31, hx("ef5d4c"))
    res = outline(res)
    bg.alpha_composite(res)
    return bg

# ---------- master + previews ----------
def squircle_mask(size=1024, body=832, n=5.0, aa=4):
    big = Image.new("L", (size * aa, size * aa), 0); px_ = big.load()
    c = size * aa / 2; a = body * aa / 2
    for y in range(size * aa):
        for x in range(size * aa):
            if abs((x + .5 - c) / a) ** n + abs((y + .5 - c) / a) ** n <= 1: px_[x, y] = 255
    return big.resize((size, size), Image.LANCZOS)
MASK = None
def master(art):
    global MASK
    if MASK is None: MASK = squircle_mask()
    bevel_frame(art)
    body = art.resize((G * CELL, G * CELL), Image.NEAREST)
    canvas = Image.new("RGBA", (1024, 1024), (0, 0, 0, 0)); canvas.paste(body, (OFF, OFF))
    canvas.putalpha(MASK); return canvas

def dock(icons, bg, label_col):
    """A mock Dock strip: rounded translucent bar with the icon at 128, 64, 32, 16 and a 4x zoom of the 16."""
    W, H = 560, 180; im = Image.new("RGBA", (W, H), bg)
    bar = Image.new("RGBA", (W, H), (0, 0, 0, 0)); d = ImageDraw.Draw(bar)
    barc = (255, 255, 255, 120) if sum(bg[:3]) > 400 else (255, 255, 255, 36)
    d.rounded_rectangle([14, 30, W - 14, H - 14], 26, fill=barc, outline=(255, 255, 255, 90) if sum(bg[:3]) > 400 else (255, 255, 255, 50))
    im.alpha_composite(bar)
    x = 34
    for s in (128, 64, 32, 16):
        ic = icons.resize((s, s), Image.LANCZOS); im.alpha_composite(ic, (x, 30 + (H - 44 - s) // 2 + 0)); x += s + 26
    z = icons.resize((16, 16), Image.LANCZOS).resize((64, 64), Image.NEAREST); im.alpha_composite(z, (W - 14 - 64 - 14, 30 + (H - 44 - 64) // 2))
    return im

def main():
    os.makedirs(OUT, exist_ok=True)
    opts = {"A-nib": option_a(), "B-bezier-tail": option_b(), "C-ink-drop": option_c()}
    sheet = []
    for name, art in opts.items():
        m = master(art); m.save(f"{OUT}/{name}-1024.png")
        for s in (16, 32, 64, 128): m.resize((s, s), Image.LANCZOS).save(f"{OUT}/{name}-{s}.png")
        light = dock(m, (226, 228, 232, 255), None); dark = dock(m, (22, 24, 28, 255), None)
        light.save(f"{OUT}/{name}-dock-light.png"); dark.save(f"{OUT}/{name}-dock-dark.png")
        sheet.append((m, light, dark))
    W = 1024 + 20 + 560; S = Image.new("RGBA", (W + 40, 3 * 1044 + 20), (11, 12, 11, 255))
    for i, (m, l, d) in enumerate(sheet):
        y = 20 + i * 1044; S.alpha_composite(m, (20, y)); S.alpha_composite(l, (1064, y + 300)); S.alpha_composite(d, (1064, y + 520))
    S.resize((S.width // 2, S.height // 2), Image.LANCZOS).save(f"{OUT}/overview.png")
if __name__ == "__main__":
    main()
