import React, { useEffect, useRef, useState } from "react";
import useWebSocket from "../hooks/useWebSocket";
import { Line } from "react-chartjs-2";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
} from "chart.js";

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
);

/*
  PriceChart:
  - Keeps a small in-memory time-series for each symbol using WS messages
  - Displays last N points for selected symbol
*/

const MAX_POINTS = 60;

export default function PriceChart({ symbol = "AAPL" }) {
  const { messages } = useWebSocket();
  const [series, setSeries] = useState({}); // { symbol: [{t:Date, p:number}, ...] }
  const [chartData, setChartData] = useState({ labels: [], datasets: [] });

  // update series with incoming messages
  useEffect(() => {
    if (!messages || messages.length === 0) return;
    const last = messages[messages.length - 1];
    const now = new Date();
    setSeries(prev => {
      const arr = prev[last.symbol] ? [...prev[last.symbol]] : [];
      arr.push({ t: now, p: Number(last.price) });
      if (arr.length > MAX_POINTS) arr.shift();
      return { ...prev, [last.symbol]: arr };
    });
  }, [messages]);

  // update chart data when symbol or series changes
  useEffect(() => {
    const data = series[symbol] || [];
    setChartData({
      labels: data.map((_, idx) => idx),
      datasets: [
        {
          label: symbol,
          data: data.map(d => d.p),
          borderColor: '#10b981',
          backgroundColor: 'rgba(16, 185, 129, 0.1)',
          fill: true,
          tension: 0.2,
          pointRadius: 0,
          borderWidth: 2,
        }
      ]
    });
  }, [series, symbol]);

  return (
    <div style={{ height: 360, position: 'relative' }}>
      <Line 
        data={chartData}
        options={{
          animation: { duration: 0 },
          responsive: true,
          maintainAspectRatio: false,
          scales: {
            y: { 
              beginAtZero: false,
            }
          },
          plugins: { 
            legend: { display: false },
            tooltip: { enabled: true }
          }
        }}
      />
    </div>
  );
}
