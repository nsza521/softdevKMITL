// /components/Navbar.tsx
"use client";
import styles from "../styles/Nav.module.css";
import { useRouter } from "next/navigation";

export default function Navbar({ title }: { title: string }) {
  const router = useRouter();
  const cursor = { cursor: "pointer" };

  return (
    <div className={styles.nav}>
      <span
      className="material-symbols-outlined"
      onClick={() => router.back()}
      style={cursor}
      >
        arrow_back_ios
      </span>
      <span>{title}</span>
    </div>
  );
}
