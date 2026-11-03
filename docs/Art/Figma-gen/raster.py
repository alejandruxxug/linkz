# Splits a generated frame SVG into named layers: pixel UI parts become exact PNGs
# (cropped, positioned), smooth vector parts (the drawing, previews, thumbnails)
# stay SVG. Writes out/layers/<key>/NN.png|svg and an index.json the importer reads.
import sys, os, re, json, zlib, struct
import xml.etree.ElementTree as ET
NS = "{http://www.w3.org/2000/svg}"
ET.register_namespace("", NS[1:-1])

DESCEND = {"layers panel (the code)", "rows", "canvas area", "grid all drawings", "row working on",
           "ui"}

def parse_tf(s):
    sx, tx, ty = 1.0, 0.0, 0.0
    if not s: return (1.0, 0.0, 0.0)
    m = re.search(r"translate\(([-\d.]+)[ ,]+([-\d.]+)\)", s)
    if m: tx, ty = float(m.group(1)), float(m.group(2))
    m = re.search(r"scale\(([-\d.]+)\)", s)
    if m: sx = float(m.group(1))
    if re.search(r"rotate|matrix|skew", s): return None
    return (sx, tx, ty)

def compose(a, b):  # apply b inside a: p -> a(b(p))
    return (a[0]*b[0], a[1] + a[0]*b[1], a[2] + a[0]*b[2])

def hexrgb(c):
    c = c.strip()
    return int(c[1:3], 16), int(c[3:5], 16), int(c[5:7], 16)

RECTPATH = re.compile(r"^[MmHhVvZz\d\s.\-]*$")

def is_raster(e, tf):
    tag = e.tag.replace(NS, "")
    if tf is None: return False
    if tag == "rect": return True
    if tag == "path": return bool(RECTPATH.match(e.get("d", ""))) and e.get("stroke") is None
    if tag == "g":
        t = parse_tf(e.get("transform"))
        if t is None or (e.get("transform") and "translate" in (e.get("transform") or "")): return False
        return all(is_raster(k, compose(tf, t)) for k in e)
    return False

def path_rects(d):
    toks = re.findall(r"[MmHhVvZz]|-?\d+(?:\.\d+)?", d)
    i, x, y, sx, sy = 0, 0.0, 0.0, 0.0, 0.0
    pts, rects = [], []
    cmd = None
    while i < len(toks):
        t = toks[i]
        if t in "MmHhVvZz":
            cmd = t; i += 1
            if t in "Zz":
                if pts:
                    xs = [p[0] for p in pts]; ys = [p[1] for p in pts]
                    rects.append((min(xs), min(ys), max(xs), max(ys)))
                pts = []; x, y = sx, sy
            continue
        if cmd in "Mm":
            a, b = float(toks[i]), float(toks[i+1]); i += 2
            if cmd == "m": x, y = x + a, y + b
            else: x, y = a, b
            sx, sy = x, y; pts = [(x, y)]
        elif cmd in "Hh":
            a = float(toks[i]); i += 1
            x = x + a if cmd == "h" else a; pts.append((x, y))
        elif cmd in "Vv":
            a = float(toks[i]); i += 1
            y = y + a if cmd == "v" else a; pts.append((x, y))
    return rects

class Buf:
    def __init__(s, w, h):
        s.w, s.h = w, h
        s.px = bytearray(w * h * 4)
    def fill(s, x0, y0, x1, y1, rgb, a=1.0, xor_mask=None):
        x0, y0, x1, y1 = max(0, round(x0)), max(0, round(y0)), min(s.w, round(x1)), min(s.h, round(y1))
        if x1 <= x0 or y1 <= y0: return
        r, g, b = rgb
        if a >= 1:
            row = bytes([r, g, b, 255]) * (x1 - x0)
            for y in range(y0, y1):
                o = (y * s.w + x0) * 4
                s.px[o:o + len(row)] = row
        else:
            for y in range(y0, y1):
                for x in range(x0, x1):
                    o = (y * s.w + x) * 4
                    da = s.px[o+3] / 255
                    oa = a + da * (1 - a)
                    for k, c in enumerate((r, g, b)):
                        s.px[o+k] = int((c * a + s.px[o+k] * da * (1 - a)) / oa) if oa else 0
                    s.px[o+3] = int(oa * 255)
    def bbox(s):
        xs0, ys0, xs1, ys1 = s.w, s.h, -1, -1
        for y in range(s.h):
            row = s.px[y*s.w*4:(y+1)*s.w*4]
            alph = row[3::4]
            if any(alph):
                ys0 = min(ys0, y); ys1 = max(ys1, y)
                first = next(i for i, v in enumerate(alph) if v)
                last = len(alph) - 1 - next(i for i, v in enumerate(reversed(alph)) if v)
                xs0 = min(xs0, first); xs1 = max(xs1, last)
        return None if ys1 < 0 else (xs0, ys0, xs1 + 1, ys1 + 1)
    def png(s, path, box):
        x0, y0, x1, y1 = box
        raw = bytearray()
        for y in range(y0, y1):
            raw.append(0); raw += s.px[(y*s.w + x0)*4:(y*s.w + x1)*4]
        def chunk(t, d): return struct.pack(">I", len(d)) + t + d + struct.pack(">I", zlib.crc32(t + d) & 0xffffffff)
        open(path, "wb").write(b"\x89PNG\r\n\x1a\n" + chunk(b"IHDR", struct.pack(">IIBBBBB", x1-x0, y1-y0, 8, 6, 0, 0, 0))
                               + chunk(b"IDAT", zlib.compress(bytes(raw), 9)) + chunk(b"IEND", b""))

