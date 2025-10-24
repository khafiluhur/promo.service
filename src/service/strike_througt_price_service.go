package service

import (
	"context"
	"time"

	ctx2 "github.com/Golden-Rama-Digital/library-core-go/context"
	"github.com/harryosmar/generic-gorm/base"
	errorSvc "github.com/tripdeals/cms.backend.tripdeals.id/src/error"
	errorLib "github.com/tripdeals/library-service.go/error"
	"github.com/tripdeals/promo.service/src/entity"
	"github.com/tripdeals/promo.service/src/repository"
)

//go:generate mockgen -destination=mocks/mock_StrikeThroughtPriceService.go -package=mocks . StrikeThroughtPriceService
type StrikeThroughtPriceService interface {
	List(ctx context.Context, page int, pageSize int, orders []base.OrderBy, wheres []base.Where) ([]entity.StrikeThroughtPrice, *base.Paginator, error)
	Detail(ctx context.Context, id int64) (*entity.StrikeThroughtPrice, error)
	BySlugs(ctx context.Context, slugs []string) ([]entity.StrikeThroughtPrice, error)
	Create(ctx context.Context, record *entity.StrikeThroughtPrice) (*entity.StrikeThroughtPrice, error)
	Update(ctx context.Context, record *entity.StrikeThroughtPrice) error
	Activate(ctx context.Context, id int64) error
	Deactivate(ctx context.Context, id int64) error
}

type StrikeThroughtPriceServiceV1 struct {
	Repo *repository.StrikeThroughtPriceRepositoryMySQL
}

func NewStrikeThroughtPriceServiceV1(repo *repository.StrikeThroughtPriceRepositoryMySQL) *StrikeThroughtPriceServiceV1 {
	return &StrikeThroughtPriceServiceV1{Repo: repo}
}

func (s *StrikeThroughtPriceServiceV1) List(ctx context.Context, page int, pageSize int, orders []base.OrderBy, wheres []base.Where) ([]entity.StrikeThroughtPrice, *base.Paginator, error) {
	return s.Repo.List(ctx, page, pageSize, orders, wheres)
}

func (s *StrikeThroughtPriceServiceV1) Detail(ctx context.Context, id int64) (*entity.StrikeThroughtPrice, error) {
	detail, err := s.Repo.Detail(ctx, id)
	if err != nil {
		return nil, err
	}
	if detail == nil {
		return nil, errorSvc.ErrRecordTourNotFound
	}
	return detail, nil
}

func (s *StrikeThroughtPriceServiceV1) BySlugs(ctx context.Context, slugs []string) ([]entity.StrikeThroughtPrice, error) {
	return s.Repo.ByIds(ctx, slugs)
}

func (s *StrikeThroughtPriceServiceV1) Create(ctx context.Context, record *entity.StrikeThroughtPrice) (*entity.StrikeThroughtPrice, error) {
	record.MetaData = record.SetMetaCreate(ctx)
	create, err := s.Repo.Create(ctx, record)
	if err != nil {
		if errDuplicate, ok := errorLib.ToErrDuplicateRecordsV2(err); ok {
			return nil, errDuplicate
		}
	}
	return create, err
}

func (s *StrikeThroughtPriceServiceV1) Update(ctx context.Context, record *entity.StrikeThroughtPrice) error {
	userId, _ := ctx2.GetUserIdInt64FromSession(ctx)
	record.UpdatedAt = uint64(time.Now().Unix())
	record.UpdatedBy = uint64(userId)
	_, err := s.Repo.Update(ctx, record, []string{
		"description", "product_slug", "departure",
		"termCondition", "startDate", "endDate",
		"promoAction", "type", "status",
		"platform", "updatedat", "updatedby",
	})
	if err != nil {
		if errDuplicate, ok := errorLib.ToErrDuplicateRecordsV2(err); ok {
			return errDuplicate
		}
	}
	return err
}

func (s *StrikeThroughtPriceServiceV1) Activate(ctx context.Context, id int64) error {
	record := &entity.StrikeThroughtPrice{Id: id, Status: true}
	record.MetaData = record.SetMetaUpdate(ctx)
	_, err := s.Repo.Update(ctx, record, []string{"is_active"})
	return err
}

func (s *StrikeThroughtPriceServiceV1) Deactivate(ctx context.Context, id int64) error {
	record := &entity.StrikeThroughtPrice{Id: id, Status: false}
	record.MetaData = record.SetMetaUpdate(ctx)
	_, err := s.Repo.Update(ctx, record, []string{"is_active"})
	return err
}
