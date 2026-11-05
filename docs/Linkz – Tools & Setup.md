# 🛠️ Linkz: Tools & Setup

## 🧰 Tools

| Tool | Version | For |
|---|---|---|
| lizardscript | the lizard-script commit in `../.lizard-stable` (or `../lizard-script`) | Runs, tests, builds and drives Linkz |
| Go | 1.27.0 | Only to build `lizardscript` itself (and the other systems' runtimes for a cross build) |
| gh | 2.91.0 | Private repo, releases |

## 📁 Layout on disk

`linkz/` sits next to `lizard-script/`: the CI builds `lizardscript` from it, and `scripts/crossbuild.sh` finds the runtimes' source there (or in `$LIZARD_SRC`).

## ⌨️ Commands

```bash
(cd ../lizard-script && go build -o ~/bin/lizardscript .)  # once, and after the language moves
lizardscript run --vm . edit [--no-splash] [file.lzdraw]   # also render, trace, paint, calibrate, mcp
lizardscript build --vm . -o linkz                         # one executable (bundle.txt lists the data files)
scripts/test.sh                                            # tests/ on both engines, HOME in a temp folder
scripts/parity.sh                                          # the outputs against the last Go host's goldens
scripts/smoke.sh                                           # the commands and a driven editor window
scripts/crossbuild.sh [dist]                               # the release archives for every system
lizardscript tests/editor/editor_bench.liz 5               # the editor benchmarks (M10 fixtures)
lizardscript fmt --check main.liz lib tests tools          # the layout
lizardscript check main.liz lib tests tools                # lib/lizard-check.liz lists the host's names
lizardscript tools/uses.liz --check lib                    # every library file uses what it needs
(cd ../lizard-script && go test ./...)                     # after any language change
```

`LINKZ_ENGINE=tree` runs the loaded code on the evaluator. `LIZARDSCRIPT=path` tells the scripts which `lizardscript` to use.

## ⚙️ Settings

`~/.linkz/settings.liz` (the folder `linkz.config_dir()` gives) is lizard-script: one `let` per setting. The editor evaluates it once at startup in its own small, step-limited environment and hands the variables to the UI as a `settings` hash. A missing file means defaults; a file with a syntax or runtime error is ignored with a warning on stderr and the defaults are used.

```
-~> Linkz settings
let language = "es"
let skip_loading_anim = true
```

| Setting | Values | Does |
|---|---|---|
| `language` | `"en"`, `"es"` | The interface language. Without it Linkz uses the system's language, else English. The editor's language command rewrites this line and keeps the others. |
| `skip_loading_anim` | `true` / `false` (default `false`) | `true` makes `linkz edit` start without the splash screen, like `--no-splash`. |
| `trackpad_pressure` | a hash `{"low": 0.8, "high": 1.0, "gamma": 1.5}` (these are the defaults) | The pressure curve of a Force Touch trackpad. A trackpad needs a firm press to click, so strokes start thick; the editor sees `clamp((raw - low) / (high - low), 0, 1) ** gamma` instead of the raw pressure, before the pens' own `pressure` blocks. `low` is the raw pressure that counts as none, `high` the one that counts as full, `gamma` above 1 keeps light pressure small. Keys left out keep their default; valid curves have `0 <= low < high <= 1` and `0.1 <= gamma <= 10`. A wrong curve is ignored with a warning on stderr and the defaults stay. The defaults are a starting point: tune them for your trackpad. |
| `tablet_pressure` | a hash like `trackpad_pressure` (default `{"low": 0, "high": 1, "gamma": 1}`, no change) | The same curve for a pen tablet. A mouse has none. |

To tune a pressure curve, look at the raw pressure of your strokes (the `raw` of each sample of a `move` event, see lizard-script's TERRARIUM.md → pen) while pressing as lightly and as hard as you draw; the curve takes it as input. (The Go `pressureprobe` tool went with the Go host in B042.) Set `low` just below your lightest stroke and `high` at your hardest.

`linkz edit --slow-splash` plays the splash at a quarter speed (to watch the animation); it plays even when `skip_loading_anim` is on.

## 📦 Releasing

Nothing is published until a person pushes a tag and publishes the draft. The pieces:

| File | Does |
|---|---|
| `scripts/crossbuild.sh` | `lizardscript build --vm --os OS --arch ARCH --icon assets/icon/linkz-1024.png --icon-small assets/icon/linkz-16.png . -o linkz` for darwin/arm64, darwin/amd64, linux/amd64 and windows/amd64 (the other systems' runtimes are built from lizard-script's source with Go), then `tar.gz` archives (`zip` on Windows) named `linkz_<version>_<os>_<arch>` with the README, the license and the examples, and `checksums.txt` (sha256). The icon is in each build: the darwin archives hold `Linkz.app` (not a bare `linkz`: `Contents/MacOS/Linkz`, the program's files in `Contents/Resources/program.lzb`, `icon.icns`, an `Info.plist`; signed ad hoc when `codesign` is there), the Windows `.exe` carries it as a resource, and the Linux archive has `icon.png`, `icon-256.png` and `linkz.desktop` next to `linkz`. The version is `linkz_version` in `lib/app/version.liz`. |
| `.github/workflows/ci.yml` | By hand only for now: checks out linkz and lizard-script side by side, builds `lizardscript`, runs `scripts/test.sh`, `scripts/parity.sh`, fmt and check, and the cross build. |
| `.github/workflows/release.yml` | On tags `v*` (and by hand): builds `lizardscript`, the tests, writes the tag's version into `lib/app/version.liz`, the cross build, then `gh release create --draft` with the archives and `checksums.txt`. |
| `scripts/install.sh`, `scripts/install.ps1` | Download the latest release's archive for this system, verify its checksum, install. On macOS `install.sh` puts `Linkz.app` in `~/Applications` (`/Applications` with `--system`, or `--apps DIR`) and makes a link `linkz` in the bin folder to the program inside it. `--from-file` / `-FromFile` installs a local archive without the network. |

**Flow:** `git tag v0.1.0 && git push origin v0.1.0` → release workflow → draft release with the four archives and `checksums.txt` → read it, then publish it on GitHub.

**Version:** `linkz_version` is the release version (`1.0.0`); between releases it may say `-dev`. A new version makes Linkz measure the machine again.

**Private repo:** while `alejandruxxug/linkz` is private, `install.sh` / `install.ps1` only work with `GITHUB_TOKEN` set to a token that can read it; once a release is public they work without one.

**App icon:** the approved icon is R006 E2 (`docs/Art/Reviews/R006/E2-*`, generator `docs/Art/Figma-gen/r006c.py`); `assets/icon/linkz-1024.png` and `linkz-16.png` are copies of its 1024 master and 16 px sprite (copy again after changing the generator). The editor and the splash set the master as the window's `icon` (`app_icon()` in `lib/app/editor.liz`); `bundle.txt` carries `assets/icon/*.png`. Opening `Linkz.app` with no arguments starts the editor (`main.liz`).

**Not done yet:** macOS notarization (`Linkz.app` is signed ad hoc, so `codesign --verify --strict` passes, but `spctl -a` rejects it: Gatekeeper wants a Developer ID signature and a notarization ticket; the installer prints the `xattr` workaround), Windows signing (the installers print the Gatekeeper and SmartScreen workarounds), the Homebrew tap, linux/arm64.

## 🔌 Every session

1. `git pull` in **both** repos (the language may have moved).
2. `scripts/test.sh` here (with `LIZARDSCRIPT` pointing at a fresh build).

## ⚠️ Things only the developer can do

- Approve commits in lizard-script (the team sees them).

## 🔧 Backdoor quick reference

`lizardscript drive [--pid N] <Command>` drives the newest running window (it was `linkz debug`); start the editor with `LIZARD_AUTOMATION=1`. Commands: `List`, `Dump`, `Screenshot [path]`, `Click x y`, `DoubleClick x y`, `Drag x1 y1 x2 y2 [steps] [--pressure a,b] [--full]`, `Scroll dy`, `Key name [mods]`, `Type text`, `Resize w h` (clamped to at least 640×420), `Focus true|false`, `Leave`, `Enter`, `Stats`, `Wait`, `Quit`, `Call Name args…` (editor/home commands such as `Open`, `Save`, `SetCode`, `SelectRow`, `SetProp`, `Tool`, `PencilStroke`, `Undo`, `Redo`, `Language`, `Home`, `OpenCard`, `NewDrawing`, `Complete`). Testers run with `HOME=$(mktemp -d)`, `LIZARD_AUTOMATION=1` and `--no-splash`.
