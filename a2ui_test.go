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
