package a2ui

import "fmt"

// Surface builds A2UI messages for a UI surface.
type Surface struct {
	id         string
	root       string
	components []any
	data       map[string]any
}

// NewSurface creates a new surface with the given ID.
// The root component ID defaults to "root".
func NewSurface(id string) *Surface {
	return &Surface{
		id:   id,
		root: "root",
		data: make(map[string]any),
	}
}

// SetRoot sets the root component ID.
func (s *Surface) SetRoot(id string) *Surface {
	s.root = id
	return s
}

// Add appends a component to the surface.
// The component can be of type Component or any custom struct with embedded Component.
func (s *Surface) Add(c any) *Surface {
	s.components = append(s.components, c)
	return s
}

// AddAll appends multiple standard components to the surface.
// For adding custom components, use Add() individually.
func (s *Surface) AddAll(components ...Component) *Surface {
	for _, c := range components {
		s.components = append(s.components, c)
	}
	return s
}

// SetData sets a value at the given JSON Pointer path.
func (s *Surface) SetData(path string, value any) *Surface {
	s.data[path] = value
	return s
}

// Messages returns the complete message sequence for this surface.
func (s *Surface) Messages() []Message {
	messages := []Message{
		{BeginRendering: &BeginRendering{SurfaceID: s.id, Root: s.root}},
		{UpdateComponents: &UpdateComponents{SurfaceID: s.id, Components: s.components}},
	}

	if len(s.data) > 0 {
		messages = append(messages, Message{
			DataModelUpdate: &DataModelUpdate{SurfaceID: s.id, Contents: s.data},
		})
	}

	return messages
}

// UpdateComponentsMessage returns an UpdateComponents message with current components.
func (s *Surface) UpdateComponentsMessage() Message {
	return Message{
		UpdateComponents: &UpdateComponents{SurfaceID: s.id, Components: s.components},
	}
}

// DataModelUpdateMessage returns a DataModelUpdate message with current data.
func (s *Surface) DataModelUpdateMessage() Message {
	return Message{
		DataModelUpdate: &DataModelUpdate{SurfaceID: s.id, Contents: s.data},
	}
}

// Components returns the current component list.
func (s *Surface) Components() []any {
	return s.components
}

// ValidationError represents a validation error for a component.
type ValidationError struct {
	ComponentID string
	Field       string
	Message     string
}

func (e ValidationError) Error() string {
	if e.ComponentID != "" {
		return fmt.Sprintf("%s: %s - %s", e.ComponentID, e.Field, e.Message)
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Validate checks the surface for structural errors.
// It returns a slice of validation errors (empty slice means valid).
// Note: Validation only works for standard Component types. Custom components
// with embedded Component are validated for their embedded fields only.
func (s *Surface) Validate() []ValidationError {
	var errors []ValidationError

	// Build a map of all component IDs
	componentIDs := make(map[string]bool)
	duplicates := make(map[string]bool)

	// Helper to extract Component from any type
	getComponent := func(c any) *Component {
		switch v := c.(type) {
		case Component:
			return &v
		case *Component:
			return v
		default:
			return nil
		}
	}

	// First pass: collect IDs and detect duplicates/empty IDs
	for _, c := range s.components {
		comp := getComponent(c)
		if comp == nil {
			// Custom component - skip validation but try to get ID via reflection would be complex
			// For now, custom components are not validated
			continue
		}

		if comp.ID == "" {
			errors = append(errors, ValidationError{
				ComponentID: "",
				Field:       "ID",
				Message:     "component ID must not be empty",
			})
			continue
		}

		if componentIDs[comp.ID] {
			if !duplicates[comp.ID] {
				errors = append(errors, ValidationError{
					ComponentID: comp.ID,
					Field:       "ID",
					Message:     "duplicate component ID",
				})
				duplicates[comp.ID] = true
			}
		} else {
			componentIDs[comp.ID] = true
		}
	}

	// Check if root exists
	if !componentIDs[s.root] {
		errors = append(errors, ValidationError{
			ComponentID: s.root,
			Field:       "Root",
			Message:     "root component not found",
		})
	}

	// Second pass: check children references based on component type
	for _, c := range s.components {
		comp := getComponent(c)
		if comp == nil {
			continue
		}

		// Skip components with empty IDs (already reported)
		if comp.ID == "" {
			continue
		}

		switch comp.Component {
		case "Column", "Row":
			for _, child := range comp.Children {
				if !componentIDs[child] {
					errors = append(errors, ValidationError{
						ComponentID: comp.ID,
						Field:       comp.Component + ".Children",
						Message:     fmt.Sprintf("child '%s' not found", child),
					})
				}
			}

		case "Card":
			if comp.Child != "" && !componentIDs[comp.Child] {
				errors = append(errors, ValidationError{
					ComponentID: comp.ID,
					Field:       "Card.Child",
					Message:     fmt.Sprintf("child '%s' not found", comp.Child),
				})
			}

		case "Button":
			if comp.Child != "" && !componentIDs[comp.Child] {
				errors = append(errors, ValidationError{
					ComponentID: comp.ID,
					Field:       "Button.Child",
					Message:     fmt.Sprintf("child '%s' not found", comp.Child),
				})
			}

		case "List":
			if comp.Template != "" && !componentIDs[comp.Template] {
				errors = append(errors, ValidationError{
					ComponentID: comp.ID,
					Field:       "List.Template",
					Message:     fmt.Sprintf("template '%s' not found", comp.Template),
				})
			}

		case "Tabs":
			for i, tab := range comp.Tabs {
				if !componentIDs[tab.Child] {
					errors = append(errors, ValidationError{
						ComponentID: comp.ID,
						Field:       fmt.Sprintf("Tabs.Tabs[%d].Child", i),
						Message:     fmt.Sprintf("child '%s' not found", tab.Child),
					})
				}
			}

		case "Modal":
			if comp.EntryPointChild != "" && !componentIDs[comp.EntryPointChild] {
				errors = append(errors, ValidationError{
					ComponentID: comp.ID,
					Field:       "Modal.EntryPointChild",
					Message:     fmt.Sprintf("entryPointChild '%s' not found", comp.EntryPointChild),
				})
			}
			if comp.ContentChild != "" && !componentIDs[comp.ContentChild] {
				errors = append(errors, ValidationError{
					ComponentID: comp.ID,
					Field:       "Modal.ContentChild",
					Message:     fmt.Sprintf("contentChild '%s' not found", comp.ContentChild),
				})
			}
		}
	}

	return errors
}
