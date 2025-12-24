# Using Official A2UI Renderers

This Go library generates A2UI messages. To render them, use Google's official client renderers.

## Flutter (GenUI SDK)

```bash
# Add to pubspec.yaml
dependencies:
  google_genui: ^0.1.0
```

```dart
import 'package:google_genui/google_genui.dart';

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return GenUIRenderer(
      // Point to your Go backend
      streamUrl: 'http://localhost:8080/plan',
      onAction: (action) {
        // Handle button clicks, form submits
        print('Action: ${action.type} ${action.data}');
      },
    );
  }
}
```

Docs: https://docs.flutter.dev/ai/genui

## Web (Lit Renderer)

```bash
git clone https://github.com/google/A2UI.git
cd A2UI/renderers/lit
npm install
```

```html
<script type="module">
  import { A2UIRenderer } from './a2ui-renderer.js';

  const renderer = new A2UIRenderer({
    container: document.getElementById('app'),
    onAction: (action) => {
      fetch(action.data.endpoint, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({event: {type: action.type, data: formData}})
      });
    }
  });

  // Connect to Go backend
  fetch('http://localhost:8080/plan')
    .then(r => r.text())
    .then(text => {
      text.split('\n').filter(l => l).forEach(line => {
        renderer.processMessage(JSON.parse(line));
      });
    });
</script>
```

Source: https://github.com/google/A2UI/tree/main/renderers/lit

## Angular Renderer

```bash
git clone https://github.com/google/A2UI.git
cd A2UI/renderers/angular
npm install
ng serve
```

```typescript
import { A2UIModule } from './a2ui.module';

@Component({
  template: `<a2ui-surface [streamUrl]="'http://localhost:8080/plan'"></a2ui-surface>`
})
export class AppComponent {}
```

Source: https://github.com/google/A2UI/tree/main/renderers/angular

## CORS Setup

Your Go backend needs CORS headers for browser clients:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    if r.Method == "OPTIONS" {
        return
    }

    // ... your A2UI code
}
```

## Testing Without Official Renderers

Our `examples/` include simple HTML/JS clients for testing:

```bash
cd examples/streaming && go run main.go   # http://localhost:8080
cd examples/interactive && go run main.go # http://localhost:8080
```

These are minimal implementations for development/debugging.
