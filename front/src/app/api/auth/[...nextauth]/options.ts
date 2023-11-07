import { JWT } from "next-auth/jwt";
import ZitadelProvider from "next-auth/providers/zitadel";
import { NextAuthOptions } from "next-auth";
import Redis from 'ioredis';
import Redlock from "redlock";

const redis = new Redis({
  host: process.env.REDIS_HOST,
  port: 6379,
  password: process.env.REDIS_PASSWORD,
  db: 1,
});
const redlock = new Redlock([redis], {
  driftFactor: 0.01,
  retryCount: 10,
  retryDelay: 200,
  retryJitter: 200,
  automaticExtensionThreshold: 500
});

export const authOptions: NextAuthOptions = {
  providers: [
    ZitadelProvider({
      clientId: process.env.ZITADEL_CLIENT_ID as string,
      clientSecret: process.env.ZITADEL_CLIENT_SECRET as string,
      issuer: process.env.ZITADEL_URL,
      authorization: { params: { scope: "openid email profile offline_access" } },
    }),
  ],
  callbacks: {
    jwt: async ({
      token,
      user,
      account,
      profile,
      isNewUser,
    }: {
      token: JWT;
      user?: any;
      account?: any;
      profile?: any;
      isNewUser?: boolean;
      session?: any;
    }) => {
      if (user) {
        token.role = user.role;
      }
      if (account) {
        token.refreshToken = account.refresh_token;
        token.idToken = account.id_token;
        token.expiresAt = account.expires_at;
      }
      else if (new Date() > new Date(token.expiresAt as number * 1000)) {
        try {
          const { id_token, refresh_token, expires_at } = await refreshIDToken(token.refreshToken as string);
          token.idToken = id_token;
          token.refreshToken = refresh_token;
          token.expiresAt = expires_at;
          console.log("Refreshed token", token);
          return token;
        } catch (e) {
          console.error("Error refreshing token", e);
          return { ...token, error: "RefreshAccessTokenError" as const }
        }
      }
      return token;
    },
    session: ({ session, token }: { token: JWT; session?: any }) => {
      session.error = token.error;
      session.user.role = token.role;
      session.user.idToken = token.idToken;
      session.user.sub = token.sub;
      //console.debug(session);
      return session;
    },
  },
  pages: {
    signIn: '/signin',
  },
};


const refreshIDToken = async (refreshToken: string) => {
  const response = await fetch(`${process.env.ZITADEL_URL}/oauth/v2/token`, {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: new URLSearchParams({
      grant_type: "refresh_token",
      client_id: process.env.ZITADEL_CLIENT_ID as string,
      client_secret: process.env.ZITADEL_CLIENT_SECRET as string,
      refresh_token: refreshToken,
    }),
  });
  const data = await response.json();
  //console.log("Refreshed Data:", data);
  if (!response.ok) {
    throw new Error(data.error_description || data.error || "Unknown error");
  }

  return {
    id_token: data.id_token,
    refresh_token: data.refresh_token,
    expires_at: data.expires_at,
  }
}