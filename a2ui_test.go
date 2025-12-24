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

	// Check surfaceUpdate
	if msgs[1].SurfaceUpdate == nil {
		t.Error("expected SurfaceUpdate message")
	}
	if len(msgs[1].SurfaceUpdate.Components) != 2 {
		t.Errorf("expected 2 components, got %d", len(msgs[1].SurfaceUpdate.Components))
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
	if c.Column == nil {
		t.Fatal("expected Column to be set")
	}
	if len(c.Column.Children) != 3 {
		t.Errorf("expected 3 children, got %d", len(c.Column.Children))
	}
}

func TestRowHelper(t *testing.T) {
	r := Row("row", "x", "y")

	if r.ID != "row" {
		t.Errorf("expected ID 'row', got '%s'", r.ID)
	}
	if r.Row == nil {
		t.Fatal("expected Row to be set")
	}
	if len(r.Row.Children) != 2 {
		t.Errorf("expected 2 children, got %d", len(r.Row.Children))
	}
}

func TestCardHelper(t *testing.T) {
	c := Card("card", "content")

	if c.ID != "card" {
		t.Errorf("expected ID 'card', got '%s'", c.ID)
	}
	if c.Card == nil {
		t.Fatal("expected Card to be set")
	}
	if c.Card.Child != "content" {
		t.Errorf("expected child 'content', got '%s'", c.Card.Child)
	}
}

func TestTextStaticHelper(t *testing.T) {
	txt := TextStatic("txt", "Hello World")

	if txt.ID != "txt" {
		t.Errorf("expected ID 'txt', got '%s'", txt.ID)
	}
	if txt.Text == nil {
		t.Fatal("expected Text to be set")
	}
	if txt.Text.Text != "Hello World" {
		t.Errorf("expected text 'Hello World', got '%s'", txt.Text.Text)
	}
	if txt.Text.DataBinding != nil {
		t.Error("expected no DataBinding for static text")
	}
}

func TestTextBoundHelper(t *testing.T) {
	txt := TextBound("txt", "/user/name")

	if txt.Text == nil {
		t.Fatal("expected Text to be set")
	}
	if txt.Text.DataBinding == nil {
		t.Fatal("expected DataBinding to be set")
	}
	if txt.Text.DataBinding.Path != "/user/name" {
		t.Errorf("expected path '/user/name', got '%s'", txt.Text.DataBinding.Path)
	}
}

