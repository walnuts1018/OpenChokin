import { MoneyPoolSum, MoneyPool, MoneyProviderSum } from "./type";
import { Swiper, SwiperSlide, SwiperClass } from "swiper/react";
import { useState } from "react";
import Checkbox from "@mui/material/Checkbox";
import { ThemeProvider, createTheme, styled } from "@mui/material/styles";
import Modal from "react-modal";
import { Plus } from "react-feather";
import { useRef } from "react";
import { useSession } from "next-auth/react";
import Picker from "emoji-picker-react";
import { EmojiClickData } from "emoji-picker-react";
import Image from "next/image";

const tabColors = ["#f5c33f", "#31aedd"];
const theme1 = createTheme({
  palette: {
    primary: {
      main: "#f5c33f",
      dark: "#c29a31",
    },
  },
  typography: {
    fontFamily: "var(--font-Noto)",
  },
});

const theme2 = createTheme({
  palette: {
    primary: {
      main: "#31aedd",
      dark: "#2585aa",
    },
  },
  typography: {
    fontFamily: "var(--font-Noto)",
  },
});

export function Balance({
  className,
  userID,
  moneypoolSums,
  moneyProviders,
}: {
  children?: React.ReactNode;
  className?: string;
  userID: string | undefined;
  moneypoolSums: MoneyPoolSum[];
  moneyProviders: MoneyProviderSum[];
}) {
  console.log("render moneyProviders", moneyProviders);
  const { data: session } = useSession();
  const [swiper, setSwiper] = useState<SwiperClass>();
  const [swiperIndex, setSwiperIndex] = useState(0);
  const [isAddMode, setIsAddMode] = useState(false);
  const [isAddMode2, setIsAddMode2] = useState(false);
  const inputMoneyPoolEmojiElement = useRef<HTMLInputElement>(null!);
  const inputMoneyPoolNameElement = useRef<HTMLInputElement>(null!);
  const [newMoneyPoolEmoji, setNewMoneyPoolEmoji] = useState("");
  const [newMoneyPoolName, setNewMoneyPoolName] = useState("");

  const [newMoneyProviderName, setNewMoneyProviderName] = useState("");
  const [newMoneyProviderBalance, setNewMoneyProviderBalance] = useState("");

  const [isEmojiPicking, setIsEmojiPicking] = useState(false);

  const onEmojiClick = (emoji: EmojiClickData, event: MouseEvent) => {
    setNewMoneyPoolEmoji(emoji.emoji);
    setIsEmojiPicking(false);
    if (inputMoneyPoolNameElement.current) {
      console.debug("NewMoneyPool, focus to name");
      inputMoneyPoolNameElement.current.focus();
    }
  };

  async function addMoneyPool() {
    if (session && session.user) {
      const res = await fetch(`/api/back/moneypools`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${session.user.idToken}`,
        },
        body: JSON.stringify({
          emoji: newMoneyPoolEmoji,
          name: newMoneyPoolName,
          description: "",
          type: "private",
        }),
      });
      if (res.ok) {
        const data = await res.json();
        setNewMoneyPoolEmoji("");
        setNewMoneyPoolName("");
        console.log(data);
      }
    }
  }

  async function addMoneyProvider() {
    console.log("addMoneyProvider");
    if (session && session.user) {
      const res = await fetch(`/api/back/moneyproviders`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${session.user.idToken}`,
        },
        body: JSON.stringify({
          name: newMoneyProviderName,
          balance: Number(newMoneyProviderBalance),
        }),
      });
      if (res.ok) {
        const data = await res.json();
        setNewMoneyProviderName("");
        setNewMoneyProviderBalance("");
        console.log("create money provider raw", data);
      }
    }
  }

  Modal.setAppElement("body");
  if (userID) {
    return (
      <div className={`h-full ${className}`}>
        <div className="h-10 flex items-center ml-6 gap-x-1">
          <SwiperTabs
            swiper={swiper}
            swiperIndex={swiperIndex}
            index={0}
            color={tabColors[0]}
            title="Money Pools"
          />
          {userID === session?.user.sub ? (
            <SwiperTabs
              swiper={swiper}
              swiperIndex={swiperIndex}
              index={1}
              color={tabColors[1]}
              title="Money Providers"
            />
          ) : (
            <></>
          )}
        </div>
        <div
          className={
            "border-2 rounded-3xl p-1 h-[calc(100%-2rem)] overflow-hidden w-full py-3"
          }
          style={{ borderColor: tabColors[swiperIndex] }}
        >
          <Swiper
            spaceBetween={1}
            slidesPerView={1}
            onSlideChange={(i) => setSwiperIndex(i.activeIndex)}
            onSwiper={(swiper) => {
              const swiperInstance = swiper;
              setSwiper(swiperInstance);
            }}
            className="flex w-full h-[90%]"
          >
            <SwiperSlide className="">
              <ThemeProvider theme={theme1}>
                <div className="border-2 border-transparent h-full mx-2 overflow-auto overflow-x-hidden">
                  <div
                    className="balance-header flex w-full h-4 justify-end px-0 items-center font-Noto font-bold gap-x-4"
                    style={{ color: tabColors[0] }}
                  >
                    <div>公開</div>
                    {session?.user.sub === userID ? <div>削除</div> : <></>}
                  </div>
                  {moneypoolSums.map((moneyPool, index) => (
                    <BalanceItem
                      key={moneyPool.id}
                      moneyPool={moneyPool}
                      userID={userID}
                    />
                  ))}
                </div>
              </ThemeProvider>
            </SwiperSlide>
            {userID !== undefined ? (
              <SwiperSlide className="">
                <ThemeProvider theme={theme2}>
                  <div className="border-2 border-transparent h-full mx-2">
                    <div className="border-2 border-transparent h-full mx-2 overflow-auto overflow-x-hidden">
                      {moneyProviders.map((MoneyProvider, index) => (
                        <MoneyProviderItems
                          key={MoneyProvider.id}
                          MoneyProvider={MoneyProvider}
                        />
                      ))}
                    </div>
                  </div>
                </ThemeProvider>
              </SwiperSlide>
            ) : (
              <></>
            )}
          </Swiper>
          <div className="flex justify-center items-center h-16 w-full">
            {swiperIndex === 0 ? (
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
                      e.currentTarget.style.borderColor = tabColors[0];
                    }}
                  >
                    <button
                      type="button"
                      className="h-5/6"
                      style={{ color: tabColors[0] }}
                      tabIndex={0}
                      onClick={async (e) => {
                        e.preventDefault();
                        await addMoneyPool();
                      }}
                      title="MoneyPoolを追加"
                    >
                      <Plus className="h-full w-full" />
                    </button>
                    <div className="w-11/12 flex gap-2 justify-start items-center p-1 relative">
                      {isEmojiPicking ? (
                        <>
                          <div className="absolute bottom-0 z-20 bg-white">
                            <Picker onEmojiClick={onEmojiClick} />
                          </div>
                          <div
                            className="h-screen w-screen fixed left-0 top-0 z-10 bg-transparent"
                            onClick={() => {
                              setIsEmojiPicking(false);
                            }}
                          />
                        </>
                      ) : (
                        <></>
                      )}
                      <input
                        type="text"
                        ref={inputMoneyPoolEmojiElement}
                        className="h-[80%] hover:border-0 focus:outline-none w-[15%] px-0"
                        defaultValue={newMoneyPoolEmoji}
                        onClick={() => {
                          setIsEmojiPicking(true);
                        }}
                        placeholder="絵文字"
                        readOnly={true}
                      />
                      <input
                        type="text"
                        className="h-[80%] hover:border-0 focus:outline-none w-[75%]"
                        ref={inputMoneyPoolNameElement}
                        onKeyDown={async (e) => {
                          if (e.key === "Enter") {
                            if (e.currentTarget) {
                              e.currentTarget.blur();
                            }
                            e.preventDefault();
                            await addMoneyPool();
                          }
                        }}
                        placeholder="名前"
                        value={newMoneyPoolName}
                        onChange={(e) => {
                          setNewMoneyPoolName(e.target.value);
                        }}
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
                      e.currentTarget.style.borderColor = tabColors[0];
                    }}
                  >
                    <div className="h-5/6  border-primary-default aspect-square">
                      <div
                        className="h-full w-full"
                        style={{ color: tabColors[0] }}
                      >
                        <Plus className="h-full w-full" />
                      </div>
                    </div>
                    追加
                  </div>
                )}
              </div>
            ) : (
              <div
                className="w-[95%] h-12 cursor-pointer"
                onClick={() => {
                  setIsAddMode2(true);
                }}
                onBlur={(fe) => {
                  console.log(fe.relatedTarget);
                  console.log(fe.currentTarget.contains(fe.relatedTarget));
                  if (!fe.currentTarget.contains(fe.relatedTarget)) {
                    setIsAddMode2(false);
                  }
                }}
                tabIndex={0}
              >
                {isAddMode2 ? (
                  <div
                    className={`flex h-12 items-center gap-2 w-full border-2 border-gray-200 rounded-full shadow-md px-2 font-Noto`}
                    onMouseOut={(e) => {
                      e.currentTarget.style.borderColor = "transparent";
                    }}
                    onMouseOver={(e) => {
                      e.currentTarget.style.borderColor = tabColors[1];
                    }}
                  >
                    <button
                      type="button"
                      className="h-5/6"
                      style={{ color: tabColors[1] }}
                      tabIndex={0}
                      onClick={async (e) => {
                        e.preventDefault();
                        await addMoneyProvider();
                      }}
                      title="MoneyProviderを追加"
                    >
                      <Plus className="h-full w-full" />
                    </button>
                    <div className="w-11/12 flex gap-2 justify-start items-center p-1">
                      <input
                        type="text"
                        className="h-[80%] hover:border-0 focus:outline-none px-0"
                        onKeyDown={(e) => {
                          if (e.key === "Enter") {
                            e.preventDefault();
                          }
                        }}
                        placeholder="名前"
                        value={newMoneyProviderName}
                        onChange={(e) => {
                          setNewMoneyProviderName(e.target.value);
                        }}
                      />
                      <input
                        type="text"
                        className="h-[80%] hover:border-0 focus:outline-none w-[30%]"
                        onKeyDown={async (e) => {
                          if (e.key === "Enter") {
                            if (e.currentTarget) {
                              e.currentTarget.blur();
                            }
                            e.preventDefault();
                            await addMoneyProvider();
                          }
                        }}
                        placeholder="残高"
                        value={newMoneyProviderBalance}
                        onChange={(e) => {
                          setNewMoneyProviderBalance(e.target.value);
                        }}
                      />
                    </div>
                  </div>
                ) : (
                  <div
                    className="flex h-12 items-center gap-2 w-full border-2 border-transparent hover:bg-gray-50  rounded-full hover:shadow-md px-2 font-Noto"
                    onMouseOut={(e) => {
                      e.currentTarget.style.borderColor = "transparent";
                    }}
                    onMouseOver={(e) => {
                      e.currentTarget.style.borderColor = tabColors[1];
                    }}
                  >
                    <div
                      className="h-5/6   aspect-square"
                      style={{ borderColor: tabColors[1] }}
                    >
                      <div
                        className="h-full w-full"
                        style={{ color: tabColors[1] }}
                      >
                        <Plus
                          className="h-full w-full"
                          style={{ borderColor: tabColors[1] }}
                        />
                      </div>
                    </div>
                    追加
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      </div>
    );
  } else {
    return <div></div>;
  }
}

