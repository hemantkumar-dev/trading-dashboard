// simple hook to maintain websocket connection and incoming messages
import { useEffect, useRef, useState } from "react";

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";
const WS_URL = (import.meta.env.VITE_WS_URL) ? import.meta.env.VITE_WS_URL : API_URL.replace(/^http/, "ws") + "/ws";

export default function useWebSocket() {
  const wsRef = useRef(null);
  const [connected, setConnected] = useState(false);
  const [messages, setMessages] = useState([]); // latest Price messages as objects

  useEffect(() => {
    const ws = new WebSocket(WS_URL);
    wsRef.current = ws;

    ws.onopen = () => {
      setConnected(true);
      console.info("WS connected", WS_URL);
    };

    ws.onmessage = (ev) => {
      try {
        const payload = JSON.parse(ev.data);
        // payload is expected: { symbol, price }
        setMessages((m) => {
          // keep last 500 messages to avoid memory blow
          const next = [...m, payload];
          if (next.length > 500) next.shift();
          return next;
        });
      } catch (err) {
        console.warn("WS parse error", err);
      }
    };

    ws.onclose = () => setConnected(false);
    ws.onerror = (e) => console.warn("WS error", e);

    return () => {
      try { ws.close(); } catch (_) {}
    };
  }, []);

  return { connected, messages, ws: wsRef.current };
}
