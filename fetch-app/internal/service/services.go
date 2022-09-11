package service

import (
	"fetch-app/internal/model"
	"fetch-app/internal/response"
	"go.uber.org/zap"
)

type CommodityService interface {
	GetAllCommodities() (code response.Code, commodities []model.Commodity, err error)
}

type CurrencyService interface {
	GetExchangeRate(baseCode, targetCode string) (float64, error)
}

type Config interface {
	Logger() *zap.Logger
}
