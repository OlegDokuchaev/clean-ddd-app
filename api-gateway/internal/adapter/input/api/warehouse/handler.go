package warehouse

import (
	commonRequest "api-gateway/internal/adapter/input/api/request"
	commonResponse "api-gateway/internal/adapter/input/api/response"
	request "api-gateway/internal/adapter/input/api/warehouse/request"
	response "api-gateway/internal/adapter/input/api/warehouse/response"
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

// IncreaseQuantity godoc
// @Summary Increase items quantity
// @Description Increase the quantity of items in the warehouse (admin only)
// @Tags items
// @Accept json
// @Produce json
// @Param request body warehouse_request.ReserveItemsRequest true "Items to increase quantity"
// @Success 200 "" "Quantity increased successfully"
// @Failure 400 {object} response.ErrorResponseDetail "Invalid request format"
// @Failure 401 {object} response.ErrorResponseDetail "Missing or invalid access token"
// @Failure 404 {object} response.ErrorResponseDetail "Item not found"
// @Failure 422 {object} response.ErrorResponseDetail "Invalid item data"
// @Failure 500 {object} response.ErrorResponseDetail "Server error"
// @Security AdminAccessToken
// @Router /items/increase [patch]
func (h *Handler) IncreaseQuantity(c *gin.Context) {
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

// DecreaseQuantity godoc
// @Summary Decrease items quantity
// @Description Decrease the quantity of items in the warehouse (admin only)
// @Tags items
// @Accept json
// @Produce json
// @Param request body warehouse_request.ReleaseItemsRequest true "Items to decrease quantity"
// @Success 200 "" "Quantity decreased successfully"
// @Failure 400 {object} response.ErrorResponseDetail "Invalid request format"
// @Failure 401 {object} response.ErrorResponseDetail "Missing or invalid access token"
// @Failure 404 {object} response.ErrorResponseDetail "Item not found"
// @Failure 422 {object} response.ErrorResponseDetail "Invalid item data or insufficient quantity"
// @Failure 500 {object} response.ErrorResponseDetail "Server error"
// @Security AdminAccessToken
// @Router /items/decrease [patch]
func (h *Handler) DecreaseQuantity(c *gin.Context) {
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

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product in the warehouse (admin only)
// @Tags products
// @Accept json
// @Produce json
// @Param request body warehouse_request.CreateProductRequest true "Product details"
// @Success 201 {object} warehouse_response.CreateProductResponse "Product created successfully"
// @Failure 400 {object} response.ErrorResponseDetail "Invalid request format"
// @Failure 401 {object} response.ErrorResponseDetail "Missing or invalid access token"
// @Failure 409 {object} response.ErrorResponseDetail "Product already exists"
// @Failure 422 {object} response.ErrorResponseDetail "Invalid product data"
// @Failure 500 {object} response.ErrorResponseDetail "Server error"
// @Security AdminAccessToken
// @Router /products [post]
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

// GetAllItems godoc
// @Summary Get all warehouse items
// @Description Get a list of all items available in the warehouse
// @Tags items
// @Accept json
// @Produce json
// @Success 200 {object} warehouse_response.ItemsResponse "List of items"
// @Failure 500 {object} response.ErrorResponseDetail "Server error"
// @Router /items [get]
func (h *Handler) GetAllItems(c *gin.Context) {
	items, err := h.uc.GetAllItems(c)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.ToItemsResponse(items))
}
