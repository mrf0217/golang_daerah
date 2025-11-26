# Golang Daerah API - Complete Function Call Chain Documentation

This document provides a comprehensive trace of **every function call** from `cmd/app/main.go` all the way through to repositories, showing the complete branching paths, file locations, line numbers, and explanations of how and why functions are connected.

## Table of Contents
1. [Overview](#overview)
2. [Complete Function Call Chains](#complete-function-call-chains)
3. [Request Flow Chains](#request-flow-chains)
4. [Environment Variables](#environment-variables)
5. [Step-by-Step Guide: Creating New API Endpoints](#step-by-step-guide-creating-new-api-endpoints)

---

## Overview

This is a Go-based REST API application that manages traffic tickets across multiple databases (PostgreSQL and MySQL). The application follows a clean architecture pattern with separate layers for delivery (HTTP handlers), use cases (business logic), and repositories (data access).

**Key Features:**
- User authentication (register/login) with JWT tokens
- Traffic ticket management across multiple databases
- Rate limiting middleware
- JWT authentication middleware
- Support for PostgreSQL and MySQL databases
- HTTP server timeouts (prevents hanging requests)
- Database query timeouts (prevents slow queries from blocking connections)
- Consolidated pagination validation
- Consolidated query error handling

---

## Complete Function Call Chains

This section traces **every function call** from `main.go` through all layers, showing complete branching paths.

### 1. Application Entry Point

**File:** `cmd/app/main.go`  
**Function:** `main()` (Line 42)  
**Purpose:** Entry point that orchestrates all initialization and setup

---

### 2. Database Initialization Chains

#### 2.1 Golang Database (PostgreSQL) Initialization Chain

**Starting Point:** `cmd/app/main.go:44`
```go
golangDB := initDB("golang", config.InitGolangDB)
```

**Complete Call Chain:**

1. **`initDB()`** - `cmd/app/main.go:36`
   - **Purpose:** Wrapper function for consistent database initialization logging
   - **Parameters:** 
     - `name`: "golang" (display name for logging)
     - `initFunc`: `config.InitGolangDB` (function pointer)
   - **Calls:** `initFunc()` → which calls `config.InitGolangDB()`
   - **Returns:** `*sql.DB` connection object
   - **Why:** Centralizes logging and allows future enhancements (retry logic, metrics, etc.)

2. **`config.InitGolangDB()`** - `config/config.go:87`
   - **Purpose:** Initializes connection to PostgreSQL "golang" database
   - **Calls:** `InitDB("golang")` → `config/config.go:26`
   - **Returns:** `*sql.DB` connection object
   - **Why:** Provides a named function for the golang database, making code more readable

3. **`InitDB(dbName string)`** - `config/config.go:26`
   - **Purpose:** Core PostgreSQL database initialization function
   - **Parameter:** `dbName` = "golang"
   - **Function Body Execution:**
     
     a. **`getenv("DB_USER", "")`** - `config/config.go:27` → `config/config.go:114`
        - **Purpose:** Retrieves PostgreSQL username from environment
        - **Calls:** `os.Getenv("DB_USER")` - Standard library function
        - **Returns:** Username string or empty string
        - **Why:** Environment variables provide flexible configuration without code changes
     
     b. **`getenv("DB_PASSWORD", "")`** - `config/config.go:28` → `config/config.go:114`
        - **Purpose:** Retrieves PostgreSQL password from environment
        - **Calls:** `os.Getenv("DB_PASSWORD")`
        - **Returns:** Password string or empty string
     
     c. **`getenv("DB_HOST", "")`** - `config/config.go:29` → `config/config.go:114`
        - **Purpose:** Retrieves PostgreSQL host from environment
        - **Calls:** `os.Getenv("DB_HOST")`
        - **Returns:** Host string or empty string
     
     d. **`getenv("DB_PORT", "")`** - `config/config.go:30` → `config/config.go:114`
        - **Purpose:** Retrieves PostgreSQL port from environment
        - **Calls:** `os.Getenv("DB_PORT")`
        - **Returns:** Port string or empty string
     
     e. **`getenv("DB_SSLMODE", "")`** - `config/config.go:31` → `config/config.go:114`
        - **Purpose:** Retrieves PostgreSQL SSL mode from environment
        - **Calls:** `os.Getenv("DB_SSLMODE")`
        - **Returns:** SSL mode string or empty string
     
     f. **`isRunningInDocker()`** - `config/config.go:34` → `config/config.go:123`
        - **Purpose:** Detects if application is running inside Docker container
        - **Function Body:**
          - **Calls:** `os.Stat("/.dockerenv")` - Standard library, checks for Docker file
          - **If exists:** Returns `true`
          - **If not:** **Calls:** `os.ReadFile("/proc/1/cgroup")` - Standard library
            - **Calls:** `strings.Contains(string(data), "docker")` - Standard library
            - **Returns:** `true` if contains "docker", `false` otherwise
        - **Why:** Docker networking requires `host.docker.internal` instead of `localhost`
        - **If Docker detected:** Modifies `host` variable to `"host.docker.internal"`
     
     g. **`fmt.Sprintf()`** - `config/config.go:38` - Standard library
        - **Purpose:** Builds PostgreSQL connection string
        - **Format:** `"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s"`
        - **Returns:** Connection string
     
     h. **`sql.Open("postgres", connStr)`** - `config/config.go:39` - Standard library
        - **Purpose:** Opens database connection (doesn't actually connect yet)
        - **Driver:** Uses `"github.com/lib/pq"` (imported as `_` for side effects)
        - **Returns:** `*sql.DB` object and error
        - **Why:** Lazy connection - actual connection happens on first query
     
     i. **`configureConnectionPool(db)`** - `config/config.go:44` → `config/config.go:97`
        - **Purpose:** Configures database connection pool settings
        - **Function Body:**
          - **Calls:** `db.SetMaxOpenConns(25)` - Standard library method
            - **Purpose:** Sets maximum open connections to 25
            - **Why:** Prevents too many connections overwhelming the database
          - **Calls:** `db.SetMaxIdleConns(10)` - Standard library method
            - **Purpose:** Sets maximum idle connections to 10
            - **Why:** Keeps some connections ready for reuse, improving performance
          - **Calls:** `db.SetConnMaxLifetime(5 * time.Minute)` - Standard library method
            - **Purpose:** Sets maximum connection lifetime to 5 minutes
            - **Why:** Prevents stale connections from causing issues
          - **Calls:** `db.SetConnMaxIdleTime(1 * time.Minute)` - Standard library method
            - **Purpose:** Sets maximum idle time to 1 minute
            - **Why:** Closes idle connections to free resources
        - **Returns:** Nothing (void function)
        - **Why:** Centralizes pool configuration for consistency across all databases
     
     j. **`db.Ping()`** - `config/config.go:46` - Standard library method
        - **Purpose:** Actually establishes connection to database
        - **Why:** Verifies connection works before proceeding
        - **If error:** Calls `log.Fatal()` - Application exits
     
     k. **`log.Printf()`** - `config/config.go:50` - Standard library
        - **Purpose:** Logs successful connection
        - **Message:** "Database connection established for {dbName} with optimized pool settings"
     
     l. **Returns:** `*sql.DB` connection object

4. **Back to `initDB()`** - `cmd/app/main.go:36`
   - **Calls:** `log.Printf("✓ %s database initialized", name)` - Standard library
   - **Returns:** `*sql.DB` connection object
   - **Why:** Provides consistent logging for all database initializations

5. **Back to `main()`** - `cmd/app/main.go:44`
   - **Stores:** Connection in `golangDB` variable
   - **Calls:** `defer golangDB.Close()` - Standard library method
     - **Purpose:** Ensures database connection is closed when application exits
     - **Why:** Prevents resource leaks

**Connection:** `golangDB` is now ready and will be used for user operations.

**Note:** The `configureConnectionPool()` function (called in step 3.i) sets up connection pool settings that are shared across all databases. This ensures consistent performance characteristics.

---

#### 2.2 Traffic Ticket Database (PostgreSQL) Initialization Chain

**Starting Point:** `cmd/app/main.go:47`
```go
trafficDB := initDB("traffic_ticket", config.InitTrafficDB)
```

**Complete Call Chain:**

1. **`initDB()`** - `cmd/app/main.go:36`
   - **Calls:** `config.InitTrafficDB()` → `config/config.go:92`

2. **`config.InitTrafficDB()`** - `config/config.go:92`
   - **Purpose:** Initializes connection to PostgreSQL "traffic_ticket" database
   - **Calls:** `InitDB("traffic_ticket")` → `config/config.go:26`
   - **Follows:** Same flow as 2.1 above, but connects to "traffic_ticket" database instead of "golang"
   - **Why:** Separate database for traffic ticket operations, isolating data

**Connection:** `trafficDB` is now ready and will be used for traffic ticket PostgreSQL operations.

---

#### 2.3 MySQL Database Initialization Chain

**Starting Point:** `cmd/app/main.go:50`
```go
mysqlDB := initDB("MySQL", config.InitMySQLDB)
```

**Complete Call Chain:**

1. **`initDB()`** - `cmd/app/main.go:36`
   - **Calls:** `config.InitMySQLDB()` → `config/config.go:105`

2. **`config.InitMySQLDB()`** - `config/config.go:105`
   - **Purpose:** Initializes connection to MySQL database
   - **Calls:** `initMySQLDB("MYSQL", "MySQL")` → `config/config.go:56`
   - **Why:** Provides a named function for MySQL database

3. **`initMySQLDB(prefix, displayName string)`** - `config/config.go:56`
   - **Purpose:** Core MySQL database initialization function
   - **Parameters:** 
     - `prefix` = "MYSQL" (for environment variable prefix)
     - `displayName` = "MySQL" (for logging)
   - **Function Body Execution:**
     
     a. **`getenv("MYSQL_HOST", "")`** - `config/config.go:57` → `config/config.go:114`
        - **Purpose:** Retrieves MySQL host from environment
        - **Calls:** `os.Getenv("MYSQL_HOST")`
        - **Returns:** Host string
     
     b. **`getenv("MYSQL_PORT", "")`** - `config/config.go:58` → `config/config.go:114`
        - **Purpose:** Retrieves MySQL port from environment
        - **Calls:** `os.Getenv("MYSQL_PORT")`
     
     c. **`getenv("MYSQL_USER", "")`** - `config/config.go:59` → `config/config.go:114`
        - **Purpose:** Retrieves MySQL username from environment
        - **Calls:** `os.Getenv("MYSQL_USER")`
     
     d. **`getenv("MYSQL_PASSWORD", "")`** - `config/config.go:60` → `config/config.go:114`
        - **Purpose:** Retrieves MySQL password from environment
        - **Calls:** `os.Getenv("MYSQL_PASSWORD")`
     
     e. **`getenv("MYSQL_DATABASE", "")`** - `config/config.go:61` → `config/config.go:114`
        - **Purpose:** Retrieves MySQL database name from environment
        - **Calls:** `os.Getenv("MYSQL_DATABASE")`
     
     f. **`isRunningInDocker()`** - `config/config.go:64` → `config/config.go:123`
        - **Same as PostgreSQL flow** - Detects Docker and adjusts host if needed
     
     g. **`fmt.Sprintf()`** - `config/config.go:68` - Standard library
        - **Purpose:** Builds MySQL connection string
        - **Format:** `"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"`
        - **Why:** MySQL uses different connection string format than PostgreSQL
     
     h. **`sql.Open("mysql", connStr)`** - `config/config.go:71` - Standard library
        - **Purpose:** Opens MySQL database connection
        - **Driver:** Uses `"github.com/go-sql-driver/mysql"` (imported as `_` for side effects)
        - **Returns:** `*sql.DB` object and error
     
     i. **`configureConnectionPool(db)`** - `config/config.go:76` → `config/config.go:97`
        - **Same as PostgreSQL flow** - Configures connection pool settings
        - **Why:** Shared function ensures consistent pool configuration
     
     j. **`db.Ping()`** - `config/config.go:78` - Standard library method
        - **Purpose:** Establishes actual MySQL connection
        - **If error:** Calls `log.Fatalf()` - Application exits
     
     k. **`log.Printf()`** - `config/config.go:82` - Standard library
        - **Purpose:** Logs successful MySQL connection
        - **Message:** "{displayName} connection established for {database} with optimized pool settings"
     
     l. **Returns:** `*sql.DB` connection object

4. **Back to `initDB()`** - `cmd/app/main.go:36`
   - **Calls:** `log.Printf("✓ MySQL database initialized")`
   - **Returns:** `*sql.DB` connection object

5. **Back to `main()`** - `cmd/app/main.go:50`
   - **Stores:** Connection in `mysqlDB` variable
   - **Calls:** `defer mysqlDB.Close()`

**Connection:** `mysqlDB` is now ready and will be used for MySQL traffic ticket operations.

---

#### 2.4 Passenger Plane Database (MySQL) Initialization Chain

**Starting Point:** `cmd/app/main.go:53`
```go
passengerDB := initDB("Passenger Plane MySQL", config.InitPassengerPlaneDB)
```

**Complete Call Chain:**

1. **`initDB()`** - `cmd/app/main.go:36`
   - **Calls:** `config.InitPassengerPlaneDB()` → `config/config.go:110`

2. **`config.InitPassengerPlaneDB()`** - `config/config.go:110`
   - **Purpose:** Initializes connection to passenger plane MySQL database
   - **Calls:** `initMySQLDB("PASSENGER_MYSQL", "Passenger Plane MySQL")` → `config/config.go:56`
   - **Why:** Separate MySQL database for passenger plane operations

3. **`initMySQLDB(prefix, displayName string)`** - `config/config.go:56`
   - **Same flow as 2.3**, but uses:
     - `prefix` = "PASSENGER_MYSQL" (for environment variables)
     - `displayName` = "Passenger Plane MySQL" (for logging)
   - **Environment Variables Used:**
     - `PASSENGER_MYSQL_HOST`
     - `PASSENGER_MYSQL_PORT`
     - `PASSENGER_MYSQL_USER`
     - `PASSENGER_MYSQL_PASSWORD`
     - `PASSENGER_MYSQL_DATABASE`

**Connection:** `passengerDB` is now ready and will be used for passenger plane operations.

---

### 3. Repository Layer Creation Chains

#### 3.1 User Repository Creation Chain

**Starting Point:** `cmd/app/main.go:57`
```go
userRepo := repository.NewUserRepository(golangDB)
```

**Complete Call Chain:**

1. **`repository.NewUserRepository(db *sql.DB)`** - `internal/repository/user_repo.go:14`
   - **Purpose:** Creates a new UserRepository instance
   - **Parameter:** `golangDB` - PostgreSQL connection for user operations
   - **Function Body:**
     - **Creates:** `UserRepository` struct with `DB` field set to `golangDB`
     - **Returns:** `*UserRepository` pointer
   - **Why:** Encapsulates database operations for user-related queries
   - **Connection:** Repository holds reference to `golangDB` connection

**Repository Methods Available:**
- `CreateUser(username, passwordHash string) error` - `internal/repository/user_repo.go:20`
  - **Uses:** Context timeout via `config.GetQueryTimeout()` → `config/config.go:110`
  - **Uses:** `handleQueryError()` for timeout error handling → `internal/repository/pagination.go:22`
- `GetUserByUsername(username string) (*entities.User, error)` - `internal/repository/user_repo.go:37`
  - **Uses:** Context timeout via `config.GetQueryTimeout()` → `config/config.go:110`
  - **Uses:** `handleQueryError()` for timeout error handling → `internal/repository/pagination.go:22`

---

#### 3.2 Traffic Ticket Repository Creation Chain

**Starting Point:** `cmd/app/main.go:62`
```go
trafficTicketRepo := repository.NewTrafficTicketRepository(trafficDB)
```

**Complete Call Chain:**

1. **`repository.NewTrafficTicketRepository(db *sql.DB)`** - `internal/repository/traffic_ticket_repo.go:12`
   - **Purpose:** Creates a new TrafficTicketRepository instance
   - **Parameter:** `trafficDB` - PostgreSQL connection for traffic ticket operations
   - **Function Body:**
     - **Creates:** `TrafficTicketRepository` struct with `db` field set to `trafficDB`
     - **Returns:** `*TrafficTicketRepository` pointer
   - **Why:** Encapsulates database operations for traffic ticket PostgreSQL queries
   - **Connection:** Repository holds reference to `trafficDB` connection

**Repository Methods Available:**
- `Insert(ticket *entities.TrafficTicket) error` - `internal/repository/traffic_ticket_repo.go:19`
  - **Uses:** Context timeout via `config.GetQueryTimeout()` → `config/config.go:110`
  - **Uses:** `handleQueryError()` for timeout error handling → `internal/repository/pagination.go:22`
- `List(limit, offset int) ([]*entities.TrafficTicket, error)` - `internal/repository/traffic_ticket_repo.go:72`
  - **Uses:** `normalizePagination()` for pagination validation → `internal/repository/pagination.go:11`
  - **Uses:** Context timeout via `config.GetQueryTimeout()` → `config/config.go:110`
  - **Uses:** `handleQueryError()` for timeout error handling → `internal/repository/pagination.go:22`
- `GetPaginated(limit, offset int) ([]entities.TrafficTicket, error)` - `internal/repository/traffic_ticket_repo.go:137`
  - **Uses:** Context timeout via `config.GetQueryTimeout()` → `config/config.go:110`
  - **Uses:** `handleQueryError()` for timeout error handling → `internal/repository/pagination.go:22`

---

#### 3.3 MySQL Traffic Ticket Repository Creation Chain

**Starting Point:** `cmd/app/main.go:66`
```go
mysqlRepo := repository.NewMySQLTrafficTicketRepository(mysqlDB)
```

**Complete Call Chain:**

1. **`repository.NewMySQLTrafficTicketRepository(db *sql.DB)`** - `internal/repository/mysql_traffic_ticket_repo.go:16`
   - **Purpose:** Creates a new MySQLTrafficTicketRepository instance
   - **Parameter:** `mysqlDB` - MySQL connection for traffic ticket operations
   - **Function Body:**
     - **Creates:** `MySQLTrafficTicketRepository` struct with `db` field set to `mysqlDB`
     - **Returns:** `*MySQLTrafficTicketRepository` pointer
   - **Why:** Encapsulates database operations for traffic ticket MySQL queries (different SQL syntax)
   - **Connection:** Repository holds reference to `mysqlDB` connection

**Repository Methods Available:**
- `Insert(ticket *entities.TrafficTicket) error` - `internal/repository/mysql_traffic_ticket_repo.go:20`
  - **Uses:** MySQL placeholders (`?` instead of `$1`)
  - **Uses:** Context timeout via `config.GetQueryTimeout()` → `config/config.go:110`
  - **Uses:** `handleQueryError()` for timeout error handling → `internal/repository/pagination.go:22`
- `List(limit, offset int) ([]*entities.TrafficTicket, error)` - `internal/repository/mysql_traffic_ticket_repo.go:65`
  - **Uses:** `normalizePagination()` for pagination validation → `internal/repository/pagination.go:11`
  - **Uses:** Context timeout via `config.GetQueryTimeout()` → `config/config.go:110`
  - **Uses:** `handleQueryError()` for timeout error handling → `internal/repository/pagination.go:22`

---

#### 3.4 Passenger Plane Repository Creation Chain

**Starting Point:** `cmd/app/main.go:70`
```go
passengerRepo := repository.NewPassengerPlaneRepository(passengerDB)
```

**Complete Call Chain:**

1. **`repository.NewPassengerPlaneRepository(db *sql.DB)`** - `internal/repository/passenger_plane_repo.go:16`
   - **Purpose:** Creates a new PassengerPlaneRepository instance
   - **Parameter:** `passengerDB` - MySQL connection for passenger plane operations
   - **Function Body:**
     - **Creates:** `PassengerPlaneRepository` struct with `db` field set to `passengerDB`
     - **Returns:** `*PassengerPlaneRepository` pointer
   - **Why:** Encapsulates database operations for passenger plane MySQL queries
   - **Connection:** Repository holds reference to `passengerDB` connection

**Repository Methods Available:**
- `Insert(p *entities.Passenger) error` - `internal/repository/passenger_plane_repo.go:20`
  - **Uses:** Context timeout via `config.GetQueryTimeout()` → `config/config.go:110`
  - **Uses:** `handleQueryError()` for timeout error handling → `internal/repository/pagination.go:22`
- `List(limit, offset int) ([]*entities.Passenger, error)` - `internal/repository/passenger_plane_repo.go:65`
  - **Uses:** `normalizePagination()` for pagination validation → `internal/repository/pagination.go:11`
  - **Uses:** Context timeout via `config.GetQueryTimeout()` → `config/config.go:110`
  - **Uses:** `handleQueryError()` for timeout error handling → `internal/repository/pagination.go:22`

---

### 3.5 Helper Functions in Repository Layer

#### 3.5.1 Pagination Helper Functions

**File:** `internal/repository/pagination.go`

**Purpose:** Provides consolidated helper functions used across all repositories

**Functions Available:**

1. **`normalizePagination(limit, offset int) (int, int)`** - `internal/repository/pagination.go:11`
   - **Purpose:** Validates and normalizes pagination parameters
   - **Parameters:**
     - `limit`: Maximum number of items to return
     - `offset`: Number of items to skip
   - **Function Body:**
     - **If limit <= 0:** Sets limit to 10 (default)
     - **If offset < 0:** Sets offset to 0 (default)
   - **Returns:** Normalized limit and offset
   - **Why:** Centralizes pagination validation logic, ensuring consistent behavior
   - **Used by:** All `List()` methods in repositories:
     - `TrafficTicketRepository.List()` - `internal/repository/traffic_ticket_repo.go:74`
     - `MySQLTrafficTicketRepository.List()` - `internal/repository/mysql_traffic_ticket_repo.go:80`
     - `PassengerPlaneRepository.List()` - `internal/repository/passenger_plane_repo.go:80`
     - `SQLServerTrafficTicketRepository.List()` - `internal/repository/sqlserver_traffic_ticket_repo.go:86`

2. **`handleQueryError(err error) error`** - `internal/repository/pagination.go:22`
   - **Purpose:** Centralizes timeout error handling for database queries
   - **Parameter:** `err` - Error from database operation
   - **Function Body:**
     - **Calls:** `context.DeadlineExceeded` comparison - Standard library
     - **If timeout:** Returns `errors.New("database query timeout: request took too long")`
     - **Otherwise:** Returns original error
   - **Returns:** Error (timeout message or original error)
   - **Why:** Provides consistent timeout error messages across all repositories
   - **Used by:** All repository methods that execute queries:
     - All `Insert()` methods
     - All `List()` methods
     - All `GetPaginated()` methods
     - `UserRepository.CreateUser()` and `GetUserByUsername()`

3. **`scanRows[T any](rows *sql.Rows, scanFunc func(*sql.Rows) (T, error)) ([]T, error)`** - `internal/repository/pagination.go:31`
   - **Purpose:** Generic helper for scanning database rows (currently defined but not yet used)
   - **Parameters:**
     - `rows`: Database rows from query
     - `scanFunc`: Function to scan a single row
   - **Function Body:**
     - **Calls:** `defer rows.Close()` - Standard library
     - **Iterates:** `rows.Next()` - Standard library
     - **Calls:** `scanFunc(rows)` for each row
     - **Calls:** `rows.Err()` - Standard library (checks for iteration errors)
   - **Returns:** Slice of scanned items and error
   - **Why:** Consolidates common row iteration pattern (available for future use)

---

### 3.6 Repository Method Execution Pattern (Example: TrafficTicketRepository.List)

**Starting Point:** Handler calls `trafficTicketRepo.List(limit, offset)`

**Complete Call Chain:**

1. **`TrafficTicketRepository.List(limit, offset int)`** - `internal/repository/traffic_ticket_repo.go:72`
   - **Purpose:** Retrieves paginated list of traffic tickets from PostgreSQL
   - **Parameters:** `limit` and `offset` for pagination
   
2. **Create Context with Timeout**
   - **Calls:** `context.WithTimeout(context.Background(), config.GetQueryTimeout())` - Standard library
     - **Calls:** `config.GetQueryTimeout()` → `config/config.go:110`
       - **Calls:** `getenvInt("DB_QUERY_TIMEOUT_SECONDS", 10)` → `config/config.go:116`
         - **Calls:** `getenv("DB_QUERY_TIMEOUT_SECONDS", "")` → `config/config.go:114`
           - **Calls:** `os.Getenv("DB_QUERY_TIMEOUT_SECONDS")` - Standard library
         - **Calls:** `strconv.Atoi(value)` - Standard library (if value exists)
         - **Returns:** Integer (10 or environment value)
       - **Calls:** `time.Duration(timeoutSeconds) * time.Second` - Standard library
       - **Returns:** `time.Duration` (default: 10 seconds)
     - **Returns:** Context with timeout and cancel function
   - **Calls:** `defer cancel()` - Standard library
     - **Purpose:** Ensures context is cancelled when function returns
     - **Why:** Prevents context leaks

3. **Normalize Pagination Parameters**
   - **Calls:** `normalizePagination(limit, offset)` → `internal/repository/pagination.go:11`
     - **Function:** `normalizePagination(limit, offset int) (int, int)`
     - **Function Body:**
       - **If limit <= 0:** Sets limit = 10 (default)
       - **If offset < 0:** Sets offset = 0 (default)
     - **Returns:** Normalized limit and offset
     - **Why:** Ensures valid pagination parameters before query execution

4. **Execute Database Query**
   - **Calls:** `r.db.QueryContext(ctx, query, limit, offset)` - Standard library
     - **Purpose:** Executes SELECT query with timeout context
     - **Database:** `trafficDB` connection (PostgreSQL "traffic_ticket" database)
     - **Context:** Uses timeout context (query automatically cancelled after timeout)
     - **Returns:** `*sql.Rows` and error
     - **Why:** `QueryContext` respects the timeout context

5. **Handle Query Error**
   - **If error:**
     - **Calls:** `handleQueryError(err)` → `internal/repository/pagination.go:22`
       - **Function:** `handleQueryError(err error) error`
       - **Function Body:**
         - **Checks:** `err == context.DeadlineExceeded` - Standard library
         - **If timeout:**
           - **Calls:** `errors.New("database query timeout: request took too long")` - Standard library
           - **Returns:** Timeout error message
         - **Otherwise:** Returns original error
       - **Returns:** Error (timeout message or original error)
     - **Returns:** `nil, error` (stops execution)

6. **Close Rows (Deferred)**
   - **Calls:** `defer rows.Close()` - Standard library
     - **Purpose:** Ensures rows are closed when function returns
     - **Why:** Prevents resource leaks

7. **Iterate Through Rows**
   - **Calls:** `rows.Next()` - Standard library (in loop)
     - **Purpose:** Advances to next row
     - **Returns:** `true` if more rows, `false` if done
   - **For each row:**
     - **Calls:** `rows.Scan(&t.ID, &t.DetectedSpeed, ...)` - Standard library
       - **Purpose:** Scans row data into `TrafficTicket` struct
       - **If error:** Returns `nil, error` (stops iteration)
     - **Appends:** Ticket to slice

8. **Check for Iteration Errors**
   - **Calls:** `rows.Err()` - Standard library
     - **Purpose:** Checks for errors during iteration
     - **If error:** Returns `nil, error`
     - **If no error:** Continues

9. **Return Results**
   - **Returns:** `[]*entities.TrafficTicket, nil` (success)
   - **Or:** `nil, error` (if any error occurred)

**Complete Function Call Chain:**
```
TrafficTicketRepository.List() → 
  context.WithTimeout() → config.GetQueryTimeout() → getenvInt() → getenv() → os.Getenv() →
  normalizePagination() →
  db.QueryContext() → PostgreSQL Database (trafficDB) →
  handleQueryError() (if error) →
  rows.Next() → rows.Scan() → rows.Err() →
  Return results
```

**Why This Pattern:**
- **Context Timeout:** Prevents queries from hanging indefinitely
- **Pagination Normalization:** Ensures valid pagination parameters
- **Error Handling:** Provides consistent timeout error messages
- **Resource Management:** Properly closes database resources

**Same Pattern Used By:**
- All `List()` methods in all repositories
- All `GetPaginated()` methods
- All `Insert()` methods (without pagination normalization)

---

### 4. Service Layer Creation Chain

#### 4.1 User Service Creation Chain

**Starting Point:** `cmd/app/main.go:58`
```go
userService := usecases.NewUserService(userRepo)
```

**Complete Call Chain:**

1. **`usecases.NewUserService(repo *repository.UserRepository)`** - `internal/usecases/user_service.go:18`
   - **Purpose:** Creates a new UserService instance
   - **Parameter:** `userRepo` - UserRepository instance created in step 3.1
   - **Function Body:**
     - **Creates:** `UserService` struct with `Repo` field set to `userRepo`
  - **Returns:** `*UserService` pointer
   - **Why:** Contains business logic for user operations (password hashing, validation)
   - **Connection:** Service holds reference to `userRepo`, which holds `golangDB` connection

**Service Methods Available:**
- `Register(creds entities.Credentials) error` - `internal/usecases/user_service.go:22`
- `Login(creds entities.Credentials) (string, error)` - `internal/usecases/user_service.go:35`

---

### 5. Handler Layer Creation Chains

#### 5.1 User Handler Creation Chain

**Starting Point:** `cmd/app/main.go:59`
```go
userHandler := httpDelivery.NewUserHandler(userService)
```

**Complete Call Chain:**

1. **`httpDelivery.NewUserHandler(service *usecases.UserService)`** - `internal/delivery/http/user_handles.go:16`
   - **Purpose:** Creates a new UserHandler instance
   - **Parameter:** `userService` - UserService instance created in step 4.1
   - **Function Body:**
     - **Creates:** `UserHandler` struct with `Service` field set to `userService`
  - **Returns:** `*UserHandler` pointer
   - **Why:** Handles HTTP requests/responses for user operations
   - **Connection Chain:** Handler → Service → Repository → Database (`golangDB`)

**Handler Methods Available:**
- `Register(w http.ResponseWriter, r *http.Request)` - `internal/delivery/http/user_handles.go:20`
- `Login(w http.ResponseWriter, r *http.Request)` - `internal/delivery/http/user_handles.go:40`
- `Dashboard(w http.ResponseWriter, r *http.Request)` - `internal/delivery/http/user_handles.go:61`

---

#### 5.2 Traffic Ticket Handler Creation Chain

**Starting Point:** `cmd/app/main.go:63`
```go
trafficHandler := httpDelivery.NewTrafficTicketHandler(trafficTicketRepo)
```

**Complete Call Chain:**

1. **`httpDelivery.NewTrafficTicketHandler(repo *repository.TrafficTicketRepository)`** - `internal/delivery/http/traffic_ticket_handler.go:17`
   - **Purpose:** Creates a new TrafficTicketHandler instance
   - **Parameter:** `trafficTicketRepo` - TrafficTicketRepository instance created in step 3.2
   - **Function Body:**
     - **Creates:** `TrafficTicketHandler` struct with `repo` field set to `trafficTicketRepo`
  - **Returns:** `*TrafficTicketHandler` pointer
   - **Why:** Handles HTTP requests/responses for traffic ticket PostgreSQL operations
   - **Connection Chain:** Handler → Repository → Database (`trafficDB`)

**Handler Methods Available:**
- `Create(w http.ResponseWriter, r *http.Request)` - `internal/delivery/http/traffic_ticket_handler.go:21`
  - **Uses:** `CreateHandler[entities.TrafficTicket]()` from `handler_helpers.go`
- `List(w http.ResponseWriter, r *http.Request)` - Uses `ListHandler[entities.TrafficTicket]()`
- `GetPaginated(w http.ResponseWriter, r *http.Request)` - Uses `GetPaginatedHandler[entities.TrafficTicket]()`

---

#### 5.3 MySQL Traffic Ticket Handler Creation Chain

**Starting Point:** `cmd/app/main.go:67`
```go
mysqlHandler := httpDelivery.NewMySQLTrafficTicketHandler(mysqlRepo)
```

**Complete Call Chain:**

1. **`httpDelivery.NewMySQLTrafficTicketHandler(repo *repository.MySQLTrafficTicketRepository)`** - `internal/delivery/http/mysql_traffic_ticket_handler.go`
   - **Purpose:** Creates a new MySQLTrafficTicketHandler instance
   - **Parameter:** `mysqlRepo` - MySQLTrafficTicketRepository instance created in step 3.3
   - **Connection Chain:** Handler → Repository → Database (`mysqlDB`)

---

#### 5.4 Passenger Plane Handler Creation Chain

**Starting Point:** `cmd/app/main.go:71`
```go
passengerHandler := httpDelivery.NewPassengerPlaneHandler(passengerRepo)
```

**Complete Call Chain:**

1. **`httpDelivery.NewPassengerPlaneHandler(repo *repository.PassengerPlaneRepository)`** - `internal/delivery/http/passenger_plane_handler.go`
   - **Purpose:** Creates a new PassengerPlaneHandler instance
   - **Parameter:** `passengerRepo` - PassengerPlaneRepository instance created in step 3.4
   - **Connection Chain:** Handler → Repository → Database (`passengerDB`)

---

### 6. Middleware Configuration Chain

#### 6.1 Rate Limit Configuration Chain

**Starting Point:** `cmd/app/main.go:77`
```go
rateLimitRequests := getEnvInt("RATE_LIMIT_REQUESTS", 100)
```

**Complete Call Chain:**

1. **`getEnvInt(key string, defaultValue int)`** - `cmd/app/main.go:118`
   - **Purpose:** Retrieves integer environment variable with fallback
   - **Parameter:** `key` = "RATE_LIMIT_REQUESTS", `defaultValue` = 100
   - **Function Body:**
     a. **`os.Getenv("RATE_LIMIT_REQUESTS")`** - Standard library
        - **Purpose:** Gets environment variable value
        - **Returns:** String value or empty string
     b. **If value exists:**
        - **`strconv.Atoi(value)`** - Standard library
          - **Purpose:** Converts string to integer
          - **Returns:** Integer value and error
          - **If success:** Returns integer value
          - **If error:** Returns default value
     c. **If value doesn't exist:** Returns default value (100)
   - **Returns:** Integer (100 or environment value)
   - **Why:** Provides safe integer parsing with fallback

2. **`getEnvInt("RATE_LIMIT_BURST", 10)`** - `cmd/app/main.go:78`
   - **Same flow as above**, but for "RATE_LIMIT_BURST" with default 10

3. **`middleware.RateLimitMiddleware(rateLimitRequests, rateLimitBurst)`** - `cmd/app/main.go:79` → `pkg/middleware/ratelimit.go:99`
   - **Purpose:** Creates rate limiting middleware function
   - **Parameters:** 
     - `rate` = rateLimitRequests (e.g., 100)
     - `burst` = rateLimitBurst (e.g., 10)
   - **Function Body:**
     a. **`NewRateLimiter(rate, burst)`** - `pkg/middleware/ratelimit.go:28`
        - **Purpose:** Creates a new RateLimiter instance
        - **Function Body:**
          - **Creates:** `RateLimiter` struct with:
            - `requests`: Empty map for token buckets
            - `rate`: Requests per minute
            - `burst`: Maximum burst capacity
            - `cleanupInterval`: 5 minutes
          - **Calls:** `rl.startCleanup()` - `pkg/middleware/ratelimit.go:37` → `pkg/middleware/ratelimit.go:81`
            - **Purpose:** Starts background goroutine for cleanup
            - **Function Body:**
              - **Calls:** `time.NewTicker(rl.cleanupInterval)` - Standard library
                - **Purpose:** Creates ticker that fires every 5 minutes
              - **Infinite loop:**
                - **Waits:** For ticker to fire
                - **Calls:** `rl.mutex.Lock()` - Standard library (mutex lock)
                - **Iterates:** Through `rl.requests` map
                - **Calls:** `time.Now()` - Standard library
                - **Checks:** If bucket is older than 1 hour
                - **Calls:** `delete(rl.requests, key)` - Standard library (removes old entries)
                - **Calls:** `rl.mutex.Unlock()` - Standard library (mutex unlock)
            - **Why:** Prevents memory leaks by removing old rate limit entries
          - **Returns:** `*RateLimiter` pointer
        - **Why:** Token bucket algorithm for rate limiting
     b. **Returns:** Middleware function `func(http.HandlerFunc) http.HandlerFunc`
        - **This function:**
          - **Calls:** `getClientIP(r)` - `pkg/middleware/ratelimit.go:105` → `pkg/middleware/ratelimit.go:124`
            - **Purpose:** Extracts client IP from request
            - **Function Body:**
              - **Calls:** `r.Header.Get("X-Forwarded-For")` - Standard library
                - **If exists:** Returns that value
              - **Calls:** `r.Header.Get("X-Real-IP")` - Standard library
                - **If exists:** Returns that value
              - **Returns:** `r.RemoteAddr` - Standard library (fallback)
            - **Why:** Handles proxies/load balancers that modify IP
          - **Calls:** `limiter.Allow(clientIP)` - `pkg/middleware/ratelimit.go:107` → `pkg/middleware/ratelimit.go:43`
            - **Purpose:** Checks if request is allowed
            - **Function Body:**
              - **Calls:** `rl.mutex.Lock()` - Standard library (mutex lock)
              - **Calls:** `time.Now()` - Standard library
              - **Checks:** If bucket exists for client IP
              - **If not exists:**
                - **Creates:** New `TokenBucket` with tokens = burst - 1
                - **Returns:** `true` (allowed)
              - **If exists:**
                - **Calculates:** Time passed since last refill
                - **Calculates:** Tokens to add based on rate
                - **Calls:** `min(bucket.tokens+tokensToAdd, bucket.burst)` - `pkg/middleware/ratelimit.go:67` → `pkg/middleware/ratelimit.go:140`
                  - **Purpose:** Returns minimum of two integers
                - **Updates:** Bucket tokens
                - **If tokens > 0:**
                  - **Decrements:** Tokens
                  - **Returns:** `true` (allowed)
                - **If tokens = 0:**
                  - **Returns:** `false` (rate limited)
              - **Calls:** `rl.mutex.Unlock()` - Standard library (mutex unlock)
            - **Returns:** `bool` (true if allowed, false if rate limited)
          - **If not allowed:**
            - **Calls:** `w.Header().Set("Content-Type", "application/json")` - Standard library
            - **Calls:** `w.WriteHeader(http.StatusTooManyRequests)` - Standard library (429)
            - **Calls:** `fmt.Fprintf(w, ...)` - Standard library (writes JSON error)
            - **Returns:** (stops request)
          - **If allowed:**
            - **Calls:** `next.ServeHTTP(w, r)` - Standard library (continues to next middleware/handler)
   - **Returns:** Middleware function
   - **Why:** Protects API from abuse by limiting requests per IP

---

#### 6.2 JWT Authentication Middleware Chain

**Used in:** `cmd/app/main.go:27` (RouteRegistrar.Protected method)
```go
rr.rateLimit(jwtutil.AuthMiddleware(handler))
```

**Complete Call Chain:**

1. **`jwtutil.AuthMiddleware(next http.HandlerFunc)`** - `pkg/jwtutil/middleware.go:8`
   - **Purpose:** Creates JWT authentication middleware function
   - **Parameter:** `next` - Next handler in chain
   - **Function Body:**
     a. **Returns:** Middleware function `func(w http.ResponseWriter, r *http.Request)`
        - **This function:**
          - **Calls:** `r.Header.Get("Authorization")` - Standard library
            - **Purpose:** Gets Authorization header
            - **If empty:**
              - **Calls:** `http.Error(w, "Missing Authorization header", http.StatusUnauthorized)` - Standard library
              - **Returns:** (stops request)
          - **Calls:** `VerifyToken(authHeader)` - `pkg/jwtutil/middleware.go:16` → `pkg/jwtutil/jwt.go:30`
            - **Purpose:** Verifies JWT token
            - **Function Body:**
              - **Calls:** `strings.TrimPrefix(authHeader, "Bearer ")` - Standard library
                - **Purpose:** Removes "Bearer " prefix from token
              - **Creates:** Empty `Claims` struct
              - **Calls:** `jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) { return jwtSecret, nil })` - Third-party library (`github.com/golang-jwt/jwt/v5`)
                - **Purpose:** Parses and validates JWT token
                - **Validates:** Token signature using `jwtSecret`
                - **Validates:** Token expiration
                - **Returns:** Parsed token and error
              - **If error or invalid:**
                - **Returns:** Error
              - **If valid:**
                - **Returns:** `claims.Username` (string)
            - **Returns:** Username string and error
          - **If error:**
            - **Calls:** `http.Error(w, "Invalid or expired token", http.StatusUnauthorized)` - Standard library
            - **Returns:** (stops request)
          - **If valid:**
            - **Calls:** `next.ServeHTTP(w, r)` - Standard library (continues to handler)
   - **Returns:** Middleware function
   - **Why:** Protects endpoints by requiring valid JWT token

---

### 7. Route Registration Chain

#### 7.1 Router Creation

**Starting Point:** `cmd/app/main.go:73`
```go
router := http.NewServeMux()
```

**Complete Call Chain:**

1. **`http.NewServeMux()`** - Standard library
   - **Purpose:** Creates new HTTP request multiplexer (router)
   - **Returns:** `*http.ServeMux` pointer
   - **Why:** Routes HTTP requests to appropriate handlers

---

#### 7.2 RouteRegistrar Creation

**Starting Point:** `cmd/app/main.go:82`
```go
routes := &RouteRegistrar{router: router, rateLimit: rateLimit}
```

**Purpose:** Creates RouteRegistrar instance for consistent route registration
- **Fields:**
  - `router`: HTTP router created in step 7.1
  - `rateLimit`: Rate limiting middleware created in step 6.1

---

#### 7.3 Public Route Registration

**Starting Point:** `cmd/app/main.go:85`
```go
routes.Public("/api/register", userHandler.Register)
```

**Complete Call Chain:**

1. **`RouteRegistrar.Public(path string, handler http.HandlerFunc)`** - `cmd/app/main.go:31`
   - **Purpose:** Registers public route with rate limiting only
   - **Parameters:**
     - `path` = "/api/register"
     - `handler` = `userHandler.Register`
   - **Function Body:**
     - **Calls:** `rr.rateLimit(handler)` - Rate limiting middleware wraps handler
     - **Calls:** `rr.router.HandleFunc(path, wrappedHandler)` - Standard library
       - **Purpose:** Registers route in HTTP router
   - **Why:** Public endpoints only need rate limiting, not JWT auth

**Same flow for:** `/api/login` route

---

#### 7.4 Protected Route Registration

**Starting Point:** `cmd/app/main.go:91`
```go
routes.Protected("/api/traffic_tickets/inputpostgre", trafficHandler.Create)
```

**Complete Call Chain:**

1. **`RouteRegistrar.Protected(path string, handler http.HandlerFunc)`** - `cmd/app/main.go:26`
   - **Purpose:** Registers protected route with rate limiting and JWT auth
   - **Parameters:**
     - `path` = "/api/traffic_tickets/inputpostgre"
     - `handler` = `trafficHandler.Create`
   - **Function Body:**
     - **Calls:** `jwtutil.AuthMiddleware(handler)` - JWT middleware wraps handler
       - **Returns:** JWT-protected handler
     - **Calls:** `rr.rateLimit(jwtProtectedHandler)` - Rate limiting wraps JWT handler
       - **Returns:** Rate-limited + JWT-protected handler
     - **Calls:** `rr.router.HandleFunc(path, finalHandler)` - Standard library
   - **Why:** Protected endpoints need both rate limiting and authentication

**Middleware Chain Order:**
1. Rate Limit Middleware (outermost - checks first)
2. JWT Auth Middleware (middle - checks second)
3. Handler Function (innermost - executes last)

**Same flow for:** All other protected routes

---

### 8. Server Startup Chain

**Starting Point:** `cmd/app/main.go:105`
```go
port := getEnv("APP_PORT", "8080")
```

**Complete Call Chain:**

1. **`getEnv(key, defaultValue string)`** - `cmd/app/main.go:111`
   - **Purpose:** Retrieves environment variable with fallback
   - **Parameter:** `key` = "APP_PORT", `defaultValue` = "8080"
   - **Function Body:**
     - **Calls:** `os.Getenv("APP_PORT")` - Standard library
     - **If value exists and not empty:** Returns value
     - **If value doesn't exist or empty:** Returns "8080"
   - **Returns:** Port string
   - **Why:** Provides flexible port configuration without code changes

2. **`getEnvInt("HTTP_READ_TIMEOUT_SECONDS", 15)`** - `cmd/app/main.go:109` → `cmd/app/main.go:118`
   - **Purpose:** Retrieves HTTP read timeout from environment
   - **Calls:** `getEnvInt()` → `cmd/app/main.go:118`
     - **Function:** `getEnvInt(key string, defaultValue int)`
     - **Calls:** `os.Getenv("HTTP_READ_TIMEOUT_SECONDS")` - Standard library
     - **Calls:** `strconv.Atoi(value)` - Standard library (if value exists)
     - **Returns:** Integer (15 or environment value)
   - **Returns:** Integer timeout in seconds
   - **Why:** Prevents requests from reading indefinitely

3. **`getEnvInt("HTTP_WRITE_TIMEOUT_SECONDS", 15)`** - `cmd/app/main.go:110` → `cmd/app/main.go:118`
   - **Same flow as above**, but for write timeout
   - **Returns:** Integer timeout in seconds
   - **Why:** Prevents responses from writing indefinitely

4. **`getEnvInt("HTTP_IDLE_TIMEOUT_SECONDS", 60)`** - `cmd/app/main.go:111` → `cmd/app/main.go:118`
   - **Same flow as above**, but for idle timeout
   - **Returns:** Integer timeout in seconds
   - **Why:** Closes idle keep-alive connections to free resources

5. **`http.Server{}` struct creation** - `cmd/app/main.go:115`
   - **Purpose:** Creates HTTP server with timeout configuration
   - **Fields:**
     - `Addr`: ":" + port
     - `Handler`: router
     - `ReadTimeout`: `time.Duration(readTimeout) * time.Second` - Standard library conversion
     - `WriteTimeout`: `time.Duration(writeTimeout) * time.Second` - Standard library conversion
     - `IdleTimeout`: `time.Duration(idleTimeout) * time.Second` - Standard library conversion
   - **Why:** Configures server timeouts to prevent hanging requests

6. **`fmt.Printf("Server listening on port %s with timeouts...")`** - `cmd/app/main.go:123` - Standard library
   - **Purpose:** Logs server startup with timeout information

7. **`server.ListenAndServe()`** - `cmd/app/main.go:125` - Standard library method
   - **Purpose:** Starts HTTP server with configured timeouts
   - **Behavior:**
     - **Binds:** To specified port
     - **Listens:** For incoming HTTP requests
     - **Applies:** ReadTimeout, WriteTimeout, and IdleTimeout to all connections
     - **Blocks:** Until server stops or error
     - **Routes:** Requests to registered handlers
   - **If error:** Calls `log.Fatal()` - Application exits
   - **Why:** Timeouts ensure requests don't hang indefinitely, protecting server resources

---

### 8.5 Configuration Helper Functions

#### 8.5.1 Query Timeout Configuration

**File:** `config/config.go`

**Function:** `GetQueryTimeout() time.Duration` - `config/config.go:110`

**Purpose:** Returns the query timeout duration for database operations

**Complete Call Chain:**

1. **`config.GetQueryTimeout()`** - `config/config.go:110`
   - **Purpose:** Gets query timeout from environment variable
   - **Function Body:**
     a. **`getenvInt("DB_QUERY_TIMEOUT_SECONDS", 10)`** - `config/config.go:111` → `config/config.go:116`
        - **Purpose:** Retrieves integer environment variable with fallback
        - **Function:** `getenvInt(key string, defaultValue int)`
        - **Function Body:**
          - **Calls:** `getenv(key, "")` - `config/config.go:117` → `config/config.go:114`
            - **Calls:** `os.Getenv("DB_QUERY_TIMEOUT_SECONDS")` - Standard library
            - **Returns:** String value or empty string
          - **If value exists:**
            - **Calls:** `strconv.Atoi(value)` - Standard library
              - **Purpose:** Converts string to integer
              - **Returns:** Integer value and error
          - **If value doesn't exist or conversion fails:** Returns default (10)
        - **Returns:** Integer timeout in seconds (default: 10)
     b. **`time.Duration(timeoutSeconds) * time.Second`** - Standard library
        - **Purpose:** Converts seconds to `time.Duration`
        - **Returns:** `time.Duration` (e.g., 10 seconds)
   - **Returns:** `time.Duration` (default: 10 seconds)
   - **Why:** Centralizes query timeout configuration, used by all repository methods
   - **Used by:** All repository methods that create context timeouts:
     - `UserRepository.CreateUser()` - `internal/repository/user_repo.go:22`
     - `UserRepository.GetUserByUsername()` - `internal/repository/user_repo.go:39`
     - `TrafficTicketRepository.Insert()` - `internal/repository/traffic_ticket_repo.go:20`
     - `TrafficTicketRepository.List()` - `internal/repository/traffic_ticket_repo.go:75`
     - `TrafficTicketRepository.GetPaginated()` - `internal/repository/traffic_ticket_repo.go:140`
     - And all other repository methods

---

## Request Flow Chains

This section traces complete request flows from HTTP request arrival through all layers to database and back.

### Request Flow: POST /api/register (Public Endpoint)

**Complete Request Flow:**

1. **HTTP Request Arrives**
- **Location:** HTTP server receives request
   - **Router:** `cmd/app/main.go:73` (router created in main)

2. **Route Matching**
   - **Location:** `cmd/app/main.go:85`
   - **Route:** `/api/register` matches request
- **Handler Chain:** `rateLimit(userHandler.Register)`

3. **Rate Limit Middleware Execution**
   - **File:** `pkg/middleware/ratelimit.go:99` (middleware function returned)
   - **Execution:**
     a. **Extract Client IP**
        - **Calls:** `getClientIP(r)` → `pkg/middleware/ratelimit.go:124`
        - **Returns:** Client IP string
     b. **Check Rate Limit**
        - **Calls:** `limiter.Allow(clientIP)` → `pkg/middleware/ratelimit.go:43`
        - **Returns:** `true` (allowed) or `false` (rate limited)
     c. **If rate limited:**
        - **Calls:** `w.WriteHeader(http.StatusTooManyRequests)` - 429 status
        - **Returns:** (stops request)
     d. **If allowed:**
        - **Calls:** `next.ServeHTTP(w, r)` → Continues to handler

4. **Handler Function Execution**
   - **File:** `internal/delivery/http/user_handles.go:20`
   - **Function:** `UserHandler.Register(w http.ResponseWriter, r *http.Request)`
   - **Execution:**
     a. **Check HTTP Method**
        - **Calls:** `r.Method != http.MethodPost` - Standard library
        - **If not POST:**
          - **Calls:** `WriteMethodNotAllowed(w)` → `internal/delivery/http/response.go:70`
            - **Calls:** `WriteErrorResponse(w, http.StatusMethodNotAllowed, "Only POST method allowed")` → `internal/delivery/http/response.go:18`
              - **Calls:** `w.Header().Set("Content-Type", "application/json")` - Standard library
              - **Calls:** `w.WriteHeader(http.StatusMethodNotAllowed)` - Standard library (405)
              - **Calls:** `json.NewEncoder(w).Encode(Response{...})` - Standard library
          - **Returns:** (stops request)
     b. **Decode Request Body**
        - **Calls:** `json.NewDecoder(r.Body).Decode(&creds)` - Standard library
          - **Purpose:** Parses JSON body into `entities.Credentials` struct
          - **File:** `internal/entities/user.go` (Credentials struct definition)
        - **If error:**
          - **Calls:** `WriteBadRequest(w, "Invalid request body: "+err.Error())` → `internal/delivery/http/response.go:62`
            - **Calls:** `WriteErrorResponse(w, http.StatusBadRequest, message)` → `internal/delivery/http/response.go:18`
          - **Returns:** (stops request)
     c. **Call Service Layer**
        - **Calls:** `h.Service.Register(creds)` → `internal/usecases/user_service.go:22`
          - **Function:** `UserService.Register(creds entities.Credentials) error`
          - **Execution:**
            - **Validate Input**
              - **Checks:** `creds.Username == "" || creds.Password == ""`
              - **If empty:** Returns error
            - **Hash Password**
              - **Calls:** `bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)` - Third-party library (`golang.org/x/crypto/bcrypt`)
                - **Purpose:** Hashes password using bcrypt algorithm
                - **Returns:** Hashed password bytes and error
            - **Call Repository**
              - **Calls:** `s.Repo.CreateUser(creds.Username, string(hashedPassword))` → `internal/repository/user_repo.go:20`
                - **Function:** `UserRepository.CreateUser(username, passwordHash string) error`
                - **Execution:**
                  - **Create Context with Timeout**
                    - **Calls:** `context.WithTimeout(context.Background(), config.GetQueryTimeout())` - Standard library
                      - **Calls:** `config.GetQueryTimeout()` → `config/config.go:110`
                        - **Calls:** `getenvInt("DB_QUERY_TIMEOUT_SECONDS", 10)` → `config/config.go:116`
                          - **Calls:** `getenv("DB_QUERY_TIMEOUT_SECONDS", "")` → `config/config.go:114`
                            - **Calls:** `os.Getenv("DB_QUERY_TIMEOUT_SECONDS")` - Standard library
                          - **Calls:** `strconv.Atoi(value)` - Standard library (if value exists)
                        - **Returns:** `time.Duration` (default: 10 seconds)
                      - **Returns:** Context with timeout and cancel function
                    - **Calls:** `defer cancel()` - Standard library (ensures context cleanup)
                  - **SQL Query:** `INSERT INTO users (username, password) VALUES ($1, $2) ON CONFLICT (username) DO NOTHING RETURNING id;`
                  - **Calls:** `r.DB.QueryRowContext(ctx, query, username, passwordHash).Scan(&id)` - Standard library
                    - **Purpose:** Executes SQL query on PostgreSQL database with timeout
                    - **Database:** `golangDB` connection (PostgreSQL "golang" database)
                    - **Context:** Uses timeout context (query will be cancelled after timeout)
                    - **Returns:** Error if username exists, timeout, or other DB error
                  - **If error == sql.ErrNoRows:**
                    - **Returns:** Error("username already exists")
                  - **Calls:** `handleQueryError(err)` → `internal/repository/pagination.go:22`
                    - **Purpose:** Checks for timeout and returns consistent error message
                    - **Function Body:**
                      - **Checks:** `err == context.DeadlineExceeded` - Standard library
                      - **If timeout:** Returns `errors.New("database query timeout: request took too long")`
                      - **Otherwise:** Returns original error
                    - **Returns:** Error (timeout message or original error)
                  - **Returns:** Error (if any) or nil (success)
              - **Returns:** Error (if any) or nil (success)
          - **If error:**
            - **Calls:** `WriteBadRequest(w, err.Error())` → `internal/delivery/http/response.go:62`
            - **Returns:** (stops request)
     d. **Return Success Response**
        - **Calls:** `WriteSuccessResponseCreated(w, []interface{}{}, "Registration successful")` → `internal/delivery/http/response.go:45`
          - **Calls:** `WriteSuccessResponse(w, http.StatusCreated, data, message)` → `internal/delivery/http/response.go:29`
            - **Calls:** `w.Header().Set("Content-Type", "application/json")` - Standard library
            - **Calls:** `w.WriteHeader(http.StatusCreated)` - Standard library (201)
            - **Calls:** `json.NewEncoder(w).Encode(Response{...})` - Standard library
        - **Returns:** (request complete)

**Complete Call Chain Summary:**
```
HTTP Request → Router → Rate Limit Middleware → UserHandler.Register → 
UserService.Register → UserRepository.CreateUser → PostgreSQL Database (golangDB) → 
Response back through chain
```

---

### Request Flow: POST /api/traffic_tickets/inputpostgre (Protected Endpoint)

**Complete Request Flow:**

1. **HTTP Request Arrives**
- **Location:** HTTP server receives request
   - **Router:** `cmd/app/main.go:73`

2. **Route Matching**
   - **Location:** `cmd/app/main.go:91`
- **Route:** `/api/traffic_tickets/inputpostgre` matches
   - **Handler Chain:** `rateLimit(jwtutil.AuthMiddleware(trafficHandler.Create))`

3. **Rate Limit Middleware Execution**
   - **Same as Public Endpoint flow above**
- **If allowed:** Continues to JWT middleware

4. **JWT Authentication Middleware Execution**
   - **File:** `pkg/jwtutil/middleware.go:8` (middleware function)
   - **Execution:**
     a. **Extract Authorization Header**
        - **Calls:** `r.Header.Get("Authorization")` - Standard library
        - **If empty:**
          - **Calls:** `http.Error(w, "Missing Authorization header", http.StatusUnauthorized)` - Standard library (401)
          - **Returns:** (stops request)
     b. **Verify Token**
        - **Calls:** `VerifyToken(authHeader)` → `pkg/jwtutil/jwt.go:30`
          - **Calls:** `strings.TrimPrefix(authHeader, "Bearer ")` - Standard library
          - **Calls:** `jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) { return jwtSecret, nil })` - Third-party library
       - **Validates:** Token signature and expiration
       - **If invalid:** Returns error
          - **If valid:** Returns username
        - **If error:**
          - **Calls:** `http.Error(w, "Invalid or expired token", http.StatusUnauthorized)` - Standard library (401)
          - **Returns:** (stops request)
     c. **If valid:**
        - **Calls:** `next.ServeHTTP(w, r)` → Continues to handler

5. **Handler Function Execution**
   - **File:** `internal/delivery/http/traffic_ticket_handler.go:21`
   - **Function:** `TrafficTicketHandler.Create(w http.ResponseWriter, r *http.Request)`
   - **Execution:**
     a. **Calls Generic Handler**
        - **Calls:** `CreateHandler[entities.TrafficTicket](w, r, h.repo, CreateConfig{...})` → `internal/delivery/http/handler_helpers.go:42`
          - **Function:** Generic handler for creating entities
          - **Execution:**
            - **Decode Request Body**
              - **Calls:** `json.NewDecoder(r.Body).Decode(&items)` - Standard library
                - **Tries:** Decode as array `[]entities.TrafficTicket`
                - **If error:**
                  - **Tries:** Decode as single `entities.TrafficTicket`
                  - **If error:**
                    - **Calls:** `WriteBadRequest(w, config.InvalidBodyMessage+err.Error())` → `internal/delivery/http/response.go:62`
                    - **Returns:** (stops request)
            - **Insert Each Item**
              - **For each item:**
                - **Calls:** `repo.Insert(&item)` → `internal/repository/traffic_ticket_repo.go:19`
                  - **Function:** `TrafficTicketRepository.Insert(ticket *entities.TrafficTicket) error`
                  - **Execution:**
                    - **Create Context with Timeout**
                      - **Calls:** `context.WithTimeout(context.Background(), config.GetQueryTimeout())` - Standard library
                        - **Calls:** `config.GetQueryTimeout()` → `config/config.go:110`
                          - **Same flow as UserRepository.CreateUser above**
                        - **Returns:** Context with timeout and cancel function
                      - **Calls:** `defer cancel()` - Standard library
                    - **SQL Query:** `INSERT INTO traffic_tickets (...) VALUES ($1, $2, ..., $23)`
                    - **Calls:** `r.db.ExecContext(ctx, query, ticket.DetectedSpeed, ticket.LegalSpeed, ...)` - Standard library
                      - **Purpose:** Executes SQL query on PostgreSQL database with timeout
                      - **Database:** `trafficDB` connection (PostgreSQL "traffic_ticket" database)
                      - **Context:** Uses timeout context (query will be cancelled after timeout)
                      - **Returns:** Error if insertion fails or timeout occurs
                    - **Calls:** `handleQueryError(err)` → `internal/repository/pagination.go:22`
                      - **Purpose:** Checks for timeout and returns consistent error message
                      - **Returns:** Error (timeout message or original error)
                    - **Returns:** Error (if any) or nil (success)
                - **If error:**
                  - **Calls:** `WriteInternalServerError(w, fmt.Sprintf(config.InsertErrorMessage, err))` → `internal/delivery/http/response.go:74`
                  - **Returns:** (stops request)
            - **Return Success Response**
              - **Calls:** `WriteSuccessResponseCreated(w, []interface{}{}, fmt.Sprintf(config.SuccessMessage, len(items)))` → `internal/delivery/http/response.go:45`
              - **Returns:** (request complete)

**Complete Call Chain Summary:**
```
HTTP Request → Router → Rate Limit Middleware → JWT Auth Middleware → 
TrafficTicketHandler.Create → CreateHandler (generic) → TrafficTicketRepository.Insert → 
PostgreSQL Database (trafficDB) → Response back through chain
```

---

### Request Flow: GET /api/traffic_tickets/listpostgre (Protected Endpoint with Pagination)

**Complete Request Flow:**

**HTTP Request Arrives** → Router → Rate Limit → JWT Auth (same as above)

1. **HTTP Request Arrives**
- **Location:** HTTP server receives request
   - **Router:** `cmd/app/main.go:137`

2. **Route Matching**
   - **Location:** `cmd/app/main.go:175`
- **Route:** `/api/traffic_tickets/listpostgre` matches
   - **Handler Chain:** `rateLimit(jwtutil.AuthMiddleware(trafficHandler.GetPaginated))`

3. **Rate Limit Middleware Execution**
  - **File:** `pkg/middleware/ratelimit.go:103` (middleware function returned)
   - **Execution:**
     a. **Extract Client IP**
        - **Calls:** `getClientIP(r)` → `pkg/middleware/ratelimit.go:127`
        - **Returns:** Client IP string
     b. **Check Rate Limit**
        - **Calls:** `limiter.Allow(clientIP)` → `pkg/middleware/ratelimit.go:47`
        - **Returns:** `true` (allowed) or `false` (rate limited)
     c. **If rate limited:**
        - **Calls:** `w.WriteHeader(http.StatusTooManyRequests)` - 429 status
        - **Returns:** (stops request)
     d. **If allowed:**
        - **Calls:** `next.ServeHTTP(w, r)` → Continues to handler
   - **Same as Public Endpoint flow above**
- **If allowed:** Continues to JWT middleware
  **JWT Authentication Middleware Execution**
   - **File:** `pkg/jwtutil/middleware.go:8` (middleware function)
   - **Execution:**
     a. **Extract Authorization Header**
        - **Calls:** `r.Header.Get("Authorization")` - Standard library
        - **If empty:**
          - **Calls:** `http.Error(w, "Missing Authorization header", http.StatusUnauthorized)` - Standard library (401)
          - **Returns:** (stops request)
     b. **Verify Token**
        - **Calls:** `VerifyToken(authHeader)` → `pkg/jwtutil/jwt.go:34`
          - **Calls:** `strings.TrimPrefix(authHeader, "Bearer ")` - Standard library
          - **Calls:** `jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) { return jwtSecret, nil })` - Third-party library
       - **Validates:** Token signature and expiration
       - **If invalid:** Returns error
          - **If valid:** Returns username
        - **If error:**
          - **Calls:** `http.Error(w, "Invalid or expired token", http.StatusUnauthorized)` - Standard library (401)
          - **Returns:** (stops request)
     c. **If valid:**
        - **Calls:** `next.ServeHTTP(w, r)` → Continues to handler

  

4. **Handler Function Execution**
   - **File:** `internal/delivery/http/traffic_ticket_handler.go`
   - **Function:** `TrafficTicketHandler.GetPaginated(w http.ResponseWriter, r *http.Request)`
   - **Execution:**
     a. **Calls Generic Handler**
        - **Calls:** `GetPaginatedHandler[entities.TrafficTicket](w, r, h.repo, PaginatedConfig{...})` → `internal/delivery/http/handler_helpers.go:85`
          - **Function:** Generic handler for paginated listing
          - **Execution:**
            - **Parse Query Parameters**
              - **Calls:** `r.URL.Query().Get("page")` - Standard library
              - **Calls:** `strconv.Atoi(pageStr)` - Standard library (converts to int)
              - **Calls:** `r.URL.Query().Get("perPage")` - Standard library
              - **Calls:** `strconv.Atoi(perPageStr)` - Standard library
            - **Call Repository**
              - **Type Assertion:** Checks if repo implements `PaginatableRepositoryWithGetPaginated`
              - **Calls:** `repo.GetPaginated(limit, offset)` → `internal/repository/traffic_ticket_repo.go:131`
                - **Function:** `TrafficTicketRepository.GetPaginated(limit, offset int) ([]entities.TrafficTicket, error)`
                - **Execution Flow:**
                  - **Calls:** `context.WithTimeout(context.Background(), config.GetQueryTimeout())` → `config/config.go:110`
                    - **Gets:** Query timeout from environment (default: 10 seconds)
                  - **Calls:** `r.db.QueryContext(ctx, query, limit, offset)` - Standard library
                    - **Executes:** SQL query with timeout context
                  - **Calls:** `handleQueryError(err)` → `internal/repository/pagination.go:26` (if error)
                    - **Checks:** For timeout and returns consistent error message
                - **Function:** `TrafficTicketRepository.GetPaginated(limit, offset int) ([]entities.TrafficTicket, error)`
                - **Execution:**
                  - **Create Context with Timeout**
                    - **Calls:** `context.WithTimeout(context.Background(), config.GetQueryTimeout())` - Standard library
                      - **Calls:** `config.GetQueryTimeout()` → `config/config.go:114`
                        - **Same flow as above** (gets timeout from environment, default: 10 seconds)
                      - **Returns:** Context with timeout and cancel function
                    - **Calls:** `defer cancel()` - Standard library
                  - **SQL Query:** `SELECT ... FROM traffic_tickets ORDER BY id LIMIT $1 OFFSET $2`
                  - **Calls:** `r.db.QueryContext(ctx, query, limit, offset)` - Standard library
                    - **Purpose:** Executes SQL query on PostgreSQL database with timeout
                    - **Database:** `trafficDB` connection (PostgreSQL "traffic_ticket" database)
                    - **Context:** Uses timeout context (query will be cancelled after timeout)
                    - **Returns:** Rows and error
                  - **Calls:** `handleQueryError(err)` → `internal/repository/pagination.go:26`
                    - **If error:** Returns timeout message or original error
                  - **Calls:** `rows.Close()` - Standard library (deferred)
                  - **Iterates:** `rows.Next()` - Standard library
                    - **Calls:** `rows.Scan(&t.ID, &t.DetectedSpeed, ...)` - Standard library
                      - **Purpose:** Scans each row into `entities.TrafficTicket` struct
                  - **Returns:** Slice of tickets and error
            - **Return Paginated Response**
              - **Calls:** `WritePaginatedResponse(w, tickets, page, perPage, "Success")` → `internal/delivery/http/response.go:54`
                - **Calls:** `w.Header().Set("Content-Type", "application/json")` - Standard library
                - **Calls:** `json.NewEncoder(w).Encode(Response{...})` - Standard library
              - **Returns:** (request complete)

**Complete Call Chain Summary:**
```
HTTP Request → Router → Rate Limit → JWT Auth → TrafficTicketHandler.GetPaginated → 
GetPaginatedHandler (generic) → TrafficTicketRepository.GetPaginated → 
PostgreSQL Database (trafficDB) → Response back through chain
```

---

## Environment Variables

The application uses the following environment variables:

### Database Configuration

**PostgreSQL (Golang DB & Traffic DB):**
- `DB_USER` - PostgreSQL username
- `DB_PASSWORD` - PostgreSQL password
- `DB_HOST` - PostgreSQL host (defaults to "localhost", auto-adjusts to "host.docker.internal" in Docker)
- `DB_PORT` - PostgreSQL port
- `DB_SSLMODE` - SSL mode for PostgreSQL connection

**MySQL:**
- `MYSQL_HOST` - MySQL host
- `MYSQL_PORT` - MySQL port
- `MYSQL_USER` - MySQL username
- `MYSQL_PASSWORD` - MySQL password
- `MYSQL_DATABASE` - MySQL database name

**Passenger Plane MySQL:**
- `PASSENGER_MYSQL_HOST` - Passenger MySQL host
- `PASSENGER_MYSQL_PORT` - Passenger MySQL port
- `PASSENGER_MYSQL_USER` - Passenger MySQL username
- `PASSENGER_MYSQL_PASSWORD` - Passenger MySQL password
- `PASSENGER_MYSQL_DATABASE` - Passenger MySQL database name

### Application Configuration

- `APP_PORT` - Port number for the HTTP server (default: "8080")
- `RATE_LIMIT_REQUESTS` - Maximum requests per minute (default: 100)
- `RATE_LIMIT_BURST` - Maximum burst capacity (default: 10)

### HTTP Server Timeout Configuration

- `HTTP_READ_TIMEOUT_SECONDS` - Maximum time to read request (default: 15 seconds)
- `HTTP_WRITE_TIMEOUT_SECONDS` - Maximum time to write response (default: 15 seconds)
- `HTTP_IDLE_TIMEOUT_SECONDS` - Maximum idle time on keep-alive connections (default: 60 seconds)

### Database Query Timeout Configuration

- `DB_QUERY_TIMEOUT_SECONDS` - Maximum time for a single database query (default: 10 seconds)

---

## Architecture Pattern

The application follows a **Clean Architecture** pattern with clear separation of concerns:

```
┌─────────────────────────────────────┐
│   HTTP Delivery Layer (Handlers)    │  ← Handles HTTP requests/responses
│   - user_handles.go                 │
│   - traffic_ticket_handler.go       │
│   - handler_helpers.go (generic)    │
│   - response.go (JSON formatting)   │
├─────────────────────────────────────┤
│   Use Cases Layer (Services)        │  ← Contains business logic
│   - user_service.go                 │
├─────────────────────────────────────┤
│   Repository Layer (Data Access)    │  ← Handles database operations
│   - user_repo.go                    │
│   - traffic_ticket_repo.go          │
│   - mysql_traffic_ticket_repo.go    │
│   - passenger_plane_repo.go         │
│   - pagination.go (helpers)         │
├─────────────────────────────────────┤
│   Database Layer                    │  ← PostgreSQL & MySQL databases
│   - golangDB (PostgreSQL)           │
│   - trafficDB (PostgreSQL)          │
│   - mysqlDB (MySQL)                  │
│   - passengerDB (MySQL)              │
└─────────────────────────────────────┘
```

**Benefits:**
- **Separation of Concerns:** Each layer has a single responsibility
- **Testability:** Each layer can be tested independently
- **Maintainability:** Changes in one layer don't affect others
- **Flexibility:** Easy to swap databases or change business logic

---

## Key Concepts

### Dependency Injection
The application uses dependency injection where:
- Repositories are injected into services/handlers
- Services are injected into handlers
- Database connections are injected into repositories
- This makes the code more testable and flexible

### Middleware Pattern
Middleware functions wrap HTTP handlers to add cross-cutting concerns:
- Rate limiting (protects against abuse)
- Authentication (ensures user is authorized)
- Logging (tracks requests)
- Error handling (standardizes error responses)

### Generic Handlers
The application uses Go generics (Go 1.18+) to create reusable handler functions:
- `CreateHandler[T]` - Generic create handler
- `ListHandler[T]` - Generic list handler
- `GetPaginatedHandler[T]` - Generic paginated handler
- **Why:** Reduces code duplication across similar handlers

### Connection Pooling
All database connections use optimized connection pool settings:
- **Max Open Connections:** 25
- **Max Idle Connections:** 10
- **Connection Max Lifetime:** 5 minutes
- **Connection Max Idle Time:** 1 minute
- **Why:** Prevents connection exhaustion and improves performance

### Timeout Protection
The application implements multiple layers of timeout protection:
- **HTTP Server Timeouts:** Prevent requests from hanging indefinitely
  - ReadTimeout: 15 seconds (default)
  - WriteTimeout: 15 seconds (default)
  - IdleTimeout: 60 seconds (default)
- **Database Query Timeouts:** Prevent slow queries from blocking connections
  - QueryTimeout: 10 seconds (default) per query
  - Applied via `context.WithTimeout()` in all repository methods
- **Why:** Ensures system responsiveness and prevents resource exhaustion

### Consolidated Helper Functions
The repository layer uses consolidated helper functions:
- **`normalizePagination()`** - Centralizes pagination validation (used by all List methods)
- **`handleQueryError()`** - Centralizes timeout error handling (used by all query methods)
- **Why:** Reduces code duplication and ensures consistent behavior

---

## SQLC Code Generation Workflow

This project uses **sqlc** to automatically generate type-safe Go code from SQL queries. This eliminates redundancy in both entities and repositories.

### What is SQLC?

SQLC is a code generator that:
- ✅ Reads SQL schema files and generates Go structs automatically
- ✅ Reads SQL query files and generates Go functions automatically
- ✅ Eliminates all field-by-field assignments in repositories
- ✅ Provides compile-time type safety
- ✅ No runtime reflection overhead

### SQLC Project Structure

```
golang_daerah/
├── sqlc.yaml                    # SQLC configuration
├── sqlc/
│   ├── schema/                  # SQL schema files (table definitions)
│   │   ├── postgres/
│   │   ├── mysql/
│   │   ├── passenger/
│   │   └── laut/
│   └── queries/                 # SQL query files
│       ├── postgres/
│       ├── mysql/
│       ├── passenger/
│       └── laut/
└── internal/repository/
    └── generated/               # Auto-generated code (DO NOT EDIT)
        ├── postgres/
        ├── mysql/
        ├── passenger/
        └── laut/
```

### Installing SQLC

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

### Generating Code

**Windows:**
```bash
generate_sqlc.bat
```

**Or manually (works on all platforms):**
```bash
sqlc generate
```

**Note:** If you're using Docker or WSL (Windows Subsystem for Linux), you can use the manual command `sqlc generate` directly.

This will generate Go code in `internal/repository/generated/` based on your SQL files.

### SQLC Workflow

1. **Define Schema** - Create SQL schema file in `sqlc/schema/{database}/`
2. **Write Queries** - Create SQL query file in `sqlc/queries/{database}/`
3. **Generate Code** - Run `sqlc generate`
4. **Use Generated Code** - Import and use generated types/functions in repositories

### Example: How SQLC Generates Code

**SQL Query File** (`sqlc/queries/mysql/traffic_tickets.sql`):
```sql
-- name: InsertTrafficTicket :execresult
INSERT INTO traffic_tickets (detected_speed, legal_speed, ...)
VALUES (?, ?, ...);

-- name: ListTrafficTickets :many
SELECT * FROM traffic_tickets
ORDER BY id DESC
LIMIT ? OFFSET ?;
```

**Generated Go Code** (`internal/repository/generated/mysql/traffic_tickets.sql.go`):
```go
// Auto-generated by sqlc - DO NOT EDIT

type TrafficTicket struct {
    ID                         int32   `json:"id"`
    DetectedSpeed              float64 `json:"detected_speed"`
    // ... all fields automatically generated
}

type InsertTrafficTicketParams struct {
    DetectedSpeed float64 `json:"detected_speed"`
    // ... all fields automatically generated
}

func (q *Queries) InsertTrafficTicket(ctx context.Context, arg InsertTrafficTicketParams) (sql.Result, error) {
    // All field assignments automatically generated
    return q.db.ExecContext(ctx, insertTrafficTicket, arg.DetectedSpeed, ...)
}

func (q *Queries) ListTrafficTickets(ctx context.Context, arg ListTrafficTicketsParams) ([]TrafficTicket, error) {
    // All scanning code automatically generated
    // ...
}
```

**Repository Usage** (after refactoring):
```go
func (r *MySQLTrafficTicketRepository) Insert(ticket *entities.TrafficTicket) error {
    ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
    defer cancel()

    // Convert to generated params struct
    params := generated.InsertTrafficTicketParams{
        DetectedSpeed: ticket.DetectedSpeed,
        // ... map all fields
    }

    // Use generated function - no manual field assignments!
    _, err := r.queries.InsertTrafficTicket(ctx, params)
    return handleQueryError(err)
}
```

### Benefits of SQLC

- **Zero Redundancy:** All code is generated from SQL
- **Type Safety:** Compile-time checking, no runtime errors
- **Maintainability:** Change SQL → regenerate → done
- **Performance:** No reflection overhead
- **Multi-Database:** One SQL file per database engine

---

## Step-by-Step Guide: Creating New API Endpoints (With SQLC)

This guide shows you how to create new API endpoints using **sqlc** for automatic code generation. We'll create two endpoints:
- `GET /api/terminal/laut` - List port/terminal data with pagination
- `POST /api/terminal/tambah_Laut` - Add new port/terminal data

### Prerequisites
- Database table `Laut` exists in the `terminal` database
- Table has all required columns (see schema below)
- Database connection is configured in `config/config.go`
- **sqlc is installed** (`go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`)

---

### Step 1: Create SQL Schema File

**File:** `sqlc/schema/laut/laut.sql`

Create the SQL schema file that defines your table structure. SQLC will generate Go structs from this:

```sql
CREATE TABLE Laut (
    id INT AUTO_INCREMENT PRIMARY KEY,
    port_name VARCHAR(100) NOT NULL,
    port_code VARCHAR(20) NOT NULL,
    port_address TEXT NOT NULL,
    city VARCHAR(100) NOT NULL,
    province VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL,
    operator_name VARCHAR(100) NOT NULL,
    operator_contact VARCHAR(100) NOT NULL,
    harbor_master_name VARCHAR(100) NOT NULL,
    harbor_master_id VARCHAR(50) NOT NULL,
    harbor_master_rank VARCHAR(50) NOT NULL,
    harbor_master_office_address TEXT NOT NULL,
    number_of_piers INT NOT NULL,
    main_pier_length DOUBLE NOT NULL,
    max_ship_draft DOUBLE NOT NULL,
    max_ship_length DOUBLE NOT NULL,
    terminal_capacity_passenger INT NOT NULL,
    terminal_capacity_cargo INT NOT NULL,
    operational_hours VARCHAR(100) NOT NULL,
    emergency_contact VARCHAR(100) NOT NULL,
    security_office_name VARCHAR(100) NOT NULL,
    security_officer_id VARCHAR(50) NOT NULL,
    security_level VARCHAR(50) NOT NULL,
    checkin_counter_count INT NOT NULL,
    special_facilities TEXT NOT NULL
);
```

**Note:** SQLC will automatically generate a `Laut` struct with all fields and JSON tags from this schema!

**Why:** The entity struct defines the data structure that will be used throughout the application layers.

---

### Step 2: Initialize Database Connection (if needed)

**File:** `config/config.go`

If you need a new database connection for the `terminal` database, add an initialization function:

```go
// InitTerminalDB initializes connection to PostgreSQL "terminal" database
func InitTerminalDB() *sql.DB {
	return InitDB("terminal")
}
```

**File:** `cmd/app/main.go`

Add database initialization in the `main()` function (around line 80-100):

```go
// Initialize terminal database connection
terminalDB := initDB("terminal", config.InitTerminalDB)
defer terminalDB.Close()
```

**Why:** Each database needs its own connection pool. The `initDB()` wrapper ensures consistent logging and error handling.

---

### Step 3: Create the Repository

**File:** `internal/repository/laut_repo.go`

Create the repository with three methods: `Insert`, `List`, and `GetPaginated`:

```go
package repository

import (
	"context"
	"database/sql"
	"golang_daerah/config"
	"golang_daerah/internal/entities"
)

type LautRepository struct {
	db *sql.DB
}

func NewLautRepository(db *sql.DB) *LautRepository {
	return &LautRepository{db: db}
}

// Insert adds a new port/terminal record to the database
func (r *LautRepository) Insert(laut *lautgen.Laut) error {
	// Create context with timeout to prevent queries from hanging indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	defer cancel()

	query := `
		INSERT INTO Laut (
			port_name, port_code, port_address, city, province, country,
			operator_name, operator_contact, harbor_master_name, harbor_master_id,
			harbor_master_rank, harbor_master_office_address, number_of_piers,
			main_pier_length, max_ship_draft, max_ship_length,
			terminal_capacity_passenger, terminal_capacity_cargo, operational_hours,
			emergency_contact, security_office_name, security_officer_id,
			security_level, checkin_counter_count, special_facilities
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25)
	`

	_, err := r.db.ExecContext(ctx, query,
		laut.PortName,
		laut.PortCode,
		laut.PortAddress,
		laut.City,
		laut.Province,
		laut.Country,
		laut.OperatorName,
		laut.OperatorContact,
		laut.HarborMasterName,
		laut.HarborMasterID,
		laut.HarborMasterRank,
		laut.HarborMasterOfficeAddress,
		laut.NumberOfPiers,
		laut.MainPierLength,
		laut.MaxShipDraft,
		laut.MaxShipLength,
		laut.TerminalCapacityPassenger,
		laut.TerminalCapacityCargo,
		laut.OperationalHours,
		laut.EmergencyContact,
		laut.SecurityOfficeName,
		laut.SecurityOfficerID,
		laut.SecurityLevel,
		laut.CheckinCounterCount,
		laut.SpecialFacilities,
	)

	return handleQueryError(err)
}

// List returns port/terminal records with pagination controls (ordered by id DESC)
func (r *LautRepository) List(limit, offset int) ([]*lautgen.Laut, error) {
	// Create context with timeout to prevent queries from hanging indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	defer cancel()

	limit, offset = normalizePagination(limit, offset)

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, port_name, port_code, port_address, city, province, country,
		       operator_name, operator_contact, harbor_master_name, harbor_master_id,
		       harbor_master_rank, harbor_master_office_address, number_of_piers,
		       main_pier_length, max_ship_draft, max_ship_length,
		       terminal_capacity_passenger, terminal_capacity_cargo, operational_hours,
		       emergency_contact, security_office_name, security_officer_id,
		       security_level, checkin_counter_count, special_facilities
		FROM Laut
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)

	if err != nil {
		return nil, handleQueryError(err)
	}
	defer rows.Close()

	var lauts []*lautgen.Laut
	for rows.Next() {
		l := &lautgen.Laut{}
		if err := rows.Scan(
			&l.ID,
			&l.PortName,
			&l.PortCode,
			&l.PortAddress,
			&l.City,
			&l.Province,
			&l.Country,
			&l.OperatorName,
			&l.OperatorContact,
			&l.HarborMasterName,
			&l.HarborMasterID,
			&l.HarborMasterRank,
			&l.HarborMasterOfficeAddress,
			&l.NumberOfPiers,
			&l.MainPierLength,
			&l.MaxShipDraft,
			&l.MaxShipLength,
			&l.TerminalCapacityPassenger,
			&l.TerminalCapacityCargo,
			&l.OperationalHours,
			&l.EmergencyContact,
			&l.SecurityOfficeName,
			&l.SecurityOfficerID,
			&l.SecurityLevel,
			&l.CheckinCounterCount,
			&l.SpecialFacilities,
		); err != nil {
			return nil, err
		}
		lauts = append(lauts, l)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lauts, nil
}

// GetPaginated returns port/terminal records with pagination controls, ordered by id ASC (starting from ID 1)
func (r *LautRepository) GetPaginated(limit, offset int) ([]lautgen.Laut, error) {
	// Create context with timeout to prevent queries from hanging indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), config.GetQueryTimeout())
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, port_name, port_code, port_address, city, province, country,
		       operator_name, operator_contact, harbor_master_name, harbor_master_id,
		       harbor_master_rank, harbor_master_office_address, number_of_piers,
		       main_pier_length, max_ship_draft, max_ship_length,
		       terminal_capacity_passenger, terminal_capacity_cargo, operational_hours,
		       emergency_contact, security_office_name, security_officer_id,
		       security_level, checkin_counter_count, special_facilities
		FROM Laut
		ORDER BY id
		LIMIT $1 OFFSET $2
	`, limit, offset)

	if err != nil {
		return nil, handleQueryError(err)
	}
	defer rows.Close()

	var lauts []lautgen.Laut
	for rows.Next() {
		var l lautgen.Laut
		if err := rows.Scan(
			&l.ID,
			&l.PortName,
			&l.PortCode,
			&l.PortAddress,
			&l.City,
			&l.Province,
			&l.Country,
			&l.OperatorName,
			&l.OperatorContact,
			&l.HarborMasterName,
			&l.HarborMasterID,
			&l.HarborMasterRank,
			&l.HarborMasterOfficeAddress,
			&l.NumberOfPiers,
			&l.MainPierLength,
			&l.MaxShipDraft,
			&l.MaxShipLength,
			&l.TerminalCapacityPassenger,
			&l.TerminalCapacityCargo,
			&l.OperationalHours,
			&l.EmergencyContact,
			&l.SecurityOfficeName,
			&l.SecurityOfficerID,
			&l.SecurityLevel,
			&l.CheckinCounterCount,
			&l.SpecialFacilities,
		); err != nil {
			return nil, err
		}
		lauts = append(lauts, l)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lauts, nil
}
```

**Key Points:**
- Use `context.WithTimeout()` with `config.GetQueryTimeout()` for all queries
- Use `normalizePagination()` for `List()` method
- Use `handleQueryError()` for consistent error handling
- `List()` returns `[]*lautgen.Laut` (pointers) with `ORDER BY id DESC`
- `GetPaginated()` returns `[]lautgen.Laut` (values) with `ORDER BY id ASC`
- Use `$1, $2, ...` placeholders for PostgreSQL (use `?` for MySQL)

**Why:** The repository layer handles all database operations and provides a clean interface for the handler layer.

---

### Step 4: Create the Handler

**File:** `internal/delivery/http/laut_handler.go`

Create the handler using the generic helper functions:

```go
package http

import (
	"net/http"
	"golang_daerah/internal/entities"
	"golang_daerah/internal/repository"
)

type LautHandler struct {
	repo *repository.LautRepository
}

func NewLautHandler(repo *repository.LautRepository) *LautHandler {
	return &LautHandler{repo: repo}
}

// Create handles adding new port/terminal data
func (h *LautHandler) Create(w http.ResponseWriter, r *http.Request) {
	CreateHandler(w, r, h.repo, CreateConfig{
		RequirePOST:        true,
		InvalidBodyMessage: "Invalid request body: ",
		InsertErrorMessage: "Failed to insert port/terminal data: %v",
		SuccessMessage:     "%d port/terminal record(s) created successfully",
	})
}

// GetPaginated handles paginated listing of port/terminal data
func (h *LautHandler) GetPaginated(w http.ResponseWriter, r *http.Request) {
	GetPaginatedHandler[lautgen.Laut](w, r, h.repo, PaginatedConfig{
		DefaultPerPage:  10,
		UseGetPaginated: true, // Use GetPaginated method which orders by id ASC (starting from ID 1)
		ErrorMessage:    "Failed to get port/terminal data: %v",
	})
}
```

**Key Points:**
- `CreateHandler()` handles POST requests with automatic JSON parsing and validation
- `GetPaginatedHandler()` handles GET requests with query parameter parsing (`page`, `perPage`)
- Both functions automatically apply JWT authentication and rate limiting when registered as protected routes
- `UseGetPaginated: true` ensures ascending order (starting from ID 1)

**Why:** The handler layer uses generic helper functions to reduce code duplication and ensure consistent behavior.

---

### Step 5: Register Repository, Handler, and Routes

**File:** `cmd/app/main.go`

Add the repository and handler initialization in the "PHASE 2: APPLICATION LAYER SETUP" section (around line 110-130):

```go
// --- Laut (Port/Terminal) Layer (uses terminal PostgreSQL database) ---
// Request Flow: Handler → Repository → Database
lautRepo := repository.NewLautRepository(terminalDB)
lautHandler := httpDelivery.NewLautHandler(lautRepo)
```

Add route registration in the "PHASE 5: ROUTE REGISTRATION" section (around line 180-190):

```go
// Terminal/Port Endpoints
routes.Protected("/api/terminal/tambah_Laut", lautHandler.Create)      // POST: Add port/terminal data
routes.Protected("/api/terminal/laut", lautHandler.GetPaginated)       // GET: List port/terminal data (paginated)
```

**Key Points:**
- `routes.Protected()` automatically applies:
  - Rate limiting middleware
  - JWT authentication middleware
  - Both are applied in sequence before the handler
- The handler methods (`Create`, `GetPaginated`) are passed as function references
- Routes are registered after middleware configuration

**Why:** The `RouteRegistrar` struct ensures consistent middleware application and simplifies route registration.

---

### Step 6: Complete Request Flow

When a request comes in, here's the complete flow:

#### For `POST /api/terminal/tambah_Laut`:

1. **HTTP Request** arrives at server
2. **Router** (`http.ServeMux`) matches route `/api/terminal/tambah_Laut`
3. **Rate Limit Middleware** checks if IP has exceeded request limit
   - If exceeded: Returns `429 Too Many Requests`
   - If OK: Continues to next middleware
4. **JWT Authentication Middleware** validates `Authorization: Bearer <token>` header
   - If invalid: Returns `401 Unauthorized`
   - If valid: Continues to handler
5. **Handler** (`LautHandler.Create`) receives request
6. **Generic Handler** (`CreateHandler`) processes request:
   - Validates HTTP method (must be POST)
   - Parses JSON body into `lautgen.Laut` struct
   - Calls `lautRepo.Insert()`
7. **Repository** (`LautRepository.Insert`) executes database query:
   - Creates context with timeout
   - Executes INSERT query
   - Handles errors with `handleQueryError()`
8. **Response** is sent back through the chain

#### For `GET /api/terminal/laut?page=1&perPage=10`:

1. **HTTP Request** arrives at server
2. **Router** matches route `/api/terminal/laut`
3. **Rate Limit Middleware** checks request limit
4. **JWT Authentication Middleware** validates token
5. **Handler** (`LautHandler.GetPaginated`) receives request
6. **Generic Handler** (`GetPaginatedHandler`) processes request:
   - Parses `page` and `perPage` from query parameters (or headers: `Page`, `PerPage`)
   - Calculates offset: `offset = (page - 1) * perPage`
   - Calls `lautRepo.GetPaginated(perPage, offset)`
7. **Repository** (`LautRepository.GetPaginated`) executes database query:
   - Creates context with timeout
   - Executes SELECT query with `ORDER BY id ASC`
   - Returns results starting from ID 1
8. **Response** is formatted as JSON with pagination metadata:
   ```json
   {
     "status": true,
     "data": [...],
     "page": 1,
     "perPage": 10,
     "message": "OK"
   }
   ```

---

### Step 9: Testing the Endpoints

#### Test `POST /api/terminal/tambah_Laut`:

**Request:**
```http
POST /api/terminal/tambah_Laut
Authorization: Bearer <your_jwt_token>
Content-Type: application/json

{
  "port_name": "Jakarta Port",
  "port_code": "JKT001",
  "port_address": "Jl. Pelabuhan No. 1",
  "city": "Jakarta",
  "province": "DKI Jakarta",
  "country": "Indonesia",
  "operator_name": "PT Pelabuhan Indonesia",
  "operator_contact": "+62-21-12345678",
  "harbor_master_name": "John Doe",
  "harbor_master_id": "HM001",
  "harbor_master_rank": "Captain",
  "harbor_master_office_address": "Jl. Kantor Pelabuhan",
  "number_of_piers": 5,
  "main_pier_length": 500.5,
  "max_ship_draft": 15.0,
  "max_ship_length": 300.0,
  "terminal_capacity_passenger": 1000,
  "terminal_capacity_cargo": 5000,
  "operational_hours": "24/7",
  "emergency_contact": "+62-21-99999999",
  "security_office_name": "Security Office",
  "security_officer_id": "SO001",
  "security_level": "High",
  "checkin_counter_count": 10,
  "special_facilities": "Crane, Container Yard"
}
```

**Response:**
```json
{
  "status": true,
  "data": {
    "id": 1,
    "port_name": "Jakarta Port",
    ...
  },
  "message": "1 port/terminal record(s) created successfully"
}
```

#### Test `GET /api/terminal/laut`:

**Request:**
```http
GET /api/terminal/laut?page=1&perPage=10
Authorization: Bearer <your_jwt_token>
```

**Or using headers:**
```http
GET /api/terminal/laut
Authorization: Bearer <your_jwt_token>
Page: 1
PerPage: 10
```

**Response:**
```json
{
  "status": true,
  "data": [
    {
      "id": 1,
      "port_name": "Jakarta Port",
      ...
    },
    ...
  ],
  "page": 1,
  "perPage": 10,
  "message": "OK"
}
```

---

### Summary Checklist (SQLC Workflow)

- [ ] **Create SQL schema file** in `sqlc/schema/laut/laut.sql`
- [ ] **Create SQL query file** in `sqlc/queries/laut/laut.sql` with:
  - [ ] `InsertLaut :execresult` query
  - [ ] `ListLauts :many` query (with `ORDER BY id DESC`)
  - [ ] `GetPaginatedLauts :many` query (with `ORDER BY id ASC`)
- [ ] **Generate sqlc code** by running `sqlc generate` (or `generate_sqlc.bat`)
- [ ] **Use generated `lautgen.Laut` type** throughout handlers and repositories (no separate entity struct)
- [ ] **Create repository** in `internal/repository/laut_repo.go` backed by SQLC:
  - [ ] Wraps generated `lautgen.Queries`
  - [ ] Converts `lautgen.Laut` into `InsertLautParams` (handled inside repository)
- [ ] **Initialize database connection** in `config/config.go` (if new database)
- [ ] **Add database initialization** in `cmd/app/main.go`
- [ ] **Create handler** in `internal/delivery/http/laut_handler.go` with:
  - [ ] `Create()` method using `CreateHandler()`
  - [ ] `GetPaginated()` method using `GetPaginatedHandler()`
- [ ] **Register repository and handler** in `cmd/app/main.go` (PHASE 2)
- [ ] **Register routes** in `cmd/app/main.go` (PHASE 5) using `routes.Protected()`
- [ ] **Test endpoints** with Postman or curl

**SQLC Benefits:**
- ✅ **No manual repository code** - All generated from SQL
- ✅ **No field-by-field assignments** - SQLC handles everything
- ✅ **Type-safe** - Compile-time checking
- ✅ **Maintainable** - Change SQL → regenerate → done

**All endpoints automatically include:**
- ✅ Rate limiting (configurable via `RATE_LIMIT_REQUESTS` and `RATE_LIMIT_BURST`)
- ✅ JWT authentication (requires `Authorization: Bearer <token>` header)
- ✅ Pagination support (via query parameters or headers)
- ✅ HTTP server timeouts (prevents hanging requests)
- ✅ Database query timeouts (prevents slow queries from blocking)

---

## Summary

This application is a well-structured Go REST API that:
1. Manages multiple database connections (PostgreSQL and MySQL)
2. Implements user authentication with JWT tokens
3. Provides rate limiting to prevent abuse
4. Follows clean architecture principles
5. Uses environment variables for configuration
6. Handles traffic ticket operations across different databases
7. Uses generic handlers to reduce code duplication
8. Implements proper connection pooling for database performance
9. Implements HTTP server timeouts to prevent hanging requests
10. Implements database query timeouts to prevent slow queries from blocking connections
11. Uses consolidated helper functions for pagination validation and error handling

The `main()` function orchestrates the entire application setup, creating a complete dependency chain from HTTP handlers through services and repositories to database connections. Every function call is traceable through the layers, ensuring maintainability and testability.

**Key Consolidation Achievements:**
- **Route Registration:** Consolidated via `RouteRegistrar` struct
- **JSON Responses:** Consolidated via `response.go` helper functions
- **Handler Methods:** Consolidated via generic handlers in `handler_helpers.go`
- **Database Initialization:** Consolidated via `initDB()` wrapper and shared `configureConnectionPool()`
- **MySQL Database Init:** Consolidated via `initMySQLDB()` helper
- **Pagination Validation:** Consolidated via `normalizePagination()` in `pagination.go`
- **Query Error Handling:** Consolidated via `handleQueryError()` in `pagination.go`
- **Context Timeouts:** Applied consistently across all repository methods
