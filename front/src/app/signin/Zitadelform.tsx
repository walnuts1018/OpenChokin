import { getCsrfToken, getProviders } from "next-auth/react";

export default async function ZitadelForm({
  callbackUrl,
}: {
  callbackUrl: string | string[] | undefined;
}) {
  const csrfToken = await getCsrfToken();
  const providers = await getProviders();
  var signinUrl = "";
  var hasZitadel = providers?.hasOwnProperty("zitadel") ?? providers !== null;
  if (providers !== null && hasZitadel) {
    signinUrl = providers.zitadel.signinUrl;
  }

  return (
    <>
      {callbackUrl !== null && csrfToken !== undefined && hasZitadel ? (
        <form action={signinUrl} method="POST">
          <input type="hidden" name="csrfToken" value={csrfToken} />
          <input type="hidden" name="callbackUrl" value={callbackUrl} />
          <button type="submit" className="button">
            <span>Sign in with ZITADEL</span>
          </button>
        </form>
      ) : (
        <div>no zitadel</div>
      )}
    </>
  );
}
