# B032 – `lily3d`: 3D graphics for every lizard program

Repo: `../lizard-script` · Lane: lizard (+ one Go triangle primitive) · Parallel with B031

## Goal (v0, software renderer)
`use "lily3d"` renders 3D into a window, surface or image:

```
use "lily3d"
let s = lily3d.scene()
lily3d.add(s, lily3d.cube(1), lily3d.rotate_y(30deg), {"color": #f2a65a})
lily3d.render(s, lily3d.camera((0, 1.5, 4), (0, 0, 0), 60deg), w)
```

- Math: 3D vectors (POINTs with 3 coordinates), 4×4 matrices, `translate/rotate_x/y/z/scale/compose`, look-at camera, perspective and orthographic projection, normals.
- Meshes: `cube`, `sphere`, `plane`, `cylinder`, `cone`, `torus`; `lily3d.load_obj(path)` (vertices, faces, normals).
- Rendering: back-face culling, near-plane clipping, a depth buffer, flat and smooth (Gouraud) shading, ambient + directional + point lights, per-mesh color, wireframe mode, fog.
- **Speed**: all math in lizard; only the inner per-pixel step is Go — one Terrarium primitive that fills a triangle with depth test and interpolated color into a target + depth buffer (`terrarium/internal/pixels/triangles.go` + a canvas function, golden-tested). Note B043 is making the VM much faster at the same time.

## Where it lives
Standard module: `terrarium/std/lily3d.liz` + `terrarium/std/lily3d/*.liz` (the lookup mechanism is built by B031 at the same time; until it lands, test with relative `use` paths). Don't edit `evaluator/` or B031's files. Add only your own lines to `terrarium/terrarium.go` if needed.

## Clean code
One concept per file, documented, readable top to bottom. LEP 0001 style.

## Done when
Tests for the math (matrices, projection, clipping), golden renders on both engines, `examples/lily3d/cube.liz` (lit spinning cube in a window — report fps at 600×400 on the VM) and `teapot.liz` (a bundled .obj rendered to PNG), docs `docs/LILY3D.md`, guide chapter 20 "lily3d" (en + es; edit the guide index only for your row). Don't commit.
