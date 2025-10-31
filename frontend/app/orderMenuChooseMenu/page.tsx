"use client";
import { SearchParamsContext } from "next/dist/shared/lib/hooks-client-context.shared-runtime";
import styles from "./orderMenuChooseMenu.module.css";
import { useSearchParams } from "next/navigation";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

type UUID = string;

export default function MenuPage() {
  // const restaurant_id = searchParam.get("id") || "";
  const restaurant_id = "7e3073db-fae6-4520-be71-2f8088aa15fc";

  const [menuItems, setMenuItems] = useState<MenuItem[]>([]);
  const [cart, setCart] = useState<CartItem[]>([])
  const [isPopupOpen, setIsPopupOpen] = useState(false);
  const [selectedItem, setSelectedItem] = useState<MenuItem | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const storedToken = localStorage.getItem("token");

    if (!storedToken) {
      setError("ไม่พบโทเค็น กรุณาเข้าสู่ระบบใหม่");
      return;
    }

    async function fetchMenu() {
      try {
        const res = await fetch(`http://localhost:8080/restaurant/menu/${restaurant_id}/items`, {
            headers: {"Authorization": `Bearer ${storedToken}`}
          });

        if (!res.ok) {
          setError("ไม่สามารถโหลดเมนูได้");
          return;
        }

        const data = await res.json();

        const formatted: MenuItem[] = data.items.map((item: any) => ({
          id: item.id,
          name: item.name,
          price: item.price.toString(),
          menu_pic: item.menu_pic ?? undefined,
          description: item.description ?? "",
          types: item.types.map((t: any) => ({
            id: t.id,
            name: t.name ?? t.type
          })),
          addons: [],
        }));

        setMenuItems(formatted);
      } catch (err) {
        console.error(err);
        setError("เกิดข้อผิดพลาดในการโหลดเมนู");
      }
    }

    fetchMenu();
  }, []);

  const handleAddItem = (menuItem: MenuItem) => {
    setSelectedItem(menuItem);
    setIsPopupOpen(true);
  };

  const handleAddToCart = (item: MenuItem, quantity: number, selectedAddons: Addon[]) => {
    setCart((prev) => {
      if (quantity === 0) {
        // ลบออกจาก cart
        return prev.filter(ci => ci.item.id !== item.id)
      }

      const index = prev.findIndex(ci => ci.item.id === item.id)
      if (index >= 0) {
        // update
        const newCart = [...prev]
        newCart[index] = { item, quantity, selectedAddons }
        return newCart
      } else {
        return [...prev, { item, quantity, selectedAddons }]
      }
    })
  }


  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <h1>Welcome to [ชื่อร้าน]</h1>
      </header>

      <FilterGroup />

      <div className={styles.menuCon}>
        {menuItems.map((item) => (
          <MenuItem
            key={item.id}
            menuItem={item}
            onAdd={() => handleAddItem(item)}
            cart={cart}
          />
        ))}
      </div>

      {selectedItem && (
        <MenuPopup
          isOpen={isPopupOpen}
          onClose={() => setIsPopupOpen(false)}
          item={selectedItem}
          cartItem={cart.find(ci => ci.item.id === selectedItem.id) ?? null}
          onAddToCart={handleAddToCart}
        />
      )}

      {cart.length > 0 && <Cart cart={cart} />}
    </div>
  );
}

interface FilterButtonProps {
  label: string
  isActive?: boolean
  onClick?: () => void
}

function FilterButton({ label, isActive = false, onClick }: FilterButtonProps) {
  return (
    <button
      className={isActive ? styles.filterOn : styles.filterOff}
      onClick={onClick}
    >
      {label}
    </button>
  )
}

