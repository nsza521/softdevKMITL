"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

export default function Home() {
  const router = useRouter();

  useEffect(() => {
    router.replace("/login"); // root → login
  }, [router]);

  return null; // ไม่ต้อง render อะไร
}
