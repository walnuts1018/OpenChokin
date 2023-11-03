package psql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/walnuts1018/openchokin/back/config"
)

const (
	sslmode = "disable"
)

type DB struct {
	db *sqlx.DB
}

func dbInit() error {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%v port=%v user=%v password=%v sslmode=%v", config.Config.PostgresHost, config.Config.PostgresPort, config.Config.PostgresAdminUser, config.Config.PostgresAdminPassword, sslmode))
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %v", config.Config.PostgresDb))
	if err != nil {
		return fmt.Errorf("failed to create db: %w", err)
	}

	return nil
}

func NewDB() (*DB, error) {
	err := dbInit()
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v", config.Config.PostgresHost, config.Config.PostgresPort, config.Config.PostgresUser, config.Config.PostgresPassword, config.Config.PostgresDb, sslmode))
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	createTables(db)

	return &DB{db: db}, nil
}

func createTables(db *sqlx.DB) error {

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users(id text PRIMARY KEY)`)
	if err != nil {
		return fmt.Errorf("failed to create table users: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS money_pools(id BIGSERIAL PRIMARY KEY, name text NOT NULL, description text, is_world_public boolean NOT NULL DEFAULT false, owner_id BIGINT NOT NULL, version BIGINT NOT NULL DEFAULT 1)`)
	if err != nil {
		return fmt.Errorf("failed to create table money_pools: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS share_users(id BIGSERIAL PRIMARY KEY, money_pool_id BIGINT NOT NULL, user_id BIGINT NOT NULL)`)
	if err != nil {
		return fmt.Errorf("failed to create table share_users: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS money_transactions(id BIGSERIAL PRIMARY KEY, money_pool_id BIGINT NOT NULL, money_transaction_date date NOT NULL, title text NOT NULL, amount float8 NOT NULL, description text, is_expectation boolean NOT NULL DEFAULT false, store_id BIGINT, version BIGINT NOT NULL DEFAULT 1)`)
	if err != nil {
		return fmt.Errorf("failed to create table money_transactions: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS stores(id BIGSERIAL PRIMARY KEY, name text NOT NULL, user_id BIGINT NOT NULL)`)
	if err != nil {
		return fmt.Errorf("failed to create table stores: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS items(id BIGSERIAL PRIMARY KEY, name text NOT NULL, price_per_unit float8 NOT NULL, user_id BIGINT NOT NULL, version BIGINT NOT NULL DEFAULT 1)`)
	if err != nil {
		return fmt.Errorf("failed to create table items: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS money_transaction_items(id BIGSERIAL PRIMARY KEY, money_transaction_id BIGINT NOT NULL, item_id BIGINT NOT NULL, amount float8 NOT NULL)`)
	if err != nil {
		return fmt.Errorf("failed to create table money_transaction_items: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS money_provider(id BIGSERIAL PRIMARY KEY, name text NOT NULL, user_id BIGINT NOT NULL, balance float8 NOT NULL, version BIGINT NOT NULL DEFAULT 1)`)
	if err != nil {
		return fmt.Errorf("failed to create table money_provider: %w", err)
	}

	return nil
}

func (db *DB) Close() error {
	return db.db.Close()
}
