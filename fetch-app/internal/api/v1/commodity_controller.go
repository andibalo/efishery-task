package v1

import (
	voerrors "fetch-app/internal/autherrors"
	"fetch-app/internal/config"
	"fetch-app/internal/constants"
	"fetch-app/internal/response"
	"fetch-app/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
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

	cr.GET("/", h.getAllCommodities)
}

func (h *CommodityController) getAllCommodities(c *gin.Context) {

	code, commodities, err := h.commodityService.GetAllCommodities()

	if err != nil {
		h.cfg.Logger().Error("getAllCommodities: error getting commodities", zap.Error(err))

		h.failedCommodityResponse(c, code, err, "error fetching commodities")

		return
	}

	resp := response.NewResponse(response.Success, commodities)

	resp.SetResponseMessage("Successfully get commodities")

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
