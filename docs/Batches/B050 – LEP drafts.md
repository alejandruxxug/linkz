# B050 – LEP drafts: classes, named arguments, FFI (design only)

Repo: `../lizard-script` · New files only: `proposals/0003-*.md`, `0004-*.md`, `0005-*.md` (+ one row each in proposals/README.md's index) · No code

Status **Draft** for all three: the developer and the team decide later. Follow `proposals/0000-template.md` and the depth of 0001/0002 (comparison with other languages, what breaks, implementation on both engines, alternatives, open questions).
1. **0003 – Classes** (for game nodes and editor objects): small classes as a new value kind — `class Node { function init(name) { self.name = name } }`, `class Sprite extends Node`, `self`, `super.method()`, single inheritance; why auto-`self` can't apply to plain hashes (modules like `math` are hashes); comparison with Lua metatables, Python, JS prototypes, Godot/GDScript, Luau; how the VM can make fields fast (shapes/hidden classes).
2. **0004 – Named arguments**: `circle(center: (1, 2), r: 5)` — a syntax error today, so additive; mixing positional + named, defaults, errors, how builtins declare names, interaction with `lizard.docs`; comparison with Python, Kotlin, Swift, C#.
3. **0005 – FFI** (calling C from lizard through the Terrarium, built on purego): `ffi.open`, `ffi.func(lib, name, args, ret)`, types (ints, floats, pointers, strings, structs, callbacks), memory blocks, safety (opt-in module, never in sandboxed `lizard.run`), per-OS library names, how window/pen could move onto it; comparison with Python ctypes, LuaJIT FFI, Deno FFI.
Real lizard examples in each (run what can run today to show the current workaround). Don't commit.
