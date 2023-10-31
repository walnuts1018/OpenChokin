import type { Metadata } from "next";
import { Nunito, Noto_Sans_JP } from "next/font/google";
import "./globals.css";
import AppleTouchIcon from "../../public/favicons/apple-touch-icon.png";
import Favicon16 from "../../public/favicons/favicon-16x16.png";
import Favicon32 from "../../public/favicons/favicon-32x32.png";
import Favicon from "../../public/favicons/favicon.ico";
import { Header } from "./Header";
import { SessionProvider } from "next-auth/react";

const title = "OpenChokin";
const description = "OpenChokinは家計簿を公開できるサービスです。";
const url = "https://openchokin.walnuts.dev";

const NunitoFont = Nunito({
  subsets: ["latin"],
  variable: "--font-Nunito",
});
const NotoFont = Noto_Sans_JP({
  weight: ["400", "500"],
  subsets: ["latin"],
  variable: "--font-Noto",
});

export const metadata: Metadata = {
  metadataBase: new URL(url),
  title: {
    default: title,
    template: `%s - ${title}`,
  },
  description: description,
  icons: [
    { rel: "icon", url: Favicon.src },
    { rel: "apple-touch-icon", sizes: "180x180", url: AppleTouchIcon.src },
    { rel: "icon", type: "image/png", sizes: "32x32", url: Favicon32.src },
    { rel: "icon", type: "image/png", sizes: "16x16", url: Favicon16.src },
  ],
  openGraph: {
    title: title,
    description,
    url,
    siteName: title,
    locale: "ja_JP",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: title,
    description,
    site: "@walnuts1018",
    creator: "@walnuts1018",
  },
  manifest: "/favicons/site.webmanifest",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ja">
      <head>
        <link rel="mask-icon" href="/safari-pinned-tab.svg" color="#f9842c" />
        <meta name="msapplication-TileColor" content="#f9842c" />
        <meta name="msapplication-TileImage" content="/mstile-144x144.png" />
        <meta name="theme-color" content="#ffffff" />
      </head>
      <body className={`${NunitoFont.variable} ${NotoFont.variable}`}>
        <Header />
        {children}
      </body>
    </html>
  );
}
