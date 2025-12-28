// Package a2ui provides types and utilities for generating A2UI messages.
// A2UI is a protocol for AI agents to generate declarative UI definitions
// that clients render with native widgets.
package a2ui

// Message represents an A2UI protocol message.
// Exactly one field should be set per message.
type Message struct {
	BeginRendering   *BeginRendering   `json:"beginRendering,omitempty"`
	UpdateComponents *UpdateComponents `json:"updateComponents,omitempty"`
	DataModelUpdate  *DataModelUpdate  `json:"dataModelUpdate,omitempty"`
	DeleteSurface    *DeleteSurface    `json:"deleteSurface,omitempty"`
}

// BeginRendering initializes a new UI surface.
type BeginRendering struct {
	SurfaceID string `json:"surfaceId"`
	Root      string `json:"root"`
}

// UpdateComponents sends component definitions to the client.
type UpdateComponents struct {
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
// Uses a flat structure with "component" field indicating the type.
type Component struct {
	ID        string `json:"id"`
	Component string `json:"component"`

	// Layout properties (Column, Row)
	Children     []string     `json:"children,omitempty"`
	Distribution Distribution `json:"distribution,omitempty"`
	Alignment    Alignment    `json:"alignment,omitempty"`

	// Card property
	Child string `json:"child,omitempty"`

	// List properties
	Template    string       `json:"template,omitempty"`
	DataBinding *DataBinding `json:"dataBinding,omitempty"`
	Direction   string       `json:"direction,omitempty"` // "vertical" or "horizontal"

	// Tabs property
	Tabs []TabDef `json:"tabs,omitempty"`

	// Modal properties
	EntryPointChild string `json:"entryPointChild,omitempty"`
	ContentChild    string `json:"contentChild,omitempty"`

	// Text/Image properties
	Text      string    `json:"text,omitempty"`
	URL       string    `json:"url,omitempty"`
	Alt       string    `json:"alt,omitempty"`
	Fit       ImageFit  `json:"fit,omitempty"`
	UsageHint UsageHint `json:"usageHint,omitempty"`

	// Icon property
	Icon IconName `json:"icon,omitempty"`

	// AudioPlayer property
	Description string `json:"description,omitempty"`

	// Divider property
	Orientation string `json:"orientation,omitempty"`

	// Button properties
	Action  *Action `json:"action,omitempty"`
	Primary bool    `json:"primary,omitempty"`

	// TextField properties
	Label            string        `json:"label,omitempty"`
	Placeholder      string        `json:"placeholder,omitempty"`
	TextFieldType    TextFieldType `json:"textFieldType,omitempty"`
	ValidationRegexp string        `json:"validationRegexp,omitempty"`

	// CheckBox property
	Checked bool `json:"checked,omitempty"`

	// DateTimeInput properties
	EnableDate bool `json:"enableDate,omitempty"`
	EnableTime bool `json:"enableTime,omitempty"`

	// MultipleChoice properties
	Options              []ChoiceOption `json:"options,omitempty"`
	Selections           []string       `json:"selections,omitempty"`
	MaxAllowedSelections int            `json:"maxAllowedSelections,omitempty"`

	// Slider properties
	MinValue    float64 `json:"minValue,omitempty"`
	MaxValue    float64 `json:"maxValue,omitempty"`
	SliderValue float64 `json:"value,omitempty"`
}

// Distribution defines how children are distributed along the main axis.
type Distribution string

const (
	DistributionStart        Distribution = "start"
	DistributionCenter       Distribution = "center"
	DistributionEnd          Distribution = "end"
	DistributionSpaceAround  Distribution = "spaceAround"
	DistributionSpaceBetween Distribution = "spaceBetween"
	DistributionSpaceEvenly  Distribution = "spaceEvenly"
)

// Alignment defines how children are aligned along the cross axis.
type Alignment string

const (
	AlignmentStart   Alignment = "start"
	AlignmentCenter  Alignment = "center"
	AlignmentEnd     Alignment = "end"
	AlignmentStretch Alignment = "stretch"
)

// TabDef represents a single tab in a Tabs component.
type TabDef struct {
	Title string `json:"title"`
	Child string `json:"child"`
}

// UsageHint provides styling hints for text and images.
type UsageHint string

const (
	UsageHintH1      UsageHint = "h1"
	UsageHintH2      UsageHint = "h2"
	UsageHintH3      UsageHint = "h3"
	UsageHintH4      UsageHint = "h4"
	UsageHintH5      UsageHint = "h5"
	UsageHintBody    UsageHint = "body"
	UsageHintCaption UsageHint = "caption"
)

// ImageFit defines how an image fits within its container.
type ImageFit string

const (
	ImageFitContain   ImageFit = "contain"
	ImageFitCover     ImageFit = "cover"
	ImageFitFill      ImageFit = "fill"
	ImageFitNone      ImageFit = "none"
	ImageFitScaleDown ImageFit = "scale-down"
)

// IconName represents a predefined icon.
type IconName string

const (
	IconAccountCircle IconName = "accountCircle"
	IconAdd           IconName = "add"
	IconArrowBack     IconName = "arrowBack"
	IconCheck         IconName = "check"
	IconClose         IconName = "close"
	IconDelete        IconName = "delete"
	IconEdit          IconName = "edit"
	IconFavorite      IconName = "favorite"
	IconHome          IconName = "home"
	IconMenu          IconName = "menu"
	IconSearch        IconName = "search"
	IconSettings      IconName = "settings"
	IconStar          IconName = "star"
	IconWarning       IconName = "warning"
)

// Action defines what happens when a button is clicked.
type Action struct {
	Type string         `json:"type"`
	Data map[string]any `json:"data,omitempty"`
}

// TextFieldType defines the type of text input.
type TextFieldType string

const (
	TextFieldTypeShortText TextFieldType = "shortText"
	TextFieldTypeLongText  TextFieldType = "longText"
	TextFieldTypeNumber    TextFieldType = "number"
	TextFieldTypeDate      TextFieldType = "date"
	TextFieldTypeObscured  TextFieldType = "obscured"
)

// ChoiceOption represents an option in a multiple choice component.
type ChoiceOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// ClientMessage represents a message sent from client to server.
// Used for handling user interactions.
type ClientMessage struct {
	Event *Event `json:"event,omitempty"`
}

// Event represents a user interaction with a component.
type Event struct {
	SurfaceID   string         `json:"surfaceId"`
	ComponentID string         `json:"componentId"`
	Type        string         `json:"type"` // "action", "input", "change"
	Data        map[string]any `json:"data,omitempty"`
}

// DataBinding binds a component to a JSON Pointer path in the data model.
type DataBinding struct {
	Path string `json:"path"`
}
