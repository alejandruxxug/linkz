import re, zlib, struct, sys, os
sys.path.insert(0, os.path.dirname(__file__))
import gen
src = open(os.path.join(gen.ROOT, "lib/ui/splash_art.liz"), encoding="utf-8").read()
sprites = []
for m in re.finditer(r'\{"size": \((\d+), (\d+)\), "runs": \[(.*?)\]\}', src):
    runs = [(int(a), int(b), int(c), d) for a, b, c, d in re.findall(r'\[(\d+), (\d+), (\d+), (#[0-9a-f]{6})\]', m.group(3))]
    sprites.append((int(m.group(1)), int(m.group(2)), runs))
logo, lizards, vines, bloom = sprites[0], sprites[1:7], sprites[7:31], sprites[31:]
print("sprites", len(sprites), "vines", len(vines), "bloom", len(bloom))
BW, BH, PX = 200, 80, 4
img = [[None]*BW for _ in range(BH)]
def put(sp, ox, oy, remap={}):
    for (y, x0, x1, col) in sp[2]:
        for x in range(x0, x1):
            X, Y = ox + x, oy + y
            if 0 <= X < BW and 0 <= Y < BH:
                img[Y][X] = remap.get(col, col)
INK = {"#1b1a17": gen.C["void"]}
put(lizards[2], 132, 36, INK)
put(logo, 48, 14, INK)
put(vines[-1], 0, 0, INK)
put(bloom[-1], 0, 0, INK)
tag = "powered by lizard-script"
tx = (BW - gen.tw(tag)) // 2
for col, o in ((gen.C["void"], 1), (gen.C["mute"], 0)):
    for i, ch in enumerate(tag):
        for (gx, gy, w, h) in gen.glyph_rects(ch):
            for yy in range(h):
                for xx in range(w):
                    X, Y = tx + i*7 + gx + xx + o, 58 + gy + yy + o
                    if 0 <= X < BW and 0 <= Y < BH: img[Y][X] = col
def png(path, rows, scale):
    h, w = len(rows) * scale, len(rows[0]) * scale
    raw = bytearray()
    for r in rows:
        line = bytearray([0])
        for c in r:
            px = bytes([0, 0, 0, 0]) if c is None else bytes([int(c[1:3], 16), int(c[3:5], 16), int(c[5:7], 16), 255])
            line += px * scale
        for _ in range(scale): raw += line
    def chunk(t, d): return struct.pack(">I", len(d)) + t + d + struct.pack(">I", zlib.crc32(t + d) & 0xffffffff)
    open(path, "wb").write(b"\x89PNG\r\n\x1a\n" + chunk(b"IHDR", struct.pack(">IIBBBBB", w, h, 8, 6, 0, 0, 0)) + chunk(b"IDAT", zlib.compress(bytes(raw), 9)) + chunk(b"IEND", b""))
png(os.path.join(gen.OUT, "splash-4x.png"), img, 4)
png(os.path.join(gen.OUT, "splash-1x.png"), img, 1)
