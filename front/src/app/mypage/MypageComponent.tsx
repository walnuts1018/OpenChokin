"use client";
import { useSession } from "next-auth/react";
import { NextAuthProvider } from "../providers";
import { TransactionTable } from "./transactionTable";
import { useState, useRef, useEffect, ReactElement } from "react";
import { Swiper, SwiperSlide, SwiperClass } from "swiper/react";
import { Balance } from "./Balance";
import { MoneyProviderSum } from "./type";
import { ThemeProvider, createTheme } from "@mui/material/styles";
import Image from "next/image";
import { AddButton } from "./AddButton";
import { MoneyPoolSum } from "./type";
import { createContext } from "react";
import { useCallback } from "react";

export const TransactionContext = createContext({});

export function MypageComponent({ userID }: { userID?: string }) {
  return (
    <NextAuthProvider>
      <MypageContents userID={userID} />
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

function MypageContents({ userID }: { userID?: string }) {
  const { data: session } = useSession();
  const [moneyPoolIndex, setMoneyPoolIndex] = useState(0);
  const [swiper, setSwiper] = useState<SwiperClass>();
  const [moneyPoolSums, setMoneyPoolSums] = useState<MoneyPoolSum[]>([]);
  const [moneyProviders, setMoneyProviders] = useState<MoneyProviderSum[]>([]);
  const [tableReloadContext, setTableReloadContext] = useState({});

  if (userID === undefined) {
    userID = session?.user.sub;
  }

  const getMoneyPools = useCallback(async () => {
    console.log("start getMoneyPools by", userID);
    if (userID) {
      const authHeader =
        userID === session?.user.sub ? `Bearer ${session.user.idToken}` : "";
      const res = await fetch(
        `/api/back/moneypools?type=summary&user_id=${userID}`,
        {
          method: "GET",
          headers: {
            Authorization: authHeader,
          },
        }
      );
      console.log("getMoneyPools res", res);
      if (res.ok) {
        const json = await res.json();
        console.log("getMoneyPools json raw", json);
        if (json.pools !== null && json.pools !== undefined) {
          const mps: MoneyPoolSum[] = json.pools as MoneyPoolSum[];
          setMoneyPoolSums(mps);
        }
      } else {
        console.log("getMoneyPools error", res);
      }
    }
  }, [userID, session]);

  const getMoneyProviders = useCallback(async () => {
    if (session && session?.user && userID) {
      const authHeader =
        userID === session?.user.sub ? `Bearer ${session.user.idToken}` : "";
      const res = await fetch(`/api/back/moneyproviders?type=summary`, {
        method: "GET",
        headers: {
          Authorization: authHeader,
        },
      });
      if (res.ok) {
        const json = await res.json();
        console.log(json);
        if (json.provider !== null && json.provider !== undefined) {
          const provider = json.provider;
          const mps: MoneyProviderSum[] = provider.map(
            (p: { id: string; name: string; balance: Number }) => {
              return {
                id: p.id,
                name: p.name,
                balance: p.balance.toLocaleString(undefined, {
                  maximumFractionDigits: 5,
                }),
              };
            }
          );

          setMoneyProviders(mps);
        }
      }
    }
  }, [session, userID]);

  useEffect(() => {
    getMoneyPools();
    getMoneyProviders();
  }, [getMoneyPools, getMoneyProviders]);

  useEffect(() => {
    const interval = setInterval(() => {
      getMoneyPools();
      getMoneyProviders();
    }, 1000 * 60);

    return () => clearInterval(interval);
  }, [getMoneyPools, getMoneyProviders]);

  if (userID !== undefined) {
    console.log(typeof moneyPoolSums);
    console.log(moneyPoolSums);
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
                        return sum + moneypool.sum;
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
                  href={`https://twitter.com/intent/tweet?text=${
                    userID === session?.user.sub ? "%E7%A7%81" : userID
                  }%E3%81%AE%E6%AE%8B%E9%AB%98%E3%81%AF${moneyPoolSums
                    .reduce(function (sum, moneypool) {
                      return sum + moneypool.sum;
                    }, 0)
                    .toLocaleString(undefined, {
                      maximumFractionDigits: 5,
                    })}%E5%86%86%E3%81%A7%E3%81%99%EF%BC%81%0D%0AOpenChokin%E3%81%A7%E5%AE%B6%E8%A8%88%E7%B0%BF%E3%82%92%E5%85%A8%E4%B8%96%E7%95%8C%E3%81%AB%E5%85%AC%E9%96%8B%EF%BC%81&url=https://openchokin.walnuts.dev/${userID}&hashtags=OpenChokin&via=walnuts1018`}
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
                userID={userID}
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
                      {userID !== undefined ? (
                        <TransactionTable
                          moneyPoolID={moneyPool.id}
                          scroll={index === moneyPoolIndex}
                          userID={userID}
                          reloadContext={tableReloadContext}
                        />
                      ) : (
                        <div></div>
                      )}
                    </div>
                  </SwiperSlide>
                ))}
              </Swiper>
              <div className="flex justify-center items-center h-16 w-full">
                {moneyPoolSums[moneyPoolIndex] !== undefined ? (
                  <AddButton
                    color={
                      moneyPoolColors[moneyPoolIndex % moneyPoolColors.length]
                    }
                    moneyPoolID={moneyPoolSums[moneyPoolIndex].id}
                    setReloadContext={setTableReloadContext}
                  />
                ) : (
                  <div></div>
                )}
              </div>
            </div>
          </div>
        </div>
      </ThemeProvider>
    );
  }
}
