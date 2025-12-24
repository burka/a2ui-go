# A2UI Angular Client

Minimal Angular client for rendering A2UI surfaces from a Go backend.

## Setup

```bash
npm install
```

## Run

First, start the Go backend:
```bash
cd ../../examples/interactive
go run main.go
```

Then start the Angular app:
```bash
npm start
# Open http://localhost:4200
```

## Features

- Renders A2UI components (Column, Row, Card, Text, TextField, Button)
- Data binding support
- Form handling with client events
- Navigation between surfaces

## Structure

- `src/app/a2ui.service.ts` - Fetches JSONL and sends events
- `src/app/a2ui-renderer.component.ts` - Renders A2UI components
- `src/app/app.ts` - Main app connecting to Go backend
