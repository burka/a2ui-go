package a2ui

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewSurface(t *testing.T) {
	s := NewSurface("test-surface")
	if s == nil {
		t.Fatal("NewSurface returned nil")
	}
	if s.id != "test-surface" {
		t.Errorf("expected id 'test-surface', got '%s'", s.id)
	}
	if s.root != "root" {
		t.Errorf("expected root 'root', got '%s'", s.root)
	}
}

func TestSurfaceAdd(t *testing.T) {
	s := NewSurface("test")
	s.Add(Column("root", "child1", "child2"))
	s.Add(TextStatic("child1", "Hello"))
	s.Add(TextStatic("child2", "World"))

	if len(s.components) != 3 {
		t.Errorf("expected 3 components, got %d", len(s.components))
	}
}

func TestSurfaceSetData(t *testing.T) {
	s := NewSurface("test")
	s.SetData("/user/name", "Alice")
	s.SetData("/items", []string{"a", "b", "c"})

	if s.data["/user/name"] != "Alice" {
		t.Errorf("expected 'Alice', got '%v'", s.data["/user/name"])
	}

	items, ok := s.data["/items"].([]string)
	if !ok || len(items) != 3 {
		t.Errorf("expected 3 items, got '%v'", s.data["/items"])
	}
}

func TestSurfaceMessages(t *testing.T) {
	s := NewSurface("test")
	s.Add(Column("root", "text"))
	s.Add(TextStatic("text", "Hello"))
	s.SetData("/value", 42)

	msgs := s.Messages()

	if len(msgs) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(msgs))
	}

	// Check beginRendering
	if msgs[0].BeginRendering == nil {
		t.Error("expected BeginRendering message")
	}
	if msgs[0].BeginRendering.SurfaceID != "test" {
		t.Errorf("expected surface ID 'test', got '%s'", msgs[0].BeginRendering.SurfaceID)
	}
	if msgs[0].BeginRendering.Root != "root" {
		t.Errorf("expected root 'root', got '%s'", msgs[0].BeginRendering.Root)
	}

	// Check updateComponents
	if msgs[1].UpdateComponents == nil {
		t.Error("expected UpdateComponents message")
	}
	if len(msgs[1].UpdateComponents.Components) != 2 {
		t.Errorf("expected 2 components, got %d", len(msgs[1].UpdateComponents.Components))
	}

	// Check dataModelUpdate
	if msgs[2].DataModelUpdate == nil {
		t.Error("expected DataModelUpdate message")
	}
	if msgs[2].DataModelUpdate.Contents["/value"] != 42 {
		t.Errorf("expected value 42, got '%v'", msgs[2].DataModelUpdate.Contents["/value"])
	}
}

func TestSurfaceMessagesNoData(t *testing.T) {
	s := NewSurface("test")
	s.Add(TextStatic("root", "Hello"))

	msgs := s.Messages()

	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages when no data, got %d", len(msgs))
	}
}

func TestColumnHelper(t *testing.T) {
	c := Column("col", "a", "b", "c")

	if c.ID != "col" {
		t.Errorf("expected ID 'col', got '%s'", c.ID)
	}
	if c.Component != "Column" {
		t.Errorf("expected Component 'Column', got '%s'", c.Component)
	}
	if len(c.Children) != 3 {
		t.Errorf("expected 3 children, got %d", len(c.Children))
	}
}

func TestRowHelper(t *testing.T) {
	r := Row("row", "x", "y")

	if r.ID != "row" {
		t.Errorf("expected ID 'row', got '%s'", r.ID)
	}
	if r.Component != "Row" {
		t.Errorf("expected Component 'Row', got '%s'", r.Component)
	}
	if len(r.Children) != 2 {
		t.Errorf("expected 2 children, got %d", len(r.Children))
	}
}

func TestCardHelper(t *testing.T) {
	c := Card("card", "content")

	if c.ID != "card" {
		t.Errorf("expected ID 'card', got '%s'", c.ID)
	}
	if c.Component != "Card" {
		t.Errorf("expected Component 'Card', got '%s'", c.Component)
	}
	if c.Child != "content" {
		t.Errorf("expected child 'content', got '%s'", c.Child)
	}
}

func TestTextStaticHelper(t *testing.T) {
	txt := TextStatic("txt", "Hello World")

	if txt.ID != "txt" {
		t.Errorf("expected ID 'txt', got '%s'", txt.ID)
	}
	if txt.Component != "Text" {
		t.Errorf("expected Component 'Text', got '%s'", txt.Component)
	}
	if txt.Text != "Hello World" {
		t.Errorf("expected text 'Hello World', got '%s'", txt.Text)
	}
	if txt.DataBinding != nil {
		t.Error("expected no DataBinding for static text")
	}
}

func TestTextBoundHelper(t *testing.T) {
	txt := TextBound("txt", "/user/name")

	if txt.Component != "Text" {
		t.Errorf("expected Component 'Text', got '%s'", txt.Component)
	}
	if txt.DataBinding == nil {
		t.Fatal("expected DataBinding to be set")
	}
	if txt.DataBinding.Path != "/user/name" {
		t.Errorf("expected path '/user/name', got '%s'", txt.DataBinding.Path)
	}
}

func TestListTemplateHelper(t *testing.T) {
	lst := ListTemplate("list", "item-template", "/items")

	if lst.ID != "list" {
		t.Errorf("expected ID 'list', got '%s'", lst.ID)
	}
	if lst.Component != "List" {
		t.Errorf("expected Component 'List', got '%s'", lst.Component)
	}
	if lst.Template != "item-template" {
		t.Errorf("expected template 'item-template', got '%s'", lst.Template)
	}
	if lst.DataBinding == nil {
		t.Fatal("expected DataBinding to be set")
	}
	if lst.DataBinding.Path != "/items" {
		t.Errorf("expected path '/items', got '%s'", lst.DataBinding.Path)
	}
}

