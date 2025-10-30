package entity

type PromoCodeUsed struct {
	Id             int64    `gorm:"primaryKey;autoIncrement;column:id"`
	PromoCodeID    uint     `gorm:"column:promo_code_id;not null" json:"promo_code_id"`
	PromoCode      string   `gorm:"column:promo_code;type:varchar(255);not null" json:"promo_code"`
	CustomerID     string   `gorm:"column:customer_id;not null" json:"customer_id"`
	OrderID        *string  `gorm:"column:order_id" json:"order_id,omitempty"`
	Platform       string   `gorm:"column:platform;type:varchar(50);not null" json:"platform"`
	DiscountAmount *float64 `gorm:"column:discount_amount;type:decimal(13,2)" json:"discount_amount,omitempty"`
	OrderTotal     *float64 `gorm:"column:order_total;type:decimal(13,2)" json:"order_total,omitempty"`
	Status         string   `gorm:"column:status;type:enum('used','cancelled','refunded');default:'used'" json:"status"`
	UsedAt         uint64   `gorm:"column:used_at;not null" json:"used_at"`
	CreatedAt      uint64   `gorm:"column:created_at;not null;default:0" json:"created_at"`
	UpdatedAt      uint64   `gorm:"column:updated_at;not null;default:0" json:"updated_at"`
}

func (p PromoCodeUsed) TableName() string {
	return "promo_code_used"
}

func (p PromoCodeUsed) PrimaryKey() string {
	return "id"
}

func (p PromoCodeUsed) GetAllowedOrderFields() []string {
	return []string{"id", "promo_code", "customer_id", "order_id", "used_at", "created_at", "updated_at"}
}

func (p PromoCodeUsed) GetAllowedWhereFields() []string {
	return []string{"id", "promo_code_id", "promo_code", "customer_id", "order_id", "platform", "status"}
}
