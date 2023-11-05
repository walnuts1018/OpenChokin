package domain

import "time"

type PublicType string

const (
	PublicTypePrivate    PublicType = "private"
	PublicTypePublic     PublicType = "public"
	PublicTypeRestricted PublicType = "restricted"
)

type User struct {
	ID string `db:"id" json:"id"`
}

type UserGroup struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	CreatorID string `db:"creator_id"`
}

type MoneyPool struct {
	ID          string     `db:"id"`
	Name        string     `db:"name"`
	Description string     `db:"description"`
	Type        PublicType `db:"type"`
	OwnerID     string     `db:"owner_id"`
}

type MoneyProvider struct {
	ID        string  `db:"id"`
	Name      string  `db:"name"`
	CreatorID string  `db:"creator_id"`
	Balance   float64 `db:"balance"`
}

type Store struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	CreatorID string `db:"creator_id"`
}

type Item struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	CreatorID string `db:"creator_id"`
}

type Label struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	CreatorID string `db:"creator_id"`
}

type Payment struct {
	ID          string    `db:"id"`
	MoneyPoolID string    `db:"money_pool_id"`
	Date        time.Time `db:"date"`
	Title       string    `db:"title"`
	Amount      float64   `db:"amount"`
	Description string    `db:"description"`
	IsPlanned   bool      `db:"is_planned"`
	StoreID     *string   `db:"store_id"`
}

type ItemPayment struct {
	PaymentID string `db:"payment_id"`
	ItemID    string `db:"item_id"`
	Quantity  int64  `db:"quantity"`
}

type UserGroupMembership struct {
	GroupID string `db:"group_id"`
	UserID  string `db:"user_id"`
}

type RestrictedPublicationScope struct {
	PoolID  string `db:"pool_id"`
	GroupID string `db:"group_id"`
}
