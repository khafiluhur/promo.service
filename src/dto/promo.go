package dto

type PromoCode struct {
	ID             uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	TermCondition  string  `json:"term_condition"`
	StartDate      uint    `json:"start_date"`
	EndDate        uint    `json:"end_date"`
	Banner         string  `json:"banner"`
	Rules          string  `json:"rules" gorm:"type:enum('otg','sp','pw')"`
	Amount         float64 `json:"amount"`
	DiscountAmount float64 `json:"discountAmount"`
	ProductSlug    string  `json:"productSlug"`
	PaymentMethod  float64 `json:"paymentMethod"`
	PromoAction    string  `json:"promoAction" gorm:"type:enum('fixed','percent')"`
	PromoType      string  `json:"promoType"`
	Type           string  `json:"type"`
	PromoCode      string  `json:"promoCode"`
	CustomerLimit  *int    `json:"customer_limit,omitempty"`
	NewCustomer    bool    `json:"newCustomer"`
	Quantity       int     `json:"quantity"`
	IsActive       string  `json:"is_active"`
	SpecialPromo   *int    `json:"special_promo,omitempty"`
	IsDisplay      bool    `json:"isDisplay"`
	Platform       string  `json:"platform"`
	CreatedAt      int64   `json:"created_at"`
	UpdatedAt      int64   `json:"updated_at"`
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
	Code          string `json:"code" validate:"required"`
	ProductId     int64  `json:"product_id"`
	Amount        int64  `json:"amount"`
	PaymentMethod string `json:"payment_method"`
}

type ApplyPromoCodeResponse struct {
	Valid          bool    `json:"valid"`
	Code           string  `json:"code,omitempty"`
	DiscountType   string  `json:"discount_type,omitempty"`
	DiscountAmount float64 `json:"discount_amount,omitempty"`
	Message        string  `json:"message"`
}

type RedeemPromoRequest struct {
	Code    string `json:"code" validate:"required"`
	OrderID string `json:"orderId" validate:"required"`
	UserID  string `json:"-"`
}

type RedeemPromoResponse struct {
	Code       string `json:"code"`
	Status     string `json:"status"`
	RedeemedAt string `json:"redeemedAt"`
}