func TestWriteJSONL(t *testing.T) {
	s := NewSurface("test")
	s.Add(TextStatic("root", "Hello"))

	var buf bytes.Buffer
	err := WriteJSONL(&buf, s.Messages())
	if err != nil {
		t.Fatalf("WriteJSONL failed: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}

	// Verify each line is valid JSON
	for i, line := range lines {
		var msg Message
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			t.Errorf("line %d is not valid JSON: %v", i, err)
		}
	}
}

func TestWriteMessage(t *testing.T) {
	msg := Message{
		BeginRendering: &BeginRendering{SurfaceID: "test", Root: "root"},
	}

	var buf bytes.Buffer
	err := WriteMessage(&buf, msg)
	if err != nil {
		t.Fatalf("WriteMessage failed: %v", err)
	}

	var decoded Message
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if decoded.BeginRendering == nil {
		t.Error("expected BeginRendering to be set")
	}
}

func TestWritePretty(t *testing.T) {
	s := NewSurface("test")
	s.Add(TextStatic("root", "Hello"))

	var buf bytes.Buffer
	err := WritePretty(&buf, s.Messages())
	if err != nil {
		t.Fatalf("WritePretty failed: %v", err)
	}

	output := buf.String()
	// Pretty output should contain indentation
	if !strings.Contains(output, "  ") {
		t.Error("expected indented output")
	}
}

func TestJSONSerialization(t *testing.T) {
	s := NewSurface("demo")
	s.Add(Column("root", "title", "list"))
	s.Add(TextStatic("title", "Products"))
	s.Add(ListTemplate("list", "item", "/products"))
	s.Add(TextBound("item", "/name"))
	s.SetData("/products", []map[string]string{
		{"name": "Widget"},
		{"name": "Gadget"},
	})

	msgs := s.Messages()

	// Serialize to JSON
	for _, msg := range msgs {
		data, err := json.Marshal(msg)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}

		// Verify it can be deserialized
		var decoded Message
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}
	}
}

func TestButtonHelper(t *testing.T) {
	btns := Button("btn", "Click Me", "submit")

	if len(btns) != 2 {
		t.Fatalf("expected 2 components, got %d", len(btns))
	}

	btn := btns[0]
	if btn.ID != "btn" {
		t.Errorf("expected ID 'btn', got '%s'", btn.ID)
	}
	if btn.Component != "Button" {
		t.Errorf("expected Component 'Button', got '%s'", btn.Component)
	}
	if btn.Child != "btn_text" {
		t.Errorf("expected child 'btn_text', got '%s'", btn.Child)
	}
	if btn.Action == nil {
		t.Fatal("expected Action to be set")
	}
	if btn.Action.Type != "submit" {
		t.Errorf("expected action type 'submit', got '%s'", btn.Action.Type)
	}

	// Check the text child
	txt := btns[1]
	if txt.ID != "btn_text" {
		t.Errorf("expected ID 'btn_text', got '%s'", txt.ID)
	}
	if txt.Component != "Text" {
		t.Errorf("expected Component 'Text', got '%s'", txt.Component)
	}
	if txt.Text != "Click Me" {
		t.Errorf("expected text 'Click Me', got '%s'", txt.Text)
	}
}

func TestTextFieldHelper(t *testing.T) {
	tf := TextField("input", "Email", "Enter email")

	if tf.ID != "input" {
		t.Errorf("expected ID 'input', got '%s'", tf.ID)
	}
	if tf.Component != "TextField" {
		t.Errorf("expected Component 'TextField', got '%s'", tf.Component)
	}
	if tf.Label != "Email" {
		t.Errorf("expected label 'Email', got '%s'", tf.Label)
	}
	if tf.Placeholder != "Enter email" {
		t.Errorf("expected placeholder 'Enter email', got '%s'", tf.Placeholder)
	}
}

func TestImageHelpers(t *testing.T) {
	img := ImageStatic("img", "https://example.com/photo.jpg", "A photo")

	if img.ID != "img" {
		t.Errorf("expected ID 'img', got '%s'", img.ID)
	}
	if img.Component != "Image" {
		t.Errorf("expected Component 'Image', got '%s'", img.Component)
	}
	if img.URL != "https://example.com/photo.jpg" {
		t.Errorf("expected URL 'https://example.com/photo.jpg', got '%s'", img.URL)
	}
	if img.Alt != "A photo" {
		t.Errorf("expected alt 'A photo', got '%s'", img.Alt)
	}

	// Test bound image
	imgBound := ImageBound("img2", "/photo/url", "Dynamic photo")
	if imgBound.DataBinding == nil {
		t.Fatal("expected DataBinding to be set")
	}
	if imgBound.DataBinding.Path != "/photo/url" {
		t.Errorf("expected path '/photo/url', got '%s'", imgBound.DataBinding.Path)
	}
}

func TestSetRoot(t *testing.T) {
	s := NewSurface("test")
	s.SetRoot("main")

	if s.root != "main" {
		t.Errorf("expected root 'main', got '%s'", s.root)
	}

	msgs := s.Messages()
	if msgs[0].BeginRendering.Root != "main" {
		t.Errorf("expected BeginRendering root 'main', got '%s'", msgs[0].BeginRendering.Root)
	}
}

func TestAddAll(t *testing.T) {
	s := NewSurface("test")
	s.AddAll(
		Column("root", "a", "b"),
		TextStatic("a", "Hello"),
		TextStatic("b", "World"),
	)

	if len(s.components) != 3 {
		t.Errorf("expected 3 components, got %d", len(s.components))
	}
}

