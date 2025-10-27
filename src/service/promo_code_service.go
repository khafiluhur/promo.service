package service

import (
	"context"
	"fmt"
	"time"

	ctx2 "github.com/Golden-Rama-Digital/library-core-go/context"
	"github.com/harryosmar/generic-gorm/base"
	errorSvc "github.com/tripdeals/cms.backend.tripdeals.id/src/error"
	errorLib "github.com/tripdeals/library-service.go/error"
	"github.com/tripdeals/promo.service/src/dto"
	"github.com/tripdeals/promo.service/src/entity"
	"github.com/tripdeals/promo.service/src/repository"
)

//go:generate mockgen -destination=mocks/mock_PromoCodeService.go -package=mocks . PromoCodeService
type PromoCodeService interface {
	List(ctx context.Context, page int, pageSize int, orders []base.OrderBy, wheres []base.Where) ([]entity.PromoCode, *base.Paginator, error)
	Detail(ctx context.Context, id int64) (*entity.PromoCode, error)
	BySlugs(ctx context.Context, promoCodes []string) ([]entity.PromoCode, error)
	Create(ctx context.Context, record *entity.PromoCode) (*entity.PromoCode, error)
	Update(ctx context.Context, record *entity.PromoCode) error
	Activate(ctx context.Context, id int64) error
	Deactivate(ctx context.Context, id int64) error
	Apply(ctx context.Context, req dto.ApplyPromoCodeRequest) (*dto.ApplyPromoCodeResponse, error)
}

type PromoCodeServiceV1 struct {
	Repo *repository.PromoCodeRepositoryMySQL
}

func NewPromoCodeServiceV1(repo *repository.PromoCodeRepositoryMySQL) *PromoCodeServiceV1 {
	return &PromoCodeServiceV1{Repo: repo}
}

func (s *PromoCodeServiceV1) List(ctx context.Context, page int, pageSize int, orders []base.OrderBy, wheres []base.Where) ([]entity.PromoCode, *base.Paginator, error) {
	return s.Repo.List(ctx, page, pageSize, orders, wheres)
}

func (s *PromoCodeServiceV1) Detail(ctx context.Context, id int64) (*entity.PromoCode, error) {
	detail, err := s.Repo.Detail(ctx, id)
	if err != nil {
		return nil, err
	}
	if detail == nil {
		return nil, errorSvc.ErrRecordTourNotFound
	}
	return detail, nil
}

func (s *PromoCodeServiceV1) BySlugs(ctx context.Context, promoCodes []string) ([]entity.PromoCode, error) {
	return s.Repo.ByPromoCodes(ctx, promoCodes)
}

func (s *PromoCodeServiceV1) Create(ctx context.Context, record *entity.PromoCode) (*entity.PromoCode, error) {
	record.MetaData = record.SetMetaCreate(ctx)
	create, err := s.Repo.Create(ctx, record)
	if err != nil {
		if errDuplicate, ok := errorLib.ToErrDuplicateRecordsV2(err); ok {
			return nil, errDuplicate
		}
	}
	return create, err
}

func (s *PromoCodeServiceV1) Update(ctx context.Context, record *entity.PromoCode) error {
	userId, _ := ctx2.GetUserIdInt64FromSession(ctx)
	record.UpdatedAt = uint64(time.Now().Unix())
	record.UpdatedBy = uint64(userId)
	_, err := s.Repo.Update(ctx, record, []string{
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

func (s *PromoCodeServiceV1) Activate(ctx context.Context, id int64) error {
	record := &entity.PromoCode{Id: id, Status: true}
	record.MetaData = record.SetMetaUpdate(ctx)
	_, err := s.Repo.Update(ctx, record, []string{"status"})
	return err
}

func (s *PromoCodeServiceV1) Deactivate(ctx context.Context, id int64) error {
	record := &entity.PromoCode{Id: id, Status: false}
	record.MetaData = record.SetMetaUpdate(ctx)
	_, err := s.Repo.Update(ctx, record, []string{"status"})
	return err
}

func (s *PromoCodeServiceV1) Apply(ctx context.Context, req dto.ApplyPromoCodeRequest) (*dto.ApplyPromoCodeResponse, error) {
	promos, err := s.Repo.ByPromoCodes(ctx, []string{req.Code})
	if err != nil {
		return &dto.ApplyPromoCodeResponse{Valid: false, Message: "Kode promo tidak ditemukan"}, nil
	}

	promo := promos[0]
	now := time.Now()
	start := time.Unix(int64(promo.StartDate), 0)
	end := time.Unix(int64(promo.EndDate), 0)
	if now.Before(start) || now.After(end) {
		return &dto.ApplyPromoCodeResponse{Valid: false, Message: "Kode promo sudah tidak berlaku"}, nil
	}

	if !promo.Status {
		return &dto.ApplyPromoCodeResponse{Valid: false, Message: "Kode promo tidak aktif"}, nil
	}

	if promo.Quantity > 0 {
		return &dto.ApplyPromoCodeResponse{Valid: false, Message: "Kode promo sudah habis kuotanya"}, nil
	}

	return &dto.ApplyPromoCodeResponse{
		Valid:          true,
		Code:           promo.PromoCode,
		DiscountType:   promo.PromoAction,
		DiscountAmount: promo.Amount,
		Message:        fmt.Sprintf("Promo %s berhasil diterapkan", promo.PromoCode),
	}, nil
}
