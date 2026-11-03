# Contributing to Linkz 🦎

Linkz is open source under the [MIT license](LICENSE): fork it, tweak it, make your own art app out of it. This page is about sending changes back.

## The idea to keep

**The code is the drawing.** A `.lzdraw` file is lizard-script, the code panel is the layers and properties panel, and editing on the canvas rewrites the code. And Linkz is **built on lizard**: everything you see is lizard code in `lib/`; Go in `host/` only opens the window, reads input and pressure, fills pixels and talks to the OS. When something is missing or slow, the fix belongs in the language ([lizard-script](https://github.com/alejandruxxug/lizard-script)), not in a Go workaround.

## Kinds of changes

| You want to… | Do this |
|---|---|
| Fix a bug | An issue with steps, what happened and what you expected (the **Bug** template), or a pull request with the fix and a test. |
| Change how the editor looks | Designs come first: the UI is drawn in Figma and approved before it's built (see [docs/Linkz – Style Guide.md](docs/Linkz%20–%20Style%20Guide.md)). Open an issue with a sketch or mockup. |
| Add a feature | Open an issue (the **Feature** template) so we agree on the shape first. |
| Need something from the language | Propose it to lizard-script (its `proposals/`, the LEP process). |

## Building and testing

```sh
go build -o linkz .       # the app
./linkz edit              # open the editor (--no-splash to skip the splash)
go test ./...             # every package
go vet ./...
scripts/crossbuild.sh     # macOS, Windows and Linux builds
scripts/smoke.sh          # renders a drawing end to end
```

`LINKZ_ENGINE=tree` runs everything on lizard-script's tree-walker instead of the bytecode VM.

**Testing the real app:** `linkz debug` drives a running window (click, drag with pressure, type, screenshot, read its state). Use a throwaway home so your own drawings and settings are never touched:

```sh
HOME=$(mktemp -d) ./linkz edit --no-splash &
./linkz debug Screenshot shot.png
```

Commands are listed in [docs/Linkz – Tools & Setup.md](docs/Linkz%20–%20Tools%20&%20Setup.md).

## How the code is organized

Start at [docs/Linkz.md](docs/Linkz.md), then [Architecture](docs/Linkz%20–%20Architecture.md) and [Editor Design](docs/Linkz%20–%20Editor%20Design.md). In short: `lib/draw` (the drawing language), `lib/gfx` (the rasterizer), `lib/ui` (the editor, home screen, splash), `lib/mcp` (the agent server), `lib/trace` (image tracing), and `host/` (the thin Go host).

## Code style

Clean and readable above all: one concept per file, small functions, every public function documented with a `-~>` comment (what it does, an example), no dead code, no one-off hacks. Every visible text goes through `t()` with English and Spanish strings in `lib/ui/lang_en.liz` and `lang_es.liz`.

## Pull requests

One topic per pull request, tests included, `go test ./...` green, and a line in the docs when users can see the change.