function SwiperTabs({
  swiper,
  swiperIndex,
  index,
  color,
  title,
}: {
  swiper: SwiperClass | undefined;
  swiperIndex: number;
  index: number;
  color: string;
  title: string;
}) {
  return (
    <div
      className="flex border-0 rounded-t-2xl h-full justify-center px-4 font-bold font-Noto items-center min-w-max w-20 border-b-0 cursor-pointer "
      style={
        swiperIndex === index
          ? { backgroundColor: color, color: "#ffffff" }
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
      <div>{title}</div>
    </div>
  );
}

function BalanceItem({
  moneyPool,
  userID,
}: {
  moneyPool: MoneyPoolSum;
  userID: string;
}) {
  const { data: session } = useSession();
  const [isPublic, setIsPublic] = useState(moneyPool.type === "public");
  const [changePublicCheckIsOpen, setChangePublicCheckIsOpen] = useState(false);
  const [deleteCheckIsOpen, setDeleteCheckIsOpen] = useState(false);
  const [isEditEmoji, setIsEditEmoji] = useState(false);
  const [isEditName, setIsEditName] = useState(false);
  const inputName = useRef<HTMLInputElement>(null!);

  const [moneyPoolName, setMoneyPoolName] = useState(moneyPool.name);
  const [moneyPoolEmoji, setMoneyPoolEmoji] = useState(moneyPool.emoji);
  const [moneyPoolType, setMoneyPoolType] = useState(moneyPool.type);

  async function onEmojiChange(newValue: string) {
    if (session && session.user) {
      const res = await fetch(`/api/back/moneypools/${moneyPool.id}`, {
        method: "PATCH",
        headers: {
          Authorization: `Bearer ${session.user.idToken}`,
        },
        body: JSON.stringify({
          emoji: newValue,
          name: moneyPoolName,
          description: "",
          type: moneyPoolType,
        }),
      });
      if (res.ok) {
        const data = await res.json();
        console.log(data);
        setMoneyPoolEmoji(newValue);
      }
    }
  }

  async function onNameChange(newValue: string) {
    if (session && session.user) {
      const res = await fetch(`/api/back/moneypools/${moneyPool.id}`, {
        method: "PATCH",
        headers: {
          Authorization: `Bearer ${session.user.idToken}`,
        },
        body: JSON.stringify({
          emoji: moneyPoolEmoji,
          name: newValue,
          description: "",
          type: moneyPoolType,
        }),
      });
      if (res.ok) {
        const data = await res.json();
        console.log(data);
      }
    }
    setMoneyPoolName(newValue);
  }

  async function onPublicChange(newValue: boolean) {
    const type = newValue ? "public" : "private";
    if (session && session.user) {
      const res = await fetch(`/api/back/moneypools/${moneyPool.id}`, {
        method: "PATCH",
        headers: {
          Authorization: `Bearer ${session.user.idToken}`,
        },
        body: JSON.stringify({
          emoji: moneyPoolEmoji,
          name: moneyPoolName,
          description: "",
          type: type,
        }),
      });
      if (res.ok) {
        const data = await res.json();
        console.log("change public", data);
      }
    }
    setMoneyPoolType(type);
  }

  async function deleteMoneyPool(id: string) {
    if (session && session.user) {
      const res = await fetch(`/api/back/moneypools/${id}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${session.user.idToken}`,
        },
      });
      if (res.ok) {
        console.log("delete moneypool", "");
      }
    }
  }

  return (
    <div className="flex gap-4 font-Noto font-normal py-2 text-4xl items-center justify-between px-0 border-b-2 border-gray-300">
      <div className="w-10 h-10">
        {isEditEmoji ? (
          <>
            <div
              className="fixed z-20 flex justify-center items-center top-1/3 left-40 shadow-xl rounded-xl"
              onBlur={(fe) => {
                if (!fe.currentTarget.contains(fe.relatedTarget)) {
                  setIsEditEmoji(false);
                }
              }}
            >
              <Picker
                onEmojiClick={(emoji) => {
                  onEmojiChange(emoji.emoji);
                  setIsEditEmoji(false);
                }}
              />
            </div>
            <div
              className="h-screen w-screen z-10 fixed bg-transparent left-0 top-0"
              onClick={() => {
                setIsEditEmoji(false);
              }}
            />
          </>
        ) : (
          <></>
        )}
        <div
          className="cursor-pointer h-full"
          onClick={() => {
            if (!isEditEmoji) {
              setIsEditEmoji(() => {
                return true;
              });
            }
          }}
          tabIndex={0}
        >
          {moneyPoolEmoji}
        </div>
      </div>
      <div className="flex items-center justify-between w-9/12">
        <div className="w-1/2 h-10 min-w-[10px]">
          {isEditName ? (
            <div className="w-full h-full">
              <input
                aria-label="money pool name"
                ref={inputName}
                type="text"
                className="w-full"
                defaultValue={moneyPoolName}
                onBlur={(fe) => {
                  onNameChange(fe.currentTarget.value);
                  if (!fe.currentTarget.contains(fe.relatedTarget)) {
                    setIsEditName(false);
                  }
                }}
                onKeyDown={(e) => {
                  if (e.key === "Enter") {
                    onNameChange(e.currentTarget.value);
                    e.currentTarget.blur();
                  }
                }}
                autoFocus
                tabIndex={0}
              />
            </div>
          ) : (
            <div
              className="cursor-pointer w-full h-full"
              onClick={() => {
                setIsEditName((v) => {
                  return !v;
                });
              }}
              tabIndex={0}
            >
              {moneyPoolName}
            </div>
          )}
        </div>
        <div className="w-1/2 text-right">
          {moneyPool.sum.toLocaleString(undefined, {
            maximumFractionDigits: 5,
          })}
          円
        </div>
      </div>

      <div>
        <Checkbox
          onClick={(e) => {
            e.preventDefault();
            if (isPublic == false) {
              setChangePublicCheckIsOpen(true);
            } else {
              onPublicChange(false);
              setIsPublic((v) => {
                return !v;
              });
            }
          }}
          checked={isPublic}
          className="w-10"
          sx={{ "& .MuiSvgIcon-root": { fontSize: 28 } }}
        />
      </div>
      <div>
        {session?.user.sub === userID ? (
          <div
            className="w-full"
            onClick={() => {
              setDeleteCheckIsOpen(true);
            }}
            tabIndex={0}
          >
            <Image
              src="/icons/delete_FILL0_wght400_GRAD0_opsz24.svg"
              alt="x"
              width={30}
              height={30}
              style={{ objectFit: "contain", color: tabColors[0] }}
              className="min-w-[30px] max-w-[30px] cursor-pointer"
            />
          </div>
        ) : (
          <></>
        )}
      </div>
      <Modal
        isOpen={changePublicCheckIsOpen}
        className="flex justify-center items-center t-0 l-0 w-full h-full"
      >
        <div
          className="bg-transparent w-full h-full absolute z-10"
          onClick={() => setChangePublicCheckIsOpen(false)}
        />
        <div className="w-1/2 h-1/3 bg-gray-50  transform bg-opacity-90 shadow-2xl rounded-3xl flex flex-col justify-center items-center font-Noto font-semibold text-xl gap-y-20 z-50 border-primary-default border-2">
          <div>このMoney Poolを公開してもよろしいですか？</div>
          <div className="flex justify-between gap-x-8">
            <button
              type="button"
              className="bg-primary-default hover:bg-primary-dark rounded-full  text-white px-4 py-1 border-primary-default border-2 hover:border-primary-dark font-Noto font-semibold text-xl"
              onClick={() => {
                setChangePublicCheckIsOpen(false);
                onPublicChange(true);
                setIsPublic((v) => {
                  return !v;
                });
              }}
            >
              公開する
            </button>
            <button
              type="button"
              className="bg-white hover:bg-gray-100 rounded-full  text-primary-default px-4 py-1 border-primary-default border-2 font-Noto font-semibold text-xl"
              onClick={() => setChangePublicCheckIsOpen(false)}
            >
              キャンセル
            </button>
          </div>
        </div>
      </Modal>
      <Modal
        isOpen={deleteCheckIsOpen}
        className="flex justify-center items-center t-0 l-0 w-full h-full"
      >
        <div
          className="bg-transparent w-full h-full absolute z-10"
          onClick={() => setDeleteCheckIsOpen(false)}
        />
        <div className="w-1/2 h-1/3 bg-gray-50  transform bg-opacity-90 shadow-2xl rounded-3xl flex flex-col justify-center items-center font-Noto font-semibold text-xl gap-y-20 z-50 border-primary-default border-2">
          <div>Money Pool {moneyPool.name} を削除してもよろしいですか？</div>
          <div className="flex justify-between gap-x-8">
            <button
              type="button"
              className="bg-primary-default hover:bg-primary-dark rounded-full  text-white px-4 py-1 border-primary-default border-2 hover:border-primary-dark font-Noto font-semibold text-xl"
              onClick={() => {
                setDeleteCheckIsOpen(false);
                deleteMoneyPool(moneyPool.id);
                setIsPublic((v) => {
                  return !v;
                });
              }}
            >
              削除する
            </button>
            <button
              type="button"
              className="bg-white hover:bg-gray-100 rounded-full  text-primary-default px-4 py-1 border-primary-default border-2 font-Noto font-semibold text-xl"
              onClick={() => setDeleteCheckIsOpen(false)}
            >
              キャンセル
            </button>
          </div>
        </div>
      </Modal>
    </div>
  );
}

