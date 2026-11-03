# 🛠️ Linkz: Tools & Setup

## 🧰 Tools

| Tool | Version | For |
|---|---|---|
| Go | 1.27.0 | Host + lizard-script |
| gh | 2.91.0 | Private repo, later releases |
| Ebitengine | 2.10.5 | Opens the window only |
| purego | 0.11.0 | Trackpad/tablet pressure (`go run ./tools/pressureprobe` to check) |

## 📁 Layout on disk

`linkz/` must sit next to `lizard-script/`, because `go.mod` points at `../lizard-script`.

## ⌨️ Commands

```bash
go test ./...                                  # host tests
go build -o linkz . && ./linkz render examples/hello.lzdraw
(cd ../lizard-script && go test ./...)         # after any language change
```

## ⚙️ Settings

`~/.linkz/settings.liz` (the folder `os.config_dir()` gives) is lizard-script: one `let` per setting. The editor evaluates it once at startup in its own small, step-limited environment and hands the variables to the UI as a `settings` hash. A missing file means defaults; a file with a syntax or runtime error is ignored with a warning on stderr and the defaults are used.

```
-~> Linkz settings
let language = "es"
let skip_loading_anim = true
```

| Setting | Values | Does |
|---|---|---|
| `language` | `"en"`, `"es"` | The interface language. Without it Linkz uses the system's language, else English. The editor's language command rewrites this line and keeps the others. |
| `skip_loading_anim` | `true` / `false` (default `false`) | `true` makes `linkz edit` start without the splash screen, like `--no-splash`. |

`linkz edit --slow-splash` plays the splash at a quarter speed (to watch the animation); it plays even when `skip_loading_anim` is on.

## 📦 Releasing

Nothing is published until a person pushes a tag and publishes the draft. The pieces:

| File | Does |
|---|---|
| `.goreleaser.yaml` | Builds `linkz` for darwin/arm64, darwin/amd64, windows/amd64 and linux/amd64 (`CGO_ENABLED=0`, `-s -w`, version injected with `-X main.version={{.Version}}`), makes `tar.gz` archives (`zip` on Windows) with the README and examples, `checksums.txt` (sha256) and a changelog from the commits. The release is created as a **draft**. The Homebrew tap block is commented out, with the steps to turn it on. |
| `.github/workflows/ci.yml` | On every push and PR: checks out linkz and lizard-script side by side, `go test ./...`, `scripts/crossbuild.sh`. |
| `.github/workflows/release.yml` | On tags `v*` (and by hand): the same checks, then `goreleaser release --clean`. |
| `scripts/install.sh`, `scripts/install.ps1` | Download the latest release's archive for this system, verify its checksum, install. `--from-file` / `-FromFile` installs a local archive without the network. |

**Flow:** `git tag v0.1.0 && git push origin v0.1.0` → release workflow → draft release with the four archives and `checksums.txt` → read it, then publish it on GitHub.

**Version:** `main.version` is `0.1.0-dev` in source and set by the release build. A new version makes the host measure the machine again.

**The lizard-script problem:** `go.mod` has `replace lizardscript => ../lizard-script`, so a plain `go install` or a release build can't resolve it. The workflows check out `alejandruxxug/lizard-script` into `lizard-script/` beside `linkz/` (`path:` on both checkouts), which makes the `replace` work. If lizard-script's `go.mod` module is ever renamed to `github.com/alejandruxxug/lizard-script` (and the imports with it), Linkz can require it normally, the `replace` and the second checkout can go, and `go install github.com/alejandruxxug/linkz@latest` works directly. That rename touches every lizard-script import, so it needs the team's agreement.

**Private repo:** while `alejandruxxug/linkz` is private, `install.sh` / `install.ps1` only work with `GITHUB_TOKEN` set to a token that can read it; once a release is public they work without one.

**Local dry run** (publishes nothing): `goreleaser release --snapshot --clean --skip=publish` from a linkz folder that has `../lizard-script`. Needs goreleaser v2 installed (`brew install goreleaser`).

**Not done yet:** macOS notarization, Windows signing (the installers print the Gatekeeper and SmartScreen workarounds), the Homebrew tap, linux/arm64.

## 🔌 Every session

1. `git pull` in **both** repos (the language may have moved).
2. `go test ./...` here.

## ⚠️ Things only the developer can do

- Approve commits in lizard-script (the team sees them).

## 🔧 Backdoor quick reference

`linkz debug [--pid N] <Command>` drives the newest running window. Commands: `List`, `Dump`, `Screenshot [path]`, `Click x y`, `DoubleClick x y`, `Drag x1 y1 x2 y2 [steps] [--pressure a,b] [--full]`, `Scroll dy`, `Key name [mods]`, `Type text`, `Resize w h` (clamped to at least 640×420), `Focus true|false`, `Leave`, `Enter`, `Stats`, `Wait`, `Quit`, `Call Name args…` (editor/home commands such as `Open`, `Save`, `SetCode`, `SelectRow`, `SetProp`, `Tool`, `PencilStroke`, `Undo`, `Redo`, `Language`, `Home`, `OpenCard`, `NewDrawing`, `Complete`). Testers run with `HOME=$(mktemp -d)` and `--no-splash`.
