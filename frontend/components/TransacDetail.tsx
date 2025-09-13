import { title } from "process"
import styles from "../styles/TransacDetail.module.css"

type TransacDetailProps = {
  head: string;
  detail?: string;
  date: string;
  price?: string;
  imgsrc: string;
  imgalt?: string;
};

export default function TransacDetail({ head, detail, date, price, imgsrc, imgalt }: TransacDetailProps) {
    return(
        <div className={styles.container}>
            <div className={styles.left_content}>
                <img src={imgsrc} alt={imgalt} />
            </div>
            <div className={styles.middle_content}>
                <p>{head}</p>
                <p>{detail}</p>
                <p>{date}</p>
            </div>
            <div className={styles.right_content}>
                <p>{price}</p>
            </div>
        </div>
    )
}