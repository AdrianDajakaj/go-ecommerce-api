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

type UserHandler struct {
	Usecase usecase.UserUsecase
}

func NewUserHandler(uc usecase.UserUsecase) *UserHandler {
	return &UserHandler{Usecase: uc}
}

func (h *UserHandler) GetByID(c echo.Context) error {
	uidToken, err := auth.UserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}
	role, err := auth.RoleFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}

	id, err := parseUintParam(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	if role != "admin" && uidToken != id {
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	user, err := h.Usecase.GetByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetAll(c echo.Context) error {
	role, err := auth.RoleFromContext(c)
	if err != nil || role != "admin" {
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	users, err := h.Usecase.GetAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) Search(c echo.Context) error {
	role, err := auth.RoleFromContext(c)
	if err != nil || role != "admin" {
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	filters := map[string]string{}
	for key, vals := range c.QueryParams() {
		if len(vals) > 0 {
			filters[key] = vals[0]
		}
	}
	users, err := h.Usecase.GetWithFilters(filters)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

type registerInput struct {
	Email    string        `json:"email"`
	Password string        `json:"password"`
	Name     string        `json:"name"`
	Surname  string        `json:"surname"`
	Address  model.Address `json:"address"`
}

func (h *UserHandler) Register(c echo.Context) error {
	var input registerInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	user := &model.User{
		Email:   input.Email,
		Name:    input.Name,
		Surname: input.Surname,
	}

	createdUser, err := h.Usecase.Register(user, input.Password, &input.Address)
	if err != nil {
		if err.Error() == "email already in use" {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, createdUser)
}

type loginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) Login(c echo.Context) error {
	var input loginInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}
	if err := c.Validate(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.Usecase.Login(input.Email, input.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not generate token")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
		"user":  user,
	})
}

func (h *UserHandler) Update(c echo.Context) error {
	uidToken, err := auth.UserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}
	role, err := auth.RoleFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}

	id, err := parseUintParam(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	if role != "admin" && uidToken != id {
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	var input model.User
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	input.ID = id
	input.Role = ""

	updated, err := h.Usecase.Update(&input)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, updated)
}

func (h *UserHandler) Delete(c echo.Context) error {
	uidToken, err := auth.UserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}
	role, err := auth.RoleFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}

	id, err := parseUintParam(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	if role != "admin" && uidToken != id {
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	err = h.Usecase.Delete(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func parseUintParam(c echo.Context, name string) (uint, error) {
	idParam := c.Param(name)
	parsed, err := strconv.ParseUint(idParam, 10, 64)
	return uint(parsed), err
}
