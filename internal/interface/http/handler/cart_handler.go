package handler

import (
	"errors"
	"net/http"

	"go-ecommerce-api/internal/usecase"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CartHandler struct {
	Usecase usecase.CartUsecase
}

func NewCartHandler(uc usecase.CartUsecase) *CartHandler {
	return &CartHandler{Usecase: uc}
}

func (h *CartHandler) GetByUserID(c echo.Context) error {
	userID, err := parseUintParam(c, "user_id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	cart, err := h.Usecase.GetByUserID(userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "cart not found")
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
	userID, err := parseUintParam(c, "user_id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	var req addReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	cart, err := h.Usecase.AddProduct(userID, req.ProductID, req.Quantity)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "product not found")
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
		return echo.NewHTTPError(http.StatusBadRequest, "invalid item ID")
	}

	var req updateReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	cart, err := h.Usecase.UpdateItem(itemID, req.Quantity)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "item or cart not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) RemoveItem(c echo.Context) error {
	itemID, err := parseUintParam(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid item ID")
	}

	cart, err := h.Usecase.RemoveItem(itemID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "item or cart not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) ClearCart(c echo.Context) error {
	userID, err := parseUintParam(c, "user_id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	cart, err := h.Usecase.ClearCart(userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "cart not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, cart)
}
