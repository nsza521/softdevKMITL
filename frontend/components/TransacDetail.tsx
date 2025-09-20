import { title } from "process"
import styles from "../styles/TransacDetail.module.css"

type TransacDetailProps = {
  head: string;
  detail?: string;
  date: string;
  viewdetail?: string;
  price?: string;
  imgsrc: string;
  imgalt?: string;
};

export default function TransacDetail({ head, detail, date, viewdetail, price, imgsrc, imgalt }: TransacDetailProps) {
    return(
        <div className={styles.container}>
            <div className={styles.left_content}>
                <img src={imgsrc} alt={imgalt} />
            </div>
            <div className={styles.middle_content}>
                <p>{head}</p>
                <p>{detail}</p>
                <p>{date}</p>
                <div className={styles.viewdetail}>
                    <p>{viewdetail}</p>
                    <img src="Arrow_Right_LG.svg" alt="Arrow_Right_LG" />
                </div>
            </div>
            <div className={styles.right_content}>
                <p>{price}</p>
            </div>
        </div>
    )
}