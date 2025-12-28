package a2ui

// Column creates a vertical layout component.
func Column(id string, children ...string) Component {
	return Component{
		ID:        id,
		Component: "Column",
		Children:  children,
	}
}

// Row creates a horizontal layout component.
func Row(id string, children ...string) Component {
	return Component{
		ID:        id,
		Component: "Row",
		Children:  children,
	}
}

// Card creates a card container component.
func Card(id, child string) Component {
	return Component{
		ID:        id,
		Component: "Card",
		Child:     child,
	}
}

// TextStatic creates a text component with static content.
func TextStatic(id, text string) Component {
	return Component{
		ID:        id,
		Component: "Text",
		Text:      text,
	}
}

// TextBound creates a text component bound to a data path.
func TextBound(id, path string) Component {
	return Component{
		ID:          id,
		Component:   "Text",
		DataBinding: &DataBinding{Path: path},
	}
}

// ImageStatic creates an image component with a static URL.
func ImageStatic(id, url, alt string) Component {
	return Component{
		ID:        id,
		Component: "Image",
		URL:       url,
		Alt:       alt,
	}
}

// ImageBound creates an image component bound to a data path.
func ImageBound(id, path, alt string) Component {
	return Component{
		ID:          id,
		Component:   "Image",
		DataBinding: &DataBinding{Path: path},
		Alt:         alt,
	}
}

// Button creates a button component with a child text component.
// Returns []Component: the button and its text child.
func Button(id, text, actionType string) []Component {
	childID := id + "_text"
	return []Component{
		{
			ID:        id,
			Component: "Button",
			Child:     childID,
			Action:    &Action{Type: actionType},
		},
		{
			ID:        childID,
			Component: "Text",
			Text:      text,
		},
	}
}

// ButtonWithData creates a button with action data.
// Returns []Component: the button and its text child.
func ButtonWithData(id, text, actionType string, data map[string]any) []Component {
	childID := id + "_text"
	return []Component{
		{
			ID:        id,
			Component: "Button",
			Child:     childID,
			Action:    &Action{Type: actionType, Data: data},
		},
		{
			ID:        childID,
			Component: "Text",
			Text:      text,
		},
	}
}

// ButtonPrimary creates a primary styled button.
// Returns []Component: the button and its text child.
func ButtonPrimary(id, text, actionType string) []Component {
	childID := id + "_text"
	return []Component{
		{
			ID:        id,
			Component: "Button",
			Child:     childID,
			Action:    &Action{Type: actionType},
			Primary:   true,
		},
		{
			ID:        childID,
			Component: "Text",
			Text:      text,
		},
	}
}

// ButtonOnly creates just the button component without the child text.
// Use this when you want to manage the child component separately.
func ButtonOnly(id, childID, actionType string) Component {
	return Component{
		ID:        id,
		Component: "Button",
		Child:     childID,
		Action:    &Action{Type: actionType},
	}
}

// TextField creates a text input component.
func TextField(id, label, placeholder string) Component {
	return Component{
		ID:          id,
		Component:   "TextField",
		Label:       label,
		Placeholder: placeholder,
	}
}

// TextFieldBound creates a text input bound to a data path.
func TextFieldBound(id, label, placeholder, path string) Component {
	return Component{
		ID:          id,
		Component:   "TextField",
		Label:       label,
		Placeholder: placeholder,
		DataBinding: &DataBinding{Path: path},
	}
}

// ListTemplate creates a list component that renders items from data.
func ListTemplate(id, templateID, dataPath string) Component {
	return Component{
		ID:          id,
		Component:   "List",
		Template:    templateID,
		DataBinding: &DataBinding{Path: dataPath},
	}
}

// Tabs creates a tabbed container component.
func Tabs(id string, tabs ...TabDef) Component {
	return Component{
		ID:        id,
		Component: "Tabs",
		Tabs:      tabs,
	}
}

// Tab creates a tab definition for use with Tabs.
func Tab(title, child string) TabDef {
	return TabDef{Title: title, Child: child}
}

// Modal creates a modal overlay component.
func Modal(id, entryPointChild, contentChild string) Component {
	return Component{
		ID:              id,
		Component:       "Modal",
		EntryPointChild: entryPointChild,
		ContentChild:    contentChild,
	}
}

// Icon creates an icon component.
func Icon(id string, icon IconName) Component {
	return Component{
		ID:        id,
		Component: "Icon",
		Icon:      icon,
	}
}

// Video creates a video player component with a static URL.
func Video(id, url string) Component {
	return Component{
		ID:        id,
		Component: "Video",
		URL:       url,
	}
}

// VideoBound creates a video player bound to a data path.
func VideoBound(id, path string) Component {
	return Component{
		ID:          id,
		Component:   "Video",
		DataBinding: &DataBinding{Path: path},
	}
}

