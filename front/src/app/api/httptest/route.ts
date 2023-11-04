import { NextResponse, NextRequest } from "next/server";

const endpoint = "https://httpbin.org/get";
export async function POST(request: NextRequest) {
  const result = await fetch(endpoint, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ hello: "world" }),
  });
  const text = await result.text()
  return NextResponse.json({ text });
}

export async function GET(request: NextRequest) {
  const query = new URLSearchParams(new URL(request.url).searchParams);
  const result = await fetch(new URL(`/moneypools`, endpoint) + "?" + query, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });
  const data = await result.json();
  console.log(data);
  return NextResponse.json(data);
}
