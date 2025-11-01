'use client'

import { useState, useEffect } from "react"
import { useSearchParams } from "next/navigation";
import styles from "./orderMenuSummary.module.css"

type UUID = string;

interface Option {
    group_name: string;
    option_name: string;
}
interface OrderItem {
    item_id: UUID;
    menu_name: string;
    quantity: number;
    total_price: number;
    options: Option[];
}

interface Order {
    order_id: UUID;
    orders: OrderItem[];
}

interface Addon {
    group_name: string;
    option_name: string;
    qty?: number;
}

interface RawItem {
    order_item_id: UUID;
    menu_name: string;
    line_subtotal: number;
    options?: Addon[];
}

export default function OrderMenuSummaryPage() {
    const searchParams = useSearchParams();
    const order_id = searchParams.get("order_id") || ""
    const reservation_id = searchParams.get("reservationId") || ""

    const [order, setOrder] = useState<Order | null>(null)
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

                const formattedOrder: Order = {
                    order_id: data.order_id,
                    orders: data.items.map((item: RawItem): OrderItem => ({
                        item_id: item.order_item_id,
                        menu_name: item.menu_name,
                        quantity: 1, // ถ้า qty อยู่ที่ item ต้องเปลี่ยนเป็น item.qty
                        total_price: item.line_subtotal,
                        options: item.options?.map((add: Addon): Option => ({
                            group_name: add.group_name,
                            option_name: add.option_name
                        })) || []
                    }))
                };

                setOrder(formattedOrder);

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

    const total_price = order.orders.reduce((sum, item) => sum + item.total_price, 0);

    console.log("Rendered Order:", order);

    return (
        <div className={styles.container}>
            <div className={styles.myOrder}>
            <h2>My Order</h2>

            {/* ไม่มีรายการอาหาร */}
            {order.orders.length === 0 && (
                <p style={{ opacity: 0.6 }}>ไม่มีรายการอาหาร</p>
            )}

            {/* แสดงรายการอาหาร */}
            {order.orders.map((item: OrderItem) => (
                <div key={item.item_id} className={styles.blogItem}>

                {/* จำนวน */}
                <div className={styles.quantity}>
                    <p>{item.quantity}</p>
                </div>

                {/* ชื่อเมนู + ตัวเลือก */}
                <div className={styles.menu}>
                    <p>{item.menu_name}</p>

                    {/* ตัวเลือก addons */}
                    {item.options.length > 0 && (
                    <ul className={styles.addOnList}>
                        {item.options.map((add: Option, index: number) => (
                        <li key={index}>
                            {add.option_name}
                        </li>
                        ))}
                    </ul>
                    )}
                </div>

                {/* ราคาของ item */}
                <div className={styles.price}>
                    <span>฿</span> <span>{item.total_price}</span>
                </div>
                </div>
            ))}

            {/* ยอดรวมทั้งหมด */}
            <div className={styles.totalPrice}>
                <p>รวมทั้งหมด</p>
                <p>฿ {total_price}</p>
            </div>
            </div>

            {/* Balance ด้านล่าง */}
            <div className={styles.myBalance}>
            <div className={styles.content}>
                <h2>My Balance</h2>
                <div className={styles.blogBalance}>
                <p>ยอดเงินคงเหลือ  บาท</p>
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
    );
}