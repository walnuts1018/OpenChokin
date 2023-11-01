"use client";
import { useSession } from "next-auth/react";
import { NextAuthProvider } from "../providers";

export default function Mypage() {
  return (
    <NextAuthProvider>
      <MypageContents />
    </NextAuthProvider>
  );
}

function MypageContents() {
  const { data: session } = useSession();
  if (session && session.user) {
    return (
      <div>
        Signed in as {session.user.email} <br />
      </div>
    );
  }
}
