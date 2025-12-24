package a2ui

import "fmt"

// Surface builds A2UI messages for a UI surface.
type Surface struct {
	id         string
	root       string
	components []Component
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
func (s *Surface) Add(c Component) *Surface {
	s.components = append(s.components, c)
	return s
}

// AddAll appends multiple components to the surface.
func (s *Surface) AddAll(components ...Component) *Surface {
	s.components = append(s.components, components...)
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
		{SurfaceUpdate: &SurfaceUpdate{SurfaceID: s.id, Components: s.components}},
	}

	if len(s.data) > 0 {
		messages = append(messages, Message{
			DataModelUpdate: &DataModelUpdate{SurfaceID: s.id, Contents: s.data},
		})
	}

	return messages
}

// SurfaceUpdateMessage returns a SurfaceUpdate message with current components.
func (s *Surface) SurfaceUpdateMessage() Message {
	return Message{
		SurfaceUpdate: &SurfaceUpdate{SurfaceID: s.id, Components: s.components},
	}
}

// DataModelUpdateMessage returns a DataModelUpdate message with current data.
func (s *Surface) DataModelUpdateMessage() Message {
	return Message{
		DataModelUpdate: &DataModelUpdate{SurfaceID: s.id, Contents: s.data},
	}
}

// Components returns the current component list.
func (s *Surface) Components() []Component {
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
func (s *Surface) Validate() []ValidationError {
	var errors []ValidationError

	// Build a map of all component IDs
	componentIDs := make(map[string]bool)
	duplicates := make(map[string]bool)

	// First pass: collect IDs and detect duplicates/empty IDs
	for _, c := range s.components {
		if c.ID == "" {
			errors = append(errors, ValidationError{
				ComponentID: "",
				Field:       "ID",
				Message:     "component ID must not be empty",
			})
			continue
		}

		if componentIDs[c.ID] {
			if !duplicates[c.ID] {
				errors = append(errors, ValidationError{
					ComponentID: c.ID,
					Field:       "ID",
					Message:     "duplicate component ID",
				})
				duplicates[c.ID] = true
			}
		} else {
			componentIDs[c.ID] = true
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

	// Second pass: check children references
	for _, c := range s.components {
		// Skip components with empty IDs (already reported)
		if c.ID == "" {
			continue
		}

		// Check Column children
		if c.Column != nil {
			for _, child := range c.Column.Children {
				if !componentIDs[child] {
					errors = append(errors, ValidationError{
						ComponentID: c.ID,
						Field:       "Column.Children",
						Message:     fmt.Sprintf("child '%s' not found", child),
					})
				}
			}
		}

		// Check Row children
		if c.Row != nil {
			for _, child := range c.Row.Children {
				if !componentIDs[child] {
					errors = append(errors, ValidationError{
						ComponentID: c.ID,
						Field:       "Row.Children",
						Message:     fmt.Sprintf("child '%s' not found", child),
					})
				}
			}
		}

		// Check Card child
		if c.Card != nil {
			if !componentIDs[c.Card.Child] {
				errors = append(errors, ValidationError{
					ComponentID: c.ID,
					Field:       "Card.Child",
					Message:     fmt.Sprintf("child '%s' not found", c.Card.Child),
				})
			}
		}

		// Check List template
		if c.List != nil {
			if !componentIDs[c.List.Template] {
				errors = append(errors, ValidationError{
					ComponentID: c.ID,
					Field:       "List.Template",
					Message:     fmt.Sprintf("template '%s' not found", c.List.Template),
				})
			}
		}

		// Check Tabs children
		if c.Tabs != nil {
			for i, tab := range c.Tabs.Tabs {
				if !componentIDs[tab.Child] {
					errors = append(errors, ValidationError{
						ComponentID: c.ID,
						Field:       fmt.Sprintf("Tabs.Tabs[%d].Child", i),
						Message:     fmt.Sprintf("child '%s' not found", tab.Child),
					})
				}
			}
		}

		// Check Modal children
		if c.Modal != nil {
			if !componentIDs[c.Modal.EntryPointChild] {
				errors = append(errors, ValidationError{
					ComponentID: c.ID,
					Field:       "Modal.EntryPointChild",
					Message:     fmt.Sprintf("entryPointChild '%s' not found", c.Modal.EntryPointChild),
				})
			}
			if !componentIDs[c.Modal.ContentChild] {
				errors = append(errors, ValidationError{
					ComponentID: c.ID,
					Field:       "Modal.ContentChild",
					Message:     fmt.Sprintf("contentChild '%s' not found", c.Modal.ContentChild),
				})
			}
		}
	}

	return errors
}
