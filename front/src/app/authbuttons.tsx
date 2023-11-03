"use client";
import { signIn, signOut } from "next-auth/react";
import { usePathname } from "next/navigation";
import Link from "next/link";
import { useState } from "react";
import Modal from "react-modal";

export const LoginButton = () => {
  const pathname = usePathname();
  const isSigninPage = pathname === "/signin";

  return (
    <>
      {isSigninPage ? (
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
  const pathname = usePathname();
  const isMypage = pathname === "/mypage";
  const [LogoutCheckIsOpen, setLogoutCheckIsOpen] = useState(false);
  return (
    <>
      {!isMypage ? (
        <Link
          className="bg-white hover:bg-gray-100 rounded-full  text-primary-default px-4 py-1 border-primary-default border-2 font-Noto font-semibold text-xl"
          href="/mypage"
        >
          マイページ
        </Link>
      ) : (
        <>
          <button
            className="bg-white hover:bg-gray-100 rounded-full  text-primary-default px-4 py-1 border-primary-default border-2 font-Noto font-semibold text-xl"
            onClick={() => setLogoutCheckIsOpen(true)}
          >
            ログアウト
          </button>
          <Modal
            isOpen={LogoutCheckIsOpen}
            className="flex justify-center items-center t-0 l-0 w-full h-full"
          >
            <div
              className="bg-transparent w-full h-full absolute z-10"
              onClick={() => setLogoutCheckIsOpen(false)}
            />
            <div className="w-1/2 h-1/3 bg-gray-50  transform bg-opacity-90 shadow-2xl rounded-3xl flex flex-col justify-center items-center font-Noto font-semibold text-xl gap-y-20 z-50 border-primary-default border-2">
              <div>ログアウトしますか？</div>
              <div className="flex justify-between gap-x-8">
                <button
                  className="bg-primary-default hover:bg-primary-dark rounded-full  text-white px-4 py-1 border-primary-default border-2 hover:border-primary-dark font-Noto font-semibold text-xl"
                  onClick={() => signOut()}
                >
                  ログアウト
                </button>
                <button
                  className="bg-white hover:bg-gray-100 rounded-full  text-primary-default px-4 py-1 border-primary-default border-2 font-Noto font-semibold text-xl"
                  onClick={() => setLogoutCheckIsOpen(false)}
                >
                  キャンセル
                </button>
              </div>
            </div>
          </Modal>
        </>
      )}
    </>
  );
};
