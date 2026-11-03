# Linkz

Vector art written in lizard-script. **Built entirely in lizard-script: the repo has no Go.** Linkz is `main.liz` + `lib/`, run by `lizardscript` and shipped with `lizardscript build`; everything that talks to the OS lives in the language (the Terrarium). Private until the reveal. Linkz = lizard inkz.

## Read first

| For | Read |
|---|---|
| What we're building, scope, design log | `docs/Linkz – Product Design.md` |
| Where code goes (lib vs language) | `docs/Linkz – Architecture.md` |
| What lizard-script is missing | `docs/Linkz – Language Needs.md` |
| Milestones, ⏳ waiting on you, progress | `docs/Linkz – Build Plan.md` |
| Commands, setup | `docs/Linkz – Tools & Setup.md` |
| The editor (rows, layers, two-way canvas, host↔lizard contract, backdoor) | `docs/Linkz – Editor Design.md` |
| The language itself | `../lizard-script/docs/GRAMMAR.md`, `SEMANTICS.md`, `DESIGN.md`, `TERRARIUM.md`, `STYLE.md` |

## Rules

1. **Everything in lizard.** New features go in `lib/*.liz`. If Linkz needs something from the OS or something lizard can't do fast enough, it goes into lizard-script (a Terrarium module or a language feature, general enough for any program), never a Go workaround here: log it in Language Needs and fix the language.
2. **Language changes go to `../lizard-script`**, following that repo's CLAUDE.md (both engines identical, tests, docs and the guide in English and Spanish, a row in `docs/DESIGN.md` §7). The team sees that repo: never mention Linkz there before the reveal; describe changes as language features.
3. **Clean code**: one concept per file, every public function documented with `-~>`, LEP 0001 style; `lizardscript fmt --check` and `lizardscript check` clean on `main.liz lib tests tools`.
4. **`.lzdraw` is plain lizard-script.** Agents and the MCP write lizard-script, never raw SVG.
5. Commit only when asked; never push lizard-script without the developer's yes.
6. Docs: decisions → Product Design's design log; finished milestones → Build Plan's progress log.

## Commands

```bash
scripts/test.sh                                     # every lizard test suite, on both engines
scripts/parity.sh                                   # outputs against the last Go host's goldens
scripts/smoke.sh                                    # an end-to-end run
lizardscript run --vm . edit                        # the editor (render, paint, trace, mcp, calibrate too)
lizardscript build --vm . -o linkz                  # one executable; scripts/crossbuild.sh for every OS
(cd ../lizard-script && go test ./...)              # after language changes
```
Each script takes `LIZARDSCRIPT=/path/to/lizardscript` (default: the one on PATH).

## Agent workflow

The main session orchestrates: it splits a milestone into small batches, writes a brief in `docs/Batches/B<nnn> – <name>.md`, starts a **fresh** `builder`, then a **fresh** `tester` per batch, and commits each green batch by path.

| Agent | Model | Does |
|---|---|---|
| `builder` | Sonnet | Implements one batch brief, runs the tests |
| `tester` | Sonnet | Renders examples, writes edge-case programs, drives the editor through `lizardscript drive` (`LIZARD_AUTOMATION=1`, isolated HOME and config dir); reports BUG / WEIRD / VISUAL |
| `lead` | Opus | WEIRD issues, BUGs whose first fix failed, and review of risky batches (language semantics, file formats, anything that writes the user's files) |
| `design` | Opus | VISUAL issues, editor UI, default look of drawings, library API ergonomics |

- **Fast lane:** a one-file tweak with no language or file-format change: the orchestrator edits, tests, commits. No agents.
- Two failed fix rounds on one issue → stop and ask the developer.
- The orchestrator owns shared files: `lib/app/host.liz` (the names the libraries get), `lib/lizard-check.liz` and `bundle.txt`.
