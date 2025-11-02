"use client";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import styles from "./reserveSelectTime.module.css";

type UUID = string;

type Timeslot = { 
  time_slot_id: UUID;
  start_time: string; 
  end_time: string;
};

export default function ReserveSelectTimePage() {
    const [slots, setSlots] = useState<Timeslot[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchSlots = async () => {
            try {
                const res = await fetch("http://localhost:8080/table/timeslot/all");
                if (!res.ok) throw new Error("ไม่สามารถดึงข้อมูลได้");

                const json = await res.json();
                const data: Timeslot[] = Array.isArray(json.timeslots)
                ? json.timeslots.map((slot: any) => ({
                    time_slot_id: slot.timeslot_id,
                    start_time: slot.start_time,
                    end_time: slot.end_time,
                    }))
                : [];

                setSlots(data);
            } catch (err: any) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };
        fetchSlots();
    }, []);

    if (loading) return <div className={styles.container}><p>กำลังโหลดข้อมูล...</p></div>;
    if (error) return <div className={styles.container}><p>เกิดข้อผิดพลาด: {error}</p></div>;

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

function TimeSlot({ slots }: { slots: Timeslot[] }) {
  if (slots.length === 0) {
    return <p>ไม่มีช่วงเวลาให้เลือก</p>;
  }

  const grouped: Record<string, Timeslot[]> = {};

  slots.forEach((slot) => {
    const [hour] = slot.start_time.split(":");
    const hourLabel = `${hour.padStart(2, "0")}:00`;
    if (!grouped[hourLabel]) grouped[hourLabel] = [];
    grouped[hourLabel].push(slot);
  });

  const sortedHours = Object.keys(grouped).sort(
    (a, b) => Number(a.split(":")[0]) - Number(b.split(":")[0])
  );

  return (
    <div>
      {sortedHours.map((hour, i) => (
        <div
          key={hour}
          className={`${i == 0 ? styles.timeContainer1 : styles.timeContainer2} ${styles.timeContainerBase}`}
        >
          <p>{hour} น.</p>
          <div className={styles.timeBtContainer}>
            {grouped[hour].map((slot) => (
              <TimeBt
                key={slot.time_slot_id}
                timeSlotId={slot.time_slot_id}
                time={slot.start_time}
              />
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}

interface TimeBtProps {
  timeSlotId: UUID;
  time: string;
}

function TimeBt({ timeSlotId, time }: TimeBtProps) {
  const router = useRouter();

  return (
    <button
      className={`${styles.timeBtAvl} ${styles.timeBtBase}`}
      onClick={() => router.push(`/reserveSelectTable?timeSlotId=${encodeURIComponent(timeSlotId)}`)}
    >
      {time}
    </button>
  );
}