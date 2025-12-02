package main

import (
	"golang_daerah/config"
	httpDelivery "golang_daerah/internal/delivery/http"

	"golang_daerah/pkg/jwtutil"
	"golang_daerah/pkg/middleware"
	"log"
	"net/http"
)

func main() {
	// Initialize SQLX databases
	trafficDB := config.InitTrafficDBX()
	defer trafficDB.Close()

	mysqlDB := config.InitMySQLDBX()
	defer mysqlDB.Close()

	passengerDB := config.InitPassengerPlaneDBX()
	defer passengerDB.Close()

	terminalDB := config.InitTerminalDBX()
	defer terminalDB.Close()

	authDB := config.InitGolangDBX() // NEW: Dedicated auth database
	defer authDB.Close()

	passangerlocalDB := config.InitMySQLDBX_passanger()
	defer passangerlocalDB.Close()

	// Initialize repositories
	trafficHandler := httpDelivery.NewPostgresTrafficTicketSQLXRepository(trafficDB)
	mysqlHandler := httpDelivery.NewMySQLTrafficTicketSQLXRepository(mysqlDB)
	passengerHandler := httpDelivery.NewPassengerPlaneSQLXRepository(passengerDB)
	lautHandler := httpDelivery.NewLautRepository(terminalDB)
	userRepo := httpDelivery.NewUserRepository(authDB.DB) // NEW: User repository using sql.DB

	// Initialize services
	userService := httpDelivery.NewUserService(userRepo)
	authHandler := httpDelivery.NewUserHandler(userService)
	// Initialize handlers
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
		middleware.RateLimitMiddleware(100, 10)(jwtutil.AuthMiddleware(mysqlHandler.GetPaginated)))
	router.HandleFunc("/api/traffic_tickets/mysql_create",
		middleware.RateLimitMiddleware(100, 10)(jwtutil.AuthMiddleware(mysqlHandler.Create)))

	router.HandleFunc("/api/passengers",
		middleware.RateLimitMiddleware(100, 10)(jwtutil.AuthMiddleware(passengerHandler.GetPaginated)))
	router.HandleFunc("/api/passengers/create",
		middleware.RateLimitMiddleware(100, 10)(jwtutil.AuthMiddleware(passengerHandler.Create)))

	router.HandleFunc("/api/terminals",
		middleware.RateLimitMiddleware(100, 10)(jwtutil.AuthMiddleware(lautHandler.GetPaginated)))
	router.HandleFunc("/api/terminals/create",
		middleware.RateLimitMiddleware(100, 10)(jwtutil.AuthMiddleware(lautHandler.Create)))
	router.HandleFunc("/api/terminals/showall",
		middleware.RateLimitMiddleware(100, 10)(jwtutil.AuthMiddleware(lautHandler.GetCompleteDataHandler)))

	router.HandleFunc("/api/register",
		middleware.RateLimitMiddleware(100, 10)(authHandler.Register))
	router.HandleFunc("/api/login",
		middleware.RateLimitMiddleware(100, 10)(authHandler.Login))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
