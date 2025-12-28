// Package main demonstrates A2UI streaming with progressive UI updates.
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	a2ui "github.com/burka/a2ui-go"
)

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/plan", handlePlan)

	fmt.Println("A2UI Streaming Demo")
	fmt.Println("Open http://localhost:8080 in your browser")
	fmt.Println("Or: curl http://localhost:8080/plan")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(indexHTML))
}

func handlePlan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	surface := a2ui.NewSurface("itinerary")

	// Step 1: Initial structure with loading message
	surface.Add(a2ui.Column("root", "header", "content", "footer"))
	surface.Add(a2ui.TextStatic("header", "Your Travel Itinerary"))
	surface.Add(a2ui.Column("content", "loading"))
	surface.Add(a2ui.TextStatic("loading", "Planning your trip..."))
	surface.Add(a2ui.TextStatic("footer", ""))

	a2ui.WriteJSONL(w, surface.Messages())
	flusher.Flush()
	time.Sleep(800 * time.Millisecond)

	// Step 2: Day 1
	surface.Add(a2ui.Card("day1", "day1-content"))
	surface.Add(a2ui.Column("day1-content", "day1-title", "day1-activities"))
	surface.Add(a2ui.TextBound("day1-title", "/day1/title"))
	surface.Add(a2ui.ListTemplate("day1-activities", "activity", "/day1/activities"))
	surface.Add(a2ui.TextBound("activity", "/name"))

	surface.SetData("/day1/title", "Day 1: Arrival")
	surface.SetData("/day1/activities", []map[string]string{
		{"name": "Airport pickup at 10:00 AM"},
		{"name": "Hotel check-in at Grand Hotel"},
		{"name": "Welcome dinner at La Terrazza"},
	})

	// Update content to show day1 instead of loading
	surface.Add(a2ui.Column("content", "day1"))

	a2ui.WriteMessage(w, surface.UpdateComponentsMessage())
	a2ui.WriteMessage(w, surface.DataModelUpdateMessage())
	flusher.Flush()
	time.Sleep(1000 * time.Millisecond)

	// Step 3: Day 2
	surface.Add(a2ui.Card("day2", "day2-content"))
	surface.Add(a2ui.Column("day2-content", "day2-title", "day2-activities"))
	surface.Add(a2ui.TextBound("day2-title", "/day2/title"))
	surface.Add(a2ui.ListTemplate("day2-activities", "activity2", "/day2/activities"))
	surface.Add(a2ui.TextBound("activity2", "/name"))

	surface.SetData("/day2/title", "Day 2: City Exploration")
	surface.SetData("/day2/activities", []map[string]string{
		{"name": "Breakfast at hotel"},
		{"name": "City walking tour (9:00 AM - 12:00 PM)"},
		{"name": "Lunch at local market"},
		{"name": "Museum visit (2:00 PM - 5:00 PM)"},
		{"name": "Free evening"},
	})

	// Update content to show both days
	surface.Add(a2ui.Column("content", "day1", "day2"))

	a2ui.WriteMessage(w, surface.UpdateComponentsMessage())
	a2ui.WriteMessage(w, surface.DataModelUpdateMessage())
	flusher.Flush()
	time.Sleep(1000 * time.Millisecond)

	// Step 4: Day 3
	surface.Add(a2ui.Card("day3", "day3-content"))
	surface.Add(a2ui.Column("day3-content", "day3-title", "day3-activities"))
	surface.Add(a2ui.TextBound("day3-title", "/day3/title"))
	surface.Add(a2ui.ListTemplate("day3-activities", "activity3", "/day3/activities"))
	surface.Add(a2ui.TextBound("activity3", "/name"))

	surface.SetData("/day3/title", "Day 3: Departure")
	surface.SetData("/day3/activities", []map[string]string{
		{"name": "Breakfast and checkout"},
		{"name": "Souvenir shopping"},
		{"name": "Airport transfer at 2:00 PM"},
		{"name": "Flight departure at 5:00 PM"},
	})

	surface.Add(a2ui.Column("content", "day1", "day2", "day3"))

	a2ui.WriteMessage(w, surface.UpdateComponentsMessage())
	a2ui.WriteMessage(w, surface.DataModelUpdateMessage())
	flusher.Flush()
	time.Sleep(800 * time.Millisecond)

	// Step 5: Summary
	surface.Add(a2ui.Card("summary", "summary-content"))
	surface.Add(a2ui.Column("summary-content", "summary-title", "summary-total"))
	surface.Add(a2ui.TextStatic("summary-title", "Trip Summary"))
	surface.Add(a2ui.TextBound("summary-total", "/summary/total"))

	surface.SetData("/summary/total", "Total estimated cost: $1,250")

	surface.Add(a2ui.Column("content", "day1", "day2", "day3", "summary"))
	surface.Add(a2ui.TextStatic("footer", "Have a great trip!"))

	a2ui.WriteMessage(w, surface.UpdateComponentsMessage())
	a2ui.WriteMessage(w, surface.DataModelUpdateMessage())
	flusher.Flush()
}