function FilterGroup() {
  const defaultFilter = "All"
  const filters = ["อาหารจานเดียว", "เมนูเส้น", "เครื่องดื่ม", "ของหวาน", "ของทานเล่น", "ชุดอาหาร", "สลัด", "ซุป", "อาหารว่าง"]
  
  const [activeIndex, setActiveIndex] = useState<number>(0)
  const [otherButtons, setOtherButtons] = useState<string[]>([])
  const [loading, setLoading] = useState<boolean>(true)

  useEffect(() => {
    
    setOtherButtons(filters)
    setLoading(false)
  }, [])

  const handleClick = (index: number) => {
    if (activeIndex === index) {
      setActiveIndex(0)
    } else {
      setActiveIndex(index)
    }
  }

  return (
    <div className={styles.filterBar}>
      <FilterButton
        label={defaultFilter}
        isActive={activeIndex === 0}
        onClick={() => handleClick(0)}
      />

      <div className={styles.filterBarScroll}>
        {loading ? (
          <span>Loading...</span>
        ) : (
          otherButtons.map((label, i) => (
            <FilterButton
              key={label}
              label={label}
              isActive={activeIndex === i + 1}
              onClick={() => handleClick(i + 1)}
            />
          ))
        )}
      </div>
    </div>
  )
}

interface Type {
  id: UUID
  name: string
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

interface MenuItem {
  id: UUID
  name: string
  price: string
  menu_pic?: string
  description?: string
  types: Type[]
  addons?: Addon[]
  onAdd?: () => void
}

function MenuItem({ menuItem, onAdd, cart }: { menuItem: MenuItem, onAdd: () => void, cart: CartItem[] } ) {
  const cartItem = cart.find(ci => ci.item.id === menuItem.id)
  const quantity = cartItem?.quantity ?? 0

  return (
    <div className={styles.menuItemCon}>
      <div className={styles.menuItemInfoCon}>
        <img src={menuItem.menu_pic ?? "/placeholder.png"} />
        <div>
          <h3>{menuItem.name}</h3>
          <p>฿ {menuItem.price}</p>
        </div>
      </div>

    <button className={styles.bt} onClick={onAdd}>
      {quantity === 0 ? (
        <img src="/Add_Plus_Circle.svg" />
      ) : (
        <span className={styles.cartQtyCircle}>{quantity}</span>
      )}
    </button>

    </div>
  );
}

interface MenuPopupProps {
  isOpen: boolean
  onClose: () => void
  item: MenuItem
  cartItem?: CartItem | null
  onAddToCart: (item: MenuItem, quantity: number, selectedAddons: Addon[]) => void
}

function MenuPopup({ isOpen, onClose, item, cartItem, onAddToCart }: MenuPopupProps) {
  const restaurant_id = "7e3073db-fae6-4520-be71-2f8088aa15fc";

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
    console.log("Added to cart:", item, quantity, selectedAddons);
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
          <h3>ตัวเลือกเพิ่มเติมสำหรับ {detail.name}</h3>
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
                      <div style={{ display: "flex", justifyContent: "space-between", width: "150px" }}>
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

interface CartItem {
  item: MenuItem
  quantity: number
  selectedAddons: Addon[]
}

interface CartProps {
  cart: CartItem[]
}

function Cart({ cart }: CartProps) {
  if (cart.length === 0) return null

  const itemNum = cart.reduce((sum, ci) => sum + ci.quantity, 0)

  const itemCost = cart.reduce((sum, ci) => {
    const basePrice = parseFloat(ci.item.price)

    // ✅ รวมราคา addon (option ทั้งหมด)
    const addonTotal = ci.selectedAddons?.reduce((addonSum, addon) => {
      const optionTotal = addon.options.reduce((optSum, opt) => {
        return optSum + parseFloat(opt.price_delta)
      }, 0)
      return addonSum + optionTotal
    }, 0) ?? 0

    const finalItemPrice = (basePrice + addonTotal) * ci.quantity
    return sum + finalItemPrice
  }, 0)

  return (
    <div className={styles.cartCon}>
      <div className={styles.cartItemCon}>
        <img src="/Shopping_Cart_Black.svg" />
        <span>{itemNum}</span>
      </div>
      <span>ตะกร้าของฉัน</span>
      <span>฿ {itemCost}</span>
    </div>
  )
}