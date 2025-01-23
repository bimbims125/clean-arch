package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	mysqlRepo "github.com/bimbims125/clean-arch/internal/repository/mysql"
	postgresRepo "github.com/bimbims125/clean-arch/internal/repository/postgresql"
	"github.com/bimbims125/clean-arch/internal/rest"
	"github.com/bimbims125/clean-arch/internal/rest/middleware"
	_ "github.com/go-sql-driver/mysql" // Import driver MySQL
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	defaultTimeout = 30
	defaultAddress = ":3300"
)

func init() {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// Load configuration from environment variables
	dbType := os.Getenv("DB_TYPE")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	val := url.Values{}
	val.Add("parseTime", "true")
	val.Add("loc", "Asia/Jakarta")

	var dbConn *sql.DB
	var err error
	var userRepo rest.UserService // General interface for both repositories
	var categoryRepo rest.CategoryService

	// Choose database from .env setup DB_TYPE
	switch dbType {
	case "postgres":
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s", dbUser, dbPass, dbHost, dbPort, dbName, val.Encode())
		dbConn, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Fatal("failed to open connection to Postgres: ", err)
		}
		userRepo = postgresRepo.NewPostgresUserRepository(dbConn)
		// categoryRepo = postgresRepo.NewPostgresCategoryRepository(dbConn)
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", dbUser, dbPass, dbHost, dbPort, dbName, val.Encode())
		dbConn, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatal("failed to open connection to MySQL: ", err)
		}
		userRepo = mysqlRepo.NewMySQLUserRepository(dbConn)
		categoryRepo = mysqlRepo.NewMySQLCategoryRepository(dbConn)
	default:
		log.Fatal("unsupported database type. Please set DB_TYPE to 'postgres' or 'mysql'")
	}

	// Check DB connection
	err = dbConn.Ping()
	if err != nil {
		log.Fatal("failed to ping database: ", err)
	}

	defer dbConn.Close()

	// Create a main router
	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api/v1").Subrouter()

	// Register user handlers to the subrouter
	rest.NewUserHandler(apiRouter, userRepo)
	rest.NewCategoryHandler(apiRouter, categoryRepo)

	// Wrap the main router with CORS middleware
	corsWrappedRouter := middleware.CORSMiddleware(r)

	// Start HTTP server
	addr := os.Getenv("APP_ADDRESS")
	if addr == "" {
		addr = defaultAddress
	}
	fmt.Printf("Server is running at %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, corsWrappedRouter))
}
