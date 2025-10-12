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

export default function ReserveFillUsrPage() {
    const searchParam = useSearchParams();
    const random = searchParam.get("random") === "true" || false;
    const table_id = searchParam.get("table_timeslot_id");

    const [table, setTable] = useState<Table | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [selectedMemberIndex, setSelectedMemberIndex] = useState<number>(3);
    
    const reservation: Reservation = {
        table_timeslot_id: table_id || "",
        reserve_people: selectedMemberIndex,
        random: random,
        members: [],
    };

    useEffect(() => {
        const fetchSlots = async () => {
        try {
            const token = localStorage.getItem("token");
            const res = await fetch(`http://localhost:8080/table/tabletimeslot/${table_id}`, {
            headers: {
                "Content-Type": "application/json",
                ...(token ? { Authorization: `Bearer ${token}` } : {}),
            },
            });

            if (!res.ok) throw new Error("ไม่สามารถดึงข้อมูลได้");

            const json = await res.json();
            const t = json.table_timeslot;

            if (!t) throw new Error("ไม่พบข้อมูลโต๊ะ");

            const data: Table = {
                id: t.id,
                row: t.table_row,
                col: t.table_col,
                max_seats: t.max_seats,
                status: t.status,
                reserved_seats: t.reserved_seats,
            };

            setTable(data);
        } catch (err: any) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
        };

        if (table_id) fetchSlots();
    }, [table_id]);

    if (loading) return <p>กำลังโหลดข้อมูล...</p>;
    if (error) return <p style={{ color: "red" }}>เกิดข้อผิดพลาด: {error}</p>;
    if (!table) return <p>ไม่พบข้อมูลโต๊ะ</p>;

    return (
        <div className={styles.container}>
            <TableInfo table={table} selectedMemberIndex={selectedMemberIndex}/>
            <Members table={table} onSelectMember={setSelectedMemberIndex}/>
            <div className={styles.infoDiv}>
                <img src="/info.svg"/>
                <p>สมาชิกทุกท่านจะมีเวลาในการสั่งอาหาร 5 นาที หากทุกท่านไม่ทำการสั่งอาหารภายใน 5 นาที จะถือว่าสละสิทธิ์</p>
            </div>
            <button className={styles.createReserveBt}>
                เชิญเพื่อนและเริ่มสั่งอาหาร
                <img src="/Arrow_Right_MD.svg"/>
            </button>
        </div>
    );
}

type Table = {
    id: UUID;
    row: string;
    col: string;
    max_seats: number;
    status: string;
    reserved_seats: number;
};

