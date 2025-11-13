package main

import "net/http"

func createRoutes(
	companyH *companyHandler,
	authH *authHandler,
	eventH *eventHandler,
) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /company/{id}", companyH.handleGetCompany)
	mux.HandleFunc("GET /company", companyH.handleGetCompanies)
	mux.HandleFunc("POST /company", companyH.handlePostCompany)
	mux.HandleFunc("PUT /company/{id}", companyH.handlePutCompany)
	mux.HandleFunc("DELETE /company/{id}", companyH.handleDeleteCompany)

	mux.HandleFunc("POST /auth/register", authH.handleRegister)
	mux.HandleFunc("POST /auth/login", authH.handleLogin)

	mux.HandleFunc("GET /event", eventH.handleGetUpcomingEvents)
	mux.Handle("GET /event/all", AuthMiddleware(http.HandlerFunc(eventH.handleGetAllEvents)))

	return mux
}
