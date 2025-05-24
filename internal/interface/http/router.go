package http

import (
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

	addressRepository := repository.NewAddressRepository(db)
	userRepository := repository.NewUserRepository(db)
	categoryRepository := repository.NewCategoryRepository(db)
	productRepository := repository.NewProductRepository(db)
	cartItemRepository := repository.NewCartItemRepository(db)
	cartRepository := repository.NewCartRepository(db)
	orderRepository := repository.NewOrderRepository(db)

	userUsecase := usecase.NewUserUsecase(userRepository, addressRepository)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepository)
	productUsecase := usecase.NewProductUsecase(productRepository)
	cartUsecase := usecase.NewCartUsecase(cartRepository, cartItemRepository, productRepository)
	orderUsecase := usecase.NewOrderUsecase(
		orderRepository,
		cartRepository,
		cartItemRepository,
		productRepository,
		userRepository,
		addressRepository,
	)

	userHandler := handler.NewUserHandler(userUsecase)
	categoryHandler := handler.NewCategoryHandler(categoryUsecase)
	productHandler := handler.NewProductHandler(productUsecase)
	cartHandler := handler.NewCartHandler(cartUsecase)
	orderHandler := handler.NewOrderHandler(orderUsecase)

	e.POST("/users/register", userHandler.Register)
	e.POST("/users/login", userHandler.Login)
	e.GET("/users/:id", userHandler.GetByID)
	e.GET("/users", userHandler.GetAll)
	e.GET("/users/search", userHandler.Search)
	e.PUT("/users/:id", userHandler.Update)
	e.DELETE("/users/:id", userHandler.Delete)

	e.GET("/categories", categoryHandler.GetAll)
	e.GET("/categories/:id", categoryHandler.GetByID)
	e.POST("/categories", categoryHandler.Create)
	e.PUT("/categories/:id", categoryHandler.Update)
	e.DELETE("/categories/:id", categoryHandler.Delete)
	e.GET("/categories/search", categoryHandler.Search)

	e.GET("/products", productHandler.GetAll)
	e.GET("/products/search", productHandler.Search)
	e.GET("/products/:id", productHandler.GetByID)
	e.POST("/products", productHandler.Create)
	e.PUT("/products/:id", productHandler.Update)
	e.DELETE("/products/:id", productHandler.Delete)

	e.GET("/cart/:user_id", cartHandler.GetByUserID)
	e.POST("/cart/:user_id/add", cartHandler.AddProduct)
	e.PUT("/cart/item/:id", cartHandler.UpdateItem)
	e.DELETE("/cart/item/:id", cartHandler.RemoveItem)
	e.DELETE("/cart/:user_id/clear", cartHandler.ClearCart)
	e.GET("/cart/search", cartHandler.Search)

	e.POST("/users/:user_id/orders", orderHandler.CreateOrder)
	e.GET("/orders/:id", orderHandler.GetOrder)
	e.GET("/users/:user_id/orders", orderHandler.GetUserOrders)
	e.GET("/orders", orderHandler.GetAllOrders)
	e.PUT("/orders/:id/status", orderHandler.UpdateStatus)
	e.PUT("/orders/:id/cancel", orderHandler.CancelOrder)
	e.GET("/orders/search", orderHandler.Search)

	return e
}