function TableInfo({ table, selectedMemberIndex }: { table: Table, selectedMemberIndex: number }) {
    const [allowOthers, setAllowOthers] = useState(false);

    const occupied = selectedMemberIndex;
    const min_allow = table.max_seats * 0.8;
    const canClick = occupied >= min_allow;
    const isChecked = !canClick ? true : allowOthers;

    return (
        <div>
            <div className={styles.tableInfoCon}>
                <p>โต๊ะของคุณ : {table.row}{table.col}</p>
                <TableIcon table={table} occupied={occupied}/>
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
    const available = table.status === "available";

    return (
        <div
        className={`${styles.tableIcon} ${
            available ? styles.tableAvailable : styles.tableUnavailable
        }`}
        >
            <img src={available ? "./table_layout_aval.svg" : "./table_layout_notaval.svg"}/>
            <p>{occupied}/{table.max_seats}</p>
        </div>
    );
}

type MemberInfo = {
    username: string;
    first_name: string;
};

interface MembersProps {
    table: Table;
    onSelectMember: (memberNumber: number) => void;
}

function Members({ table, onSelectMember }: MembersProps) {
    const INITIAL_MEMBERS = Array.from({ length: 2 }, () => ({
        username: "",
        first_name: "",
    }));

    const [members, setMembers] = useState<MemberInfo[]>(INITIAL_MEMBERS);
    const [myProfile, setMyProfile] = useState<MemberInfo | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [token, setToken] = useState<string | null>(null);
    
    useEffect(() => {
        if (typeof window === "undefined") return;

        const storedToken = localStorage.getItem("token");
        setToken(storedToken);

        if (!storedToken) {
        setError("ไม่พบโทเค็น กรุณาเข้าสู่ระบบใหม่");
        return;
        }

        const fetchMyProfile = async () => {
        try {
            const res = await fetch("http://localhost:8080/customer/profile", {
                headers: { Authorization: `Bearer ${storedToken}` },
            });
            if (!res.ok) throw new Error("ไม่สามารถดึงข้อมูลโปรไฟล์ได้");
            const data = await res.json();
            setMyProfile(data);
        } catch (err: any) {
            setError(err.message);
        }
        };

        fetchMyProfile();
    }, []); 

    const handleRemoveMember = (index: number) => {
        setMembers((prev) => prev.filter((_, i) => i !== index));
    };

    const handleChangeMember = (index: number, field: keyof MemberInfo, value: string) => {
        const updated = [...members];
        updated[index][field] = value;
        setMembers(updated);
    };

    const handleUsernameBlur = async (index: number) => {
        const username = members[index].username.trim();
        if (!username) return;

        const isDuplicate = members.some((m, i) => i !== index && m.username === username);
        if (isDuplicate) {
            console.log("username ซ้ำ, ไม่ต้องทำอะไร");
            return;
        }

        try {
        const res = await fetch("http://localhost:8080/customer/firstname", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                ...(token ? { Authorization: `Bearer ${token}` } : {}),
            },
            body: JSON.stringify({ username }),
        });

        if (!res.ok) throw new Error("ไม่พบผู้ใช้");
        const data = await res.json();

        handleChangeMember(index, "first_name", data.first_name || "");
        } catch (err) {
            console.error("fetch username error:", err);
        }
    };

    if (error) return <p>เกิดข้อผิดพลาด: {error}</p>;
    if (!myProfile) return <p>กำลังโหลดข้อมูล...</p>;

    return (
        <div className={styles.membersCon}>
            <div>
                <div className={styles.formTitleWrapper}>
                    <div className={styles.titleLine}></div>
                    <h1>สมาชิกที่เข้าร่วม</h1>
                    <div className={styles.titleLine}></div>
                </div>
                <label className={styles.sectionLabel}>คุณ :</label>
                <div className={styles.myInputCon}>
                    <input
                    className={styles.disabInput}
                    type="text"
                    value={`@ ${myProfile.username}`}
                    disabled
                    />
                    <input
                    className={styles.disabInput}
                    type="text"
                    value={`ชื่อ : ${myProfile.first_name}`}
                    disabled
                    />
                </div>
            </div>

            {members.map((m, i) => (
                <div key={i} className={styles.formSection}>
                    <div className={styles.labelRow}>
                        <label className={styles.sectionLabel}>สมาชิกคนที่ {i + 2}</label>
                        {i >= 2 && (
                            <button
                            className={styles.removeUserBt}
                            onClick={() => {
                                handleRemoveMember(i);
                                onSelectMember(members.length);
                            }}
                            type="button"
                            >
                            ลบ
                            </button>
                        )}
                    </div>
                    <div className={styles.inputCon}>
                        <input
                        className={styles.usnInput}
                        type="text"
                        placeholder="@ username"
                        value={m.username ? `@ ${m.username}` : ""}
                        onChange={(e) => {
                            const cleanValue = e.target.value.replace(/^@ ?/, "");
                            handleChangeMember(i, "username", cleanValue);
                        }}
                        onBlur={() => handleUsernameBlur(i)}
                        />
                        <input
                        className={styles.disabInput}
                        type="text"
                        placeholder="ชื่อจะถูกกรอกอัตโนมัติ"
                        value={`ชื่อ : ${m.first_name}`}
                        onChange={(e) => {
                            const cleanValue = e.target.value.replace(/^ชื่อ : ?/, "");
                            handleChangeMember(i, "first_name", cleanValue);
                        }}
                        disabled
                        />
                    </div>
                </div>    
            ))}

            {members.length < table.max_seats - 1 && (
                <button
                className={styles.addUserBt}
                onClick={() => {
                    setMembers([...members, { username: "", first_name: "" }]);
                    onSelectMember(members.length + 2);
                    console.log(myProfile);
                }}
                >
                    <img src="/add_user.svg" />
                    เพิ่มสมาชิก
                </button>
            )}
        </div>
    );
}