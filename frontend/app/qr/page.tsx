"use client";

import { useEffect, useState } from "react";
import styles from "./qr.module.css";

export default function QRPage() {
  const [photo, setPhoto] = useState("");
  useEffect(() => {
    const fetchQR = async () => {
      try {
        const token = localStorage.getItem("token");
        if (!token) {
          alert("กรุณาเข้าสู่ระบบก่อน");
          return;
        }

        const res = await fetch("http://localhost:8080/customer/qrcode", {
          headers: { Authorization: `Bearer ${token}` },
        });

        if (!res.ok) {
          throw new Error("ไม่สามารถดึงข้อมูล QR ได้");
        }

        const data = await res.json();
        setPhoto(data.qrcode_url); 
        console.log(data.qrcode_url);
      } catch (err) {
        console.error(err);
      }
    };

    fetchQR();
  }, []); 

  return (
    <div className={styles.container}>
      <p>สแกนที่ทางเข้าเพื่อยืนยันตัว</p>
      <img src={photo} alt="QR Code" className={styles.qrImage} />
    </div>
  );
}
