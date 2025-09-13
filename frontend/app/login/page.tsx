import styles from "./login.module.css";
import { Noto_Sans_Thai } from "next/font/google";

const notoThai = Noto_Sans_Thai({
  subsets: ["thai"],
  weight: ["400", "700"],
  variable: "--font-noto-thai",
});

export default function LoginPage() {
  return (
    <div className={`${styles.container} ${notoThai.variable}`}>

        <div className={styles.loginbox}>
            <form className={styles.logininbox}>
                <h2>Login</h2>
                <div className={styles.logininputbox}>
                    <p>Username</p>
                    <input type="text" name="" id=""  placeholder="Username / Email"/>
                </div>
                <div className={styles.logininputbox}>
                    <p>Password</p>
                    <input type="password" name="" id=""  placeholder="Password"/>
                </div>
                <button type="submit" className={styles.submitbtn}>
                     <span>Log In</span> <span className="material-symbols-outlined">arrow_forward</span>
                </button>
                <div className={styles.doyouhaveacc}>   
                    {/* <div> <span>Do you have any account yet ?  <a href="">&nbsp;Sign up here</a></span></div>
                    <div> <span>or</span></div> */}
                    <div> <span> <a href="">Forget Password ?</a></span></div>
                </div>
            </form>
        </div>

        <div className={styles.logo}>Logo Banner</div>
    </div>
  );
}
