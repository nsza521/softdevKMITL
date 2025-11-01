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
  const [activePage, setActivePage] = useState("order");
  const [username, setUsername] = useState("");
  const [isOnline, setIsOnline] = useState(true);
  const [selectedMenu, setSelectedMenu] = useState(null);

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) return;
    console.log(token);
    try {
      const payload = JSON.parse(atob(token.split(".")[1]));
      if (payload.role === "restaurant") {
        setUsername(payload.username);
      }
    } catch {}
  }, []);

  const handleToggleStatus = async () => {
    try {
      const token = localStorage.getItem("token");
      if (!token) return alert("❌ ไม่มี token");

      const newStatus = isOnline ? "closed" : "open";
      const res = await fetch(`http://localhost:8080/restaurant/status`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ status: newStatus }),
      });

      if (!res.ok) throw new Error("เปลี่ยนสถานะไม่สำเร็จ");
      setIsOnline(!isOnline);
    } catch (err) {
      console.error("❌ Error:", err);
    }
  };

      
   const renderContent = () => {
    switch (activePage) {
      case "order":
        return <OrderMenu  isOnline={isOnline} onToggleStatus={handleToggleStatus} setSelectedMenu={setSelectedMenu} setActivePage={setActivePage}   />;
      case "queue":
        return <QueuePage />;
      case "sales":
        return <TotalSales username={username}/>;
      case "manage":
        return <ManagePage username={username} isOnline={isOnline} onToggleStatus={handleToggleStatus}  setSelectedMenu={setSelectedMenu} setActivePage={setActivePage} />;
      case "addmenu":
        return <AddmenuPage />;
      case "menuDetail":
        return (
          <MenuDetailPage menu={selectedMenu} onBack={() => setActivePage("order")}/>
        );
      default:
        return <OrderMenu username={username} isOnline={isOnline} onToggleStatus={handleToggleStatus} />;
    }
  };

  return (
    <div className={`${styles.container} ${notoThai.variable}`}>
      {/* Sidebar */}
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

      {/* Main Content */}
      <section className={styles.shopcontent}>{renderContent()}</section>

      {/* ปุ่มลอยสำหรับไปหน้า Add Menu */}
      <button
        className={styles.floatingBtn}
        onClick={() => setActivePage("addmenu")}
      >
        <span className="material-symbols-outlined">add_2</span>
      </button>
    </div>
  );
}
/* -------------------------
   เนื้อหาของแต่ละหน้า
-------------------------- */
interface MenuItem {
  time_taken:number;
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
    const token = localStorage.getItem("token");
    const res = await fetch("http://localhost:8080/restaurant/logout", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
    });
    if (!res.ok) throw new Error("Logout failed");
    localStorage.removeItem("token");
    alert("ออกจากระบบเรียบร้อย");
    window.location.href = "/login";
  } catch (err) {
    console.error("❌ Error:", err);
    alert("เกิดข้อผิดพลาดตอนออกจากระบบ");
  }
};
function OrderMenu({ isOnline, onToggleStatus, setActivePage, setSelectedMenu }: any)  {
  const [types, setTypes] = useState<MenuType[]>([]);
  const [data, setData] = useState<MenuData | null>(null);
  const [error, setError] = useState("");
  const [username, setUsername] = useState<string>("");
  const [selectedType, setSelectedType] = useState<string>("All"); // เพิ่ม state กรอง type
  const [restaurantID, setRestaurantID] = useState<string | null>(null);

  const [restaurantPic, setRestaurantPic] = useState<string>("");
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
        setRestaurantID(jsonPayload.user_id);
        const restaurantID = jsonPayload.user_id;
        
        // -----------------------------------

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
        
         // -----------------------------------
        fetch(`http://localhost:8080/restaurant/get_pic`, {
          method: 'GET',
          headers: { 'Authorization': `Bearer ${token}` }
        })
        .then(res => res.json())
        .then(json => {
          if (json.profile_picture) setRestaurantPic(json.profile_picture);
          console.log("📄 /image data:", json.profile_picture);
        })
        .catch(console.error);
        // -----------------------------------

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
            <img  src={restaurantPic || "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQ3XvhUCa5aaC8-riZfbBSudQ_nfCHJA-lbAw&s"}  alt="" />
            <span className={isOnline ? styles.statusdot : styles.statusdotoff}></span>
            <h2>Welcome To {username || "[ชื่อร้านจ้า]"} 
              <div>
                  <p className={isOnline ? styles.online : styles.offline}>
                    {isOnline ? "ออนไลน์" : "ออฟไลน์"}
                  </p>
                  <label className={styles.switch}>
                  <input
                    type="checkbox"
                    checked={isOnline}
                    onChange={onToggleStatus}
                  />
                  <span className={styles.slider}></span>
                </label>
              </div> </h2>
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
        {filteredItems && filteredItems.map(item => {
            return (
              <div 
                key={item.id} 
                className={styles.menu}
                onClick={async () => {
                  console.log("👉 Clicked item id:", item.id);

                  try {
                    const token = localStorage.getItem("token");
                    const res = await fetch(`http://localhost:8080/restaurant/menu/${restaurantID}/${item.id}/detail`, {
                      headers: { 
                        'Authorization': `Bearer ${token}` 
                      }
                    });
                    if (!res.ok) throw new Error("Failed to fetch menu detail");
                    const data = await res.json();
                    console.log("📦 menu detail:", data);

                    setSelectedMenu(data); // ส่งข้อมูลเต็มไป MenuDetailPage
                    setActivePage("menuDetail");
                  } catch (err) {
                    console.error(err);
                    alert("เกิดข้อผิดพลาดในการโหลดข้อมูลเมนู");
                  }
                }}
              >
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
            );
          })}
      </div>
    </section>
  );
}
function QueuePage() {
  const baseUrl = "http://localhost:8080";
  const [orders, setOrders] = useState<any[]>([]);
  const [filteredOrders, setFilteredOrders] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [current, setCurrent] = useState(0);
  const [activeChannel, setActiveChannel] = useState("walk_in");

  const visibleQueues = 7;
  const half = Math.floor(visibleQueues / 2);

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) {
      setError("❌ ไม่มี token กรุณา login ก่อน");
      setLoading(false);
      return;
    }

    async function fetchQueue() {
      try {
        const res = await fetch(`${baseUrl}/restaurant/order/queue`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        if (!res.ok) throw new Error("ไม่สามารถโหลดข้อมูลได้");
        const data = await res.json();
        setOrders(data.orders || []);
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    }

    fetchQueue();
  }, []);

  useEffect(() => {
    const filtered = orders.filter(o => o.channel === activeChannel);
    setFilteredOrders(filtered);
    setCurrent(0);
  }, [orders, activeChannel]);

  if (loading) return <p>กำลังโหลด...</p>;
  if (error) return <p>เกิดข้อผิดพลาด: {error}</p>;

  const totalQueues = filteredOrders.length;

  const displayQueues = Array.from({ length: visibleQueues }, (_, i) => {
    const index = current - half + i;
    if (index < 0 || index >= totalQueues) return null;
    return index + 1;
  });

  return (
    <div className={styles.queuepagemanagement}>
      {/* 🔹 Header — ไม่หายไม่ว่ามีคิวหรือไม่ */}
      <div className={styles.headerqueue}>
        <button
          className={
            activeChannel === "walk_in" ? styles.activebtn : styles.noactivebtn
          }
          onClick={() => setActiveChannel("walk_in")}
        >
          Walk-in
        </button>
        <button
          className={
            activeChannel === "reservation"
              ? styles.activebtn
              : styles.noactivebtn
          }
          onClick={() => setActiveChannel("reservation")}
        >
          Table
        </button>
      </div>

      {/* 🔹 ถ้าไม่มีคิว */}
      {filteredOrders.length === 0 ? (
        <div className={styles.queueall}>
           <div className={styles.queueno}>
            <p className={styles.activeQueue}>ยังไม่มีคิวในช่อง {activeChannel === "walk_in" ? "Walk-in" : "Reservation"}</p>
          </div>
        </div>
      ) : (
        /* 🔹 ถ้ามีคิวค่อยแสดงส่วนนี้ */
        <div className={styles.queueall}>
          <div className={styles.queueno}>
            {displayQueues.map((q, i) =>
              q ? (
                <button
                  key={q}
                  className={q === current + 1 ? styles.activeQueue : ""}
                  onClick={() => setCurrent(q - 1)}
                >
                  คิวที่ {String(q).padStart(3, "0")}
                  <p>{filteredOrders[current].status}</p>
                </button>
              ) : (
                <button key={`empty-${i}`} className={styles.emptyBtn} disabled />
              )
            )}
          </div>

          <div className={styles.Notesofthisreseve}>
            <p className={styles.description}>
              NOTE : {filteredOrders[current].note}
            </p>
          </div>

          <div className={styles.queuesectiondetail}>
            <div
              className={styles.sliderclickleft}
              onClick={() => setCurrent(prev => Math.max(prev - 1, 0))}
            >
              <span className="material-symbols-outlined">arrow_back_ios</span>
            </div>

            <div className={styles.therealmenudetailed}>
              {filteredOrders[current] && (
                <div key={filteredOrders[current].id} className={styles.order_n}>
                  <div className={styles.imageorderholder}>
                    <img
                      src="https://www.jmthaifood.com/wp-content/uploads/2020/01/%E0%B8%95%E0%B9%89%E0%B8%A1%E0%B8%A2%E0%B8%B3%E0%B8%81%E0%B8%B8%E0%B9%89%E0%B8%87-1.jpg"
                      alt="order"
                    />
                  </div>

                  <div className={styles.detailoforder}>
                    <div className={styles.price2}>
                      <p>฿ {filteredOrders[current].total_amount}</p>
                    </div>

                    {/* 🔹 วนลูปแสดงทุกเมนูในคิวนี้ */}
                    {filteredOrders[current].items.map((item: any, i: number) => (
                      <div key={i} className={styles.menuItem}>
                        <p className={styles.mmmmmenu}>
                          {item.menu_name}
                          {item.time_taken_min && (
                            <span>&nbsp;(&nbsp;{item.time_taken_min} นาที&nbsp;)</span>
                          )}
                        </p>

                        {item.note && (
                          <p className={styles.description}>Note: {item.note}</p>
                        )}

                        <div className={styles.handlerwhateveristhisshit}>
                          {item.options?.map((opt: any, j: number) => (
                            <button key={j}>{opt.option_name}</button>
                          ))}
                        </div>
                      </div>
                    ))}
                  </div>

                  <div className={styles.statusofsomethingidontknow}>
                    <button>
                      {filteredOrders[current].status === "pending"
                        ? "กำลังทำ"
                        : filteredOrders[current].status}
                      <span className="material-symbols-outlined">
                        arrow_drop_down
                      </span>
                    </button>
                    <button>
                      ยกเลิก{" "}
                      <span className="material-symbols-outlined">close_small</span>
                    </button>
                  </div>
                </div>
              )}
            </div>

            <div
              className={styles.sliderclickright}
              onClick={() =>
                setCurrent(prev => Math.min(prev + 1, totalQueues - 1))
              }
            >
              <span className="material-symbols-outlined">
                arrow_forward_ios
              </span>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
function TotalSales({username}:any) {
  const [showMoney, setShowMoney] = useState(true);
  const [activeTab, setActiveTab] = useState("history");

  return (
    <section className={styles.shopcontent}>
        <div className={styles.sectionofcirclemoney}>
              <h2 className={styles.headerstotalsales}>บัญชีของ {username}</h2>

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
              activeTab === "withdraw" ? styles.activeTab : ""
            }`}
            onClick={() => setActiveTab("withdraw")}
          >
            ประวัติการถอนเงิน
          </button>
                    <button
            className={`${styles.tabBtn2} ${
              activeTab === "summary" ? styles.activeTab : ""
            }`}
            onClick={() => setActiveTab("summary")}
          >
            {/* สรุปรายรับ */}
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
function ManagePage({ username, isOnline, onToggleStatus ,setActivePage, setSelectedMenu}: any) {
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
  const [editFile, setEditFile] = useState<File | null>(null);
  
  const token = localStorage.getItem("token");
  const restaurantID = token ? JSON.parse(atob(token.split('.')[1])).user_id : null;
  const [editRestaurantFile, setEditRestaurantFile] = useState<File | null>(null);
  const [editingRestaurant, setEditingRestaurant] = useState(false);

  const [restaurantPic, setRestaurantPic] = useState<string>("");

  //Popup edit 
  const [editingMenu, setEditingMenu] = useState<MenuItem | null>(null);
  const [editName, setEditName] = useState("");
  const [editPrice, setEditPrice] = useState<number | "">("");
  const [editTimeTaken, setEditTimeTaken] = useState<number | "">("");
  const [editDescription, setEditDescription] = useState("");
  const [editSelectedTypes, setEditSelectedTypes] = useState<string[]>([]);
  const [editMenuPic, setEditMenuPic] = useState<File | null>(null);
  useEffect(() => {
    if (!restaurantID) return;

    fetch(`http://localhost:8080/restaurant/menu/${restaurantID}/types`, {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then(res => res.json())
      .then(json => setTypes(Array.isArray(json.types) ? json.types : []))
      .catch(console.error);

    fetch(`http://localhost:8080/restaurant/get_pic`, {
          method: 'GET',
          headers: { 'Authorization': `Bearer ${token}` }
        })
        .then(res => res.json())
        .then(json => {
          if (json.profile_picture) setRestaurantPic(json.profile_picture);
          console.log("📄 /image data:", json.profile_picture);
        })
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
      setError("กรุณากรอกข้อมูลให้ครบ");
      return;
    }

    try {
      setError("");
      const body = {
        name,
        price,
        description,
        time_taken: timeTaken,
        menu_pic: null,
        menu_type_ids: selectedTypes,
      };

      const res = await fetch(
        `http://localhost:8080/restaurant/menu/${restaurantID}/items`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(body),
        }
      );

      if (!res.ok) throw new Error(await res.text());
      const json = await res.json();
      console.log("✅ Add Menu Response:", json);

      // upload picture if provided
      if (menuPic) {
        const formData = new FormData();
        formData.append("menu_item_picture", menuPic);
        await fetch(
          `http://localhost:8080/restaurant/menu/items/${json.id}/upload_pic`,
          { method: "POST", headers: { Authorization: `Bearer ${token}` }, body: formData }
        );
      }

      setSuccess("เพิ่มเมนูสำเร็จ!");
      setName("");
      setPrice("");
      setTimeTaken("");
      setDescription("");
      setMenuPic(null);
      setSelectedTypes([]);
    } catch (err) {
      console.error(err);
      setError("❌ เพิ่มเมนูไม่สำเร็จ");
    }
  };
  const handleAddType = async () => {
  const newType = prompt("กรอกชื่อประเภทอาหารใหม่:");
  if (!newType || newType.trim() === "") return ;

  try {
    const res = await fetch(
      `http://localhost:8080/restaurant/menu/${restaurantID}/types`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ type: newType.trim() }),
      }
    );

    if (!res.ok) throw new Error(await res.text());
    const json = await res.json();

    // อัปเดต state ทันทีโดยไม่ต้อง reload
    setTypes((prev) => [...prev, json]);
  } catch (err) {
    console.error(err);
  }
};
  // ฟังก์ชันลบ type
  const handleDeleteType = async (typeId: string) => {
  if (!restaurantID) return;
  if (!confirm("คุณแน่ใจจะลบประเภทนี้?")) return;

  try {
    const res = await fetch(
      `http://localhost:8080/restaurant/menu/types/${typeId}`,
      {
        method: "DELETE",
        headers: { Authorization: `Bearer ${token}` },
      }
    );

    if (!res.ok) throw new Error(await res.text());

    // ลบออกจาก state ทันที
    setTypes(prev => prev.filter(t => t.id !== typeId));

    // ถ้า type ที่ลบเป็น type ที่เลือกอยู่ ก็เปลี่ยนเป็น "All"
    if (selectedType === types.find(t => t.id === typeId)?.type) {
      setSelectedType("All");
    }

    alert("ลบประเภทเรียบร้อยแล้ว");
  } catch (err) {
    console.error(err);
    alert("❌ ลบประเภทไม่สำเร็จ");
  }
  };

  const filteredItems = menuList.filter(item => {
    if (selectedType === "All") return true;
    return item.types?.some(t => t.type === selectedType);
  });
  const handleEditRestaurantPic = async () => {
  if (!editRestaurantFile) return alert("กรุณาเลือกไฟล์ก่อน");

  try {
    const formData = new FormData();
    formData.append("restaurant_picture", editRestaurantFile);

    const res = await fetch(
      `http://localhost:8080/restaurant/upload_pic`,
      {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
        body: formData,
      }
    );

    if (!res.ok) throw new Error(await res.text());
    alert("อัปโหลดรูปร้านเรียบร้อย!");

    setEditRestaurantFile(null);
    setEditingRestaurant(false);
  } catch (err) {
    console.error(err);
  }
  };
  const openEditPopup = (item: MenuItem) => {
    setEditingMenu(item);
    setEditName(item.name);
    setEditPrice(item.price);
    setEditTimeTaken(item.time_taken || "");
    setEditDescription(item.description);
    setEditSelectedTypes(item.types?.map(t => t.id) || []);
    setEditMenuPic(null);
  };
  const handleEditMenuSubmit = async () => {
  if (!editingMenu) return;

  try {
    const body = {
      name: editName,
      price: editPrice,
      description: editDescription,
      time_taken: editTimeTaken,
      menu_type_ids: editSelectedTypes,
    };

    // PATCH menu item
    const res = await fetch(
      `http://localhost:8080/restaurant/menu/items/${editingMenu.id}`,
      {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(body),
      }
    );
    if (!res.ok) throw new Error(await res.text());

    // upload picture ถ้ามี
    if (editMenuPic) {
      const formData = new FormData();
      formData.append("menu_item_picture", editMenuPic);
      const picRes = await fetch(
        `http://localhost:8080/restaurant/menu/items/${editingMenu.id}/upload_pic`,
        {
          method: "POST",
          headers: { Authorization: `Bearer ${token}` },
          body: formData,
        }
      );
      if (!picRes.ok) throw new Error(await picRes.text());
    }

    alert("แก้ไขเมนูเรียบร้อย!");
    setEditingMenu(null);

    // refresh list
    const newRes = await fetch(
      `http://localhost:8080/restaurant/menu/${restaurantID}/items`,
      { headers: { Authorization: `Bearer ${token}` } }
    );
    const json = await newRes.json();
    setMenuList(json.items || []);
  } catch (err) {
    console.error(err);
    alert("❌ แก้ไขเมนูไม่สำเร็จ");
  }
};
  return (
    <section className={styles.shopcontent2}>
      {/* header เหมือน order */}
      
      <div className={styles.shophead}>
        <div className={styles.restaurantname}>
          <div>
            <img  src={restaurantPic || "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQ3XvhUCa5aaC8-riZfbBSudQ_nfCHJA-lbAw&s"}  alt="" />
            <span className={isOnline ? styles.statusdot : styles.statusdotoff}></span>
            
            <h2>Welcome To {username || "[ชื่อร้านจ้า]"} 
              
              <div>
                   <p className={isOnline ? styles.online : styles.offline}>
                    {isOnline ? "ออนไลน์" : "ออฟไลน์"}
                  </p>
                  <label className={styles.switch}>
                  <input
                    type="checkbox"
                    checked={isOnline}
                    onChange={onToggleStatus}
                  />
                  <span className={styles.slider}></span>
                </label>
              </div> </h2>
              <button onClick={() => setEditingRestaurant(true)}>
              <span className="material-symbols-outlined">edit</span>
            </button>
            {editingRestaurant && (
              <div>
                <input type="file" onChange={e => e.target.files && setEditRestaurantFile(e.target.files[0])} />
                <button onClick={handleEditRestaurantPic}>อัปโหลด</button>
                <button onClick={() => { setEditingRestaurant(false); setEditRestaurantFile(null); }}>ยกเลิก</button>
              </div>
            )}
          {/* <button><span className="material-symbols-outlined">edit</span></button> */}
        </div>
        <div></div>
        </div>

        <section className={styles.category}>
          <section className={styles.all}>
            <button onClick={() => setSelectedType("All")}>All</button>
          </section>
          <section className={styles.cate}>
{types.length > 0 ? types.map((type) => (
  <button
    key={type.id}
    className={`${selectedType === type.type ? styles.activeTypeBtn : ""} ${styles.typeBtnWithDelete}`}
    onClick={() => setSelectedType(type.type)}
    style={{ position: "relative" }} // ทำให้ span position absolute อยู่บนปุ่มนี้
  >
    {type.type}
    {/* ปุ่มกากบาท */}
    <span
      onClick={(e) => {
        e.stopPropagation(); // ป้องกันการ trigger เลือก type
        handleDeleteType(type.id);
      }}
      style={{
        position: "absolute",
        top: "0px",
        right: "0px",
        cursor: "pointer",
        color: "red",
        fontWeight: "800",  
        fontSize: "12px",
      }}
    >
      ✕
    </span>
  </button>
)) : <p>ไม่มีประเภทเมนู</p>}
        </section>

          <span className={`material-symbols-outlined ${styles.addtypeBTN}`} onClick={handleAddType}>add_circle</span>
        </section>
      </div>

      {/* ปุ่มเปลี่ยนโหมด */}
      {/* <div style={{ margin: "20px 0", display: "flex", gap: "10px" }}>
        <button onClick={() => setMode("manage")}>จัดการเมนูเดิม</button>
        <button onClick={() => setMode("add")}>เพิ่มเมนูใหม่</button>
      </div> */}

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
              <div key={item.id} className={styles.menu22}
                //   onClick={async () => {
                //   console.log("👉 Clicked item id:", item.id);

                //   try {
                //     const token = localStorage.getItem("token");
                //     const res = await fetch(`http://localhost:8080/restaurant/menu/${restaurantID}/${item.id}/detail`, {
                //       headers: { 
                //         'Authorization': `Bearer ${token}` 
                //       }
                //     });
                //     if (!res.ok) throw new Error("Failed to fetch menu detail");
                //     const data = await res.json();
                //     console.log("📦 menu detail:", data);

                //     setSelectedMenu(data); 
                //     setActivePage("menuDetail");
                //   } catch (err) {
                //     console.error(err);
                //     alert("เกิดข้อผิดพลาดในการโหลดข้อมูลเมนู");
                //   }
                // }}
              >
                <div className={styles.menuimg}>
                  {item.menu_pic && <img src={item.menu_pic} alt={item.name} />}
                  <button className={styles.editBtn} onClick={() => openEditPopup(item)}>
                    <span className="material-symbols-outlined">edit</span>
                  </button>
                </div>
                <div className={styles.menudetail}>
                  <p className={styles.price}>฿{item.price}</p>
                  <p>{item.name}</p>
                  <p className={styles.description}>{item.description}</p>
                </div>
                {editingMenu && (
                  <div className={styles.popupOverlay}>
                    <div className={styles.popupForm}>
                      <h3>แก้ไขเมนู</h3>
                      <input value={editName} onChange={e => setEditName(e.target.value)} placeholder="ชื่อเมนู" />
                      <input type="number" value={editPrice} onChange={e => setEditPrice(Number(e.target.value))} placeholder="ราคา" />
                      <input type="number" value={editTimeTaken} onChange={e => setEditTimeTaken(Number(e.target.value))} placeholder="เวลา (นาที)" />
                      <textarea value={editDescription} onChange={e => setEditDescription(e.target.value)} placeholder="รายละเอียด" />

                      <div>
                        {/* {types.map(t => (
                          <label key={t.id} style={{ marginRight: "10px" }}>
                            <input
                              type="checkbox"
                              value={t.id}
                              checked={editSelectedTypes.includes(t.id)}
                              onChange={e => {
                                const id = e.target.value;
                                setEditSelectedTypes(prev => prev.includes(id) ? prev.filter(x => x !== id) : [...prev, id]);
                              }}
                            />
                            {t.type}
                          </label>
                        ))} */}
                      </div>

                      <input type="file" onChange={e => e.target.files && setEditMenuPic(e.target.files[0])} />

                      <div className={styles.popupActions}>
                        <button onClick={handleEditMenuSubmit}>บันทึก</button>
                        <button onClick={() => setEditingMenu(null)}>ยกเลิก</button>
                      </div>
                    </div>
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
function AddmenuPage() {
  const [types, setTypes] = useState<MenuType[]>([]);
  const [name, setName] = useState("");
  const [price, setPrice] = useState<number | "">("");
  const [timeTaken, setTimeTaken] = useState<number | "">("");
  const [description, setDescription] = useState("");
  const [menuPic, setMenuPic] = useState<File | null>(null);
  const [selectedTypes, setSelectedTypes] = useState<string[]>([]);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const token = localStorage.getItem("token");
  const restaurantID = token ? JSON.parse(atob(token.split(".")[1])).user_id : null;

  useEffect(() => {
    if (!restaurantID) return;
    fetch(`http://localhost:8080/restaurant/menu/${restaurantID}/types`, {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then((res) => res.json())
      .then((json) => setTypes(Array.isArray(json.types) ? json.types : []))
      .catch(console.error);
  }, [restaurantID]);

  const handleAddMenu = async () => {
    if (!name || !price || !timeTaken || selectedTypes.length === 0) {
      setError("กรุณากรอกข้อมูลให้ครบ");
      return;
    }

    try {
      setError("");
      const body = {
        name,
        price,
        description,
        time_taken: timeTaken,
        menu_pic: null,
        menu_type_ids: selectedTypes,
      };

      const res = await fetch(
        `http://localhost:8080/restaurant/menu/${restaurantID}/items`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(body),
        }
      );

      if (!res.ok) throw new Error(await res.text());
      const json = await res.json();
      console.log("✅ Add Menu Response:", json);

      // upload picture if provided
      if (menuPic) {
        const formData = new FormData();
        formData.append("menu_item_picture", menuPic);
        await fetch(
          `http://localhost:8080/restaurant/menu/items/${json.id}/upload_pic`,
          { method: "POST", headers: { Authorization: `Bearer ${token}` }, body: formData }
        );
      }

      setSuccess("เพิ่มเมนูสำเร็จ!");
      setName("");
      setPrice("");
      setTimeTaken("");
      setDescription("");
      setMenuPic(null);
      setSelectedTypes([]);
    } catch (err) {
      console.error(err);
      setError("❌ เพิ่มเมนูไม่สำเร็จ");
    }
  };

  return (
    <section className={styles.shopcontent2}>
      {/* <h2>เพิ่มเมนูใหม่</h2> */}
      <div className={styles.addform}>
        <section>
          <input type="file"onChange={(e) => e.target.files && setMenuPic(e.target.files[0])}/>
                  <div>
          {types.map((t) => (
            <label key={t.id} style={{ marginRight: "10px" }}>
              <input
                type="checkbox"
                value={t.id}
                checked={selectedTypes.includes(t.id)}
                onChange={(e) => {
                  const id = e.target.value;
                  setSelectedTypes((prev) =>
                    prev.includes(id)
                      ? prev.filter((x) => x !== id)
                      : [...prev, id]
                  );
                }}
              />
              {t.type}
            </label>
          ))}
        </div>
        </section>
        <section className={styles.sectiongapaddmenu}>
          <div className={styles.Contenthandler}>
              <div>
                <p>ชื่ออาหาร : </p>
                <input className={styles.menunameinput} placeholder="ชื่อเมนู"value={name}onChange={(e) => setName(e.target.value)}/>
             </div>
              <div className={styles.numprice}>
                  <p>ราคา : </p> <input type="number"value={price}onChange={(e) => {const value = e.target.value;setPrice(value === "" ? "" : Number(value));}}/> <p> บาท </p>
              </div>
              <div className={styles.numprice}>
                  <p>เวลา : </p><input type="number"value={timeTaken}onChange={(e) => {const value = e.target.value;setTimeTaken(value === "" ? "" : Number(value));}}/> <p> นาที </p>
              </div>
              <div>
                    <p>รายละเอียด : </p>
                    <textarea
                      className={styles.menuadddescriptin}
                      placeholder="รายละเอียด"
                      value={description}
                      onChange={(e) => setDescription(e.target.value)}
                    />
              </div>
          </div>
          <button className={styles.submitBTNaddmenu} onClick={handleAddMenu}>ยืนยัน</button>
        </section>
  
   





        {error && <p style={{ color: "red" }}>{error}</p>}
        {success && <p style={{ color: "green" }}>{success}</p>}
      </div>
    </section>
  );
}
function MenuDetailPage({ menu, onBack }: any) {
  if (!menu) return <p>ไม่พบข้อมูลเมนู</p>;

  return (
    <div className={styles.menuDetailPageWrapper}>
      <button onClick={onBack} className={styles.menuDetailBackBtn}>
        ← กลับ
      </button>

      <div className={styles.menuDetailContainer}>
        <img src={menu.menu_pic || "https://via.placeholder.com/200"} alt={menu.name} />

        <div className={styles.menuDetailInfo}>
          <h2>{menu.name}</h2>
          <p className={styles.menuDetailPrice}>฿{menu.price}</p>
          <p>{menu.description}</p>
          <p>⏱ ใช้เวลา {menu.time_taken} นาที</p>

          <div className={styles.menuDetailTypeList}>
            <h4>ประเภทเมนู:</h4>
            {menu.types?.map((t: any, idx: number) => (
              <span key={`${menu.id}-type-${t.id}-${idx}`} className={styles.menuDetailTypeTag}>{t.name}</span>
            ))}
          </div>

          <div className={styles.menuDetailAddonSection}>
            <h4>🍳 Add-ons (ตัวเลือกเพิ่มเติม)</h4>
            {menu.addons && menu.addons.length > 0 ? (
              menu.addons.map((a: any) => (
                <div key={a.id} className={styles.menuDetailAddonItem}>
                  <p><strong>{a.name}</strong></p>
                  {a.options?.length > 0 && (
                    <div>
                      <ul>
                        {a.options.map((o: any, idxO: number) => (
                        <li key={`${menu.id}-addon-${a.id}-option-${o.id}-${idxO}`}>
                          {o.name} {o.price ? `+฿${o.price}` : ""}
                        </li>
                      ))}
                      </ul>
                    </div>
                  )}
                  <p>Required: {a.required ? "✅" : "❌"}</p>
                  <p>From: {a.from}</p>
                  <p>Max select: {a.max_select}, Min select: {a.min_select}</p>
                  {a.allow_qty && <p>Allow quantity selection</p>}
                  
                </div>
              ))
            ) : (
              <p>ไม่มี Add-on</p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
