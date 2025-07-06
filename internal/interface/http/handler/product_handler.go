package handler

import (
	"errors"
	"net/http"

	"go-ecommerce-api/internal/domain/model"
	"go-ecommerce-api/internal/infrastructure/auth"
	"go-ecommerce-api/internal/usecase"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Error message constants
const (
	errInvalidProductID = "invalid product ID"
	errProductNotFound  = "product not found"
	errAccessDenied     = "access denied"
	errInvalidBody      = "invalid request body"
)

type ProductHandler struct {
	Usecase usecase.ProductUsecase
}

func NewProductHandler(uc usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{Usecase: uc}
}

// checkAdminRole verifies if the user has admin role
func (h *ProductHandler) checkAdminRole(c echo.Context) error {
	role, err := auth.RoleFromContext(c)
	if err != nil || role != "admin" {
		return echo.NewHTTPError(http.StatusForbidden, errAccessDenied)
	}
	return nil
}

func (h *ProductHandler) GetByID(c echo.Context) error {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errInvalidProductID)
	}
	prod, err := h.Usecase.GetByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, errProductNotFound)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, prod)
}

func (h *ProductHandler) GetAll(c echo.Context) error {
	prods, err := h.Usecase.GetAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, prods)
}

func (h *ProductHandler) Search(c echo.Context) error {
	filters := map[string]string{}
	for key, vals := range c.QueryParams() {
		if len(vals) > 0 {
			filters[key] = vals[0]
		}
	}
	prods, err := h.Usecase.GetWithFilters(filters)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, prods)
}

func (h *ProductHandler) Create(c echo.Context) error {
	if err := h.checkAdminRole(c); err != nil {
		return err
	}

	var input model.Product
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errInvalidBody)
	}
	created, err := h.Usecase.Create(&input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, created)
}

func (h *ProductHandler) Update(c echo.Context) error {
	if err := h.checkAdminRole(c); err != nil {
		return err
	}

	id, err := parseUintParam(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errInvalidProductID)
	}
	var input model.Product
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errInvalidBody)
	}
	input.ID = id
	updated, err := h.Usecase.Update(&input)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, errProductNotFound)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, updated)
}

func (h *ProductHandler) Delete(c echo.Context) error {
	if err := h.checkAdminRole(c); err != nil {
		return err
	}

	id, err := parseUintParam(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errInvalidProductID)
	}
	if err := h.Usecase.Delete(id); errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, errProductNotFound)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
