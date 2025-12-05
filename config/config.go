// config/config.go
package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// InitGolangDBX initializes PostgreSQL connection for golang database with sqlx
func InitGolangDBX() *sqlx.DB {
	user := getenv("DB_USER", "postgres")
	pass := getenv("DB_PASSWORD", "mraffa0217")
	host := getenv("DB_HOST", "localhost")
	port := getenv("DB_PORT", "5432")
	ssl := getenv("DB_SSLMODE", "disable")

	if host == "localhost" && isRunningInDocker() {
		host = "host.docker.internal"
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=golang sslmode=%s",
		host, port, user, pass, ssl)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to golang DB:", err)
	}

	configureConnectionPool(db.DB)
	log.Printf("Database connection established for golang with optimized pool settings")

	return db
}

// InitTrafficDBX initializes PostgreSQL connection for traffic_ticket database with sqlx
func InitTrafficDBX() *sqlx.DB {
	user := getenv("DB_USER", "postgres")
	pass := getenv("DB_PASSWORD", "mraffa0217")
	host := getenv("DB_HOST", "localhost")
	port := getenv("DB_PORT", "5432")
	ssl := getenv("DB_SSLMODE", "disable")

	if host == "localhost" && isRunningInDocker() {
		host = "host.docker.internal"
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=traffic_ticket sslmode=%s",
		host, port, user, pass, ssl)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to traffic_ticket DB:", err)
	}

	configureConnectionPool(db.DB)
	log.Printf("Database connection established for traffic_ticket with optimized pool settings")

	return db
}

func InitMySQLDBX_passanger() *sqlx.DB {
	host := getenv("MYSQL_HOST", "localhost")
	port := getenv("MYSQL_PORT", "3306")
	user := getenv("MYSQL_USER", "root")
	password := getenv("MYSQL_PASSWORD", "")
	database := getenv("MYSQL_DATABASE", "passanger")

	if host == "localhost" && isRunningInDocker() {
		host = "host.docker.internal"
	}

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, database)

	db, err := sqlx.Connect("mysql", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL DB: %v", err)
	}

	configureConnectionPool(db.DB)
	log.Printf("MySQL connection established for %s with optimized pool settings", database)

	return db
}

// InitMySQLDBX initializes MySQL connection for traffic_tickets database with sqlx
func InitMySQLDBX() *sqlx.DB {
	host := getenv("MYSQL_HOST", "localhost")
	port := getenv("MYSQL_PORT", "3306")
	user := getenv("MYSQL_USER", "root")
	password := getenv("MYSQL_PASSWORD", "")
	database := getenv("MYSQL_DATABASE", "traffic_ticket")

	if host == "localhost" && isRunningInDocker() {
		host = "host.docker.internal"
	}

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, database)

	db, err := sqlx.Connect("mysql", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL DB: %v", err)
	}

	configureConnectionPool(db.DB)
	log.Printf("MySQL connection established for %s with optimized pool settings", database)

	return db
}

// InitPassengerPlaneDBX initializes MySQL connection for passenger database with sqlx
func InitPassengerPlaneDBX() *sqlx.DB {
	host := getenv("PASSENGER_MYSQL_HOST", "localhost")
	port := getenv("PASSENGER_MYSQL_PORT", "3307")
	user := getenv("PASSENGER_MYSQL_USER", "root")
	password := getenv("PASSENGER_MYSQL_PASSWORD", "")
	database := getenv("PASSENGER_MYSQL_DATABASE", "passenger")

	if host == "localhost" && isRunningInDocker() {
		host = "host.docker.internal"
	}

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, database)

	db, err := sqlx.Connect("mysql", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to Passenger Plane MySQL DB: %v", err)
	}

	configureConnectionPool(db.DB)
	log.Printf("Passenger Plane MySQL connection established for %s with optimized pool settings", database)

	return db
}

// InitTerminalDBX initializes MySQL connection for terminal database with sqlx
func InitTerminalDBX() *sqlx.DB {
	host := getenv("LAUT_MYSQL_HOST", "localhost")
	port := getenv("LAUT_MYSQL_PORT", "3306")
	user := getenv("LAUT_MYSQL_USER", "root")
	password := getenv("LAUT_MYSQL_PASSWORD", "")
	database := getenv("LAUT_MYSQL_DATABASE", "terminal")

	if host == "localhost" && isRunningInDocker() {
		host = "host.docker.internal"
	}

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, database)

	db, err := sqlx.Connect("mysql", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to terminal MySQL DB: %v", err)
	}

	configureConnectionPool(db.DB)
	log.Printf("Terminal MySQL connection established for %s with optimized pool settings", database)

	return db
}

// InitAuthDBX initializes MySQL connection for auth database with sqlx
func InitAuthDBX() *sqlx.DB {
	host := getenv("AUTH_MYSQL_HOST", "localhost")
	port := getenv("AUTH_MYSQL_PORT", "3306")
	user := getenv("AUTH_MYSQL_USER", "root")
	password := getenv("AUTH_MYSQL_PASSWORD", "")
	database := getenv("AUTH_MYSQL_DATABASE", "golang")

	if host == "localhost" && isRunningInDocker() {
		host = "host.docker.internal"
	}

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, database)

	db, err := sqlx.Connect("mysql", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to Auth MySQL DB: %v", err)
	}

	configureConnectionPool(db.DB)
	log.Printf("Auth MySQL connection established for %s with optimized pool settings", database)

	return db
}

// configureConnectionPool sets up shared connection pool settings
func configureConnectionPool(db *sql.DB) {
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)
}

// GetQueryTimeout returns the query timeout duration from environment variable
func GetQueryTimeout() time.Duration {
	timeoutSeconds := getenvInt("DB_QUERY_TIMEOUT_SECONDS", 10)
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

// getenv retrieves string environment variable with fallback
func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

// isRunningInDocker detects if the application is running inside a Docker container
func isRunningInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	if data, err := os.ReadFile("/proc/1/cgroup"); err == nil {
		return strings.Contains(string(data), "docker")
	}

	return false
}
