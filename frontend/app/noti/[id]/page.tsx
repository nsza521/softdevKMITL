"use client";


import styles from "./[id].module.css"
import { title } from "process";
import { useParams } from "next/navigation";

const mockUsers = [
  {
    id: 1,
    head: "คุณได้จองโต๊ะร่วมกับ Username",
    date: "19 ส.ค. 2025",
    imgsrc: "/mail.svg",
  },
  {
    id: 2,
    head: "จองโต๊ะไม่สำเร็จ",
    date: "20 ส.ค. 2025",
    imgsrc: "/unsuccess.svg",
  },
   {
    id: 3,
    head: "อาหารพร้อมแล้ว! คุณสามารถรับอาหารได้ที่ร้านค้า",
    date: "20 ส.ค. 2025",
    imgsrc: "/food.svg",
  },
   {
    id: 4,
    head: "อาหารที่คุณสั่งถูกยกเลิก",
    date: "20 ส.ค. 2025",
    imgsrc: "/orderchange.svg",
  },
   {
    id: 5,
    head: "คุณสร้างคำสั่งการจองโต๊ะ",
    date: "20 ส.ค. 2025",
    imgsrc: "/create.svg",
  },
   {
    id: 6,
    head: "จองโต๊ะไม่สำเร็จ",
    date: "20 ส.ค. 2025",
    imgsrc: "/unsuccess.svg",
  },

  {
    id: 7,
    head: "จองโต๊ะไม่สำเร็จ",
    date: "20 ส.ค. 2025",
    imgsrc: "/unsuccess.svg",
  },
];

export default function NotificationDetailPage (){
  const params = useParams();
  const id = Number(params.id);

  const user = mockUsers.find(u => u.id === id);

  const handleConfirm = () => {
    alert("คุณกดยืนยันเรียบร้อยแล้ว!");
    // TODO: เรียก API เพื่อ update สถานะได้ที่นี่
  };

  if (!user) return <p>ไม่พบ notification</p>;
  return (
    <div>
      <div className={styles.container}>
        <div className={styles.content}>
          <div className={styles.header}>
            <h3>{user.head}</h3>
          </div>
          <div className={styles.detail}>
            <p>{user.date}</p>
          </div>
        </div>

         {user.head === "คุณได้จองโต๊ะร่วมกับ Username" &&(
            <div className={styles.confirmBtn}>
              <button  onClick={handleConfirm}>ยืนยัน</button>
              <button  onClick={handleConfirm}>ยกเลิก</button>
            </div>
          )}

          {user.head === "อาหารที่คุณสั่งถูกยกเลิก" && (
            <div className={styles.confirmBtn}>
              <button onClick={handleConfirm}>เปลี่ยนคำสั่งซื้อ</button>
            </div>
          )}

          {user.head === "จองโต๊ะไม่สำเร็จ" && (
            <div className={styles.confirmBtn}>
              <button onClick={handleConfirm}>จองใหม่อีกครั้ง</button>  
            </div>
          )}
      </div>
    </div>
  );
}
