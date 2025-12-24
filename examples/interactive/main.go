// Package main demonstrates A2UI interactive forms with client events.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	a2ui "github.com/burka/a2ui-go"
)

type Booking struct {
	ID    string
	Name  string
	Date  string
	Time  string
	Party int
}

var (
	bookings = make(map[string]Booking)
	mu       sync.RWMutex
)

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/form", handleForm)
	http.HandleFunc("/submit", handleSubmit)

	fmt.Println("A2UI Interactive Demo")
	fmt.Println("Open http://localhost:8080 in your browser")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(indexHTML))
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	surface := a2ui.NewSurface("booking-form")

	// Form layout
	surface.Add(a2ui.Column("root", "header", "form-card", "status"))
	surface.Add(a2ui.TextStatic("header", "Restaurant Booking"))

	// Form card with fields
	surface.Add(a2ui.Card("form-card", "form-content"))
	surface.Add(a2ui.Column("form-content",
		"name-field", "date-field", "time-field", "party-field", "submit-btn"))

	// Input fields
	surface.Add(a2ui.TextFieldBound("name-field", "Name", "Your name", "/form/name"))
	surface.Add(a2ui.TextFieldBound("date-field", "Date", "YYYY-MM-DD", "/form/date"))
	surface.Add(a2ui.TextFieldBound("time-field", "Time", "HH:MM", "/form/time"))
	surface.Add(a2ui.TextFieldBound("party-field", "Party Size", "Number of guests", "/form/party"))

	// Submit button with endpoint data
	surface.Add(a2ui.ButtonWithData("submit-btn", "Book Table", "submit",
		map[string]any{"endpoint": "/submit"}))

	surface.Add(a2ui.TextStatic("status", ""))

	// Set default form values
	surface.SetData("/form/name", "")
	surface.SetData("/form/date", time.Now().AddDate(0, 0, 1).Format("2006-01-02"))
	surface.SetData("/form/time", "19:00")
	surface.SetData("/form/party", "2")

	a2ui.WriteJSONL(w, surface.Messages())
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}

	// Parse client event
	var clientMsg a2ui.ClientMessage
	if err := json.NewDecoder(r.Body).Decode(&clientMsg); err != nil {
		log.Printf("Error decoding: %v", err)
		sendError(w, "Invalid request")
		return
	}

	log.Printf("Received event: %+v", clientMsg.Event)

	// Extract form data from event
	data := clientMsg.Event.Data
	name, _ := data["name"].(string)
	date, _ := data["date"].(string)
	timeStr, _ := data["time"].(string)
	partyFloat, _ := data["party"].(float64)
	party := int(partyFloat)

	if name == "" {
		name = "Guest"
	}
	if party == 0 {
		party = 2
	}

	// Create booking
	booking := Booking{
		ID:    fmt.Sprintf("BK-%d", time.Now().Unix()),
		Name:  name,
		Date:  date,
		Time:  timeStr,
		Party: party,
	}

	mu.Lock()
	bookings[booking.ID] = booking
	mu.Unlock()

	log.Printf("Created booking: %+v", booking)

	// Send confirmation UI
	surface := a2ui.NewSurface("confirmation")

	surface.Add(a2ui.Column("root", "success-card"))
	surface.Add(a2ui.Card("success-card", "confirm-content"))
	surface.Add(a2ui.Column("confirm-content",
		"title", "booking-id", "details", "back-btn"))

	surface.Add(a2ui.TextStatic("title", "Booking Confirmed!"))
	surface.Add(a2ui.TextBound("booking-id", "/booking/id"))
	surface.Add(a2ui.TextBound("details", "/booking/details"))
	surface.Add(a2ui.ButtonWithData("back-btn", "New Booking", "navigate",
		map[string]any{"url": "/form"}))

	surface.SetData("/booking/id", fmt.Sprintf("Confirmation: %s", booking.ID))
	surface.SetData("/booking/details",
		fmt.Sprintf("%s - %s at %s for %d guests", booking.Name, booking.Date, booking.Time, booking.Party))

	a2ui.WriteJSONL(w, surface.Messages())
}

func sendError(w http.ResponseWriter, msg string) {
	surface := a2ui.NewSurface("error")
	surface.Add(a2ui.Column("root", "error-msg"))
	surface.Add(a2ui.TextStatic("error-msg", "Error: "+msg))
	a2ui.WriteJSONL(w, surface.Messages())
}

