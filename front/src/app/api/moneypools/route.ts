import { NextResponse, NextRequest } from "next/server";
import { BackendEndpoint } from "../endpoint";

export async function GET(request: NextRequest) {
  const query = new URLSearchParams(new URL(request.url).searchParams);
  const result = await fetch(new URL(`/moneypools`, BackendEndpoint) + "?" + query, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });
  const data = await result.json();
  console.log(data);
  return NextResponse.json(data);
}
