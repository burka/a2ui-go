# a2ui-go

Minimal Go implementation of Google's A2UI protocol for agent-driven interfaces.

## What is A2UI?

A2UI lets AI agents generate native UIs safely by sending declarative JSON instead of code. Agents describe what components to show, clients render them with native widgets.

**Key benefits:**
- Secure - no code execution, just data
- Cross-platform - same JSON renders on web/mobile/desktop
- Streaming - progressive UI updates in real-time
- Simple - flat component list, easy for LLMs to generate

## Installation

```bash
go get github.com/burka/a2ui-go
```

**Zero dependencies** - pure stdlib only.

## Quick Start

```go
package main

import (
    "net/http"
    a2ui "github.com/burka/a2ui-go"
)

func main() {
    http.HandleFunc("/ui", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/x-ndjson")

        // Build UI
        surface := a2ui.NewSurface("demo")
        surface.Add(a2ui.Column("root", "title", "message"))
        surface.Add(a2ui.TextStatic("title", "Hello A2UI"))
        surface.Add(a2ui.TextStatic("message", "This is a minimal example"))

        // Stream JSONL
        a2ui.WriteJSONL(w, surface.Messages())
    })

    http.ListenAndServe(":8080", nil)
}
```

Test it:
```bash
curl http://localhost:8080/ui
```

## Core API

### Building UIs

```go
// Create surface
surface := a2ui.NewSurface("my-surface")

// Add components
surface.Add(a2ui.Column("root", "child1", "child2"))
surface.Add(a2ui.TextStatic("child1", "Hello"))
surface.Add(a2ui.Card("child2", "content"))

// Bind data
surface.Add(a2ui.TextBound("name", "/user/name"))
surface.SetData("/user/name", "Alice")

// List with template
surface.Add(a2ui.ListTemplate("list", "item-card", "/items"))
surface.SetData("/items", []Product{{Name: "Item 1"}, {Name: "Item 2"}})

// Get messages
messages := surface.Messages()
```

### Component Helpers

```go
// Layout
a2ui.Column(id, children...)                          // Vertical layout
a2ui.ColumnWithLayout(id, dist, align, children...)   // With distribution/alignment
a2ui.Row(id, children...)                             // Horizontal layout
a2ui.RowWithLayout(id, dist, align, children...)      // With distribution/alignment
a2ui.Card(id, child)                                  // Card container
a2ui.ListTemplate(id, templateID, path)               // Data-bound list
a2ui.Tabs(id, tabs...)                                // Tabbed container
a2ui.Tab(title, child)                                // Tab definition
a2ui.Modal(id, entryPoint, content)                   // Modal overlay

// Content
a2ui.TextStatic(id, text)                             // Static text
a2ui.TextBound(id, path)                              // Data-bound text
a2ui.TextWithHint(id, text, hint)                     // Text with usage hint
a2ui.ImageStatic(id, url, alt)                        // Static image
a2ui.ImageBound(id, path, alt)                        // Data-bound image
a2ui.ImageWithFit(id, url, alt, fit)                  // Image with fit option
a2ui.Icon(id, iconName)                               // Predefined icon
a2ui.Video(id, url)                                   // Video player
a2ui.VideoBound(id, path)                             // Data-bound video
a2ui.AudioPlayer(id, url, description)                // Audio player
a2ui.AudioPlayerBound(id, path, description)          // Data-bound audio
a2ui.Divider(id)                                      // Horizontal divider
a2ui.DividerVertical(id)                              // Vertical divider

// Form
a2ui.Button(id, text, actionType)                     // Button with action
a2ui.ButtonPrimary(id, text, actionType)              // Primary button
a2ui.ButtonWithData(id, text, type, data)             // Button with action data
a2ui.TextField(id, label, placeholder)                // Text input
a2ui.TextFieldBound(id, label, ph, path)              // Data-bound text input
a2ui.TextFieldWithType(id, label, ph, type)           // Text input with type
a2ui.CheckBox(id, label, value)                       // Checkbox
a2ui.CheckBoxBound(id, label, path)                   // Data-bound checkbox
a2ui.DateTimeInput(id, label, date, time)             // Date/time picker
a2ui.DateTimeInputBound(id, label, path, date, time)  // Data-bound date/time
a2ui.MultipleChoice(id, label, options)               // Multiple choice
a2ui.MultipleChoiceBound(id, label, path, options)    // Data-bound choice
a2ui.Choice(label, value)                             // Choice option
a2ui.Slider(id, label, min, max, value)               // Numeric slider
a2ui.SliderBound(id, label, path, min, max)           // Data-bound slider
```

### Writing Output

```go
// JSONL streaming (for A2UI clients)
a2ui.WriteJSONL(w, messages)

// Single message
a2ui.WriteMessage(w, msg)

// Pretty JSON (for debugging)
a2ui.WritePretty(w, messages)
```

## Examples

### Static UI

```go
surface := a2ui.NewSurface("greeting")
surface.Add(a2ui.Column("root", "title", "subtitle"))
surface.Add(a2ui.TextStatic("title", "Welcome"))
surface.Add(a2ui.TextStatic("subtitle", "Getting started with A2UI"))
```

### Data-Bound List

