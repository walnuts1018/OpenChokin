"use client";
import ZitadelForm from "./Zitadelform";

export default function Page({
  searchParams,
}: {
  searchParams: { [key: string]: string | string[] | undefined };
}) {
  const callbackUrl = searchParams.callbackUrl;

  return (
    <div className="page h-screen text-black bg-white flex flex-col">
      <div className="header-space h-20" />
      <div className="signin">
        <div className="card">
          <div className="provider">
            <ZitadelForm callbackUrl={callbackUrl} />
          </div>
        </div>
      </div>
    </div>
  );
}
