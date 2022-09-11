package external

import (
	"encoding/json"
	"fetch-app/internal/config"
	"fetch-app/internal/model"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

type currencyService struct {
	config config.Config
}

func NewCurrencyService(config config.Config) *currencyService {

	return &currencyService{
		config: config,
	}
}

func (s *currencyService) GetExchangeRate(baseCode, targetCode string) (float64, error) {

	var conversionRate *model.ConversionRate

	s.config.Logger().Info(fmt.Sprintf("Getting exhange rate for %s to %s", baseCode, targetCode))

	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/pair/%s/%s", s.config.CurrencyServiceAPIKey(), baseCode, targetCode)
	resp, err := http.Get(url)
	if err != nil {
		s.config.Logger().Error("GetExchangeRate: Failed to fetch exchange rate", zap.Error(err))
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, conversionRate)
	if err != nil {
		s.config.Logger().Error("GetExchangeRate: Failed to parse exchange rate", zap.Error(err))
		return 0, err
	}

	return conversionRate.Rate, nil
}
