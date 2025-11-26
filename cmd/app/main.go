package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	config "golang_daerah/config"
	httpDelivery "golang_daerah/internal/delivery/http"
	"golang_daerah/internal/repository"
	"golang_daerah/internal/usecases"
	"golang_daerah/pkg/jwtutil"
	"golang_daerah/pkg/middleware"
)

// ============================================================================
// RouteRegistrar: Helper for Consistent Route Registration
// ============================================================================
// This struct helps register routes with consistent middleware patterns.
// Instead of repeating the middleware chain everywhere, we use helper methods.
// ============================================================================

// RouteRegistrar helps register routes with consistent middleware
type RouteRegistrar struct {
	router    *http.ServeMux                          // HTTP router that matches requests to handlers
	rateLimit func(http.HandlerFunc) http.HandlerFunc // Rate limiting middleware function
}

// Protected registers a route with rate limiting AND JWT authentication
// Request Flow: HTTP Request → Rate Limit Middleware → JWT Auth Middleware → Handler
// Users must provide valid JWT token in "Authorization: Bearer <token>" header
func (rr *RouteRegistrar) Protected(path string, handler http.HandlerFunc) {
	// Middleware chain (executed from outside to inside):
	// 1. Rate Limit checks request rate per IP
	// 2. JWT Auth validates the token
	// 3. Handler processes the business logic
	rr.router.HandleFunc(path, rr.rateLimit(jwtutil.AuthMiddleware(handler)))
}

// Public registers a route with ONLY rate limiting (no JWT authentication)
// Request Flow: HTTP Request → Rate Limit Middleware → Handler
// No authentication required, but still rate limited to prevent abuse
func (rr *RouteRegistrar) Public(path string, handler http.HandlerFunc) {
	// Only rate limiting, no authentication
	rr.router.HandleFunc(path, rr.rateLimit(handler))
}

// ============================================================================
// initDB: Database Initialization Wrapper
// ============================================================================
// Wraps database initialization with consistent logging.
// This allows us to add retry logic, metrics, or other enhancements in one place.
// ============================================================================

