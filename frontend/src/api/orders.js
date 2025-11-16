const API_BASE = import.meta.env.VITE_API_URL || "http://localhost:8080";

export async function fetchOrders() {
  try {
    const r = await fetch(`${API_BASE}/orders`);
    if (!r.ok) return [];
    return await r.json();
  } catch (e) {
    console.error(e);
    return [];
  }
}

export async function placeOrder(order, token) {
  try {
    const headers = { "Content-Type": "application/json" };
    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }
    const r = await fetch(`${API_BASE}/orders`, {
      method: "POST",
      headers,
      body: JSON.stringify(order),
    });
    if (!r.ok) {
      const err = await r.json();
      console.error("Order error:", err);
      throw new Error(err.error || "Order placement failed");
    }
    return r.ok;
  } catch (e) {
    console.error("Order exception:", e);
    alert(`Order failed: ${e.message}`);
    return false;
  }
}

export async function login(username) {
  try {
    const r = await fetch(`${API_BASE}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username }),
    });
    if (!r.ok) return null;
    const data = await r.json();
    return data.token;
  } catch (e) {
    console.error(e);
    return null;
  }
}

export async function getHoldings(token) {
  try {
    const headers = {};
    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }
    const r = await fetch(`${API_BASE}/holdings`, {
      headers,
    });
    if (!r.ok) return [];
    return await r.json();
  } catch (e) {
    console.error(e);
    return [];
  }
}

export async function cancelOrder(orderId, token) {
  try {
    const headers = {};
    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }
    const r = await fetch(`${API_BASE}/orders/${orderId}/cancel`, {
      method: "POST",
      headers,
    });
    return r.ok;
  } catch (e) {
    console.error(e);
    return false;
  }
}

export async function fetchPrices() {
  try {
    const r = await fetch(`${API_BASE}/prices`);
    if (!r.ok) return {};
    return await r.json();
  } catch (e) {
    console.error(e);
    return {};
  }
}