func TestButtonWithData(t *testing.T) {
	btns := ButtonWithData("btn", "Submit", "submit", map[string]any{
		"endpoint": "/api/submit",
		"method":   "POST",
	})

	if len(btns) != 2 {
		t.Fatalf("expected 2 components, got %d", len(btns))
	}

	btn := btns[0]
	if btn.ID != "btn" {
		t.Errorf("expected ID 'btn', got '%s'", btn.ID)
	}
	if btn.Component != "Button" {
		t.Errorf("expected Component 'Button', got '%s'", btn.Component)
	}
	if btn.Action == nil {
		t.Fatal("expected Action to be set")
	}
	if btn.Action.Type != "submit" {
		t.Errorf("expected action type 'submit', got '%s'", btn.Action.Type)
	}
	if btn.Action.Data == nil {
		t.Fatal("expected Action.Data to be set")
	}
	if btn.Action.Data["endpoint"] != "/api/submit" {
		t.Errorf("expected endpoint '/api/submit', got '%v'", btn.Action.Data["endpoint"])
	}
}

func TestClientMessage(t *testing.T) {
	msg := ClientMessage{
		Event: &Event{
			SurfaceID:   "form",
			ComponentID: "submit-btn",
			Type:        "action",
			Data: map[string]any{
				"name": "John",
				"age":  30,
			},
		},
	}

	// Serialize and deserialize
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded ClientMessage
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Event == nil {
		t.Fatal("expected Event to be set")
	}
	if decoded.Event.SurfaceID != "form" {
		t.Errorf("expected surfaceId 'form', got '%s'", decoded.Event.SurfaceID)
	}
	if decoded.Event.ComponentID != "submit-btn" {
		t.Errorf("expected componentId 'submit-btn', got '%s'", decoded.Event.ComponentID)
	}
	if decoded.Event.Type != "action" {
		t.Errorf("expected type 'action', got '%s'", decoded.Event.Type)
	}
	if decoded.Event.Data["name"] != "John" {
		t.Errorf("expected name 'John', got '%v'", decoded.Event.Data["name"])
	}
}

func TestActionWithData(t *testing.T) {
	action := Action{
		Type: "navigate",
		Data: map[string]any{
			"url":    "/dashboard",
			"params": map[string]any{"tab": "settings"},
		},
	}

	data, err := json.Marshal(action)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded Action
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Type != "navigate" {
		t.Errorf("expected type 'navigate', got '%s'", decoded.Type)
	}
	if decoded.Data["url"] != "/dashboard" {
		t.Errorf("expected url '/dashboard', got '%v'", decoded.Data["url"])
	}
}

func TestTabsHelper(t *testing.T) {
	tabs := Tabs("tabs",
		Tab("Tab 1", "content1"),
		Tab("Tab 2", "content2"),
	)

	if tabs.ID != "tabs" {
		t.Errorf("expected ID 'tabs', got '%s'", tabs.ID)
	}
	if tabs.Component != "Tabs" {
		t.Errorf("expected Component 'Tabs', got '%s'", tabs.Component)
	}
	if len(tabs.Tabs) != 2 {
		t.Errorf("expected 2 tabs, got %d", len(tabs.Tabs))
	}
	if tabs.Tabs[0].Title != "Tab 1" {
		t.Errorf("expected title 'Tab 1', got '%s'", tabs.Tabs[0].Title)
	}
	if tabs.Tabs[0].Child != "content1" {
		t.Errorf("expected child 'content1', got '%s'", tabs.Tabs[0].Child)
	}
}

func TestModalHelper(t *testing.T) {
	modal := Modal("modal", "trigger-btn", "dialog-content")

	if modal.ID != "modal" {
		t.Errorf("expected ID 'modal', got '%s'", modal.ID)
	}
	if modal.Component != "Modal" {
		t.Errorf("expected Component 'Modal', got '%s'", modal.Component)
	}
	if modal.EntryPointChild != "trigger-btn" {
		t.Errorf("expected entryPointChild 'trigger-btn', got '%s'", modal.EntryPointChild)
	}
	if modal.ContentChild != "dialog-content" {
		t.Errorf("expected contentChild 'dialog-content', got '%s'", modal.ContentChild)
	}
}

func TestIconHelper(t *testing.T) {
	icon := Icon("icon", IconSearch)

	if icon.ID != "icon" {
		t.Errorf("expected ID 'icon', got '%s'", icon.ID)
	}
	if icon.Component != "Icon" {
		t.Errorf("expected Component 'Icon', got '%s'", icon.Component)
	}
	if icon.Icon != IconSearch {
		t.Errorf("expected icon 'search', got '%s'", icon.Icon)
	}
}

func TestVideoHelpers(t *testing.T) {
	video := Video("video", "https://example.com/video.mp4")

	if video.ID != "video" {
		t.Errorf("expected ID 'video', got '%s'", video.ID)
	}
	if video.Component != "Video" {
		t.Errorf("expected Component 'Video', got '%s'", video.Component)
	}
	if video.URL != "https://example.com/video.mp4" {
		t.Errorf("expected URL 'https://example.com/video.mp4', got '%s'", video.URL)
	}

	// Test bound video
	videoBound := VideoBound("video2", "/media/video")
	if videoBound.DataBinding == nil {
		t.Fatal("expected DataBinding to be set")
	}
	if videoBound.DataBinding.Path != "/media/video" {
		t.Errorf("expected path '/media/video', got '%s'", videoBound.DataBinding.Path)
	}
}

