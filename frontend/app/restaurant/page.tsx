"use client";

import { useState, useEffect } from "react";
import styles from "./restaurant.module.css";
import { Noto_Sans_Thai } from "next/font/google";
import { useSearchParams } from "next/navigation";
import { useRouter } from "next/navigation";


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
          <MenuDetailPage menu={selectedMenu} onBack={() => setActivePage("manage")}/>
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
  addons: Addon[]
  onAdd?: () => void
}
interface MenuData {
  items: MenuItem[];
}
interface MenuType {
  id: string;
  restaurant_id: string;
  type: string;
}

type UUID = string;

interface CartItem {
  item: MenuItem
  quantity: number
  selectedAddons: Addon[]
}

interface Addon {
  id: UUID
  name: string
  required: boolean
  options: Option[]
}

interface Option {
  id: UUID
  name: string
  price_delta: string
  is_default?: boolean
}

interface CartProps {
  cart: CartItem[]
}

interface MenuPopupProps {
  isOpen: boolean
  onClose: () => void
  item: MenuItem
  cartItem?: CartItem | null
  onAddToCart: (item: MenuItem, quantity: number, selectedAddons: Addon[]) => void
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
  const [cart, setCart] = useState<CartItem[]>([])
  const [isPopupOpen, setIsPopupOpen] = useState(false);
  const [selectedItem, setSelectedItem] = useState<MenuItem | null>(null);


  const [restaurantPic, setRestaurantPic] = useState<string>("");
  
  const handleAddItem = (menuItem: MenuItem) => {
    setSelectedItem(menuItem);
    setIsPopupOpen(true);
  };



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

  const handleAddToCart = (item: MenuItem, quantity: number, selectedAddons: Addon[]) => {
    setCart((prev) => {
      if (quantity === 0) {
        return prev.filter(ci => ci.item.id !== item.id);
      }

      const index = prev.findIndex(ci => ci.item.id === item.id);

      if (index >= 0) {
        const updated = [...prev];
        updated[index] = { item, quantity, selectedAddons };
        return updated;
      }

      return [...prev, { item, quantity, selectedAddons }];
    });
  };

