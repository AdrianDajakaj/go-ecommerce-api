package handler

import (
	"errors"
	"net/http"
	"strconv"

	"go-ecommerce-api/internal/domain/model"
	"go-ecommerce-api/internal/infrastructure/auth"
	"go-ecommerce-api/internal/usecase"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type OrderHandler struct {
	usecase usecase.OrderUsecase
}

func NewOrderHandler(uc usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{usecase: uc}
}

func (h *OrderHandler) GetOrder(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid order ID")
	}

	order, err := h.usecase.GetByID(uint(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "order not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	uid, err := auth.UserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}
	role, err := auth.RoleFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}

	if role != "admin" && order.UserID != uid {
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	return c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetUserOrders(c echo.Context) error {
	uid, err := auth.UserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}

	orders, err := h.usecase.GetByUserID(uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetAllOrders(c echo.Context) error {
	role, err := auth.RoleFromContext(c)
	if err != nil || role != "admin" {
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	orders, err := h.usecase.GetAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) Search(c echo.Context) error {
	role, errRole := auth.RoleFromContext(c)
	if errRole != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}

	filters := map[string]string{}

	if role == "admin" {
		for key, vals := range c.QueryParams() {
			if len(vals) > 0 {
				filters[key] = vals[0]
			}
		}
	} else {
		uid, errUID := auth.UserIDFromContext(c)
		if errUID != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}
		filters["user_id"] = strconv.FormatUint(uint64(uid), 10)
	}

	orders, err := h.usecase.GetWithFilters(filters)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, orders)
}

type createOrderRequest struct {
	PaymentMethod     model.PaymentMethod `json:"payment_method" validate:"required"`
	ShippingAddressID uint                `json:"shipping_address_id" validate:"required"`
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {
	uid, err := auth.UserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}

	var req createOrderRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	order, err := h.usecase.CreateFromCart(uid, req.PaymentMethod, req.ShippingAddressID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, order)
}

type updateStatusRequest struct {
	Status model.OrderStatus `json:"status" validate:"required"`
}

func (h *OrderHandler) UpdateStatus(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid order ID")
	}

	role, err := auth.RoleFromContext(c)
	if err != nil || role != "admin" {
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	var req updateStatusRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	order, err := h.usecase.UpdateStatus(uint(id), req.Status)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "order not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) CancelOrder(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid order ID")
	}

	order, err := h.usecase.GetByID(uint(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "order not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	uid, err := auth.UserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}
	role, err := auth.RoleFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}

	if role != "admin" && order.UserID != uid {
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	updatedOrder, err := h.usecase.CancelOrder(uint(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "order not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, updatedOrder)
}
