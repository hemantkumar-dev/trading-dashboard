import React, { useEffect, useState } from "react";
import { placeOrder, fetchPrices } from "../api/orders";

const SYMBOLS = ["AAPL", "TSLA", "AMZN", "INFY", "TCS"];

export default function OrderForm({ onPlaced, token }) {
  const [form, setForm] = useState({
    symbol: "AAPL",
    side: "buy",
    quantity: 1,
    price: "",
  });

  useEffect(() => {
    (async () => {
      const p = await fetchPrices();
      if (!form.price) setForm((f) => ({ ...f, price: p[form.symbol] ?? "" }));
    })();
    // eslint-disable-next-line
  }, []);

  const submit = async (e) => {
    e?.preventDefault();
    // basic validation
    if (!form.symbol || !form.side || !form.quantity || !form.price) {
      alert("fill all fields");
      return;
    }

    const ok = await placeOrder(
      {
        symbol: form.symbol,
        side: form.side,
        quantity: Number(form.quantity),
        price: Number(form.price),
      },
      token
    );
    if (ok) {
      alert("Order placed");
      setForm((s) => ({ ...s, quantity: 1 }));
      onPlaced?.();
    } else {
      alert("Failed to place order");
    }
  };

  const onSymbolChange = async (s) => {
    setForm((f) => ({ ...f, symbol: s }));
    const p = await fetchPrices();
    setForm((f) => ({ ...f, price: p[s] ?? f.price }));
  };

  return (
    <form onSubmit={submit}>
      <label>Symbol</label>
      <select value={form.symbol} onChange={(e) => onSymbolChange(e.target.value)}>
        {SYMBOLS.map((s) => <option key={s} value={s}>{s}</option>)}
      </select>

      <label>Side</label>
      <select value={form.side} onChange={(e) => setForm({ ...form, side: e.target.value })}>
        <option value="buy">Buy</option>
        <option value="sell">Sell</option>
      </select>

      <label>Quantity</label>
      <input type="number" min="1" value={form.quantity} onChange={(e) => setForm({ ...form, quantity: Number(e.target.value) })} />

      <label>Price</label>
      <input type="number" step="0.01" value={form.price} onChange={(e) => setForm({ ...form, price: e.target.value })} />

      <div className="row">
        <button type="submit">Place Order</button>
        <button type="button" onClick={() => { setForm({ symbol: "AAPL", side: "buy", quantity: 1, price: "" }); }}>Reset</button>
      </div>
    </form>
  );
}
