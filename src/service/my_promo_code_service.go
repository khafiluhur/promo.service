package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

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
	PromoCodeRepo     *repository.PromoCodeRepositoryMySQL
	PromoCodeUsedRepo *repository.PromoCodeUsedRepositoryMySQL
}

func NewMyPromoCodeServiceV1(promoCodeRepo *repository.PromoCodeRepositoryMySQL, promoCodeUsedRepo *repository.PromoCodeUsedRepositoryMySQL) *MyPromoCodeServiceV1 {
	return &MyPromoCodeServiceV1{PromoCodeRepo: promoCodeRepo, PromoCodeUsedRepo: promoCodeUsedRepo}
}

func (m *MyPromoCodeServiceV1) MyList(ctx context.Context, userID string) ([]entity.PromoCode, error) {
	platformCfg, _ := utilspayment.GetPlatformConfig(ctx)

	wheres := []base.Where{
		{Name: "platform", Value: platformCfg.System},
		{Name: "is_active", Value: true},
	}

	list, _, err := m.PromoCodeRepo.List(ctx, 1, 100, []base.OrderBy{{Field: "created_at", Direction: "desc"}}, wheres)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (m *MyPromoCodeServiceV1) MyDetail(ctx context.Context, code string) (*entity.PromoCode, error) {
	platformCfg, _ := utilspayment.GetPlatformConfig(ctx)
	promos, err := m.PromoCodeRepo.ByPromoCodes(ctx, platformCfg.System, []string{code})
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
	promos, err := m.PromoCodeRepo.ByPromoCodes(ctx, platformCfg.System, []string{req.Code})
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

	switch promo.Rules {
	case "otg":
		if req.Amount < promo.Amount {
			return &dto.ApplyPromoCodeResponse{
				Valid:   false,
				Message: fmt.Sprintf("Minimal pembelian %.2f untuk menggunakan promo ini", promo.Amount),
			}, nil
		}

	case "sp":
		if req.ProductSlug == "" || req.ProductSlug != promo.ProductSlug {
			return &dto.ApplyPromoCodeResponse{
				Valid:   false,
				Message: "Promo ini hanya berlaku untuk produk tertentu",
			}, nil
		}

	case "pw":
		if req.PaymentMethod == "" || req.PaymentMethod != promo.PaymentMethod {
			return &dto.ApplyPromoCodeResponse{
				Valid:   false,
				Message: "Promo ini hanya berlaku untuk metode pembayaran tertentu",
			}, nil
		}
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
	promos, err := m.PromoCodeRepo.ByPromoCodes(ctx, platformCfg.System, []string{req.Code})
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

	switch promo.Rules {
	case "otg":
		if req.Amount < promo.Amount {
			return nil, fmt.Errorf("Minimal pembelian %.2f untuk menggunakan promo ini", promo.Amount)
		}
	case "sp":
		if req.ProductSlug == "" || req.ProductSlug != promo.ProductSlug {
			return nil, fmt.Errorf("Promo ini hanya berlaku untuk produk tertentu")
		}

	case "pw":
		if req.PaymentMethod == "" || req.PaymentMethod != promo.PaymentMethod {
			return nil, fmt.Errorf("Promo ini hanya berlaku untuk metode pembayaran tertentu")
		}
	}

	promo.Quantity = promo.Quantity - 1
	promo.UpdatedAt = uint64(time.Now().Unix())
	userID, err := strconv.ParseUint(req.UserID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}
	promo.UpdatedBy = userID

	_, err = m.PromoCodeRepo.Update(ctx, &promo, []string{"quantity", "updatedat", "updatedby"})
	if err != nil {
		if errDuplicate, ok := errorLib.ToErrDuplicateRecordsV2(err); ok {
			return nil, errDuplicate
		}
		return nil, err
	}

	used := entity.PromoCodeUsed{
		PromoCodeID:    uint(promo.Id),
		PromoCode:      promo.PromoCode,
		CustomerID:     req.UserID,
		OrderID:        &req.OrderID,
		Platform:       platformCfg.System,
		DiscountAmount: &promo.DiscountAmount,
		OrderTotal:     &req.Amount,
		Status:         "used",
		UsedAt:         uint64(time.Now().Unix()),
		CreatedAt:      uint64(time.Now().Unix()),
		UpdatedAt:      uint64(time.Now().Unix()),
	}
	_, err = m.PromoCodeUsedRepo.Create(ctx, &used)
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
