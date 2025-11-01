"use client";

import styles from "./history.module.css";
import { useEffect, useState } from "react";
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
    wallet : "/wallet.svg",
};

interface Histroy {
    reservation_id : string;
    table_row : string;
    table_col : string;
    create_at : string;
}

interface Topup {
    transaction_id : string;
    payment_method : string;
    amount : number;
    created_at :string;
}

export default function HistoryPage(){
    const [active, setActive] = useState("จองโต๊ะ");
    const [history, setHistroy] = useState<Histroy[]>([]);
    const [topup, setTopup] = useState<Topup[]>([]);

    const [error, setError] = useState("");

    useEffect(() => {
        const fetchHistory = async () =>{
            try {
                const token = localStorage.getItem("token");
                const resHistory = await fetch("http://localhost:8080/table/reservation/history",{
                    headers:{"Authorization": `Bearer ${token}`,},
                })
                const dataHistory = await resHistory.json();
                setHistroy(dataHistory.reservations);

                const resTopup = await fetch("http://localhost:8080/payment/transaction/all",{
                    headers:{"Authorization": `Bearer ${token}`,},
                })
                const dataTopup = await resTopup.json();
                setTopup(dataTopup.transactions);
            }catch (err) {
                console.error(err);
                setError("โหลดข้อมูลไม่สำเร็จ");
            }
        }
        fetchHistory();
},[]);

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
                    {[...history].reverse().map((n) => (
                        <TransacDetail
                            key={n.reservation_id}  
                            head={"โต๊ะ " + n.table_row + n.table_col}
                            date={n.create_at}
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
                    {[...topup].reverse().map((n) => (
                        <TransacDetail
                            key={n.transaction_id}  
                            head="เติมเงิน"
                            detail={"ผ่าน " + n.payment_method}
                            date={n.created_at}
                            price={n.amount + " บาท"}
                            imgsrc={icon.wallet}
                        />
                        ))}
                </div>
            )}


        </div>
    )
}