export type MoneyPool = {
  id: number;
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
  Payments: MoneyTransaction[];
};

export type MoneyPoolSum = {
  id: string;
  name: string;
  Sum: number;
  Type: "private" | "public" | "restricted";
  emoji: string;
};

export type MoneyTransaction = {
  id: number;
  date: Date;
  title: string;
  amount: number;
};

export type MoneyProvider = {
  id: number;
  name: string;
  balance: number;
};
