"use client";
import { useSession } from "next-auth/react";
import { NextAuthProvider } from "../providers";
import { TransactionTable } from "./transactionTable";
import { useState, useRef, useEffect } from "react";
import { Swiper, SwiperSlide, SwiperClass } from "swiper/react";
import { Plus } from "react-feather";
import { Balance } from "./Balance";
import { MoneyPool } from "./type";
export default function Mypage() {
  return (
    <NextAuthProvider>
      <MypageContents />
    </NextAuthProvider>
  );
}

const moneyPools: MoneyPool[] = [
  {
    id: 1,
    name: "食費",
    description: "食費",
    is_world_public: true,
    owner_id: 1,
    color: "#f9842c",
    amount: 1000,
  },
  {
    id: 2,
    name: "生活費",
    description: "生活費",
    is_world_public: true,
    owner_id: 1,
    color: "#f20027",
    amount: 1000,
  },
  {
    id: 3,
    name: "飲み",
    description: "飲み",
    is_world_public: true,
    owner_id: 1,
    color: "#1010f0",
    amount: 1000,
  },
];

function MypageContents() {
  const { data: session } = useSession();
  const [moneyPoolIndex, setMoneyPoolIndex] = useState(0);
  const [swiper, setSwiper] = useState<SwiperClass>();
  const [isAddMode, setIsAddMode] = useState(false);
  const inputEl = useRef<HTMLInputElement>(null!);

  async function addTransaction() {
    const res = await fetch("/api/httptest", {
      method: "POST",
      body: JSON.stringify({
        name: inputEl.current.value,
      }),
    });
    const data = await res.json();
    console.log(data);
  }

  useEffect(() => {
    if (inputEl.current) {
      console.log("focus");
      inputEl.current.focus();
    }
  }, [isAddMode]);

  if (session && session.user) {
    return (
      <div className="flex p-5 h-[calc(100vh-5rem)]">
        <div className="col-span-1 w-4/12">
          <Balance user={session.user} moneypools={moneyPools} />
        </div>

        <div className="col-span-1 w-8/12">
          <div className="h-10 flex items-center ml-6 gap-x-1">
            {moneyPools.map((moneyPool, index) => (
              <div
                key={moneyPool.id}
                className="flex border-0 rounded-t-2xl h-full justify-center px-2 font-bold font-Noto items-center min-w-max w-20 border-b-0 cursor-pointer"
                style={
                  moneyPoolIndex === index
                    ? { backgroundColor: moneyPool.color, color: "#ffffff" }
                    : {
                        backgroundColor: "#f4f4f4",
                      }
                }
                onClick={() => {
                  if (swiper) {
                    swiper.slideTo(index);
                  }
                }}
              >
                <div>{moneyPool.name}</div>
              </div>
            ))}
          </div>
          <div
            className={
              "border-2 rounded-3xl p-1 h-[calc(100%-2rem)] overflow-hidden w-full"
            }
            style={{ borderColor: moneyPools[moneyPoolIndex].color }}
          >
            <Swiper
              spaceBetween={1}
              slidesPerView={1}
              onSlideChange={(i) => setMoneyPoolIndex(i.activeIndex)}
              onSwiper={(swiper) => {
                const swiperInstance = swiper;
                setSwiper(swiperInstance);
              }}
              className="flex w-full h-[calc(100%-4rem)] ml-4"
            >
              {moneyPools.map((moneyPool) => (
                <SwiperSlide key={moneyPool.id} className="px-2 ">
                  <div className="border-2 border-transparent h-full">
                    <TransactionTable />
                  </div>
                </SwiperSlide>
              ))}
            </Swiper>
            <div className="flex justify-center items-center h-16 w-full">
              <div
                className="w-[95%] h-12 cursor-pointer"
                onClick={() => {
                  setIsAddMode(true);
                }}
                onBlur={(fe) => {
                  if (!fe.currentTarget.contains(fe.relatedTarget)) {
                    setIsAddMode(false);
                  }
                }}
                tabIndex={0}
              >
                {isAddMode ? (
                  <div
                    className={`flex h-12 items-center gap-2 w-full border-2 border-gray-200 hover:border-primary-default rounded-full shadow-md px-2 font-Noto`}
                    onMouseOut={(e) => {
                      e.currentTarget.style.borderColor = "transparent";
                    }}
                    onMouseOver={(e) => {
                      e.currentTarget.style.borderColor =
                        moneyPools[moneyPoolIndex].color;
                    }}
                  >
                    <button
                      className="h-5/6"
                      style={{ color: moneyPools[moneyPoolIndex].color }}
                      tabIndex={0}
                      onClick={async (e) => {
                        e.preventDefault();
                        await addTransaction();
                      }}
                    >
                      <Plus className="h-full w-full" />
                    </button>
                    <div className="w-11/12 flex gap-2 justify-start items-center p-1">
                      <input
                        type="date"
                        className="h-[80%] hover:border-0 focus:outline-none w-[15%] px-0"
                        onKeyDown={(e) => {
                          if (e.key === "Enter") {
                            e.preventDefault();
                          }
                        }}
                        placeholder="日付"
                      />
                      <input
                        className="h-[80%] hover:border-0 focus:outline-none w-[75%]"
                        ref={inputEl}
                        onKeyDown={(e) => {
                          if (e.key === "Enter") {
                            e.preventDefault();
                          }
                        }}
                        placeholder="タイトル"
                      />
                      <input
                        className="h-[80%] hover:border-0 focus:outline-none w-[10%]"
                        onKeyDown={async (e) => {
                          if (e.key === "Enter") {
                            if (e.currentTarget) {
                              e.currentTarget.blur();
                            }
                            e.preventDefault();
                            await addTransaction();
                          }
                        }}
                        placeholder="金額"
                      />
                    </div>
                  </div>
                ) : (
                  <div
                    className="flex h-12 items-center gap-2 w-full border-2 border-transparent hover:bg-gray-50 hover:border-primary-default rounded-full hover:shadow-md px-2 font-Noto"
                    onMouseOut={(e) => {
                      e.currentTarget.style.borderColor = "transparent";
                    }}
                    onMouseOver={(e) => {
                      e.currentTarget.style.borderColor =
                        moneyPools[moneyPoolIndex].color;
                    }}
                  >
                    <div className="h-5/6  border-primary-default aspect-square">
                      <div
                        className="h-full w-full"
                        style={{ color: moneyPools[moneyPoolIndex].color }}
                      >
                        <Plus className="h-full w-full" />
                      </div>
                    </div>
                    追加
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
}
