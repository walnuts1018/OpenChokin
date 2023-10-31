import NextAuth from "next-auth";
import ZitadelProvider from "next-auth/providers/zitadel";
import { custom } from 'openid-client';

const client_id = process.env.ZITADEL_CLIENT_ID;
const client_secret = process.env.ZITADEL_CLIENT_SECRET;
const zitadel_url = process.env.ZITADEL_URL;

if (!client_id || !client_secret) {
  throw new Error("ZITADEL_CLIENT_ID and ZITADEL_CLIENT_SECRET must be set");
}

custom.setHttpOptionsDefaults({
  timeout: 10000,
});

export const authOptions = {
  providers: [
    ZitadelProvider({
      clientId: client_id,
      clientSecret: client_secret,
      issuer: zitadel_url,
    }),
  ],
  callbacks: {
    jwt: async ({ token, user, account, profile, isNewUser }) => {
      if (user) {
        token.user = user;
        const u = user as any
        token.role = u.role;
      }
      if (account) {
        token.accessToken = account.access_token
      }
      return token;
    },
  },
  session: ({ session, token }) => {
    token.accessToken
    return {
      ...session,
      user: {
        ...session.user,
        role: token.role,
      },
    };
  },
};

const handler = NextAuth(authOptions);

export { handler as GET, handler as POST };