// initDB initializes a database connection with consistent error handling and logging
// Parameters:
//   - name: Display name for logging (e.g., "golang", "MySQL")
//   - initFunc: Function that actually initializes the database connection
//
// Returns: *sql.DB connection object
func initDB(name string, initFunc func() *sql.DB) *sql.DB {
	db := initFunc() // Calls the actual database initialization function
	log.Printf("✓ %s database initialized", name)
	return db
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
	// ============================================================================
	// PHASE 1: DATABASE CONNECTION INITIALIZATION
	// ============================================================================
	// Initialize all database connections that the application will use.
	// Each database connection is set up with connection pooling and timeouts.
	// The defer statements ensure connections are properly closed when app exits.
	// ============================================================================

	// PostgreSQL Database: "golang" - Used for user authentication and management
	golangDB := initDB("golang", config.InitGolangDB)
	defer golangDB.Close()

	// PostgreSQL Database: "traffic_ticket" - Used for traffic ticket operations
	trafficDB := initDB("traffic_ticket", config.InitTrafficDB)
	defer trafficDB.Close()

	// MySQL Database - Used for traffic ticket operations (alternative to PostgreSQL)
	mysqlDB := initDB("MySQL", config.InitMySQLDB)
	defer mysqlDB.Close()

	// MySQL Database: Passenger Plane - Used for passenger plane operations
	passengerDB := initDB("Passenger Plane MySQL", config.InitPassengerPlaneDB)
	defer passengerDB.Close()

	// Initialize terminal database connection
	terminalDB := initDB("terminal", config.InitTerminalDB)
	defer terminalDB.Close()

	// ============================================================================
	// PHASE 2: APPLICATION LAYER SETUP (Dependency Injection)
	// ============================================================================
	// Build the application layers from bottom to top:
	// Database → Repository → Service → Handler
	// This follows Clean Architecture pattern for separation of concerns.
	// ============================================================================

	// --- User Layer (uses golang database) ---
	// Request Flow: Handler → Service → Repository → Database
	userRepo := repository.NewUserRepository(golangDB)      // Repository: Handles database operations
	userService := usecases.NewUserService(userRepo)        // Service: Contains business logic (password hashing, validation)
	userHandler := httpDelivery.NewUserHandler(userService) // Handler: Handles HTTP requests/responses

	// --- Traffic Ticket Layer (uses traffic_ticket PostgreSQL database) ---
	// Request Flow: Handler → Repository → Database
	trafficTicketRepo := repository.NewTrafficTicketRepository(trafficDB)
	trafficHandler := httpDelivery.NewTrafficTicketHandler(trafficTicketRepo)

	// --- MySQL Traffic Ticket Layer (uses MySQL database) ---
	// Request Flow: Handler → Repository → Database
	mysqlRepo := repository.NewMySQLTrafficTicketRepository(mysqlDB)
	mysqlHandler := httpDelivery.NewMySQLTrafficTicketHandler(mysqlRepo)

	// --- Passenger Plane Layer (uses passenger MySQL database) ---
	// Request Flow: Handler → Repository → Database
	passengerRepo := repository.NewPassengerPlaneRepository(passengerDB)
	passengerHandler := httpDelivery.NewPassengerPlaneHandler(passengerRepo)

	// --- Laut (Port/Terminal) Layer (uses terminal PostgreSQL database) ---
	// Request Flow: Handler → Repository → Database
	lautRepo := repository.NewLautRepository(terminalDB)
	lautHandler := httpDelivery.NewLautHandler(lautRepo)

	// ============================================================================
	// PHASE 3: HTTP ROUTER SETUP
	// ============================================================================
	// Create the HTTP router that will match incoming requests to handlers.
	// ============================================================================

	router := http.NewServeMux()

	// ============================================================================
	// PHASE 4: MIDDLEWARE CONFIGURATION
	// ============================================================================
	// Configure rate limiting middleware to protect against abuse.
	// Rate limiting checks requests per IP address before allowing them through.
	// ============================================================================

	// Get rate limit settings from environment variables (with defaults)
	rateLimitRequests := getEnvInt("RATE_LIMIT_REQUESTS", 100) // Max 100 requests per minute
	rateLimitBurst := getEnvInt("RATE_LIMIT_BURST", 10)        // Allow burst of 10 requests
	rateLimit := middleware.RateLimitMiddleware(rateLimitRequests, rateLimitBurst)

	// Create route registrar for consistent route registration
	// This helper ensures all routes use the same middleware pattern
	routes := &RouteRegistrar{router: router, rateLimit: rateLimit}

	// ============================================================================
	// PHASE 5: ROUTE REGISTRATION
	// ============================================================================
	// Register all API endpoints with their appropriate middleware.
	// Request Flow: HTTP Request → Router → Middleware → Handler → Service/Repository → Database
	// ============================================================================

	// --- Public Endpoints (No Authentication Required) ---
	// These endpoints only use rate limiting, no JWT authentication.
	// Request Flow: HTTP Request → Router → Rate Limit Middleware → Handler
	routes.Public("/api/register", userHandler.Register) // POST: User registration
	routes.Public("/api/login", userHandler.Login)       // POST: User login (returns JWT token)

	// --- Protected Endpoints (Authentication Required) ---
	// These endpoints use both rate limiting AND JWT authentication.
	// Request Flow: HTTP Request → Router → Rate Limit Middleware → JWT Auth Middleware → Handler
	// Users must include "Authorization: Bearer <token>" header to access these endpoints.

	// PostgreSQL Traffic Ticket Endpoints
	routes.Protected("/api/traffic_tickets/inputpostgre", trafficHandler.Create)      // POST: Create traffic ticket in PostgreSQL
	routes.Protected("/api/traffic_tickets/listpostgre", trafficHandler.GetPaginated) // GET: List traffic tickets from PostgreSQL (paginated)

	// MySQL Traffic Ticket Endpoints
	routes.Protected("/api/traffic_tickets/inputsql", mysqlHandler.Create)      // POST: Create traffic ticket in MySQL
	routes.Protected("/api/traffic_tickets/listsql", mysqlHandler.GetPaginated) // GET: List traffic tickets from MySQL (paginated)

	// Passenger Plane MySQL Endpoints
	routes.Protected("/api/traffic_tickets/inputpassenger", passengerHandler.Create)      // POST: Create passenger plane record
	routes.Protected("/api/traffic_tickets/listpassenger", passengerHandler.GetPaginated) // GET: List passenger plane records (paginated)

	// Terminal/Port Endpoints
	routes.Protected("/api/terminal/tambah_Laut", lautHandler.Create) // POST: Add port/terminal data
	routes.Protected("/api/terminal/laut", lautHandler.GetPaginated)  // GET: List port/terminal data (paginated)

	// ============================================================================
	// PHASE 6: HTTP SERVER CONFIGURATION
	// ============================================================================
	// Configure HTTP server with timeouts to prevent requests from hanging indefinitely.
	// This protects the server from resource exhaustion when many users are waiting.
	// ============================================================================

	// Get port from environment variable (default: 8080)
	port := getEnv("APP_PORT", "8080")

	// Get timeout settings from environment variables (with defaults)
	// These timeouts ensure requests don't wait forever for database connections
	readTimeout := getEnvInt("HTTP_READ_TIMEOUT_SECONDS", 15)   // Max 15 seconds to read request
	writeTimeout := getEnvInt("HTTP_WRITE_TIMEOUT_SECONDS", 15) // Max 15 seconds to write response
	idleTimeout := getEnvInt("HTTP_IDLE_TIMEOUT_SECONDS", 60)   // Max 60 seconds idle on keep-alive connections

	// Configure HTTP server with timeouts
	// If a request takes longer than these timeouts, it will be cancelled automatically
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,                                    // Routes requests to registered handlers
		ReadTimeout:  time.Duration(readTimeout) * time.Second,  // Cancel request if reading takes too long
		WriteTimeout: time.Duration(writeTimeout) * time.Second, // Cancel request if writing takes too long
		IdleTimeout:  time.Duration(idleTimeout) * time.Second,  // Close idle keep-alive connections
	}

	// ============================================================================
	// PHASE 7: SERVER STARTUP
	// ============================================================================
	// Start the HTTP server and begin accepting incoming requests.
	// The server will block here until it's stopped or encounters an error.
	// ============================================================================

	fmt.Printf("Server listening on port %s with timeouts (Read: %ds, Write: %ds, Idle: %ds)...\n",
		port, readTimeout, writeTimeout, idleTimeout)
	log.Fatal(server.ListenAndServe()) // Blocks here, serving requests until error or shutdown
}

// ============================================================================
// Helper Functions for Environment Variables
// ============================================================================
// These functions provide safe access to environment variables with fallback defaults.
// This ensures the application always has valid configuration values.
// ============================================================================

// getEnv retrieves a string environment variable with a fallback default value
// Parameters:
//   - key: Environment variable name (e.g., "APP_PORT")
//   - defaultValue: Value to return if environment variable is not set or empty
//
// Returns: Environment variable value or default value
// Example: getEnv("APP_PORT", "8080") returns "8080" if APP_PORT is not set
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt retrieves an integer environment variable with a fallback default value
// Parameters:
//   - key: Environment variable name (e.g., "RATE_LIMIT_REQUESTS")
//   - defaultValue: Integer value to return if environment variable is not set, empty, or invalid
//
// Returns: Integer value from environment or default value
// Example: getEnvInt("RATE_LIMIT_REQUESTS", 100) returns 100 if variable is not set
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
