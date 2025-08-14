"use client";

export default function Navbar() {
  return (
    <nav style={{ padding: 12, background: "#0070f3", color: "white" }}>
      <span>MyApp</span>
      <a href="/" style={{ marginLeft: 20 }}>Home</a>
      <a href="/about" style={{ marginLeft: 10 }}>About</a>
    </nav>
  );
}
