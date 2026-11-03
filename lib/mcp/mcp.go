// Package mcp bundles the libraries of the MCP server (linkz mcp): a JSON
// parser and writer, and the JSON-RPC server with its tools. They are written
// in lizard-script; this file only embeds them.
package mcp

import "embed"

// Files holds the source of the MCP libraries.
//
//go:embed *.liz
var Files embed.FS

// Order is the order the libraries load in, since a library can use what the
// ones before it define.
var Order = []string{"json.liz", "tools.liz", "server.liz"}
