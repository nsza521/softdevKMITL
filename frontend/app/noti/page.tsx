"use client";

import styles from "./noti.module.css"
import TransacDetail from "@/components/TransacDetail";
import Link from "next/link";


const icon = {
    success : "/success.svg",
    food : "/food.svg",
    wallet : "/wallet.png",
    mail : "/mail.svg",
    order_cancle : "/orderchange.svg",
    create : "/create.svg",
    unsuccess : "/unsuccess.svg",
};

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



function getViewDetail(head:string): string {
  if (head === "จองโต๊ะสำเร็จ") return "กดยืนยัน";
  else if (head === "จองโต๊ะไม่สำเร็จ") return "ดูรายละเอียด";
  else if (head === "กำลังตรวจสอบ") return "รอการยืนยัน";
  else return "ไม่ทราบสถานะ";
}

export default function NotificationPage() {

    return(
        <div className={styles.container}>
            <div className={styles.content}>
                              {mockUsers.map((user) => (
                <Link key={user.id} href={`/noti/${user.id}`} className={styles.link}>
                  <TransacDetail
                    head={user.head}
                    date={user.date}
                    viewdetail={getViewDetail(user.head)}
                    imgsrc={user.imgsrc}
                  />
                </Link>
            ))}
                  
            </div>
            
        </div>
    );
}