"use client";
import { useSession } from "next-auth/react";
import { NextAuthProvider } from "../providers";
import { TransactionTable } from "./transactionTable";
import { useState, useRef, useEffect, ReactElement } from "react";
import { Swiper, SwiperSlide, SwiperClass } from "swiper/react";
import { Balance } from "./Balance";
import { MoneyPool, MoneyProviderSum } from "./type";
import { ThemeProvider, createTheme } from "@mui/material/styles";
import Image from "next/image";
import { AddButton } from "./AddButton";
import { MoneyPoolSum } from "./type";

export default function Mypage() {
  return (
    <NextAuthProvider>
      <MypageContents />
    </NextAuthProvider>
  );
}

const moneyPoolColors = [
  "#00BFFF",
  "#FFA500",
  "#FF69B4",
  "#ADFF2F",
  "#D8BFD8",
  "#CD853F",
];

const theme = createTheme({
  palette: {
    primary: {
      main: "#f9842c",
      dark: "#FA6C28",
    },
  },
  typography: {
    fontFamily: "var(--font-Noto)",
  },
});

function MypageContents() {
  const { data: session } = useSession();
  const [moneyPoolIndex, setMoneyPoolIndex] = useState(0);
  const [swiper, setSwiper] = useState<SwiperClass>();
  const [moneyPoolSums, setMoneyPoolSums] = useState<MoneyPoolSum[]>([
    {
      id: "1",
      name: "食費",
      Sum: 1000,
      Type: "public",
      emoji: "🍣",
    },
  ]);
  const [moneyProviders, setMoneyProviders] = useState<MoneyProviderSum[]>([
    {
      id: "1",
      name: "PayPay",
      balance: 1000,
    },
  ]);
  useEffect(() => {
    const getMoneyPools = async () => {
      if (session && session?.user) {
        const res = await fetch(
          `/api/back/moneypools?type=summary&user_id=${session.user.sub}`,
          {
            method: "GET",
            headers: {
              Authorization: `Bearer ${session.user.idToken}`,
            },
          }
        );
        if (res.ok) {
          const mps: MoneyPoolSum[] = await res.json();
          setMoneyPoolSums(mps);
        }
      }
    };

    const getMoneyProviders = async () => {
      if (session && session?.user) {
        const res = await fetch(
          `/api/back/moneyproviders?type=summary&user_id=${session.user.sub}`,
          {
            method: "GET",
            headers: {
              Authorization: `Bearer ${session.user.idToken}`,
            },
          }
        );
        if (res.ok) {
          const mps: MoneyProviderSum[] = await res.json();
          setMoneyProviders(mps);
        }
      }
    };

    getMoneyPools();
    getMoneyProviders();
  }, [session]);

  if (session && session.user) {
    return (
      <ThemeProvider theme={theme}>
        <div className="flex p-5 h-[calc(100vh-5rem)]">
          <div className="col-span-1 w-4/12 px-2">
            <div className="h-1/6 flex justify-center items-center pb-4 gap-2">
              <div className="flex gap-6 h-7/12 w-3/4 p-2 font-medium text-5xl font-Noto justify-between items-center border-b-4 border-cyan-600">
                <div className="font-light text-3xl t-0 l-0">総残高</div>
                <div className="flex gap-4 w-fit r-0">
                  <p>
                    {moneyPoolSums
                      .reduce(function (sum, moneypool) {
                        return sum + moneypool.Sum;
                      }, 0)
                      .toLocaleString(undefined, {
                        maximumFractionDigits: 5,
                      })}
                  </p>
                  <p>円</p>
                </div>
              </div>
              <div className="h-4/12 flex items-center justify-center pt-4 pl-3 b-0">
                <a
                  href={`https://twitter.com/intent/tweet?text=%E7%A7%81%E3%81%AE%E6%AE%8B%E9%AB%98%E3%81%AF${moneyPoolSums
                    .reduce(function (sum, moneypool) {
                      return sum + moneypool.Sum;
                    }, 0)
                    .toLocaleString(undefined, {
                      maximumFractionDigits: 5,
                    })}%E5%86%86%E3%81%A7%E3%81%99%EF%BC%81%0D%0AOpenChokin%E3%81%A7%E5%AE%B6%E8%A8%88%E7%B0%BF%E3%82%92%E5%85%A8%E4%B8%96%E7%95%8C%E3%81%AB%E5%85%AC%E9%96%8B%EF%BC%81&url=https://openchokin.walnuts.dev&hashtags=OpenChokina&via=walnuts1018`}
                  rel="nofollow"
                  target="_blank"
                  className="h-full font-Nunito text-xl"
                >
                  <Image
                    src="/icons/twitter-x-line.svg"
                    alt="x"
                    width={30}
                    height={30}
                    style={{ objectFit: "contain" }}
                    className="min-w-[30px] max-w-[30px]"
                  />
                </a>
              </div>
            </div>
            <div className="h-5/6">
              <Balance
                user={session.user}
                moneypoolSums={moneyPoolSums}
                moneyProviders={moneyProviders}
              />
            </div>
          </div>

          <div className="col-span-1 w-8/12">
            <div className="h-10 flex items-center ml-6 gap-x-1">
              {moneyPoolSums.map((moneyPool, index) => (
                <div
                  key={moneyPool.id}
                  className="flex border-0 rounded-t-2xl h-full justify-center px-2 font-bold font-Noto items-center min-w-max w-20 border-b-0 cursor-pointer"
                  style={
                    moneyPoolIndex === index
                      ? {
                          backgroundColor:
                            moneyPoolColors[
                              moneyPoolIndex % moneyPoolColors.length
                            ],
                          color: "#ffffff",
                        }
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
              style={{
                borderColor:
                  moneyPoolColors[moneyPoolIndex % moneyPoolColors.length],
              }}
            >
              <Swiper
                spaceBetween={1}
                slidesPerView={1}
                onSlideChange={(i) => setMoneyPoolIndex(i.activeIndex)}
                onSwiper={(swiper) => {
                  const swiperInstance = swiper;
                  setSwiper(swiperInstance);
                }}
                initialSlide={0}
                className="flex w-full h-[calc(100%-4rem)]"
              >
                {moneyPoolSums.map((moneyPool, index) => (
                  <SwiperSlide key={moneyPool.id} className="">
                    <div className="border-2 border-transparent h-full mx-2">
                      <TransactionTable
                        moneyPoolID={moneyPool.id}
                        scroll={index === moneyPoolIndex}
                      />
                    </div>
                  </SwiperSlide>
                ))}
              </Swiper>
              <div className="flex justify-center items-center h-16 w-full">
                <AddButton
                  color={
                    moneyPoolColors[moneyPoolIndex % moneyPoolColors.length]
                  }
                  moneyPoolID={moneyPoolSums[moneyPoolIndex].id}
                />
              </div>
            </div>
          </div>
        </div>
      </ThemeProvider>
    );
  }
}
