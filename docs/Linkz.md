# 🦎 Linkz (lizard inkz)

> Vector art you write in lizard-script. The program is built on lizard-script itself: Go is only a thin host underneath.

Private until the reveal.

## 🗂️ Knowledge base

| Note | What's in it |
|---|---|
| [[Linkz – Product Design]] | What Linkz is, the scope tiers, the design log |
| [[Linkz – Editor Design]] | The native editor: code = layers + properties, the v1 goal, the backdoor |
| [[Linkz – Brushes]] | Note: Procreate-style brush creator, `.lzbrush` import/export, importing Procreate/ABR brushes |
| [[Linkz – Tracing v2]] | Design: line-art tracing for anime art (strokes with width, smooth Bézier edges, hair spikes, shape snapping) |
| [[Linkz – Architecture]] | The thin host, the lizard libraries, where new code goes |
| [[Linkz – Language Needs]] | Every gap in lizard-script that Linkz hits, and the extensions they become |
| [[Linkz – Build Plan]] | Milestones, ⏳ waiting on you, progress log |
| [[Linkz – Style Guide]] | How the editor looks and what drawings look like by default |
| [[Linkz – Tools & Setup]] | Tools, commands, every-session checklist |

## ⚡ Linkz in 30 seconds

| | |
|---|---|
| What | A vector art program where a drawing is a lizard-script program |
| File | `.lzdraw`: plain-text lizard-script with a `-~> lzdraw 1` header |
| Output | SVG, then raster (PNG/JPEG/WebP) |
| Killer features | Image tracing (picture → editable shapes), an MCP so agents draw by writing lizard-script, 3D projection |
| Built with | lizard-script for everything; Go only for the host primitives |
| Ships as | One executable per OS on GitHub releases, plus an install script |

## 🧭 Non-negotiables

1. **Built on lizard.** Logic goes in lizard-script. Go only when lizard *can't* do it (talking to the OS, the runtime itself). Lizard being too slow is a reason to make lizard faster, not to move code to Go.
2. **`.lzdraw` is plain text.** Readable, diffable, writable by agents.
3. **Agents write lizard-script, not SVG.** The MCP runs code; it never accepts raw SVG.
4. **Every language change is logged** as an extension in lizard-script's `docs/DESIGN.md`. Extensions are additive: the Go foundation stays and existing programs keep working.
5. **One executable** to install.

## 🤔 Assumptions (took the recommended option; overturn any)

- Repo `linkz`, private on GitHub.
- The editor is a web page served by the `linkz` binary on localhost; a WebAssembly build for an online demo comes later.
- Until lizard-script has `use`, the host loads every library in `lib/` in a fixed order.

## ❓ Open questions

- When is the reveal, and is it tied to a course deadline?
- Do Maria José and Samuel join the repo after the reveal?
