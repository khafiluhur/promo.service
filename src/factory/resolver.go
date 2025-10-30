package factory

import (
	"context"
	"fmt"
	"time"

	loggerCore "github.com/Golden-Rama-Digital/library-core-go/logger"
	cache "github.com/harryosmar/cache-go"
	"github.com/redis/go-redis/v9"
	paymentAuthenticator "github.com/tripdeals/library-service.go"
	"github.com/tripdeals/promo.service/config"
	"github.com/tripdeals/promo.service/src/repository"
	"github.com/tripdeals/promo.service/src/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ResolveRedisClient(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDbIndex,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return client
}

func ResolveCache(client *redis.Client) cache.CacheRepo {
	return cache.NewRedisCacheV2(client)
}

func MakeMySQLDBFromConfig(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBMySQLDbUsername, cfg.DBMySQLPassword, cfg.DBMySQLHost, cfg.DBMySQLPort, cfg.DBMySQLDbName,
	)

	dbLogger := loggerCore.MakeGormSingleLineLoggerWithDefaultInterface()

	db, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{
			Logger: dbLogger,
		},
	)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(cfg.DBMySQLMaxIdleConnection)
	sqlDB.SetMaxOpenConns(cfg.DBMySQLMaxOpenConnection)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.DBMySQLConnMaxLifetime) * time.Minute)

	return db, nil
}

func ResolveDatabase(cfg *config.Config) (*gorm.DB, error) {
	return MakeMySQLDBFromConfig(cfg)
}

func ResolvePromoCodeService(db *gorm.DB) service.PromoCodeServiceV1 {
	return *service.NewPromoCodeServiceV1(repository.NewPromoCodeRepositoryMySQL(db))
}

func ResolveMyPromoCodeService(db *gorm.DB) service.MyPromoCodeServiceV1 {
	return *service.NewMyPromoCodeServiceV1(repository.NewPromoCodeRepositoryMySQL(db), repository.NewPromoCodeUsedRepositoryMySQL(db))
}

func ResolveStrikeThroughtPriceService(db *gorm.DB) service.StrikeThroughtPriceServiceV1 {
	return *service.NewStrikeThroughtPriceServiceV1(repository.NewStrikeThroughtPriceRepositoryMySQL(db))
}

func ResolvePaymentAuthenticator(cfg *config.Config) paymentAuthenticator.PaymentInternalTokenAuthenticatorV1 {
	return *paymentAuthenticator.NewPaymentInternalTokenAuthenticator(
		paymentAuthenticator.NewDynamicSecretAuthenticatorHmac512(),
		time.Duration(cfg.ExpiredService)*time.Second,
	)
}
