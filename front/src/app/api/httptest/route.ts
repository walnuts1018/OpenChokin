import { NextResponse, NextRequest } from "next/server";

export async function POST(request: NextRequest) {
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
