"use client";
import { signIn, signOut } from "next-auth/react";

export const LoginButton = () => {
  return (
    <button
      className="bg-[#f9842c] hover:bg-[#FA6C28] rounded-full  text-white px-4 py-1 border-[#f9842c] border-2 hover:border-[#FA6C28] font-Noto font-semibold text-xl"
      style={{ marginRight: 10 }}
      onClick={() => signIn()}
    >
      ログイン
    </button>
  );
};

export const LogoutButton = () => {
  return (
    <button
      className="bg-white hover:bg-gray-100 rounded-full  text-[#f9842c] px-4 py-1 border-[#f9842c] border-2 font-Noto font-semibold text-xl"
      onClick={() => signOut()}
    >
      ログアウト
    </button>
  );
};
