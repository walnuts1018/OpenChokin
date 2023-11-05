-- 公開タイプの列挙型を定義
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'public_type') THEN
        CREATE TYPE public_type AS ENUM ('private', 'public', 'restricted');
    END IF;
END$$;

-- ユーザーテーブル
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY
);

-- ユーザーグループテーブル
CREATE TABLE IF NOT EXISTS user_groups (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    creator_id BIGINT NOT NULL,
    FOREIGN KEY (creator_id) REFERENCES users(id)
);

-- マネープールテーブル
CREATE TABLE IF NOT EXISTS money_pool (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL,
    owner_id BIGINT NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES users(id)
);

-- マネープロバイダーテーブル
CREATE TABLE IF NOT EXISTS money_provider (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    creator_id BIGINT NOT NULL,
    balance DECIMAL(19,4) NOT NULL CHECK (balance >= 0),
    FOREIGN KEY (creator_id) REFERENCES users(id)
);

-- 店舗テーブル
CREATE TABLE IF NOT EXISTS store (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    creator_id BIGINT NOT NULL,
    FOREIGN KEY (creator_id) REFERENCES users(id)
);

-- 商品テーブル
CREATE TABLE IF NOT EXISTS item (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    creator_id BIGINT NOT NULL,
    FOREIGN KEY (creator_id) REFERENCES users(id)
);

-- ラベルテーブル
CREATE TABLE IF NOT EXISTS label (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    creator_id BIGINT NOT NULL,
    FOREIGN KEY (creator_id) REFERENCES users(id)
);

-- 取引テーブル
CREATE TABLE IF NOT EXISTS payment (
    id BIGSERIAL PRIMARY KEY,
    money_pool_id BIGINT NOT NULL,
    date DATE NOT NULL,
    title VARCHAR(255) NOT NULL,
    amount DECIMAL(19,4) NOT NULL,
    description TEXT,
    is_planned BOOLEAN NOT NULL,
    store_id BIGINT,
    FOREIGN KEY (money_pool_id) REFERENCES money_pool(id),
    FOREIGN KEY (store_id) REFERENCES store(id)
);

-- 商品取引テーブル
CREATE TABLE IF NOT EXISTS item_payment (
    payment_id BIGINT NOT NULL,
    item_id BIGINT NOT NULL,
    quantity BIGINT NOT NULL CHECK (quantity > 0),
    PRIMARY KEY (payment_id, item_id),
    FOREIGN KEY (payment_id) REFERENCES payment(id),
    FOREIGN KEY (item_id) REFERENCES item(id)
);

-- ユーザーグループ所属テーブル
CREATE TABLE IF NOT EXISTS user_group_membership (
    group_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    PRIMARY KEY (group_id, user_id),
    FOREIGN KEY (group_id) REFERENCES user_groups(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 限定公開範囲テーブル
CREATE TABLE IF NOT EXISTS restricted_publication_scope (
    pool_id BIGINT NOT NULL,
    group_id BIGINT NOT NULL,
    PRIMARY KEY (pool_id, group_id),
    FOREIGN KEY (pool_id) REFERENCES money_pool(id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES user_groups(id)
);
