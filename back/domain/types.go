package domain

import "time"

type MoneyTransaction struct {
	ID              int64     `db:"id"`
	MoneyPoolID     int64     `db:"money_pool_id"`
	TransactionDate time.Time `db:"transaction_date"`
	Title           string    `db:"title"`
	Amount          float64   `db:"amount"` //金額
	Description     string    `db:"description"`
	IsWorldPublic   bool      `db:"is_world_public"`
	IsExpectation   bool      `db:"is_expectation"`
	StoreID         int64     `db:"store_id"`
	Version         int64     `db:"version"`
}

type User struct {
	ID string `db:"id"`
}

type Store struct {
	ID      int64  `db:"id"`
	Name    string `db:"name"`
	UserID  int64  `db:"user_id"`
	Version int64  `db:"version"`
}

type MoneyPool struct {
	ID            int64  `db:"id"`
	Name          string `db:"name"`
	Description   string `db:"description"`
	Color         string `db:"color"`
	IsWorldPublic bool   `db:"is_world_public"`
	OwnerID       int64  `db:"owner_id"`
	Version       int64  `db:"version"`
}

type Item struct {
	ID           int64   `db:"id"`
	Name         string  `db:"name"`
	PricePerUnit float64 `db:"price_per_unit"`
	UserID       int64   `db:"user_id"`
	Version      int64   `db:"version"`
}

type MoneyProvider struct {
	ID      int64   `db:"id"`
	Name    string  `db:"name"`
	UserID  int64   `db:"user_id"`
	Balance float64 `db:"balance"`
	Version int64   `db:"version"`
}
