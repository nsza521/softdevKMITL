import styles from "./restaurant.module.css";
import { Noto_Sans_Thai } from "next/font/google";

const notoThai = Noto_Sans_Thai({
  subsets: ["thai"],
  weight: ["400", "700"],
  variable: "--font-noto-thai",
});

export default function LoginPage() {
  return (
    <div className={`${styles.container} ${notoThai.variable}`}>
        <section className={styles.sidebar}>
            <div className={styles.sidebarsection}>
                <div><h2>[ชื่อร้านจ้า]</h2></div>
            </div>
            <div className={styles.sidebarsection}>
                <button><span className="material-symbols-outlined">shopping_cart</span><span>Order Menu</span></button>
            </div>
            <div className={styles.sidebarsection}>
                <button><span className="material-symbols-outlined">star</span><span>Queue</span></button>
            </div>
            <div className={styles.sidebarsection}>
                <button><span className="material-symbols-outlined">document_search</span><span>Total Sales</span></button>
            </div>
            <div className={styles.sidebarsection}>
                <button><span className="material-symbols-outlined">edit</span><span>Manage</span></button>
            </div>
            <div className={styles.sidebarsection} id={styles.logoutbtn}>
                <button><span className="material-symbols-outlined">logout</span><span>Logout</span></button>
            </div>
        </section>

        <section className={styles.shopcontent}>
            <div className={styles.shophead}>
                <div className={styles.restaurantname}>
                    <div><h2>Welcome To [ชื่อร้านจ้า]</h2> <button><span className="material-symbols-outlined">edit</span></button></div>
                    <div></div>
                </div>
                <section className={styles.category}>
                    <section className={styles.all}>
                        <button>All</button>
                    </section>
                    <section className={styles.cate}>
                        <button>อาหารตามสั่ง</button>
                        <button>เมนูเส้น</button>
                        <button>เมนูข้าว</button>
                        <button>ขาหมู</button>
                        <button>ก๋วยเตี๋ยว</button>

                    </section>
                </section>
            </div>
            <div className={styles.s_content_detail}>
                <div className={styles.menu}> 
                    <div className={styles.menuimg}></div>
                    <div className={styles.menudetail}>
                        <p>฿ xx.xx</p>
                        <p>xxxxxxxxxxxxx</p>
                    </div>
                </div>
                <div className={styles.menu}> 
                    <div className={styles.menuimg}></div>
                    <div className={styles.menudetail}>
                        <p>฿ xx.xx</p>
                        <p>xxxxxxxxxxxxx</p>
                    </div>
                </div>
                <div className={styles.menu}> 
                    <div className={styles.menuimg}></div>
                    <div className={styles.menudetail}>
                        <p>฿ xx.xx</p>
                        <p>xxxxxxxxxxxxx</p>
                    </div>
                </div>
                <div className={styles.menu}> 
                    <div className={styles.menuimg}></div>
                    <div className={styles.menudetail}>
                        <p>฿ xx.xx</p>
                        <p>xxxxxxxxxxxxx</p>
                    </div>
                </div>
                <div className={styles.menu}> 
                    <div className={styles.menuimg}></div>
                    <div className={styles.menudetail}>
                        <p>฿ xx.xx</p>
                        <p>xxxxxxxxxxxxx</p>
                    </div>
                </div>
                <div className={styles.menu}> 
                    <div className={styles.menuimg}></div>
                    <div className={styles.menudetail}>
                        <p>฿ xx.xx</p>
                        <p>xxxxxxxxxxxxx</p>
                    </div>
                </div>
                <div className={styles.menu}> 
                    <div className={styles.menuimg}></div>
                    <div className={styles.menudetail}>
                        <p>฿ xx.xx</p>
                        <p>xxxxxxxxxxxxx</p>
                    </div>
                </div>
                <div className={styles.menu}> 
                    <div className={styles.menuimg}></div>
                    <div className={styles.menudetail}>
                        <p>฿ xx.xx</p>
                        <p>xxxxxxxxxxxxx</p>
                    </div>
                </div>
                <div className={styles.menu}> 
                    <div className={styles.menuimg}></div>
                    <div className={styles.menudetail}>
                        <p>฿ xx.xx</p>
                        <p>xxxxxxxxxxxxx</p>
                    </div>
                </div>
                <div className={styles.menu}> 
                    <div className={styles.menuimg}></div>
                    <div className={styles.menudetail}>
                        <p>฿ xx.xx</p>
                        <p>xxxxxxxxxxxxx</p>
                    </div>
                </div>
                <div className={styles.menu}> 
                    <div className={styles.menuimg}></div>
                    <div className={styles.menudetail}>
                        <p>฿ xx.xx</p>
                        <p>xxxxxxxxxxxxx</p>
                    </div>
                </div>
                <div className={styles.menu}> 
                    <div className={styles.menuimg}></div>
                    <div className={styles.menudetail}>
                        <p>฿ xx.xx</p>
                        <p>xxxxxxxxxxxxx</p>
                    </div>
                </div>
                <div className={styles.menu}> 
                    <div className={styles.menuimg}></div>
                    <div className={styles.menudetail}>
                        <p>฿ xx.xx</p>
                        <p>xxxxxxxxxxxxx</p>
                    </div>
                </div>
                <div className={styles.menu}> 
                    <div className={styles.menuimg}></div>
                    <div className={styles.menudetail}>
                        <p>฿ xx.xx</p>
                        <p>xxxxxxxxxxxxx</p>
                    </div>
                </div>
                
            </div>
        </section>
        {/* <div className={styles.sidebarright}> */}

        {/* </div> */}
      <button className={styles.floatingBtn}><span className="material-symbols-outlined">add_2</span></button>
    </div>
  );
}
