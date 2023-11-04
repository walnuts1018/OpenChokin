import { JWT } from "next-auth/jwt";
import ZitadelProvider from "next-auth/providers/zitadel";
import { NextAuthOptions } from "next-auth";

export const authOptions: NextAuthOptions = {
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
        token.idToken = account.id_token;
      }
      return token;
    },
    session: ({ session, token }: { token: JWT; session?: any }) => {
      session.user.role = token.role;
      session.user.idToken = token.idToken;
      session.user.sub = token.sub;
      console.log(session);
      return session;
    },
  },
  pages: {
    signIn: '/signin',
  },
};
