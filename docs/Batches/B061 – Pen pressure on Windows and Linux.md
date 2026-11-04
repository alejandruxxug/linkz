# B061 – Pen pressure on Windows and Linux (lizard-script)

Model: Opus. Repo: `../lizard-script`, following ITS `CLAUDE.md` (a §7 row in DESIGN.md, TERRARIUM.md, the guide en+es if users can see the change). Read only: `terrarium/pen/` (all files) and the `window` code that delivers drag samples to `pen`.

Today, macOS reads real pen pressure (`pen_darwin.go`), while Windows and Linux use stubs (mouse, pressure 1.0). Make both real, with no cgo beyond what Ebitengine already needs, and keep the `pen` API unchanged.

- **Windows:** Windows Ink pointer messages (`WM_POINTER*`, `GetPointerPenInfo`: pressure, tilt, eraser, and pen vs touch vs mouse) via `golang.org/x/sys/windows`. Hook the window procedure of the Ebitengine window (subclass it), then feed the same per-drag samples the darwin code feeds.
- **Linux:** research and choose. Options are XInput2 valuators on X11 (through the GLFW/X11 handle Ebitengine already links), and/or Wayland `tablet-unstable-v2`. Pick what works on a normal X11/XWayland desktop, and degrade to the stub if it's unavailable. Write the choice and its reason into DESIGN.md.
- Device kind (`"tablet"`, `"trackpad"`, `"mouse"`) reported as on macOS, so per-input overrides work.
- Nothing can be tested on real hardware here. Isolate the OS calls behind a small interface and unit-test the conversion (raw values → samples, ranges, tilt, eraser) with fake events. The build must pass `GOOS=windows`, `GOOS=linux` and darwin; `go vet` must be clean for all three.
- Write a manual check list (in `docs/`) for someone with a Windows tablet and with a Linux tablet.
- Never mention Linkz. Don't commit.

Report: the approach per OS, the files, the test and build output for all three OSes, and the hardware checks still needed.