func TestAudioPlayerHelpers(t *testing.T) {
	audio := AudioPlayer("audio", "https://example.com/song.mp3", "My Song")

	if audio.ID != "audio" {
		t.Errorf("expected ID 'audio', got '%s'", audio.ID)
	}
	if audio.Component != "AudioPlayer" {
		t.Errorf("expected Component 'AudioPlayer', got '%s'", audio.Component)
	}
	if audio.URL != "https://example.com/song.mp3" {
		t.Errorf("expected URL 'https://example.com/song.mp3', got '%s'", audio.URL)
	}
	if audio.Description != "My Song" {
		t.Errorf("expected description 'My Song', got '%s'", audio.Description)
	}

	// Test bound audio
	audioBound := AudioPlayerBound("audio2", "/media/audio", "Dynamic Audio")
	if audioBound.DataBinding == nil {
		t.Fatal("expected DataBinding to be set")
	}
	if audioBound.DataBinding.Path != "/media/audio" {
		t.Errorf("expected path '/media/audio', got '%s'", audioBound.DataBinding.Path)
	}
}

func TestDividerHelpers(t *testing.T) {
	div := Divider("div")

	if div.ID != "div" {
		t.Errorf("expected ID 'div', got '%s'", div.ID)
	}
	if div.Component != "Divider" {
		t.Errorf("expected Component 'Divider', got '%s'", div.Component)
	}

	// Test vertical divider
	divV := DividerVertical("divV")
	if divV.Orientation != "vertical" {
		t.Errorf("expected orientation 'vertical', got '%s'", divV.Orientation)
	}
}

func TestCheckBoxHelpers(t *testing.T) {
	cb := CheckBox("cb", "Accept Terms", true)

	if cb.ID != "cb" {
		t.Errorf("expected ID 'cb', got '%s'", cb.ID)
	}
	if cb.Component != "CheckBox" {
		t.Errorf("expected Component 'CheckBox', got '%s'", cb.Component)
	}
	if cb.Label != "Accept Terms" {
		t.Errorf("expected label 'Accept Terms', got '%s'", cb.Label)
	}
	if cb.Checked != true {
		t.Errorf("expected checked true, got false")
	}

	// Test bound checkbox
	cbBound := CheckBoxBound("cb2", "Subscribe", "/prefs/subscribe")
	if cbBound.DataBinding == nil {
		t.Fatal("expected DataBinding to be set")
	}
	if cbBound.DataBinding.Path != "/prefs/subscribe" {
		t.Errorf("expected path '/prefs/subscribe', got '%s'", cbBound.DataBinding.Path)
	}
}

func TestDateTimeInputHelpers(t *testing.T) {
	dt := DateTimeInput("dt", "Select Date", true, false)

	if dt.ID != "dt" {
		t.Errorf("expected ID 'dt', got '%s'", dt.ID)
	}
	if dt.Component != "DateTimeInput" {
		t.Errorf("expected Component 'DateTimeInput', got '%s'", dt.Component)
	}
	if dt.Label != "Select Date" {
		t.Errorf("expected label 'Select Date', got '%s'", dt.Label)
	}
	if dt.EnableDate != true {
		t.Error("expected enableDate true")
	}
	if dt.EnableTime != false {
		t.Error("expected enableTime false")
	}

	// Test bound datetime
	dtBound := DateTimeInputBound("dt2", "Appointment", "/booking/datetime", true, true)
	if dtBound.DataBinding == nil {
		t.Fatal("expected DataBinding to be set")
	}
	if dtBound.DataBinding.Path != "/booking/datetime" {
		t.Errorf("expected path '/booking/datetime', got '%s'", dtBound.DataBinding.Path)
	}
}

func TestMultipleChoiceHelpers(t *testing.T) {
	mc := MultipleChoice("mc", "Select Size", []ChoiceOption{
		Choice("Small", "s"),
		Choice("Medium", "m"),
		Choice("Large", "l"),
	})

	if mc.ID != "mc" {
		t.Errorf("expected ID 'mc', got '%s'", mc.ID)
	}
	if mc.Component != "MultipleChoice" {
		t.Errorf("expected Component 'MultipleChoice', got '%s'", mc.Component)
	}
	if mc.Label != "Select Size" {
		t.Errorf("expected label 'Select Size', got '%s'", mc.Label)
	}
	if len(mc.Options) != 3 {
		t.Errorf("expected 3 options, got %d", len(mc.Options))
	}
	if mc.Options[0].Label != "Small" {
		t.Errorf("expected first option label 'Small', got '%s'", mc.Options[0].Label)
	}
	if mc.Options[0].Value != "s" {
		t.Errorf("expected first option value 's', got '%s'", mc.Options[0].Value)
	}

	// Test bound multiple choice
	mcBound := MultipleChoiceBound("mc2", "Color", "/prefs/color", []ChoiceOption{
		Choice("Red", "red"),
		Choice("Blue", "blue"),
	})
	if mcBound.DataBinding == nil {
		t.Fatal("expected DataBinding to be set")
	}
	if mcBound.DataBinding.Path != "/prefs/color" {
		t.Errorf("expected path '/prefs/color', got '%s'", mcBound.DataBinding.Path)
	}
}

func TestSliderHelpers(t *testing.T) {
	slider := Slider("slider", "Volume", 0, 100, 50)

	if slider.ID != "slider" {
		t.Errorf("expected ID 'slider', got '%s'", slider.ID)
	}
	if slider.Component != "Slider" {
		t.Errorf("expected Component 'Slider', got '%s'", slider.Component)
	}
	if slider.Label != "Volume" {
		t.Errorf("expected label 'Volume', got '%s'", slider.Label)
	}
	if slider.MinValue != 0 {
		t.Errorf("expected minValue 0, got %f", slider.MinValue)
	}
	if slider.MaxValue != 100 {
		t.Errorf("expected maxValue 100, got %f", slider.MaxValue)
	}
	if slider.SliderValue != 50 {
		t.Errorf("expected value 50, got %f", slider.SliderValue)
	}

	// Test bound slider
	sliderBound := SliderBound("slider2", "Brightness", "/settings/brightness", 0, 255)
	if sliderBound.DataBinding == nil {
		t.Fatal("expected DataBinding to be set")
	}
	if sliderBound.DataBinding.Path != "/settings/brightness" {
		t.Errorf("expected path '/settings/brightness', got '%s'", sliderBound.DataBinding.Path)
	}
}