  const filteredItems = data?.items.filter(item => {
    if (selectedType === "All") return true;
    return item.types.some(t => t.type === selectedType);
  });
 return (
    <section className={styles.pangContain}>
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
              const cartItem = cart.find(ci => ci.item.id === item.id);
              const quantity = cartItem?.quantity ?? 0;
              return (
                <div key={item.id} className={styles.menu}>
                  
                  <div className={styles.menuimg}>
                    {item.menu_pic && <img src={item.menu_pic} alt={item.name} />}
                  </div>
                  <div className={styles.menuCon}>
                    <div className={styles.menudetail}>
                      <p className={styles.price}>฿{item.price}</p>
                      <p>{item.name}</p>
                      <p className={styles.description}>{item.description}</p>
                    </div>
                      
                    <button
                      className={styles.addBtn}
                      onClick={() => handleAddItem(item)}
                    >
                      {quantity === 0 ? (
                        <img src="/Add_Plus_Circle.svg" />
                      ) : (
                        <span className={styles.cartQtyCircle}>{quantity}</span>
                      )}
                    </button>
                  </div>
                </div>
              );
            })}
            {selectedItem && (
              <MenuPopup
                isOpen={isPopupOpen}
                onClose={() => setIsPopupOpen(false)}
                item={selectedItem}
                cartItem={cart.find(ci => ci.item.id === selectedItem.id) ?? null}
                onAddToCart={handleAddToCart}
              />
            )}
        </div>
      </section>
      <Cart cart={cart} />
    </section>
  );
}
function Cart({ cart }: CartProps) {
  const router = useRouter();
  const searchParam = useSearchParams();
  const restaurant_id = searchParam.get("id") || "";
  // const reservation_id = searchParam.get("reservationId") || "";

  if (cart.length === 0) return null;

  const itemNum = cart.reduce((sum, ci) => sum + ci.quantity, 0);
  const itemCost = cart.reduce((sum, ci) => {
    // ✅ รองรับทั้ง number และ string
    const basePrice = typeof ci.item.price === "number"
      ? ci.item.price
      : parseFloat(ci.item.price);

    const addonTotal =
      ci.selectedAddons?.reduce((addonSum, addon) => {
        const optionTotal = addon.options.reduce((optSum, opt) => {
          const delta = typeof opt.price_delta === "number"
            ? opt.price_delta
            : parseFloat(opt.price_delta);
          return optSum + delta;
        }, 0);

        return addonSum + optionTotal;
      }, 0) ?? 0;

    return sum + (basePrice + addonTotal) * ci.quantity;
  }, 0);


  const handleCheckout = async () => {
    const token = localStorage.getItem("token");
    if (!token) {
      alert("Token not found");
      return;
    }

    const items = cart.map(ci => ({
      menu_item_id: ci.item.id,
      quantity: ci.quantity,
      selections: ci.selectedAddons.flatMap(addon =>
        addon.options.map(opt => {
          if (addon.required === false) {
            return {
              group_id: addon.id,
              option_id: opt.id,
              qty: 1
            };
          }

          return {
            group_id: addon.id,
            option_id: opt.id
          };
        })
      )
    }));

    const body = {
      items
    };

    try {
      const res = await fetch("http://localhost:8080/restaurant/order", {
        method: "POST",
        headers: {
          "Authorization": `Bearer ${token}`,
          "Content-Type": "application/json"
        },
        body: JSON.stringify(body)
      });

      if (!res.ok) {
        console.error(await res.text());
        alert("ส่งคำสั่งซื้อไม่สำเร็จ");
        return;
      }
      const resp = await res.json();
      console.log("Order response:", resp);
      alert("สั่งอาหารสำเร็จ")
      window.location.reload();

    } catch (err) {
      console.error(err);
      alert("เกิดข้อผิดพลาด");
    }
  };

return (
  <div className={styles.cartBox}>
    <h3 className={styles.cartHeader}>My Order</h3>

    <div className={styles.orderList}>
      {cart.map((ci) => {
        const basePrice = typeof ci.item.price === "number"
          ? ci.item.price
          : parseFloat(ci.item.price);

        const addonTotal = ci.selectedAddons.reduce((sum, addon) => {
          return sum + addon.options.reduce((optSum, opt) => {
            const delta = typeof opt.price_delta === "number"
              ? opt.price_delta
              : parseFloat(opt.price_delta);
            return optSum + delta;
          }, 0);
        }, 0);

        const finalPrice = (basePrice + addonTotal) * ci.quantity;

        return (
          <div key={ci.item.id} className={styles.orderItem}>
            <div className={styles.qtyBox}>{ci.quantity}</div>

            <div className={styles.menuInfo}>
              <p className={styles.menuName}>{ci.item.name}</p>

              <p className={styles.addonText}>
                {ci.selectedAddons
                  .flatMap(addon => addon.options.map(opt => opt.name))
                  .join(", ")}
              </p>

            </div>

            <p className={styles.menuPrice}>฿ {finalPrice}</p>
          </div>
        );
      })}
    </div>

    <div className={styles.summaryBox}>
      <p className={styles.totalText}>รวมทั้งหมด</p>
      <p className={styles.totalAmount}>฿ {itemCost}</p>
    </div>

    <button className={styles.payBtn} onClick={handleCheckout}>
      ชำระเงิน
    </button>
  </div>
);

}
function MenuPopup({ isOpen, onClose, item, cartItem, onAddToCart }: MenuPopupProps) {
  const searchParam = useSearchParams();
  const restaurant_id = searchParam.get("id") || "";

  const [quantity, setQuantity] = useState(1);
  const [selectedAddons, setSelectedAddons] = useState<Addon[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [detail, setDetail] = useState<MenuItem | null>(null);

  useEffect(() => {
    const storedToken = localStorage.getItem("token");
    if (!storedToken || !isOpen) return;

    async function fetchDetail() {
      try {
        const res = await fetch(
          `http://localhost:8080/restaurant/menu/${restaurant_id}/${item.id}/detail`,
          { headers: { Authorization: `Bearer ${storedToken}` } }
        )
        if (!res.ok) throw new Error("Failed to fetch detail");

        const data: MenuItem = await res.json();
        setDetail(data);

        // ถ้ามี cartItem ให้ใช้ cartItem.selectedAddons ไม่ใช้ default
        if (!cartItem) {
          const init = data.addons?.map(addon => ({
            ...addon,
            options: addon.options.filter(opt => !!opt.is_default)
          })) ?? []
          setSelectedAddons(init)
        }
        // ถ้ามี cartItem จะใช้ useEffect ตัวที่ handle cartItem อยู่แล้ว
      } catch (err) {
        console.error(err)
        setError("โหลดรายละเอียดเมนูไม่สำเร็จ");
      }
    }

    fetchDetail()
  }, [isOpen, item.id, cartItem])

  useEffect(() => {
    if (!isOpen) return;

    if (cartItem) {
      setQuantity(cartItem.quantity)
      setSelectedAddons(cartItem.selectedAddons.map(a => ({
        ...a,
        options: [...a.options]
      })))
    } else {
      setQuantity(1)
      setSelectedAddons([])
    }
  }, [isOpen, cartItem])

  const isInCart = !!cartItem
  const isRemove = isInCart && quantity === 0
  const buttonLabel = !isInCart
    ? "เพิ่มลงตะกร้า"
    : isRemove
      ? "นำออก"
      : "อัปเดตตะกร้า"

  const handleAdd = () => {
    onAddToCart(item, isRemove ? 0 : quantity, selectedAddons)
    onClose()
    // console.log("Added to cart:", item, quantity, selectedAddons);
  }

  if (error) return <p style={{ color: "red" }}>{error}</p>
  if (!isOpen || !detail) return null

  return (
    <div className={styles.popupOverlay}>
      <div className={styles.popupCon}>
        <button className={styles.closeBt} onClick={onClose}>
          <img src="/Close_LG.svg" />
        </button>

        <img src={detail.menu_pic ?? "/placeholder.png"} className={styles.menuImg} />

        <div className={styles.popupOptionsCon}>
          <h3>
            {detail.addons?.length > 0
              ? `ตัวเลือกเพิ่มเติมสำหรับ ${detail.name}`
              : detail.name}
          </h3>
          <div className={styles.optionsList}>
            {detail.addons?.map((addon) => (
              <div key={addon.id} className={styles.addonGroup}>
                <p className={styles.addonTitle}>{addon.name}{addon.required ? "" : " [optional]"}</p>

                {addon.options.map(opt => {
                  const isSelected = selectedAddons
                    .find(a => a.id === addon.id)
                    ?.options.some(o => o.id === opt.id)

                  return (
                    <label key={opt.id}>
                      <input
                        type={addon.required ? "radio" : "checkbox"}
                        name={addon.required ? addon.id : opt.id}
                        checked={!!isSelected}
                        onChange={() => {
                          setSelectedAddons(prev => prev.map(a => {
                            if (a.id !== addon.id) return a
                            if (addon.required) {
                              return { ...a, options: [opt] } // radio
                            } else {
                              const exists = a.options.find(o => o.id === opt.id)
                              const newOptions = exists
                                ? a.options.filter(o => o.id !== opt.id)
                                : [...a.options, opt]
                              return { ...a, options: newOptions }
                            }
                          }))
                        }}
                      />
                      <div className={styles.addonOptionRow}>
                        <span>{opt.name}</span>
                        <span>{opt.price_delta} ฿</span>
                      </div>

                    </label>
                  )
                })}
              </div>
            ))}
          </div>
        </div>

        <div className={styles.popupBottomCon}>
          <div className={styles.quantityCon}>
            <button className={styles.bt} onClick={() => setQuantity(q => Math.max(0, q - 1))}>
              <img src="/Remove_Minus.svg" />
            </button>
            <span>{quantity}</span>
            <button className={styles.bt} onClick={() => setQuantity(q => q + 1)}>
              <img src="/Add_Plus.svg" />
            </button>
          </div>

          <button className={styles.addCartBt} onClick={handleAdd}>
            <img src="/Shopping_Cart_White.svg" />
            {buttonLabel}
          </button>
        </div>
      </div>
    </div>
  )
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
        console.log("somethingwhateveridontknowfuckmaybethisisqueue",data);
        setOrders(data.orders || []);
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    }

    fetchQueue();
  }, []);





