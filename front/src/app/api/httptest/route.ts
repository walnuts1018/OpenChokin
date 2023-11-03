import type { NextApiRequest } from "next";
import { NextResponse } from "next/server";

export async function POST(request: NextApiRequest) {
  const result = await fetch("https://httptest.walnuts.dev", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ hello: "world" }),
  });
  const text = await result.text()
  return NextResponse.json({ text });
}
