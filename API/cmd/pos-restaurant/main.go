package main

import (
	"log"
	"pos-restaurant/database"
	"pos-restaurant/handlers"
	"pos-restaurant/repositories"
	"pos-restaurant/server"
	"pos-restaurant/services"
)

func main() {
	// Connect Database
	err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Gagal inisialisasi database: %v", err)
	}
	defer database.DB.Close()

	// Repo Init
	menuRepo := repositories.NewMenuItemRepository(database.DB)
	categoryRepo := repositories.NewMenuCategoryRepository(database.DB)
	ingredientRepo := repositories.NewIngredientRepository(database.DB)

	// Service Init
	menuService := services.NewMenuService(menuRepo)
	categoryService := services.NewMenuCategoryService(categoryRepo)
	ingredientService := services.NewIngredientService(ingredientRepo)

	// Handler init
	menuHandler := handlers.NewMenuItemHandler(menuService)
	categoryHandler := handlers.NewMenuCategoryHandler(categoryService)
	ingredientHandler := handlers.NewIngredientHandler(ingredientService)

	// Create and Start server
	srv := server.NewServer(
		menuHandler,
		categoryHandler,
		ingredientHandler,
	)
	log.Printf("Server starting on port 8080")
	if err := srv.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}
