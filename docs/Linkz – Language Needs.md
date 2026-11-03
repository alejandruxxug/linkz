# 🦎 Linkz: Language Needs

Every time Linkz hits something lizard-script can't do (or does slowly), it goes here. When it's built, the change lands in the lizard-script repo and gets a row in its `docs/DESIGN.md` → *Extensions beyond the course design*.

## 🧱 Gaps

| Gap | Hit by | Workaround now | Priority |
|---|---|---|---|
| No `\"` or `\n` escapes in strings | `lib/svg.liz` | Single-quoted XML attributes; a literal line break inside a string | High (JSON and the MCP need `"`) |
| No floats | Curves, rotation, color math, 3D | Integer coordinates only | High |
| No dot access (`shapes.circle`) | Library API | Flat global names | High |
| No `use` / modules | Every library | The host loads `lib/` in a fixed order into one environment | High |
| No way to raise an error from lizard | `svg_shape` on an unknown type | Writes an SVG comment instead | Medium |
| No error handling (`try`) | The MCP server must survive bad requests | None yet | Medium (before the MCP) |
| Tree-walking speed | Rasterizer, tracing | Not hit yet; measure first | Will be high |
| No string builder / bytes | PNG and JSON output | `+` on strings | Medium |

## 📓 Log

| Date | Need | What happened |
|---|---|---|
| 2026-11-02 | Escapes, floats, dot, `use`, raise | Found while writing the first `svg.liz` / `shapes.liz` |
