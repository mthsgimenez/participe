package main

import (
	"errors"
	"net/http"
	"strconv"

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

func (h *eventHandler) handlePostCheckin(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		RespondJSONError(w, "id must be an int", http.StatusBadRequest)
		return
	}

	ev, err := h.eventService.GetEvent(id)
	if err != nil {
		if errors.Is(err, event.ErrEventNotFound) {
			RespondJSONError(w, "event not found", http.StatusNotFound)
			return
		}

		RespondJSONError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

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

	if err := h.eventService.CheckinUserInEvent(ev, u); err != nil {
		RespondJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *eventHandler) handleGetEvent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		RespondJSONError(w, "id must be an int", http.StatusBadRequest)
		return
	}

	ev, err := h.eventService.GetEvent(id)
	if err != nil {
		if errors.Is(err, event.ErrEventNotFound) {
			RespondJSONError(w, "event not found", http.StatusInternalServerError)
			return
		}

		RespondJSONError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	RespondJSON(w, ev, http.StatusOK)
}

func (h *eventHandler) handleGetCheckins(w http.ResponseWriter, r *http.Request) {
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

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		RespondJSONError(w, "id must be an int", http.StatusBadRequest)
		return
	}

	ev, err := h.eventService.GetEvent(id)
	if err != nil {
		if errors.Is(err, event.ErrEventNotFound) {
			RespondJSONError(w, "event not found", http.StatusInternalServerError)
			return
		}

		RespondJSONError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	list, err := h.eventService.GetCheckedUsers(ev)
	if err != nil {
		RespondJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RespondJSON(w, list, http.StatusOK)
}

func (h *eventHandler) handlePostEvent(w http.ResponseWriter, r *http.Request) {
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

	ev, problems, err := BindJSONValid[*event.Event](r)
	if err != nil {
		if len(problems) > 0 {
			RespondJSONErrorWithProblems(w, "invalid request body", http.StatusBadRequest, problems)
			return
		}

		RespondJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newEvent, err := h.eventService.CreateEvent(ev)
	if err != nil {
		RespondJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RespondJSON(w, newEvent, http.StatusOK)
}
