"use client";

import React from "react";
import { useRouter } from "next/navigation";
import styles from "./orderMenuChooseRes.module.css";
import { useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";  

interface allRestaurant {
  id: string;
  username: string;
  status: string;
  picture_url: string;
}

interface Detail {
  create_at : string;
  table_row : string;
  table_col : string;
  members: { username: string }[];
}

export default function OrderMenuChooseRes() {
  const searchParams = useSearchParams();
  const reservationId = searchParams.get("reservationId");
  const router = useRouter();
  const [restaurant, setRestaurant] = useState<allRestaurant[]>([]);
  const [detail, setDetail] = useState<Detail | null>(null);

  const [time, setTime] = useState("00:00");
  const [timeout, setTimeoutStatus] = useState(false);

  useEffect(() => {
    const fetchRestaurant = async () => {
      try {
        const token = localStorage.getItem("token");
        const res = await fetch("http://localhost:8080/restaurant/all", {
          headers: { Authorization: `Bearer ${token}` },
        });
        const data = await res.json();
        setRestaurant(data.restaurants);

        //fetch reserve detail
        const resDetail = await fetch(`http://localhost:8080/table/reservation/${reservationId}/detail`,{
          headers: { Authorization: `Bearer ${token}` },
        })
        const dataDetail = await resDetail.json();
        setDetail(dataDetail.reservation);

        const resTimer = await fetch(`http://localhost:8080/table/reservation/${reservationId}/time`,{
          headers: { Authorization: `Bearer ${token}` },
        })

        const dataTime = await resTimer.json();
        const timeRemaining = dataTime?.time_detail?.time_remaining ?? "00:00";
        const isTimeout = dataTime?.time_detail?.timeout ?? false;
        console.log(isTimeout)
        setTime(timeRemaining);
        setTimeoutStatus(isTimeout);
        
      } catch (err) {
        console.error(err);
      }
    };
    fetchRestaurant();
    const interval = setInterval(fetchRestaurant, 1000);
    return () => clearInterval(interval);
  }, [reservationId]);

  useEffect(() => {
      if (timeout) {
        (async () => {
          try {
            const token = localStorage.getItem("token");
            if (!token) return;
  
            const res = await fetch(
              `http://localhost:8080/table/reservation/${reservationId}`,
              {
                method: "DELETE",
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
  
            alert("หมดเวลา — ระบบได้ยกเลิกการจองแล้ว");
            router.push("/home");
          } catch (error) {
            console.error("Error deleting reservation:", error);
            alert("เกิดข้อผิดพลาดในการยกเลิกการจอง");
            router.push("/");
          }
        })();
      }
    }, [timeout, router, reservationId]);


  return (
    <div className={styles.container}>
      <div className={styles.timer}>
        <img src="Clock.svg" alt="" />
        <p>{time}</p>
      </div>
      {restaurant.map((Rstr) => (
        <div
          key={Rstr.id}
          className={styles.blog_item}
          onClick={() =>{
              if (Rstr.status === "open"){
                  router.push(`/orderMenuChooseMenu?id=${encodeURIComponent(Rstr.id)}&reservationId=${encodeURIComponent(reservationId || "")}`);
                }
            }
          }
          style={{ cursor: Rstr.status === "open" ? "pointer" : "not-allowed", opacity: Rstr.status === "open" ? 1 : 0.5 }}
        >
          <div className={styles.image}>
            <img src={Rstr.picture_url || "./Rectangle.svg"} alt="ResPicture" />
          </div>
          <div className={styles.content}>
            <h3>{Rstr.username}</h3>
            {Rstr.status === "closed" && (
                <div className={styles.closedStatus}>
                    <p>{Rstr.status}</p>
                </div>
            )}
            {Rstr.status === "open" && (
                <div className={styles.openStatus}>
                    <p>{Rstr.status}</p>
                </div>
            )}

          </div>
        </div>
      ))}
    </div>
  );
}
