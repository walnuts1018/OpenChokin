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
