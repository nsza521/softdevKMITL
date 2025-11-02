"use client";

import { useEffect, useState } from "react";
import styles from "./topup.module.css"
import { Noto_Sans_Thai } from "next/font/google";

const notoThai = Noto_Sans_Thai({
  subsets: ["thai"],
  weight: ["400", "700"],
  variable: "--font-noto-thai",
});


const amounts = [100, 200, 300, 500];
const bankImages: Record<string, string> = {
  KBANK: "/kbank.png",
  SCB: "/scb.png",
  Promtpay: "/promtpay.png",
};


interface Balance {
     wallet_balance : number;
}

interface PaymentMethod {
    payment_method_id : string;
    name : string;
}
const initialBalance: Balance = { wallet_balance: 0 };

export default function TopUpPage(){
    const [amount, setAmount] = useState<number | null>(null);
    const [bank, setBank] = useState<string | null>(null);
    const [nameBank, setNameBank] = useState<string | null>(null);
    const [error, setError] = useState("");
    const [success, setSuccess] = useState(false);
    const [custom, setCustom] = useState(false);

   const [balance,setBalance] = useState<Balance>(initialBalance);
   const [methods, setMethods] = useState<PaymentMethod[]>([]);
   const [isLoading, setIsLoading] = useState(false);
    

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
    const confirmTopup = async () => {
    try {
      const token = localStorage.getItem("token");
      if (!token) {
        alert("กรุณาเข้าสู่ระบบก่อนทำรายการ");
        return;
      }

      setIsLoading(true);

      const res = await fetch("http://localhost:8080/payment/topup/wallet", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          payment_method_id: bank,
          amount: amount,
        }),
      });

      if (!res.ok) {
        throw new Error("Top-up failed");
      }

      const data = await res.json();
      console.log("Top-up success:", data);
      alert("เติมเงินสำเร็จ!");

      // ✅ อัปเดตยอดเงินใหม่ทันที
      setBalance((prev) => ({
        wallet_balance: prev.wallet_balance + (amount ?? 0),
      }));
    } catch (error) {
      console.error(error);
      alert("เกิดข้อผิดพลาดในการเติมเงิน");
    } finally {
      setIsLoading(false);
      setSuccess(false);
    }
  };

    useEffect(() => {
        const fetchTopUp = async() => {
            try{
                const token = localStorage.getItem("token");
                const resProfile = await fetch("http://localhost:8080/customer/profile",{
                    headers: { Authorization: `Bearer ${token}` },
                });
                const dataProfile = await resProfile.json();
                setBalance({ wallet_balance: dataProfile.wallet_balance });

                const resPayment = await fetch("http://localhost:8080/payment/topup/method/all",{
                    headers: { Authorization: `Bearer ${token}` },
                });
                const dataPayment = await resPayment.json();
                setMethods(dataPayment.payment_methods);
            }catch(err){
                console.error(err)
            }
        }
        fetchTopUp(); 
    },[])

    return(
        <div className={`${styles.container} ${notoThai.variable}`}>
            <div className={styles.content}>
                <div className={styles.balance}>
                    <p>จำนวนเงินของคุณ</p>
                    <p>{balance.wallet_balance} บาท</p>
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
                             bank === mth.payment_method_id ? styles.activepayment : ""
                            }`}
                            onClick={() => [setBank(mth.payment_method_id),setNameBank(mth.name)]}
                            >
                           <img
                            src={bankImages[mth.name] || "/default.png"}
                            alt={mth.name}
                            className={styles.bankLogo}
                            />   
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
                        <p>ต้องการเติมเงิน {amount} บาท ผ่าน {nameBank} ใช่หรือไม่</p>
                        <div className={styles.actionBtn}>
                        <button onClick={() => setSuccess(false)}>ยกเลิก</button>
                        <button
                            onClick={confirmTopup}
                        >
                            ตกลง
                        </button>
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


