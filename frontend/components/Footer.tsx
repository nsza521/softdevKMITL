"use client";
import styles from "../styles/Footer.module.css";
import Image from "next/image";

export default function Footer() {
  return (
    <footer className={styles.footer}>
        <div className={styles.footercontent}>
            <button className={styles.flogo}><Image src="/home.svg" alt="Home" width={80} height={56} /></button>
            <button className={styles.flogo}><Image src="/table.svg" alt="Home" width={80} height={56} /></button>
            <button className={styles.flogo}><Image src="/history.svg" alt="Home" width={80} height={56} /></button>
            <button className={styles.flogo}><Image src="/bell.svg" alt="Home" width={80} height={56} /></button>
            <button className={styles.flogo}><Image src="/user.svg" alt="Home" width={80} height={56} /></button>
        </div>
    </footer>
  );
}
