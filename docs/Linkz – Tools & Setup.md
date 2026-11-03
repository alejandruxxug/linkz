# 🛠️ Linkz: Tools & Setup

## 🧰 Tools

| Tool | Version | For |
|---|---|---|
| Go | 1.27.0 | Host + lizard-script |
| gh | 2.91.0 | Private repo, later releases |
| Chrome + Claude in Chrome | | The tester drives the editor (from M3) |

## 📁 Layout on disk

`linkz/` must sit next to `lizard-script/`, because `go.mod` points at `../lizard-script`.

## ⌨️ Commands

```bash
go test ./...                                  # host tests
go build -o linkz . && ./linkz render examples/hello.lzdraw
(cd ../lizard-script && go test ./...)         # after any language change
```

## 🔌 Every session

1. `git pull` in **both** repos (the language may have moved).
2. `go test ./...` here.

## ⚠️ Things only the developer can do

- Approve commits in lizard-script (the team sees them).
