import NextAuth from "next-auth";
import { JWT } from "next-auth/jwt";
import ZitadelProvider from "next-auth/providers/zitadel";
import { custom } from 'openid-client';

var client_id = process.env.ZITADEL_CLIENT_ID;
var client_secret = process.env.ZITADEL_CLIENT_SECRET;
var zitadel_url = process.env.ZITADEL_URL;

if (!client_id || !client_secret) {
  throw new Error("ZITADEL_CLIENT_ID and ZITADEL_CLIENT_SECRET must be set");
}

custom.setHttpOptionsDefaults({
  timeout: 10000,
});

const authOptions = {
  providers: [
    ZitadelProvider({
      clientId: client_id,
      clientSecret: client_secret,
      issuer: zitadel_url,
    }),
  ],
  callbacks: {
    jwt: async ({ token, user, account, profile, isNewUser }: { token: JWT, user?: any, account?: any, profile?: any, isNewUser?: boolean }) => {
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
    session: ({ session, token }: { token: JWT, session?: any }) => {
      token.accessToken
      return {
        ...session,
        user: {
          ...session.user,
          role: token.role,
        },
      };
    },
  },
};

const handler = NextAuth(authOptions);

export { handler as GET, handler as POST };
