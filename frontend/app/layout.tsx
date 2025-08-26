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

export default function RootLayout({ children }: { children: ReactNode }) {
  const pathname = usePathname();
  const isAuthPage = pathname === "/login"; // เช็คว่าหน้านี้เป็น login มั้ย

  return (
    <html lang="th" className={notoThai.variable}>
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
      >
        {isAuthPage ? (
          children
        ) : (
          <>
            {/* <Navbar /> */}
            <main style={{ flex: 1 }}>{children}</main>
            {/* <Footer /> */}
          </>
        )}
      </body>
    </html>
  );
}
