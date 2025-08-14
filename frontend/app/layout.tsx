import { Noto_Sans_Thai } from "next/font/google";
import type { ReactNode } from "react";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import "../styles/globals.css";

const notoThai = Noto_Sans_Thai({
  subsets: ["thai"],
  weight: ["400", "700"], // เลือกน้ำหนักที่ต้องการ
  variable: "--font-noto-thai",
});

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="th" className={notoThai.variable}>
      <head>
        <link
          href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined"
          rel="stylesheet"
        />
      </head>

      <body style={{ display: "flex", flexDirection: "column", minHeight: "100vh" }}>
        <Navbar />
        <main style={{ flex: 1 }}>{children}</main>
        <Footer /> {/* Footer อยู่ที่นี่เพียงแห่งเดียว */}
      </body>
    </html>
  );
}
