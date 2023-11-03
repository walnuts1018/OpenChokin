import { MoneyPool } from "./type";

export function Balance({
  className,
  user,
  moneypools,
}: {
  children?: React.ReactNode;
  className?: string;
  user: {
    name?: string | null | undefined;
    email?: string | null | undefined;
    image?: string | null | undefined;
  };
  moneypools: MoneyPool[];
}) {
  return (
    <div>
      Signed in as {user.email} <br />
    </div>
  );
}
