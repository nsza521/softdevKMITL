"use client"; // 👈 ถ้าใช้ Next.js 13+ (App Router) ต้องใส่

import { useState } from "react";
import styles from "./login.module.css";
import { Noto_Sans_Thai } from "next/font/google";

const notoThai = Noto_Sans_Thai({
  subsets: ["thai"],
  weight: ["400", "700"],
  variable: "--font-noto-thai",
});

export default function LoginPage() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleLogin = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault(); // กัน reload หน้า
    setError("");
    setLoading(true);

    try {
      const res = await fetch("http://localhost:8080/customer/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      });

      if (!res.ok) {
        throw new Error("Login failed");
      }

      const data = await res.json();
      console.log("✅ Login success:", data);

      // เก็บ token ลง localStorage ก็ได้
      localStorage.setItem("token", data.token);

      alert("Login Success!");
      // redirect ไปหน้าอื่นก็ได้ เช่น /dashboard
      window.location.href = "/home";

    } catch (err) {
      console.error("❌ Error:", err);
      setError("Username หรือ Password ไม่ถูกต้อง");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={`${styles.container} ${notoThai.variable}`}>
      <div className={styles.loginbox}>
        <form className={styles.logininbox} onSubmit={handleLogin}>
          <h2>Login</h2>

          <div className={styles.logininputbox}>
            <p>Username</p>
            <input
              type="text"
              placeholder="Username / Email"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
          </div>

          <div className={styles.logininputbox}>
            <p>Password</p>
            <input
              type="password"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </div>

          <button type="submit" className={styles.submitbtn} disabled={loading}>
            {loading ? "Logging in..." : <><span>Log In</span> <span className="material-symbols-outlined">arrow_forward</span></>}
          </button>

          {error && <p style={{ color: "red" }}>{error}</p>}

          <div className={styles.doyouhaveacc}>   
            <div><span><a href="">Forget Password ?</a></span></div>
          </div>
        </form>
      </div>
      <div className={styles.logo}>Logo Banner</div>
    </div>
  );
}