func TestButtonPrimaryHelper(t *testing.T) {
	btns := ButtonPrimary("btn", "Save", "save")

	if len(btns) != 2 {
		t.Fatalf("expected 2 components, got %d", len(btns))
	}

	btn := btns[0]
	if btn.Component != "Button" {
		t.Errorf("expected Component 'Button', got '%s'", btn.Component)
	}
	if btn.Primary != true {
		t.Error("expected primary true")
	}
}

func TestTextWithHintHelper(t *testing.T) {
	txt := TextWithHint("title", "Welcome", UsageHintH1)

	if txt.Component != "Text" {
		t.Errorf("expected Component 'Text', got '%s'", txt.Component)
	}
	if txt.UsageHint != UsageHintH1 {
		t.Errorf("expected usageHint 'h1', got '%s'", txt.UsageHint)
	}
}

func TestImageWithFitHelper(t *testing.T) {
	img := ImageWithFit("img", "https://example.com/photo.jpg", "Photo", ImageFitCover)

	if img.Component != "Image" {
		t.Errorf("expected Component 'Image', got '%s'", img.Component)
	}
	if img.Fit != ImageFitCover {
		t.Errorf("expected fit 'cover', got '%s'", img.Fit)
	}
}

func TestTextFieldWithTypeHelper(t *testing.T) {
	tf := TextFieldWithType("password", "Password", "Enter password", TextFieldTypeObscured)

	if tf.Component != "TextField" {
		t.Errorf("expected Component 'TextField', got '%s'", tf.Component)
	}
	if tf.TextFieldType != TextFieldTypeObscured {
		t.Errorf("expected textFieldType 'obscured', got '%s'", tf.TextFieldType)
	}
}

func TestColumnWithLayoutHelper(t *testing.T) {
	col := ColumnWithLayout("col", DistributionSpaceBetween, AlignmentCenter, "a", "b")

	if col.Component != "Column" {
		t.Errorf("expected Component 'Column', got '%s'", col.Component)
	}
	if col.Distribution != DistributionSpaceBetween {
		t.Errorf("expected distribution 'spaceBetween', got '%s'", col.Distribution)
	}
	if col.Alignment != AlignmentCenter {
		t.Errorf("expected alignment 'center', got '%s'", col.Alignment)
	}
}

func TestRowWithLayoutHelper(t *testing.T) {
	row := RowWithLayout("row", DistributionSpaceEvenly, AlignmentStretch, "x", "y")

	if row.Component != "Row" {
		t.Errorf("expected Component 'Row', got '%s'", row.Component)
	}
	if row.Distribution != DistributionSpaceEvenly {
		t.Errorf("expected distribution 'spaceEvenly', got '%s'", row.Distribution)
	}
	if row.Alignment != AlignmentStretch {
		t.Errorf("expected alignment 'stretch', got '%s'", row.Alignment)
	}
}

func TestTextFieldBoundHelper(t *testing.T) {
	tf := TextFieldBound("email", "Email", "Enter email", "/user/email")

	if tf.Component != "TextField" {
		t.Errorf("expected Component 'TextField', got '%s'", tf.Component)
	}
	if tf.DataBinding == nil {
		t.Fatal("expected DataBinding to be set")
	}
	if tf.DataBinding.Path != "/user/email" {
		t.Errorf("expected path '/user/email', got '%s'", tf.DataBinding.Path)
	}
}

func TestAllComponentsJSONSerialization(t *testing.T) {
	s := NewSurface("all-components")

	// Add all component types (buttons are added via AddAll since they return []Component)
	s.AddAll(
		Column("root", "row", "tabs", "modal", "list"),
		Row("row", "card", "text", "image", "icon"),
		Card("card", "divider"),
		Tabs("tabs", Tab("Tab1", "content1")),
		Modal("modal", "trigger", "dialog"),
		ListTemplate("list", "template", "/items"),
		TextStatic("text", "Hello"),
		ImageStatic("image", "https://example.com/img.jpg", "Alt"),
		Icon("icon", IconStar),
		Video("video", "https://example.com/video.mp4"),
		AudioPlayer("audio", "https://example.com/audio.mp3", "Audio"),
		Divider("divider"),
		TextField("textfield", "Label", "Placeholder"),
		CheckBox("checkbox", "Check", false),
		DateTimeInput("datetime", "Date", true, true),
		MultipleChoice("choice", "Select", []ChoiceOption{Choice("A", "a")}),
		Slider("slider", "Slide", 0, 100, 50),
	)
	// Add button components
	s.AddAll(Button("button", "Click", "action")...)

	// Add placeholder components for references
	s.AddAll(
		TextStatic("content1", "Tab Content"),
		TextStatic("trigger", "Trigger"),
		TextStatic("dialog", "Dialog"),
		TextStatic("template", "Template"),
	)

	msgs := s.Messages()

	// Verify serialization works for all components
	for _, msg := range msgs {
		data, err := json.Marshal(msg)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}

		var decoded Message
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}
	}

	// Verify we have the expected number of components
	if msgs[1].UpdateComponents == nil {
		t.Fatal("expected UpdateComponents")
	}
	// 17 base components + 2 button components + 4 placeholder components = 23
	if len(msgs[1].UpdateComponents.Components) != 23 {
		t.Errorf("expected 23 components, got %d", len(msgs[1].UpdateComponents.Components))
	}
}

