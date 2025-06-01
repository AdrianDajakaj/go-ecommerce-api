package http

import (
	"go-ecommerce-api/internal/infrastructure/auth"
	"go-ecommerce-api/internal/infrastructure/persistence/repository"
	"go-ecommerce-api/internal/interface/http/handler"
	"go-ecommerce-api/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func NewRouter(db *gorm.DB) *echo.Echo {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	addressRepo := repository.NewAddressRepository(db)
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	cartItemRepo := repository.NewCartItemRepository(db)
	cartRepo := repository.NewCartRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	userUC := usecase.NewUserUsecase(userRepo, addressRepo)
	catUC := usecase.NewCategoryUsecase(categoryRepo)
	prodUC := usecase.NewProductUsecase(productRepo)
	cartUC := usecase.NewCartUsecase(cartRepo, cartItemRepo, productRepo)
	orderUC := usecase.NewOrderUsecase(orderRepo, cartRepo, cartItemRepo, productRepo, userRepo, addressRepo)

	userHandler := handler.NewUserHandler(userUC)
	catHandler := handler.NewCategoryHandler(catUC)
	prodHandler := handler.NewProductHandler(prodUC)
	cartHandler := handler.NewCartHandler(cartUC)
	orderHandler := handler.NewOrderHandler(orderUC)

	e.Static("/images", "assets/images/")

	e.POST("/users/register", userHandler.Register)
	e.POST("/users/login", userHandler.Login)
	userGroup := e.Group("/users")
	userGroup.Use(auth.JWTMiddleware())
	userGroup.GET("/:id", userHandler.GetByID)
	userGroup.GET("", userHandler.GetAll)
	userGroup.GET("/search", userHandler.Search)
	userGroup.PUT("/:id", userHandler.Update)
	userGroup.DELETE("/:id", userHandler.Delete)

	e.GET("/categories", catHandler.GetAll)
	e.GET("/categories/:id", catHandler.GetByID)
	e.GET("/categories/:id/subcategories", catHandler.GetSubcategories)
	e.GET("/categories/search", catHandler.Search)
	categoryGroup := e.Group("/categories")
	categoryGroup.Use(auth.JWTMiddleware())
	categoryGroup.POST("", catHandler.Create)
	categoryGroup.PUT("/:id", catHandler.Update)
	categoryGroup.DELETE("/:id", catHandler.Delete)

	e.GET("/products", prodHandler.GetAll)
	e.GET("/products/search", prodHandler.Search)
	e.GET("/products/:id", prodHandler.GetByID)
	productGroup := e.Group("/products")
	productGroup.Use(auth.JWTMiddleware())
	productGroup.POST("", prodHandler.Create)
	productGroup.PUT("/:id", prodHandler.Update)
	productGroup.DELETE("/:id", prodHandler.Delete)

	r := e.Group("")
	r.Use(auth.JWTMiddleware())
	r.GET("/cart", cartHandler.GetByUserID)
	r.POST("/cart/add", cartHandler.AddProduct)
	r.PUT("/cart/item/:id", cartHandler.UpdateItem)
	r.DELETE("/cart/item/:id", cartHandler.RemoveItem)
	r.DELETE("/cart/clear", cartHandler.ClearCart)
	r.GET("/cart/search", cartHandler.Search)

	r.POST("/orders", orderHandler.CreateOrder)
	r.GET("/orders/:id", orderHandler.GetOrder)
	r.GET("/orders", orderHandler.GetAllOrders)
	r.GET("/orders/user", orderHandler.GetUserOrders)
	r.PUT("/orders/:id/status", orderHandler.UpdateStatus)
	r.PUT("/orders/:id/cancel", orderHandler.CancelOrder)
	r.GET("/orders/search", orderHandler.Search)

	return e
}
