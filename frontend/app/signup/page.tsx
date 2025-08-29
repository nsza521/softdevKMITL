import styles from "./signup.module.css";
import { Noto_Sans_Thai } from "next/font/google";

const notoThai = Noto_Sans_Thai({
  subsets: ["thai"],
  weight: ["400", "700"],
  variable: "--font-noto-thai",
});

export default function LoginPage() {
  return (
    <div className={`${styles.container} ${notoThai.variable}`}>
        <div className={styles.logo}>
            <div>
                Logo Banner
            </div>
        </div>
        <div className={styles.registerbox}>
            <form>
                <h2>Sign UP</h2>
                <input type="text" placeholder="Name"/>
                <input type="text" placeholder="Surname"/>
                <input type="text" placeholder="Email"/>
                <input type="text" placeholder="Username"/>
                <input type="text" placeholder="Password"/>
                <button type="submit" className={styles.submitbtn}>
                     <span>Sign Up</span> <span className="material-symbols-outlined">arrow_forward</span>
                </button>
            </form>
        </div>  

    </div>
  );
}
