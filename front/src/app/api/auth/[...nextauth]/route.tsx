import NextAuth from "next-auth";
import { JWT } from "next-auth/jwt";
import ZitadelProvider from "next-auth/providers/zitadel";
import { custom } from "openid-client";

custom.setHttpOptionsDefaults({
  timeout: 10000,
});

const authOptions = {
  providers: [
    ZitadelProvider({
      clientId: process.env.ZITADEL_CLIENT_ID as string,
      clientSecret: process.env.ZITADEL_CLIENT_SECRET as string,
      issuer: process.env.ZITADEL_URL,
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
    }) => {
      if (user) {
        token.user = user;
        const u = user as any;
        token.role = u.role;
      }
      if (account) {
        token.accessToken = account.access_token;
      }
      return token;
    },
    session: ({ session, token }: { token: JWT; session?: any }) => {
      token.accessToken;
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
