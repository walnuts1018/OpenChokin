"use client";
import { signIn, signOut } from "next-auth/react";
import { usePathname } from "next/navigation";

export const LoginButton = () => {
  const pathname = usePathname();
  const isSignin = pathname === "/signin";

  return (
    <>
      {isSignin ? (
        <div />
      ) : (
        <button
          className="bg-primary-default hover:bg-primary-dark rounded-full  text-white px-4 py-1 border-primary-default border-2 hover:border-primary-dark font-Noto font-semibold text-xl"
          style={{ marginRight: 10 }}
          onClick={() => signIn()}
        >
          ログイン / 新規登録
        </button>
      )}
    </>
  );
};

export const LogoutButton = () => {
  return (
    <button
      className="bg-white hover:bg-gray-100 rounded-full  text-primary-default px-4 py-1 border-primary-default border-2 font-Noto font-semibold text-xl"
      onClick={() => signOut()}
    >
      ログアウト
    </button>
  );
};
