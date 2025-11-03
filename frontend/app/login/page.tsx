"use client";

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
  const [userType, setUserType] = useState<"customer" | "restaurant">("customer");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleLogin = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      // 🟢 เลือก endpoint ตามประเภทผู้ใช้
      const url =
        userType === "restaurant"
          ? "http://localhost:8080/restaurant/login"
          : "http://localhost:8080/customer/login";

      const res = await fetch(url, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      });

      if (!res.ok) throw new Error("Login failed");

      const data = await res.json();
      console.log("✅ Login success:", data);

      localStorage.setItem("token", data.token);
      localStorage.setItem("userType", userType);

      alert("Login Success!");

      // 🔀 Redirect ตามประเภทผู้ใช้
      if (userType === "restaurant") {
        window.location.href = "/restaurant";
      } else {
        window.location.href = "/home";
      }
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

          <div className={styles.roleSection}>
            {/* <p>เข้าสู่ระบบในฐานะ</p> */}

            <div className={styles.roleToggle}>
              <button
                type="button"
                className={`${styles.roleBtn} ${userType === "customer" ? styles.activeRole : ""
                  }`}
                onClick={() => setUserType("customer")}
              >
                
                <p className={styles.roleSub1}>ลูกค้า</p>
                <p className={styles.roleSub}>Customer</p>
              </button>

              <button
                type="button"
                className={`${styles.roleBtn} ${userType === "restaurant" ? styles.activeRole : ""
                  }`}
                onClick={() => setUserType("restaurant")}
              >
                <p className={styles.roleSub1}>ร้านอาหาร</p>
                <p className={styles.roleSub}>Restaurant</p>
              </button>

            </div>
          </div>


          <button type="submit" className={styles.submitbtn} disabled={loading}>
            {loading ? (
              "Logging in..."
            ) : (
              <>
                <span>Log In</span>{" "}
                <span className="material-symbols-outlined">arrow_forward</span>
              </>
            )}
          </button>

          {error && <p style={{ color: "red" }}>{error}</p>}

          <div className={styles.doyouhaveacc}>
            <div>
              <span>
                <a href="/signup">No account? Sign up</a>
              </span>
            </div>
          </div>
        </form>
      </div>
      <div className={styles.logo}><img src="/jongtohlogo.png" alt="" /></div>
    </div>
  );
}
