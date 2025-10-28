package service

import (
	"context"
	"time"

	ctx2 "github.com/Golden-Rama-Digital/library-core-go/context"
	"github.com/harryosmar/generic-gorm/base"
	errorSvc "github.com/tripdeals/cms.backend.tripdeals.id/src/error"
	errorLib "github.com/tripdeals/library-service.go/error"
	utilspayment "github.com/tripdeals/payment.service/src/utils"
	"github.com/tripdeals/promo.service/src/entity"
	"github.com/tripdeals/promo.service/src/repository"
)

//go:generate mockgen -destination=mocks/mock_PromoCodeService.go -package=mocks . PromoCodeService
type PromoCodeService interface {
	List(ctx context.Context, page int, pageSize int, orders []base.OrderBy, wheres []base.Where) ([]entity.PromoCode, *base.Paginator, error)
	Detail(ctx context.Context, id int64) (*entity.PromoCode, error)
	ByPromoCodes(ctx context.Context, promoCodes []string) ([]entity.PromoCode, error)
	Create(ctx context.Context, record *entity.PromoCode) (*entity.PromoCode, error)
	Update(ctx context.Context, record *entity.PromoCode) error
	Activate(ctx context.Context, id int64) error
	Deactivate(ctx context.Context, id int64) error
}

type PromoCodeServiceV1 struct {
	Repo *repository.PromoCodeRepositoryMySQL
}

func NewPromoCodeServiceV1(repo *repository.PromoCodeRepositoryMySQL) *PromoCodeServiceV1 {
	return &PromoCodeServiceV1{Repo: repo}
}

func (p *PromoCodeServiceV1) List(ctx context.Context, page int, pageSize int, orders []base.OrderBy, wheres []base.Where) ([]entity.PromoCode, *base.Paginator, error) {
	platformCfg, _ := utilspayment.GetPlatformConfig(ctx)
	wheres = append(wheres, base.Where{
		Name:  "platform",
		Value: platformCfg.System,
	})
	return p.Repo.List(ctx, page, pageSize, orders, wheres)
}

func (p *PromoCodeServiceV1) Detail(ctx context.Context, id int64) (*entity.PromoCode, error) {
	detail, err := p.Repo.Detail(ctx, id)
	if err != nil {
		return nil, err
	}
	if detail == nil {
		return nil, errorSvc.ErrRecordTourNotFound
	}
	return detail, nil
}

func (p *PromoCodeServiceV1) ByPromoCodes(ctx context.Context, promoCodes []string) ([]entity.PromoCode, error) {
	return p.Repo.ByPromoCodes(ctx, promoCodes)
}

func (p *PromoCodeServiceV1) Create(ctx context.Context, record *entity.PromoCode) (*entity.PromoCode, error) {
	record.MetaData = record.SetMetaCreate(ctx)
	create, err := p.Repo.Create(ctx, record)
	if err != nil {
		if errDuplicate, ok := errorLib.ToErrDuplicateRecordsV2(err); ok {
			return nil, errDuplicate
		}
	}
	return create, err
}

func (p *PromoCodeServiceV1) Update(ctx context.Context, record *entity.PromoCode) error {
	userId, _ := ctx2.GetUserIdInt64FromSession(ctx)
	record.UpdatedAt = uint64(time.Now().Unix())
	record.UpdatedBy = uint64(userId)
	_, err := p.Repo.Update(ctx, record, []string{
		"name", "description", "termCondition",
		"amount", "product_slug", "paymentMethod",
		"promoAction", "promoType", "type",
		"status", "isDisplay", "endDate",
		"updatedat", "updatedby",
	})
	if err != nil {
		if errDuplicate, ok := errorLib.ToErrDuplicateRecordsV2(err); ok {
			return errDuplicate
		}
	}
	return err
}

func (p *PromoCodeServiceV1) Activate(ctx context.Context, id int64) error {
	record := &entity.PromoCode{Id: id, Status: true}
	record.MetaData = record.SetMetaUpdate(ctx)
	_, err := p.Repo.Update(ctx, record, []string{"status"})
	return err
}

func (p *PromoCodeServiceV1) Deactivate(ctx context.Context, id int64) error {
	record := &entity.PromoCode{Id: id, Status: false}
	record.MetaData = record.SetMetaUpdate(ctx)
	_, err := p.Repo.Update(ctx, record, []string{"status"})
	return err
}
