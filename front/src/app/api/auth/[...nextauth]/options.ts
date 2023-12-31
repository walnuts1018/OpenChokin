import { JWT } from "next-auth/jwt";
import ZitadelProvider from "next-auth/providers/zitadel";
import { NextAuthOptions } from "next-auth";
import Redis from 'ioredis';
import crypto from 'crypto';
import AsyncLock from 'async-lock';

let redis: Redis;

if (process.env.REDIS_SENTINEL_HOST !== undefined) {
  redis = new Redis({
    sentinels: [{
      host: process.env.REDIS_SENTINEL_HOST,
      port: Number(process.env.REDIS_SENTINEL_PORT),
    }],
    sentinelPassword: process.env.REDIS_PASSWORD,
    password: process.env.REDIS_PASSWORD,
    name: process.env.REDIS_SENTINEL_NAME || "mymaster",
    db: Number(process.env.REDIS_DB),
  });
} else if (process.env.REDIS_HOST !== undefined) {
  redis = new Redis({
    host: process.env.REDIS_HOST,
    port: Number(process.env.REDIS_PORT),
    password: process.env.REDIS_PASSWORD,
    db: Number(process.env.REDIS_DB),
  });
}

const cachePassword = process.env.CACHE_PASSWORD || "password";
const cacheKey = crypto.scryptSync(cachePassword, "salt", 32);
const lock = new AsyncLock({ timeout: 10000 });

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
      lock.acquire("refreshToken" + token.sub as string, async function (done) {
        console.debug("time:", new Date(), "firsttoken:", token.refreshToken);
        try {
          if (account) {
            token.refreshToken = account.refresh_token;
            token.idToken = account.id_token;
            token.expiresAt = account.expires_at;
          }
          else if (new Date() > new Date(token.expiresAt as number * 1000)) {
            //else if (true) {
            const cachedJsonData = await redis.get("openchokin-" + token.sub as string);
            if (cachedJsonData) {
              const cachedData = JSON.parse(cachedJsonData);
              const iv = Buffer.from(cachedData.iv, 'hex');
              const decipher = crypto.createDecipheriv('aes-256-cbc', cacheKey, iv);
              const decryptedRefreshToken = Buffer.concat([decipher.update(Buffer.from(cachedData.refreshToken, 'hex')), decipher.final()]);
              const cachedRefreshToken = decryptedRefreshToken.toString();


              console.debug("time:", new Date(), "cached:", cachedRefreshToken, "token:", token.refreshToken);
              if (cachedRefreshToken === token.refreshToken) {
                const { id_token, refresh_token, expires_at } = await refreshIDToken(token.refreshToken as string);
                token.idToken = id_token;
                token.refreshToken = refresh_token;
                token.expiresAt = expires_at;
              } else {
                console.debug("refresh token mismatch")
                const cachedExpiresAt = cachedData.expiresAt;
                const cachedIdToken = cachedData.idToken;
                token.idToken = cachedIdToken;
                token.refreshToken = cachedRefreshToken;
                token.expiresAt = cachedExpiresAt;
              }
              console.debug("time:", new Date(), "NewCached:", token.refreshToken, "NewToken:", token.refreshToken);
            }
          }
          const iv = crypto.randomBytes(16);
          const cipher = crypto.createCipheriv('aes-256-cbc', cacheKey, iv);
          const encryptedRefreshToken = Buffer.concat([cipher.update(token.refreshToken as string), cipher.final()]);
          const newCachedData = JSON.stringify({
            refreshToken: encryptedRefreshToken.toString('hex'),
            idToken: token.idToken,
            expiresAt: token.expiresAt,
            iv: iv.toString('hex'),
          })
          await redis.set("openchokin-" + token.sub as string, newCachedData, "EX", 60 * 60 * 24 * 30);
          token.error = undefined;
          return token;
        } catch (e) {
          console.error("Error refreshing token", e);
          return { ...token, error: "RefreshAccessTokenError" as const }
        } finally {
          console.debug("time:", new Date(), "finalToken:", token.refreshToken, "\n\n");
          done();
        }
      });
      return token;
    },
    session: ({ session, token }: { token: JWT; session?: any }) => {
      session.error = token.error;
      session.user.role = token.role;
      session.user.idToken = token.idToken;
      session.user.sub = token.sub;
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
    console.debug("time:", new Date(), "refresh error", data.error_description || data.error || "Unknown error");
    throw new Error(data.error_description || data.error || "Unknown error");
  }
  console.debug("token refreshed");
  return {
    id_token: data.id_token,
    refresh_token: data.refresh_token,
    expires_at: data.expires_at,
  }
}