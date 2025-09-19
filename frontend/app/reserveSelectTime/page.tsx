"use client"
import styles from "./reserveSelectTime.module.css";
import { useRouter } from "next/navigation";

type TimeBtProps = { 
    time: string; 
    available: boolean;
};

export default function ReserveSelectTimePage() {
    const slots = ["10:00", "11:00", "12:00", "13:00"];
    return (
        <div className={styles.container}>
            <h1 className={styles.title}>เลือกช่วงเวลาที่ต้องการจองโต๊ะ</h1>
            <div>
                <div className={styles.infoContainer}>
                    <img src="/map_pin.svg" className={styles.infoIcon}/>
                    <p>โรงอาหารอาคารเรียนรวมสมเด็จพระเทพฯ ชั้น 2 ห้องแอร์</p>
                </div>
                <div className={styles.infoContainer}>
                    <img src="/info.svg" className={styles.infoIcon}/>
                    <p>มีเวลาในการใช้งานโต๊ะได้ slot ละ 15 นาที และต้องทำการสั่งอาหารภายใน 5 นาที</p>
                </div>
            </div>
            <TimeSlot slots={slots}/>
        </div>
    );
}

function TimeSlot({slots} : any) {
    return (
        <div>
            {slots.map((time : any, i : any) => (
            <div
                key={i}
                className={`${i === 0 ? styles.timeContainer1 : styles.timeContainer2} ${styles.timeContainerBase}`}
            >
                <p>{time} น.</p>
                <div className={styles.timeBtContainer}>
                    {[0, 15, 30, 45].map((m) => {
                    const hh = time.split(":")[0];
                    const mm = String(m).padStart(2, "0");
                    const t = `${hh}:${mm}`;
                    return <TimeBt key={t} time={t} available={true} />
                    })}
                </div>
            </div>
            ))}
        </div>
    );
}

function TimeBt({ time, available }: TimeBtProps) {
    const router = useRouter();

    return (
        <button
            className={`${available ? styles.timeBtAvl : styles.timeBtNotAvl} ${styles.timeBtBase}`}
            onClick={() => router.push(`/reserveSelectTable?time=${encodeURIComponent(time)}`)}
        >
            {time}          
        </button>
    );
}