func TestEnumValues(t *testing.T) {
	// Distribution values
	distributions := []Distribution{
		DistributionStart, DistributionCenter, DistributionEnd,
		DistributionSpaceAround, DistributionSpaceBetween, DistributionSpaceEvenly,
	}
	for _, d := range distributions {
		if d == "" {
			t.Error("expected non-empty distribution value")
		}
	}

	// Alignment values
	alignments := []Alignment{
		AlignmentStart, AlignmentCenter, AlignmentEnd, AlignmentStretch,
	}
	for _, a := range alignments {
		if a == "" {
			t.Error("expected non-empty alignment value")
		}
	}

	// UsageHint values
	hints := []UsageHint{
		UsageHintH1, UsageHintH2, UsageHintH3, UsageHintH4, UsageHintH5,
		UsageHintBody, UsageHintCaption,
	}
	for _, h := range hints {
		if h == "" {
			t.Error("expected non-empty usage hint value")
		}
	}

	// ImageFit values
	fits := []ImageFit{
		ImageFitContain, ImageFitCover, ImageFitFill, ImageFitNone, ImageFitScaleDown,
	}
	for _, f := range fits {
		if f == "" {
			t.Error("expected non-empty image fit value")
		}
	}

	// IconName values
	icons := []IconName{
		IconAccountCircle, IconAdd, IconArrowBack, IconCheck, IconClose,
		IconDelete, IconEdit, IconFavorite, IconHome, IconMenu,
		IconSearch, IconSettings, IconStar, IconWarning,
	}
	for _, i := range icons {
		if i == "" {
			t.Error("expected non-empty icon name value")
		}
	}

	// TextFieldType values
	types := []TextFieldType{
		TextFieldTypeShortText, TextFieldTypeLongText, TextFieldTypeNumber,
		TextFieldTypeDate, TextFieldTypeObscured,
	}
	for _, tp := range types {
		if tp == "" {
			t.Error("expected non-empty text field type value")
		}
	}
}

func TestValidateEmptyID(t *testing.T) {
	s := NewSurface("test")
	s.Add(TextStatic("root", "Hello"))
	s.Add(Component{ID: "", Component: "Text", Text: "Empty ID"})

	errors := s.Validate()

	if len(errors) == 0 {
		t.Fatal("expected validation errors for empty ID")
	}

	found := false
	for _, err := range errors {
		if err.Field == "ID" && err.Message == "component ID must not be empty" {
			found = true
			break
		}
	}

	if !found {
		t.Error("expected error about empty component ID")
	}
}

func TestValidateDuplicateID(t *testing.T) {
	s := NewSurface("test")
	s.Add(TextStatic("root", "Hello"))
	s.Add(TextStatic("duplicate", "First"))
	s.Add(TextStatic("duplicate", "Second"))

	errors := s.Validate()

	if len(errors) == 0 {
		t.Fatal("expected validation errors for duplicate IDs")
	}

	found := false
	for _, err := range errors {
		if err.ComponentID == "duplicate" && err.Field == "ID" && err.Message == "duplicate component ID" {
			found = true
			break
		}
	}

	if !found {
		t.Error("expected error about duplicate component ID")
	}
}

func TestValidateRootNotFound(t *testing.T) {
	s := NewSurface("test")
	s.SetRoot("nonexistent")
	s.Add(TextStatic("text", "Hello"))

	errors := s.Validate()

	if len(errors) == 0 {
		t.Fatal("expected validation errors for missing root")
	}

	found := false
	for _, err := range errors {
		if err.ComponentID == "nonexistent" && err.Field == "Root" && err.Message == "root component not found" {
			found = true
			break
		}
	}

	if !found {
		t.Error("expected error about root component not found")
	}
}

func TestValidateChildNotFound(t *testing.T) {
	// Test Column with missing child
	t.Run("Column", func(t *testing.T) {
		s := NewSurface("test")
		s.Add(Column("root", "missing-child"))

		errors := s.Validate()

		if len(errors) == 0 {
			t.Fatal("expected validation errors for missing Column child")
		}

		found := false
		for _, err := range errors {
			if err.ComponentID == "root" && err.Field == "Column.Children" {
				found = true
				break
			}
		}

		if !found {
			t.Error("expected error about missing Column child")
		}
	})

	// Test Row with missing child
	t.Run("Row", func(t *testing.T) {
		s := NewSurface("test")
		s.Add(Row("root", "missing-child"))

		errors := s.Validate()

		if len(errors) == 0 {
			t.Fatal("expected validation errors for missing Row child")
		}

		found := false
		for _, err := range errors {
			if err.ComponentID == "root" && err.Field == "Row.Children" {
				found = true
				break
			}
		}

		if !found {
			t.Error("expected error about missing Row child")
		}
	})

	// Test Card with missing child
	t.Run("Card", func(t *testing.T) {
		s := NewSurface("test")
		s.Add(Card("root", "missing-child"))

		errors := s.Validate()

		if len(errors) == 0 {
			t.Fatal("expected validation errors for missing Card child")
		}

		found := false
		for _, err := range errors {
			if err.ComponentID == "root" && err.Field == "Card.Child" {
				found = true
				break
			}
		}

		if !found {
			t.Error("expected error about missing Card child")
		}
	})

	// Test Button with missing child
	t.Run("Button", func(t *testing.T) {
		s := NewSurface("test")
		s.Add(TextStatic("root", "Root"))
		s.Add(ButtonOnly("btn", "missing-child", "click"))

		errors := s.Validate()

		found := false
		for _, err := range errors {
			if err.ComponentID == "btn" && err.Field == "Button.Child" {
				found = true
				break
			}
		}

		if !found {
			t.Error("expected error about missing Button child")
		}
	})

	// Test List with missing template
	t.Run("List", func(t *testing.T) {
		s := NewSurface("test")
		s.Add(ListTemplate("root", "missing-template", "/items"))

		errors := s.Validate()

		if len(errors) == 0 {
			t.Fatal("expected validation errors for missing List template")
		}

		found := false
		for _, err := range errors {
			if err.ComponentID == "root" && err.Field == "List.Template" {
				found = true
				break
			}
		}

		if !found {
			t.Error("expected error about missing List template")
		}
	})

	// Test Tabs with missing child
	t.Run("Tabs", func(t *testing.T) {
		s := NewSurface("test")
		s.Add(Tabs("root", Tab("Tab 1", "missing-child"), Tab("Tab 2", "also-missing")))

		errors := s.Validate()

		if len(errors) == 0 {
			t.Fatal("expected validation errors for missing Tabs children")
		}

		foundTab0 := false
		foundTab1 := false
		for _, err := range errors {
			if err.ComponentID == "root" && err.Field == "Tabs.Tabs[0].Child" {
				foundTab0 = true
			}
			if err.ComponentID == "root" && err.Field == "Tabs.Tabs[1].Child" {
				foundTab1 = true
			}
		}

		if !foundTab0 || !foundTab1 {
			t.Error("expected errors about missing Tabs children")
		}
	})

	// Test Modal with missing children
	t.Run("Modal", func(t *testing.T) {
		s := NewSurface("test")
		s.Add(Modal("root", "missing-entry", "missing-content"))

		errors := s.Validate()

		if len(errors) == 0 {
			t.Fatal("expected validation errors for missing Modal children")
		}

		foundEntry := false
		foundContent := false
		for _, err := range errors {
			if err.ComponentID == "root" && err.Field == "Modal.EntryPointChild" {
				foundEntry = true
			}
			if err.ComponentID == "root" && err.Field == "Modal.ContentChild" {
				foundContent = true
			}
		}

		if !foundEntry || !foundContent {
			t.Error("expected errors about missing Modal children")
		}
	})
}

