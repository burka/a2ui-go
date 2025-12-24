import { Component } from '@angular/core';
import { A2UIRendererComponent } from './a2ui-renderer.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [A2UIRendererComponent],
  template: `
    <div class="container">
      <h1>A2UI Angular Client</h1>
      <p>Connected to Go backend at <code>{{ backendUrl }}</code></p>
      <a2ui-renderer [streamUrl]="backendUrl + '/form'"></a2ui-renderer>
    </div>
  `,
  styles: [`
    .container {
      max-width: 500px;
      margin: 40px auto;
      padding: 20px;
      font-family: system-ui, sans-serif;
    }
    h1 { color: #333; margin-bottom: 5px; }
    p { color: #666; margin-bottom: 20px; }
    code { background: #f0f0f0; padding: 2px 6px; border-radius: 3px; }
  `]
})
export class App {
  backendUrl = 'http://localhost:8080';
}
