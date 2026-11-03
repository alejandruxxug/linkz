# 🪟 Linkz: Windows test checklist

Linkz builds for Windows (`scripts/crossbuild.sh`), but nobody has run it there yet. Someone with a Windows PC (ideally with a pen) should go through this and write results in the Build Plan's progress log.

## App basics (B023)
1. `linkz edit` starts (splash floats, then the editor); `linkz debug Dump` finds the window; killing it leaves a stale `~/.linkz/debug-<pid>.json` that the next start cleans up.
2. Copy/paste in the editor round-trips text, including á ñ, and to/from Notepad (~0.3 s PowerShell latency per paste is expected).
3. Home → delete a drawing → it's in the Recycle Bin and restorable; no console window flashes.
4. Home → ⋯ → Show in folder → Explorer opens with the file selected.
5. With Windows set to Spanish, the UI starts in Spanish.
6. `linkz edit C:\path\x.lzdraw` shows `x.lzdraw` in the tab (known gap: needs `fs.slash` in `view_layout.liz`).
7. Home shows `Documents\Linkz`, thumbnails appear in `%USERPROFILE%\.linkz\thumbs`, names with spaces open.

## Pen pressure via Windows Ink (B024)
1. `go run ./tools/pressureprobe` prints `monitor installed: true` (else: the window class or user32 procs weren't found).
2. Light → hard pressure moves 0 → 1 smoothly, back to 0 on lift.
3. Tilt both axes → `tilt` near ±1 at the extremes (report if a sign looks flipped).
4. Barrel rotation (if the pen reports it) changes `rot`.
5. Eraser end → `eraser=true`; normal tip → false.
6. Hover in/out of range → nothing stuck; a mouse click afterwards reads pressure 1, source mouse.
7. Fast strokes in the editor aren't polygonal (coalesced samples kept).
8. Ink lands under the tip at 100 / 150 / 200 % display scaling and across two monitors with different scales.
9. Mouse and touch work exactly as before (clicks, drags, wheel, keys), also when alternating with the pen.