func TestValidateValid(t *testing.T) {
	// Test simple valid surface
	t.Run("SimpleValid", func(t *testing.T) {
		s := NewSurface("test")
		s.Add(Column("root", "text"))
		s.Add(TextStatic("text", "Hello World"))

		errors := s.Validate()

		if len(errors) != 0 {
			t.Errorf("expected no validation errors, got %d: %v", len(errors), errors)
		}
	})

	// Test complex valid surface with all component types
	t.Run("ComplexValid", func(t *testing.T) {
		s := NewSurface("test")
		s.Add(Column("root", "row", "card", "tabs", "modal", "list"))
		s.Add(Row("row", "text1", "text2"))
		s.Add(TextStatic("text1", "Hello"))
		s.Add(TextStatic("text2", "World"))
		s.Add(Card("card", "text3"))
		s.Add(TextStatic("text3", "Card Content"))
		s.Add(Tabs("tabs", Tab("Tab 1", "tab1-content"), Tab("Tab 2", "tab2-content")))
		s.Add(TextStatic("tab1-content", "Tab 1 Content"))
		s.Add(TextStatic("tab2-content", "Tab 2 Content"))
		s.Add(Modal("modal", "modal-trigger", "modal-dialog"))
		s.AddAll(Button("modal-trigger", "Open", "open")...)
		s.Add(TextStatic("modal-dialog", "Modal Content"))
		s.Add(ListTemplate("list", "list-item", "/items"))
		s.Add(TextBound("list-item", "/name"))

		errors := s.Validate()

		if len(errors) != 0 {
			t.Errorf("expected no validation errors, got %d: %v", len(errors), errors)
		}
	})

	// Test nested layout components
	t.Run("NestedLayouts", func(t *testing.T) {
		s := NewSurface("test")
		s.Add(Column("root", "row1", "row2"))
		s.Add(Row("row1", "col1", "col2"))
		s.Add(Column("col1", "text1"))
		s.Add(Column("col2", "text2"))
		s.Add(TextStatic("text1", "Cell 1"))
		s.Add(TextStatic("text2", "Cell 2"))
		s.Add(Row("row2", "text3", "text4"))
		s.Add(TextStatic("text3", "Cell 3"))
		s.Add(TextStatic("text4", "Cell 4"))

		errors := s.Validate()

		if len(errors) != 0 {
			t.Errorf("expected no validation errors for nested layouts, got %d: %v", len(errors), errors)
		}
	})
}

func TestValidationErrorFormat(t *testing.T) {
	// Test error format with ComponentID
	err1 := ValidationError{
		ComponentID: "my-component",
		Field:       "Child",
		Message:     "not found",
	}

	expected1 := "my-component: Child - not found"
	if err1.Error() != expected1 {
		t.Errorf("expected error '%s', got '%s'", expected1, err1.Error())
	}

	// Test error format without ComponentID
	err2 := ValidationError{
		ComponentID: "",
		Field:       "ID",
		Message:     "must not be empty",
	}

	expected2 := "ID: must not be empty"
	if err2.Error() != expected2 {
		t.Errorf("expected error '%s', got '%s'", expected2, err2.Error())
	}
}

func TestValidateMultipleErrors(t *testing.T) {
	s := NewSurface("test")
	s.SetRoot("missing-root")
	s.Add(TextStatic("duplicate", "First"))
	s.Add(TextStatic("duplicate", "Second"))
	s.Add(Component{ID: "", Component: "Text", Text: "Empty"})
	s.Add(Column("col", "missing-child"))

	errors := s.Validate()

	if len(errors) < 4 {
		t.Errorf("expected at least 4 errors, got %d", len(errors))
	}

	// Check we have various error types
	hasRootError := false
	hasDuplicateError := false
	hasEmptyIDError := false
	hasMissingChildError := false

	for _, err := range errors {
		if err.Field == "Root" {
			hasRootError = true
		}
		if err.Field == "ID" && err.Message == "duplicate component ID" {
			hasDuplicateError = true
		}
		if err.Field == "ID" && err.Message == "component ID must not be empty" {
			hasEmptyIDError = true
		}
		if err.Field == "Column.Children" {
			hasMissingChildError = true
		}
	}

	if !hasRootError {
		t.Error("expected root error")
	}
	if !hasDuplicateError {
		t.Error("expected duplicate ID error")
	}
	if !hasEmptyIDError {
		t.Error("expected empty ID error")
	}
	if !hasMissingChildError {
		t.Error("expected missing child error")
	}
}

