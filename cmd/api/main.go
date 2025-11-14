package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mthsgimenez/participe/internal/company"
	"github.com/mthsgimenez/participe/internal/db"
	"github.com/mthsgimenez/participe/internal/env"
	"github.com/mthsgimenez/participe/internal/event"
	"github.com/mthsgimenez/participe/internal/user"
)

var (
	companyRepository company.Repository
	companyService    *company.Service
	companyH          *companyHandler
	userRepository    user.Repository
	userService       *user.Service
	userH             *userHandler
	authH             *authHandler
	eventRepository   event.Repository
	eventService      *event.Service
	eventH            *eventHandler
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

	// Dependencies

	companyRepository = company.NewRepositoryPostgres(conn)
	companyService = company.NewService(companyRepository)
	companyH = newCompanyHandler(companyService)

	userRepository = user.NewRepositoryPostgres(conn)
	userService = user.NewService(userRepository, companyRepository)
	userH = NewUserHandler(userService)

	authH = NewAuthHandler(userService, companyService)

	eventRepository = event.NewRepositoryPostgres(conn)
	eventService = event.NewService(eventRepository)
	eventH = NewEventHandler(eventService, userService)

	mux := createRoutes(companyH, authH, eventH, userH)
	muxWithCors := CorsMiddleware(mux)

	// -------------

	port := env.GetStringFallback("PORT", "8000")
	fmt.Printf("Server started at localhost:%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), muxWithCors))
}
