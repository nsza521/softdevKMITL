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

export default function OrderMenuChooseRes() {
  const searchParams = useSearchParams();
  const reservationId = searchParams.get("reservationId");
  const router = useRouter();
  const [restaurant, setRestaurant] = useState<allRestaurant[]>([]);

  useEffect(() => {
    const fetchRestaurant = async () => {
      try {
        const token = localStorage.getItem("token");
        const res = await fetch("http://localhost:8080/restaurant/all", {
          headers: { Authorization: `Bearer ${token}` },
        });
        const data = await res.json();
        setRestaurant(data.restaurants);
      } catch (err) {
        console.error(err);
      }
    };
    fetchRestaurant();
  }, []);

  return (
    <div className={styles.container}>
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
