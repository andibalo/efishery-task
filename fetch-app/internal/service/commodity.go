package service

import (
	"encoding/json"
	"fetch-app/internal/model"
	"fetch-app/internal/response"
	"net/http"
)

type commodityService struct {
	config Config
}

func NewCommodityService(config Config) *commodityService {

	return &commodityService{
		config: config,
	}
}

func (s *commodityService) GetAllCommodities() (code response.Code, commodities []model.Commodity, err error) {
	resp, err := http.Get("https://stein.efishery.com/v1/storages/5e1edf521073e315924ceab4/list")
	if err != nil {

		s.config.Logger().Error("GetAllCommodities: Failed to fetch from efishery endpoint")
		return response.ServerError, nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&commodities)
	if err != nil {
		s.config.Logger().Error("GetAllCommodities: Failed to decode response from efishery endpoint")
		return response.ServerError, nil, err
	}

	return response.Success, commodities, err
}
