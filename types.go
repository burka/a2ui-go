// Package a2ui provides types and utilities for generating A2UI messages.
// A2UI is a protocol for AI agents to generate declarative UI definitions
// that clients render with native widgets.
package a2ui

// Message represents an A2UI protocol message.
// Exactly one field should be set per message.
type Message struct {
	BeginRendering  *BeginRendering  `json:"beginRendering,omitempty"`
	SurfaceUpdate   *SurfaceUpdate   `json:"surfaceUpdate,omitempty"`
	DataModelUpdate *DataModelUpdate `json:"dataModelUpdate,omitempty"`
	DeleteSurface   *DeleteSurface   `json:"deleteSurface,omitempty"`
}

// BeginRendering initializes a new UI surface.
type BeginRendering struct {
	SurfaceID string `json:"surfaceId"`
	Root      string `json:"root"`
}

// SurfaceUpdate sends component definitions to the client.
type SurfaceUpdate struct {
	SurfaceID  string      `json:"surfaceId"`
	Components []Component `json:"components"`
}

// DataModelUpdate sends data model contents to the client.
type DataModelUpdate struct {
	SurfaceID string         `json:"surfaceId"`
	Contents  map[string]any `json:"contents"`
}

// DeleteSurface removes a surface from the client.
type DeleteSurface struct {
	SurfaceID string `json:"surfaceId"`
}

// Component represents a UI component in the adjacency list.
// Exactly one component type field should be set.
type Component struct {
	ID        string        `json:"id"`
	Column    *ColumnDef    `json:"Column,omitempty"`
	Row       *RowDef       `json:"Row,omitempty"`
	Card      *CardDef      `json:"Card,omitempty"`
	Text      *TextDef      `json:"Text,omitempty"`
	Image     *ImageDef     `json:"Image,omitempty"`
	Button    *ButtonDef    `json:"Button,omitempty"`
	TextField *TextFieldDef `json:"TextField,omitempty"`
	List      *ListDef      `json:"List,omitempty"`
}

// ColumnDef is a vertical layout container.
type ColumnDef struct {
	Children []string `json:"children"`
}

// RowDef is a horizontal layout container.
type RowDef struct {
	Children []string `json:"children"`
}

// CardDef is a container with visual styling.
type CardDef struct {
	Child string `json:"child"`
}

// TextDef displays text content.
type TextDef struct {
	Text        string       `json:"text,omitempty"`
	DataBinding *DataBinding `json:"dataBinding,omitempty"`
}

// ImageDef displays an image.
type ImageDef struct {
	URL         string       `json:"url,omitempty"`
	Alt         string       `json:"alt,omitempty"`
	DataBinding *DataBinding `json:"dataBinding,omitempty"`
}

// ButtonDef triggers an action when clicked.
type ButtonDef struct {
	Text   string `json:"text"`
	Action Action `json:"action"`
}

// Action defines what happens when a button is clicked.
type Action struct {
	Type string `json:"type"`
}

// TextFieldDef accepts user text input.
type TextFieldDef struct {
	Label       string       `json:"label,omitempty"`
	Placeholder string       `json:"placeholder,omitempty"`
	DataBinding *DataBinding `json:"dataBinding,omitempty"`
}

// ListDef renders items from a data array using a template.
type ListDef struct {
	Template    string       `json:"template"`
	DataBinding *DataBinding `json:"dataBinding,omitempty"`
}

// DataBinding binds a component to a JSON Pointer path in the data model.
type DataBinding struct {
	Path string `json:"path"`
}