def paint(buf, e, tf):
    tag = e.tag.replace(NS, "")
    sc, tx, ty = tf
    if tag == "g":
        t = parse_tf(e.get("transform"))
        for k in e: paint(buf, k, compose(tf, t))
        return
    fill = e.get("fill")
    if not fill or fill == "none": return
    a = float(e.get("fill-opacity", "1"))
    rgb = hexrgb(fill)
    if tag == "rect":
        x, y, w, h = (float(e.get(k, "0")) for k in ("x", "y", "width", "height"))
        buf.fill(tx + sc*x, ty + sc*y, tx + sc*(x+w), ty + sc*(y+h), rgb, a)
    elif tag == "path":
        rects = path_rects(e.get("d", ""))
        if e.get("fill-rule") == "evenodd" and len(rects) == 2:
            (ax0, ay0, ax1, ay1), (bx0, by0, bx1, by1) = rects
            # outline ring: four bars
            for (x0, y0, x1, y1) in ((ax0, ay0, ax1, by0), (ax0, by1, ax1, ay1), (ax0, by0, bx0, by1), (bx1, by0, ax1, by1)):
                buf.fill(tx + sc*x0, ty + sc*y0, tx + sc*x1, ty + sc*y1, rgb, a)
        else:
            for (x0, y0, x1, y1) in rects:
                buf.fill(tx + sc*x0, ty + sc*y0, tx + sc*x1, ty + sc*y1, rgb, a)

def run(svgpath, key, outroot):
    root = ET.parse(svgpath).getroot()
    W, H = int(float(root.get("width"))), int(float(root.get("height")))
    outdir = os.path.join(outroot, key); os.makedirs(outdir, exist_ok=True)
    layers = []
    pend = {"buf": None, "names": []}
    def flush_raster():
        if pend["buf"] is None: return
        box = pend["buf"].bbox()
        if box:
            n = len(layers)
            p = os.path.join(outdir, f"{n:02d}.png")
            pend["buf"].png(p, box)
            layers.append({"kind": "png", "file": p, "name": " + ".join(pend["names"]), "x": box[0], "y": box[1], "w": box[2]-box[0], "h": box[3]-box[1]})
        pend["buf"], pend["names"] = None, []
    def emit_raster(e, tf, name, alone=True):
        if alone: flush_raster()
        if pend["buf"] is None: pend["buf"] = Buf(W, H)
        paint(pend["buf"], e, tf)
        pend["names"].append(name)
        if alone: flush_raster()
    def emit_vector(e, tf, name):
        flush_raster()
        n = len(layers)
        p = os.path.join(outdir, f"{n:02d}.svg")
        sc, tx, ty = tf
        if e.get("data-size"):
            dw, dh = (float(v) for v in e.get("data-size").split())
            t = parse_tf(e.get("transform"))
            full = compose(tf, t)
            inner = "".join(ET.tostring(k, encoding="unicode") for k in e)
            inner = re.sub(r' xmlns="[^"]*"', "", inner)
            w, h = dw * full[0], dh * full[0]
            open(p, "w").write(f'<svg xmlns="http://www.w3.org/2000/svg" width="{w:.2f}" height="{h:.2f}" viewBox="0 0 {dw:g} {dh:g}">{inner}</svg>')
            layers.append({"kind": "svg", "file": p, "name": name, "x": round(full[1], 2), "y": round(full[2], 2), "size": os.path.getsize(p)})
            return
        body = ET.tostring(e, encoding="unicode")
        open(p, "w").write(f'<svg xmlns="http://www.w3.org/2000/svg" width="{W}" height="{H}" viewBox="0 0 {W} {H}"><g transform="translate({tx} {ty}) scale({sc})">{body}</g></svg>')
        layers.append({"kind": "svg", "file": p, "name": name, "x": 0, "y": 0, "size": os.path.getsize(p)})
    def walk(e, tf, depth):
        tag = e.tag.replace(NS, "")
        name = e.get("id", tag)
        if tag == "g" and (name in DESCEND or name.startswith("ui (")):
            t = parse_tf(e.get("transform"))
            sub = compose(tf, t)
            for k in e: walk(k, sub, depth + 1)
            return
        if is_raster(e, tf):
            emit_raster(e, tf, name)
        else:
            # mixed group: walk children, raster ones grouped individually
            if tag == "g":
                t = parse_tf(e.get("transform"))
                if t is not None and all(k.tag.replace(NS, "") in ("g", "rect", "path") for k in e) and not (e.get("transform") or "").startswith("translate"):
                    sub = compose(tf, t)
                    flush_raster()
                    part = 0
                    for k in e:
                        if is_raster(k, sub):
                            if pend["buf"] is None:
                                part += 1
                            emit_raster(k, sub, f"{name} · pixels {part}", alone=False)
                            pend["names"] = [f"{name} · pixels {part}"]
                        else:
                            walk(k, sub, depth + 1)
                    flush_raster()
                    return
            emit_vector(e, tf, name)
    for k in root:
        walk(k, (1.0, 0.0, 0.0), 0)
    flush_raster()
    json.dump(layers, open(os.path.join(outdir, "index.json"), "w"), indent=1, ensure_ascii=False)
    return layers

if __name__ == "__main__":
    d = sys.argv[1]
    for key in sys.argv[2:]:
        L = run(os.path.join(d, key + ".svg"), key, os.path.join(d, "layers"))
        print(key, len(L), [(l["kind"], l["name"][:30], l.get("size", "")) for l in L])
