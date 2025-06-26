package server

import (
	"pos-restaurant/handlers"

	"github.com/gin-gonic/gin"
)

func NewServer(
	menuHandler *handlers.MenuItemHandler,
	categoryHandler *handlers.MenuCategoryHandler,
	ingredientHandler *handlers.IngredientHandler,
) *gin.Engine {

	r := gin.Default()
	api := r.Group("/api")

	// Grup route untuk API v1
	menu := api.Group("/menu")
	{
		menu.POST("/menu-items", menuHandler.CreateMenuItem)
		menu.PUT("/menu-items/:id", menuHandler.UpdateMenuItem)
		menu.GET("/menu-items", menuHandler.ListMenuItems) // Admin
		menu.GET("/menu-items/:id", menuHandler.ListMenuItemsById)
		menu.GET("/menu-items-active", menuHandler.ListActiveMenuItems)
		menu.GET("/menu-items/search", menuHandler.SearchMenuItems) // ?search=ayam
		menu.DELETE("/menu-items/:id", menuHandler.DeleteMenuItem)

		menu.POST("/category", categoryHandler.CreateCategory)
		menu.GET("/category", categoryHandler.ListCategories)
		menu.DELETE("/category/:id", categoryHandler.DeleteCategory)

	}

	ingredient := api.Group("/ingredients")
	{
		ingredient.POST("/", ingredientHandler.CreateIngredient)
		ingredient.GET("/", ingredientHandler.ListIngredients)
		ingredient.PUT("/:id", ingredientHandler.UpdateIngredient)
		ingredient.DELETE("/:id", ingredientHandler.DeleteIngredient)
	}

	return r
}
