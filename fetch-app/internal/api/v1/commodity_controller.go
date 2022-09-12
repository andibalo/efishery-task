package v1

import (
	"errors"
	voerrors "fetch-app/internal/autherrors"
	"fetch-app/internal/config"
	"fetch-app/internal/constants"
	middleware "fetch-app/internal/middlewares"
	"fetch-app/internal/model"
	"fetch-app/internal/response"
	"fetch-app/internal/service"
	"fetch-app/internal/util"
	"fmt"
	"github.com/bluele/gcache"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type CommodityController struct {
	cfg              config.Config
	commodityService service.CommodityService
	currencyService  service.CurrencyService
}

func NewCommodityController(cfg config.Config, commodityService service.CommodityService, currencyService service.CurrencyService) *CommodityController {
	return &CommodityController{
		cfg:              cfg,
		commodityService: commodityService,
		currencyService:  currencyService,
	}
}

func (h *CommodityController) AddRoutes(r *gin.Engine) {
	cr := r.Group(constants.V1BasePath + constants.CommodityBasePath)

	cr.GET("/", middleware.CheckJWTToken(), h.getAllCommodities)
	cr.GET("/aggregated", middleware.CheckJWTTokenAdmin(), h.getAggregratedCommodities)
}

func (h *CommodityController) getAllCommodities(c *gin.Context) {

	h.cfg.Logger().Info("getAllCommodities:  getting all commodities")

	code, commodities, err := h.commodityService.GetAllCommodities()

	if err != nil {
		h.cfg.Logger().Error("getAllCommodities: error getting commodities", zap.Error(err))

		h.failedCommodityResponse(c, code, err, "error fetching commodities")

		return
	}

	var conversionRate float64

	cachedRate, err := h.cfg.GCache().Get(constants.IDRUSDConversionRate)

	if err != nil {
		if errors.Is(err, gcache.KeyNotFoundError) {
			h.cfg.Logger().Info("getAllCommodities: idr usd conversion rate cache not set")
		} else {
			h.cfg.Logger().Error("getAllCommodities: error getting from cache", zap.Error(err))
		}
	}

	if cachedRate == nil {
		code, conversionRate, err = h.currencyService.GetExchangeRate(constants.IndonesiaCurrencyCode, constants.USCurrencyCode)

		if err != nil {
			h.cfg.Logger().Error("getAllCommodities: error getting currency conversion rate", zap.Error(err))

			h.failedCommodityResponse(c, code, err, "error getting exchange rate")

			return
		}

		err = h.cfg.GCache().SetWithExpire(constants.IDRUSDConversionRate, conversionRate, constants.IDRUSDConversionRateTTL)

		if err != nil {
			h.cfg.Logger().Error("getAllCommodities: failed setting conversion rate to cache", zap.Error(err))
		}

	} else {
		conversionRate, err = strconv.ParseFloat(fmt.Sprintf("%v", cachedRate), 64)
	}

	transformedCommodities := h.addUSDPriceToList(commodities, conversionRate)

	resp := response.NewResponse(response.Success, transformedCommodities)

	resp.SetResponseMessage("Successfully get commodities")

	c.JSON(http.StatusOK, resp)
}

type tempData struct {
	provinsi string
	amount   int
	minggu   string
}

func (h *CommodityController) getAggregratedCommodities(c *gin.Context) {

	h.cfg.Logger().Info("getAggregratedCommodities: getting aggregrated commodities")

	code, commodities, err := h.commodityService.GetAllCommodities()

	if err != nil {
		h.cfg.Logger().Error("getAllCommodities: error getting commodities", zap.Error(err))

		h.failedCommodityResponse(c, code, err, "error fetching commodities")

		return
	}

	var data []tempData

	for _, val := range commodities {
		if val.Tanggal == "" || val.Price == "" || val.Size == "" || val.Provinsi == "" {
			continue
		}

		dateTime := util.ParseStringToTime(val.Tanggal)

		_, week := dateTime.ISOWeek()

		price, _ := strconv.Atoi(val.Price)
		size, _ := strconv.Atoi(val.Size)
		amount := price * size

		temp := tempData{
			provinsi: val.Provinsi,
			amount:   amount,
			minggu:   fmt.Sprintf("week_%s", strconv.Itoa(week)),
		}

		data = append(data, temp)
	}

	tempMap := make(map[string]map[string]map[int]int)

	for _, val := range data {

		if prov, ok := tempMap[val.provinsi]; !ok {
			priceMap := map[int]int{val.amount: val.amount}
			minggu := map[string]map[int]int{val.minggu: priceMap}
			tempMap[val.provinsi] = minggu
		} else {

			if week, ok := prov[val.minggu]; !ok {
				priceMap := map[int]int{val.amount: val.amount}
				prov[val.minggu] = priceMap
			} else {

				if _, ok := week[val.amount]; !ok {
					week[val.amount] = val.amount
				}
			}
		}

	}

	var result []model.AggregratedCommodity

	for key, val := range tempMap {

		data := model.AggregratedCommodity{
			Provinsi: key,
			Profit:   val,
			Max:      findMaxProfit(val),
			Min:      findMinProfit(val),
			Avg:      findAvgProfit(val),
			Median:   findMedianProfit(val),
		}

		result = append(result, data)
	}

	resp := response.NewResponse(response.Success, result)

	resp.SetResponseMessage("Successfully get aggregrated commodities")

	c.JSON(http.StatusOK, resp)
}

func (h *CommodityController) failedCommodityResponse(c *gin.Context, code response.Code, err error, errorMsg string) {
	if code == "" {
		code = voerrors.MapErrorsToCode(err)
	}

	resp := response.Wrapper{
		ResponseCode: code,
		Status:       code.GetStatus(),
		Message:      code.GetMessage(),
	}

	if errorMsg != "" {
		resp.SetResponseMessage(errorMsg)
	}

	c.JSON(voerrors.MapErrorsToStatusCode(err), resp)
}

func (h *CommodityController) addUSDPriceToList(listData []model.Commodity, conversionRate float64) (transformedCommodity []model.Commodity) {

	for _, val := range listData {
		if val.Price == "" {
			val.USDPrice = ""

			transformedCommodity = append(transformedCommodity, val)
			continue
		}
		priceFloat, err := strconv.ParseFloat(val.Price, 64)
		if err != nil {
			h.cfg.Logger().Error(fmt.Sprintf("Failed to convert price to float for id %v", val.ID), zap.Error(err))
			val.USDPrice = "error converting"

			transformedCommodity = append(transformedCommodity, val)
			continue
		}

		usdPrice := priceFloat * conversionRate

		val.USDPrice = fmt.Sprintf("%f", usdPrice)
		transformedCommodity = append(transformedCommodity, val)
	}

	return transformedCommodity
}

func findMaxProfit(data map[string]map[int]int) float64 {
	var max int
	for _, val := range data {
		for _, amount := range val {
			if amount >= max {
				max = amount
			}
		}
	}

	return float64(max)
}

func findMinProfit(data map[string]map[int]int) float64 {

	min := int(^uint(0) >> 1)
	for _, val := range data {
		for _, amount := range val {
			if amount <= min {
				min = amount
			}
		}
	}

	return float64(min)
}

func findAvgProfit(data map[string]map[int]int) float64 {
	var sum, counter int
	for _, val := range data {
		for _, amount := range val {
			sum += amount
			counter++
		}
	}

	return float64(sum / counter)
}

func findMedianProfit(data map[string]map[int]int) float64 {
	var arr []int
	for _, val := range data {
		for _, amount := range val {
			arr = append(arr, amount)
		}
	}

	counter := len(arr)

	if counter+1%2 == 0 {
		a := arr[(counter / 2)]
		b := arr[(counter/2)+1]
		return float64((a + b) / 2)
	} else {
		return float64(arr[counter/2])
	}
}
