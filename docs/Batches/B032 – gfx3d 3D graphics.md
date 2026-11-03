# B032 – `gfx3d`: 3D graphics in the Terrarium

Milestone: std · Repo: `../lizard-script` · Lane: lizard (+ one Go triangle/depth span primitive) · Needs: B031

## Goal (v0, software renderer)
- Math: vectors (POINTs with 3 coordinates), 4×4 matrices, transforms (translate, rotate, scale), look-at camera, perspective and orthographic projection.
- Meshes: `cube`, `sphere`, `plane`, `cylinder`; `gfx3d.load_obj(path)`.
- Rendering into a `window` or `img`: back-face culling, near-plane clipping, a depth buffer, flat and smooth (Gouraud) shading, ambient + directional light, per-mesh color, wireframe mode.
- Scene API: `let s = gfx3d.scene()`, `gfx3d.add(s, mesh, transform, material)`, `gfx3d.render(s, camera, target)`.
- Speed: all math in lizard; only the inner "fill a span with depth test and interpolated color" is a Go primitive (the per-pixel exception, like gradient spans).

## Later
GPU path (Ebitengine shaders, then Metal / Direct3D through the FFI), textures, shadows.

## Clean code
Same rules as B031: one concept per file, documented, readable top to bottom.

## Done when
`examples/terrarium/cube.liz` (a lit spinning cube in a window, ≥ 30 fps at 600×400 on the VM) and `teapot.liz` (load an .obj, render a PNG); tests for the math (matrices, projection, clipping) and a golden-image render; docs in TERRARIUM.md and docs/learn.
