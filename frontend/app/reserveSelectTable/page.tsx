"use client";
import styles from "./reserveSelectTable.module.css";
import { useSearchParams } from "next/navigation";

type Table = {
    row: string;
    col: number;
    occupied: number;
    capacity: number;
    available: boolean;
}

const tables: Table[] = [
  { row: "A", col: 1, occupied: 0, capacity: 6, available: true },
  { row: "A", col: 2, occupied: 0, capacity: 6, available: true },
  { row: "A", col: 3, occupied: 0, capacity: 6, available: true },
  { row: "B", col: 1, occupied: 0, capacity: 6, available: true },
  { row: "B", col: 2, occupied: 4, capacity: 6, available: false }, 
  { row: "B", col: 3, occupied: 0, capacity: 6, available: true },
];

export default function ReserveSelectTablePage() {
    const searchParam = useSearchParams();
    const time = searchParam.get("time");

    return (
        <div className={styles.container}>
            <h1 className={styles.title}>เลือกโต๊ะที่ต้องการจอง</h1>
            <TableLayout />
            <p className={styles.txt}>หรือ</p>
            <button className={styles.soloBt}>Join with Anyone</button>
        </div>
    );
}

function TableLayout() {
    const rows = Array.from(new Set(tables.map((t) => t.row)));
    const cols = Array.from(new Set(tables.map((t) => t.col)));

    return (
        <div className={styles.layoutContainer}>
            <div className={styles.tableContainer}>
                {rows.map((c) => (
                <span key={c} className={styles.rowLabel}>{c}</span>
                ))}
            </div>

            {rows.map((r) => (
                <div key={r} className={`${styles.tableContainer} ${styles.tableRowContainer}`}>
                    <span className={styles.rowLabel}>{r}</span>
                    {cols.map((c) => {
                        const table = tables.find((t) => t.row === r && t.col === c);
                        return table ? (
                        <TableIcon key={`${r}${c}`} table={table} />
                        ) : (
                        <span key={`${r}${c}`}></span>
                        );
                    })}
                </div>
            ))}
        </div>
    );
}

function TableIcon({ table }: { table: Table }) {
    const { occupied, capacity, available } = table;
    return (
        <div
            className={`${styles.tableIcon} ${available ? styles.tableAvailable : styles.tableUnavailable}`}>
            <img
                src={
                available
                    ? "./table_layout_aval.svg"
                    : "./table_layout_notaval.svg" }/>
            <p>{occupied}/{capacity}</p>
        </div>
  );
}