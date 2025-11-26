package config

// Request Flow Link:
// The init* functions in this package are called from main.go during startup to create
// each database connection that the downstream repositories (and thus every request) depend on.

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// const (
// 	defaultUser     = "postgres"
// 	defaultPassword = "mraffa0217"
// 	defaultHost     = "localhost" // Changed from "db" to "localhost"
// 	defaultPort     = "5432"
// 	defaultSSLMode  = "disable"
// 	golangDB        = "golang"
// 	trafficDB       = "traffic_ticket"
// )

// InitDB initializes database connection with specified database name and optimized settings
func InitDB(dbName string) *sql.DB {
	user := getenv("DB_USER", "")
	pass := getenv("DB_PASSWORD", "")
	host := getenv("DB_HOST", "")
	port := getenv("DB_PORT", "")
	ssl := getenv("DB_SSLMODE", "")

	// Auto-detect Docker environment and use host.docker.internal
	if host == "localhost" && isRunningInDocker() {
		host = "host.docker.internal"
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, pass, dbName, ssl)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}

	configureConnectionPool(db)

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to connect to postgres DB:", err)
	}

	log.Printf("Database connection established for %s with optimized pool settings", dbName)
	return db

}

// initMySQLDB initializes a MySQL database connection using the provided prefix and display name
func initMySQLDB(prefix, displayName string) *sql.DB {
	host := getenv(prefix+"_HOST", "")
	port := getenv(prefix+"_PORT", "")
	user := getenv(prefix+"_USER", "")
	password := getenv(prefix+"_PASSWORD", "")
	database := getenv(prefix+"_DATABASE", "")

	// Auto-detect Docker environment and use host.docker.internal
	if host == "localhost" && isRunningInDocker() {
		host = "host.docker.internal"
	}

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, database)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatalf("Failed to open %s DB: %v", displayName, err)
	}

	configureConnectionPool(db)

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to connect to sql %s DB: %v", displayName, err)
	}

	log.Printf("%s connection established for %s with optimized pool settings", displayName, database)
	return db
}

// InitGolangDB initializes connection to golang database (for user operations)
func InitGolangDB() *sql.DB {
	return InitDB("golang")
}

// InitTrafficDB initializes connection to traffic_ticket database (for traffic ticket operations)
func InitTrafficDB() *sql.DB {
	return InitDB("traffic_ticket")
}

// configureConnectionPool sets up shared connection pool settings
func configureConnectionPool(db *sql.DB) {
	db.SetMaxOpenConns(25)                 // Maximum number of open connections
	db.SetMaxIdleConns(10)                 // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Maximum lifetime of a connection
	db.SetConnMaxIdleTime(1 * time.Minute) // Maximum idle time of a connection
	// Note: Go's database/sql doesn't have a direct "wait timeout" for connections
	// Use context.WithTimeout in queries to add query-level timeouts
	// HTTP server timeouts (in main.go) will cancel requests that wait too long
}

// GetQueryTimeout returns the query timeout duration from environment variable
// This timeout applies to individual database queries to prevent slow queries from blocking connections
func GetQueryTimeout() time.Duration {
	timeoutSeconds := getenvInt("DB_QUERY_TIMEOUT_SECONDS", 10) // Default: 10 seconds
	return time.Duration(timeoutSeconds) * time.Second
}

// getenvInt retrieves integer environment variable with fallback
func getenvInt(key string, defaultValue int) int {
	value := getenv(key, "")
	if value == "" {
		return defaultValue
	}
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	return defaultValue
}

// InitMySQLDB initializes connection to MySQL database
func InitMySQLDB() *sql.DB {
	return initMySQLDB("MYSQL", "MySQL")
}
func InitTerminalDB() *sql.DB {
	return initMySQLDB("LAUT_MYSQL", "terminal")
}

// InitPassengerPlaneDB initializes connection to passenger MySQL database (no password)
func InitPassengerPlaneDB() *sql.DB {
	return initMySQLDB("PASSENGER_MYSQL", "Passenger Plane MySQL")
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

// isRunningInDocker detects if the application is running inside a Docker container
func isRunningInDocker() bool {
	// Check for Docker-specific files that exist in containers
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	// Check for Docker in cgroup
	if data, err := os.ReadFile("/proc/1/cgroup"); err == nil {
		return strings.Contains(string(data), "docker")
	}

	return false
}
