# B064 – The app icon (E2) in every build

Model: Sonnet. Approved: R006 **E2** (`docs/Art/Reviews/R006/E2-cool-brush-in-teeth-one-eye-*`, generator `docs/Art/Figma-gen/r006c.py`), "Icon E2 is the best" — the developer, 2026-11-04.

**Step 1 – lizard-script** (its CLAUDE.md rules; read only `tools/build/` and `scripts/` in lizard-script, plus the window options in `terrarium/window/window.go`): `lizardscript build` takes `--icon <png>` (a 1024 master; optional `--icon-small <png>` for the hand-drawn 16 px), and:
- darwin: writes an `.app` bundle (`Name.app/Contents/{MacOS,Resources,Info.plist}`) with an `.icns` made in Go (PNG-based icns entries: ic07–ic10, ic11–ic14; no `iconutil` needed so it builds on any OS), and ad-hoc signs it when `codesign` is there (`codesign -s - --deep`). That also addresses the known "strict codesign fails" note: report what `codesign --verify --strict` and `spctl -a -vv` say now;
- windows: embeds an `.ico` (16, 32, 48, 256) as the exe's icon resource (a `.syso` written in Go before the build, or the PE resource added after; choose and explain);
- linux: writes `icon.png` (256 and 512) next to the binary, plus a `.desktop` template.
- a window option `"icon": path` (`ebiten.SetWindowIcon` on Windows/Linux; on macOS the Dock icon when run unbundled, via NSApp `setApplicationIconImage` with purego as `pen_darwin.go` does).
Tests for the icns/ico writers (parse back what was written). `go test ./...`, `go vet`, cross-builds. Never mention Linkz. Don't commit.

**Step 2 – Linkz** (rules `_v11-rules.md`; read only `scripts/crossbuild.sh`, the install scripts it names, `main.liz`'s window opening): export the E2 master and 16 px sprite to `assets/icon/` (1024 + 16; keep the generator as the source), pass `--icon` in `crossbuild.sh` (darwin archives now hold `Linkz.app`; keep the install scripts working with it), set the window `icon` option (splash and editor), and update `docs/Linkz – Tools & Setup.md`. Run `scripts/crossbuild.sh` (LIZARD_SRC = your step-1 tree), unpack the darwin archive and check the bundle opens (`open Linkz.app` with an isolated HOME, then quit it) and shows the icon (`sips`/`mdls` or a Finder screenshot is fine). Don't commit.
