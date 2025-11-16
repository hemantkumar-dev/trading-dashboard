import React, { useEffect, useState } from "react";
import LivePrices from "./components/LivePrices";
import OrderForm from "./components/OrderForm";
import OrdersTable from "./components/OrdersTable";
import Holdings from "./components/Holdings";
import PriceChart from "./components/PriceChart";
import { fetchOrders, login } from "./api/orders";

export default function App() {
  const [orders, setOrders] = useState([]);
  const [selectedSymbol, setSelectedSymbol] = useState("AAPL");
  const [token, setToken] = useState(null);
  const [username, setUsername] = useState("");
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  // Load token from localStorage on mount
  useEffect(() => {
    const savedToken = localStorage.getItem("token");
    const savedUsername = localStorage.getItem("username");
    
    // Check if token is expired before using it
    if (savedToken) {
      try {
        // Decode JWT to check expiration
        const parts = savedToken.split('.');
        if (parts.length === 3) {
          const decoded = JSON.parse(atob(parts[1]));
          const expTime = decoded.exp * 1000; // Convert to milliseconds
          const now = Date.now();
          
          if (now < expTime) {
            // Token is still valid
            setToken(savedToken);
            setIsLoggedIn(true);
            setUsername(savedUsername || "User");
          } else {
            // Token is expired, clear it
            console.log("Stored token is expired, clearing localStorage");
            localStorage.clear();
          }
        }
      } catch (e) {
        console.error("Error checking token expiration:", e);
        localStorage.clear();
      }
    }
  }, []);

  const handleLogin = async (e) => {
    e.preventDefault();
    if (!username.trim()) {
      alert("Enter username");
      return;
    }
    const newToken = await login(username);
    if (newToken) {
      // Clear old localStorage data before storing new token
      localStorage.clear();
      setToken(newToken);
      setIsLoggedIn(true);
      localStorage.setItem("token", newToken);
      localStorage.setItem("username", username);
      await loadOrders();
    } else {
      alert("Login failed");
    }
  };

  const handleLogout = () => {
    setToken(null);
    setIsLoggedIn(false);
    setUsername("");
    localStorage.removeItem("token");
    localStorage.removeItem("username");
    setOrders([]);
  };

  const loadOrders = async () => {
    const data = await fetchOrders();
    setOrders(data || []);
  };

  useEffect(() => {
    if (isLoggedIn) {
      loadOrders();
    }
  }, [isLoggedIn]);

  const onOrderPlaced = async () => {
    await loadOrders();
  };

  if (!isLoggedIn) {
    return (
      <div className="app">
        <header>
          <h1>Trading Dashboard</h1>
        </header>
        <main style={{ display: "flex", justifyContent: "center", alignItems: "center", height: "60vh" }}>
          <div className="card" style={{ width: "400px" }}>
            <h2>Login</h2>
            <form onSubmit={handleLogin}>
              <label>Username</label>
              <input
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder="Enter your username"
              />
              <button type="submit" style={{ marginTop: "20px", width: "100%" }}>
                Login
              </button>
            </form>
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className="app">
      <header>
        <h1>Trading Dashboard</h1>
        <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginTop: "10px" }}>
          <span>Welcome, {username}!</span>
          <button onClick={handleLogout} style={{ padding: "8px 16px", cursor: "pointer" }}>
            Logout
          </button>
        </div>
      </header>

      <main>
        <section className="left">
          <div className="card">
            <h2>Live Prices</h2>
            <LivePrices onSymbolClick={(s) => setSelectedSymbol(s)} />
          </div>

          <div className="card">
            <h2>Place Order</h2>
            <OrderForm onPlaced={onOrderPlaced} token={token} />
          </div>
        </section>

        <section className="center">
          <div className="card">
            <h2>Price Chart — {selectedSymbol}</h2>
            <PriceChart symbol={selectedSymbol} />
          </div>

          <div className="card">
            <h2>Holdings</h2>
            <Holdings token={token} />
          </div>
        </section>

        <section className="right">
          <div className="card">
            <h2>Orders</h2>
            <OrdersTable orders={orders} refresh={loadOrders} token={token} />
          </div>
        </section>
      </main>

      <footer>
        <small>Built with Go (backend) + React (frontend)</small>
      </footer>
    </div>
  );
}
