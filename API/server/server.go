package server

import (
	"pos-restaurant/handlers"

	"github.com/gin-gonic/gin"
)

func NewServer(
	menuHandler *handlers.MenuItemHandler,
	categoryHandler *handlers.MenuCategoryHandler,
	ingredientHandler *handlers.IngredientHandler,
	menuIngredientHandler *handlers.MenuIngredientHandler,

	outletHandler *handlers.OutletHandler,
	tableHandler *handlers.TableHandler,
	staffHandler *handlers.StaffHandler,
) *gin.Engine {

	r := gin.Default()
	api := r.Group("/api")

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

	return r
}
