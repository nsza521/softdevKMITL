"use client";


import styles from "./[id].module.css"
import { title } from "process";
import { useEffect, useState } from "react";
import styles from "./[id].module.css";
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
interface NotiCon {
  id : string;
  title : string;
  type : string;
  content : string;
  createdAt : string;
  attributes : NotiAttributes;
}

interface NotiAttributes {
  members : string[];
  tableNo : string;
  reserveId : string;
  when : string;
  queueNo? : string;
  restaurant? : string;
}

export default function NotificationDetailPage (){
  const params = useParams();
   const id = params.id as string;
  const [notiContent, setNotiContent] = useState<NotiCon | null>(null);

  const handleConfirm = () => {
    alert("คุณกดยืนยันเรียบร้อยแล้ว!");
    
  };

  useEffect(() =>{
    const fetchNotiContent = async () => {
      try{
        const token = localStorage.getItem("token");
        const res = await fetch("http://localhost:8080/notification/1",{
          headers: { Authorization: `Bearer ${token}` },
        });
        const data = await res.json();
         console.log("🔍 Data from API:", data);
        const found = data.items.find((item: NotiCon) => item.id === id);
        setNotiContent(found || null);
        console.log("reserveId:", found?.attributes.reserveId);
      }
      catch(err){
        console.error(err);
      }
    }
    fetchNotiContent();
  },[id])
  if (!notiContent) return <p>ไม่พบ notification</p>;
  return (
    <div>
      <div className={styles.container}>
        <div className={styles.content}>
          <div className={styles.header}>
            <h2>{notiContent.title}</h2>
          </div>

            {notiContent.type === "RESERVE_WITH" && (
              <div className={styles.detail}>
                <p>รายละเอียด :&nbsp;{notiContent.content}</p>
                <p>โต๊ะที่ {notiContent.attributes.tableNo}</p>
                <p>วันที่ {notiContent.attributes.when}</p>
                <div className={styles.member}>
                  <p>สมาชิก :&nbsp;</p>
                  <div>
                    {notiContent.attributes.members.map((member, index) => (
                      <p key={index}>{member}</p>
                    ))}
                  </div>
                </div>
                <p className={styles.descibe}>*  หากคุณได้ทำการจองโต๊ะร่วมกับรายชื่อดังกล่าว
                      โปรดยืนยันเพื่อดำเนินการต่อ</p>
                <div className={styles.confirmBtn}>
                  <button
            className={styles.acceptBtn}
            onClick={async () => {
              try {
                const token = localStorage.getItem("token");
                if (!token) {
                  alert("กรุณาเข้าสู่ระบบก่อน");
                  return;
                }
                const reserveId = notiContent.attributes.reserveId;
               
                const res = await fetch(
                  "http://localhost:8080/table/reservation/${reserveId}/confirm_member",
                  {
                    method: "POST",
                    headers: {
                      Authorization: `Bearer ${token}`,
                      "Content-Type": "application/json",
                    },
                  }
                );

                if (!res.ok) {
                  const err = await res.text();
                  throw new Error(err);
                }

                alert("ยืนยันการจองโต๊ะสำเร็จ!");
              } catch (error) {
                console.error(error);
                alert("เกิดข้อผิดพลาดในการยืนยันการจองโต๊ะ");
              }
            }}
          >
        ยืนยัน
      </button>
                  <button className={styles.cancleBtn}>ยกเลิก</button>
                </div>
              </div>
            )}

            {notiContent.type === "ORDER_FINISHED" && (
              <div className={styles.detail}>
                <p>รายละเอียด :&nbsp;{notiContent.content}</p>
                <p>โต๊ะที่ {notiContent.attributes.tableNo}</p>
                <p>วันที่ {notiContent.attributes.when}</p>
                <p>ร้านอาหาร : {notiContent.attributes.restaurant}</p>
                <p>คิวที่ {notiContent.attributes.queueNo}</p>
              </div>
            )}
            {notiContent.type === "ORDER_CANCELED" && (
              <div className={styles.detail}>
                <p>รายละเอียด :&nbsp;{notiContent.content}</p>
                <p>โต๊ะที่ {notiContent.attributes.tableNo}</p>
                <p>วันที่ {notiContent.attributes.when}</p>
                <p>ร้านอาหาร : {notiContent.attributes.restaurant}</p>
                <p>คิวที่ {notiContent.attributes.queueNo}</p>
                <p className={styles.descibe}>* คิวของคุณจะไม่ถูกเลื่อนออกไปแต่อาหารที่คุณเปลี่ยน
                    หากราคาแตกต่างเราจะทำการหักเงิน/คืนของคุณใน
                    ระบบ</p>
              </div>
            )}
            {notiContent.type === "RESERVE_SUCCESS" && (
               <div className={styles.detail}>
                <p>รายละเอียด :&nbsp;{notiContent.content}</p>
                <p>โต๊ะที่ {notiContent.attributes.tableNo}</p>
                <p>วันที่ {notiContent.attributes.when}</p>
                <p>ร้านอาหาร : {notiContent.attributes.restaurant}</p>
                <p>คิวที่ {notiContent.attributes.queueNo}</p>
                <div className={styles.member}>
                  <p>สมาชิก :&nbsp;</p>
                  <div>
                    {notiContent.attributes.members.map((member, index) => (
                      <p key={index}>{member}</p>
                    ))}
                  </div>
                </div>
              </div>
            )}
            {notiContent.type === "RESERVE_FAILED" && (
              <div className={styles.detail}>
                <p>รายละเอียด :&nbsp;{notiContent.content}</p>
                <p>โต๊ะที่ {notiContent.attributes.tableNo}</p>
                <p>วันที่ {notiContent.attributes.when}</p>
                <div className={styles.member}>
                  <p>สมาชิก :&nbsp;</p>
                  <div>
                    {notiContent.attributes.members.map((member, index) => (
                      <p key={index}>{member}</p>
                    ))}
                  </div>
                </div>
              </div>
            )}
        </div>
      </div>
    </div>
  );
}
