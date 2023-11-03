import { MoneyPool } from "./type";
import { Swiper, SwiperSlide, SwiperClass } from "swiper/react";
import { useState } from "react";
import Checkbox from "@mui/material/Checkbox";
import { ThemeProvider, createTheme, styled } from "@mui/material/styles";
import Modal from "react-modal";

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
  user,
  moneypools,
}: {
  children?: React.ReactNode;
  className?: string;
  user: {
    name?: string | null | undefined;
    email?: string | null | undefined;
    image?: string | null | undefined;
  };
  moneypools: MoneyPool[];
}) {
  const [swiper, setSwiper] = useState<SwiperClass>();
  const [swiperIndex, setSwiperIndex] = useState(0);
  const [changePublicCheckIsOpen, setChangePublicCheckIsOpen] = useState(false);
  const [forceReload, setForceReload] = useState(0);

  Modal.setAppElement("body");
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
        <SwiperTabs
          swiper={swiper}
          swiperIndex={swiperIndex}
          index={1}
          color={tabColors[1]}
          title="Money Providers"
        />
      </div>
      <div
        className={
          "border-2 rounded-3xl p-1 h-[calc(100%-2rem)] overflow-hidden w-full"
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
          className="flex w-full h-full"
        >
          <SwiperSlide className="">
            <ThemeProvider theme={theme1}>
              <div className="border-2 border-transparent h-full mx-2 overflow-auto overflow-x-hidden">
                {moneypools.map((moneyPool, index) => (
                  <div
                    key={moneyPool.id}
                    className="flex gap-4 font-Noto font-normal py-2 text-4xl items-center justify-between px-4 overflow-hidden border-b-2 border-gray-300"
                  >
                    <div className="w-10">{moneyPool.emoji}</div>
                    <div className="flex items-center justify-between w-9/12">
                      <div className="w-1/2">{moneyPool.name}</div>
                      <div className="w-1/2 text-right">
                        {moneyPool.amount.toLocaleString(undefined, {
                          maximumFractionDigits: 5,
                        })}
                        円
                      </div>
                    </div>

                    <div>
                      <Checkbox
                        onChange={(e) => {
                          if (e.target.checked) {
                            setChangePublicCheckIsOpen(true);
                            setForceReload((forceReload + 1) % 2);
                          }
                        }}
                        value={moneyPool.is_world_public}
                        className="w-10"
                        sx={{ "& .MuiSvgIcon-root": { fontSize: 28 } }}
                      />
                    </div>
                  </div>
                ))}
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
                      className="bg-primary-default hover:bg-primary-dark rounded-full  text-white px-4 py-1 border-primary-default border-2 hover:border-primary-dark font-Noto font-semibold text-xl"
                      onClick={() => setChangePublicCheckIsOpen(false)}
                    >
                      公開する
                    </button>
                    <button
                      className="bg-white hover:bg-gray-100 rounded-full  text-primary-default px-4 py-1 border-primary-default border-2 font-Noto font-semibold text-xl"
                      onClick={() => setChangePublicCheckIsOpen(false)}
                    >
                      キャンセル
                    </button>
                  </div>
                </div>
              </Modal>
            </ThemeProvider>
          </SwiperSlide>
          <SwiperSlide className="">
            <ThemeProvider theme={theme2}>
              <div className="border-2 border-transparent h-full mx-2">add</div>
            </ThemeProvider>
          </SwiperSlide>
        </Swiper>
      </div>
    </div>
  );
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
