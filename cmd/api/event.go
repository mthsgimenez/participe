package main

import (
	"net/http"

	"github.com/mthsgimenez/participe/internal/event"
	"github.com/mthsgimenez/participe/internal/user"
)

type eventHandler struct {
	eventService *event.Service
	userService  *user.Service
}

func NewEventHandler(e *event.Service, u *user.Service) *eventHandler {
	return &eventHandler{e, u}
}

func (h *eventHandler) handleGetUpcomingEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.eventService.GetUpcomingEvents()
	if err != nil {
		RespondJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RespondJSON(w, events, http.StatusOK)
}

func (h *eventHandler) handleGetAllEvents(w http.ResponseWriter, r *http.Request) {
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

	if u.Role != user.ROLE_ADMIN {
		RespondJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	events, err := h.eventService.GetEvents()
	if err != nil {
		RespondJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RespondJSON(w, events, http.StatusOK)
}
