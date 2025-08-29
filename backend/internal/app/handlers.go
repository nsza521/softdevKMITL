package app

import (
	customerHttp"backend/internal/customer/delivery/http"
	customerUsecase "backend/internal/customer/usecase"
	customerRepository "backend/internal/customer/repository"

	// foodOrderHttp "backend/internal/order/delivery/http"
	// foodOrderUsecase "backend/internal/order/usecase"
	// foodOrderRepository "backend/internal/order/repository"
)

func (s *App) MapHandlers() error {

	// userGroup := s.gin.Group("/user")
	customerGroup := s.gin.Group("/user/customer")
	// restaurantGroup := s.gin.Group("/user/restaurant")
	// tableGroup := s.gin.Group("/table")
	// tableReservationGroup := s.gin.Group("/table/reservation")
	// foodOrderGroup := s.gin.Group("/food/order")
	// notificationGroup := s.gin.Group("/notification")
	// paymentGroup := s.gin.Group("/payment")
	// menuGroup := s.gin.Group("/food/menu")

	// Customer Group
	customerRepository := customerRepository.NewCustomerRepository(s.db)
	customerUsecase := customerUsecase.NewCustomerUsecase(customerRepository)
	customerHandler := customerHttp.NewCustomerHandler(customerUsecase)
	customerHttp.MapCustomerRoutes(customerGroup, customerHandler)

	// Food Order Group
	// foodOrderRepository := foodOrderRepository.NewFoodOrderRepository(s.db)
	// foodOrderUsecase := foodOrderUsecase.NewFoodOrderUsecase(foodOrderRepository)
	// foodOrderHandler := foodOrderHttp.NewFoodOrderHandler(foodOrderUsecase)
	// foodOrderHttp.MapFoodOrderRoutes(foodOrderGroup, foodOrderHandler)

	return nil
}