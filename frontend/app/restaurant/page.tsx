"use client";

import { useState, useEffect } from "react";
import styles from "./restaurant.module.css";
import { Noto_Sans_Thai } from "next/font/google";

const notoThai = Noto_Sans_Thai({
  subsets: ["thai"],
  weight: ["400", "700"],
  variable: "--font-noto-thai",
});

export default function RestaurantPage() {
  // state สำหรับเก็บว่าหน้าปัจจุบันคืออะไร
  const [activePage, setActivePage] = useState("order");

  // ฟังก์ชันเปลี่ยนหน้า
  const renderContent = () => {
    switch (activePage) {
      case "order":
        return <OrderMenu />;
      case "queue":
        return <QueuePage />;
      case "sales":
        return <TotalSales />;
      case "manage":
        return <ManagePage />;
      default:
        return <OrderMenu />;
    }
  };

  return (
    <div className={`${styles.container} ${notoThai.variable}`}>
      {/* -------- Sidebar -------- */}
      <section className={styles.sidebar}>
        <div className={styles.sidebarsection}>
          <h2>[ชื่อร้านจ้า]</h2>
        </div>

        <div className={styles.sidebarsection}>
          <button onClick={() => setActivePage("order")}>
            <span className="material-symbols-outlined">shopping_cart</span>
            <span>Order Menu</span>
          </button>
        </div>

        <div className={styles.sidebarsection}>
          <button onClick={() => setActivePage("queue")}>
            <span className="material-symbols-outlined">star</span>
            <span>Queue</span>
          </button>
        </div>

        <div className={styles.sidebarsection}>
          <button onClick={() => setActivePage("sales")}>
            <span className="material-symbols-outlined">document_search</span>
            <span>Total Sales</span>
          </button>
        </div>

        <div className={styles.sidebarsection}>
          <button onClick={() => setActivePage("manage")}>
            <span className="material-symbols-outlined">edit</span>
            <span>Manage</span>
          </button>
        </div>

        <div className={styles.sidebarsection} id={styles.logoutbtn}>
          <button onClick={handleLogout}>
            <span className="material-symbols-outlined">logout</span>
            <span>Logout</span>
          </button>
        </div>
      </section>

      {/* -------- Main Content -------- */}
      <section className={styles.shopcontent}>{renderContent()}</section>

      <button className={styles.floatingBtn}>
        <span className="material-symbols-outlined">add_2</span>
      </button>
    </div>
  );
}

