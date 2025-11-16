import React from "react";
import { cancelOrder } from "../api/orders";

export default function OrdersTable({ orders = [], refresh, token }) {
  const handleCancel = async (orderId) => {
    if (!window.confirm("Cancel this order?")) return;
    const ok = await cancelOrder(orderId, token);
    if (ok) {
      alert("Order cancelled");
      refresh();
    } else {
      alert("Failed to cancel order");
    }
  };

  return (
    <div>
      <table>
        <thead><tr>
          <th>#</th><th>Symbol</th><th>Side</th><th>Qty</th><th>Price</th><th>Status</th><th>Action</th>
        </tr></thead>
        <tbody>
          {orders.length === 0 && <tr><td colSpan="7" className="small">No orders yet</td></tr>}
          {orders.map((o) => (
            <tr key={o.id}>
              <td>{o.id}</td>
              <td>{o.symbol}</td>
              <td>{o.side}</td>
              <td>{o.quantity}</td>
              <td>{Number(o.price).toFixed(2)}</td>
              <td>{o.status}</td>
              <td>
                {o.status === "open" && (
                  <button onClick={() => handleCancel(o.id)} style={{ fontSize: "12px", padding: "4px 8px" }}>
                    Cancel
                  </button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      <div style={{marginTop:8}}>
        <button onClick={refresh}>Refresh</button>
      </div>
    </div>
  );
}
