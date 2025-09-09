"use client";

import { useState } from "react";
import styles from "./home.module.css";
import { Noto_Sans_Thai } from "next/font/google";

const notoThai = Noto_Sans_Thai({
  subsets: ["thai"],
  weight: ["400", "700"],
  variable: "--font-noto-thai",
});

export default function LoginPage() {
  const [showPopup, setShowPopup] = useState(false);

  return (
    <div className={`${styles.container} ${notoThai.variable}`}>
      <div className={styles.headername}>
        <span className={styles.headernameedit}>
          สวัสดี ปกรณ์ บุญเกษม
          <button onClick={() => setShowPopup(true)}>
            <img src="/editpencil.svg" width={25} height={25} />
          </button>
        </span>
        <button className={styles.logoutbtn}>
          <img src="/logout.svg" width={25} height={25} />
          logout
        </button>
      </div>

      <div className={styles.boxs}>
        <span>ยอดเงินคงเหลือ xxxx.xx บาท</span>
        <button>
          <img src="/plus.svg" width={15} height={15} />
          เติมเงิน
        </button>
      </div>

      <div className={styles.boxs}>
        <span className={styles.boxspan}>
          <img src="/qr.svg" width={20} height={20} />
          คิวอาโค้ดของฉัน
        </span>
        <button>
          <img src="/show.svg" width={25} height={25} />
          ดู
        </button>
      </div>

      <div className={styles.headername}>โต๊ะตอนี้ </div>
      <div className={styles.table}></div>
      <button className={styles.tablebtn}>จองโต๊ะ</button>




        {/* Popup  */}
        {showPopup && (
        <div className={styles.popupbg} onClick={() => setShowPopup(false)}>
          <div className={styles.popup} onClick={(e) => e.stopPropagation()}>
            <h2>แก้ไขข้อมูลส่วนตัว</h2>
            <form className={styles.form} onSubmit={(e) => {
              e.preventDefault();
              // TODO: handle submit
              setShowPopup(false);
            }}>
              <div className={styles.formGroup}>
                <label>Name</label>
                <input type="text" name="name" placeholder="กรอกชื่อ" />
              </div>

              <div className={styles.formGroup}>
                <label>Surname</label>
                <input type="text" name="surname" placeholder="กรอกนามสกุล" />
              </div>

              <div className={styles.formGroup}>
                <label>Username</label>
                <input type="text" name="username" placeholder="กรอกชื่อผู้ใช้" />
              </div>

              <div className={styles.formGroup}>
                <label>Password</label>
                <input type="password" name="password" placeholder="กรอกรหัสผ่าน" />
              </div>

              <div className={styles.buttonGroup}>
                <button type="button" className={styles.cancelBtn} onClick={() => setShowPopup(false)}>ยกเลิก</button>
                <button type="submit" className={styles.submitBtn}>ยืนยัน</button>
              </div>
            </form>
          </div>
        </div>
      )}

    </div>
  );
}
