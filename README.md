# Linkz 🦎

*lizard inkz*

Vector art you write in [lizard-script](https://github.com/alejandruxxug/lizard-script). Linkz itself is built on lizard-script: the shapes, the SVG writer and everything after them are lizard code, on top of a thin Go host.

```
-~> lzdraw 1
canvas = {"width": 400, "height": 200, "background": "#f4f1e8"}
line(20, 100, 380, 100)
for (let i = 0; i < 6; i = i + 1) {
  circle(40 + i * 64, 100, 6 + i * 4)
}
```

## Install

From a release (once one is public; while the repo is private, set `GITHUB_TOKEN` to a token that can read it):

```bash
# macOS / Linux: installs to ~/.local/bin (add --system for /usr/local/bin)
scripts/install.sh
```

```powershell
# Windows
.\scripts\install.ps1
```

Both download the archive for your system from the latest release, check its sha256, and install `linkz`. To install an archive you already downloaded: `scripts/install.sh --from-file linkz_<version>_<os>_<arch>.tar.gz --checksums checksums.txt`. Builds: macOS (Apple Silicon and Intel), Windows (amd64), Linux (amd64).

## Try it

Linkz needs `lizard-script` checked out next to it.

```bash
go build -o linkz .
./linkz render examples/hello.lzdraw   # writes examples/hello.svg
```

The plan, the architecture and the design log are in [docs/](docs/Linkz.md).
