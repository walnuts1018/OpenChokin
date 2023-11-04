import { DefaultSession } from "next-auth";

declare module "next-auth" {
  interface Session {
    user: {
      idToken?: string;
      sub?: string;
    } & DefaultSession["user"];
  }
}