"use client";

import styles from "./noti.module.css"
import TransacDetail from "@/components/TransacDetail";
import Link from "next/link";
import { useEffect, useState } from "react";


const icon = {
    success : "/success.svg",
    food : "/food.svg",
    wallet : "/wallet.png",
    mail : "/mail.svg",
    order_cancle : "/orderchange.svg",
    create : "/create.svg",
    unsuccess : "/unsuccess.svg",
};

interface allNoti {
  id : string;
  content : string;
  type : string;
  isRead : boolean;
  createdAt : string;
}

function getType(type:string) {
  if (type === "PAYMENT") return icon.wallet;
  else if (type === "RESERVE_SUCCESS") return icon.success;
  else if (type === "RESERVE_FAILED") return icon.unsuccess;
  else if (type === "ORDER_CANCELED") return icon.order_cancle;
  else if (type === "ORDER_FINISHED") return icon.food;
  else if (type === "BOOKING") return icon.create
  else return icon.mail
}



export default function NotificationPage() {
  const [noti, setNoti] = useState<allNoti[]>([]);

  useEffect(() => {
    const fetchNoti = async() => {
      try{
        const token = localStorage.getItem("token");
        const res = await fetch("http://localhost:8080/notification/1",{
          headers: { Authorization: `Bearer ${token}` },
        });

        const data = await res.json();
        setNoti(data.items);
        console.log(data)
      }
      catch(err){
        console.error(err);
      }
    }
    fetchNoti();
  },[])
    return(
        <div className={styles.container}>
            <div className={styles.content}>
                              {noti.map((user) => (
                <Link key={user.id} href={`/noti/${user.id}`} className={styles.link}>
                  <TransacDetail
                    head={user.content}
                    date={user.createdAt}
                    viewdetail={'ดูรายละเอียด'}
                    imgsrc={getType(user.type)}
                  />
                </Link>
            ))}
                  
            </div>
            
        </div>
    );
}