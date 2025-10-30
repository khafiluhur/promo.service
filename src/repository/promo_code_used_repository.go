package repository

import (
	"github.com/harryosmar/generic-gorm/base"
	entityCms "github.com/tripdeals/cms.backend.tripdeals.id/src/entity"
	"github.com/tripdeals/promo.service/src/entity"
	"gorm.io/gorm"
)

type PromoCodeUsedRepositoryMySQL struct {
	*base.BaseGorm[entity.PromoCodeUsed, int64]
	sectionItemPosition *base.BaseGorm[entityCms.SectionItemPosition, int64]
}

func NewPromoCodeUsedRepositoryMySQL(db *gorm.DB) *PromoCodeUsedRepositoryMySQL {
	return &PromoCodeUsedRepositoryMySQL{
		BaseGorm:            base.NewBaseGorm[entity.PromoCodeUsed, int64](db),
		sectionItemPosition: base.NewBaseGorm[entityCms.SectionItemPosition, int64](db),
	}
}
