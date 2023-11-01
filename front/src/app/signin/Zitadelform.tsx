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
          <button
            type="submit"
            className="flex gap-x-2 button bg-primary-default  p-5 rounded-full hover:bg-primary-dark shadow-md shadow-gray-200 items-center px-5"
          >
            <div>
              <span className="font-Noto text-xl font-semibold text-white">
                Walnuts.dev アカウントでサインイン
              </span>
            </div>
          </button>
        </form>
      ) : (
        <div className="flex flex-col justify-center items-center gap-y-2">
          <div className="font-Noto">内部エラー</div>
          <div className="font-Noto">
            サインイン方法が見つかりませんでした。
          </div>
        </div>
      )}
    </>
  );
}
