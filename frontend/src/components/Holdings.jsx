import React, { useEffect, useState } from "react";
import { getHoldings } from "../api/orders";

export default function Holdings({ token }) {
  const [holdings, setHoldings] = useState([]);

  useEffect(() => {
    if (token) {
      (async () => {
        const data = await getHoldings(token);
        setHoldings(Array.isArray(data) ? data : []);
      })();
    }
  }, [token]);

  return (
    <div>
      <table>
        <thead><tr><th>Symbol</th><th>Qty</th><th>Avg Cost</th></tr></thead>
        <tbody>
          {holdings.length === 0 && <tr><td colSpan="3" className="small">No holdings</td></tr>}
          {holdings.map((r) => (
            <tr key={r.symbol}>
              <td>{r.symbol}</td>
              <td>{r.net_qty}</td>
              <td>{r.avg_cost ? Number(r.avg_cost).toFixed(2) : "-"}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
