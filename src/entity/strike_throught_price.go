package entity

import "github.com/tripdeals/cms.backend.tripdeals.id/src/entity"

type StrikeThroughtPrice struct {
	Id            int64  `gorm:"primaryKey;autoIncrement;column:id"`
	Description   string `gorm:"column:description;type:varchar(255)"`
	ProductSlug   string `gorm:"column:product_slug;type:varchar(255)"`
	Departure     *int   `gorm:"column:departure;type:int"`
	TermCondition string `gorm:"column:termCondition;type:varchar(255)"`
	StartDate     uint   `gorm:"column:startDate;type:int(10) unsigned"`
	EndDate       uint   `gorm:"column:endDate;type:int(10) unsigned"`
	PromoAction   string `gorm:"column:promoAction;type:enum('fixed','percent')"`
	Type          string `gorm:"column:type;type:varchar(255)"`
	Status        bool   `gorm:"column:status;type:boolean"`
	Platform      string `gorm:"column:platform;type:varchar(50)"`
	entity.MetaData
}

func (s StrikeThroughtPrice) TableName() string {
	return "strike_throught_price"
}

func (s StrikeThroughtPrice) PrimaryKey() string {
	return "id"
}

func (s StrikeThroughtPrice) GetAllowedOrderFields() []string {
	return []string{"id", "product_slug", "created_at", "updated_at"}
}

func (s StrikeThroughtPrice) GetAllowedWhereFields() []string {
	return []string{"id", "product_slug", "created_by", "updated_by"}
}
