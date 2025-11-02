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
  const baseUrl = "http://localhost:8080";

  const [role, setRole] = useState<"customer" | "restaurant">("customer");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [email, setEmail] = useState("");
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [bankName, setBankName] = useState("");
  const [accountNo, setAccountNo] = useState("");
  const [accountName, setAccountName] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    const endpoint =
      role === "customer"
        ? `${baseUrl}/customer/register`
        : `${baseUrl}/restaurant/register`;

    const bodyData =
      role === "customer"
        ? {
            username,
            password,
            first_name: firstName,
            last_name: lastName,
            email,
          }
        : {
            username,
            email,
            password,
            bank_name: bankName,
            account_no: accountNo,
            account_name: accountName,
          };

    try {
      const res = await fetch(endpoint, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(bodyData),
      });

      if (!res.ok) throw new Error(await res.text());
      const data = await res.json();

      console.log("✅ Register success:", data);
      alert("สมัครสมาชิกสำเร็จ!");
      window.location.href = "/login";
    } catch (err: any) {
      console.error("❌ Error:", err);
      setError("สมัครไม่สำเร็จ: " + err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={`${styles.container} ${notoThai.variable}`}>
      <div className={styles.signupbox}>
        <form className={styles.signupinbox} onSubmit={handleRegister}>
          <h2>Sign Up</h2>

          {/* เลือกประเภทผู้ใช้ */}
          <div className={styles.signupinputbox}>
            <p>Register as</p>
            <select
              value={role}
              onChange={(e) => setRole(e.target.value as "customer" | "restaurant")}
              style={{ width: "100%", padding: "8px", borderRadius: "8px" }}
            >
              <option value="customer">Customer</option>
              <option value="restaurant">Restaurant</option>
            </select>
          </div>

          {/* ช่องกรอกที่เหมือนกัน */}
          <div className={styles.signupinputbox}>
            <p>Username</p>
            <input
              type="text"
              placeholder="Username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
          </div>

          <div className={styles.signupinputbox}>
            <p>Email</p>
            <input
              type="email"
              placeholder="Email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>

          <div className={styles.signupinputbox}>
            <p>Password</p>
            <input
              type="password"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>

          {/* ถ้าเป็น Customer */}
          {role === "customer" && (
            <>
              <div className={styles.signupinputbox}>
                <p>First Name</p>
                <input
                  type="text"
                  placeholder="First name"
                  value={firstName}
                  onChange={(e) => setFirstName(e.target.value)}
                />
              </div>

              <div className={styles.signupinputbox}>
                <p>Last Name</p>
                <input
                  type="text"
                  placeholder="Last name"
                  value={lastName}
                  onChange={(e) => setLastName(e.target.value)}
                />
              </div>
            </>
          )}

          {/* ถ้าเป็น Restaurant */}
          {role === "restaurant" && (
            <>
              <div className={styles.signupinputbox}>
                <p>Bank Name</p>
                <input
                  type="text"
                  placeholder="เช่น scb, kbank, bbl"
                  value={bankName}
                  onChange={(e) => setBankName(e.target.value)}
                />
              </div>

              <div className={styles.signupinputbox}>
                <p>Account No</p>
                <input
                  type="text"
                  placeholder="0123456789"
                  value={accountNo}
                  onChange={(e) => setAccountNo(e.target.value)}
                />
              </div>

              <div className={styles.signupinputbox}>
                <p>Account Name</p>
                <input
                  type="text"
                  placeholder="ชื่อเจ้าของบัญชี"
                  value={accountName}
                  onChange={(e) => setAccountName(e.target.value)}
                />
              </div>
            </>
          )}

          <button type="submit" className={styles.submitbtn} disabled={loading}>
            {loading ? "Signing up..." : "Sign Up"}
          </button>

          {error && <p style={{ color: "red" }}>{error}</p>}

          <div className={styles.doyouhaveacc}>
            <a href="/login">Already have an account? Log in</a>
          </div>
        </form>
      </div>

      <div className={styles.logo}>
        <img src="/jongtohlogo.png" alt="Logo" />
      </div>
    </div>
  );
}
