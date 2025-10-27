package repository

import (
	"context"

	"github.com/harryosmar/generic-gorm/base"
	entityCms "github.com/tripdeals/cms.backend.tripdeals.id/src/entity"
	"github.com/tripdeals/promo.service/src/entity"
	"gorm.io/gorm"
)

type PromoCodeRepositoryMySQL struct {
	*base.BaseGorm[entity.PromoCode, int64]
	sectionItemPosition *base.BaseGorm[entityCms.SectionItemPosition, int64]
}

func NewPromoCodeRepositoryMySQL(db *gorm.DB) *PromoCodeRepositoryMySQL {
	return &PromoCodeRepositoryMySQL{
		BaseGorm:            base.NewBaseGorm[entity.PromoCode, int64](db),
		sectionItemPosition: base.NewBaseGorm[entityCms.SectionItemPosition, int64](db),
	}
}

func (repo *PromoCodeRepositoryMySQL) ByPromoCodes(ctx context.Context, promoCodes []string) ([]entity.PromoCode, error) {
	result := []entity.PromoCode{}
	err := repo.DB(ctx).Model(&entity.PromoCode{}).Where("promoCode IN ?", promoCodes).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