const updateOrderStatus = async (orderId: string, newStatus: string) => {
  const token = localStorage.getItem("token");
  if (!token) return alert("❌ ไม่มี token");

  console.log("🛰️ updateOrderStatus ->", `${baseUrl}/restaurant/order/orders/${orderId}/status`, "status:", newStatus);

  try {
    const res = await fetch(`${baseUrl}/restaurant/order/orders/${orderId}/status`, {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify({ status: newStatus }),
    });

    if (!res.ok) {
      const text = await res.text();
      console.error("❌ Backend response:", text);
      throw new Error("ไม่สามารถอัปเดตสถานะได้");
    }

    setOrders(prev =>
      prev.map(o => (o.id === orderId ? { ...o, status: newStatus } : o))
    );

  } catch (err) {
    console.error("🔥 updateOrderStatus error:", err);
    alert("❌ ไม่สามารถอัปเดตสถานะได้");
  }
};





  
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
    return index; // <-- เก็บ index แทน
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
            {displayQueues.map((idx, i) =>
              idx !== null ? (
                  <button
                      key={filteredOrders[idx].id}
                      className={idx === current ? styles.activeQueue : styles.activeQueue2}
                      onClick={() => setCurrent(idx)}
                  >
                      คิวที่ {String(idx + 1).padStart(3, "0")}
                      <p></p>
                          <select className={styles.selectofstauts}
                                value={filteredOrders[idx].status}
                           onChange={(e) => updateOrderStatus(filteredOrders[idx].id, e.target.value)}>
                            <option value="paid">Paid</option>
                            <option value="cancelled">Cancelled</option>
                            <option value="served">Served</option>
                      </select>
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
  {filteredOrders[current] && 
    filteredOrders[current].items.map((item: any, i: number) => (
      <div key={i} className={styles.order_n}>
        {/* รูป */}
        <div className={styles.imageorderholder}>
          <img
            src={item.menu_pic}
            alt="order"
          />
        </div>

        {/* รายละเอียดเมนู */}
        <div className={styles.detailoforder}>
          <div className={styles.price2}>
            <p>฿ {filteredOrders[current].total_amount}</p>
          </div>

          <div className={styles.menuItem}>
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
        </div>

        {/* สถานะรวม */}
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
    ))
  }
</div>

                      {/* dfjdshisaodpsadlpadposa */}
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
function TotalSales({ username }: any) {
  const [showMoney, setShowMoney] = useState(true);
  const [activeTab, setActiveTab] = useState("history");
  const [balance, setBalance] = useState<number | null>(null);
  const [transactions, setTransactions] = useState<any[]>([]);
    const [orders, setOrders] = useState<any[]>([]);
  // ✅ popup state
  const [showPopupoftiHisButtonIsAmazaing, setShowPopupoftiHisButtonIsAmazaing] = useState(false);
  const [withdrawData, setWithdrawData] = useState({
    full_name: "",
    bank_name: "KBANK",
    bank_account_number: "",
    withdraw_amount: "",
  });

  const token = localStorage.getItem("token");
  console.log("token",token);
  const handleChange = (e: any) => {
    const { name, value } = e.target;
    setWithdrawData((prev) => ({ ...prev, [name]: value }));
  };

    useEffect(() => {
    if (!token) return;

    const fetchBalance = async () => {
      try {
        const res = await fetch("http://localhost:8080/restaurant/balance", {
          headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`,
          },
        });
        if (!res.ok) throw new Error("ไม่สามารถโหลดยอดเงินได้");
        const data = await res.json();
        setBalance(data.balance); // สมมติ API คืน { balance: 12540.75 }
      } catch (err) {
        console.error("❌ Fetch balance error:", err);
      }
    };

    fetchBalance();
  }, [token]);


    useEffect(() => {
    if (!token) return;

    const fetchTransactions = async () => {
      try {
        const res = await fetch("http://localhost:8080/payment/transaction/all", {
          headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`,
          },
        });
        if (!res.ok) throw new Error("ไม่สามารถโหลดรายการธุรกรรมได้");
        const data = await res.json();
        setTransactions(data.transactions || []);
      } catch (err) {
        console.error("❌ Fetch transactions error:", err);
      }
      };
      fetchTransactions();
    }, [token]);  


  const handleWithdraw = async () => {
    if (!withdrawData.full_name || !withdrawData.bank_account_number || !withdrawData.withdraw_amount) {
      alert("กรอกข้อมูลให้ครบก่อนครับ");
      return;
    }

    try {
      const res = await fetch("http://localhost:8080/payment/withdraw/wallet", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${token}`,
        },
        body: JSON.stringify({
          ...withdrawData,
          withdraw_amount: Number(withdrawData.withdraw_amount),
        }),
      });

      const data = await res.json();
      console.log("📦 Withdraw response:", data);

      if (!res.ok) {
        alert("ถอนเงินไม่สำเร็จ: " + (data.message || "เกิดข้อผิดพลาด"));
      } else {
        alert("ถอนเงินสำเร็จ!");
        setShowPopupoftiHisButtonIsAmazaing(false);
        setWithdrawData({
          full_name: "",
          bank_name: "KBANK",
          bank_account_number: "",
          withdraw_amount: "",
        });
      }
    } catch (err) {
      console.error("❌ Error:", err);
      alert("เกิดข้อผิดพลาดในการเชื่อมต่อเซิร์ฟเวอร์");
    }
  };


   useEffect(() => {
    if (!token) return;

    const fetchOrders = async () => {
      try {
        const res = await fetch(`http://localhost:8080/restaurant/order/history?date=2025-11-03`, {
          headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`,
          },
        });
        if (!res.ok) throw new Error("ไม่สามารถโหลดรายการสั่งอาหารได้");
        const data = await res.json();
        setOrders(data.orders || []);
      } catch (err) {
        console.error("❌ Fetch orders error:", err);
      }
    };
    fetchOrders();
  }, [token]);
  return (
    <section className={styles.shopcontent}>
      <div className={styles.sectionofcirclemoney}>
        <h2 className={styles.headerstotalsales}>บัญชีของ {username}</h2>

        <div className={styles.moneyCircle}>
          <p className={styles.subText}>ยอดเงินคงเหลือ</p>
          <h1 className={styles.totalAmount}>
             {showMoney
            ? balance !== null
              ? `${balance.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })} ฿`
              : "กำลังโหลด..."
            : "********"} 
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

      {/* ปุ่มยื่นคำขอถอนเงิน */}
      <button
        className={styles.withdrawButton}
        onClick={() => setShowPopupoftiHisButtonIsAmazaing(true)}
      >
        ยื่นคำขอถอนเงิน
      </button>

      {/* popup ถอนเงิน */}
      {showPopupoftiHisButtonIsAmazaing && (
        <div
          className={styles.popupOverlay}
          onClick={() => setShowPopupoftiHisButtonIsAmazaing(false)}
        >
          <div
            className={styles.popupForm}
            onClick={(e) => e.stopPropagation()}
          >
            <h3>ยื่นคำขอถอนเงิน</h3>

            <label>
              ชื่อ-นามสกุล:
              <input
                type="text"
                name="full_name"
                value={withdrawData.full_name}
                onChange={handleChange}
                placeholder="กรอกชื่อผู้ถอน"
                required
              />
            </label>

            <label>
              ธนาคาร:
              <select
                name="bank_name"
                value={withdrawData.bank_name}
                onChange={handleChange}
              >
                <option value="KBANK"> กสิกรไทย (KBANK)</option>
                <option value="SCB">ไทยพาณิชย์ (SCB)</option>
                <option value="BBL">กรุงเทพ (BBL)</option>
                <option value="KTB">กรุงไทย (KTB)</option>
              </select>
            </label>

            <label>
              เลขบัญชี:
              <input
                type="text"
                name="bank_account_number"
                value={withdrawData.bank_account_number}
                onChange={handleChange}
                placeholder="0123456789"
              />
            </label>

            <label>
              จำนวนเงินที่ต้องการถอน:
              <input
                type="number"
                name="withdraw_amount"
                value={withdrawData.withdraw_amount}
                onChange={handleChange}
                placeholder="10"
                min="1"
              />
            </label>
            
            <div className={styles.popupActions}>
              <button className={styles.confirmBtn} onClick={handleWithdraw}>
                ยืนยันถอนเงิน
              </button>
              <button
                className={styles.cancelBtnnnnnn}
                onClick={() => setShowPopupoftiHisButtonIsAmazaing(false)}
              >
                ยกเลิก
              </button>
            </div>
          </div>
        </div>
      )}

      {/* footer section */}
      <div className={styles.footerSection}>
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
          >
            {/* สรุปรายรับ */}
          </button>
        </div>

          <div className={styles.tabContent}>

            {activeTab === "history" && (
              <div className={styles.orderHistory}>
                {orders.length === 0 ? (
                  <p>ยังไม่มีรายการสั่งซื้อ</p>
                ) : (
                  orders.map((order) => (
                    <div key={order.order_id} className={styles.orderCardSSSS}>
                      {/* <h4>Order #{order.order_id.slice(0, 8)}</h4> */}
                        {order.items.map((item: any, idx: number) => (
                          <span key={idx} className={styles.data1ofwhatvevrearasd}>
                            {item.menu_name} x{item.quantity} ({item.options.map((o: any) => o.option_name).join(", ")})
                          </span>
                        ))}
                      <p className={styles.adasdsadsssssssssssssssa}>{new Date(order.order_time).toLocaleString("th-TH")}</p>
                      <p  className={styles.asdasdsadsasassssqqq}>รวม {order.total_amount.toLocaleString()} ฿</p>
                      <ul>

                      </ul>
                    </div>
                  ))
                )}
              </div>
            )}

            {activeTab === "withdraw" && (
              <div className={styles.withdrawHistoryWrapper}>
                {transactions.filter(tx => tx.type === "withdraw").length === 0 ? (
                  <p>ยังไม่มีรายการถอนเงิน</p>
                ) : (
                  
                  transactions
                    .filter(tx => tx.type === "withdraw")
                    .map((tx) => (
                      <div key={tx.transaction_id} className={styles.withdrawItem}>
                        <div className={styles.withdrawDate}>
                          {new Date(tx.created_at).toLocaleDateString("th-TH", {
                            day: "2-digit",
                            month: "long",
                            year: "numeric",
                          })}
                        </div>
                        <div className={styles.withdrawInfo}>
                          <span className={styles.withdrawBank}>
                            ({tx.payment_method})
                          </span>
                          <span className={styles.withdrawTime}>
                            {new Date(tx.created_at).toLocaleTimeString("th-TH", {
                              hour: "2-digit",
                              minute: "2-digit",
                            })}
                          </span>
                          <span className={styles.withdrawAmount}>
                            {tx.amount.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })} ฿
                          </span>

                          <span className={styles.withdrawStatus}>เสร็จสิ้น</span>
                        </div>
                      </div>
                    ))
                )}
              </div>
            )}
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
      // menu_type_ids: editSelectedTypes,
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

                    setSelectedMenu(data); 
                    setActivePage("menuDetail");
                  } catch (err) {
                    console.error(err);
                    alert("เกิดข้อผิดพลาดในการโหลดข้อมูลเมนู");
                  }
                }}
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
                  <div className={styles.popupOverlay2}>
                    <div className={styles.popupForm2}>
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
<div className={styles.asdasdsadsadsadsaaaaaaaa}>
  <div className={styles.imageUploadBox}>
    <label htmlFor="menuPic" className={styles.uploadLabel}>
      {menuPic ? (
        <img
          src={URL.createObjectURL(menuPic)}
          alt="Preview"
          className={styles.previewImage}
        />
      ) : (
        <span className={styles.uploadText}>📷 เลือกรูปเมนู</span>
      )}
    </label>
    <input
      id="menuPic"
      type="file"
      accept="image/*"
      style={{ display: "none" }}
      onChange={(e) => e.target.files && setMenuPic(e.target.files[0])}
    />
  </div>

  <div style={{ marginTop: "15px" }}>
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
</div>

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
  const [showGroupPopup, setShowGroupPopup] = useState(false);
  const [showOptionPopup, setShowOptionPopup] = useState(false);
  const [groupID, setGroupID] = useState<string | null>(null);
  const [restaurantID, setRestaurantID] = useState<string | null>(null);
  const [types, setTypes] = useState<any[]>([]);
  const [selectedTypeID, setSelectedTypeID] = useState<string | null>(null);





  const token = localStorage.getItem("token");




  useEffect(() => {
  if (!restaurantID || !token) return;

  fetch(`http://localhost:8080/restaurant/menu/${restaurantID}/types`, {
    headers: { Authorization: `Bearer ${token}` },
  })
    .then((res) => res.json())
    .then((json) => {
      console.log("📄 Available types:", json);
      setTypes(json.types || []);
    })
    .catch((err) => console.error("❌ Fetch types error:", err));
}, [restaurantID]);






  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) return;

    try {
      const payload = token.split('.')[1];
      const base64 = payload.replace(/-/g, '+').replace(/_/g, '/');
      const jsonPayload = JSON.parse(atob(base64));

      if (jsonPayload.role === "restaurant") {
        setRestaurantID(jsonPayload.user_id);
        console.log("Restaurant ID:", jsonPayload.user_id);
      } else {
        console.error("Token ไม่ใช่ร้านอาหาร");
      }
    } catch (err) {
      console.error("❌ JWT decode error:", err);
    }
  }, []);

