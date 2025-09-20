"use client";

import styles from "./history.module.css";
import { useState } from "react";
import { Noto_Sans_Thai } from "next/font/google";
import TransacDetail from "@/components/TransacDetail";


const notoThai = Noto_Sans_Thai({
  subsets: ["thai"],
  weight: ["400", "700"],
  variable: "--font-noto-thai",
});

const icon = {
    success : "/success.svg",
    food : "/food.svg",
    wallet : "wallet.png",
    mail : "",
    order_cancle : "",
    create : "",
    unsuccess : "",
};




export default function TransactionPage(){
    const [active, setActive] = useState("จองโต๊ะ");
    return(
        <div className={`${styles.content} ${notoThai.variable}`}>
            <div className={styles.catagories}>
                <button className={active === "จองโต๊ะ" ? styles.active : ""} onClick={()=> setActive("จองโต๊ะ")}>
                    จองโต๊ะ
                </button>
                <button className={active === "อาหาร" ? styles.active : ""} onClick={()=> setActive("อาหาร")}>
                    อาหาร
                </button>
                <button className={active === "การเติมเงิน" ? styles.active : ""} onClick={()=> setActive("การเติมเงิน")}>
                    การเติมเงิน
                </button>
            </div>

            {active === "จองโต๊ะ" &&(
                <div className={styles.detail_container}>
                    {Array.from({ length: 10 }).map((_, index) => (
                        <TransacDetail
                            key={index} // key ต้องไม่ซ้ำ
                            head={`โต๊ะ 3`}
                            date="19 ส.ค. 2025"
                            imgsrc={icon.success}
                        />
                        ))}
                </div>
            )}

            {active === "อาหาร" &&(
                <div className={styles.detail_container}>
                    {Array.from({ length: 10 }).map((_, index) => (
                        <TransacDetail
                            key={index} // key ต้องไม่ซ้ำ
                            head="กะเพราหมูกรอบ"
                            detail="จำนวน 1"
                            date="19 ส.ค. 2025"
                            price="100 บาท"
                            imgsrc={icon.food}
                        />
                        ))}
                </div>
            )}

            {active === "การเติมเงิน" &&(
                <div className={styles.detail_container}>
                    {Array.from({ length: 10 }).map((_, index) => (
                        <TransacDetail
                            key={index} // key ต้องไม่ซ้ำ
                            head="เติมเงิน"
                            detail="ผ่าน QR"
                            date="19 ส.ค. 2025"
                            price="100 บาท"
                            imgsrc={icon.wallet}
                        />
                        ))}
                </div>
            )}


        </div>
    )
}