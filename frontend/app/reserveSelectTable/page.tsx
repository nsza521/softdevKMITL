"use client";
import styles from "./reserveSelectTable.module.css";
import { useSearchParams } from "next/navigation";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

type UUID = string;

type TableTimeSlot = {
  id: UUID;
  row: string;
  col: string;
  max_seats: number;
  status: string;
  reserved_seats: number;
};

export default function ReserveSelectTablePage() {
    const searchParam = useSearchParams();
    const time_slot_id = searchParam.get("timeSlotId");
    const router = useRouter();

    const [tableTimeSlots, setTableTimeSlots] = useState<TableTimeSlot[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchSlots = async () => {
            try {
                const token = localStorage.getItem("token");
                const res = await fetch(`http://localhost:8080/table/table_timeslot/all/${time_slot_id}`, {
                    headers: {
                        "Content-Type": "application/json",
                        ...(token ? { Authorization: `Bearer ${token}` } : {}),
                    },
                });

                if (!res.ok) throw new Error("ไม่สามารถดึงข้อมูลได้");

                const json = await res.json();
                const data: TableTimeSlot[] = Array.isArray(json.table_timeslots)
                ? json.table_timeslots.map((t: any) => ({
                    id: t.id,
                    row: t.table_row,
                    col: t.table_col,
                    max_seats: t.max_seats,
                    status: t.status,
                    reserved_seats: t.reserved_seats,
                    }))
                : [];

                setTableTimeSlots(data);
            } catch (err: any) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchSlots();
    }, [time_slot_id]);

    if (loading) return <div className={styles.container}><p>กำลังโหลดข้อมูล...</p></div>;
    if (error) return <div className={styles.container}><p>เกิดข้อผิดพลาด: {error}</p></div>;

    return (
        <div className={styles.container}>
            <h1 className={styles.title}>เลือกโต๊ะที่ต้องการจอง</h1>
            <TableLayout tables={tableTimeSlots} />
            <p className={styles.txt}>หรือ</p>
            <button
                className={styles.soloBt}
                onClick={() => router.push(`/reserveFillUsr?random=${encodeURIComponent(true)}`)}
            >
                Join with Anyone
            </button>
        </div>
    );
}

function TableLayout({ tables }: { tables: TableTimeSlot[] }) {
    const rows = Array.from(new Set(tables.map((t) => t.row))).sort();
    const cols = Array.from(new Set(tables.map((t) => t.col))).sort(
        (a, b) => Number(a) - Number(b)
    );

    return (
        <div className={styles.layoutContainer}>
            <div className={styles.tableContainer}>
                <span></span>
                {cols.map((c) => (
                <span key={c}>{c}</span>
                ))}
            </div>

            {rows.map((r) => (
                <div key={r} className={`${styles.tableContainer} ${styles.tableRowContainer}`}>
                    <span className={styles.rowLabel}>{r}</span>
                    {cols.map((c) => {
                        const table = tables.find((t) => t.row === r && t.col === c);
                        return table ? (
                        <TableIcon key={table.id} table={table} />
                        ) : (
                        <span key={`${r}${c}`}></span>
                        );
                    })}
                </div>
            ))}

            <div className={styles.compassCon}>
                <img src="./compass.svg" />
                <p>W</p>
            </div>
        </div>
    );
}


function TableIcon({ table }: { table: TableTimeSlot }) {
    const router = useRouter();
    const available = table.status === "available";

    return (
        <button
        className={`${styles.tableIcon} ${
            available ? styles.tableAvailable : styles.tableUnavailable
        }`}
        onClick={() => router.push(`/reserveFillUsr?random=${encodeURIComponent(false)}&table_timeslot_id=${encodeURIComponent(table.id)}`)}
        >
            <img src={available ? "./table_layout_aval.svg" : "./table_layout_notaval.svg"}/>
            <p>{table.reserved_seats}/{table.max_seats}</p>
        </button>
    );
}