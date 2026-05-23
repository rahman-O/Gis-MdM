/**
 * WebSocket Client Example
 * ملف مثال لاستخدام WebSocket مع TypeScript
 */

import * as React from 'react';


export interface WebSocketMessage {
  type: string;
  messageId?: string;
  timestamp: number;
  userId?: string;
  deviceId?: string;
  data?: any;
  error?: string;
  metadata?: Record<string, string>;
}

export interface ConnectionOptions {
  reconnect?: boolean;
  reconnectInterval?: number;
  maxReconnectAttempts?: number;
  pingInterval?: number;
}

export class MDMWebSocketClient {
  private url: string;
  private token: string;
  private ws: WebSocket | null = null;
  private listeners: Map<string, Set<(msg: WebSocketMessage) => void>> = new Map();
  private isConnecting = false;
  private isConnected = false;
  private reconnectAttempts = 0;
  private options: Required<ConnectionOptions>;
  private reconnectTimeout?: NodeJS.Timeout;
  private pingInterval?: NodeJS.Timeout;

  constructor(url: string, token: string, options: ConnectionOptions = {}) {
    this.url = url;
    this.token = token;
    this.options = {
      reconnect: options.reconnect ?? true,
      reconnectInterval: options.reconnectInterval ?? 3000,
      maxReconnectAttempts: options.maxReconnectAttempts ?? 5,
      pingInterval: options.pingInterval ?? 30000,
    };
  }

  /**
   * الاتصال بـ WebSocket | Connect to WebSocket
   */
  public connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      if (this.isConnected) {
        resolve();
        return;
      }

      if (this.isConnecting) {
        reject(new Error('Already connecting'));
        return;
      }

      this.isConnecting = true;

