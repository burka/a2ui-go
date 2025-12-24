package a2ui

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
