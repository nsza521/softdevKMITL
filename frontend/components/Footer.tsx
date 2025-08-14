"use client";
import styles from "../styles/Footer.module.css";

export default function Footer() {
  return (
    <footer className={styles.footer}>
        <div className={styles.footercontent}>
            <button className={styles.flogo}><span className="material-symbols-outlined">home</span><p>หน้าหลัก</p></button>
            <button className={styles.flogo}><span className="material-symbols-outlined">apps</span><p>บริการ</p></button>
            <button className={styles.flogo}><span className="material-symbols-outlined">chat</span><p>กล่องข้อความ</p></button>
            <button className={styles.flogo}><span className="material-symbols-outlined">settings</span><p>การตั้งค่า</p></button>
        </div>
        <div className={styles.footerdec}></div>
    </footer>
  );
}
