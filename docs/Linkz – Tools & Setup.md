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

## 🔌 Every session

1. `git pull` in **both** repos (the language may have moved).
2. `go test ./...` here.

## ⚠️ Things only the developer can do

- Approve commits in lizard-script (the team sees them).
