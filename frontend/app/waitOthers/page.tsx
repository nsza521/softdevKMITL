"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useSearchParams } from "next/navigation";
import styles from "./waitOthers.module.css";

export default function WaitOthers() {
  const searchParams = useSearchParams();
  const reservation_id = searchParams.get("reservationId")  || "";
  const [mode, setMode] = useState<1 | 2>(1); 
  const [confirmed_paid_people, setConfirmed_paid_people] = useState<number>(0)
  const [total_people, setTotal_people] = useState<number>(0)

  useEffect(() => {
    const token = localStorage.getItem("token")
    //polling ทุก 2 วินาที
    const interval = setInterval(async () => {
      try {
        const res = await fetch(`http://localhost:8080/table/reservation/${reservation_id}/status`, {
            headers: {
            "Authorization": `Bearer ${token}`,
            "Content-Type": "application/json"
        }});
        if (!res.ok) throw new Error("โหลดข้อมูลการจองไม่สำเร็จ")
        const data = await res.json();

        console.log(data)
        const reserve_status = data.status_detail.reservation_status
        setConfirmed_paid_people(data.status_detail.confirmed_paid_people)
        setTotal_people(data.status_detail.total_people)
        
        if (reserve_status === "paid") {
          const confirm = await fetch(`http://localhost:8080/table/reservation/${reservation_id}/confirm`, {
            method: "POST",
            headers: {
              "Authorization": `Bearer ${token}`,
              "Content-Type": "application/json"
            }});
          if (!confirm.ok) throw new Error("คอนเฟิร์มการจองไม่สำเร็จ")
              
          const confirm_resp = await confirm.json();
          console.log(confirm_resp)

          // const myprofile = await fetch("http://localhost:8080/customer/profile", {
          //   headers: {
          //     "Authorization": `Bearer ${token}`,
          //     "Content-Type": "application/json"
          //   }});
          // if(!myprofile) throw new Error("ดึงข้อมูลของฉันไม่สำเร็จ")
          

          // const myprofile_resp = await myprofile.json();
          // const my_usrname = myprofile_resp.username
        
          // noti part
          // const members = data.status_detail.members
          // // const targetMembers = members.slice(1); 
          // if(my_usrname == members[0]) {
          //   for (const member of members) {
          //     const noti = {
          //         event: "reserve_success",
          //         receiverUsername: member.username,
          //         receiverType: "customer",
          //         data: {
          //             // tableNo: table.row + table.col,
          //             // when: result.reservation.create_at,
          //             members: members.map((m: { username: string }) => m.username),
          //             reserveId: reservation_id,
          //         },
          //     };

          //     const notificationRes = await fetch("http://localhost:8080/notification/event", {
          //         method: "POST",
          //         headers: {  
          //             "Content-Type": "application/json",
          //             ...(token ? { Authorization: `Bearer ${token}` } : {}),
          //         },
          //         body: JSON.stringify(noti),
          //     });

          //     if (!notificationRes.ok) {
          //         console.error("Failed to send notification to", member.username);
          //     }

          //     console.log("Notification sent to", member.username);
          //     const notires = await notificationRes.json();
          //     console.log(notires)
          //   }
          // }

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
        <Mode1 confirmed_paid_people={confirmed_paid_people} total_people={total_people}/>
      ) : (
        <Mode2 />
      )}
    </div>
  );
}

function Mode1( { confirmed_paid_people, total_people }: { confirmed_paid_people: number, total_people: number } ) {
  return (
    <div className={styles.modeCon1}>
      <h2>ระบบกำลังรอสมาชิกท่านอื่นสั่งอาหาร</h2>
      <h2>{confirmed_paid_people}/{total_people}</h2>
    </div>
  );
}

function Mode2() {
  const router = useRouter();
  const [countdown, setCountdown] = useState(5); // จำนวนวินาทีเริ่มต้น

  useEffect(() => {
    if (countdown <= 0) {
      router.push("/home");    // ถึง 0 ให้ redirect
      return;
    }

    const timer = setTimeout(() => {
      setCountdown(prev => prev - 1);
    }, 1000);

    return () => clearTimeout(timer);
  }, [countdown, router]);

  return (
    <div className={styles.modeCon2}>
      <div>
        <h2>จองโต๊ะและสั่งอาหารสำเร็จ!</h2>
        <h2>ระบบจะทำการหักเงินในกระเป๋าอัตโนมัติ</h2>
      </div>

      <div className={styles.buttonCon}>
        <button className={styles.histBt} onClick={() => router.push("/history")}>
          ดูประวัติการจอง <img src="/Arrow_Right_MD.svg" />
        </button>

        <p>กำลังกลับไปที่หน้าหลักในอีก {countdown} วินาที</p>
      </div>
    </div>
  );
}