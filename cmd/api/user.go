package main

import (
	"net/http"

	"github.com/mthsgimenez/participe/internal/user"
)

type userHandler struct {
	userService *user.Service
}

func NewUserHandler(s *user.Service) *userHandler {
	return &userHandler{s}
}

func (h *userHandler) handleGetMe(w http.ResponseWriter, r *http.Request) {
	claims := GetUserClaims(r)
	if claims == nil {
		RespondJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	u, err := h.userService.GetUserByEmail(claims.Email)
	if err != nil {
		RespondJSONError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	RespondJSON(w, u, http.StatusOK)
}
