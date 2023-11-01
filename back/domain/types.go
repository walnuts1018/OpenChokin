package domain

import "time"

type Transaction struct {
	ID              string    `db:"id"`
	MoneyPoolID     string    `db:"money_pool_id"`
	TransactionDate time.Time `db:"transaction_date"`
	Title           string    `db:"title"`
	Amount          float64   `db:"amount"` //金額
	Labels          []string  `db:"labels"`
	IsWorldPublic   bool      `db:"is_world_public"`
	ShareUserIDs    []string  `db:"share_user_ids"`
	Expectation     bool      `db:"expectation"`
	StoreID         string    `db:"store_id"`
	ItemIDs         []string  `db:"item_ids"`
	Description     string    `db:"description"`
}

type User struct {
	ID           string   `db:"id"`
	MoneyPoolIDs []string `db:"money_pool_ids"`
}

type Store struct {
	ID     string `db:"id"`
	Name   string `db:"name"`
	UserID string `db:"user_id"`
}

type MoneyPool struct {
	ID            string   `db:"id"`
	Name          string   `db:"name"`
	Color         string   `db:"color"`
	IsWorldPublic bool     `db:"is_world_public"`
	ShareUserIDs  []string `db:"share_user_ids"`
}

type Item struct {
	ID           string  `db:"id"`
	Name         string  `db:"name"`
	PricePerUnit float64 `db:"price_per_unit"`
	UserID       string  `db:"user_id"`
}

type GetTransactionHints struct {
	PartitioningKeys []string // YYYY-MM
}
