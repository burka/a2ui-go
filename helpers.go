package a2ui

// Column creates a vertical layout component.
func Column(id string, children ...string) Component {
	return Component{
		ID:     id,
		Column: &ColumnDef{Children: children},
	}
}

// Row creates a horizontal layout component.
func Row(id string, children ...string) Component {
	return Component{
		ID:  id,
		Row: &RowDef{Children: children},
	}
}

// Card creates a card container component.
func Card(id, child string) Component {
	return Component{
		ID:   id,
		Card: &CardDef{Child: child},
	}
}

// TextStatic creates a text component with static content.
func TextStatic(id, text string) Component {
	return Component{
		ID:   id,
		Text: &TextDef{Text: text},
	}
}

// TextBound creates a text component bound to a data path.
func TextBound(id, path string) Component {
	return Component{
		ID:   id,
		Text: &TextDef{DataBinding: &DataBinding{Path: path}},
	}
}

// ImageStatic creates an image component with a static URL.
func ImageStatic(id, url, alt string) Component {
	return Component{
		ID:    id,
		Image: &ImageDef{URL: url, Alt: alt},
	}
}

// ImageBound creates an image component bound to a data path.
func ImageBound(id, path, alt string) Component {
	return Component{
		ID:    id,
		Image: &ImageDef{DataBinding: &DataBinding{Path: path}, Alt: alt},
	}
}

// Button creates a button component.
func Button(id, text, actionType string) Component {
	return Component{
		ID:     id,
		Button: &ButtonDef{Text: text, Action: Action{Type: actionType}},
	}
}

// ButtonWithData creates a button with action data.
func ButtonWithData(id, text, actionType string, data map[string]any) Component {
	return Component{
		ID:     id,
		Button: &ButtonDef{Text: text, Action: Action{Type: actionType, Data: data}},
	}
}

// TextField creates a text input component.
func TextField(id, label, placeholder string) Component {
	return Component{
		ID:        id,
		TextField: &TextFieldDef{Label: label, Placeholder: placeholder},
	}
}

// TextFieldBound creates a text input bound to a data path.
func TextFieldBound(id, label, placeholder, path string) Component {
	return Component{
		ID: id,
		TextField: &TextFieldDef{
			Label:       label,
			Placeholder: placeholder,
			DataBinding: &DataBinding{Path: path},
		},
	}
}

// ListTemplate creates a list component that renders items from data.
func ListTemplate(id, templateID, dataPath string) Component {
	return Component{
		ID: id,
		List: &ListDef{
			Template:    templateID,
			DataBinding: &DataBinding{Path: dataPath},
		},
	}
}