const indexHTML = `<!DOCTYPE html>
<html>
<head>
    <title>A2UI Interactive Demo</title>
    <style>
        * { box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            max-width: 500px;
            margin: 40px auto;
            padding: 20px;
            background: #f5f5f5;
        }
        h1 { color: #333; margin-bottom: 5px; }
        .card {
            background: white;
            border: 1px solid #ddd;
            border-radius: 8px;
            padding: 20px;
            margin: 15px 0;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .card.success {
            background: #d4edda;
            border-color: #c3e6cb;
        }
        .field {
            margin-bottom: 15px;
        }
        label {
            display: block;
            margin-bottom: 5px;
            font-weight: 500;
            color: #333;
        }
        input {
            width: 100%;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 16px;
        }
        input:focus {
            outline: none;
            border-color: #007bff;
        }
        button {
            background: #007bff;
            color: white;
            border: none;
            padding: 12px 24px;
            border-radius: 4px;
            cursor: pointer;
            font-size: 16px;
            width: 100%;
            margin-top: 10px;
        }
        button:hover { background: #0056b3; }
        .header { font-size: 24px; font-weight: bold; margin-bottom: 10px; }
        .title { font-size: 20px; color: #155724; margin-bottom: 10px; }
        .details { color: #666; margin: 10px 0; }
        #debug {
            margin-top: 20px;
            background: #1e1e1e;
            color: #d4d4d4;
            padding: 15px;
            border-radius: 8px;
            font-family: monospace;
            font-size: 12px;
            max-height: 300px;
            overflow-y: auto;
            display: none;
        }
        #debug.show { display: block; }
        .toggle-debug {
            background: #6c757d;
            font-size: 12px;
            padding: 8px 16px;
            margin-top: 20px;
        }
    </style>
</head>
<body>
    <div id="ui"></div>
    <button class="toggle-debug" onclick="toggleDebug()">Toggle Debug</button>
    <pre id="debug"></pre>

    <script>
        let formData = {};

        loadUI('/form');

        function toggleDebug() {
            document.getElementById('debug').classList.toggle('show');
        }

        function log(msg) {
            const debug = document.getElementById('debug');
            debug.textContent += msg + '\n';
            debug.scrollTop = debug.scrollHeight;
        }

        async function loadUI(url) {
            log('Loading: ' + url);
            const response = await fetch(url);
            const text = await response.text();

            const messages = text.split('\n')
                .filter(line => line.trim())
                .map(line => JSON.parse(line));

            log('Received ' + messages.length + ' messages');
            renderUI(messages);
        }

        function renderUI(messages) {
            const ui = document.getElementById('ui');
            ui.innerHTML = '';
            formData = {};

            let components = {};
            let dataModel = {};

            messages.forEach(msg => {
                if (msg.surfaceUpdate) {
                    msg.surfaceUpdate.components.forEach(c => components[c.id] = c);
                }
                if (msg.dataModelUpdate) {
                    Object.assign(dataModel, msg.dataModelUpdate.contents);
                }
            });

            // Simple recursive render
            function render(id, container) {
                const comp = components[id];
                if (!comp) return;

                if (comp.Column) {
                    const div = document.createElement('div');
                    comp.Column.children.forEach(childId => render(childId, div));
                    container.appendChild(div);
                }
                else if (comp.Card) {
                    const card = document.createElement('div');
                    card.className = 'card';
                    if (id.includes('success')) card.className += ' success';
                    render(comp.Card.child, card);
                    container.appendChild(card);
                }
                else if (comp.Text) {
                    const div = document.createElement('div');
                    if (comp.Text.text) {
                        div.textContent = comp.Text.text;
                    } else if (comp.Text.dataBinding) {
                        div.textContent = dataModel[comp.Text.dataBinding.path] || '';
                    }
                    if (id === 'header') div.className = 'header';
                    if (id === 'title') div.className = 'title';
                    if (id.includes('details') || id.includes('booking')) div.className = 'details';
                    container.appendChild(div);
                }
                else if (comp.TextField) {
                    const field = document.createElement('div');
                    field.className = 'field';

                    const label = document.createElement('label');
                    label.textContent = comp.TextField.label;
                    field.appendChild(label);

                    const input = document.createElement('input');
                    input.placeholder = comp.TextField.placeholder || '';
                    input.id = 'input-' + id;

                    if (comp.TextField.dataBinding) {
                        const path = comp.TextField.dataBinding.path;
                        input.value = dataModel[path] || '';
                        formData[path] = input.value;
                        input.oninput = () => { formData[path] = input.value; };
                    }

                    field.appendChild(input);
                    container.appendChild(field);
                }
                else if (comp.Button) {
                    const btn = document.createElement('button');
                    btn.textContent = comp.Button.text;
                    btn.onclick = () => handleAction(comp.Button.action);
                    container.appendChild(btn);
                }
            }

            render('root', ui);
        }

        async function handleAction(action) {
            log('Action: ' + action.type + ' ' + JSON.stringify(action.data));

            if (action.type === 'submit' && action.data?.endpoint) {
                const event = {
                    event: {
                        surfaceId: 'booking-form',
                        componentId: 'submit-btn',
                        type: 'action',
                        data: {
                            name: formData['/form/name'] || '',
                            date: formData['/form/date'] || '',
                            time: formData['/form/time'] || '',
                            party: parseInt(formData['/form/party']) || 2
                        }
                    }
                };

                log('Sending: ' + JSON.stringify(event));

                const response = await fetch(action.data.endpoint, {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify(event)
                });

                const text = await response.text();
                const messages = text.split('\n')
                    .filter(line => line.trim())
                    .map(line => JSON.parse(line));

                renderUI(messages);
            }
            else if (action.type === 'navigate' && action.data?.url) {
                loadUI(action.data.url);
            }
        }
    </script>
</body>
</html>
`
