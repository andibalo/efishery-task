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

func (h *CommodityController) getAggregratedCommodities(c *gin.Context) {

	h.cfg.Logger().Info("getAllCommodities:  getting all commodities")

	c.JSON(http.StatusOK, "aggregated route")
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
