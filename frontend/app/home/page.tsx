"use client";

import { useState , useEffect } from "react";
import styles from "./home.module.css";
import { Noto_Sans_Thai } from "next/font/google";
import { useRouter } from "next/navigation"; 

const notoThai = Noto_Sans_Thai({
  subsets: ["thai"],
  weight: ["400", "700"],
  variable: "--font-noto-thai",
});

type UUID = string;

type TableTimeSlot = {
  id: UUID;
  row: string;
  col: string;
  max_seats: number;
  status: string;
  reserved_seats: number;
};

export default function HomePage() {
  const [showPopup, setShowPopup] = useState(false);
  const [profile, setProfile] = useState<any>(null);
  const [tableTimeSlots, setTableTimeSlots] = useState<TableTimeSlot[]>([]);
  const router = useRouter(); 

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const token = localStorage.getItem("token");
        if (!token) return;

        const res = await fetch("http://localhost:8080/customer/profile", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (!res.ok) throw new Error("Failed to fetch profile");

        const data = await res.json();
        console.log("📌 Profile data:", data); // <<
        setProfile(data);
      } catch (err) {
        console.error("❌ Fetch profile error:", err);
      }
    };

    fetchProfile();
    }, []);

    useEffect(() => {
    const fetchTableSlots = async () => {
      try {
        const res = await fetch("http://localhost:8080/table/table_timeslot/now");

        if (!res.ok) throw new Error("Failed to fetch table timeslots");

        const data = await res.json();

        const formatted = data.table_timeslots.map((t: any) => ({
          id: t.id,
          row: t.table_row,
          col: t.table_col,
          max_seats: t.max_seats,
          status: t.status,
          reserved_seats: t.reserved_seats,
        }));

        setTableTimeSlots(formatted);
      } catch (err) {
        console.error("❌ Error fetching table slots:", err);
      }
    };

    fetchTableSlots();
  }, []);

  const handleLogout = async () => {
    try {
      const token = localStorage.getItem("token"); // ดึง token ที่เก็บไว้ตอน login

      const res = await fetch("http://localhost:8080/customer/logout", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`, // ถ้า backend ต้องการ
        },
      });

      if (!res.ok) {
        throw new Error("Logout failed");
      }

      // เคลียร์ token ทิ้ง
      localStorage.removeItem("token");

      alert("ออกจากระบบเรียบร้อย");
      window.location.href = "/login"; // redirect กลับไปหน้า login

    } catch (err) {
      console.error("❌ Error:", err);
      alert("เกิดข้อผิดพลาดตอนออกจากระบบ");
    }
  };

  return (
    <div className={`${styles.container} ${notoThai.variable}`}>
      <div className={styles.headername}>
        <span className={styles.headernameedit}>
          สวัสดี {profile ? `${profile.first_name} ${profile.last_name}` : "กำลังโหลด..."}
          <button onClick={() => setShowPopup(true)}>
            <img src="/editpencil.svg" width={25} height={25} />
          </button>
        </span>
        <button className={styles.logoutbtn} onClick={handleLogout}>
          <img src="/logout.svg" width={25} height={25} />
          logout
        </button>
      </div>

      <div className={styles.boxs}>
        <span>ยอดเงินคงเหลือ {profile ? `${profile.wallet_balance}` : "กำลังโหลด..."} บาท</span>
        <button onClick={() => router.push("/topup")}>
          <img src="/plus.svg" width={15} height={15} />
          เติมเงิน
        </button>
      </div>

      <div className={styles.boxs}>
        <span className={styles.boxspan}>
          <img src="/qr.svg" width={20} height={20} />
          คิวอาโค้ดของฉัน
        </span>
        <button onClick={() => router.push("/qr")}>
          <img src="/show.svg" width={25} height={25} />
          ดู
        </button>
      </div>

      <div className={styles.headername}>
        สถานะโต๊ะตอนนี้ 
      </div>
      <TableLayout tables={tableTimeSlots} />
      {/* <div className={styles.table}></div> */}
      <button className={styles.tablebtn} onClick={() => router.push("/reserveSelectTime")}>จองโต๊ะ</button>

      {/* Popup */}
      {showPopup && (
        <div className={styles.popupbg} onClick={() => setShowPopup(false)}>
          <div className={styles.popup} onClick={(e) => e.stopPropagation()}>
            <h2>แก้ไขข้อมูลส่วนตัว</h2>
            <form
              className={styles.form}
              onSubmit={(e) => {
                e.preventDefault();
                setShowPopup(false);
              }}
            >
              <div className={styles.formGroup}>
                <label>Name</label>
                <input type="text" name="name" placeholder="กรอกชื่อ" />
              </div>

              <div className={styles.formGroup}>
                <label>Surname</label>
                <input type="text" name="surname" placeholder="กรอกนามสกุล" />
              </div>

              <div className={styles.formGroup}>
                <label>Username</label>
                <input type="text" name="username" placeholder="กรอกชื่อผู้ใช้" />
              </div>

              <div className={styles.formGroup}>
                <label>Password</label>
                <input type="password" name="password" placeholder="กรอกรหัสผ่าน" />
              </div>

              <div className={styles.buttonGroup}>
                <button
                  type="button"
                  className={styles.cancelBtn}
                  onClick={() => setShowPopup(false)}
                >
                  ยกเลิก
                </button>
                <button type="submit" className={styles.submitBtn}>
                  ยืนยัน
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
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
        onClick={() =>
            available &&
            router.push(
            `/reserveFillUsr?random=${encodeURIComponent(false)}&table_timeslot_id=${encodeURIComponent(table.id)}`
            )
        }
        disabled={!available}
        >
            <img src={available ? "./table_layout_aval.svg" : "./table_layout_notaval.svg"}/>
            <p>{table.reserved_seats}/{table.max_seats}</p>
        </button>
    );
}