import { useState, useRef, useEffect } from "react";
import { Plus } from "react-feather";

export function AddButton({ color }: { color: string }) {
  const [isAddMode, setIsAddMode] = useState(false);
  const inputEl = useRef<HTMLInputElement>(null!);
  useEffect(() => {
    if (inputEl.current) {
      console.log("focus");
      inputEl.current.focus();
    }
  }, [isAddMode]);

  async function addTransaction(input: HTMLInputElement) {
    const res = await fetch("/api/httptest", {
      method: "POST",
      body: JSON.stringify({
        name: input.value,
      }),
    });
    const data = await res.json();
    console.log(data);
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
              await addTransaction(inputEl.current);
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
              value={new Date().toISOString().split("T")[0]}
              placeholder="日付"
            />
            <input
              className="h-[80%] hover:border-0 focus:outline-none w-[75%]"
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
                  await addTransaction(inputEl.current);
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
