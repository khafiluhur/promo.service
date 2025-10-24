package repository

import (
	"context"

	"github.com/harryosmar/generic-gorm/base"
	entityCms "github.com/tripdeals/cms.backend.tripdeals.id/src/entity"
	"github.com/tripdeals/promo.service/src/entity"
	"gorm.io/gorm"
)

type StrikeThroughtPriceRepositoryMySQL struct {
	*base.BaseGorm[entity.StrikeThroughtPrice, int64]
	sectionItemPosition *base.BaseGorm[entityCms.SectionItemPosition, int64]
}

func NewStrikeThroughtPriceRepositoryMySQL(db *gorm.DB) *StrikeThroughtPriceRepositoryMySQL {
	return &StrikeThroughtPriceRepositoryMySQL{
		BaseGorm:            base.NewBaseGorm[entity.StrikeThroughtPrice, int64](db),
		sectionItemPosition: base.NewBaseGorm[entityCms.SectionItemPosition, int64](db),
	}
}

func (repo *StrikeThroughtPriceRepositoryMySQL) ByIds(ctx context.Context, ids []string) ([]entity.StrikeThroughtPrice, error) {
	result := []entity.StrikeThroughtPrice{}
	err := repo.DB(ctx).Model(&entity.StrikeThroughtPrice{}).Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