func TestFlatComponentStructure(t *testing.T) {
	// Test that the flat structure serializes correctly
	txt := TextStatic("my-text", "Hello World")

	data, err := json.Marshal(txt)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	// Should have "component": "Text" instead of nested "Text": {...}
	jsonStr := string(data)
	if !strings.Contains(jsonStr, `"component":"Text"`) {
		t.Errorf("expected flat structure with component field, got: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"text":"Hello World"`) {
		t.Errorf("expected text field at top level, got: %s", jsonStr)
	}
}

func TestUpdateComponentsMessageType(t *testing.T) {
	s := NewSurface("test")
	s.Add(TextStatic("root", "Hello"))

	msgs := s.Messages()

	data, err := json.Marshal(msgs[1])
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	jsonStr := string(data)
	if !strings.Contains(jsonStr, `"updateComponents"`) {
		t.Errorf("expected 'updateComponents' message type, got: %s", jsonStr)
	}
	if strings.Contains(jsonStr, `"surfaceUpdate"`) {
		t.Errorf("should not contain old 'surfaceUpdate' message type, got: %s", jsonStr)
	}
}

func TestButtonChildReference(t *testing.T) {
	btns := Button("my-btn", "Click Me", "submit")

	// Verify button has child reference
	btn := btns[0]
	if btn.Child != "my-btn_text" {
		t.Errorf("expected child 'my-btn_text', got '%s'", btn.Child)
	}

	// Verify child text component
	txt := btns[1]
	if txt.ID != "my-btn_text" {
		t.Errorf("expected ID 'my-btn_text', got '%s'", txt.ID)
	}
	if txt.Text != "Click Me" {
		t.Errorf("expected text 'Click Me', got '%s'", txt.Text)
	}

	// Verify JSON serialization
	data, err := json.Marshal(btn)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	jsonStr := string(data)
	if !strings.Contains(jsonStr, `"child":"my-btn_text"`) {
		t.Errorf("expected child field in button, got: %s", jsonStr)
	}
	// Should not have "text" field directly on button
	if strings.Contains(jsonStr, `"text":"Click Me"`) {
		t.Errorf("button should not have text field directly, got: %s", jsonStr)
	}
}

// Custom component type for testing
type Gauge struct {
	Component
	Color      string `json:"color"`
	ShowLabel  bool   `json:"showLabel"`
	Thresholds []int  `json:"thresholds,omitempty"`
}

// Another custom component type
type SparklineChart struct {
	Component
	Data   []float64 `json:"data"`
	Color  string    `json:"color"`
	Height int       `json:"height"`
}

func TestCustomComponentWithEmbedding(t *testing.T) {
	// Create a custom component with embedded Component
	gauge := Gauge{
		Component: Component{
			ID:        "temp-gauge",
			Component: "Gauge",
			Label:     "Temperature",
			MinValue:  0,
			MaxValue:  100,
		},
		Color:      "#ff5500",
		ShowLabel:  true,
		Thresholds: []int{30, 60, 90},
	}

	data, err := json.Marshal(gauge)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	jsonStr := string(data)

	// Verify standard fields from embedded Component
	if !strings.Contains(jsonStr, `"component":"Gauge"`) {
		t.Errorf("expected component field, got: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"label":"Temperature"`) {
		t.Errorf("expected label field, got: %s", jsonStr)
	}

	// Verify custom fields are at top level
	if !strings.Contains(jsonStr, `"color":"#ff5500"`) {
		t.Errorf("expected color field, got: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"showLabel":true`) {
		t.Errorf("expected showLabel field, got: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"thresholds":[30,60,90]`) {
		t.Errorf("expected thresholds field, got: %s", jsonStr)
	}
}

func TestCustomComponentInSurface(t *testing.T) {
	surface := NewSurface("custom-demo")

	// Add standard components
	surface.Add(Column("root", "gauge", "chart"))

	// Add custom Gauge component with embedding
	surface.Add(Gauge{
		Component: Component{
			ID:        "gauge",
			Component: "Gauge",
			Label:     "CPU Usage",
		},
		Color:     "#00ff00",
		ShowLabel: true,
	})

	// Add another custom component
	surface.Add(SparklineChart{
		Component: Component{
			ID:        "chart",
			Component: "SparklineChart",
		},
		Data:   []float64{10, 25, 30, 45, 60},
		Color:  "#0066ff",
		Height: 50,
	})

	var buf bytes.Buffer
	err := WriteJSONL(&buf, surface.Messages())
	if err != nil {
		t.Fatalf("WriteJSONL failed: %v", err)
	}

	output := buf.String()

	// Verify custom components are serialized correctly
	if !strings.Contains(output, `"component":"Gauge"`) {
		t.Errorf("expected Gauge component, got: %s", output)
	}
	if !strings.Contains(output, `"component":"SparklineChart"`) {
		t.Errorf("expected SparklineChart component, got: %s", output)
	}
	if !strings.Contains(output, `"color":"#00ff00"`) {
		t.Errorf("expected color field for Gauge, got: %s", output)
	}
	if !strings.Contains(output, `"height":50`) {
		t.Errorf("expected height field for SparklineChart, got: %s", output)
	}
}
