"use client";

import { useState } from "react";
import styles from "./topup.module.css"
import Navbar from "@/components/Navbar";
import TransacDetail from "@/components/TransacDetail";
import { title } from "process";
import { Noto_Sans_Thai } from "next/font/google";



const notoThai = Noto_Sans_Thai({
  subsets: ["thai"],
  weight: ["400", "700"],
  variable: "--font-noto-thai",
});

const iconPayment = {
    qr : "qr.png",
    kbank : "",
    scb : ""
} ;

const balance = 4000;
const amounts = [100, 200, 300, 500];
const methods = [
  { name: "PromptPay", icon: "/promtpay.png" },
  { name: "KBANK", icon: "/kbank.png" },
  { name: "SCB", icon: "/scb.png" },
];

export default function TopUpPage(){
    const [selected, setSelected] = useState<number | "custom" | null>(null);
    const [customAmount, setCustomAmount] = useState<number | "">("");
    const [active, setActive] = useState<string | null>(null);
    const [showPopup, setShowPopup] = useState(false);

    return(
        <div className={`${styles.container} ${notoThai.variable}`}>
            <Navbar title="เติมเงิน"/>
            <div className={styles.content}>
                <div className={styles.balance}>
                    <p>จำนวนเงินของคุณ</p>
                    <p>{balance} บาท</p>
                </div>
                <div className={styles.amount}>
                    <div className={styles.header}>จำนวนที่ต้องการเติม</div>
                    <div className={styles.options}>
                        {amounts.map((amt) => (
                        <button
                            key={amt}
                            className={`${styles.optionAmount} ${
                            selected === amt ? styles.active : ""
                            }`}
                            onClick={() => setSelected(amt)}
                        >
                            {amt} บาท
                        </button>  ))}
                    </div>

                    <div className={styles.defined}>
                        <button
                        className={`${styles.optionBtn} ${
                            selected === "custom" ? styles.active : ""
                        }`}
                        onClick={() => setSelected("custom")}
                        >
                        กำหนดเอง
                        </button>
                        
                        {selected === "custom" && (
                            <div className={styles.customInput}>
                            <input
                                type="number"
                                placeholder="กรอกจำนวนเงิน"
                                value={customAmount}
                                onChange={(e) => setCustomAmount(Number(e.target.value))}
                            />
                            <span>บาท</span>
                            </div>
                        )}
                    </div>
                </div>

                <div className={styles.paymentmethod}>
                    <div className={styles.header}>เติมเงินผ่าน</div>
                    <div className={styles.optionMethod}>
                     {methods.map((mth) => (
                         <button
                         key={mth.name}
                         className={`${styles.payment} ${
                             active === mth.name ? styles.activepayment : ""
                            }`}
                            onClick={() => setActive(mth.name)}
                            >
                        <img src={mth.icon} alt={mth.name} className={styles.icon} />
                        {mth.name}
                    </button>
                    ))}
                    </div>
                </div>

                <div className={styles.btnTopup}>
                    <button onClick={() => setShowPopup(true)}>เติมเงิน</button>
                </div>
                {showPopup && (
                    <div className={styles.overlayTopup}>
                        <p>ยืนยันเติมเงินเข้าสู่ระบบ</p>
                        <div className={styles.actionBtn}>
                            <button onClick={() => setShowPopup(false)}>ยกเลิก</button>
                            <button>ยืนยัน</button>
                        </div>
                    </div>
                )}
            </div>
        </div>
    )
}
