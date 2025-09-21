"use client";

import { useState } from "react";
import styles from "./topup.module.css"
import { Noto_Sans_Thai } from "next/font/google";



const notoThai = Noto_Sans_Thai({
  subsets: ["thai"],
  weight: ["400", "700"],
  variable: "--font-noto-thai",
});

// const iconPayment = {
//     qr : "qr.png",
//     kbank : "",
//     scb : ""
// } ;

const balance = 4000;
const amounts = [100, 200, 300, 500];
const methods = [
  { name: "PromptPay", icon: "/promtpay.png" },
  { name: "KBANK", icon: "/kbank.png" },
  { name: "SCB", icon: "/scb.png" },
];

export default function TopUpPage(){
    const [amount, setAmount] = useState<number | null>(null);
    const [bank, setBank] = useState<string | null>(null);
    const [error, setError] = useState("");
    const [success, setSuccess] = useState(false);
    const [custom, setCustom] = useState(false);
    

    const handleTopup = () => {
        if (!amount && !bank){
            setError("กรุณาเลือกจำนวนเงินและธนาคาร");
        }else if (!amount) {
            setError("กรุณาเลือกจำนวนเงิน");
        }else if (!bank) {
            setError("กรุณาเลือกธนาคาร");
        }else {
            setError("");
            //api
            setSuccess(true);
        }
    }


    return(
        <div className={`${styles.container} ${notoThai.variable}`}>
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
                            amount === amt && !custom ? styles.active : ""
                            }`}
                            onClick={() => {
                            setAmount(amt);
                            setCustom(false); // ถ้าเลือกปุ่มตัวเลข ปิดโหมดกำหนดเอง
                            }}
                        >
                            {amt} บาท
                        </button>
                        ))}
                    </div>
                <div className={styles.customAmount}>
                    <button
                    className={`${styles.optionAmount} ${custom ? styles.active : ""}`}
                    onClick={() => {
                        setCustom(true);
                        setAmount(null); // reset amount ปกติ
                    }}
                    >
                    กำหนดเอง
                    </button>

                    {custom && (
                        <div>
                            <input
                                type="number"
                                placeholder="กรอกจำนวนเงิน"
                                value={amount ?? ""} // ให้ input แสดงค่าที่พิมพ์
                                onChange={(e) => {
                                    const val = e.target.value;
                                    setAmount(val === "" ? null : Number(val));
                                }}
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
                             bank === mth.name ? styles.activepayment : ""
                            }`}
                            onClick={() => setBank(mth.name)}
                            >
                        <img src={mth.icon} alt={mth.name} className={styles.icon} />
                        {mth.name}
                    </button>
                    ))}
                    </div>
                </div>

                {error && (
                    <div className={styles.overlayTopup} onClick={() => setError("")}>
                        <p>{error}</p>
                        <div className={styles.actionBtn}>
                        <button onClick={() => setError("")}>ปิด</button>
                        </div>
                    </div>
                )}

                {success && (
                    <div className={styles.overlayTopup} onClick={() => setSuccess(false)}>
                            <p>ต้องการเติมเงิน {amount} บาท ผ่าน {bank} ใช่หรือไม่</p>
                            <div className={styles.actionBtn}>
                                <button onClick={() => setSuccess(false)}>ยกเลิก</button>
                                <button onClick={() => setSuccess(false)}>ตกลง</button>
                            </div>
                    </div>
                )}

                <div className={styles.btnTopup}>
                    <button onClick={handleTopup}>เติมเงิน</button>
                </div>
            </div>
        </div>
    )
}
