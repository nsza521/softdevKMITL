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
  const [username, setUsername] = useState("");
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
        return <ManagePage username={username} />;
      default:
        return <OrderMenu />;
    }
  };
    useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) return;

    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      if (payload.role === "restaurant") {
        setUsername(payload.username);
      }
    } catch {}
  }, []);
  return (
    <div className={`${styles.container} ${notoThai.variable}`}>
      {/* -------- Sidebar -------- */}
      <section className={styles.sidebar}>
        <section className={styles.sidebarsection}>
          <h2>{username || "[ชื่อร้านจ้า]"}</h2>
        </section>

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
            {/* <span className="material-symbols-outlined">info</span> */}
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
  types: MenuType[]; // เพิ่มตรงนี้
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
  const [selectedType, setSelectedType] = useState<string>("All"); // เพิ่ม state กรอง type

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
  const filteredItems = data?.items.filter(item => {
    if (selectedType === "All") return true;
    return item.types.some(t => t.type === selectedType);
  });
 return (
    <section className={styles.shopcontent}>
      <div className={styles.shophead}>
        <div className={styles.restaurantname}>
          <div>
          <h2>Welcome To {username || "[ชื่อร้านจ้า]"}</h2>
          {/* <button><span className="material-symbols-outlined">edit</span></button> */}
        </div>
        <div></div>
        </div>
        <section className={styles.category}>
          <section className={styles.all}>
            <button
              className={selectedType === "All" ? styles.activeTypeBtn : ""}
              onClick={() => setSelectedType("All")}
            >
              All
            </button>
          </section>

          <section className={styles.cate}>
            {types.length > 0 ? types.map((type) => (
              <button
                key={type.id}
                className={selectedType === type.type ? styles.activeTypeBtn : ""}
                onClick={() => setSelectedType(type.type)}
              >
                {type.type}
              </button>
            )) : <p>ไม่มีประเภทเมนู</p>}
          </section>
        </section>
      </div>

      <div className={styles.s_content_detail}>
        {error && <p style={{ color: "red" }}>{error}</p>}
        {!data && !error && <p>⌛ กำลังโหลดเมนู...</p>}
        {filteredItems && filteredItems.map(item => (
          <div key={item.id} className={styles.menu}>
            <div className={styles.menuimg}>
              {item.menu_pic && <img src={item.menu_pic} alt={item.name} />}
              <button className={styles.editBtn}>
                <span className="material-symbols-outlined">info</span>
              </button>
            </div>
            <div className={styles.menudetail}>
              <p className={styles.price}>฿{item.price}</p>
              <p>{item.name}</p>
              <p className={styles.description}>{item.description}</p>
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
function ManagePage({ username }: { username: string }) {
  const [mode, setMode] = useState<"add" | "manage">("manage");
  const [menuList, setMenuList] = useState<MenuItem[]>([]);
  const [types, setTypes] = useState<MenuType[]>([]);
  const [selectedType, setSelectedType] = useState<string>("All");

  // สำหรับ add menu
  const [name, setName] = useState("");
  const [price, setPrice] = useState<number | "">("");
  const [timeTaken, setTimeTaken] = useState<number | "">("");
  const [description, setDescription] = useState("");
  const [menuPic, setMenuPic] = useState<File | null>(null);
  const [selectedTypes, setSelectedTypes] = useState<string[]>([]);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [editingItemId, setEditingItemId] = useState<string | null>(null);
  const [editFile, setEditFile] = useState<File | null>(null);

  const token = localStorage.getItem("token");
  const restaurantID = token ? JSON.parse(atob(token.split('.')[1])).user_id : null;

  useEffect(() => {
    if (!restaurantID) return;

    fetch(`http://localhost:8080/restaurant/menu/${restaurantID}/types`, {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then(res => res.json())
      .then(json => setTypes(Array.isArray(json.types) ? json.types : []))
      .catch(console.error);

    fetch(`http://localhost:8080/restaurant/menu/${restaurantID}/items`, {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then(res => res.json())
      .then(json => setMenuList(json.items || []))
      .catch(console.error);
  }, [restaurantID]);

  const handleAddMenu = async () => {
    if (!name || !price || !timeTaken || selectedTypes.length === 0) {
      setError("กรอกข้อมูลให้ครบก่อน");
      return;
    }

    try {
      const bodyData = { name, price, time_taken: timeTaken, description, menu_type_ids: selectedTypes };
      const res = await fetch(`http://localhost:8080/restaurant/menu/${restaurantID}/items`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify(bodyData),
      });

      if (!res.ok) throw new Error(await res.text());
      const menuItem = await res.json();

      if (menuPic) {
        const formData = new FormData();
        formData.append("menu_pic", menuPic);
        await fetch(
          `http://localhost:8080/restaurant/menu/items/${menuItem.id}/upload_pic`,
          { method: "POST", headers: { Authorization: `Bearer ${token}` }, body: formData }
        );
      }

      setSuccess("✅ เพิ่มเมนูสำเร็จ");
      setError("");
      setMenuList(prev => [...prev, menuItem]);
      setName(""); setPrice(""); setTimeTaken(""); setDescription(""); setMenuPic(null); setSelectedTypes([]);
    } catch (err) {
      console.error(err);
      setError("❌ เพิ่มเมนูไม่สำเร็จ");
      setSuccess("");
    }
  };

  const handleEditMenuPic = async (menuItemId: string) => {
    if (!editFile) return alert("กรุณาเลือกไฟล์ก่อน");
    try {
      const formData = new FormData();
      formData.append("menu_item_picture", editFile);

      const res = await fetch(
        `http://localhost:8080/restaurant/menu/items/${menuItemId}/upload_pic`,
        { method: "POST", headers: { Authorization: `Bearer ${token}` }, body: formData }
      );
      if (!res.ok) throw new Error(await res.text());

      alert("อัปโหลดรูปเรียบร้อย!");
      setEditFile(null);
      setEditingItemId(null);

      // refresh list
      const newRes = await fetch(`http://localhost:8080/restaurant/menu/${restaurantID}/items`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      const json = await newRes.json();
      setMenuList(json.items || []);
    } catch (err) {
      console.error(err);
      alert("❌ อัปโหลดไม่สำเร็จ");
    }
  };

  const filteredItems = menuList.filter(item => {
    if (selectedType === "All") return true;
    return item.types?.some(t => t.type === selectedType);
  });

  return (
    <section className={styles.shopcontent2}>
      {/* header เหมือน order */}
      
      <div className={styles.shophead}>
        <div className={styles.restaurantname}>
          <div>
          <h2>Welcome To {username || "[ชื่อร้านจ้า]"}</h2>
          {/* <button><span className="material-symbols-outlined">edit</span></button> */}
        </div>
        <div></div>
        </div>

        <section className={styles.category}>
          <section className={styles.all}>
            <button onClick={() => setSelectedType("All")}>All</button>
          </section>
          <section className={styles.cate}>
            {types.map(t => (
              <button
                key={t.id}
                onClick={() => setSelectedType(t.type)}
                className={selectedType === t.type ? styles.activeTypeBtn : ""}
              >
                {t.type}
              </button>
            ))}
          </section>
        </section>
      </div>

      {/* ปุ่มเปลี่ยนโหมด */}
      <div style={{ margin: "20px 0", display: "flex", gap: "10px" }}>
        <button onClick={() => setMode("manage")}>จัดการเมนูเดิม</button>
        <button onClick={() => setMode("add")}>เพิ่มเมนูใหม่</button>
      </div>

      {/* เนื้อหา */}
      <div className={styles.s2_content_detail}>
        {mode === "add" ? (
          <div className={styles.addform}>
            <input placeholder="ชื่อเมนู" value={name} onChange={e => setName(e.target.value)} />
            <input type="number" placeholder="ราคา" value={price} onChange={e => setPrice(Number(e.target.value))} />
            <input type="number" placeholder="เวลา (นาที)" value={timeTaken} onChange={e => setTimeTaken(Number(e.target.value))} />
            <textarea placeholder="รายละเอียด" value={description} onChange={e => setDescription(e.target.value)} />
            <input type="file" onChange={e => e.target.files && setMenuPic(e.target.files[0])} />

            <div>
              {types.map(t => (
                <label key={t.id} style={{ marginRight: "10px" }}>
                  <input
                    type="checkbox"
                    value={t.id}
                    checked={selectedTypes.includes(t.id)}
                    onChange={e => {
                      const id = e.target.value;
                      setSelectedTypes(prev => prev.includes(id) ? prev.filter(x => x !== id) : [...prev, id]);
                    }}
                  />
                  {t.type}
                </label>
              ))}
            </div>

            <button onClick={handleAddMenu}>ยืนยัน</button>
            {error && <p style={{ color: "red" }}>{error}</p>}
            {success && <p style={{ color: "green" }}>{success}</p>}
          </div>
        ) : (
          <div className={styles.menuList}>
            {filteredItems.length === 0 ? <p>ไม่มีเมนู</p> : filteredItems.map(item => (
              <div key={item.id} className={styles.menu}>
                <div className={styles.menuimg}>
                  {item.menu_pic && <img src={item.menu_pic} alt={item.name} />}
                  <button className={styles.editBtn} onClick={() => setEditingItemId(item.id)}>
                    <span className="material-symbols-outlined">edit</span>
                  </button>
                </div>
                <div className={styles.menudetail}>
                  <p className={styles.price}>฿{item.price}</p>
                  <p>{item.name}</p>
                  <p className={styles.description}>{item.description}</p>
                </div>

                {editingItemId === item.id && (
                  <div className={styles.editSection}>
                    <input type="file" onChange={e => e.target.files && setEditFile(e.target.files[0])} />
                    <button onClick={() => handleEditMenuPic(item.id)}>อัปโหลด</button>
                    <button onClick={() => setEditingItemId(null)}>ยกเลิก</button>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </section>
  );
}
