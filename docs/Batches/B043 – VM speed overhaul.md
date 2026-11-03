# B043 – VM speed overhaul: catch Python, chase Lua (language)

Repo: `../lizard-script` · Needs: B039 committed (it owns compiler/vm until then) · Measuring stick: `python3 bench/compare.py`

> "bench it with lua … with all your musclepower of course" — the developer, 2026-11-03

## Where we are (M4, best of 5, startup included; lizard startup ~26 ms, Python ~16 ms)

| Benchmark | lizard VM | Python 3.14 | Lua 5.5 |
|---|---|---|---|
| closures (1M calls) | 0.104 s | 0.045 | 0.013 |
| fib(27) | 0.062 | 0.027 | 0.007 |
| float (5M rounds) | 0.296 | 0.112 | 0.024 |
| hashes (200k records) | 0.110 | 0.050 | 0.024 |
| loop (10M rounds) | 0.605 | 0.228 | 0.053 |
| points (2M adds) | 0.169 | 0.077 | 0.099 |
| strings (200k chars) | 0.719 | 0.023 | 0.211 |

Goal: **beat CPython on every benchmark**, and get within 2–3× of Lua on loops, floats and calls. Behaviour stays identical on both engines (differential tests are the safety net).

## Build, in the order the profiler says pays most (measure each step with bench/compare.py and `go test ./vm -bench`; keep what wins, drop what doesn't)
1. **Unboxed values in the VM.** Replace `object.Object` on the VM stack and in locals with a small tagged value (e.g. `struct{ kind; num uint64/float64 bits; ref object.Object }`), boxing only when a value leaves the VM (builtins, host env, hashes/arrays until they also hold values). The design note is DESIGN.md §8. This is the big one: integers and floats stop allocating, the GC (~45–50% of CPU in real hosts) mostly disappears.
2. **Loop scopes only when needed.** Each loop round gets a fresh scope only if a closure created in the body captures a loop variable (the compiler knows); otherwise reuse the slots. Same for blocks.
3. **Fused loop and arithmetic instructions**: integer `for` with a compare-and-jump fused (`i < n` + jump), `i += 1` as one instruction on a local, fast paths for `%`, `//`, `**`, unary minus, comparisons of mixed numbers; superinstructions for the hottest pairs the profiler shows.
4. **Calls**: no allocation per call (frames and locals in one preallocated stack), arguments passed in place, a direct path for calling a known closure, cheap returns.
5. **Strings**: make `s += x` amortized O(1) when `s` isn't shared (a builder/rope behind the string value, or Python's in-place trick when the VM knows it holds the only reference); `len(s)` and `s[i]` O(1) on ASCII and cached for UTF-8 (already partly there).
6. **Hashes and points**: values stored unboxed inside points and small arrays; the hash field cache from B038 reused for writes.
7. **Startup**: measure the 26 ms (Terrarium install, json.liz compile); lazily install modules on first use, cache compiled embedded libraries.
8. If after 1–7 the dispatch loop itself dominates: consider a register-based instruction set for hot functions (Lua's design), as a design note with a prototype benchmark before committing to it.

## Done when
`bench/compare.py` table before/after in the report and in bench/README.md; every benchmark ≥ CPython; the `go test ./vm -bench` set before/after; differential + both-engine tests green; tree-walker unaffected or faster; docs (DESIGN.md §2/§8, guide chapter 7's profiling part if numbers changed). Clean code: one concept per file, the value representation explained in plain words at the top of its file. Don't commit. Never mention Linkz in lizard-script.
