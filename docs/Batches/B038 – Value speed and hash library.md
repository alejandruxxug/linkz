# B038 – Faster values and a hash library (M10, language side)

Repo: `../lizard-script` · Lane: Go (object, evaluator, compiler, vm) · Needs: B033

## Why
- B033: dot access (1.23×) and float math (1.46×) are bound by how values are stored: `Hash.Pairs` is a Go map keyed by a struct holding a string, and every float result is a heap-boxed `*object.Float`. Go's GC is ~45% of CPU in a real host workload.
- Apps ask for: comparing two hashes/arrays **by identity** (caches today invent ids) and a fast **copy** of a hash or array (today a `keys()` loop in lizard).

## Build
1. **`hash` library** (like `text`/`list`): `hash.copy(h)`, `hash.keys(h)`, `hash.values(h)`, `hash.merge(a, b)` (b wins), `hash.has(h, key)`, `hash.remove(h, key)` (returns whether it was there); `list.copy(a)`. Both engines, docs (SEMANTICS, DESIGN §7, guide chapter 9 en + es, cheatsheet).
2. **Identity**: `same(a, b)` builtin — true when both are the same array/hash/function/native value (for numbers, strings, booleans, colors, points and null it is `===`). Tests that it never confuses equal-but-different containers.
3. **Faster hash reads**: give `object.Hash` a cheaper key path for string keys (e.g. a `map[string]` fast lane for string keys alongside the general one, or interned key strings) and a per-instruction cache in `OpGetField` keyed by the hash's identity + a mutation version. Keep iteration order and every behaviour identical (differential tests).
4. **Fewer boxed floats** where results don't escape: e.g. fuse arithmetic chains in the VM (`a * b + c` computed in Go floats before boxing once), compare without boxing, and a small-float cache where it helps. Measure; don't change semantics.
5. Design note (no code): what fully unboxed numbers would take (a tagged `Value` struct in the VM instead of `object.Object`), costs and benefits, as a section in DESIGN.md "Future: unboxed values".

## Done when
Before/after benchmark table (the B033 set, plus a hash-heavy one: build/copy/read 10k small hashes); DotAccess ≥ 1.5× over c703e44, FloatMath ≥ 1.3× more; allocations down; differential + both-engine tests green; docs and guide updated. Clean code. Don't commit. Never mention Linkz in lizard-script.
