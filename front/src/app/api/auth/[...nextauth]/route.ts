import NextAuth from "next-auth";
import { custom } from "openid-client";
import { authOptions } from "./options";

custom.setHttpOptionsDefaults({
  timeout: 10000,
});

const handler = NextAuth(authOptions);

export { handler as GET, handler as POST };
