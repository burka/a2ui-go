import { Injectable } from '@angular/core';

export interface A2UIMessage {
  beginRendering?: { surfaceId: string; root: string };
  surfaceUpdate?: { surfaceId: string; components: A2UIComponent[] };
  dataModelUpdate?: { surfaceId: string; contents: Record<string, any> };
}

export interface A2UIComponent {
  id: string;
  Column?: { children: string[] };
  Row?: { children: string[] };
  Card?: { child: string };
  Text?: { text?: string; dataBinding?: { path: string } };
  TextField?: { label?: string; placeholder?: string; dataBinding?: { path: string } };
  Button?: { text: string; action: { type: string; data?: Record<string, any> } };
}

export interface ClientEvent {
  event: {
    surfaceId: string;
    componentId: string;
    type: string;
    data?: Record<string, any>;
  };
}

@Injectable({ providedIn: 'root' })
export class A2UIService {

  async fetchSurface(url: string): Promise<A2UIMessage[]> {
    try {
      const response = await fetch(url);
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}`);
      }
      const text = await response.text();
      return text
        .split('\n')
        .filter(line => line.trim())
        .map(line => JSON.parse(line));
    } catch (error) {
      console.error('fetchSurface error:', error);
      throw error;
    }
  }

  async sendEvent(url: string, event: ClientEvent): Promise<A2UIMessage[]> {
    try {
      console.log('Sending event to:', url, event);
      const response = await fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(event)
      });
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}`);
      }
      const text = await response.text();
      console.log('Response:', text);
      return text
        .split('\n')
        .filter(line => line.trim())
        .map(line => JSON.parse(line));
    } catch (error) {
      console.error('sendEvent error:', error);
      throw error;
    }
  }
}
