import { LitElement, html, css, PropertyValues } from 'lit';

interface A2UIComponent {
  id: string;
  Column?: { children: string[] };
  Row?: { children: string[] };
  Card?: { child: string };
  Text?: { text?: string; dataBinding?: { path: string } };
  TextField?: { label?: string; placeholder?: string; dataBinding?: { path: string } };
  Button?: { text: string; action: { type: string; data?: Record<string, any> } };
}

interface A2UIMessage {
  beginRendering?: { surfaceId: string; root: string };
  surfaceUpdate?: { surfaceId: string; components: A2UIComponent[] };
  dataModelUpdate?: { surfaceId: string; contents: Record<string, any> };
}

export class A2UIRenderer extends LitElement {
  static styles = css`
    :host { display: block; font-family: system-ui, sans-serif; }
    .column { display: flex; flex-direction: column; gap: 10px; }
    .row { display: flex; flex-direction: row; gap: 10px; }
    .card {
      background: white;
      border: 1px solid #ddd;
      border-radius: 8px;
      padding: 20px;
      box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }
    .card.success { background: #d4edda; border-color: #c3e6cb; }
    .text { margin: 5px 0; }
    .text.header { font-size: 24px; font-weight: bold; margin-bottom: 15px; }
    .text.title { font-size: 20px; color: #155724; }
    .field { margin-bottom: 15px; }
    .field label { display: block; margin-bottom: 5px; font-weight: 500; }
    .field input {
      width: 100%;
      padding: 10px;
      border: 1px solid #ddd;
      border-radius: 4px;
      font-size: 16px;
      box-sizing: border-box;
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
    }
    button:hover { background: #0056b3; }
    .loading { padding: 20px; color: #666; }
  `;

  static properties = {
    streamUrl: { type: String },
    baseUrl: { type: String },
    loading: { state: true },
    surfaceId: { state: true },
    rootId: { state: true },
    components: { state: true },
    dataModel: { state: true },
  };

  streamUrl = '';
  baseUrl = 'http://localhost:8080';
  loading = true;
  surfaceId = '';
  rootId = 'root';
  components: Record<string, A2UIComponent> = {};
  dataModel: Record<string, any> = {};
  formData: Record<string, any> = {};

  connectedCallback() {
    super.connectedCallback();
    if (this.streamUrl) {
      this.loadSurface(this.streamUrl);
    }
  }

  async loadSurface(url: string) {
    this.loading = true;
    const response = await fetch(url);
    const text = await response.text();
    const messages: A2UIMessage[] = text
      .split('\n')
      .filter(line => line.trim())
      .map(line => JSON.parse(line));
    this.processMessages(messages);
    this.loading = false;
  }

  processMessages(messages: A2UIMessage[]) {
    const comps: Record<string, A2UIComponent> = {};
    let data: Record<string, any> = {};

    for (const msg of messages) {
      if (msg.beginRendering) {
        this.surfaceId = msg.beginRendering.surfaceId;
        this.rootId = msg.beginRendering.root;
      }
      if (msg.surfaceUpdate) {
        for (const comp of msg.surfaceUpdate.components) {
          comps[comp.id] = comp;
        }
      }
      if (msg.dataModelUpdate) {
        data = { ...data, ...msg.dataModelUpdate.contents };
      }
    }

    this.components = comps;
    this.dataModel = data;

    // Initialize form data
    for (const [path, value] of Object.entries(data)) {
      if (path.startsWith('/form/')) {
        this.formData[path] = value;
      }
    }
  }

  getData(path?: string): any {
    if (!path) return '';
    return this.formData[path] ?? this.dataModel[path] ?? '';
  }

  handleInput(path: string, e: Event) {
    this.formData[path] = (e.target as HTMLInputElement).value;
  }

  async handleAction(action: { type: string; data?: Record<string, any> }) {
    if (action.type === 'submit' && action.data?.endpoint) {
      const event = {
        event: {
          surfaceId: this.surfaceId,
          componentId: 'submit-btn',
          type: 'action',
          data: {
            name: this.formData['/form/name'] || '',
            date: this.formData['/form/date'] || '',
            time: this.formData['/form/time'] || '',
            party: parseInt(this.formData['/form/party']) || 2
          }
        }
      };

      this.loading = true;
      const endpoint = this.baseUrl + action.data.endpoint;
      const response = await fetch(endpoint, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(event)
      });
      const text = await response.text();
      const messages = text.split('\n').filter(l => l.trim()).map(l => JSON.parse(l));
      this.processMessages(messages);
      this.loading = false;
    }
    else if (action.type === 'navigate' && action.data?.url) {
      await this.loadSurface(this.baseUrl + action.data.url);
    }
  }

  renderComponent(id: string): any {
    const comp = this.components[id];
    if (!comp) return html``;

    if (comp.Column) {
      return html`
        <div class="column">
          ${comp.Column.children.map(childId => this.renderComponent(childId))}
        </div>
      `;
    }

    if (comp.Row) {
      return html`
        <div class="row">
          ${comp.Row.children.map(childId => this.renderComponent(childId))}
        </div>
      `;
    }

    if (comp.Card) {
      return html`
        <div class="card ${id.includes('success') ? 'success' : ''}">
          ${this.renderComponent(comp.Card.child)}
        </div>
      `;
    }

    if (comp.Text) {
      const text = comp.Text.text || this.getData(comp.Text.dataBinding?.path);
      const classes = `text ${id === 'header' ? 'header' : ''} ${id === 'title' ? 'title' : ''}`;
      return html`<div class="${classes}">${text}</div>`;
    }

    if (comp.TextField) {
      const path = comp.TextField.dataBinding?.path || '';
      return html`
        <div class="field">
          <label>${comp.TextField.label}</label>
          <input
            .value="${this.getData(path)}"
            placeholder="${comp.TextField.placeholder || ''}"
            @input="${(e: Event) => this.handleInput(path, e)}"
          />
        </div>
      `;
    }

    if (comp.Button) {
      return html`
        <button @click="${() => this.handleAction(comp.Button!.action)}">
          ${comp.Button.text}
        </button>
      `;
    }

    return html``;
  }

  render() {
    if (this.loading) {
      return html`<div class="loading">Loading...</div>`;
    }
    return this.renderComponent(this.rootId);
  }
}

customElements.define('a2ui-renderer', A2UIRenderer);
