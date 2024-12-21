package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	postgresRepo "github.com/bimbims125/clean-arch/internal/repository/postgresql"
	"github.com/bimbims125/clean-arch/internal/rest"
	"github.com/bimbims125/clean-arch/internal/rest/middleware"
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
	// Prepare database
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	val := url.Values{}
	val.Add("sslmode", "disable")
	val.Add("timezone", "Asia/Jakarta")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s", dbUser, dbPass, dbHost, dbPort, dbName, val.Encode())

	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to open connection to database", err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal("failed to ping database ", err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal("got error when closing the DB connection", err)
		}
	}()

	// Initialize repository
	userRepo := postgresRepo.NewUserRepository(dbConn)
	productRepo := postgresRepo.NewProductRepository(dbConn)
	categoryRepo := postgresRepo.NewCategoryRepository(dbConn)

	// Create a main router
	r := mux.NewRouter()

	// Create a subrouter with prefix "/api/v1"
	apiRouter := r.PathPrefix("/api/v1").Subrouter()

	// Register user handlers to the subrouter
	rest.NewUserHandler(apiRouter, userRepo)
	rest.NewProductHandler(apiRouter, productRepo)
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
