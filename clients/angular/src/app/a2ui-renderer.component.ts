import { Component, Input, OnInit, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { A2UIService, A2UIMessage, A2UIComponent, ClientEvent } from './a2ui.service';

@Component({
  selector: 'a2ui-renderer',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="a2ui-surface">
      @if (loading()) {
        <div class="loading">Loading...</div>
      } @else {
        <ng-container *ngTemplateOutlet="componentTpl; context: { id: rootId() }"></ng-container>
      }
    </div>

    <ng-template #componentTpl let-id="id">
      @if (getComponent(id); as comp) {
        <!-- Column -->
        @if (comp.Column) {
          <div class="a2ui-column">
            @for (childId of comp.Column.children; track childId) {
              <ng-container *ngTemplateOutlet="componentTpl; context: { id: childId }"></ng-container>
            }
          </div>
        }

        <!-- Row -->
        @if (comp.Row) {
          <div class="a2ui-row">
            @for (childId of comp.Row.children; track childId) {
              <ng-container *ngTemplateOutlet="componentTpl; context: { id: childId }"></ng-container>
            }
          </div>
        }

        <!-- Card -->
        @if (comp.Card) {
          <div class="a2ui-card" [class.success]="id.includes('success')">
            <ng-container *ngTemplateOutlet="componentTpl; context: { id: comp.Card.child }"></ng-container>
          </div>
        }

        <!-- Text -->
        @if (comp.Text) {
          <div class="a2ui-text"
               [class.header]="id === 'header'"
               [class.title]="id === 'title'">
            {{ comp.Text.text || getData(comp.Text.dataBinding?.path) }}
          </div>
        }

        <!-- TextField -->
        @if (comp.TextField) {
          <div class="a2ui-field">
            <label>{{ comp.TextField.label }}</label>
            <input
              [placeholder]="comp.TextField.placeholder || ''"
              [ngModel]="getData(comp.TextField.dataBinding?.path)"
              (ngModelChange)="setFormData(comp.TextField.dataBinding?.path, $event)"
            />
          </div>
        }

        <!-- Button -->
        @if (comp.Button) {
          <button class="a2ui-button" (click)="handleAction(comp.Button.action)">
            {{ comp.Button.text }}
          </button>
        }
      }
    </ng-template>
  `,
  styles: [`
    .a2ui-surface { font-family: system-ui, sans-serif; }
    .a2ui-column { display: flex; flex-direction: column; gap: 10px; }
    .a2ui-row { display: flex; flex-direction: row; gap: 10px; }
    .a2ui-card {
      background: white;
      border: 1px solid #ddd;
      border-radius: 8px;
      padding: 20px;
      box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }
    .a2ui-card.success { background: #d4edda; border-color: #c3e6cb; }
    .a2ui-text { margin: 5px 0; }
    .a2ui-text.header { font-size: 24px; font-weight: bold; margin-bottom: 15px; }
    .a2ui-text.title { font-size: 20px; color: #155724; }
    .a2ui-field { margin-bottom: 15px; }
    .a2ui-field label { display: block; margin-bottom: 5px; font-weight: 500; }
    .a2ui-field input {
      width: 100%;
      padding: 10px;
      border: 1px solid #ddd;
      border-radius: 4px;
      font-size: 16px;
      box-sizing: border-box;
    }
    .a2ui-button {
      background: #007bff;
      color: white;
      border: none;
      padding: 12px 24px;
      border-radius: 4px;
      cursor: pointer;
      font-size: 16px;
      width: 100%;
    }
    .a2ui-button:hover { background: #0056b3; }
    .loading { padding: 20px; color: #666; }
  `]
})
export class A2UIRendererComponent implements OnInit {
  @Input() streamUrl = '';
  @Input() baseUrl = 'http://localhost:8080';

  private a2uiService = new A2UIService();

  loading = signal(true);
  surfaceId = signal('');
  rootId = signal('root');
  components = signal<Record<string, A2UIComponent>>({});
  dataModel = signal<Record<string, any>>({});
  formData: Record<string, any> = {};

  async ngOnInit() {
    if (this.streamUrl) {
      await this.loadSurface(this.streamUrl);
    }
  }

  async loadSurface(url: string) {
    this.loading.set(true);
    const messages = await this.a2uiService.fetchSurface(url);
    this.processMessages(messages);
    this.loading.set(false);
  }

  processMessages(messages: A2UIMessage[]) {
    const comps: Record<string, A2UIComponent> = {};
    let data: Record<string, any> = {};

    for (const msg of messages) {
      if (msg.beginRendering) {
        this.surfaceId.set(msg.beginRendering.surfaceId);
        this.rootId.set(msg.beginRendering.root);
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

    this.components.set(comps);
    this.dataModel.set(data);

    // Initialize form data from data model
    for (const [path, value] of Object.entries(data)) {
      if (path.startsWith('/form/')) {
        this.formData[path] = value;
      }
    }
  }

  getComponent(id: string): A2UIComponent | undefined {
    return this.components()[id];
  }

  getData(path?: string): any {
    if (!path) return '';
    return this.formData[path] ?? this.dataModel()[path] ?? '';
  }

  setFormData(path: string | undefined, value: any) {
    if (path) {
      this.formData[path] = value;
    }
  }

  async handleAction(action: { type: string; data?: Record<string, any> }) {
    if (action.type === 'submit' && action.data?.['endpoint']) {
      const event: ClientEvent = {
        event: {
          surfaceId: this.surfaceId(),
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

      try {
        this.loading.set(true);
        const endpoint = this.baseUrl + action.data['endpoint'];
        const messages = await this.a2uiService.sendEvent(endpoint, event);
        this.processMessages(messages);
      } catch (error) {
        console.error('Submit failed:', error);
        alert('Submit failed - check console');
      } finally {
        this.loading.set(false);
      }
    }
    else if (action.type === 'navigate' && action.data?.['url']) {
      await this.loadSurface(this.baseUrl + action.data['url']);
    }
  }
}
