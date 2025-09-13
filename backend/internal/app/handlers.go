package app

import (
	userHttp "backend/internal/user/delivery/http"
	userUsecase "backend/internal/user/usecase"
	userRepository "backend/internal/user/repository"

	customerHttp"backend/internal/customer/delivery/http"
	customerUsecase "backend/internal/customer/usecase"
	customerRepository "backend/internal/customer/repository"

	restaurantHttp "backend/internal/restaurant/delivery/http"
	restaurantUsecase "backend/internal/restaurant/usecase"
	restaurantRepository "backend/internal/restaurant/repository"

	tableHttp "backend/internal/table/delivery/http"
	tableUsecase "backend/internal/table/usecase"
	tableRepository "backend/internal/table/repository"

	tableReservationHttp "backend/internal/reservation/delivery/http"
	tableReservationUsecase "backend/internal/reservation/usecase"
	tableReservationRepository "backend/internal/reservation/repository"

	paymentHttp "backend/internal/payment/delivery/http"
	paymentUsecase "backend/internal/payment/usecase"
	paymentRepository "backend/internal/payment/repository"

	foodOrderHttp "backend/internal/order/delivery/http"
	foodOrderUsecase "backend/internal/order/usecase"
	foodOrderRepository "backend/internal/order/repository"

	notiHttp "backend/internal/notifications/delivery/http"
	notiUsecase "backend/internal/notifications/usecase"
	notiRepository "backend/internal/notifications/repository"
)

func (s *App) MapHandlers() error {

	userGroup := s.gin.Group("/user")
	customerGroup := s.gin.Group("/customer")
	restaurantGroup := s.gin.Group("/restaurant")
	tableGroup := s.gin.Group("/table")
	tableReservationGroup := s.gin.Group("/table/reservation")
	foodOrderGroup := s.gin.Group("/food/order")
	notificationGroup := s.gin.Group("/notification")
	paymentGroup := s.gin.Group("/payment")

	// Customer Group
	customerRepository := customerRepository.NewCustomerRepository(s.db)
	customerUsecase := customerUsecase.NewCustomerUsecase(customerRepository)
	customerHandler := customerHttp.NewCustomerHandler(customerUsecase)
	customerHttp.MapCustomerRoutes(customerGroup, customerHandler)

	// Restaurant Group
	menuRepository := restaurantRepository.NewMenuRepository(s.db)
	restaurantRepository := restaurantRepository.NewRestaurantRepository(s.db)
	restaurantUsecase := restaurantUsecase.NewRestaurantUsecase(restaurantRepository, menuRepository)
	restaurantHandler := restaurantHttp.NewRestaurantHandler(restaurantUsecase)
	restaurantHttp.MapRestaurantRoutes(restaurantGroup, restaurantHandler)

	// User Group
	userRepository := userRepository.NewUserRepository(s.db)
	// userUsecase := userUsecase.NewUserUsecase(userRepository, customerUsecase, restaurantUsecase)
	userUsecase := userUsecase.NewUserUsecase(userRepository)
	userHandler := userHttp.NewUserHandler(userUsecase, customerUsecase, restaurantUsecase)
	userHttp.MapUserRoutes(userGroup, userHandler)

	// Table Group
	tableRepository := tableRepository.NewTableRepository(s.db)
	tableUsecase := tableUsecase.NewTableUsecase(tableRepository)
	tableHandler := tableHttp.NewTableHandler(tableUsecase)
	tableHttp.MapTableRoutes(tableGroup, tableHandler)
	
	// Table Reservation Group
	tableReservationRepository := tableReservationRepository.NewTableReservationRepository(s.db)
	tableReservationUsecase := tableReservationUsecase.NewTableReservationUsecase(tableReservationRepository)
	tableReservationHandler := tableReservationHttp.NewTableReservationHandler(tableReservationUsecase)
	tableReservationHttp.MapTableReservationRoutes(tableReservationGroup, tableReservationHandler)

	// Payment Group
	paymentRepository := paymentRepository.NewPaymentRepository(s.db)
	paymentUsecase := paymentUsecase.NewPaymentUsecase(paymentRepository)
	paymentHandler := paymentHttp.NewPaymentHandler(paymentUsecase)
	paymentHttp.MapPaymentRoutes(paymentGroup, paymentHandler)

	// Food Order Group
	foodOrderRepository := foodOrderRepository.NewFoodOrderRepository(s.db)
	foodOrderUsecase := foodOrderUsecase.NewFoodOrderUsecase(foodOrderRepository)
	foodOrderHandler := foodOrderHttp.NewFoodOrderHandler(foodOrderUsecase)
	foodOrderHttp.MapFoodOrderRoutes(foodOrderGroup, foodOrderHandler)

	// Notification Group
	notiRepository := notiRepository.NewNotiRepository(s.db)
	notiUsecase := notiUsecase.NewNotiUsecase(notiRepository)
	notificationHandler := notiHttp.NewNotiHandler(notiUsecase)
	notiHttp.MapNotiRoutes(notificationGroup, notificationHandler)

	return nil
}