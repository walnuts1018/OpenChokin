"use client";
import ZitadelForm from "./Zitadelform";
import { useRouter } from "next/navigation";

export default function Page({
  searchParams,
}: {
  searchParams: { [key: string]: string | string[] | undefined };
}) {
  const router = useRouter();
  const callbackUrl = searchParams.callbackUrl;

  return (
    <div className="page h-[calc(100vh-5rem)] text-black bg-gray-100 flex-col flex justify-center items-center">
      <div className="flex justify-between flex-col items-center h-2/3 bg-white w-1/2 rounded-2xl border-primary-default shadow-xl px-4 min-w-max py-20">
        <div className="pt-10 items-center h-full">
          <h1 className="text-black font-semibold font-Noto text-2xl">
            サインイン方法を選択
          </h1>
        </div>
        <div className="items-center">
          <ZitadelForm callbackUrl={callbackUrl} />
          <div className="flex justify-center mt-8">
            <button
              type="button"
              onClick={() => router.back()}
              className=" text-xl items-center bg-white hover:bg-gray-100 rounded-full font-bold text-primary-default px-8 py-1 w-2/12 border-2 border-primary-default min-w-fit"
            >
              戻る
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