const indexHTML = `<!DOCTYPE html>
<html>
<head>
    <title>A2UI Streaming Demo</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            max-width: 800px;
            margin: 40px auto;
            padding: 20px;
            background: #f5f5f5;
        }
        h1 { color: #333; }
        button {
            background: #007bff;
            color: white;
            border: none;
            padding: 12px 24px;
            font-size: 16px;
            cursor: pointer;
            border-radius: 6px;
        }
        button:hover { background: #0056b3; }
        button:disabled { background: #ccc; cursor: not-allowed; }
        #output {
            margin-top: 20px;
            background: #1e1e1e;
            color: #d4d4d4;
            padding: 20px;
            border-radius: 8px;
            font-family: 'Monaco', 'Consolas', monospace;
            font-size: 13px;
            max-height: 600px;
            overflow-y: auto;
            white-space: pre-wrap;
            word-break: break-all;
        }
        .message {
            margin-bottom: 12px;
            padding: 10px;
            background: #2d2d2d;
            border-radius: 4px;
            border-left: 3px solid #007bff;
        }
        .message.begin { border-left-color: #28a745; }
        .message.update { border-left-color: #17a2b8; }
        .message.data { border-left-color: #ffc107; }
    </style>
</head>
<body>
    <h1>A2UI Streaming Demo</h1>
    <p>Click the button to stream a travel itinerary progressively.</p>
    <button id="planBtn" onclick="planTrip()">Plan Trip</button>
    <div id="output"></div>

    <script>
        async function planTrip() {
            const btn = document.getElementById('planBtn');
            const output = document.getElementById('output');

            btn.disabled = true;
            output.innerHTML = '';

            try {
                const response = await fetch('/plan');
                const reader = response.body.getReader();
                const decoder = new TextDecoder();

                let buffer = '';
                while (true) {
                    const {done, value} = await reader.read();
                    if (done) break;

                    buffer += decoder.decode(value, {stream: true});
                    const lines = buffer.split('\n');
                    buffer = lines.pop();

                    for (const line of lines) {
                        if (line.trim()) {
                            const msg = JSON.parse(line);
                            const div = document.createElement('div');
                            div.className = 'message';

                            if (msg.beginRendering) div.className += ' begin';
                            else if (msg.updateComponents) div.className += ' update';
                            else if (msg.dataModelUpdate) div.className += ' data';

                            div.textContent = JSON.stringify(msg, null, 2);
                            output.appendChild(div);
                            output.scrollTop = output.scrollHeight;
                        }
                    }
                }
            } catch (err) {
                output.textContent = 'Error: ' + err.message;
            }

            btn.disabled = false;
        }
    </script>
</body>
</html>
`
