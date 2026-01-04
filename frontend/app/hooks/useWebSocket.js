'use client';

import { useEffect, useRef, useState, useCallback } from 'react';

const WS_URL = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080/ws';

export function useWebSocket(onMessage) {
    const [connected, setConnected] = useState(false);
    const [wsInstance, setWsInstance] = useState(null);
    const wsRef = useRef(null);
    const reconnectTimeoutRef = useRef(null);
    const onMessageRef = useRef(onMessage);

    // Keep onMessage ref updated
    useEffect(() => {
        onMessageRef.current = onMessage;
    }, [onMessage]);

    const connect = useCallback(() => {
        if (wsRef.current?.readyState === WebSocket.OPEN) return;

        try {
            console.log('[WS] Connecting to', WS_URL);
            const ws = new WebSocket(WS_URL);

            ws.onopen = () => {
                console.log('[WS] Connected');
                wsRef.current = ws;
                setWsInstance(ws);
                setConnected(true);
            };

            ws.onmessage = (event) => {
                try {
                    // Handle multiple messages split by newline
                    const messages = event.data.split('\n');
                    messages.forEach(msgStr => {
                        if (!msgStr.trim()) return;
                        const message = JSON.parse(msgStr);
                        console.log('[WS] Received:', message.type);
                        if (onMessageRef.current) {
                            onMessageRef.current(message);
                        }
                    });
                } catch (err) {
                    console.error('[WS] Failed to parse message:', err);
                }
            };

            ws.onclose = () => {
                console.log('[WS] Disconnected');
                setConnected(false);
                setWsInstance(null);
                wsRef.current = null;

                // Reconnect after 3 seconds
                reconnectTimeoutRef.current = setTimeout(() => {
                    connect();
                }, 3000);
            };

            ws.onerror = (error) => {
                console.warn('[WS] Connection error - will retry');
            };

            wsRef.current = ws;
        } catch (err) {
            console.error('[WS] Failed to connect:', err);
        }
    }, []);

    useEffect(() => {
        connect();

        return () => {
            if (reconnectTimeoutRef.current) {
                clearTimeout(reconnectTimeoutRef.current);
            }
            if (wsRef.current) {
                wsRef.current.close();
            }
        };
    }, [connect]);

    // Helper to send messages
    const send = useCallback((data) => {
        if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
            const message = typeof data === 'string' ? data : JSON.stringify(data);
            console.log('[WS] Sending:', typeof data === 'object' ? data.type : 'string message');
            wsRef.current.send(message);
            return true;
        }
        console.warn('[WS] Cannot send - not connected');
        return false;
    }, []);

    return { connected, ws: wsInstance, send };
}
