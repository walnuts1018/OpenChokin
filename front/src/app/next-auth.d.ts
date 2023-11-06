import { DefaultSession } from "next-auth";

declare module "next-auth" {
  interface Session {
    user: {
      refreshToken?: string;
      exiresAt?: Date
      idToken?: string;
      sub?: string;
    } & DefaultSession["user"];
  }
}