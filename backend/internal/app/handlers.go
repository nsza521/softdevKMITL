package app

import (
	userHttp "backend/internal/user/delivery/http"
	userRepository "backend/internal/user/repository"
	userUsecase "backend/internal/user/usecase"

	customerHttp "backend/internal/customer/delivery/http"
	customerRepository "backend/internal/customer/repository"
	customerUsecase "backend/internal/customer/usecase"

	restaurantHttp "backend/internal/restaurant/delivery/http"
	restaurantRepository "backend/internal/restaurant/repository"
	restaurantUsecase "backend/internal/restaurant/usecase"

	tableHttp "backend/internal/table/delivery/http"
	tableRepository "backend/internal/table/repository"
	tableUsecase "backend/internal/table/usecase"

	tableReservationHttp "backend/internal/reservation/delivery/http"
	tableReservationRepository "backend/internal/reservation/repository"
	tableReservationUsecase "backend/internal/reservation/usecase"

	paymentHttp "backend/internal/payment/delivery/http"
	paymentRepository "backend/internal/payment/repository"
	paymentUsecase "backend/internal/payment/usecase"

	foodOrderHttp "backend/internal/order/delivery/http"
	foodOrderRepository "backend/internal/order/repository"
	foodOrderUsecase "backend/internal/order/usecase"
	foodOrderAdapter "backend/internal/order/adapter"

	notiHttp "backend/internal/notifications/delivery/http"
	notiRepository "backend/internal/notifications/repository"
	notiUsecase "backend/internal/notifications/usecase"

	menuHttp "backend/internal/menu/delivery/http"
	menuRepo "backend/internal/menu/repository"
	menuUC "backend/internal/menu/usecase"
)

func (s *App) MapHandlers() error {

	userGroup := s.gin.Group("/user")
	customerGroup := s.gin.Group("/customer")
	restaurantGroup := s.gin.Group("/restaurant")
	menuGroup := s.gin.Group("/restaurant/menu")
	tableGroup := s.gin.Group("/table")
	tableReservationGroup := s.gin.Group("/table/reservation")
	foodOrderGroup := s.gin.Group("/restaurant/order")
	notificationGroup := s.gin.Group("/notification")
	paymentGroup := s.gin.Group("/payment")

	// Customer Group
	customerRepository := customerRepository.NewCustomerRepository(s.db)
	customerUsecase := customerUsecase.NewCustomerUsecase(customerRepository)
	customerHandler := customerHttp.NewCustomerHandler(customerUsecase)
	customerHttp.MapCustomerRoutes(customerGroup, customerHandler)

	// --- MenuItem (ของเดิม) ---
	mRepo := menuRepo.NewMenuRepository(s.db)
	mUC := menuUC.NewMenuUsecase(mRepo, s.minio)
	mH := menuHttp.NewMenuHandler(mUC)
	menuHttp.MapMenuRoutes(menuGroup, mH)

	// --- MenuType (ของใหม่) ---
	mtRepo := menuRepo.NewMenuTypeRepository(s.db)
	mtUC := menuUC.NewMenuTypeUsecase(mtRepo)
	mtH := menuHttp.NewMenuTypeHandler(mtUC)
	menuHttp.MapMenuTypeRoutes(menuGroup, mtH)

	// --- AddOn (Group + Option) ---
	addonRepo := menuRepo.NewAddOnRepository(s.db)
	addonUC := menuUC.NewAddOnUsecase(addonRepo)
	addonH := menuHttp.NewAddOnHandler(addonUC)
	menuHttp.MapAddOnRoutes(menuGroup, addonH)

	// Restaurant Group
	restaurantRepository := restaurantRepository.NewRestaurantRepository(s.db)
	restaurantUsecase := restaurantUsecase.NewRestaurantUsecase(restaurantRepository, mRepo, s.minio)
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
	foodOrderRepository := foodOrderRepository.NewOrderRepository(s.db)
	menuRead := foodOrderAdapter.NewMenuReadAdapter(mUC)
	foodOrderUsecase := foodOrderUsecase.NewOrderUsecase(foodOrderRepository, menuRead)
	foodOrderHandler := foodOrderHttp.NewOrderHandler(foodOrderUsecase)
	foodOrderHttp.MapFoodOrderRoutes(foodOrderGroup, foodOrderHandler)

	// Notification Group
	notiRepository := notiRepository.NewNotiRepository(s.db)
	notiUsecase := notiUsecase.NewNotiUsecase(notiRepository)
	notificationHandler := notiHttp.NewNotiHandler(notiUsecase)
	notiHttp.MapNotiRoutes(notificationGroup, notificationHandler)

	return nil
}
