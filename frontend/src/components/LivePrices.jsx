import React, { useEffect, useState } from "react";
import useWebSocket from "../hooks/useWebSocket";
import { fetchPrices } from "../api/orders"; // reuse

export default function LivePrices({ onSymbolClick }) {
  const { connected, messages } = useWebSocket();
  const [prices, setPrices] = useState({});
  const [prev, setPrev] = useState({});

  useEffect(() => {
    // load snapshot on mount
    (async () => {
      const p = await fetchPrices();
      setPrices(p);
    })();
  }, []);

  useEffect(() => {
    // apply incoming messages
    if (messages.length === 0) return;
    const last = messages[messages.length - 1];
    setPrev((s) => ({ ...s, [last.symbol]: prices[last.symbol] ?? last.price }));
    setPrices((s) => ({ ...s, [last.symbol]: last.price }));
    // eslint-disable-next-line
  }, [messages]);

  return (
    <div>
      <div className="small">WS: {connected ? "connected" : "disconnected"}</div>
      <table>
        <thead><tr><th>Symbol</th><th>Price</th></tr></thead>
        <tbody>
          {Object.entries(prices).map(([sym, price]) => {
            const pprev = prev[sym] ?? price;
            const up = price > pprev;
            return (
              <tr key={sym}>
                <td>
                  <span className="symbol-link" onClick={() => onSymbolClick?.(sym)}>{sym}</span>
                </td>
                <td className={up ? "price-up" : "price-down"}>{Number(price).toFixed(2)}</td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}
