package handler

import (
	"errors"
	"fmt"
	"net/http"

	"go-ecommerce-api/internal/domain/model"
	"go-ecommerce-api/internal/infrastructure/auth"
	"go-ecommerce-api/internal/usecase"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const (
	invalidCategoryIDMsg = "invalid category ID"
	categoryNotFoundMsg  = "category not found"
	accessDeniedMsg      = "access denied"
)

type CategoryHandler struct {
	Usecase usecase.CategoryUsecase
}

func NewCategoryHandler(uc usecase.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{Usecase: uc}
}

func (h *CategoryHandler) GetByID(c echo.Context) error {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, invalidCategoryIDMsg)
	}

	category, err := h.Usecase.GetByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, categoryNotFoundMsg)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) GetAll(c echo.Context) error {
	categories, err := h.Usecase.GetAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) GetSubcategories(c echo.Context) error {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, invalidCategoryIDMsg)
	}
	filters := map[string]string{
		"parent_id":          fmt.Sprint(id),
		"with_subcategories": "1",
	}
	cats, err := h.Usecase.GetWithFilters(filters)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, cats)
}

func (h *CategoryHandler) Search(c echo.Context) error {
	filters := map[string]string{}
	for key, vals := range c.QueryParams() {
		if len(vals) > 0 {
			filters[key] = vals[0]
		}
	}

	cats, err := h.Usecase.GetWithFilters(filters)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, cats)
}

func (h *CategoryHandler) Create(c echo.Context) error {
	role, err := auth.RoleFromContext(c)
	if err != nil || role != "admin" {
		return echo.NewHTTPError(http.StatusForbidden, accessDeniedMsg)
	}

	var input model.Category
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	created, err := h.Usecase.Create(&input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, created)
}

func (h *CategoryHandler) Update(c echo.Context) error {
	role, err := auth.RoleFromContext(c)
	if err != nil || role != "admin" {
		return echo.NewHTTPError(http.StatusForbidden, accessDeniedMsg)
	}

	id, err := parseUintParam(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, invalidCategoryIDMsg)
	}

	var input model.Category
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	input.ID = id

	updated, err := h.Usecase.Update(&input)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, categoryNotFoundMsg)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, updated)
}

func (h *CategoryHandler) Delete(c echo.Context) error {
	role, err := auth.RoleFromContext(c)
	if err != nil || role != "admin" {
		return echo.NewHTTPError(http.StatusForbidden, accessDeniedMsg)
	}

	id, err := parseUintParam(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, invalidCategoryIDMsg)
	}

	err = h.Usecase.Delete(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, categoryNotFoundMsg)
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
