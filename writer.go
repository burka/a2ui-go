package a2ui

import (
	"encoding/json"
	"io"
)

// WriteJSONL writes messages as JSON Lines (one JSON object per line).
func WriteJSONL(w io.Writer, messages []Message) error {
	for _, msg := range messages {
		data, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		if _, err := w.Write(data); err != nil {
			return err
		}
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
	}
	return nil
}

// WriteMessage writes a single message as JSON followed by a newline.
func WriteMessage(w io.Writer, msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	_, err = w.Write([]byte("\n"))
	return err
}

// WritePretty writes messages as indented JSON (for debugging).
func WritePretty(w io.Writer, messages []Message) error {
	for i, msg := range messages {
		data, err := json.MarshalIndent(msg, "", "  ")
		if err != nil {
			return err
		}
		if _, err := w.Write(data); err != nil {
			return err
		}
		if i < len(messages)-1 {
			if _, err := w.Write([]byte("\n")); err != nil {
				return err
			}
		}
	}
	return nil
}
