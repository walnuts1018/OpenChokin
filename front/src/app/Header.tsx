import Image from "next/image";
import Link from "next/link";
import { LoginButton, LogoutButton } from "./authbuttons";

import { getServerSession } from "next-auth/next";
import { authOptions } from "./api/auth/[...nextauth]/options";

export async function Header() {
  const session = await getServerSession(authOptions);
  const user = session?.user;

  return (
    <>
      <header className="fixed top-0 z-50 w-full">
        <div className="flex justify-center w-full h-20 bg-white">
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
              {user ? <LogoutButton /> : <LoginButton />}
            </div>
          </div>
        </div>
        <div className="flex justify-center w-full border-0 bg-white">
          <div className="w-11/12 h-[1px] bg-gray-300 px-20"></div>
        </div>
        <div className="flex justify-center w-full border-0 bg-transparent">
          <div className="w-11/12 h-[1px] bg-gray-300 px-20"></div>
        </div>
        {/*<div className="flex justify-center w-full border-0 bg-gradient-to-b from-white to-transparent h-1"></div>*/}
      </header>
      <div className="header-space h-20" />
    </>
  );
}
