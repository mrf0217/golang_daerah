package main

import (
	"golang_daerah/internal/database"
	"golang_daerah/internal/handler"
	"golang_daerah/internal/service"
	"golang_daerah/pkg/jwtutil"
	"golang_daerah/pkg/middleware"
	"log"
	"net/http"
)

func main() {

	allDBs := database.InitAllDatabases()
	defer database.CloseAllDatabases(allDBs)

	// Initialize SQLX databases
	// trafficDB := config.InitTrafficDBX()
	// defer trafficDB.Close()

	// mysqlDB := config.InitMySQLDBX()
	// defer mysqlDB.Close()

	// passengerDB := config.InitPassengerPlaneDBX()
	// defer passengerDB.Close()

	// terminalDB := config.InitTerminalDBX()
	// defer terminalDB.Close()

	// authDB := config.InitGolangDBX() // NEW: Dedicated auth database
	// defer authDB.Close()

	// passangerlocalDB := config.InitMySQLDBX_passanger()
	// defer passangerlocalDB.Close()

	lautBase := &database.BaseMultiDBRepository{Dbs: allDBs}
	passengerBase := &database.BaseMultiDBRepository{Dbs: allDBs}
	trafficBase := &database.BaseMultiDBRepository{Dbs: allDBs}
	mysqlTrafficBase := &database.BaseMultiDBRepository{Dbs: allDBs}

	// Initialize repositories
	// trafficHandler := httpDelivery.NewPostgresTrafficTicketSQLXRepository()
	// mysqlHandler := httpDelivery.NewMySQLTrafficTicketSQLXRepository()
	// passengerHandler := httpDelivery.NewPassengerPlaneSQLXRepository()
	// lautHandler := httpDelivery.NewLautSQLXRepository()
	// userRepo := httpDelivery.NewUserRepository() // NEW: User repository using sql.DB

	// Auth service (manages its own DBs internally)
	userRepo := service.NewUserRepository()
	authService := service.NewAuthService(userRepo)
	// Initialize services
	lautService := service.NewLautService(lautBase)
	passengerService := service.NewPassengerPlaneService(passengerBase)
	trafficService := service.NewPostgresTrafficTicketSQLXRepository(trafficBase)
	mysqlTrafficService := service.NewMySQLTrafficTicketService(mysqlTrafficBase)
	// userService := httpDelivery.NewUserService(userRepo)
	// authHandler := httpDelivery.NewUserHandler(userService)
	// Initialize handlers
	lautHandler := handler.NewLautHandler(lautService)
	passengerHandler := handler.NewPassengerHandler(passengerService)
	trafficHandler := handler.NewTrafficHandler(trafficService)
	mysqlTrafficHandler := handler.NewTrafficMySQLHandler(mysqlTrafficService)
	authHandler := handler.NewAuthHandler(authService)
	// trafficHandler := httpDelivery.NewTrafficTicketSQLXHandler(trafficRepo)
	// mysqlHandler := httpDelivery.NewMySQLTrafficTicketSQLXHandler(mysqlRepo)
	// passengerHandler := httpDelivery.NewPassengerPlaneSQLXHandler(passengerRepo)
	// lautHandler := httpDelivery.NewLautSQLXHandler(lautRepo)
	// authHandler := httpDelivery.NewUserHandler(userService)

	// Setup router
	router := http.NewServeMux()
	// rateLimit := middleware.RateLimitMiddleware(100, 10)
	// protected := func(handler http.HandlerFunc) http.HandlerFunc {
	// 	return rateLimit(jwtutil.AuthMiddleware(handler))
	// }

	// Register routes
	router.HandleFunc("/api/traffic_tickets/postgres",
		(jwtutil.AuthMiddleware(trafficHandler.GetPaginated)))
	router.HandleFunc("/api/traffic_tickets/postgres_create",
		jwtutil.AuthMiddleware(middleware.RateLimitMiddleware(100, 10)(trafficHandler.Create)))

	router.HandleFunc("/api/traffic_tickets/mysql",
		middleware.RateLimitMiddleware(100, 10)(jwtutil.AuthMiddleware(mysqlTrafficHandler.GetPaginated)))
	router.HandleFunc("/api/traffic_tickets/mysql_create",
		middleware.RateLimitMiddleware(100, 10)(jwtutil.AuthMiddleware(mysqlTrafficHandler.Create)))

	router.HandleFunc("/api/passengers",
		middleware.RateLimitMiddleware(100, 10)(jwtutil.AuthMiddleware(passengerHandler.GetPaginated)))
	router.HandleFunc("/api/passengers/create",
		middleware.RateLimitMiddleware(100, 10)(jwtutil.AuthMiddleware(passengerHandler.Create)))

	router.HandleFunc("/api/terminals",
		middleware.RateLimitMiddleware(100, 10)(jwtutil.AuthMiddleware(lautHandler.GetPaginated)))
	router.HandleFunc("/api/terminals/create",
		middleware.RateLimitMiddleware(100, 10)(jwtutil.AuthMiddleware(lautHandler.Create)))
	router.HandleFunc("/api/terminals/showall",
		middleware.RateLimitMiddleware(100, 10)(jwtutil.AuthMiddleware(lautHandler.LautGetCompleteDataHandler)))

	router.HandleFunc("/api/register",
		middleware.RateLimitMiddleware(100, 10)(authHandler.Register))
	router.HandleFunc("/api/login",
		middleware.RateLimitMiddleware(100, 10)(authHandler.Login))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