/* -------------------------
   เนื้อหาของแต่ละหน้า
-------------------------- */
interface MenuItem {
  id: string;
  name: string;
  price: number;
  description: string;
  menu_pic?: string;
}
interface MenuData {
  items: MenuItem[];
}
interface MenuType {
  id: string;
  restaurant_id: string;
  type: string;
}
const handleLogout = async () => {
    try {
      const token = localStorage.getItem("token"); // ดึง token ที่เก็บไว้ตอน login

      const res = await fetch("http://localhost:8080/user/logout", {
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


function OrderMenu() {
  const [types, setTypes] = useState<MenuType[]>([]);
  const [data, setData] = useState<MenuData | null>(null);
  const [error, setError] = useState("");
  const [username, setUsername] = useState<string>("");

  useEffect(() => {

    const token = localStorage.getItem("token");
    if (!token) {
      setError("❌ ไม่มี token กรุณา login ก่อน");
      return;
    }

    try {
      const payload = token.split('.')[1];
      const base64 = payload.replace(/-/g, '+').replace(/_/g, '/');
      const jsonPayload = JSON.parse(atob(base64));

      if (jsonPayload.role === "restaurant") {
        setUsername(jsonPayload.username); // เอา username มาโชว์
        const restaurantID = jsonPayload.user_id;

        fetch(`http://localhost:8080/restaurant/menu/${restaurantID}/items`, {
          method: 'GET',
          headers: { 'Authorization': `Bearer ${token}` },
        })
        .then(async (res) => {
          const text = await res.text();
          if (!res.ok) throw new Error(text);
          const json = JSON.parse(text);
          setData(json);
          console.log("📄 /menuitem data:", json);
        })
        .catch(err => {
          console.error("❌ Fetch error:", err);
          setError("โหลดข้อมูลไม่สำเร็จ");
        });

        fetch(`http://localhost:8080/restaurant/menu/${restaurantID}/types`, {
        headers: { 'Authorization': `Bearer ${token}` },
      })
        .then(res => res.json())
        .then(json => {
          console.log("📄 /types data:", json); // จะเห็น can_edit และ types
          setTypes(Array.isArray(json.types) ? json.types : []);
        })
        .catch(err => console.error("❌ Fetch /types error:", err));

        
      } else {
        setError("❌ Token ไม่ใช่ร้านอาหาร");
      }
    } catch (err) {
      console.error("❌ JWT decode error:", err);
      setError("Token ไม่ถูกต้อง");
    }
  }, []);

  return (
    <section className={styles.shopcontent}>
      <div className={styles.shophead}>
        <div className={styles.restaurantname}>
          <div>
            <h2>Welcome To {username || "[ชื่อร้านจ้า]"}</h2>
            <button><span className="material-symbols-outlined">edit</span></button>
          </div>
          <div></div>
        </div>
        <section className={styles.category}>
          <section className={styles.all}>
            <button>All</button>
          </section>
          <section className={styles.cate}>
            {types.length > 0 ? types.map((type) => (
              <button key={type.id}>{type.type}</button>
            )) : <p>ไม่มีประเภทเมนู</p>}
          </section>
        </section>
      </div>

      <div className={styles.s_content_detail}>
        {error && <p style={{ color: "red" }}>{error}</p>}
        {!data && !error && <p>⌛ กำลังโหลดเมนู...</p>}
        {data && data.items.map(item => (
          <div key={item.id} className={styles.menu}>
            <div className={styles.menuimg}>
              {item.menu_pic && <img src={item.menu_pic} alt={item.name} />}
            </div>
            <div className={styles.menudetail}>
              <p>฿ {item.price}</p>
              <p>{item.name}</p>
              <p>{item.description}</p>
            </div>
          </div>
        ))}
      </div>
    </section>
  );
}
function QueuePage() {
  return (
    <div>
      <h2>⭐ Queue</h2>
      <p>แสดงคิวของลูกค้าในร้าน</p>
    </div>
  );
}
function TotalSales() {
  const [showMoney, setShowMoney] = useState(true);
  const [activeTab, setActiveTab] = useState("history");

  return (
    <section className={styles.shopcontent}>
        <div className={styles.sectionofcirclemoney}>
              <h2 className={styles.headerstotalsales}>บัญชีของ [ชื่อร้านจ้า]</h2>

            {/* วงกลมยอดเงิน */}
            <div className={styles.moneyCircle}>
                <p className={styles.subText}>ยอดเงินคงเหลือ</p>

                <h1 className={styles.totalAmount}>
                {showMoney ? "12,540.75 ฿" : "********"}
                </h1>

                <button
                className={styles.eyeButton}
                onClick={() => setShowMoney(!showMoney)}
                >
                <span className="material-symbols-outlined">
                    {showMoney ? "visibility" : "visibility_off"}
                </span>
                </button>
            </div>
        </div>

      <button className={styles.withdrawButton}>ยื่นคำขอถอนเงิน</button>

      {/* footer ภายใน section */}
      <div className={styles.footerSection}>
        {/* ปุ่มแท็บ */}
        <div className={styles.tabButtons}>
          <button
            className={`${styles.tabBtn} ${
              activeTab === "history" ? styles.activeTab : ""
            }`}
            onClick={() => setActiveTab("history")}
          >
            รายการย้อนหลัง
          </button>

          <button
            className={`${styles.tabBtn} ${
              activeTab === "summary" ? styles.activeTab : ""
            }`}
            onClick={() => setActiveTab("summary")}
          >
            สรุปรายรับ
          </button>

          <button
            className={`${styles.tabBtn} ${
              activeTab === "withdraw" ? styles.activeTab : ""
            }`}
            onClick={() => setActiveTab("withdraw")}
          >
            ประวัติการถอนเงิน
          </button>
        </div>

        {/* เนื้อหาแท็บ */}
        <div className={styles.tabContent}>
          {activeTab === "history" && <p>📜 รายการย้อนหลังของร้านทั้งหมด</p>}
          {activeTab === "summary" && <p>📊 สรุปรายรับรายวัน / เดือน</p>}
          {activeTab === "withdraw" && 
          <div className={styles.historywithdrawflex}>
            <div>สิงหาคม 2568 ▾</div>
            <div>
                <p>dd mm yy hh:mm -xxx,xxx,xxx กำลังดำเนินการ</p>
                <p>dd mm yy hh:mm -xxx,xxx,xxx กำลังดำเนินการ</p>
                <p>dd mm yy hh:mm -xxx,xxx,xxx กำลังดำเนินการ</p>
                <p>dd mm yy hh:mm -xxx,xxx,xxx กำลังดำเนินการ</p>
            </div>
          </div>
          }
        </div>  
      </div>
    </section>
  );
}

function ManagePage() {
  return (
    <div>
      <h2>🛠️ Manage</h2>
      <p>หน้าแก้ไขข้อมูลร้าน เมนู ราคา ฯลฯ</p>
    </div>
  );
}
