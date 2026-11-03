# B012 – Splash screen

Milestone: M3 · Repo: `linkz` · Lanes: art (`design`, Opus) then host + lizard (`builder`) · Needs: B005

## Goal
While the editor loads, a frameless, windowless banner floats on the desktop: a big pixel **Linkz** logo, a pixel lizard crawling across *behind* it, and **powered by lizard-script** underneath. Like the banner an app shows while it starts.

## Art (design agent) → `lib/ui/splash_art.liz`
- **Logo**: "Linkz" in chunky custom pixel letters (not the 7×13 UI font), Lizard 16 palette (`docs/Linkz – Style Guide.md`): `leaf` letters, `moss` shadow 1 px down-right, `ink` outline. About 120×32 px at 1×.
- **Lizard**: a side-view pixel lizard walk cycle, 6 frames, ~48×16 px, `leaf`/`moss`/`lime` with a `sun` eye; legs alternate, tail sways. It crawls left → right across the banner behind the logo (the logo covers it while it passes) and wraps around.
- **Tagline**: "powered by lizard-script" in the UI font, `stone`.
- Data format: each sprite/frame as `{"size": (w, h), "runs": [[y, x0, x1, color_name], ...]}` so lizard draws it with `fb.spans` per color. Palette by name from the Style Guide.
- Make it in Aseprite through the aseprite MCP if it's available (keep `docs/Art/splash.aseprite` as the source, export PNG frames to `docs/Art/Previews/`), else author the pixel data directly; either way convert to `splash_art.liz` with a small dev tool (`tools/spritegen`, Go, reading PNGs). Look at every frame (Read the PNG) before handing off. Log it in the Style Guide's art log.

## Behaviour (builder)
- `linkz edit` starts with a **frameless, transparent** window (Ebitengine `SetWindowDecorated(false)`, `SetScreenTransparent(true)` before `Run`), centered, 3× UI scale of a ~200×80 px banner, floating above other windows. Pixels lizard leaves transparent show the desktop.
- **Order** (the interpreter runs one program at a time): (1) load only `lib/gfx/font.liz` + `lib/ui/splash.liz`; lizard renders all splash frames into offscreen buffers (`fb.offscreen`); (2) the host starts loading the editor libraries in the background (one goroutine, nothing else evaluating); (3) meanwhile the host cycles the prepared frames at ~12 fps (Go only flips buffers lizard drew); (4) when loading is done **and** at least 1.2 s have passed, the window becomes decorated, opaque content, resizes to the editor size and recenters, and the editor paints.
- A click on the splash does nothing; Escape quits. If loading fails, the splash shows the error line in `ember` under the tagline and stays until clicked.
- `linkz edit --no-splash` skips it (for agents and the smoke script); the backdoor works from the moment the editor is up.
- Verify that switching decorated/transparent → normal at runtime works on macOS with Ebitengine 2.10.5; if it doesn't, fall back to closing the splash and opening the editor window cleanly, and say which happened.

## Test
`linkz edit` by the developer (it's visual), plus: a Go test renders the splash frames headlessly via `linkz paint` and checks the logo pixels; screenshots of frames 0 and 3 Read by the tester.

## Done when
The splash floats frameless over the desktop with the lizard crawling behind the logo, the editor appears after loading, `--no-splash` works, `go test ./...` green.
