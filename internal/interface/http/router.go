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

	// Initialize repositories and use cases
	handlers := initializeHandlers(db)

	// Setup routes
	setupPublicRoutes(e, handlers)
	setupAuthenticatedRoutes(e, handlers)

	return e
}

type Handlers struct {
	User     *handler.UserHandler
	Category *handler.CategoryHandler
	Product  *handler.ProductHandler
	Cart     *handler.CartHandler
	Order    *handler.OrderHandler
}

func initializeHandlers(db *gorm.DB) *Handlers {
	// Initialize repositories
	addressRepo := repository.NewAddressRepository(db)
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	cartItemRepo := repository.NewCartItemRepository(db)
	cartRepo := repository.NewCartRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// Initialize use cases
	userUC := usecase.NewUserUsecase(userRepo, addressRepo)
	catUC := usecase.NewCategoryUsecase(categoryRepo)
	prodUC := usecase.NewProductUsecase(productRepo)
	cartUC := usecase.NewCartUsecase(cartRepo, cartItemRepo, productRepo)
	orderUC := usecase.NewOrderUsecase(orderRepo, cartRepo, cartItemRepo, productRepo, userRepo, addressRepo)

	// Initialize handlers
	return &Handlers{
		User:     handler.NewUserHandler(userUC),
		Category: handler.NewCategoryHandler(catUC),
		Product:  handler.NewProductHandler(prodUC),
		Cart:     handler.NewCartHandler(cartUC),
		Order:    handler.NewOrderHandler(orderUC),
	}
}

func setupPublicRoutes(e *echo.Echo, h *Handlers) {
	e.Static("/images", "assets/images/")

	// Public user routes
	e.POST("/users/register", h.User.Register)
	e.POST("/users/login", h.User.Login)

	// Public category routes
	e.GET("/categories", h.Category.GetAll)
	e.GET("/categories/:id", h.Category.GetByID)
	e.GET("/categories/:id/subcategories", h.Category.GetSubcategories)
	e.GET("/categories/search", h.Category.Search)

	// Public product routes
	e.GET("/products", h.Product.GetAll)
	e.GET("/products/search", h.Product.Search)
	e.GET("/products/:id", h.Product.GetByID)
}

func setupAuthenticatedRoutes(e *echo.Echo, h *Handlers) {
	setupUserRoutes(e, h)
	setupCategoryRoutes(e, h)
	setupProductRoutes(e, h)
	setupCartRoutes(e, h)
	setupOrderRoutes(e, h)
}

func setupUserRoutes(e *echo.Echo, h *Handlers) {
	userGroup := e.Group("/users")
	userGroup.Use(auth.JWTMiddleware())
	userGroup.GET("/:id", h.User.GetByID)
	userGroup.GET("", h.User.GetAll)
	userGroup.GET("/search", h.User.Search)
	userGroup.PUT("/:id", h.User.Update)
	userGroup.DELETE("/:id", h.User.Delete)
}

func setupCategoryRoutes(e *echo.Echo, h *Handlers) {
	categoryGroup := e.Group("/categories")
	categoryGroup.Use(auth.JWTMiddleware())
	categoryGroup.POST("", h.Category.Create)
	categoryGroup.PUT("/:id", h.Category.Update)
	categoryGroup.DELETE("/:id", h.Category.Delete)
}

func setupProductRoutes(e *echo.Echo, h *Handlers) {
	productGroup := e.Group("/products")
	productGroup.Use(auth.JWTMiddleware())
	productGroup.POST("", h.Product.Create)
	productGroup.PUT("/:id", h.Product.Update)
	productGroup.DELETE("/:id", h.Product.Delete)
}

func setupCartRoutes(e *echo.Echo, h *Handlers) {
	cartGroup := e.Group("")
	cartGroup.Use(auth.JWTMiddleware())
	cartGroup.GET("/cart", h.Cart.GetByUserID)
	cartGroup.POST("/cart/add", h.Cart.AddProduct)
	cartGroup.PUT("/cart/item/:id", h.Cart.UpdateItem)
	cartGroup.DELETE("/cart/item/:id", h.Cart.RemoveItem)
	cartGroup.DELETE("/cart/clear", h.Cart.ClearCart)
	cartGroup.GET("/cart/search", h.Cart.Search)
}

func setupOrderRoutes(e *echo.Echo, h *Handlers) {
	orderGroup := e.Group("")
	orderGroup.Use(auth.JWTMiddleware())
	orderGroup.POST("/orders", h.Order.CreateOrder)
	orderGroup.GET("/orders/:id", h.Order.GetOrder)
	orderGroup.GET("/orders", h.Order.GetAllOrders)
	orderGroup.GET("/orders/user", h.Order.GetUserOrders)
	orderGroup.PUT("/orders/:id/status", h.Order.UpdateStatus)
	orderGroup.PUT("/orders/:id/cancel", h.Order.CancelOrder)
	orderGroup.GET("/orders/search", h.Order.Search)
}
