# R002 import split: like raster.py, but one layer per named UI part (rows stay inside their
# panel's PNG), so a frame imports in 1-2 batches. Usage: python3 r002_split.py <dir> key...
import sys, os
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import raster
raster.DESCEND = {"ui", "canvas area"}
if __name__ == "__main__":
    d = sys.argv[1]
    for key in sys.argv[2:]:
        L = raster.run(os.path.join(d, key + ".svg"), key, os.path.join(d, "layers"))
        print(key, len(L), sum(1 for l in L if l["kind"] == "svg"))
