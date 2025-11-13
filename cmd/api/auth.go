package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/mthsgimenez/participe/internal/auth"
	"github.com/mthsgimenez/participe/internal/company"
	"github.com/mthsgimenez/participe/internal/user"
)

type UserRegisterDTO struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	CompanyId int    `json:"company_id"`
	Name      string `json:"name"`
}

type UserLoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserLoginDTO) Validate() (problems map[string]string) {
	problems = map[string]string{}

	if strings.TrimSpace(u.Email) == "" {
		problems["email"] = "email cant be empty"
	}

	if strings.TrimSpace(u.Password) == "" {
		problems["password"] = "password cant be empty"
	}

	return
}

func (u *UserRegisterDTO) Validate() (problems map[string]string) {
	problems = map[string]string{}

	if strings.TrimSpace(u.Email) == "" {
		problems["email"] = "email cant be empty"
	}

	if strings.TrimSpace(u.Name) == "" {
		problems["name"] = "name cant be empty"
	}

	if strings.TrimSpace(u.Password) == "" {
		problems["password"] = "password cant be empty"
	}

	if u.CompanyId == 0 {
		problems["company_id"] = "company_id cant be empty"
	}

	return
}

type authHandler struct {
	userService    *user.Service
	companyService *company.Service
}

func NewAuthHandler(us *user.Service, cs *company.Service) *authHandler {
	return &authHandler{us, cs}
}

func (h *authHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	d, problems, err := BindJSONValid[*UserLoginDTO](r)
	if err != nil {
		if len(problems) > 0 {
			RespondJSONErrorWithProblems(w, "invalid request body", http.StatusBadRequest, problems)
			return
		}

		RespondJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u, err := h.userService.GetUserByEmail(d.Email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			RespondJSONError(w, "email or password invalid", http.StatusUnauthorized)
			return
		}

		RespondJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !u.CheckPassword(d.Password) {
		RespondJSONError(w, "email or password invalid", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT(d.Email)
	if err != nil {
		RespondJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RespondJSON(w, fmt.Sprintf("Bearer %s", token), http.StatusOK)
}

func (h *authHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	u, problems, err := BindJSONValid[*UserRegisterDTO](r)
	if err != nil {
		if len(problems) > 0 {
			RespondJSONErrorWithProblems(w, "invalid request body", http.StatusBadRequest, problems)
			return
		}

		RespondJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cmp, err := h.companyService.GetCompany(u.CompanyId)
	if err != nil {
		RespondJSONError(w, "invalid company_id", http.StatusBadRequest)
		return
	}

	newUser := &user.User{
		Role:    user.ROLE_USER,
		Company: *cmp,
		Name:    u.Name,
		Email:   u.Email,
	}

	if err := newUser.SetPassword(u.Password); err != nil {
		RespondJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.userService.CreateUser(newUser)
	if err != nil {
		RespondJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RespondJSON(w, "user registered", http.StatusOK)
}
