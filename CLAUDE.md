# Linkz

Vector art written in lizard-script. **Built on lizard-script itself: Go is only a thin host.** Private until the reveal. Linkz = lizard inkz.

## Read first

| For | Read |
|---|---|
| What we're building, scope, design log | `docs/Linkz – Product Design.md` |
| Where code goes (host vs lib vs language) | `docs/Linkz – Architecture.md` |
| What lizard-script is missing | `docs/Linkz – Language Needs.md` |
| Milestones, ⏳ waiting on you, progress | `docs/Linkz – Build Plan.md` |
| Commands, setup | `docs/Linkz – Tools & Setup.md` |
| The editor (rows, layers, two-way canvas, host↔lizard contract, backdoor) | `docs/Linkz – Editor Design.md` |
| The language itself | `../lizard-script/docs/GRAMMAR.md`, `SEMANTICS.md`, `DESIGN.md` |

## Rules

1. **Logic in lizard, not Go.** New features go in `lib/*.liz`. Go in `host/` only for a primitive lizard can't write (OS access). Never move logic to Go because lizard is slow or lacks a feature: log it in Language Needs and fix the language.
2. **Language changes go to `../lizard-script`**, in that repo's style (phase packages, table-driven tests, `tests/valid` + `tests/errors` programs). Changes are **additive**: the Go foundation stays and existing lizard programs keep working. Every change gets a row in its `docs/DESIGN.md` → *Extensions beyond the course design* (what, why, what it adds). The team sees that repo: never mention Linkz there before the reveal; describe changes as language features.
3. **`.lzdraw` is plain lizard-script.** Agents and the MCP write lizard-script, never raw SVG.
4. Commit only when asked; never push lizard-script without the developer's yes.
5. Docs: decisions → Product Design's design log; finished milestones → Build Plan's progress log.

## Commands

```bash
go test ./...                                       # here
(cd ../lizard-script && go test ./...)              # after language changes
go build -o linkz . && ./linkz render examples/hello.lzdraw
```

## Agent workflow

The main session orchestrates: it splits a milestone into small batches, writes a brief in `docs/Batches/B<nnn> – <name>.md`, starts a **fresh** `builder`, then a **fresh** `tester` per batch, and commits each green batch by path.

| Agent | Model | Does |
|---|---|---|
| `builder` | Sonnet | Implements one batch brief, runs the tests |
| `tester` | Sonnet | Renders examples, writes edge-case programs, drives the native editor through `linkz debug`; reports BUG / WEIRD / VISUAL |
| `lead` | Opus | WEIRD issues, BUGs whose first fix failed, and review of risky batches (language semantics, file formats, anything that writes the user's files) |
| `design` | Opus | VISUAL issues, editor UI, default look of drawings, library API ergonomics |

- **Fast lane:** a one-file tweak with no language or file-format change: the orchestrator edits, tests, commits. No agents.
- Two failed fix rounds on one issue → stop and ask the developer.
- The orchestrator owns shared files: `go.mod`, `lib/lib.go` (library order), `host/host.go`'s primitive list.