func TestListTemplateHelper(t *testing.T) {
	lst := ListTemplate("list", "item-template", "/items")

	if lst.ID != "list" {
		t.Errorf("expected ID 'list', got '%s'", lst.ID)
	}
	if lst.List == nil {
		t.Fatal("expected List to be set")
	}
	if lst.List.Template != "item-template" {
		t.Errorf("expected template 'item-template', got '%s'", lst.List.Template)
	}
	if lst.List.DataBinding == nil {
		t.Fatal("expected DataBinding to be set")
	}
	if lst.List.DataBinding.Path != "/items" {
		t.Errorf("expected path '/items', got '%s'", lst.List.DataBinding.Path)
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
	btn := Button("btn", "Click Me", "submit")

	if btn.ID != "btn" {
		t.Errorf("expected ID 'btn', got '%s'", btn.ID)
	}
	if btn.Button == nil {
		t.Fatal("expected Button to be set")
	}
	if btn.Button.Text != "Click Me" {
		t.Errorf("expected text 'Click Me', got '%s'", btn.Button.Text)
	}
	if btn.Button.Action.Type != "submit" {
		t.Errorf("expected action type 'submit', got '%s'", btn.Button.Action.Type)
	}
}

func TestTextFieldHelper(t *testing.T) {
	tf := TextField("input", "Email", "Enter email")

	if tf.ID != "input" {
		t.Errorf("expected ID 'input', got '%s'", tf.ID)
	}
	if tf.TextField == nil {
		t.Fatal("expected TextField to be set")
	}
	if tf.TextField.Label != "Email" {
		t.Errorf("expected label 'Email', got '%s'", tf.TextField.Label)
	}
	if tf.TextField.Placeholder != "Enter email" {
		t.Errorf("expected placeholder 'Enter email', got '%s'", tf.TextField.Placeholder)
	}
}

func TestImageHelpers(t *testing.T) {
	img := ImageStatic("img", "https://example.com/photo.jpg", "A photo")

	if img.ID != "img" {
		t.Errorf("expected ID 'img', got '%s'", img.ID)
	}
	if img.Image == nil {
		t.Fatal("expected Image to be set")
	}
	if img.Image.URL != "https://example.com/photo.jpg" {
		t.Errorf("expected URL 'https://example.com/photo.jpg', got '%s'", img.Image.URL)
	}
	if img.Image.Alt != "A photo" {
		t.Errorf("expected alt 'A photo', got '%s'", img.Image.Alt)
	}

	// Test bound image
	imgBound := ImageBound("img2", "/photo/url", "Dynamic photo")
	if imgBound.Image.DataBinding == nil {
		t.Fatal("expected DataBinding to be set")
	}
	if imgBound.Image.DataBinding.Path != "/photo/url" {
		t.Errorf("expected path '/photo/url', got '%s'", imgBound.Image.DataBinding.Path)
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
	btn := ButtonWithData("btn", "Submit", "submit", map[string]any{
		"endpoint": "/api/submit",
		"method":   "POST",
	})

	if btn.ID != "btn" {
		t.Errorf("expected ID 'btn', got '%s'", btn.ID)
	}
	if btn.Button == nil {
		t.Fatal("expected Button to be set")
	}
	if btn.Button.Text != "Submit" {
		t.Errorf("expected text 'Submit', got '%s'", btn.Button.Text)
	}
	if btn.Button.Action.Type != "submit" {
		t.Errorf("expected action type 'submit', got '%s'", btn.Button.Action.Type)
	}
	if btn.Button.Action.Data == nil {
		t.Fatal("expected Action.Data to be set")
	}
	if btn.Button.Action.Data["endpoint"] != "/api/submit" {
		t.Errorf("expected endpoint '/api/submit', got '%v'", btn.Button.Action.Data["endpoint"])
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
	if tabs.Tabs == nil {
		t.Fatal("expected Tabs to be set")
	}
	if len(tabs.Tabs.Tabs) != 2 {
		t.Errorf("expected 2 tabs, got %d", len(tabs.Tabs.Tabs))
	}
	if tabs.Tabs.Tabs[0].Title != "Tab 1" {
		t.Errorf("expected title 'Tab 1', got '%s'", tabs.Tabs.Tabs[0].Title)
	}
	if tabs.Tabs.Tabs[0].Child != "content1" {
		t.Errorf("expected child 'content1', got '%s'", tabs.Tabs.Tabs[0].Child)
	}
}

func TestModalHelper(t *testing.T) {
	modal := Modal("modal", "trigger-btn", "dialog-content")

	if modal.ID != "modal" {
		t.Errorf("expected ID 'modal', got '%s'", modal.ID)
	}
	if modal.Modal == nil {
		t.Fatal("expected Modal to be set")
	}
	if modal.Modal.EntryPointChild != "trigger-btn" {
		t.Errorf("expected entryPointChild 'trigger-btn', got '%s'", modal.Modal.EntryPointChild)
	}
	if modal.Modal.ContentChild != "dialog-content" {
		t.Errorf("expected contentChild 'dialog-content', got '%s'", modal.Modal.ContentChild)
	}
}

func TestIconHelper(t *testing.T) {
	icon := Icon("icon", IconSearch)

	if icon.ID != "icon" {
		t.Errorf("expected ID 'icon', got '%s'", icon.ID)
	}
	if icon.Icon == nil {
		t.Fatal("expected Icon to be set")
	}
	if icon.Icon.Icon != IconSearch {
		t.Errorf("expected icon 'search', got '%s'", icon.Icon.Icon)
	}
}

func TestVideoHelpers(t *testing.T) {
	video := Video("video", "https://example.com/video.mp4")

	if video.ID != "video" {
		t.Errorf("expected ID 'video', got '%s'", video.ID)
	}
	if video.Video == nil {
		t.Fatal("expected Video to be set")
	}
	if video.Video.URL != "https://example.com/video.mp4" {
		t.Errorf("expected URL 'https://example.com/video.mp4', got '%s'", video.Video.URL)
	}

	// Test bound video
	videoBound := VideoBound("video2", "/media/video")
	if videoBound.Video.DataBinding == nil {
		t.Fatal("expected DataBinding to be set")
	}
	if videoBound.Video.DataBinding.Path != "/media/video" {
		t.Errorf("expected path '/media/video', got '%s'", videoBound.Video.DataBinding.Path)
	}
}

func TestAudioPlayerHelpers(t *testing.T) {
	audio := AudioPlayer("audio", "https://example.com/song.mp3", "My Song")

	if audio.ID != "audio" {
		t.Errorf("expected ID 'audio', got '%s'", audio.ID)
	}
	if audio.AudioPlayer == nil {
		t.Fatal("expected AudioPlayer to be set")
	}
	if audio.AudioPlayer.URL != "https://example.com/song.mp3" {
		t.Errorf("expected URL 'https://example.com/song.mp3', got '%s'", audio.AudioPlayer.URL)
	}
	if audio.AudioPlayer.Description != "My Song" {
		t.Errorf("expected description 'My Song', got '%s'", audio.AudioPlayer.Description)
	}

	// Test bound audio
	audioBound := AudioPlayerBound("audio2", "/media/audio", "Dynamic Audio")
	if audioBound.AudioPlayer.DataBinding == nil {
		t.Fatal("expected DataBinding to be set")
	}
	if audioBound.AudioPlayer.DataBinding.Path != "/media/audio" {
		t.Errorf("expected path '/media/audio', got '%s'", audioBound.AudioPlayer.DataBinding.Path)
	}
}

func TestDividerHelpers(t *testing.T) {
	div := Divider("div")

	if div.ID != "div" {
		t.Errorf("expected ID 'div', got '%s'", div.ID)
	}
	if div.Divider == nil {
		t.Fatal("expected Divider to be set")
	}

	// Test vertical divider
	divV := DividerVertical("divV")
	if divV.Divider.Orientation != "vertical" {
		t.Errorf("expected orientation 'vertical', got '%s'", divV.Divider.Orientation)
	}
}

func TestCheckBoxHelpers(t *testing.T) {
	cb := CheckBox("cb", "Accept Terms", true)

	if cb.ID != "cb" {
		t.Errorf("expected ID 'cb', got '%s'", cb.ID)
	}
	if cb.CheckBox == nil {
		t.Fatal("expected CheckBox to be set")
	}
	if cb.CheckBox.Label != "Accept Terms" {
		t.Errorf("expected label 'Accept Terms', got '%s'", cb.CheckBox.Label)
	}
	if cb.CheckBox.Value != true {
		t.Errorf("expected value true, got false")
	}

	// Test bound checkbox
	cbBound := CheckBoxBound("cb2", "Subscribe", "/prefs/subscribe")
	if cbBound.CheckBox.DataBinding == nil {
		t.Fatal("expected DataBinding to be set")
	}
	if cbBound.CheckBox.DataBinding.Path != "/prefs/subscribe" {
		t.Errorf("expected path '/prefs/subscribe', got '%s'", cbBound.CheckBox.DataBinding.Path)
	}
}

func TestDateTimeInputHelpers(t *testing.T) {
	dt := DateTimeInput("dt", "Select Date", true, false)

	if dt.ID != "dt" {
		t.Errorf("expected ID 'dt', got '%s'", dt.ID)
	}
	if dt.DateTimeInput == nil {
		t.Fatal("expected DateTimeInput to be set")
	}
	if dt.DateTimeInput.Label != "Select Date" {
		t.Errorf("expected label 'Select Date', got '%s'", dt.DateTimeInput.Label)
	}
	if dt.DateTimeInput.EnableDate != true {
		t.Error("expected enableDate true")
	}
	if dt.DateTimeInput.EnableTime != false {
		t.Error("expected enableTime false")
	}

	// Test bound datetime
	dtBound := DateTimeInputBound("dt2", "Appointment", "/booking/datetime", true, true)
	if dtBound.DateTimeInput.DataBinding == nil {
		t.Fatal("expected DataBinding to be set")
	}
	if dtBound.DateTimeInput.DataBinding.Path != "/booking/datetime" {
		t.Errorf("expected path '/booking/datetime', got '%s'", dtBound.DateTimeInput.DataBinding.Path)
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
	if mc.MultipleChoice == nil {
		t.Fatal("expected MultipleChoice to be set")
	}
	if mc.MultipleChoice.Label != "Select Size" {
		t.Errorf("expected label 'Select Size', got '%s'", mc.MultipleChoice.Label)
	}
	if len(mc.MultipleChoice.Options) != 3 {
		t.Errorf("expected 3 options, got %d", len(mc.MultipleChoice.Options))
	}
	if mc.MultipleChoice.Options[0].Label != "Small" {
		t.Errorf("expected first option label 'Small', got '%s'", mc.MultipleChoice.Options[0].Label)
	}
	if mc.MultipleChoice.Options[0].Value != "s" {
		t.Errorf("expected first option value 's', got '%s'", mc.MultipleChoice.Options[0].Value)
	}

	// Test bound multiple choice
	mcBound := MultipleChoiceBound("mc2", "Color", "/prefs/color", []ChoiceOption{
		Choice("Red", "red"),
		Choice("Blue", "blue"),
	})
	if mcBound.MultipleChoice.DataBinding == nil {
		t.Fatal("expected DataBinding to be set")
	}
	if mcBound.MultipleChoice.DataBinding.Path != "/prefs/color" {
		t.Errorf("expected path '/prefs/color', got '%s'", mcBound.MultipleChoice.DataBinding.Path)
	}
}

func TestSliderHelpers(t *testing.T) {
	slider := Slider("slider", "Volume", 0, 100, 50)

	if slider.ID != "slider" {
		t.Errorf("expected ID 'slider', got '%s'", slider.ID)
	}
	if slider.Slider == nil {
		t.Fatal("expected Slider to be set")
	}
	if slider.Slider.Label != "Volume" {
		t.Errorf("expected label 'Volume', got '%s'", slider.Slider.Label)
	}
	if slider.Slider.MinValue != 0 {
		t.Errorf("expected minValue 0, got %f", slider.Slider.MinValue)
	}
	if slider.Slider.MaxValue != 100 {
		t.Errorf("expected maxValue 100, got %f", slider.Slider.MaxValue)
	}
	if slider.Slider.Value != 50 {
		t.Errorf("expected value 50, got %f", slider.Slider.Value)
	}

	// Test bound slider
	sliderBound := SliderBound("slider2", "Brightness", "/settings/brightness", 0, 255)
	if sliderBound.Slider.DataBinding == nil {
		t.Fatal("expected DataBinding to be set")
	}
	if sliderBound.Slider.DataBinding.Path != "/settings/brightness" {
		t.Errorf("expected path '/settings/brightness', got '%s'", sliderBound.Slider.DataBinding.Path)
	}
}

func TestButtonPrimaryHelper(t *testing.T) {
	btn := ButtonPrimary("btn", "Save", "save")

	if btn.Button == nil {
		t.Fatal("expected Button to be set")
	}
	if btn.Button.Primary != true {
		t.Error("expected primary true")
	}
}

func TestTextWithHintHelper(t *testing.T) {
	txt := TextWithHint("title", "Welcome", UsageHintH1)

	if txt.Text == nil {
		t.Fatal("expected Text to be set")
	}
	if txt.Text.UsageHint != UsageHintH1 {
		t.Errorf("expected usageHint 'h1', got '%s'", txt.Text.UsageHint)
	}
}

func TestImageWithFitHelper(t *testing.T) {
	img := ImageWithFit("img", "https://example.com/photo.jpg", "Photo", ImageFitCover)

	if img.Image == nil {
		t.Fatal("expected Image to be set")
	}
	if img.Image.Fit != ImageFitCover {
		t.Errorf("expected fit 'cover', got '%s'", img.Image.Fit)
	}
}

func TestTextFieldWithTypeHelper(t *testing.T) {
	tf := TextFieldWithType("password", "Password", "Enter password", TextFieldTypeObscured)

	if tf.TextField == nil {
		t.Fatal("expected TextField to be set")
	}
	if tf.TextField.TextFieldType != TextFieldTypeObscured {
		t.Errorf("expected textFieldType 'obscured', got '%s'", tf.TextField.TextFieldType)
	}
}

func TestColumnWithLayoutHelper(t *testing.T) {
	col := ColumnWithLayout("col", DistributionSpaceBetween, AlignmentCenter, "a", "b")

	if col.Column == nil {
		t.Fatal("expected Column to be set")
	}
	if col.Column.Distribution != DistributionSpaceBetween {
		t.Errorf("expected distribution 'spaceBetween', got '%s'", col.Column.Distribution)
	}
	if col.Column.Alignment != AlignmentCenter {
		t.Errorf("expected alignment 'center', got '%s'", col.Column.Alignment)
	}
}

func TestRowWithLayoutHelper(t *testing.T) {
	row := RowWithLayout("row", DistributionSpaceEvenly, AlignmentStretch, "x", "y")

	if row.Row == nil {
		t.Fatal("expected Row to be set")
	}
	if row.Row.Distribution != DistributionSpaceEvenly {
		t.Errorf("expected distribution 'spaceEvenly', got '%s'", row.Row.Distribution)
	}
	if row.Row.Alignment != AlignmentStretch {
		t.Errorf("expected alignment 'stretch', got '%s'", row.Row.Alignment)
	}
}

func TestTextFieldBoundHelper(t *testing.T) {
	tf := TextFieldBound("email", "Email", "Enter email", "/user/email")

	if tf.TextField == nil {
		t.Fatal("expected TextField to be set")
	}
	if tf.TextField.DataBinding == nil {
		t.Fatal("expected DataBinding to be set")
	}
	if tf.TextField.DataBinding.Path != "/user/email" {
		t.Errorf("expected path '/user/email', got '%s'", tf.TextField.DataBinding.Path)
	}
}

func TestAllComponentsJSONSerialization(t *testing.T) {
	s := NewSurface("all-components")

	// Add all 18 component types
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
		Button("button", "Click", "action"),
		TextField("textfield", "Label", "Placeholder"),
		CheckBox("checkbox", "Check", false),
		DateTimeInput("datetime", "Date", true, true),
		MultipleChoice("choice", "Select", []ChoiceOption{Choice("A", "a")}),
		Slider("slider", "Slide", 0, 100, 50),
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

	// Verify we have all 18 components
	if msgs[1].SurfaceUpdate == nil {
		t.Fatal("expected SurfaceUpdate")
	}
	if len(msgs[1].SurfaceUpdate.Components) != 18 {
		t.Errorf("expected 18 components, got %d", len(msgs[1].SurfaceUpdate.Components))
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
	s.Add(Component{ID: "", Text: &TextDef{Text: "Empty ID"}})

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
		s.Add(Button("modal-trigger", "Open", "open"))
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
	s.Add(Component{ID: "", Text: &TextDef{Text: "Empty"}})
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
