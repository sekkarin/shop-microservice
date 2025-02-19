package inventoryHandler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sekkarin/shop-microservice/config"
	"github.com/sekkarin/shop-microservice/modules/inventory"
	"github.com/sekkarin/shop-microservice/modules/inventory/inventoryUsecase"
	"github.com/sekkarin/shop-microservice/pkg/request"
	"github.com/sekkarin/shop-microservice/pkg/response"
)

type (
	InventoryHttpHandlerService interface {
		FindPlayerItems(c echo.Context) error
	}

	inventoryHttpHandler struct {
		cfg              *config.Config
		inventoryUsecase inventoryUsecase.InventoryUsecaseService
	}
)

func NewInventoryHttpHandler(cfg *config.Config, inventoryUsecase inventoryUsecase.InventoryUsecaseService) InventoryHttpHandlerService {
	return &inventoryHttpHandler{
		cfg:              cfg,
		inventoryUsecase: inventoryUsecase,
	}
}

func (h *inventoryHttpHandler) FindPlayerItems(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(inventory.InventorySearchReq)
	playerId := c.Param("player_id")

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.inventoryUsecase.FindPlayerItems(ctx, h.cfg, playerId, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}
