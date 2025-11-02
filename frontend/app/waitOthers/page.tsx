"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useSearchParams } from "next/navigation";
import styles from "./waitOthers.module.css";

export default function WaitOthers() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const reservation_id = searchParams.get("reservationId")  || "";
  const [mode, setMode] = useState<1 | 2>(1); 

  useEffect(() => {
    //polling ทุก 2 วินาที
    const interval = setInterval(async () => {
      try {
        const res = await fetch(`http://localhost:8080/table/reservation/${reservation_id}/detail`);
        const data = await res.json();

        if (data.status === true) {
          setMode(2);
          clearInterval(interval); //หยุด polling ไม่จำเป็นต้องเรียกแล้ว
        }
      } catch (error) {
        console.error("Polling error:", error);
      }
    }, 2000);

    return () => clearInterval(interval); 
  }, []);

  return (
    <div className={styles.container}>
      {mode === 1 ? (
        <Mode1 />
      ) : (
        <Mode2 />
      )}
    </div>
  );
}

function Mode1() {
  return (
    <div className={styles.modeCon}>
      <h2>ระบบกำลังรอสมาชิกท่านอื่นสั่งอาหาร</h2>
      <p>ระบบกำลังตรวจสอบสถานะ...</p>
    </div>
  );
}

function Mode2() {
  return (
    <div className={styles.modeCon}>
      <h2>จองโต๊ะและสั่งอาหารสำเร็จ! ระบบจะทำการหักเงินในกระเป๋าอัตโนมัติ!</h2>
      <div>
        <button>
            ดูประวัติการจอง <img src=""/>
        </button>
        <p>กำลังกลับไปที่หน้าหลักในอีก  วินาที</p>
      </div>
    </div>
  );
}