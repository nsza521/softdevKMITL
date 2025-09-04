// /components/Navbar.tsx
import styles from "../styles/Nav.module.css";

export default function Navbar({ title }: { title: string }) {
  return (
    <div className={styles.nav}>
      <span className="material-symbols-outlined">arrow_back_ios</span>
      <span>{title}</span>
    </div>
  );
}
