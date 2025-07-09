package server

import (
	"pos-restaurant/handlers"

	"github.com/gin-gonic/gin"

	_ "pos-restaurant/cmd/pos-restaurant/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewServer(
	menuHandler *handlers.MenuItemHandler,
	categoryHandler *handlers.MenuCategoryHandler,
	ingredientHandler *handlers.IngredientHandler,
	menuIngredientHandler *handlers.MenuIngredientHandler,

	outletHandler *handlers.OutletHandler,
	tableHandler *handlers.TableHandler,
	staffHandler *handlers.StaffHandler,

	customerHandler *handlers.CustomerHandler,
	visitHandler *handlers.CustomerVisitHandler,
	reservationHandler *handlers.ReservationHandler,

	orderHandler *handlers.OrderHandler,
	billHandler *handlers.BillHandler,
	tableTransferHandler *handlers.TableTransferHandler,
) *gin.Engine {

	r := gin.Default()
	api := r.Group("/api")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Menu-items Routes
	menu := api.Group("/menu")
	{
		// Admin use
		menu.POST("/menu-items", menuHandler.CreateMenuItem)
		menu.GET("/menu-items", menuHandler.ListMenuItems)
		menu.PUT("/menu-items/:id", menuHandler.UpdateMenuItem)
		menu.DELETE("/menu-items/:id", menuHandler.DeleteMenuItem)

		// Front use
		menu.GET("/menu-items-active", menuHandler.ListActiveMenuItems)      // Show only active menu
		menu.GET("/menu-items/category", menuHandler.GetMenuItemsByCategory) // Search by category || ?category_id=2
		menu.GET("/menu-items/search", menuHandler.SearchMenuItems)          // Search by name || ?search=ayam

		menu.GET("/menu-items/detail/:id", menuHandler.GetMenuDetail) // Show selected menu detail for ordering

		menu.POST("/category", categoryHandler.CreateCategory)
		menu.GET("/category", categoryHandler.ListCategories)
		menu.DELETE("/category/:id", categoryHandler.DeleteCategory)

	}

	// Ingredients Routes
	ingredient := api.Group("/ingredients")
	{
		ingredient.POST("/", ingredientHandler.CreateIngredient)
		ingredient.GET("/", ingredientHandler.ListIngredients)
		ingredient.GET("/:id", ingredientHandler.GetIngredientByID)
		ingredient.PUT("/:id", ingredientHandler.UpdateIngredient)
		ingredient.DELETE("/:id", ingredientHandler.DeleteIngredient)
	}

	// Menu <-> Ingredients Routes
	menuIngredient := api.Group("/menu-ingredients")
	{
		menuIngredient.POST("/", menuIngredientHandler.Create)
		menuIngredient.GET("/:menu_item_id", menuIngredientHandler.ListByMenuItem)
		menuIngredient.PUT("/:id", menuIngredientHandler.UpdateMenuIngredient)
		menuIngredient.DELETE("/:id", menuIngredientHandler.DeleteMenuIngredient)
	}

	// Outlet Routes
	outlet := api.Group("/outlets")
	{
		outlet.POST("/", outletHandler.Create)
		outlet.GET("/", outletHandler.List)
		outlet.GET("/:id", outletHandler.GetByID)
		outlet.PUT("/:id", outletHandler.Update)
		outlet.DELETE("/:id", outletHandler.Delete)
	}

	// Table Routes
	table := api.Group("tables")
	{
		table.POST("/", tableHandler.Create)
		table.GET("/", tableHandler.List)
		table.GET("/:id", tableHandler.GetByID)
		table.PUT("/:id", tableHandler.Update)
		table.DELETE("/:id", tableHandler.Delete)
	}

	// Staff Routes
	staff := api.Group("/staff")
	{
		staff.POST("/", staffHandler.Create)
		staff.GET("/", staffHandler.List)
		staff.GET("/:id", staffHandler.GetByID)
		staff.PUT("/:id", staffHandler.Update)
		staff.DELETE("/:id", staffHandler.SoftDelete)
	}

	customer := api.Group("/customers")
	{
		customer.POST("/", customerHandler.Create)
		customer.GET("/", customerHandler.List)
		customer.GET("/:id", customerHandler.GetByID)
		customer.PUT("/:id", customerHandler.Update)
		customer.DELETE("/:id", customerHandler.SoftDelete)
	}

	visits := api.Group("/visits")
	{
		visits.POST("/", visitHandler.Create)
		visits.GET("/", visitHandler.List)
		visits.GET("/:id", visitHandler.GetByID)
		visits.PUT("/:id", visitHandler.Update)
		visits.DELETE("/:id", visitHandler.Delete)
	}

	group := api.Group("/reservations")
	{
		group.POST("/", reservationHandler.Create)
		group.GET("/", reservationHandler.List)
		group.GET("/:id", reservationHandler.GetByID)
		group.PUT("/:id", reservationHandler.Update)
		group.DELETE("/:id", reservationHandler.Delete)
	}

	// Orders
	orders := api.Group("/orders")
	{
		orders.POST("/", orderHandler.Create)
		orders.GET("/", orderHandler.List)
		orders.GET("/:id", orderHandler.GetByID)
		orders.PUT("/:id", orderHandler.Update)
		orders.POST("/:id/add", orderHandler.AddItem)
		orders.DELETE("/:id", orderHandler.Delete)
	}

	// Bill
	bills := api.Group("/bills")
	{
		bills.POST("/", billHandler.Create)
		bills.POST("/split", billHandler.CreateSplit)
		bills.GET("/", billHandler.List)
		bills.GET("/:id", billHandler.GetByID)
		bills.DELETE("/:id", billHandler.Delete)

		bills.POST("/pay", billHandler.Pay)
	}

	tabletf := api.Group("/table-transfer")
	{
		tabletf.POST("/", tableTransferHandler.Create)
		tabletf.GET("/", tableTransferHandler.List)
		tabletf.GET("/:id", tableTransferHandler.GetByID)
		tabletf.PUT("/:id", tableTransferHandler.Update)
		tabletf.DELETE("/:id", tableTransferHandler.Delete)
	}

	return r
}
