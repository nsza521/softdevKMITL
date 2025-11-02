"use client";
import { Noto_Sans_Thai } from "next/font/google";
import type { ReactNode } from "react";
import { usePathname } from "next/navigation";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import "../styles/globals.css";

const notoThai = Noto_Sans_Thai({
  subsets: ["thai"],
  weight: ["400", "700"],
  variable: "--font-noto-thai",
});

const titleGroups: Record<string, string[]> = {
  "จองโต๊ะ": ["/reserve", "/reserveSelectTime", "/reserveSelectTable", "/reserveFillUsr"],
  "หน้าหลัก": ["/home",],
  "ประวัติ" : ["/history"],
  "เติมเงิน": ["/topup"],
  "รายการแจ้งเตือน" : ["/noti"],
  "สั่งอาหาร": ["/orderMenuChooseRes", "/orderMenuChooseMenu"],
  "สรุปการสั่งซื้อ": ["/orderMenuSummary", "/waitOthers"]
};

function getTitle(pathname: string): string {
  for (const [title, paths] of Object.entries(titleGroups)) {
    if (paths.includes(pathname)) return title;
  }
  return "ลืมแมปหน้า";
}

export default function RootLayout({ children, metadata }: any) {
  const pathname = usePathname();

  // กำหนด path ที่ไม่ต้องการ Navbar + Footer
  const hiddenLayoutRoutes = ["/login", "/signup", "/restaurant" , "/loginrestaurant"];
  const isHiddenLayout = hiddenLayoutRoutes.includes(pathname);
  const fullScreenLayoutRoutes = ["/login", "/signup", "/restaurant","/loginrestaurant"];
  const isFullScreen = fullScreenLayoutRoutes.includes(pathname);

  return (
    <html lang="th" className={notoThai.className}>
      <head>
        <link
          href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined"
          rel="stylesheet"
        />
      </head>

      <body
        style={{
          display: "flex",
          flexDirection: "column",
          minHeight: "100vh",
        }}

        className={isFullScreen ? "fullscreen-layout" : "mobile-layout"}
      >
        {isHiddenLayout ? (
          children
        ) : (
          <>
            <Navbar title={getTitle(pathname)} />
            <main style={{ flex: 1 }}>{children}</main>
            <Footer />
          </>
        )}
      </body>
    </html>
  );
}
