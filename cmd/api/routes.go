package main

import "net/http"

func createRoutes(companyH *companyHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /company/{id}", companyH.handleGetCompany)
	mux.HandleFunc("GET /company", companyH.handleGetCompanies)
	mux.HandleFunc("POST /company", companyH.handlePostCompany)
	mux.HandleFunc("PUT /company/{id}", companyH.handlePutCompany)
	mux.HandleFunc("DELETE /company/{id}", companyH.handleDeleteCompany)

	return mux
}
