"use client";
import { useRouter } from "next/navigation";

export default function NotFoundPage() {
  const router = useRouter();
  return (
    <div className="h-screen">
      <div className="h-20 bg-white"></div>
      <div className="flex  items-center justify-center font-Noto mt-40 ">
        <div className=" flex flex-col items-center justify-center space-y-10 border-2 rounded-2xl  w-6/12 shadow-md h-60 py-10">
          <h1 className=" text-3xl ">このページは存在しません。</h1>
          <div className="m-0">
            <button
              type="button"
              onClick={() => router.back()}
              className=" text-2xl items-center bg-[#f9842c] hover:bg-[#FA6C28] rounded-full font-bold text-white px-4 py-1  w-28 p-10"
            >
              戻る
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
