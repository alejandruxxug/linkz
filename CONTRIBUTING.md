# Contributing to Linkz 🦎

Linkz is open source under the [MIT license](LICENSE): fork it, tweak it, make your own art app out of it. This page is about sending changes back.

## The idea to keep

**The code is the drawing.** A `.lzdraw` file is lizard-script, the code panel is the layers and properties panel, and editing on the canvas rewrites the code. And Linkz is **built on lizard**: it is a lizard-script program (`main.liz` and `lib/`), and the window, pixels, files and pen pressure come from lizard-script's Terrarium. When something is missing or slow, the fix belongs in the language ([lizard-script](https://github.com/alejandruxxug/lizard-script)), not in a Go workaround.

## Kinds of changes

| You want to… | Do this |
|---|---|
| Fix a bug | An issue with steps, what happened and what you expected (the **Bug** template), or a pull request with the fix and a test. |
| Change how the editor looks | Designs come first: the UI is drawn in Figma and approved before it's built (see [docs/Linkz – Style Guide.md](docs/Linkz%20–%20Style%20Guide.md)). Open an issue with a sketch or mockup. |
| Add a feature | Open an issue (the **Feature** template) so we agree on the shape first. |
| Need something from the language | Propose it to lizard-script (its `proposals/`, the LEP process). |

## Building and testing

```sh
lizardscript run --vm . edit   # open the editor (--no-splash to skip the splash)
scripts/test.sh                # every test, on the evaluator and on the VM
scripts/parity.sh              # the outputs against the last Go host's goldens
scripts/smoke.sh               # the commands and a driven editor, end to end
scripts/crossbuild.sh          # macOS, Windows and Linux builds
lizardscript fmt --check main.liz lib tests tools
lizardscript check main.liz lib tests tools
```

`LINKZ_ENGINE=tree` runs everything on lizard-script's tree-walker instead of the bytecode VM.

**Testing the real app:** `lizardscript drive` drives a running window (click, drag with pressure, type, screenshot, read its state). Use a throwaway home so your own drawings and settings are never touched:

```sh
HOME=$(mktemp -d) LIZARD_AUTOMATION=1 lizardscript run --vm . edit --no-splash &
lizardscript drive Screenshot shot.png
```

Commands are listed in [docs/Linkz – Tools & Setup.md](docs/Linkz%20–%20Tools%20&%20Setup.md).

## How the code is organized

Start at [docs/Linkz.md](docs/Linkz.md), then [Architecture](docs/Linkz%20–%20Architecture.md) and [Editor Design](docs/Linkz%20–%20Editor%20Design.md). In short: `lib/draw` (the drawing language), `lib/gfx` (the rasterizer), `lib/ui` (the editor, home screen, splash), `lib/mcp` (the agent server), `lib/trace` (image tracing), and `lib/app` (the commands, the window loop and how a drawing runs).

## Code style

Clean and readable above all: one concept per file, small functions, every public function documented with a `-~>` comment (what it does, an example), no dead code, no one-off hacks. Every visible text goes through `t()` with English and Spanish strings in `lib/ui/lang_en.liz` and `lang_es.liz`.

## Pull requests

One topic per pull request, tests included, `scripts/test.sh` green, and a line in the docs when users can see the change.
