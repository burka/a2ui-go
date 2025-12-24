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
	ID string `json:"id"`
	// Layout components
	Column *ColumnDef `json:"Column,omitempty"`
	Row    *RowDef    `json:"Row,omitempty"`
	Card   *CardDef   `json:"Card,omitempty"`
	List   *ListDef   `json:"List,omitempty"`
	Tabs   *TabsDef   `json:"Tabs,omitempty"`
	Modal  *ModalDef  `json:"Modal,omitempty"`
	// Content components
	Text        *TextDef        `json:"Text,omitempty"`
	Image       *ImageDef       `json:"Image,omitempty"`
	Icon        *IconDef        `json:"Icon,omitempty"`
	Video       *VideoDef       `json:"Video,omitempty"`
	AudioPlayer *AudioPlayerDef `json:"AudioPlayer,omitempty"`
	Divider     *DividerDef     `json:"Divider,omitempty"`
	// Form components
	Button         *ButtonDef         `json:"Button,omitempty"`
	TextField      *TextFieldDef      `json:"TextField,omitempty"`
	CheckBox       *CheckBoxDef       `json:"CheckBox,omitempty"`
	DateTimeInput  *DateTimeInputDef  `json:"DateTimeInput,omitempty"`
	MultipleChoice *MultipleChoiceDef `json:"MultipleChoice,omitempty"`
	Slider         *SliderDef         `json:"Slider,omitempty"`
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

// ColumnDef is a vertical layout container.
type ColumnDef struct {
	Children     []string     `json:"children"`
	Distribution Distribution `json:"distribution,omitempty"`
	Alignment    Alignment    `json:"alignment,omitempty"`
}

// RowDef is a horizontal layout container.
type RowDef struct {
	Children     []string     `json:"children"`
	Distribution Distribution `json:"distribution,omitempty"`
	Alignment    Alignment    `json:"alignment,omitempty"`
}

// CardDef is a container with visual styling.
type CardDef struct {
	Child string `json:"child"`
}

// ListDef renders items from a data array using a template.
type ListDef struct {
	Template    string       `json:"template"`
	DataBinding *DataBinding `json:"dataBinding,omitempty"`
	Direction   string       `json:"direction,omitempty"` // "vertical" or "horizontal"
	Alignment   Alignment    `json:"alignment,omitempty"`
}

// TabDef represents a single tab in a Tabs component.
type TabDef struct {
	Title string `json:"title"`
	Child string `json:"child"`
}

// TabsDef displays content in switchable tabs.
type TabsDef struct {
	Tabs []TabDef `json:"tabs"`
}

// ModalDef displays content in a modal overlay.
type ModalDef struct {
	EntryPointChild string `json:"entryPointChild"`
	ContentChild    string `json:"contentChild"`
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

// TextDef displays text content.
type TextDef struct {
	Text        string       `json:"text,omitempty"`
	DataBinding *DataBinding `json:"dataBinding,omitempty"`
	UsageHint   UsageHint    `json:"usageHint,omitempty"`
}

// ImageFit defines how an image fits within its container.
type ImageFit string

const (
	ImageFitContain   ImageFit = "contain"
	ImageFitCover     ImageFit = "cover"
	ImageFitFill      ImageFit = "fill"
	ImageFitNone      ImageFit = "none"
	ImageFitScaleDown ImageFit = "scale-down"
)

// ImageDef displays an image.
type ImageDef struct {
	URL         string       `json:"url,omitempty"`
	Alt         string       `json:"alt,omitempty"`
	DataBinding *DataBinding `json:"dataBinding,omitempty"`
	Fit         ImageFit     `json:"fit,omitempty"`
	UsageHint   UsageHint    `json:"usageHint,omitempty"`
}

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

// IconDef displays a predefined icon.
type IconDef struct {
	Icon IconName `json:"icon"`
}

// VideoDef displays a video player.
type VideoDef struct {
	URL         string       `json:"url,omitempty"`
	DataBinding *DataBinding `json:"dataBinding,omitempty"`
}

// AudioPlayerDef displays an audio player.
type AudioPlayerDef struct {
	URL         string       `json:"url,omitempty"`
	DataBinding *DataBinding `json:"dataBinding,omitempty"`
	Description string       `json:"description,omitempty"`
}

// DividerDef displays a visual separator.
type DividerDef struct {
	Orientation string `json:"orientation,omitempty"` // "horizontal" or "vertical"
}

// ButtonDef triggers an action when clicked.
type ButtonDef struct {
	Text    string `json:"text"`
	Action  Action `json:"action"`
	Primary bool   `json:"primary,omitempty"`
}

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

// TextFieldDef accepts user text input.
type TextFieldDef struct {
	Label            string        `json:"label,omitempty"`
	Placeholder      string        `json:"placeholder,omitempty"`
	DataBinding      *DataBinding  `json:"dataBinding,omitempty"`
	TextFieldType    TextFieldType `json:"textFieldType,omitempty"`
	ValidationRegexp string        `json:"validationRegexp,omitempty"`
}

// CheckBoxDef displays a checkbox input.
type CheckBoxDef struct {
	Label       string       `json:"label,omitempty"`
	Value       bool         `json:"value,omitempty"`
	DataBinding *DataBinding `json:"dataBinding,omitempty"`
}

// DateTimeInputDef displays a date/time picker.
type DateTimeInputDef struct {
	Label       string       `json:"label,omitempty"`
	DataBinding *DataBinding `json:"dataBinding,omitempty"`
	EnableDate  bool         `json:"enableDate,omitempty"`
	EnableTime  bool         `json:"enableTime,omitempty"`
}

// ChoiceOption represents an option in a multiple choice component.
type ChoiceOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// MultipleChoiceDef displays a multiple choice selector.
type MultipleChoiceDef struct {
	Label                string         `json:"label,omitempty"`
	Options              []ChoiceOption `json:"options"`
	Selections           []string       `json:"selections,omitempty"`
	DataBinding          *DataBinding   `json:"dataBinding,omitempty"`
	MaxAllowedSelections int            `json:"maxAllowedSelections,omitempty"`
}

// SliderDef displays a numeric slider input.
type SliderDef struct {
	Label       string       `json:"label,omitempty"`
	MinValue    float64      `json:"minValue"`
	MaxValue    float64      `json:"maxValue"`
	Value       float64      `json:"value,omitempty"`
	DataBinding *DataBinding `json:"dataBinding,omitempty"`
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
