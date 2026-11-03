# 🔌 Linkz: MCP

`linkz mcp` is an [MCP](https://modelcontextprotocol.io) server over stdio, so an AI agent (Claude Code, Claude Desktop, Cursor, any MCP client) can draw in Linkz. The agent **writes lizard-script, never raw SVG**: it reads the drawing API with `docs`, checks code with `run`, **looks at the picture** with `render` (a PNG, rasterized by lizard itself), and keeps the result with `save`, so the file shows up in the Linkz editor.

No network is used. The server talks only through stdin and stdout.

## ➕ Add it to a client

Build once, then point the client at the binary:

```bash
cd linkz && go build -o linkz .
```

**Claude Code**

```bash
claude mcp add linkz -- /absolute/path/to/linkz mcp
```

(`claude mcp add --scope user linkz -- /absolute/path/to/linkz mcp` to have it in every project, `claude mcp list` to check, `/mcp` inside a session to see its tools.)

**Claude Desktop, Cursor and other clients** take a JSON entry with a command and arguments (`claude_desktop_config.json`, `.cursor/mcp.json`, ...):

```json
{
  "mcpServers": {
    "linkz": {
      "command": "/absolute/path/to/linkz",
      "args": ["mcp"]
    }
  }
}
```

To try it by hand: `./linkz mcp`, then paste one JSON message per line (see the session below). Logs go to stderr; stdout carries only protocol messages.

## 🧰 The tools

| Tool | Arguments | Gives |
|---|---|---|
| `docs` | `name?` | The primer (canvas, layers and their order, points, colors, gradients, blend modes, pens, the language in brief) and every drawing function with its parameters, summary and example. With `name`, the full doc of one function (a near miss lists the close names). **Start here.** |
| `run` | `code` | Text JSON `{"ok", "errors": [{"line", "column", "message"}], "printed", "shapes", "steps", "ms"}`. No picture; fast check. |
| `render` | `code`, `width?`, `height?` | An **image** item (`image/png`, base64) of the drawing at the canvas size, scaled to fit `width` / `height` (1 to 4096) keeping its shape. When the code has errors, a text item with them follows, and the image shows what was drawn before the first error. |
| `svg` | `code` | The SVG text (an error result with the `errors` JSON if the code fails). |
| `save` | `path`, `code`, `overwrite?` | Writes a `.lzdraw`. Answers `{"ok", "path", "bytes", "errors"}` (a file with syntax errors is saved, and the errors are reported). |
| `open` | `path` | The code of a `.lzdraw`. |

Every tool takes the lizard-script source as a string, the same text a `.lzdraw` file holds. A drawing that has errors is **not** a tool failure: `run` and `render` report the errors as data (with line and column) so the agent can fix them and call again. Bad arguments (a missing `code`, a canvas bigger than 4096 pixels, a width of 0) come back as a result with `isError: true` and a message.

### Protocol

JSON-RPC 2.0, one JSON message per line. Implemented: `initialize` (answers the client's `protocolVersion` when it is `2025-06-18`, `2025-03-26` or `2024-11-05`, else `2025-06-18`), `notifications/*` (ignored, no reply), `ping`, `tools/list`, `tools/call`; anything else is `-32601`. A line that is not JSON gets `-32700` and the server reads on; a message that is not a request gets `-32600`; a missing tool or bad `params` is `-32602`; a failure inside the server is `-32603`. Batches (arrays) work. At the end of stdin the server says bye on stderr and exits 0.

## 🧪 An example session

Lines starting `>` are what the client sends, `<` what comes back (long values cut with `...`).

```
> {"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"me","version":"1"}}}
< {"jsonrpc":"2.0","id":1,"result":{"protocolVersion":"2025-06-18","capabilities":{"tools":{"listChanged":false}},"serverInfo":{"name":"linkz","version":"0.1.0"},"instructions":"Linkz draws vector art ..."}}
> {"jsonrpc":"2.0","method":"notifications/initialized"}
> {"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"run","arguments":{"code":"canvas = {\"size\": (200, 100), \"background\": #f4f1e8}\ncircle((100, 50), 30, #3fa34d)\nputs(\"drawn\")"}}}
< {"jsonrpc":"2.0","id":2,"result":{"content":[{"type":"text","text":"{\"ok\":true,\"errors\":[],\"printed\":\"drawn\\n\",\"shapes\":1,\"steps\":24,\"ms\":2.1}"}],"isError":false}}
> {"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"render","arguments":{"code":"canvas = ...","width":400}}}
< {"jsonrpc":"2.0","id":3,"result":{"content":[{"type":"image","data":"iVBORw0KGgo...","mimeType":"image/png"}],"isError":false}}
> {"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"save","arguments":{"path":"pond","code":"canvas = ..."}}}
< {"jsonrpc":"2.0","id":4,"result":{"content":[{"type":"text","text":"{\"ok\":true,\"path\":\"/Users/me/Documents/Linkz/pond.lzdraw\",\"bytes\":112,\"errors\":[]}"}],"isError":false}}
```

In practice the agent's loop is: `docs`, write code, `render`, look, fix, `render`, then `save`. Ask it for example: "Using the linkz tools, draw a lighthouse at dusk, render it, improve it until it looks good, and save it as lighthouse."

## 🛡️ Safety rules

- **Sandboxed runs.** Every drawing runs in a fresh lizard environment with the drawing libraries only: no files, no clipboard, no OS. It stops after 2,000,000 steps (loops and calls) with a normal error (`the program ran too long`), so an endless loop cannot hang a call; the server answers the next request. A tool's own work (painting the picture) has a larger budget of 300,000,000 steps per request. Canvases above 4096 by 4096 are refused. What a drawing prints is cut at 20,000 characters.
- **Writes are narrow.** `save` and `open` handle only `.lzdraw` files (a name without an extension gets `.lzdraw`). A relative path lives inside the drawings folder (`settings.drawings_dir`, default `~/Documents/Linkz`; `..` may not leave it). An absolute path (or one starting with `~/`) must be inside the user's home folder, or the drawings folder. No path may have a part starting with a dot, so `~/.ssh` and `~/.linkz` are out. Folders are created as needed.
- **No silent overwrites.** `save` refuses an existing file unless `overwrite` is exactly `true`. Writes go through the atomic `fs.write` (a temp file, then a rename), so a crash never leaves half a file.
- **No network**, no shell, nothing else is reachable from a tool.
- A request line over 32 MB is refused.

Known limits: paths are checked as text, so a symbolic link inside the allowed folders can still point elsewhere; the server's own step limit and a stack overflow cannot be caught inside lizard, so a pathological `render` (hundreds of thousands of huge shapes) would end the server process, and the client has to start it again. JSON integers beyond 64 bits are read as floats.

## 🧱 How it is built

All logic is lizard-script in `lib/mcp/`, bundled into the binary:

| File | Does |
|---|---|
| `json.liz` | JSON parser and writer: objects to hashes, `\uXXXX` and surrogate pairs, ints and floats, errors with line and column, a nesting limit of 100 |
| `tools.liz` | The six tools, the drawing primer, path rules, scaling shapes for `render` |
| `server.liz` | JSON-RPC: `mcp_handle_line`, `mcp_dispatch`, the `mcp_serve` loop; every request in try/catch |

Go (`host/stdio.go`) only adds the OS primitives: `io.read_line()` (stdin, `null` at the end), `io.write(s)` (stdout, flushed), `io.log(s)` (stderr), `io.step_limit(n)` (the per-request budget) and `fb.png_base64()` (the current target as a base64 PNG). `render` is `fb.offscreen` at the canvas size, the background, then `raster.paint`, so the pixels are the ones the editor paints. Scaling to `width` / `height` scales the shapes (paths, widths, pinned gradients) before painting, so the picture is sharp, not a resized bitmap.

## ✅ Tests

- `lib/mcp/tests/json_*.liz`: parsing, strings and unicode, writing, errors; `server_*.liz`: JSON-RPC, every tool, the path rules (run with `HOME` in a temp folder). Output is compared with the `.expected` files (`go test ./lib/mcp -update` rewrites them).
- `mcp_e2e_test.go`: builds `linkz`, starts `linkz mcp`, and sends initialize, tools/list, run, render (decodes the PNG, checks its size and pixels, also scaled), a malformed line (parse error, the server carries on), an unknown method, an endless drawing, a refused save, then EOF (exit 0).

Timings (Apple silicon): a first `render` of `examples/layers.lzdraw` in a fresh server takes about 25 ms, and 6 ms once the shapes are cached; startup (loading the libraries) about 30 ms.
