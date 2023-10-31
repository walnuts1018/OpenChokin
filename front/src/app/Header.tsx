import Image from "next/image";
import Link from "next/link";
import { signIn, signOut } from "next-auth/react";

export function Header() {
  return (
    <>
      <header className="fixed top-0 z-50 w-full ">
        <div className="flex justify-center w-full h-20">
          <div className="w-11/12 flex items-center justify-between   text-black font-bold font-Nunito text-2xl px-1 space-x-1">
            <div className="flex items-center">
              <Image
                src="/logo.jpg"
                alt="logo"
                width={40}
                height={40}
                style={{ objectFit: "contain" }}
                className="min-w-[40px] max-w-[40px] mx-1"
              />
              <Link className="logo" href="/">
                OpenChokin
              </Link>
            </div>
            <div className="flex items-center space-x-2 font-Noto font-semibold text-xl">
              <Link
                className="bg-white hover:bg-gray-100 rounded-full  text-[#f9842c] px-4 py-1 border-[#f9842c] border-2"
                href="/signin"
              >
                ログイン
              </Link>
              <Link
                className="bg-[#f9842c] hover:bg-[#FA6C28] rounded-full  text-white px-4 py-1 border-[#f9842c] border-2"
                href="/signup"
              >
                新規登録
              </Link>
            </div>
          </div>
        </div>
        <div className="flex justify-center w-full bg-white border-0 ">
          <div className="w-11/12 h-[3px] bg-gray-200 px-20"></div>
        </div>
      </header>
    </>
  );
}