const [groupData, setGroupData] = useState({ 
  name: "",
  required: false,
  min_select: 1,
  max_select: 1,
  allow_qty: false,
});

const [optionData, setOptionData] = useState({
  name: "",
  price_delta: 0,
  is_default: false,
  max_qty: 0,
});

const [options, setOptions] = useState<any[]>([]);

if (!menu) return <p>ไม่พบข้อมูลเมนู</p>;

              const handleGroupChange = (e: any) => {
                const { name, value, type, checked } = e.target;
                setGroupData((prev) => ({
                  ...prev,
                  [name]: type === "checkbox" ? checked : Number(value) || value, // convert number inputs
                }));
              };



              const handleOptionChange = (e: any) => {
                const { name, value, type, checked } = e.target;
                setOptionData((prev) => ({
                  ...prev,
                  [name]:
                    type === "checkbox"
                      ? checked
                      : name === "price_delta" || name === "max_qty"
                      ? Number(value)
                      : value,
                }));
              };


const handleCreateGroup = async () => {
  try {
    console.log("✅ Creating AddOn Group:", groupData);
    const res = await fetch(
      `http://localhost:8080/restaurant/menu/${restaurantID}/addon-groups`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        method: "POST",
        body: JSON.stringify(groupData),
      }
    );

    const data = await res.json();
    if (!res.ok) throw new Error(data.message || "สร้างกลุ่มไม่สำเร็จ");

    setGroupID(data.id);
    console.log("🎯 Group created:", data);

    // 🔗 ลิงก์กับ type ถ้ามีเลือกไว้
    if (selectedTypeID) {
      const linkRes = await fetch(
        `http://localhost:8080/restaurant/menu/addon-groups/${data.id}/types/${selectedTypeID}`,
        {
          method: "POST",
          headers: { Authorization: `Bearer ${token}` },
        }
      );

      const linkData = await linkRes.json();
      console.log("🔗 Group linked with type:", linkData);
    } else {
      console.warn("⚠️ ยังไม่ได้เลือก type สำหรับ group นี้");
    }

    setShowGroupPopup(false);
    setShowOptionPopup(true);

  } catch (err) {
    console.error("❌ Failed to create group:", err);
    alert("สร้างกลุ่มไม่สำเร็จ");
  }
};


          const handleAddOption = () => {
            setOptions((prev) => [
              ...prev,
              {
                ...optionData,
                price_delta: Number(optionData.price_delta),
                max_qty: Number(optionData.max_qty),
              },
            ]);
            setOptionData({
              name: "",
              price_delta: 0,
              is_default: false,
              max_qty: 0,
            });
          };




          
            const handleSubmitOptions = async () => {
              if (!groupID) return alert("ยังไม่มี group id");
              try {
                console.log("✅ Sending all options:", options);

                for (const [index, opt] of options.entries()) {
                  const payload = {
                    ...opt,
                    price_delta: Number(opt.price_delta),
                    max_qty: Number(opt.max_qty),
                  };

                  const res = await fetch(
                    `http://localhost:8080/restaurant/menu/addon-groups/${groupID}/options`,
                    {
                      method: "POST",
                      headers: {
                        "Content-Type": "application/json",
                        Authorization: `Bearer ${token}`,
                      },
                      body: JSON.stringify(payload),
                    }
                  );

                  if (!res.ok) {
                    const errData = await res.json().catch(() => ({ message: "No JSON response" }));
                    console.error(`❌ Failed to send option #${index + 1}:`, payload, errData);
                    alert(`ส่ง Option "${opt.name}" ไม่สำเร็จ`);
                    return;
                  }

                  const data = await res.json().catch(() => ({}));
                  console.log(`✅ Option #${index + 1} saved:`, data);
                }

                // reset after all options sent
                setOptions([]);
                setOptionData({
                  name: "",
                  price_delta: 0,
                  is_default: false,
                  max_qty: 0,
                });
                setGroupID(null);
                setShowOptionPopup(false);
                alert("ส่ง options สำเร็จทั้งหมด ✅");
              } catch (err) {
                console.error("❌ Failed to submit options:", err);
                alert("ส่ง options ไม่สำเร็จ");
              }
            };

  

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
            <h4 className={styles.handlerthisfkignstupidshit}>🍳 Add-ons (ตัวเลือกเพิ่มเติม) <button onClick={() => setShowGroupPopup(true)}  className={styles.addonsBTN}> <span className="material-symbols-outlined">add_circle</span>เพิ่ม Add-ons</button></h4>
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
{showGroupPopup && (
  <div className={styles.popupOverlay} onClick={() => setShowGroupPopup(false)}>
    <div className={styles.popupForm} onClick={(e) => e.stopPropagation()}>
      <h3>Add-Ons Group</h3>

      <div className={styles.inlineInputs}>
        <label>
          ลิงก์กับประเภทเมนู: <span style={{ color: "red" }}>*</span>
          <select
            required
            value={selectedTypeID || ""}
            onChange={(e) => setSelectedTypeID(e.target.value)}
          >
            <option value="">-- เลือกประเภท --</option>
            {types.map((t: any) => (
              <option key={t.id} value={t.id}>
                {t.type}
              </option>
            ))}
          </select>
        </label>
      </div>

      <label>
        ชื่อกลุ่ม: <span style={{ color: "red" }}>*</span>
        <input
          required
          name="name"
          value={groupData.name}
          onChange={handleGroupChange}
          placeholder="กรอกชื่อกลุ่ม"
        />
      </label>

      <label className={styles.checkboxRow}>
        <input
          type="checkbox"
          name="required"
          checked={groupData.required}
          onChange={handleGroupChange}
        />
        <span>
          Required (บังคับให้ลูกค้าเลือกอย่างน้อยหนึ่ง)
        </span>
      </label>

      <div className={styles.inlineInputs}>
        <label>
          Min select: <span style={{ color: "red" }}>*</span>
          <input
            required
            type="number"
            name="min_select"
            min={1}
            value={groupData.min_select}
            onChange={handleGroupChange}
          />
        </label>
        <label>
          Max select: <span style={{ color: "red" }}>*</span>
          <input
            required
            type="number"
            name="max_select"
            min={1}
            value={groupData.max_select}
            onChange={handleGroupChange}
          />
        </label>
      </div>

      <label className={styles.checkboxRow}>
        <input
          type="checkbox"
          name="allow_qty"
          checked={groupData.allow_qty}
          onChange={handleGroupChange}
        />
        <span>
          เลือกจำนวนของ Add-on ได้ไหม เช่น เพิ่มชีส 2 ชุด
        </span>
      </label>

      <div className={styles.popupActions}>
        <button className={styles.confirmBtn} onClick={handleCreateGroup}>สร้าง Group</button>
        <button className={styles.cancelBtnnnnnn} onClick={() => setShowGroupPopup(false)}>ยกเลิก</button>
      </div>
    </div>
  </div>
)}




