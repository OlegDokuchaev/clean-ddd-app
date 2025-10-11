package warehouse

import (
	commonRequest "api-gateway/internal/adapter/input/api/request"
	commonResponse "api-gateway/internal/adapter/input/api/response"
	request "api-gateway/internal/adapter/input/api/warehouse/request"
	response "api-gateway/internal/adapter/input/api/warehouse/response"
	warehouseUseCase "api-gateway/internal/domain/usecases/warehouse"
	"io"
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
// @Param request query warehouse_request.GetAllItemsRequest true "Pagination"
// @Success 200 {object} warehouse_response.ItemsResponse "List of items"
// @Failure 500 {object} response.ErrorResponseDetail "Server error"
// @Router /items [get]
func (h *Handler) GetAllItems(c *gin.Context) {
	// Limit
	limit, err := commonRequest.ParseParamInt(c, "limit")
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	// Offset
	offset, err := commonRequest.ParseParamInt(c, "offset")
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	items, err := h.uc.GetAllItems(c, limit, offset)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.ToItemsResponse(items))
}

// UpdateProductImage godoc
// @Summary      Update product image
// @Description  Update the image for a specific product (admin only)
// @Tags         products
// @Accept       mpfd
// @Produce      json
// @Param        id   path      string  true  "Product ID" format(uuid)
// @Param        file formData  file    true  "Product image file"
// @Success      204  "No Content" "Image updated successfully"
// @Failure      400  {object}  response.ErrorResponseDetail "Invalid request format or invalid product ID"
// @Failure      401  {object}  response.ErrorResponseDetail "Missing or invalid access token"
// @Failure      404  {object}  response.ErrorResponseDetail "Product not found"
// @Failure      500  {object}  response.ErrorResponseDetail "Server error"
// @Security     AdminAccessToken
// @Router       /products/{id}/image [put]
func (h *Handler) UpdateProductImage(c *gin.Context) {
	// Product ID
	productID, err := commonRequest.ParseParamUUID(c, "id")
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	// Access token
	token, err := commonRequest.ParseAccessToken(c)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	// File
	fileHeader, err := commonRequest.ParseFormFile(c, "file")
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}
	defer file.Close()

	// Content type
	contentType := fileHeader.Header.Get("Content-Type")

	err = h.uc.UpdateProductImage(c, productID, file, contentType, token)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetProductImage godoc
// @Summary      Get product image
// @Description  Get the image for a specific product (admin only)
// @Tags         products
// @Produce      image/*
// @Param        id   path      string  true  "Product ID" format(uuid)
// @Success      200  {file}    file    "Product image"
// @Failure      400  {object}  response.ErrorResponseDetail "Invalid request format or invalid product ID"
// @Failure      401  {object}  response.ErrorResponseDetail "Missing or invalid access token"
// @Failure      404  {object}  response.ErrorResponseDetail "Product or image not found"
// @Failure      500  {object}  response.ErrorResponseDetail "Server error"
// @Security     AdminAccessToken
// @Router       /products/{id}/image [get]
func (h *Handler) GetProductImage(c *gin.Context) {
	productID, err := commonRequest.ParseParamUUID(c, "id")
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	token, err := commonRequest.ParseAccessToken(c)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	reader, contentType, err := h.uc.GetProductImage(c, productID, token)
	if err != nil {
		commonResponse.HandleError(c, err)
		return
	}

	c.Writer.Header().Set("Content-Type", contentType)
	c.Status(http.StatusOK)

	if _, err = io.Copy(c.Writer, reader); err != nil {
		commonResponse.HandleError(c, err)
		return
	}
}
