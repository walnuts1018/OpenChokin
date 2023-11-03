```mermaid
erDiagram
  users ||--o{ money_pools : ""
  money_pools ||--o{ money_transactions : ""
  money_transactions ||--o{ items : ""
  money_transactions ||--o{labels: ""
  users ||--o{ money_provider : ""

  users {
    text id "UserID OIDCのSubjectと同じ"
  }

  money_pools {
    BIGSERIAL NOT NULL id "MoneyPoolID"
    text NOT NULL name "MoneyPoolの名前"
    text description "MoneyPoolの説明"
    boolean NOT NULL is_world_public "MoneyPoolが公開されているかどうか, default false"
    BIGINT NOT NULL owner_id "MoneyPoolのオーナーのID"
    BIGINT NOT NULL version "MoneyPoolのバージョン"
  }

  share_users {
    BIGSERIAL NOT NULL id "ShareUserID"
    BIGINT NOT NULL money_pool_id "共有されているMoneyPoolのID"
    BIGINT NOT NULL user_id "共有されているユーザーのID"
  }

  money_transactions {
    BIGSERIAL NOT NULL id "TransactionID"
    BIGINT NOT NULL money_pool_id "取引が属するMoneyPoolのID"
    date NOT NULL money_transaction_date "取引が発生した日付"
    title NOT NULL title "取引のタイトル"
    float8 NOT NULL amount "取引の金額"
    text description "取引の説明"
    boolean NOT NULL is_expectation "取引が予定かどうか, default:false"
    BIGINT store_id "取引が発生した店舗のID"
    BIGINT NOT NULL version "取引のバージョン"
  }

  stores {
    BIGSERIAL NOT NULL id "StoreID"
    text NOT NULL name "店舗の名前"
    BIGINT NOT NULL user_id "店舗を登録したユーザーのID"
  }

  items {
    BIGSERIAL NOT NULL id "ItemID"
    text NOT NULL name "購入品の名前"
    float8 NOT NULL price_per_unit "単価"
    BIGINT NOT NULL user_id "購入品を登録したユーザーのID"
    BIGINT NOT NULL version "購入品のバージョン"
  }

  money_transaction_items {
    BIGSERIAL NOT NULL id "TransactionItemID"
    BIGINT NOT NULL money_transaction_id "購入品が属する取引のID"
    BIGINT NOT NULL item_id "購入品のID"
    float8 NOT NULL amount "購入品の個数"
  }

  money_provider {
    BIGSERIAL NOT NULL id "MoneyProviderID"
    text NOT NULL name "MoneyProviderの名前"
    BIGINT NOT NULL user_id "MoneyProviderを登録したユーザーのID"
    float8 NOT NULL balance "MoneyProviderの残高"
    BIGINT NOT NULL version "MoneyProviderのバージョン"
  }
```
