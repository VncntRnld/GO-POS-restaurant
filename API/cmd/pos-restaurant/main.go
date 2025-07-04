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
	menuIngredientRepo := repositories.NewMenuIngredientRepository(database.DB)

	outletRepo := repositories.NewOutletRepository(database.DB)
	tableRepo := repositories.NewTableRepository(database.DB)
	staffRepo := repositories.NewStaffRepository(database.DB)

	customerRepo := repositories.NewCustomerRepository(database.DB)
	customerVisitRepo := repositories.NewCustomerVisitRepository(database.DB)
	reservationRepo := repositories.NewReservationRepository(database.DB)

	orderRepo := repositories.NewOrderRepository(database.DB)
	billRepo := repositories.NewBillRepository(database.DB)

	// Service Init
	menuService := services.NewMenuService(menuRepo)
	categoryService := services.NewMenuCategoryService(categoryRepo)
	ingredientService := services.NewIngredientService(ingredientRepo)
	menuIngredientService := services.NewMenuIngredientService(menuIngredientRepo)

	outletService := services.NewOutletService(outletRepo)
	tableService := services.NewTableService(tableRepo)
	staffService := services.NewStaffService(staffRepo)

	customerService := services.NewCustomerService(customerRepo)
	customerVisitService := services.NewCustomerVisitService(customerVisitRepo)
	reservationService := services.NewReservationService(reservationRepo)

	OrderService := services.NewOrderService(orderRepo)
	billService := services.NewBillService(billRepo)

	// Handler init
	menuHandler := handlers.NewMenuItemHandler(menuService)
	categoryHandler := handlers.NewMenuCategoryHandler(categoryService)
	ingredientHandler := handlers.NewIngredientHandler(ingredientService)
	menuIngredientHandler := handlers.NewMenuIngredientHandler(menuIngredientService)

	outletHandler := handlers.NewOutletHandler(outletService)
	tableHandler := handlers.NewTableHandler(tableService)
	staffHandler := handlers.NewStaffHandler(staffService)

	customerHandler := handlers.NewCustomerHandler(customerService)
	customerVisitHandler := handlers.NewCustomerVisitHandler(customerVisitService)
	reservationHandler := handlers.NewReservationHandler(reservationService)

	orderHandler := handlers.NewOrderHandler(OrderService)
	billHandler := handlers.NewBillHandler(billService)

	// Create and Start server
	srv := server.NewServer(
		menuHandler,
		categoryHandler,
		ingredientHandler,
		menuIngredientHandler,

		outletHandler,
		tableHandler,
		staffHandler,

		customerHandler,
		customerVisitHandler,
		reservationHandler,

		orderHandler,
		billHandler,
	)

	log.Printf("Server starting on port 8080")
	if err := srv.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}