      try {
        const protocol = this.url.startsWith('https') ? 'wss' : 'ws';
        const baseUrl = this.url
          .replace('https://', '')
          .replace('http://', '')
          .replace(/\/$/, '');
        
        const wsUrl = `${protocol}://${baseUrl}/rest/ws/connect`;

        this.ws = new WebSocket(wsUrl);

        // إضافة رمز JWT إلى رأس الاتصال
        // Note: WebSocket لا يدعم رؤوس مخصصة، استخدم معاملات الاستعلام بدلاً منها
        const wsUrlWithToken = `${wsUrl}?token=${this.token}`;
        this.ws = new WebSocket(wsUrlWithToken);

        this.ws.onopen = () => {
          this.isConnected = true;
          this.isConnecting = false;
          this.reconnectAttempts = 0;
          console.log('✓ WebSocket connected');

          // بدء النبضات
          this.startPing();

          // إصدار حدث الاتصال
          this.emit({
            type: 'ws_connected',
            timestamp: Date.now(),
          });

          resolve();
        };

        this.ws.onmessage = (event) => {
          this.handleMessage(event.data);
        };

        this.ws.onerror = (error) => {
          console.error('✗ WebSocket error:', error);
          this.emit({
            type: 'ws_error',
            timestamp: Date.now(),
            error: error.toString(),
          });
        };

        this.ws.onclose = () => {
          this.isConnected = false;
          this.isConnecting = false;
          this.stopPing();
          console.log('✗ WebSocket disconnected');

          this.emit({
            type: 'ws_disconnected',
            timestamp: Date.now(),
          });

          // محاولة إعادة الاتصال إذا كان مفعل
          if (this.options.reconnect && this.reconnectAttempts < this.options.maxReconnectAttempts) {
            this.scheduleReconnect();
          }
        };
      } catch (error) {
        this.isConnecting = false;
        reject(error);
      }
    });
  }

  /**
   * فصل الاتصال | Disconnect from WebSocket
   */
  public disconnect(): void {
    this.stopPing();
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout);
    }
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    this.isConnected = false;
  }

  /**
   * إرسال رسالة | Send message
   */
  public send(message: WebSocketMessage): boolean {
    if (!this.isConnected || !this.ws) {
      console.warn('WebSocket is not connected');
      return false;
    }

    try {
      this.ws.send(JSON.stringify(message));
      return true;
    } catch (error) {
      console.error('Failed to send message:', error);
      return false;
    }
  }

  /**
   * الاستماع لنوع رسالة معينة | Listen for specific message type
   */
  public on(
    messageType: string,
    callback: (message: WebSocketMessage) => void
  ): () => void {
    if (!this.listeners.has(messageType)) {
      this.listeners.set(messageType, new Set());
    }

    this.listeners.get(messageType)!.add(callback);

    // إرجاع دالة لإلغاء الاشتراك
    return () => this.off(messageType, callback);
  }

  /**
   * إلغاء الاستماع | Stop listening
   */
  public off(messageType: string, callback: (message: WebSocketMessage) => void): void {
    const callbacks = this.listeners.get(messageType);
    if (callbacks) {
      callbacks.delete(callback);
    }
  }

  /**
   * إرسال أمر إلى جهاز | Send command to device
   */
  public sendCommand(
    deviceId: string,
    command: string,
    parameters?: Record<string, any>
  ): boolean {
    const message: WebSocketMessage = {
      type: 'command_send',
      timestamp: Date.now(),
      deviceId,
      data: {
        commandId: this.generateId(),
        command,
        deviceId,
        parameters: parameters || {},
        createdAt: new Date().toISOString(),
      },
    };

    return this.send(message);
  }

  /**
   * طلب تزامن البيانات | Request data sync
   */
  public requestSync(entityType: string, filters?: Record<string, any>): boolean {
    const message: WebSocketMessage = {
      type: 'sync_request',
      timestamp: Date.now(),
      data: {
        syncId: this.generateId(),
        entityType,
        operation: 'fetch',
        filters: filters || {},
      },
    };

    return this.send(message);
  }

  /**
   * الحصول على حالة الاتصال | Get connection status
   */
  public isConnectedStatus(): boolean {
    return this.isConnected;
  }

  /**
   * الحصول على معلومات إحصائية | Get statistics
   */
  public getStats(): {
    connected: boolean;
    reconnectAttempts: number;
    listeners: number;
  } {
    return {
      connected: this.isConnected,
      reconnectAttempts: this.reconnectAttempts,
      listeners: this.listeners.size,
    };
  }

  /**
   * معالجة الرسالة الواردة | Handle incoming message
   */
  private handleMessage(data: string): void {
    try {
      const message: WebSocketMessage = JSON.parse(data);

      // إصدار الحدث للمستمعين
      const callbacks = this.listeners.get(message.type);
      if (callbacks) {
        callbacks.forEach((callback) => {
          try {
            callback(message);
          } catch (error) {
            console.error(`Error in callback for ${message.type}:`, error);
          }
        });
      }

      // معالجة الرسائل الخاصة
      this.handleSpecialMessages(message);
    } catch (error) {
      console.error('Failed to parse message:', error);
    }
  }

  /**
   * معالجة الرسائل الخاصة | Handle special messages
   */
  private handleSpecialMessages(message: WebSocketMessage): void {
    switch (message.type) {
      case 'heartbeat':
        // إرسال استجابة النبضة
        this.send({
          type: 'heartbeat_ok',
          timestamp: Date.now(),
        });
        break;

      case 'error':
        console.error('Server error:', message.error);
        break;

      case 'command_failed':
        console.error('Command failed:', message.data);
        break;
    }
  }

  /**
   * بدء النبضات الدورية | Start periodic pings
   */
  private startPing(): void {
    this.pingInterval = setInterval(() => {
      if (this.isConnected) {
        this.send({
          type: 'heartbeat_ok',
          timestamp: Date.now(),
        });
      }
    }, this.options.pingInterval);
  }

  /**
   * إيقاف النبضات | Stop pings
   */
  private stopPing(): void {
    if (this.pingInterval) {
      clearInterval(this.pingInterval);
      this.pingInterval = undefined;
    }
  }

  /**
   * جدولة إعادة الاتصال | Schedule reconnection
   */
  private scheduleReconnect(): void {
    this.reconnectAttempts++;
    const delay = this.options.reconnectInterval * this.reconnectAttempts;

    console.log(
      `Attempting to reconnect in ${delay}ms (attempt ${this.reconnectAttempts}/${this.options.maxReconnectAttempts})`
    );

    this.reconnectTimeout = setTimeout(() => {
      this.connect().catch((error) => {
        console.error('Reconnection failed:', error);
      });
    }, delay);
  }

  /**
   * إصدار حدث مخصص | Emit custom event
   */
  private emit(message: WebSocketMessage): void {
    const callbacks = this.listeners.get(message.type);
    if (callbacks) {
      callbacks.forEach((callback) => {
        try {
          callback(message);
        } catch (error) {
          console.error(`Error in callback for ${message.type}:`, error);
        }
      });
    }
  }

  /**
   * توليد معرّف فريد | Generate unique ID
   */
  private generateId(): string {
    return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
  }
}

