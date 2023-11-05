export type MoneyPool = {
  id: string;
  name: string;
  description: string;
  is_world_public: boolean;
  owner_id: number;
  color: string;
  amount: number;
  emoji: string;
  transactions: MoneyTransaction[];
};

export type MoneyPoolResponse = {
  id: string;
  name: string;
  description: string;
  type: "private" | "public" | "restricted";
  payments: MoneyPoolResponsePayment[];
};

export type MoneyPoolResponsePayment = {
  id: string;
  date: string;
  title: string;
  amount: number;
};

export type MoneyPoolSum = {
  id: string;
  name: string;
  sum: number;
  type: "private" | "public" | "restricted";
  emoji: string;
};

export type MoneyTransaction = {
  id: string;
  date: Date;
  title: string;
  amount: number;
};

export type MoneyProviderSum = {
  id: string;
  name: string;
  balance: number;
};
