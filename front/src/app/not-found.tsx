"use client";
import { useRouter } from "next/navigation";
import Link from "next/link";
export default function NotFoundPage() {
  const router = useRouter();
  return (
    <div className="h-screen">
      <div className="h-20 bg-white"></div>
      <div className="flex  items-center justify-center font-Noto mt-40">
        <div className=" flex flex-col items-center justify-center space-y-2 border-2 rounded-2xl w-6/12 shadow-md h-60 py-6 min-w-max px-10">
          <div className=" text-2xl flex flex-col items-center justify-center h-full">
            <h1>このページは存在しません。</h1>
          </div>

          <div className="mt-0 space-x-8 w-full flex items-center justify-center h-fit">
            <Link
              href="/"
              className=" text-xl items-center bg-[#f9842c] hover:bg-[#FA6C28] rounded-full font-bold text-white px-8 py-1 min-w-fit w-2/12 border-2 border-[#f9842c] hover:border-[#FA6C28]"
            >
              トップ
            </Link>
            <button
              type="button"
              onClick={() => router.back()}
              className=" text-xl items-center bg-white hover:bg-gray-100 rounded-full font-bold text-[#f9842c] px-8 py-1 w-2/12 border-2 border-[#f9842c] min-w-fit"
            >
              前の画面
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
