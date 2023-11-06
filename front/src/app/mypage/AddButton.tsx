import { useState, useRef, useEffect } from "react";
import { Plus } from "react-feather";
import { useSession } from "next-auth/react";

export function AddButton({
  color,
  moneyPoolID,
  setReloadContext,
}: {
  color: string;
  moneyPoolID: string;
  setReloadContext: React.Dispatch<React.SetStateAction<{}>>;
}) {
  const { data: session } = useSession();
  const [isAddMode, setIsAddMode] = useState(false);
  const inputEl = useRef<HTMLInputElement>(null!);
  const [transactionDate, setTransactionDate] = useState(
    new Date().toISOString().split("T")[0]
  );
  const [transactionTitle, setTransactionTitle] = useState("");
  const [transactionAmount, setTransactionAmount] = useState("");

  useEffect(() => {
    if (inputEl.current) {
      console.log("focus");
      inputEl.current.focus();
    }
  }, [isAddMode]);

  async function addTransaction() {
    if (session && session.user) {
      const res = await fetch(`/api/back/moneypools/${moneyPoolID}/payments`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${session.user.idToken}`,
        },
        body: JSON.stringify({
          title: transactionTitle,
          amount: Number(transactionAmount),
          description: "",
          is_planned: false,
          date: transactionDate,
        }),
      });
      if (res.ok) {
        setTransactionTitle("");
        setTransactionAmount("");
        setReloadContext({});
      } else {
        console.log("post transaction error", res);
      }
    }
  }

  return (
    <div
      className="w-[95%] h-12 cursor-pointer"
      onClick={() => {
        setIsAddMode(true);
      }}
      onBlur={(fe) => {
        console.log(fe.relatedTarget);
        console.log(fe.currentTarget.contains(fe.relatedTarget));
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
            e.currentTarget.style.borderColor = color;
          }}
        >
          <button
            title="追加"
            className="h-5/6"
            style={{ color: color }}
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
              ref={inputEl}
              className="h-[80%] hover:border-0 focus:outline-none w-[15%] px-0"
              onKeyDown={(e) => {
                if (e.key === "Enter") {
                  e.preventDefault();
                }
              }}
              value={transactionDate}
              placeholder="日付"
              onChange={(e) => {
                setTransactionDate(e.target.value);
              }}
            />
            <input
              className="h-[80%] hover:border-0 focus:outline-none w-[75%]"
              onKeyDown={(e) => {
                if (e.key === "Enter") {
                  e.preventDefault();
                }
              }}
              placeholder="タイトル"
              value={transactionTitle}
              onChange={(e) => {
                setTransactionTitle(e.target.value);
              }}
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
              value={transactionAmount}
              onChange={(e) => {
                setTransactionAmount(e.target.value);
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
            e.currentTarget.style.borderColor = color;
          }}
        >
          <div className="h-5/6  border-primary-default aspect-square">
            <div className="h-full w-full" style={{ color: color }}>
              <Plus className="h-full w-full" />
            </div>
          </div>
          追加
        </div>
      )}
    </div>
  );
}
