package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mthsgimenez/participe/internal/company"
)

type companyHandler struct {
	companyService *company.Service
}

func newCompanyHandler(service *company.Service) *companyHandler {
	return &companyHandler{service}
}

func (c *companyHandler) handleGetCompany(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		RespondJSONError(w, "id must be an int", http.StatusBadRequest)
		return
	}

	cmp, err := c.companyService.GetCompany(id)
	if err != nil {
		if errors.Is(err, company.ErrCompanyNotFound) {
			RespondJSONError(w, "company not found", http.StatusNotFound)
			return
		}

		RespondJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RespondJSON(w, cmp, http.StatusOK)
}

func (c *companyHandler) handleGetCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := c.companyService.GetCompanies()
	if err != nil {
		RespondJSONError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	RespondJSON(w, companies, http.StatusOK)
}

func (c *companyHandler) handlePostCompany(w http.ResponseWriter, r *http.Request) {
	company, problems, err := BindJSONValid[*company.Company](r)
	if err != nil {
		if len(problems) > 0 {
			RespondJSONErrorWithProblems(w, "invalid request body", http.StatusBadRequest, problems)
			return
		}

		RespondJSONError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	newCompany, err := c.companyService.CreateCompany(company)
	if err != nil {
		RespondJSONError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	RespondJSON(w, newCompany, http.StatusCreated)
}

func (c *companyHandler) handleDeleteCompany(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		RespondJSONError(w, "id must be an int", http.StatusBadRequest)
		return
	}

	if err := c.companyService.DeleteCompany(id); err != nil {
		if errors.Is(err, company.ErrCompanyNotFound) {
			RespondJSONError(w, fmt.Sprintf("company with id %d does not exist", id), http.StatusNotFound)
			return
		}

		RespondJSONError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *companyHandler) handlePutCompany(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		RespondJSONError(w, "id must be an int", http.StatusBadRequest)
		return
	}

	cmp, problems, err := BindJSONValid[*company.Company](r)
	if err != nil {
		if len(problems) > 0 {
			RespondJSONErrorWithProblems(w, "invalid request body", http.StatusBadRequest, problems)
			return
		}

		RespondJSONError(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	cmp.Id = id

	updatedCompany, err := c.companyService.UpdateCompany(id, cmp)
	if err != nil {
		RespondJSONError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	RespondJSON(w, updatedCompany, http.StatusOK)
}
