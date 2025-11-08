package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mthsgimenez/participe/internal/company"
	"github.com/mthsgimenez/participe/internal/db"
	"github.com/mthsgimenez/participe/internal/env"
)

var (
	companyRepository company.Repository
	companyService    *company.Service
	cmpnyHandler      *companyHandler
)

func main() {
	dbUser := env.GetStringFallback("DATABASE_USER", "postgres")
	dbPassword := env.GetStringFallback("DATABASE_PASSWORD", "postgres")
	dbHost := env.GetStringFallback("DATABASE_HOST", "localhost")
	dbPort := env.GetStringFallback("DATABASE_PORT", "5432")
	dbName := env.GetStringFallback("DATABASE_NAME", "postgres")

	connString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	conn, err := db.ConnectToDB(connString)
	if err != nil {
		panic("error connecting to database: " + err.Error())
	}
	defer conn.Close()

	companyRepository = company.NewRepositoryPostgres(conn)
	companyService = company.NewService(companyRepository)
	cmpnyHandler = newCompanyHandler(companyService)

	mux := createRoutes(cmpnyHandler)

	port := env.GetStringFallback("PORT", "8000")
	fmt.Printf("Server started at localhost:%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}
