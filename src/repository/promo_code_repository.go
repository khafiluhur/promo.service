package repository

import (
	"context"
	"fmt"

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

func (repo *PromoCodeRepositoryMySQL) ByPromoCodes(ctx context.Context, platform string, promoCodes []string) ([]entity.PromoCode, error) {
	result := []entity.PromoCode{}
	wheres := []base.Where{
		{Name: "platform", Value: platform},
		{Name: "is_active", Value: true},
	}

	query := repo.DB(ctx).Model(&entity.PromoCode{})

	for _, w := range wheres {
		query = query.Where(fmt.Sprintf("%s = ?", w.Name), w.Value)
	}

	query = query.Where("promoCode IN ?", promoCodes)

	if err := query.Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
