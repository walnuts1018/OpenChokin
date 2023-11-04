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
    type public_type NOT NULL,
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

-- 制約：限定公開範囲のプールIDに対応するマネープールの公開タイプは限定公開である
ALTER TABLE restricted_publication_scope 
ADD CONSTRAINT fk_money_pool_restricted
CHECK ((SELECT type FROM money_pool WHERE id = pool_id) = 'restricted');

-- 初期データ追加
-- usersテーブルにデータを挿入
INSERT INTO users (id) VALUES (DEFAULT), (DEFAULT), (DEFAULT);

-- user_groupsテーブルにデータを挿入
INSERT INTO user_groups (name, creator_id) VALUES
('Developers Group', 1),
('Designers Group', 1),
('Managers Group', 2);

-- money_poolテーブルにデータを挿入
INSERT INTO money_pool (name, description, type, owner_id) VALUES
('Fund for Emergencies', 'Emergency funds for unexpected expenses', 'private', 1),
('Office Party Fund', 'Savings for annual office parties', 'public', 2),
('Project X Budget', 'Budget allocated for Project X', 'restricted', 3);

-- money_providerテーブルにデータを挿入
INSERT INTO money_provider (name, creator_id, balance) VALUES
('John’s Wallet', 1, 1000.0000),
('Anna’s Savings', 2, 1500.0000),
('Company Petty Cash', 3, 500.0000);

-- storeテーブルにデータを挿入
INSERT INTO store (name, creator_id) VALUES
('Tech Gadgets', 1),
('Office Supplies Store', 2),
('Bookstore', 3);

-- itemテーブルにデータを挿入
INSERT INTO item (name, creator_id) VALUES
('Laptop', 1),
('Ergonomic Keyboard', 1),
('Planner Notebook', 2);

-- labelテーブルにデータを挿入
INSERT INTO label (name, creator_id) VALUES
('Electronics', 1),
('Office Equipment', 1),
('Stationery', 2);

-- paymentテーブルにデータを挿入
INSERT INTO payment (money_pool_id, date, title, amount, description, is_planned, store_id) VALUES
(1, '2023-01-15', 'Emergency Repair', 200.0000, 'Repairing office printer', TRUE, NULL),
(2, '2023-02-20', 'Office Party', 300.0000, 'Annual office party expenses', FALSE, NULL),
(3, '2023-03-12', 'Project X Software', 400.0000, 'Software purchase for Project X', TRUE, 1);

-- item_paymentテーブルにデータを挿入
INSERT INTO item_payment (payment_id, item_id, quantity) VALUES
(1, 1, 1),
(2, 2, 2),
(3, 3, 3);

-- user_group_membershipテーブルにデータを挿入
INSERT INTO user_group_membership (group_id, user_id) VALUES
(1, 1),
(1, 2),
(2, 3);

-- restricted_publication_scopeテーブルにデータを挿入
INSERT INTO restricted_publication_scope (pool_id, group_id) VALUES
(3, 1),
(3, 2),
(3, 3);
