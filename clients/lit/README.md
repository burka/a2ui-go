# A2UI Lit Client

Minimal Lit web component client for rendering A2UI surfaces from a Go backend.

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

Then start the Lit dev server:
```bash
npm run dev
# Open http://localhost:5173
```

## Features

- Renders A2UI components (Column, Row, Card, Text, TextField, Button)
- Data binding support
- Form handling with client events
- Navigation between surfaces
- Uses Lit 3 with decorators

## Structure

- `src/a2ui-renderer.ts` - Web component that renders A2UI
- `index.html` - Main HTML page
