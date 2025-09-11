"use client";
import { usePathname, useRouter } from "next/navigation";
import styles from "../styles/Footer.module.css";
import Image from "next/image";

export default function Footer() {
  const pathname = usePathname();
  const router = useRouter();

  return (
    <footer className={styles.footer}>
      <div className={styles.footercontent}>
        <button
          className={`${styles.flogo} ${pathname === "/home" ? styles.active : ""}`}
          onClick={() => router.push("/home")}
        >
          <Image
            src={pathname === "/home" ? "/home_2.svg" : "/home_1.svg"}
            alt="Home"
            width={30}
            height={30}
          />
          หน้าหลัก
        </button>

        <button
          className={`${styles.flogo} ${pathname === "/reserveSelectTime" ? styles.active : ""}`}
          onClick={() => router.push("/reserveSelectTime")}
        >
          <Image
            src={pathname === "/reserveSelectTime" ? "/table_2.svg" : "/table_1.svg"}
            alt="Table"
            width={30}
            height={30}
          />
          จองโต๊ะ
        </button>

        <button
          className={`${styles.flogo} ${pathname === "/history" ? styles.active : ""}`}
          onClick={() => router.push("/history")}
        >
          <Image
            src={pathname === "/history" ? "/history_2.svg" : "/history_1.svg"}
            alt="History"
            width={30}
            height={30}
          />
          ประวัติ
        </button>

        <button
          className={`${styles.flogo} ${pathname === "/noti" ? styles.active : ""}`}
          onClick={() => router.push("/noti")}
        >
          <Image
            src={pathname === "/noti" ? "/noti_2.svg" : "/noti_1.svg"}
            alt="Notifications"
            width={30}
            height={30}
          />
          แจ้งเตือน
        </button>
      </div>
    </footer>
  );
}
