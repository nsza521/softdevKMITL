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
    type : string;
    created_at :string;
}

interface OrderHistory {
    orders : Order[];
}

interface Order {
    order_id :string;
    channel : string;
    order_time : string;
    items : Item[];
}

interface Item {
    menu_name :string;
    quantity : number;
    subtotal : number;
}

export default function HistoryPage(){
    const [active, setActive] = useState("จองโต๊ะ");
    const [history, setHistroy] = useState<Histroy[]>([]);
    const [topup, setTopup] = useState<Topup[]>([]);
    const [order, setOrder] = useState<OrderHistory[]>([]);

    const [error, setError] = useState("");

    useEffect(() => {
        const fetchHistory = async () =>{
            try {
                const token = localStorage.getItem("token");
                // fetch Table reservation
                const resHistory = await fetch("http://localhost:8080/table/reservation/history",{
                    headers:{"Authorization": `Bearer ${token}`,},
                })
                const dataHistory = await resHistory.json();
                setHistroy(dataHistory.reservations);

                // fetch Top up 
                const resTopup = await fetch("http://localhost:8080/payment/transaction/all",{
                    headers:{"Authorization": `Bearer ${token}`,},
                })
                const dataTopup = await resTopup.json();
                setTopup(dataTopup.transactions ?? []);

                // fetch Order
                const resOrder = await fetch("http://localhost:8080/customer/orders/history",{
                    headers:{"Authorization": `Bearer ${token}`,},
                })
                const dataOrder = await resOrder.json();
                setOrder(dataOrder.history);
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

            {active === "อาหาร" && (
                <div className={styles.detail_container}>
                    {[...order].reverse().map((n) =>
                    n.orders.map((orderItem) =>
                        orderItem.items.map((item) => (
                        <TransacDetail
                            key={item.menu_name}
                            head={item.menu_name}
                            detail= {`จำนวน ${item.quantity}`}
                            date={orderItem.order_time}
                            price={`${item.subtotal} บาท`}
                            imgsrc={icon.food}
                        />
                        ))
                    )
                    )}
                </div>
            )}

            {active === "การเติมเงิน" &&(
                <div className={styles.detail_container}>
                    {[...topup].filter((n) => n.type === "topup").reverse().map((n) => (
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