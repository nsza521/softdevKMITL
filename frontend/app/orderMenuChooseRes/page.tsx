"use client";

import { stat } from "fs";
import styles from "./orderMenuChooseRes.module.css"
import { useEffect, useState } from "react";
import Link from "next/link";


const mockUpUser = [
    {
        id: 1,
        Restaurant: "ชื่อร้านค้า",
        Foodtype: ["foodtype1","foodtype2","foodtype3"],
        Status: "Open"
    },
    {
        id: 2,  
        Restaurant: "ชื่อร้านค้า",
        Foodtype: ["foodtype1","foodtype2","foodtype3"],
        Status: "Closed"
    },
];




export default function orderMenuChooseRes(){
    const [status, setStatus] = useState("");

    return(
        <div className={styles.container}>
                {mockUpUser.map((Rstr) => (
            <Link className={styles.blog_item} key={Rstr.id} 
            href={
                {
                    pathname : '/orderMenuChooseMenu/[id]',
                    query :{id: Rstr.id},
                }} as={`/orderMenuChooseMenu/${encodeURIComponent(Rstr.id)}`}> 
                <div className={styles.image}>
                    <img src="./Rectangle.svg" alt="ResPicture" />
                </div>              
                <div className={styles.content}>
                    <h3>ชื่อร้านค้า</h3>
                    <ul>
                        <li>foodtype,</li>
                        <li>foodtype,</li>
                        <li>foodtype,</li>
                    </ul>
                    {Rstr.Status === "Open" &&(
                        <p className={styles.status}>Open</p>
                    )}
                    {Rstr.Status === "Closed" &&(
                        <p className={styles.status}>Closed</p>
                    )}
                </div>
            </Link>
                ))}
        </div>
    )
}