```go
type Product struct {
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

surface := a2ui.NewSurface("products")
surface.Add(a2ui.Column("root", "list"))
surface.Add(a2ui.ListTemplate("list", "card", "/products"))

// Template for each item
surface.Add(a2ui.Card("card", "content"))
surface.Add(a2ui.Column("content", "name", "price"))
surface.Add(a2ui.TextBound("name", "/name"))
surface.Add(a2ui.TextBound("price", "/price"))

// Bind data
surface.SetData("/products", []Product{
    {"Widget", 29.99},
    {"Gadget", 49.99},
})
```

### Progressive Streaming

For real-time UI updates, build components incrementally and flush between updates:

```go
w.Header().Set("Content-Type", "application/x-ndjson")
flusher := w.(http.Flusher)

surface := a2ui.NewSurface("demo")
surface.Add(a2ui.Column("root", "loading"))
surface.Add(a2ui.TextStatic("loading", "Processing..."))

// Send initial UI
a2ui.WriteJSONL(w, surface.Messages())
flusher.Flush()

// Add components progressively
surface.Add(a2ui.Card("result", "text"))
surface.Add(a2ui.TextBound("text", "/result"))
surface.SetData("/result", "Done!")
surface.Add(a2ui.Column("root", "result"))

// Send updates
a2ui.WriteMessage(w, surface.SurfaceUpdateMessage())
a2ui.WriteMessage(w, surface.DataModelUpdateMessage())
flusher.Flush()
```

The pattern: initial surface → add components → set data → `WriteMessage()` + `Flush()` for each update.

### Interactive Forms

Handle user events with `ClientMessage` and `Event` types:

```go
// Button with action data
surface.Add(a2ui.ButtonWithData("submit", "Book", "submit",
    map[string]any{"endpoint": "/api/book"}))

// Handle client events
func handleSubmit(w http.ResponseWriter, r *http.Request) {
    var msg a2ui.ClientMessage
    json.NewDecoder(r.Body).Decode(&msg)

    // Access event data
    name := msg.Event.Data["name"].(string)

    // Send response UI
    surface := a2ui.NewSurface("confirmation")
    // ...
}
```

## Running Examples

**Streaming** - Progressive rendering:
```bash
cd examples/streaming && go run main.go
# Open http://localhost:8080
```

**Interactive** - Form with client events:
```bash
cd examples/interactive && go run main.go
# Open http://localhost:8080
```

## Project Structure

```
a2ui-go/
├── types.go         # Message & component types
├── builder.go       # Surface builder
├── helpers.go       # Component constructors
├── writer.go        # I/O functions
├── a2ui_test.go     # Tests
├── examples/
│   ├── streaming/   # Progressive rendering
│   └── interactive/ # Forms with client events
└── clients/
    ├── angular/     # Angular client renderer
    └── lit/         # Lit web component renderer
```

## Client Examples

Start the Go backend first, then run a client:

```bash
# Terminal 1: Go backend
cd examples/interactive && go run main.go

# Terminal 2: Angular client (port 4200)
cd clients/angular && npm install && npm start

# Or: Lit client (port 5173)
cd clients/lit && npm install && npm run dev
```

## Protocol Details

A2UI uses JSON Lines (JSONL) - one JSON object per line:

```json
{"beginRendering":{"surfaceId":"demo","root":"root"}}
{"surfaceUpdate":{"surfaceId":"demo","components":[...]}}
{"dataModelUpdate":{"surfaceId":"demo","contents":{"/items":[...]}}}
```

**Message flow:**
1. `beginRendering` - Initialize surface with root component ID
2. `surfaceUpdate` - Send component tree (adjacency list)
3. `dataModelUpdate` - Send data model (JSON Pointer paths)

Components reference each other by ID (flat list, not nested JSON).

## A2UI Spec

This implements [A2UI v0.8](https://a2ui.org/specification/v0.8-a2ui/) with **100% component coverage**.

**All 18 standard components supported:**
- Layout: `Column`, `Row`, `Card`, `List`, `Tabs`, `Modal`
- Content: `Text`, `Image`, `Icon`, `Video`, `AudioPlayer`, `Divider`
- Form: `Button`, `TextField`, `CheckBox`, `DateTimeInput`, `MultipleChoice`, `Slider`

**Extended properties supported:**
- Layout: `Distribution`, `Alignment`
- Text: `UsageHint` (h1-h5, body, caption)
- Image: `ImageFit` (contain, cover, fill, none, scale-down)
- TextField: `TextFieldType` (shortText, longText, number, date, obscured)
- Icons: 14 standard icons (accountCircle, add, arrowBack, check, close, delete, edit, favorite, home, menu, search, settings, star, warning)

Clients can extend with custom components.

## Design Philosophy

**Minimal** - No external dependencies, ~300 LOC core library
**Idiomatic** - Simple Go patterns, no magic
**Practical** - Built for HTTP streaming and LLM generation
**Flexible** - Use helpers or build `Component` structs directly

## Client Rendering

This is the **server library** for generating A2UI. For rendering:

- **Flutter**: [GenUI SDK](https://docs.flutter.dev/ai/genui)
- **Web**: [A2UI Lit Renderer](https://github.com/google/A2UI/tree/main/renderers/lit)
- **Angular**: [A2UI Angular Renderer](https://github.com/google/A2UI/tree/main/renderers/angular)
- **React**: Community implementations in progress

## License

MIT

## Links

- [A2UI Spec](https://a2ui.org/)
- [Google A2UI Repo](https://github.com/google/A2UI)
- [A2UI Documentation](https://a2ui.org/specification/v0.8-a2ui/)
