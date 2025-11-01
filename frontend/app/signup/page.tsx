"use client";

import { useState } from "react";
import styles from "./signup.module.css";
import { Noto_Sans_Thai } from "next/font/google";

const notoThai = Noto_Sans_Thai({
  subsets: ["thai"],
  weight: ["400", "700"],
  variable: "--font-noto-thai",
});

export default function SignupPage() {
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const res = await fetch("http://localhost:8080/customer/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          username,
          password,
          first_name: firstName,
          last_name: lastName,
          email,
        }),
      });

      if (!res.ok) {
        const errMsg = await res.text();
        throw new Error(errMsg || "Register failed");
      }

      const data = await res.json();
      console.log("✅ Register success:", data);

      alert("สมัครสมาชิกเรียบร้อย!");
      window.location.href = "/login"; // สมัครเสร็จเด้งไป login
    } catch (err: any) {
      console.error("❌ Error:", err);
      setError("สมัครไม่สำเร็จ: " + err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={`${styles.container} ${notoThai.variable}`}>
      <div className={styles.logo}>
        <div><img src="/jongtohsmall.png" alt="" /></div>
      </div>

      <div className={styles.registerbox}>
        <form onSubmit={handleRegister}>
          <h2>Sign Up</h2>

          <input
            type="text"
            placeholder="Name"
            value={firstName}
            onChange={(e) => setFirstName(e.target.value)}
          />
          <input
            type="text"
            placeholder="Surname"
            value={lastName}
            onChange={(e) => setLastName(e.target.value)}
          />
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
          <input
            type="text"
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />

          <button type="submit" className={styles.submitbtn} disabled={loading}>
            {loading ? "Signing up..." : (
              <>
                <span>Sign Up</span>{" "}
                <span className="material-symbols-outlined">arrow_forward</span>
              </>
            )}
          </button>

          {error && <p style={{ color: "red" }}>{error}</p>}
        </form>
      </div>
    </div>
  );
}
