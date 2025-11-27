package main

import (
    "log"
    "net/http"
    
    "golang_daerah/config"
    httpDelivery "golang_daerah/internal/delivery/http"
    "golang_daerah/internal/repository"
    "golang_daerah/pkg/middleware"
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

    // Initialize repositories
    trafficRepo := repository.NewPostgresTrafficTicketSQLXRepository(trafficDB)
    mysqlRepo := repository.NewMySQLTrafficTicketSQLXRepository(mysqlDB)
    passengerRepo := repository.NewPassengerPlaneSQLXRepository(passengerDB)
    lautRepo := repository.NewLautSQLXRepository(terminalDB)

    // Initialize handlers
    trafficHandler := httpDelivery.NewTrafficTicketSQLXHandler(trafficRepo)
    mysqlHandler := httpDelivery.NewMySQLTrafficTicketSQLXHandler(mysqlRepo)
    passengerHandler := httpDelivery.NewPassengerPlaneSQLXHandler(passengerRepo)
    lautHandler := httpDelivery.NewLautSQLXHandler(lautRepo)

    // Setup router
    router := http.NewServeMux()
    rateLimit := middleware.RateLimitMiddleware(100, 10)

    // Register routes
    router.HandleFunc("/api/traffic_tickets/postgres", rateLimit(trafficHandler.GetPaginated))
    router.HandleFunc("/api/traffic_tickets/postgres/create", rateLimit(trafficHandler.Create))
    
    router.HandleFunc("/api/traffic_tickets/mysql", rateLimit(mysqlHandler.GetPaginated))
    router.HandleFunc("/api/traffic_tickets/mysql/create", rateLimit(mysqlHandler.Create))
    
    router.HandleFunc("/api/passengers", rateLimit(passengerHandler.GetPaginated))
    router.HandleFunc("/api/passengers/create", rateLimit(passengerHandler.Create))
    
    router.HandleFunc("/api/terminals", rateLimit(lautHandler.GetPaginated))
    router.HandleFunc("/api/terminals/create", rateLimit(lautHandler.Create))

    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}