{/* -------- Popup: Add Options -------- */}
{showOptionPopup && (
  <div className={styles.popupOverlay} onClick={() => setShowOptionPopup(false)}>
    <div className={styles.popupForm} onClick={(e) => e.stopPropagation()}>
      <h3>เพิ่ม Option ใน Group</h3>

      <label>
        ชื่อ Option: <span style={{ color: "red" }}>*</span>
        <input
          required
          name="name"
          value={optionData.name}
          onChange={handleOptionChange}
          placeholder="กรอกชื่อ Option เช่น เพิ่มชีส"
        />
      </label>

      <label>
        ราคาเพิ่ม (฿): <span style={{ color: "red" }}>*</span>
        <input
          required
          type="number"
          name="price_delta"
          min={0}
          value={optionData.price_delta}
          onChange={handleOptionChange}
          placeholder="เช่น 10"
        />
      </label>

      <label className={styles.checkboxRow}>
        <input
          type="checkbox"
          name="is_default"
          checked={optionData.is_default}
          onChange={handleOptionChange}
        />
        <span>
          เป็นค่าเริ่มต้น (เลือกให้อัตโนมัติ)
        </span>
      </label>

      <label>
        Max Quantity: <span style={{ color: "red" }}>*</span>
        <input
          required
          type="number"
          name="max_qty"
          min={1}
          value={optionData.max_qty}
          onChange={handleOptionChange}
          placeholder="เช่น 3"
        />
      </label>

      <button className={styles.addOptionBtn} onClick={handleAddOption}>เพิ่ม Option</button>

      {options.length > 0 && (
        <ul className={styles.optionList}>
          {options.map((opt, i) => (
            <li key={i}>
              {opt.name} (+฿{opt.price_delta})
            </li>
          ))}
        </ul>
      )}

      <div className={styles.popupActions}>
        <button className={styles.confirmBtn} onClick={handleSubmitOptions}>บันทึกทั้งหมด</button>
        <button className={styles.cancelBtnnnnnn} onClick={() => setShowOptionPopup(false)}>ยกเลิก</button>
      </div>
    </div>
  </div>
)}







    </div>
  );
  
}
