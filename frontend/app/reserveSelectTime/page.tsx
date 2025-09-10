"use client"
import { useState } from "react";
import styles from "./reserveSelectTime.module.css";

export default function ReserveSelectTimePage() {
    const slots = ["10:00", "11:00", "12:00", "13:00"];
    return (
        <div className={styles.container}>
            <h1 className={styles.title}>เลือกช่วงเวลาที่ต้องการจองโต๊ะ</h1>
            <div>
                <div className={styles.infoContainer}>
                    <img src="/map_pin.svg" className={styles.infoIcon}></img>
                    <p>โรงอาหารอาคารเรียนรวมสมเด็จพระเทพฯ ชั้น 2 ห้องแอร์</p>
                </div>
                <div className={styles.infoContainer}>
                    <img src="/info.svg" className={styles.infoIcon}></img>
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
            {[...Array(4)].map((_, j) => (
              <TimeBt key={j} time={time} />
            ))}
          </div>
        </div>
      ))}
        </div>
    );
}

function TimeBt({time} : any) {
    const [isAvailable, setIsAvailable] = useState(1);

    return (
        <button 
            className={`${isAvailable ? styles.timeBtAvl : styles.timeBtNotAvl} ${styles.timeBtBase}`}
            // onClick={() => setIsAvailable(isAvailable ? 0 : 1)}
            >
                {time}          
        </button>
    );
}