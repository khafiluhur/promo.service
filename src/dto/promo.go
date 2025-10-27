package dto

type PromoCode struct {
	ID            uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	TermCondition string  `json:"term_condition"`
	StartDate     uint    `json:"start_date"`
	EndDate       uint    `json:"end_date"`
	Banner        string  `json:"banner"`
	Rules         string  `json:"rules" gorm:"type:enum('otg','sp','pw')"`
	Amount        float64 `json:"amount"`
	ProductSlug   string  `json:"product_slug"`
	PaymentMethod float64 `json:"payment_method"`
	PromoAction   string  `json:"promo_action" gorm:"type:enum('fixed','percent')"`
	PromoType     string  `json:"promo_type"`
	Type          string  `json:"type"`
	PromoCode     string  `json:"promo_code"`
	CustomerLimit *int    `json:"customer_limit,omitempty"`
	NewCustomer   bool    `json:"new_customer"`
	Quantity      int     `json:"quantity"`
	Status        string  `json:"status"`
	SpecialPromo  *int    `json:"special_promo,omitempty"`
	IsDisplay     bool    `json:"is_display"`
	Platform      string  `json:"platform"`
	CreatedAt     int64   `json:"created_at"`
	UpdatedAt     int64   `json:"updated_at"`
}

type StrikeThoughtPrice struct {
	ID            uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Description   string `json:"description"`
	ProductSlug   string `json:"product_slug"`
	Departure     *int   `json:"departure,omitempty"`
	TermCondition string `json:"term_condition"`
	StartDate     uint   `json:"start_date"`
	EndDate       uint   `json:"end_date"`
	PromoAction   string `json:"promo_action" gorm:"type:enum('fixed','percent')"`
	Type          string `json:"type"`
	Status        string `json:"status"`
	Platform      string `json:"platform"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
}

type ApplyPromoCodeRequest struct {
	Code string `json:"code" validate:"required"`
}

type ApplyPromoCodeResponse struct {
	Valid          bool    `json:"valid"`
	Code           string  `json:"code,omitempty"`
	DiscountType   string  `json:"discount_type,omitempty"`
	DiscountAmount float64 `json:"discount_amount,omitempty"`
	Message        string  `json:"message"`
}
