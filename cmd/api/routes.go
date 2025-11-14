package main

import "net/http"

func createRoutes(
	companyH *companyHandler,
	authH *authHandler,
	eventH *eventHandler,
	userH *userHandler,
) *http.ServeMux {
	root := http.NewServeMux()

	// Public routes
	root.HandleFunc("POST /auth/register", authH.handleRegister)
	root.HandleFunc("POST /auth/login", authH.handleLogin)

	// Private routes
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("GET /company/{id}", companyH.handleGetCompany)
	protectedMux.HandleFunc("GET /company", companyH.handleGetCompanies)
	protectedMux.HandleFunc("POST /company", companyH.handlePostCompany)
	protectedMux.HandleFunc("PUT /company/{id}", companyH.handlePutCompany)
	protectedMux.HandleFunc("DELETE /company/{id}", companyH.handleDeleteCompany)

	protectedMux.HandleFunc("GET /event", eventH.handleGetUpcomingEvents)
	protectedMux.HandleFunc("GET /event/all", eventH.handleGetAllEvents)
	protectedMux.HandleFunc("GET /event/{id}", eventH.handleGetEvent)
	protectedMux.HandleFunc("GET /event/{id}/checkin", eventH.handleGetCheckins)
	protectedMux.HandleFunc("POST /event/{id}/checkin", eventH.handlePostCheckin)
	protectedMux.HandleFunc("POST /event", eventH.handlePostEvent)

	protectedMux.HandleFunc("GET /user/{id}", userH.handleGetUser)

	protected := AuthMiddleware(protectedMux)

	root.Handle("/", protected)

	return root
}
