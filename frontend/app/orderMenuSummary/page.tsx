'use client'

import { useState, useEffect } from "react"
import { useSearchParams } from "next/navigation";
import styles from "./orderMenuSummary.module.css"

export default function OrderMenuSummaryPage({ params }: any) {
    const searchParams = useSearchParams();
    const order_id = searchParams.get("order_id") || ""
    const reservation_id = searchParams.get("reservationId") || ""

    const [order, setOrder] = useState<any>(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)

    useEffect(() => {
        const fetchOrder = async () => {
            const token = localStorage.getItem("token")
            try {
                const res = await fetch(`http://localhost:8080/restaurant/order/${order_id}/detail`, {
                    headers: {
                    "Authorization": `Bearer ${token}`,
                    "Content-Type": "application/json"
                }});
                if (!res.ok) throw new Error("โหลดข้อมูลออเดอร์ไม่สำเร็จ")
                
                const data = await res.json()

                console.log("Order Data:", data)


                // setOrder(data)
            } catch (err: any) {
                setError(err.message)
            } finally {
                setLoading(false)
            }
        }

        fetchOrder()
    }, [order_id])

    if (loading) return <p>กำลังโหลดข้อมูล...</p>
    if (error) return <p style={{color:"red"}}>{error}</p>
    if (!order) return <p>ไม่พบข้อมูลออเดอร์</p>


    return (
        <div className={styles.container}>
            <div className={styles.myOrder}>
                <h2>My Order</h2>

                {order.items.length === 0 && (
                    <p style={{opacity:0.6}}>ไม่มีรายการอาหาร</p>
                )}

                {order.items.map((item: any) => (
                    <div key={item.item_id} className={styles.blogItem}>
                        
                        <div className={styles.quantity}>
                            <p>{item.quantity}</p>
                        </div>

                        <div className={styles.menu}>
                            <p>{item.menu_name}</p>

                            <ul className={styles.addOnList}>
                                {item.addons.map((add: any) => (
                                    <li key={add.group_id}>
                                        {add.group_name}: {add.option_name}
                                    </li>
                                ))}
                            </ul>
                        </div>

                        <div className={styles.price}>
                            <span>฿</span> <span>{item.price}</span>
                        </div>
                    </div>
                ))}

                <div className={styles.totalPrice}>
                    <p>รวมทั้งหมด</p>
                    <p>฿ {order.total_price}</p>
                </div>
            </div>


            <div className={styles.myBalance}>
                <div className={styles.content}>
                    <h2>My Balance</h2>
                    <div className={styles.blogBalance}>
                        <p>ยอดเงินคงเหลือ {order.balance} บาท</p>
                        <button className={styles.topUpBtn}>
                            <img src="/plus.svg" width={15} height={15} />
                            เติมเงิน
                        </button>
                    </div>

                    <div className={styles.blogActionBtn}>
                        <button className={styles.cancleBtn}>Cancel</button>
                        <button className={styles.acceptBtn}>ชำระเงิน</button>
                    </div>
                </div>
            </div>

        </div>
    )
}