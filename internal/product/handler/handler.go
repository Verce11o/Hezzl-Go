package handler

import (
	"errors"
	"github.com/Verce11o/Hezzl-Go/api"
	"github.com/Verce11o/Hezzl-Go/internal/product"
	"github.com/Verce11o/Hezzl-Go/lib/request"
	"github.com/Verce11o/Hezzl-Go/lib/response"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	log     *zap.SugaredLogger
	service product.Service
}

func NewHandler(log *zap.SugaredLogger, service product.Service) *Handler {
	return &Handler{log: log, service: service}
}

func (h *Handler) CreateProduct(c *gin.Context, params api.CreateProductParams) {
	var input api.CreateProductJSONBody

	if err := request.Read(c, &input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	prod, err := h.service.CreateProduct(c.Request.Context(), params.ProjectId, input)
	if err != nil {
		h.log.Errorf("error while creating product: %v", err)
		response.WithHTTPError(c, err)
		return
	}

	c.JSON(http.StatusOK, prod)
}

func (h *Handler) UpdateProduct(c *gin.Context, params api.UpdateProductParams) {
	var input api.UpdateProductJSONBody

	if err := request.Read(c, &input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	prod, err := h.service.UpdateProduct(c.Request.Context(), params.Id, params.ProjectId, input)

	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    3,
			"message": "errors.good.notFound",
			"details": gin.H{},
		})
		return
	}

	if err != nil {
		h.log.Errorf("error while updating product: %v", err)
		response.WithHTTPError(c, err)
		return
	}

	c.JSON(http.StatusOK, prod)
}

func (h *Handler) DeleteProduct(c *gin.Context, params api.DeleteProductParams) {

	err := h.service.DeleteProduct(c.Request.Context(), params.Id, params.ProjectId)

	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    3,
			"message": "errors.good.notFound",
			"details": gin.H{},
		})
		return
	}

	if err != nil {
		h.log.Errorf("error while deleting product: %v", err)
		response.WithHTTPError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        params.Id,
		"projectId": params.ProjectId,
		"removed":   true,
	})
}

func (h *Handler) GetProducts(c *gin.Context, params api.GetProductsParams) {
	var limit, offset int

	if params.Limit != nil {
		limit = *params.Limit
	}
	if params.Offset != nil {
		offset = *params.Offset
	}

	products, err := h.service.GetProductsList(c.Request.Context(), limit, offset)
	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    3,
			"message": "errors.good.notFound",
			"details": gin.H{},
		})
		return
	}

	if err != nil {
		h.log.Errorf("error while getting products: %v", err)
		response.WithHTTPError(c, err)
		return
	}

	c.JSON(http.StatusOK, products)

}

func (h *Handler) PatchGoodsReprioritize(c *gin.Context, params api.PatchGoodsReprioritizeParams) {

	var input api.PatchGoodsReprioritizeJSONBody
	if err := request.Read(c, &input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	priorities, err := h.service.UpdateProductPriority(c.Request.Context(), params.Id, params.ProjectId, input.NewPriority)
	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    3,
			"message": "errors.good.notFound",
			"details": gin.H{},
		})
		return
	}

	if err != nil {
		h.log.Errorf("error while update product: %v", err)
		response.WithHTTPError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"priorities": priorities,
	})

}
