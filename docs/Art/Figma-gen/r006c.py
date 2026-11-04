"""R006 round 3: the developer picked E and asked to remove the small far (left) eye. Writes E2-* next to E-*.
Run: python3 r006c.py (needs Pillow). The 16 px sprite for E never had that eye, so E2-16 equals E-16."""
import os
from PIL import Image
import r006 as R
import r006b as B

def main():
    out = R.OUT; name = "E2-cool-brush-in-teeth-one-eye"
    m = R.master(B.option_e(far_eye=False)); m.save(f"{out}/{name}-1024.png")
    s16 = B.small("e"); s16.save(f"{out}/{name}-16.png")
    for s in (32, 64, 128): m.resize((s, s), Image.LANCZOS).save(f"{out}/{name}-{s}.png")
    l, d = B.dock(m, s16, (226, 228, 232, 255)), B.dock(m, s16, (22, 24, 28, 255))
    l.save(f"{out}/{name}-dock-light.png"); d.save(f"{out}/{name}-dock-dark.png")
    S = Image.new("RGBA", (128 + 64 + 32 + 16 + 100, 168), (0, 0, 0, 0)); x = 20
    for s, ic in ((128, Image.open(f"{out}/{name}-128.png")), (64, Image.open(f"{out}/{name}-64.png")), (32, Image.open(f"{out}/{name}-32.png")), (16, s16)):
        S.alpha_composite(ic, (x, 20)); x += s + 20
    S.save(f"{out}/{name}-small-sizes.png")
    e = Image.open(f"{out}/E-cool-brush-in-teeth-1024.png")
    T = Image.new("RGBA", (2 * 1024 + 60, 1024 + 40), (11, 12, 11, 255)); T.alpha_composite(e, (20, 20)); T.alpha_composite(m, (1044, 20))
    T.resize((T.width // 2, T.height // 2), Image.LANCZOS).save(f"{out}/overview-3.png")
main()
