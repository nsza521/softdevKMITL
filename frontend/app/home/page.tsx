"use client";

import { useState , useEffect } from "react";
import styles from "./home.module.css";
import { Noto_Sans_Thai } from "next/font/google";
import { useRouter } from "next/navigation"; 

const notoThai = Noto_Sans_Thai({
  subsets: ["thai"],
  weight: ["400", "700"],
  variable: "--font-noto-thai",
});

export default function HomePage() {
  const [showPopup, setShowPopup] = useState(false);
  const [profile, setProfile] = useState<any>(null);
   const router = useRouter(); 

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const token = localStorage.getItem("token");
        if (!token) return;

        const res = await fetch("http://localhost:8080/customer/profile", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (!res.ok) throw new Error("Failed to fetch profile");

        const data = await res.json();
        console.log("📌 Profile data:", data); // <<
        setProfile(data);
      } catch (err) {
        console.error("❌ Fetch profile error:", err);
      }
    };

    fetchProfile();
  }, []);

  const handleLogout = async () => {
    try {
      const token = localStorage.getItem("token"); // ดึง token ที่เก็บไว้ตอน login

      const res = await fetch("http://localhost:8080/customer/logout", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`, // ถ้า backend ต้องการ
        },
      });

      if (!res.ok) {
        throw new Error("Logout failed");
      }

      // เคลียร์ token ทิ้ง
      localStorage.removeItem("token");

      alert("ออกจากระบบเรียบร้อย");
      window.location.href = "/login"; // redirect กลับไปหน้า login

    } catch (err) {
      console.error("❌ Error:", err);
      alert("เกิดข้อผิดพลาดตอนออกจากระบบ");
    }
  };

  return (
    <div className={`${styles.container} ${notoThai.variable}`}>
      <div className={styles.headername}>
        <span className={styles.headernameedit}>
          สวัสดี {profile ? `${profile.first_name} ${profile.last_name}` : "กำลังโหลด..."}
          <button onClick={() => setShowPopup(true)}>
            <img src="/editpencil.svg" width={25} height={25} />
          </button>
        </span>
        <button className={styles.logoutbtn} onClick={handleLogout}>
          <img src="/logout.svg" width={25} height={25} />
          logout
        </button>
      </div>

      <div className={styles.boxs}>
        <span>ยอดเงินคงเหลือ {profile ? `${profile.wallet_balance}` : "กำลังโหลด..."} บาท</span>
        <button onClick={() => router.push("/topup")}>
          <img src="/plus.svg" width={15} height={15} />
          เติมเงิน
        </button>
      </div>

      <div className={styles.boxs}>
        <span className={styles.boxspan}>
          <img src="/qr.svg" width={20} height={20} />
          คิวอาโค้ดของฉัน
        </span>
        <button>
          <img src="/show.svg" width={25} height={25} />
          ดู
        </button>
      </div>

      <div className={styles.headername}>โต๊ะตอนี้ </div>
      <div className={styles.table}></div>
      <button className={styles.tablebtn}>จองโต๊ะ</button>

      {/* Popup */}
      {showPopup && (
        <div className={styles.popupbg} onClick={() => setShowPopup(false)}>
          <div className={styles.popup} onClick={(e) => e.stopPropagation()}>
            <h2>แก้ไขข้อมูลส่วนตัว</h2>
            <form
              className={styles.form}
              onSubmit={(e) => {
                e.preventDefault();
                setShowPopup(false);
              }}
            >
              <div className={styles.formGroup}>
                <label>Name</label>
                <input type="text" name="name" placeholder="กรอกชื่อ" />
              </div>

              <div className={styles.formGroup}>
                <label>Surname</label>
                <input type="text" name="surname" placeholder="กรอกนามสกุล" />
              </div>

              <div className={styles.formGroup}>
                <label>Username</label>
                <input type="text" name="username" placeholder="กรอกชื่อผู้ใช้" />
              </div>

              <div className={styles.formGroup}>
                <label>Password</label>
                <input type="password" name="password" placeholder="กรอกรหัสผ่าน" />
              </div>

              <div className={styles.buttonGroup}>
                <button
                  type="button"
                  className={styles.cancelBtn}
                  onClick={() => setShowPopup(false)}
                >
                  ยกเลิก
                </button>
                <button type="submit" className={styles.submitBtn}>
                  ยืนยัน
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