// AudioPlayer creates an audio player component.
func AudioPlayer(id, url, description string) Component {
	return Component{
		ID:          id,
		Component:   "AudioPlayer",
		URL:         url,
		Description: description,
	}
}

// AudioPlayerBound creates an audio player bound to a data path.
func AudioPlayerBound(id, path, description string) Component {
	return Component{
		ID:          id,
		Component:   "AudioPlayer",
		DataBinding: &DataBinding{Path: path},
		Description: description,
	}
}

// Divider creates a visual separator component.
func Divider(id string) Component {
	return Component{
		ID:        id,
		Component: "Divider",
	}
}

// DividerVertical creates a vertical divider.
func DividerVertical(id string) Component {
	return Component{
		ID:          id,
		Component:   "Divider",
		Orientation: "vertical",
	}
}

// CheckBox creates a checkbox input component.
func CheckBox(id, label string, checked bool) Component {
	return Component{
		ID:        id,
		Component: "CheckBox",
		Label:     label,
		Checked:   checked,
	}
}

// CheckBoxBound creates a checkbox bound to a data path.
func CheckBoxBound(id, label, path string) Component {
	return Component{
		ID:          id,
		Component:   "CheckBox",
		Label:       label,
		DataBinding: &DataBinding{Path: path},
	}
}

// DateTimeInput creates a date/time picker component.
func DateTimeInput(id, label string, enableDate, enableTime bool) Component {
	return Component{
		ID:         id,
		Component:  "DateTimeInput",
		Label:      label,
		EnableDate: enableDate,
		EnableTime: enableTime,
	}
}

// DateTimeInputBound creates a date/time picker bound to a data path.
func DateTimeInputBound(id, label, path string, enableDate, enableTime bool) Component {
	return Component{
		ID:          id,
		Component:   "DateTimeInput",
		Label:       label,
		DataBinding: &DataBinding{Path: path},
		EnableDate:  enableDate,
		EnableTime:  enableTime,
	}
}

// MultipleChoice creates a multiple choice selector component.
func MultipleChoice(id, label string, options []ChoiceOption) Component {
	return Component{
		ID:        id,
		Component: "MultipleChoice",
		Label:     label,
		Options:   options,
	}
}

// MultipleChoiceBound creates a multiple choice selector bound to a data path.
func MultipleChoiceBound(id, label, path string, options []ChoiceOption) Component {
	return Component{
		ID:          id,
		Component:   "MultipleChoice",
		Label:       label,
		Options:     options,
		DataBinding: &DataBinding{Path: path},
	}
}

// Choice creates a choice option for use with MultipleChoice.
func Choice(label, value string) ChoiceOption {
	return ChoiceOption{Label: label, Value: value}
}

// Slider creates a numeric slider component.
func Slider(id, label string, min, max, value float64) Component {
	return Component{
		ID:          id,
		Component:   "Slider",
		Label:       label,
		MinValue:    min,
		MaxValue:    max,
		SliderValue: value,
	}
}

// SliderBound creates a slider bound to a data path.
func SliderBound(id, label, path string, min, max float64) Component {
	return Component{
		ID:          id,
		Component:   "Slider",
		Label:       label,
		MinValue:    min,
		MaxValue:    max,
		DataBinding: &DataBinding{Path: path},
	}
}

// TextWithHint creates a text component with a usage hint.
func TextWithHint(id, text string, hint UsageHint) Component {
	return Component{
		ID:        id,
		Component: "Text",
		Text:      text,
		UsageHint: hint,
	}
}

// ImageWithFit creates an image with fit option.
func ImageWithFit(id, url, alt string, fit ImageFit) Component {
	return Component{
		ID:        id,
		Component: "Image",
		URL:       url,
		Alt:       alt,
		Fit:       fit,
	}
}

// TextFieldWithType creates a text field with a specific type.
func TextFieldWithType(id, label, placeholder string, fieldType TextFieldType) Component {
	return Component{
		ID:            id,
		Component:     "TextField",
		Label:         label,
		Placeholder:   placeholder,
		TextFieldType: fieldType,
	}
}

// ColumnWithLayout creates a column with distribution and alignment.
func ColumnWithLayout(id string, distribution Distribution, alignment Alignment, children ...string) Component {
	return Component{
		ID:           id,
		Component:    "Column",
		Children:     children,
		Distribution: distribution,
		Alignment:    alignment,
	}
}

// RowWithLayout creates a row with distribution and alignment.
func RowWithLayout(id string, distribution Distribution, alignment Alignment, children ...string) Component {
	return Component{
		ID:           id,
		Component:    "Row",
		Children:     children,
		Distribution: distribution,
		Alignment:    alignment,
	}
}
