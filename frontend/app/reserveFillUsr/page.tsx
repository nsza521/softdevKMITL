"use client";
import { table } from "console";
import styles from "./reserveFillUsr.module.css";
import { useSearchParams } from "next/navigation";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

type UUID = string;

type Member = {
    username: string;
};

type Reservation = {
    table_timeslot_id: UUID;
    reserve_people: number;
    random: boolean;
    members: Member[];
};

type Table = {
    id: UUID;
    row: string;
    col: string;
    max_seats: number;
    status: string;
    reserved_seats: number;
};

export default function ReserveFillUsrPage() {
    const searchParam = useSearchParams();
    const random = searchParam.get("random");
    const router = useRouter();
    const table: Table = {
        id: "ae7bbca8-4a80-4195-9b4b-ab881426f6f1",
        row: "A",
        col: "1",
        max_seats: 6,
        status: "available",
        reserved_seats: 5,
    };

    const token = localStorage.getItem("token");
    
    return (
        <div className={styles.container}>
            <TableInfo table={table} occupied={0}/>
            <Members token="token"/>
            <div>
                <img src="/info.svg"/>
                <p>สมาชิกทุกท่านจะมีเวลาในการสั่งอาหาร 5 นาที หากทุกท่านไม่ทำการสั่งอาหารภายใน 5 นาที จะถือว่าสละสิทธิ์</p>
            </div>
            <button>
                เชิญเพื่อนและเริ่มสั่งอาหาร
                <img src="/Arrow_Right_MD.svg"/>
            </button>
        </div>
    );
}

function TableInfo({ table, occupied }: { table: Table; occupied: number }) {
    const [allowOthers, setAllowOthers] = useState(false);
    const all_occupied = table.reserved_seats + occupied;
    const min_allow = table.max_seats * 0.8;

    const canClick = all_occupied >= min_allow;
    const isChecked = !canClick ? true : allowOthers;

    return (
        <div>
            <div className={styles.tableInfoCon}>
                <p>โต๊ะของคุณ : {table.row}{table.col}</p>
                <TableIcon table={table} occupied={occupied} />
            </div>
            <label htmlFor="allow_others_join" className={styles.checkbox}>
                <input
                    type="checkbox"
                    id="allow_others_join"
                    name="allow_others_join"
                    checked={isChecked}
                    disabled={!canClick}
                    onChange={(e) => setAllowOthers(e.target.checked)}
                />
                <span className={styles.checkmark}></span>
                อนุญาตให้ผู้อื่นเข้าร่วม
            </label>
        </div>
    );
}

function TableIcon({ table, occupied }: { table: Table, occupied: number }) {
    const router = useRouter();
    const available = table.status === "available";

    return (
        <div
        className={`${styles.tableIcon} ${
            available ? styles.tableAvailable : styles.tableUnavailable
        }`}
        >
            <img src={available ? "./table_layout_aval.svg" : "./table_layout_notaval.svg"}/>
            <p>{table.reserved_seats + occupied}/{table.max_seats}</p>
        </div>
    );
}

type MyProfile = {
    id: UUID;
    username: string;
    email: string;
    first_name: string;
    last_name: string;
    wallet_balance: number;
};

type MemberInfo = {
    username: string;
    first_name: string;
};

function Members({ token }: { token: string }) {
    const [myProfile, setMyProfile] = useState<MyProfile | null>(null);
    const [members, setMembers] = useState<MemberInfo[]>([{ username: "", first_name: "" }]);
    const [error, setError] = useState<string | null>(null);

    // useEffect(() => {
    //     const fetchMyProfile = async () => {
    //     try {
    //         const res = await fetch("http://localhost:8080/profile/me", {
    //         headers: { Authorization: `Bearer ${token}` },
    //         });
    //         if (!res.ok) throw new Error("ไม่สามารถดึงข้อมูลโปรไฟล์ได้");
    //         const data = await res.json();
    //         setMyProfile(data);
    //     } catch (err: any) {
    //         setError(err.message);
    //     }
    //     };
    //     fetchMyProfile();
    // }, [token]);

    // const handleChangeMember = (index: number, field: keyof MemberInfo, value: string) => {
    //     const updated = [...members];
    //     updated[index][field] = value;
    //     setMembers(updated);
    // };

    // const handleUsernameBlur = async (index: number) => {
    //     const username = members[index].username.trim();
    //     if (!username) return;

    //     try {
    //         const res = await fetch("http://localhost:8080/customer/firstname", {
    //             method: "POST",
    //             headers: {
    //                 "Content-Type": "application/json", 
    //                 Authorization: token,   
    //             },
    //             body: JSON.stringify({ username }),  
    //         });

    //         if (!res.ok) throw new Error("ไม่พบผู้ใช้");
    //         const data = await res.json();

    //         handleChangeMember(index, "first_name", data.first_name || "");
    //         } catch (err) {
    //             console.error("fetch username error:", err);
    //         }
    // };

    // if (error) return <p>เกิดข้อผิดพลาด: {error}</p>;
    // if (!myProfile) return <p>กำลังโหลดข้อมูล...</p>;

    return (
        <div className={styles.membersCon}>
            <div className={styles.formTitleWrapper}>
                <div className={styles.titleLine}></div>
                <h1>สมาชิกที่เข้าร่วม</h1>
                <div className={styles.titleLine}></div>
            </div>
            <label className={styles.sectionLabel}>คุณ :</label>
            <div className={styles.inputCon}>
                <input
                type="text"
                // value={myProfile.username}
                disabled
                />
                <input
                type="text"
                // value={`${myProfile.first_name} ${myProfile.last_name}`}
                disabled
                />
            </div>

            {members.map((m, i) => (
                <div key={i} className={styles.formSection}>
                    <label className={styles.sectionLabel}>สมาชิกคนที่ {i + 2}</label>
                    <div className={styles.inputCon}>
                        <input
                        type="text"
                        placeholder="@username"
                        value={m.username}
                        // onChange={(e) => handleChangeMember(i, "username", e.target.value)}
                        // onBlur={() => handleUsernameBlur(i)}
                        />
                        <input
                        type="text"
                        placeholder="ชื่อจะถูกกรอกอัตโนมัติ"
                        value={m.first_name}
                        // onChange={(e) => handleChangeMember(i, "first_name", e.target.value)}
                        />
                    </div>
                </div>
            ))}

            <button
            className={styles.addUserBt}
            onClick={() => setMembers([...members, { username: "", first_name: "" }])}
            >
                <img src="/add_user.svg"/>
                เพิ่มสมาชิก
            </button>
        </div>
    );
}