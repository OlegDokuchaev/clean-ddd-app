package warehouse

import (
	commonRequest "api-gateway/internal/adapter/input/api/request"
	commonResponse "api-gateway/internal/adapter/input/api/response"
	"api-gateway/internal/adapter/input/api/warehouse/request"
	"api-gateway/internal/adapter/input/api/warehouse/response"
	warehouseUseCase "api-gateway/internal/domain/usecases/warehouse"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Handler struct {
	uc warehouseUseCase.UseCase
}

func NewHandler(warehouseUseCase warehouseUseCase.UseCase) *Handler {
	return &Handler{uc: warehouseUseCase}
}

func (h *Handler) IncreaseQuantity(c *gin.Context) {
	var req request.ReserveItemsRequest
	if err := commonRequest.ParseInput(c, &req, binding.JSON); err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	token, err := commonRequest.ParseAccessToken(c)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	items := request.ToItemInfoDtoList(req.Items)
	err = h.uc.ReserveItems(c, items, token)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) DecreaseQuantity(c *gin.Context) {
	var req request.ReleaseItemsRequest
	if err := commonRequest.ParseInput(c, &req, binding.JSON); err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	token, err := commonRequest.ParseAccessToken(c)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	items := request.ToItemInfoDtoList(req.Items)
	err = h.uc.ReleaseItems(c, items, token)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) CreateProduct(c *gin.Context) {
	var req request.CreateProductRequest
	if err := commonRequest.ParseInput(c, &req, binding.JSON); err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	token, err := commonRequest.ParseAccessToken(c)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	data := request.ToCreateProductDto(&req)
	productID, err := h.uc.CreateProduct(c, data, token)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.CreateProductResponse{
		ProductID: productID,
	})
}

func (h *Handler) GetAllItems(c *gin.Context) {
	items, err := h.uc.GetAllItems(c)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.ToItemsResponse(items))
}
