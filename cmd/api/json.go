package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Message  string            `json:"error"`
	Problems map[string]string `json:"problems,omitempty"`
}

type Validator interface {
	Validate() (problems map[string]string)
}

func RespondJSON[T any](w http.ResponseWriter, v T, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("bind: %w", err)
	}
	return nil
}

func RespondJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	e := ErrorResponse{Message: message}

	json.NewEncoder(w).Encode(e)
}

func RespondJSONErrorWithProblems(w http.ResponseWriter, message string, status int, problems map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	e := ErrorResponse{Message: message}
	if problems != nil {
		e.Problems = problems
	}

	json.NewEncoder(w).Encode(e)
}

func BindJSON[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("bindjson: %w", err)
	}
	return v, nil
}

func BindJSONValid[T Validator](r *http.Request) (v T, problems map[string]string, err error) {
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, nil, fmt.Errorf("bindjson: %w", err)
	}
	if problems := v.Validate(); len(problems) > 0 {
		return v, problems, fmt.Errorf("invalid %T: %d problems", v, len(problems))
	}
	return
}
