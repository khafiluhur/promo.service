package service

import (
	"context"
	"fmt"
	"time"

	ctx2 "github.com/Golden-Rama-Digital/library-core-go/context"
	"github.com/harryosmar/generic-gorm/base"
	errorSvc "github.com/tripdeals/cms.backend.tripdeals.id/src/error"
	errorLib "github.com/tripdeals/library-service.go/error"
	utilspayment "github.com/tripdeals/payment.service/src/utils"
	"github.com/tripdeals/promo.service/src/dto"
	"github.com/tripdeals/promo.service/src/entity"
	"github.com/tripdeals/promo.service/src/repository"
)

//go:generate mockgen -destination=mocks/mock_MyPromoCodeService.go -package=mocks . MyPromoCodeService
type MyPromoCodeService interface {
	MyList(ctx context.Context, userID string) ([]entity.PromoCode, error)
	MyDetail(ctx context.Context, code string) (*entity.PromoCode, error)
	Apply(ctx context.Context, req dto.ApplyPromoCodeRequest) (*dto.ApplyPromoCodeResponse, error)
	Redeem(ctx context.Context, req dto.RedeemPromoRequest) (*dto.RedeemPromoResponse, error)
}

type MyPromoCodeServiceV1 struct {
	Repo *repository.PromoCodeRepositoryMySQL
}

func NewMyPromoCodeServiceV1(repo *repository.PromoCodeRepositoryMySQL) *MyPromoCodeServiceV1 {
	return &MyPromoCodeServiceV1{Repo: repo}
}

func (m *MyPromoCodeServiceV1) MyList(ctx context.Context, userID string) ([]entity.PromoCode, error) {
	platformCfg, _ := utilspayment.GetPlatformConfig(ctx)

	wheres := []base.Where{
		{Name: "platform", Value: platformCfg.System},
		{Name: "is_active", Value: true},
	}

	list, _, err := m.Repo.List(ctx, 1, 100, []base.OrderBy{{Field: "created_at", Direction: "desc"}}, wheres)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (m *MyPromoCodeServiceV1) MyDetail(ctx context.Context, code string) (*entity.PromoCode, error) {
	platformCfg, _ := utilspayment.GetPlatformConfig(ctx)
	promos, err := m.Repo.ByPromoCodes(ctx, platformCfg.System, []string{code})
	if err != nil {
		return nil, err
	}
	if len(promos) == 0 {
		return nil, errorSvc.ErrRecordTourNotFound
	}
	return &promos[0], nil
}

func (m *MyPromoCodeServiceV1) Apply(ctx context.Context, req dto.ApplyPromoCodeRequest) (*dto.ApplyPromoCodeResponse, error) {
	platformCfg, _ := utilspayment.GetPlatformConfig(ctx)
	promos, err := m.Repo.ByPromoCodes(ctx, platformCfg.System, []string{req.Code})
	if err != nil || len(promos) == 0 {
		return &dto.ApplyPromoCodeResponse{Valid: false, Message: "Kode promo tidak ditemukan"}, nil
	}

	promo := promos[0]
	now := time.Now()
	start := time.Unix(int64(promo.StartDate), 0)
	end := time.Unix(int64(promo.EndDate), 0)

	if now.Before(start) || now.After(end) {
		return &dto.ApplyPromoCodeResponse{Valid: false, Message: "Kode promo sudah tidak berlaku"}, nil
	}

	if !promo.IsActive {
		return &dto.ApplyPromoCodeResponse{Valid: false, Message: "Kode promo tidak aktif"}, nil
	}

	if promo.Quantity <= 0 {
		return &dto.ApplyPromoCodeResponse{Valid: false, Message: "Kode promo sudah habis kuotanya"}, nil
	}

	return &dto.ApplyPromoCodeResponse{
		Valid:          true,
		Code:           promo.PromoCode,
		DiscountType:   promo.PromoAction,
		DiscountAmount: promo.DiscountAmount,
		Message:        fmt.Sprintf("Promo %s berhasil diterapkan", promo.PromoCode),
	}, nil
}

func (m *MyPromoCodeServiceV1) Redeem(ctx context.Context, req dto.RedeemPromoRequest) (*dto.RedeemPromoResponse, error) {
	platformCfg, _ := utilspayment.GetPlatformConfig(ctx)
	promos, err := m.Repo.ByPromoCodes(ctx, platformCfg.System, []string{req.Code})
	if err != nil || len(promos) == 0 {
		return nil, errorSvc.ErrRecordTourNotFound
	}

	promo := promos[0]
	now := time.Now()
	start := time.Unix(int64(promo.StartDate), 0)
	end := time.Unix(int64(promo.EndDate), 0)

	if now.Before(start) || now.After(end) {
		return nil, fmt.Errorf("Kode promo sudah tidak berlaku")
	}

	if !promo.IsActive {
		return nil, fmt.Errorf("Kode promo tidak aktif")
	}

	if promo.Quantity <= 0 {
		return nil, fmt.Errorf("Kode promo sudah habis kuotanya")
	}

	promo.Quantity = promo.Quantity - 1
	userId, _ := ctx2.GetUserIdInt64FromSession(ctx)
	promo.UpdatedAt = uint64(time.Now().Unix())
	promo.UpdatedBy = uint64(userId)

	_, err = m.Repo.Update(ctx, &promo, []string{"quantity", "updatedat", "updatedby"})
	if err != nil {
		if errDuplicate, ok := errorLib.ToErrDuplicateRecordsV2(err); ok {
			return nil, errDuplicate
		}
		return nil, err
	}

	return &dto.RedeemPromoResponse{
		Code:       promo.PromoCode,
		Status:     "redeemed",
		RedeemedAt: time.Now().Format(time.RFC3339),
	}, nil
}
