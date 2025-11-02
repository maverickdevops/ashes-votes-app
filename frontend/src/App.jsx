import { useEffect, useState } from "react";

function App() {
  const API = "http://backend:8080";

  const [counts, setCounts] = useState({ australia: 0, england: 0 });
  const [loading, setLoading] = useState(false);

  async function fetchCounts() {
    try {
      const res = await fetch(`${API}/counts`);
      const data = await res.json();
      const map = { australia: 0, england: 0 };
      data.forEach(d => (map[d.team] = d.count));
      setCounts(map);
    } catch (e) {
      console.error("fetch counts", e);
    }
  }

  useEffect(() => {
    fetchCounts();
    const t = setInterval(fetchCounts, 1000);
    return () => clearInterval(t);
  }, []);

  async function castVote(team) {
    setLoading(true);
    try {
      await fetch(`${API}/vote`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ team }),
      });
      await fetchCounts();
    } catch (e) {
      console.error("vote error", e);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div style={{ fontFamily: "system-ui, sans-serif", maxWidth: 700, margin: "2rem auto", padding: "1rem" }}>
      <h1>Ashes – Vote for the winner</h1>
      <div style={{ display: "flex", gap: 20 }}>
        <div style={{ flex: 1, border: "1px solid #ddd", padding: 20, borderRadius: 8 }}>
          <h2>Australia</h2>
          <button onClick={() => castVote("australia")} disabled={loading}>Vote Australia</button>
          <p>Votes: <strong>{counts.australia}</strong></p>
        </div>
        <div style={{ flex: 1, border: "1px solid #ddd", padding: 20, borderRadius: 8 }}>
          <h2>England</h2>
          <button onClick={() => castVote("england")} disabled={loading}>Vote England</button>
          <p>Votes: <strong>{counts.england}</strong></p>
        </div>
      </div>

      <div style={{ marginTop: 20 }}>
        <h3>Live counts</h3>
        <div style={{ display: "flex", gap: 12 }}>
          <div style={{ padding: 10, borderRadius: 8, width: 140, textAlign: "center" }}>
            <div>Australia</div>
            <div style={{ fontSize: 24 }}>{counts.australia}</div>
          </div>
          <div style={{ padding: 10, borderRadius: 8, width: 140, textAlign: "center" }}>
            <div>England</div>
            <div style={{ fontSize: 24 }}>{counts.england}</div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