function MoneyProviderItems({
  MoneyProvider,
}: {
  MoneyProvider: MoneyProviderSum;
}) {
  const { data: session } = useSession();
  const [isEditProviderName, setIsEditProviderName] = useState(false);
  const inputProviderName = useRef<HTMLInputElement>(null!);
  const [providerName, setProviderName] = useState(MoneyProvider.name);

  const [isEditBalance, setIsEditBalance] = useState(false);
  const inputBalance = useRef<HTMLInputElement>(null!);
  const [providerBalance, setProviderBalance] = useState(MoneyProvider.balance);

  const [moneyProviderName, setMoneyProviderName] = useState(
    MoneyProvider.name
  );

  const [moneyProviderBalance, setMoneyProviderBalance] = useState(
    MoneyProvider.balance
  );

  async function onNameChanged(newValue: string) {
    if (session && session.user) {
      const res = await fetch(`/api/back/moneyproviders/${MoneyProvider.id}`, {
        method: "PATCH",
        headers: {
          Authorization: `Bearer ${session.user.idToken}`,
        },
        body: JSON.stringify({
          name: newValue,
          balance: moneyProviderBalance,
        }),
      });
      if (res.ok) {
        const data = await res.json();
        console.log(data);
      }
    }
    setMoneyProviderName(newValue);
  }

  async function onBalanceChanged(newValue: string) {
    if (session && session.user) {
      console.log("change balance data", newValue);
      const res = await fetch(`/api/back/moneyproviders/${MoneyProvider.id}`, {
        method: "PATCH",
        headers: {
          Authorization: `Bearer ${session.user.idToken}`,
        },
        body: JSON.stringify({
          name: moneyProviderName,
          balance: Number(newValue),
        }),
      });
      if (res.ok) {
        console.log("changed balance");
      }
    }
    setMoneyProviderBalance(newValue);
  }

  return (
    <div className="flex gap-4 font-Noto font-normal py-2 text-4xl items-center justify-between px-4 overflow-hidden border-b-2 border-gray-300">
      <div className="flex items-center justify-between w-full">
        <div className="w-1/2">
          {isEditProviderName ? (
            <div className="w-full h-full">
              <input
                aria-label="money provider name"
                ref={inputProviderName}
                type="text"
                className="w-full"
                defaultValue={providerName}
                value={providerName}
                onBlur={(fe) => {
                  onNameChanged(fe.currentTarget.value);
                  if (!fe.currentTarget.contains(fe.relatedTarget)) {
                    setIsEditProviderName(false);
                  }
                }}
                onKeyDown={(e) => {
                  if (e.key === "Enter") {
                    onNameChanged(e.currentTarget.value);
                    e.currentTarget.blur();
                  }
                }}
                onChange={(e) => {
                  setProviderName(e.currentTarget.value);
                }}
                autoFocus
                tabIndex={0}
              />
            </div>
          ) : (
            <div
              className="cursor-pointer h-full"
              onClick={() => {
                setIsEditProviderName((v) => {
                  return !v;
                });
              }}
              tabIndex={0}
            >
              {providerName}
            </div>
          )}
        </div>
        <div className="w-1/2 text-right">
          {isEditBalance ? (
            <div className="w-full h-full">
              <input
                aria-label="money provider balance"
                ref={inputBalance}
                type="text"
                className="w-full"
                defaultValue={providerBalance}
                value={providerBalance}
                onBlur={(fe) => {
                  onBalanceChanged(providerBalance);
                  if (!fe.currentTarget.contains(fe.relatedTarget)) {
                    setIsEditBalance(false);
                  }
                }}
                onKeyDown={(e) => {
                  if (e.key === "Enter") {
                    onBalanceChanged(providerBalance);
                    e.currentTarget.blur();
                  }
                }}
                onChange={(e) => {
                  setProviderBalance(e.currentTarget.value);
                }}
                autoFocus
                tabIndex={0}
              />
            </div>
          ) : (
            <div
              className="cursor-pointer h-full"
              onClick={() => {
                setIsEditBalance((v) => {
                  return !v;
                });
              }}
              tabIndex={0}
            >
              {providerBalance} 円
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
