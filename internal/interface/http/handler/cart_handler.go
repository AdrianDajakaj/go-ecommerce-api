package handler

import (
	"errors"
	"net/http"

	"go-ecommerce-api/internal/infrastructure/auth"
	"go-ecommerce-api/internal/usecase"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const (
	invalidTokenMsg       = "invalid token"
	cartNotFoundMsg       = "cart not found"
	productNotFoundMsg    = "product not found"
	itemNotFoundMsg       = "item or cart not found"
	invalidRequestBodyMsg = "invalid request body"
	invalidItemIDMsg      = "invalid item ID"
)

type CartHandler struct {
	Usecase usecase.CartUsecase
}

func NewCartHandler(uc usecase.CartUsecase) *CartHandler {
	return &CartHandler{Usecase: uc}
}

func (h *CartHandler) GetByUserID(c echo.Context) error {
	userID, err := auth.UserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, invalidTokenMsg)
	}

	cart, err := h.Usecase.GetByUserID(userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, cartNotFoundMsg)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) Search(c echo.Context) error {
	filters := map[string]string{}
	for key, vals := range c.QueryParams() {
		if len(vals) > 0 {
			filters[key] = vals[0]
		}
	}

	carts, err := h.Usecase.GetWithFilters(filters)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, carts)
}

type addReq struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

func (h *CartHandler) AddProduct(c echo.Context) error {
	userID, err := auth.UserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, invalidTokenMsg)
	}

	var req addReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, invalidRequestBodyMsg)
	}

	cart, err := h.Usecase.AddProduct(userID, req.ProductID, req.Quantity)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, productNotFoundMsg)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, cart)
}

type updateReq struct {
	Quantity int `json:"quantity"`
}

func (h *CartHandler) UpdateItem(c echo.Context) error {
	itemID, err := parseUintParam(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, invalidItemIDMsg)
	}

	var req updateReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, invalidRequestBodyMsg)
	}

	cart, err := h.Usecase.UpdateItem(itemID, req.Quantity)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, itemNotFoundMsg)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) RemoveItem(c echo.Context) error {
	itemID, err := parseUintParam(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, invalidItemIDMsg)
	}

	cart, err := h.Usecase.RemoveItem(itemID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, itemNotFoundMsg)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) ClearCart(c echo.Context) error {
	userID, err := auth.UserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, invalidTokenMsg)
	}

	cart, err := h.Usecase.ClearCart(userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, cartNotFoundMsg)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, cart)
}
