package entity

import "github.com/tripdeals/cms.backend.tripdeals.id/src/entity"

type PromoCode struct {
	Id             int64   `gorm:"primaryKey;autoIncrement;column:id"`
	Name           string  `gorm:"column:name;type:varchar(255)"`
	Description    string  `gorm:"column:description;type:varchar(255)"`
	TermCondition  string  `gorm:"column:termCondition;type:varchar(255)"`
	StartDate      uint    `gorm:"column:startDate;type:int(10) unsigned"`
	EndDate        uint    `gorm:"column:endDate;type:int(10) unsigned"`
	Banner         string  `gorm:"column:banner;type:varchar(255)"`
	Rules          string  `gorm:"column:rules;type:enum('otg','sp','pw')"`
	Amount         float64 `gorm:"column:amount;type:decimal(13,2)"`
	DiscountAmount float64 `gorm:"column:discountAmount;type:decimal(13,2)"`
	ProductSlug    string  `gorm:"column:productSlug;type:varchar(255)"`
	PaymentMethod  string  `gorm:"column:paymentMethod;type:varchar(255)"`
	PromoAction    string  `gorm:"column:promoAction;type:enum('fixed','percent')"`
	PromoType      string  `gorm:"column:promoType;type:varchar(255)"`
	Type           string  `gorm:"column:type;type:varchar(255)"`
	PromoCode      string  `gorm:"column:promoCode;type:varchar(255);unique"`
	CustomerLimit  *int    `gorm:"column:customerLimit;type:int(10)"`
	NewCustomer    bool    `gorm:"column:newCustomer;type:tinyint(3) unsigned"`
	Quantity       int     `gorm:"column:quantity;type:int(10)"`
	IsActive       bool    `gorm:"column:is_active;type:boolean"`
	SpecialPromo   *int    `gorm:"column:specialPromo;type:int(10) unsigned"`
	IsDisplay      bool    `gorm:"column:isDisplay;type:tinyint(3) unsigned"`
	Platform       string  `gorm:"column:platform;type:varchar(50)"`
	entity.MetaData
}

func (p PromoCode) TableName() string {
	return "promo_code"
}

func (p PromoCode) PrimaryKey() string {
	return "id"
}

func (p PromoCode) GetAllowedOrderFields() []string {
	return []string{"id", "promoCode", "created_at", "updated_at"}
}

func (p PromoCode) GetAllowedWhereFields() []string {
	return []string{"id", "is_active", "promoCode", "created_by", "updated_by"}
}
