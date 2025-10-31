'use client'

import { useState } from "react"
import styles from "./orderMenuSummary.module.css"


const orders = [
    {
        id : 1,
        menu : "ข้าวไข่เจียว",
        price : 50,
        quantity : 1,
        addOn : ['พิเศษ'],    
    },
    {
        id : 2,
        menu : "ข้าวหน้าเนื้อ",
        price : 70,
        quantity : 2,
        addOn : ['พิเศษ','ไข่ดาว'],     
    }
]




export default function orderMenuSummaryPage() {
    return (
        <div className={styles.container}>
            <div className={styles.myOrder}>
                <h2>My Order</h2>
            {orders.map ((ord) => (
                <div className={styles.blogItem}>
                    <div className={styles.quantity}>
                        <p>{ord.quantity}</p>
                    </div>
                    <div className={styles.menu}>
                        <p>{ord.menu}</p>
                        <ul className={styles.addOnList}>
                            <button><img src="Add_Plus_Circle_Vector.svg" alt="addBtn"/></button>
                            {ord.addOn.map((addOns) =>(
                                <li>{addOns}</li>
                            ))}
                        </ul>
                    </div>
                    <div className={styles.price}>
                        <span>฿</span> <span>{ord.price}</span>
                    </div>
                </div>
            ))
            }
                <div className={styles.totalPrice}>
                    <p>รวมทั้งหมด</p>
                    <p>฿ {120}</p>
                </div>
            </div>
            <div className={styles.myBalance}>
                <div className={styles.content}>
                    <h2>My Balance</h2>
                    <div className={styles.blogBalance}>
                        <p>ยอดเงินคงเหลือ {222} บาท</p>
                        <button className={styles.topUpBtn}>
                             <img src="/plus.svg" width={15} height={15} />
                             เติมเงิน
                        </button>
                    </div>
                    <div className={styles.blogActionBtn}>
                        <button className={styles.cancleBtn}>Cancle</button>
                        <button className={styles.acceptBtn}>ชำระเงิน</button>
                    </div>
                </div>
            </div>
        </div>
    )
}