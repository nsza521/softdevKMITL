'use client'

import { useState, useEffect } from "react"
import { useSearchParams } from "next/navigation";
import { useRouter } from "next/navigation";
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
    const router = useRouter();
    const order_id = searchParams.get("order_id") || ""
    const reservation_id = searchParams.get("reservationId") || ""

    const [order, setOrder] = useState<Order | null>(null)
    const [walletBalance, setWalletBalance] = useState<number>(0)
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
                        quantity: 1,
                        total_price: item.line_subtotal,
                        options: item.options?.map((add: Addon): Option => ({
                            group_name: add.group_name,
                            option_name: add.option_name
                        })) || []
                    }))
                };

                setOrder(formattedOrder);

                const walletRes = await fetch(`http://localhost:8080/customer/profile`, {
                    headers: { "Authorization": `Bearer ${token}` }
                });
                if (!walletRes.ok) throw new Error("โหลดข้อมูลกระเป๋าเงินไม่สำเร็จ")

                const walletData = await walletRes.json();
                setWalletBalance(walletData.wallet_balance);

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

    const handleCancel = () => {

        router.back();
    }

    const handlePurchase = async () => {
        const token = localStorage.getItem("token");
        if (!token) {
            alert("Token not found");
            return;
        }

        if (walletBalance < total_price) {
            alert("ยอดเงินคงเหลือไม่เพียงพอต่อการทำรายการ")
            return
        }

        try {
            const confirm = await fetch(`http://localhost:8080/table/reservation/${reservation_id}/confirm`, {
                method: "POST",
                headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json"
                },
            });

            if (!confirm.ok) {
                setError("cannot confirm reservation")
            }

            const confirm_res = await confirm.json();
            console.log(confirm_res)

            router.push(`/waitOthers?reservationId=${reservation_id}`)

        } catch (err) {
            console.log(err)
            alert("เกิดข้อผิดพลาด")
        }

    }

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

                    <div className={styles.itemDetail}>
                        {/* จำนวน */}
                        <div className={styles.quantity}>
                            <p>{item.quantity}</p>
                        </div>

                        {/* ชื่อเมนู + ตัวเลือก */}
                        <div className={styles.menu}>
                            <p>{item.menu_name}</p>

                            {/* ตัวเลือก addons */}
                            {item.options.length > 0 && (
                                <p className={styles.addOnList}>
                                    + {item.options.map(opt => opt.option_name).join(", ")}
                                </p>
                            )}

                        </div>
                    </div>

                    {/* ราคาของ item */}
                    <div className={styles.price}>
                        <p>฿ {item.total_price}</p>
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
                        <p>ยอดเงินคงเหลือ {walletBalance} บาท</p>
                        <button className={styles.topUpBtn}>
                            <img src="/plus.svg" width={15} height={15} />
                            เติมเงิน
                        </button>
                    </div>

                    <div className={styles.blogActionBtn}>
                        <button className={styles.cancleBtn} onClick={handleCancel}>Cancel</button>
                        <button className={styles.acceptBtn} onClick={handlePurchase}>ชำระเงิน</button>
                    </div>
                </div>
            </div>
        </div>
    );
}