// ==================== أمثلة الاستخدام | Usage Examples ====================

// مثال 1: الاتصال الأساسي | Basic Connection
export function basicExample() {
  const client = new MDMWebSocketClient('http://localhost:8081', 'your-jwt-token');

  client.connect().then(() => {
    console.log('Connected!');

    // الاستماع لتحديثات الأجهزة
    client.on('device_status', (msg) => {
      console.log('Device status:', msg.data);
    });

    // الاستماع للإشعارات
    client.on('notification', (msg) => {
      console.log('Notification:', msg.data);
    });
  });
}

// مثال 2: إرسال الأوامر | Send Commands
export function sendCommandExample(client: MDMWebSocketClient) {
  // إرسال أمر تثبيت تطبيق
  client.sendCommand('DEVICE001', 'install_app', {
    appId: 'com.example.app',
    version: '2.0',
  });

  // الاستماع للاستجابة
  client.on('command_response', (msg) => {
    if (msg.data.status === 'success') {
      console.log('Command succeeded:', msg.data.result);
    } else {
      console.error('Command failed:', msg.data.error);
    }
  });
}

// مثال 3: مراقبة الأجهزة | Monitor Devices
export function monitorDevicesExample(client: MDMWebSocketClient) {
  // متابعة أجهزة محددة
  const monitoredDevices = new Set<string>();

  client.on('device_online', (msg) => {
    console.log(`✓ Device online: ${msg.deviceId}`);
    monitoredDevices.add(msg.deviceId!);
  });

  client.on('device_offline', (msg) => {
    console.log(`✗ Device offline: ${msg.deviceId}`);
    monitoredDevices.delete(msg.deviceId!);
  });

  client.on('device_status', (msg) => {
    const device = msg.data;
    console.log(`Device status: ${device.deviceId} - ${device.status}`, {
      battery: device.battery,
      signal: device.signal,
      cpu: `${device.cpuUsage}%`,
      memory: `${device.memoryUsage}%`,
      storage: `${device.storageUsage}%`,
    });
  });
}

// مثال 4: عمليات دفعية | Batch Operations
export function batchOperationExample(client: MDMWebSocketClient) {
  // مراقبة تقدم العملية الدفعية
  client.on('batch_progress', (msg) => {
    const op = msg.data;
    console.log(
      `Batch Progress: ${op.operationId}`,
      `${op.completedCount}/${op.deviceCount} (${op.percent}%)`
    );
  });

  client.on('batch_complete', (msg) => {
    const op = msg.data;
    console.log(`Batch Complete: ${op.operationId}`, {
      total: op.deviceCount,
      succeeded: op.completedCount,
      failed: op.failedCount,
    });
  });
}

// مثال 5: React Hook للـ WebSocket | React Hook
export function useWebSocket(url: string, token: string) {
  const [isConnected, setIsConnected] = React.useState(false);
  const clientRef = React.useRef<MDMWebSocketClient | null>(null);

  React.useEffect(() => {
    if (!clientRef.current) {
      clientRef.current = new MDMWebSocketClient(url, token);
    }

    const client = clientRef.current;

    client.on('ws_connected', () => setIsConnected(true));
    client.on('ws_disconnected', () => setIsConnected(false));

    client.connect().catch(console.error);

    return () => {
      client.disconnect();
    };
  }, [url, token]);

  return {
    client: clientRef.current!,
    isConnected,
  };
}
