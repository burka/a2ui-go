# A2UI-Go

Minimal Go library for generating A2UI messages. Zero external dependencies.

## Constraints

- `a2ui/` package: **stdlib only** - no external dependencies allowed
- `examples/` may use external deps if needed

## Before Every Commit

```bash
go fmt ./...
go vet ./...
go mod tidy
go test ./...
```

All must pass with zero errors/warnings.

## Architecture

- `types.go` - Message & component structs (oneOf pattern)
- `builder.go` - Surface builder (`Add`, `SetData`, `Messages`)
- `helpers.go` - Component constructors (`Column`, `TextStatic`, etc.)
- `writer.go` - Output functions (`WriteJSONL`, `WritePretty`)

Components use flat adjacency list - children referenced by ID, not nested.
