"use client";
import styles from "./orderMenuChooseMenu.module.css";
import { useSearchParams } from "next/navigation";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

type UUID = string;

export default function MenuPage() {
  // const restaurant_id = useSearchParams().get("restaurantId")
  // const restaurant_id = "efec0d8e-ede6-48f0-a923-6054475b816f"
  const [isPopupOpen, setIsPopupOpen] = useState(false)
  const [selectedItem, setSelectedItem] = useState("")

  const handleAddItem = (itemName: string) => {
    setSelectedItem(itemName)
    setIsPopupOpen(true)
  }

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <h1>Welcome to [ชื่อร้าน]</h1>
      </header>
      <FilterGroup/>
      
      <Menu onAdd={handleAddItem} />
      
      <MenuPopup isOpen={isPopupOpen} onClose={() => setIsPopupOpen(false)} itemName={selectedItem} />
      <Cart/>
    </div>
  )
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

function Menu({ onAdd }: { onAdd: (itemName: string) => void }) {
  return (
    <div className={styles.menuCon}>
      <MenuItem/>
    </div>

  )
}

// MenuItem Component
interface MenuItemProps {
  id: UUID
  name: string
  price: string
  menu_pic?: string
  description: string
  types: string[]
  onAdd?: () => void
}

function MenuItem({ imageSrc, name, price, onAdd }: MenuItemProps) {
  return (
    <div className={styles.menuItemCon}>
      <div className={styles.menuItemInfoCon}>
        <img src={imageSrc ? imageSrc : "/placeholder.png"} />
        <div>
          <h3>{name}</h3>
          <p>{price}</p>
        </div>
      </div>
      <button className={styles.bt} onClick={onAdd}>
        <img src="/Add_Plus_Circle.svg" />
      </button>
    </div>
  )
}

// MenuPopup component for item customization
interface MenuPopupProps {
  isOpen: boolean
  onClose: () => void
  itemName: string
}

function MenuPopup({ isOpen, onClose, itemName }: MenuPopupProps) {
  const [quantity, setQuantity] = useState(2)
  const [selectedOptions, setSelectedOptions] = useState<{ [key: string]: boolean }>({
    option1: false,
    option2: true,
  })

  if (!isOpen) return null

  const handleQuantityChange = (delta: number) => {
    setQuantity((prev) => Math.max(1, prev + delta))
  }

  const toggleOption = (optionKey: string) => {
    setSelectedOptions((prev) => ({
      ...prev,
      [optionKey]: !prev[optionKey],
    }))
  }

  const handleAddToCart = () => {
    console.log("Added to cart:", { itemName, quantity, options: selectedOptions })
    onClose()
  }

  return (
    <div className={styles.popupOverlay}>
      <div className={styles.popupCon}>
        <button className={styles.closeBt} onClick={onClose}>
          <img src="/Close_LG.svg"/>
        </button>

        <img src="/placeholder.png" className={styles.menuImg}/>

        <div className={styles.popupOptionsCon}>
          <h3>ตัวเลือกเพิ่มเติมสำหรับ {itemName}</h3>

          <div className={styles.optionsList}>
            <label>
              <input type="checkbox" checked={selectedOptions.option1} onChange={() => toggleOption("option1")} />
              <span style={{ whiteSpace: "pre" }}>
                xxxxxxxx{"\t"} 10 ฿
              </span>
            </label>

            <label>
              <input type="checkbox" checked={selectedOptions.option2} onChange={() => toggleOption("option2")} />
              <span style={{ whiteSpace: "pre" }}>
                xxxxxxxx{"\t"} 10 ฿
              </span>
            </label>
          </div>
        </div>

        <div className={styles.popupBottomCon}>
          <div className={styles.quantityCon}>
            <button className={styles.bt} onClick={() => handleQuantityChange(-1)}>
              <img src="/Remove_Minus.svg"/>
            </button>
            <span>{quantity}</span>
            <button className={styles.bt} onClick={() => handleQuantityChange(1)}>
              <img src="/Add_Plus.svg"/>
            </button>
          </div>

          <button className={styles.addCartBt} onClick={handleAddToCart}>
            <img src="/Shopping_Cart_White.svg"/>
            เพิ่มลงตะกร้า
          </button>
        </div>
      </div>
    </div>
  )
}

function Cart() {
  

  return (
    <div className={styles.cartCon}>
      <div className={styles.cartItemCon}>
        <img src="/Shopping_Cart_Black.svg" />
        <span>itemNum</span>
      </div>
      <span>ตะกร้าของฉัน</span>
      <span>฿ itemCost</span>
    </div>
  )
}