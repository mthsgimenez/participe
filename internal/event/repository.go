package event

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/mthsgimenez/participe/internal/user"
)

var (
	ErrEventNotFound       = errors.New("event not found")
	ErrForeignKeyViolation = errors.New("foreign key constraint violated")
	ErrUniqueViolation     = errors.New("unique constraint violated")
)

type RepositoryPostgres struct {
	db *sql.DB
}

func NewRepositoryPostgres(db *sql.DB) *RepositoryPostgres {
	return &RepositoryPostgres{db}
}

func (r *RepositoryPostgres) FindById(id int) (*Event, error) {
	event := &Event{}

	row := r.db.QueryRow(`SELECT id, description, "name", "date" FROM events WHERE id = $1`, id)
	if err := row.Scan(&event.Id, &event.Description, &event.Name, &event.Date); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("event_repository: find by id: %w", ErrEventNotFound)
		}
		return nil, fmt.Errorf("event_repository: find by id: %w", err)
	}

	return event, nil
}

func (r *RepositoryPostgres) FindAll() (*[]Event, error) {
	rows, err := r.db.Query(`SELECT id, description, "name", "date" FROM events`)
	if err != nil {
		return nil, fmt.Errorf("event_repository: find all: %w", err)
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.Id, &e.Description, &e.Name, &e.Date); err != nil {
			return nil, fmt.Errorf("event_repository: find all: %w", err)
		}
		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("event_repository: find all: %w", err)
	}

	return &events, nil
}

func (r *RepositoryPostgres) Insert(e *Event) (*Event, error) {
	row := r.db.QueryRow(`INSERT INTO events (description, "name", "date") 
		VALUES ($1, $2, $3) 
		RETURNING id, description, "name", "date"`,
		e.Description, e.Name, e.Date)

	var newEvent Event
	if err := row.Scan(&newEvent.Id, &newEvent.Description, &newEvent.Name, &newEvent.Date); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return nil, fmt.Errorf("event_repository: insert: %w", ErrUniqueViolation)
			}
		}
		return nil, fmt.Errorf("event_repository: insert: %w", err)
	}

	return &newEvent, nil
}

func (r *RepositoryPostgres) Update(e *Event) (*Event, error) {
	row := r.db.QueryRow(`UPDATE events 
		SET description = $1, "name" = $2, "date" = $3 
		WHERE id = $4 
		RETURNING id, description, "name", "date"`,
		e.Description, e.Name, e.Date, e.Id)

	var updatedEvent Event
	if err := row.Scan(&updatedEvent.Id, &updatedEvent.Description, &updatedEvent.Name, &updatedEvent.Date); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return nil, fmt.Errorf("event_repository: update: %w", ErrUniqueViolation)
			}
		}
		return nil, fmt.Errorf("event_repository: update: %w", err)
	}

	return &updatedEvent, nil
}

func (r *RepositoryPostgres) DeleteById(id int) error {
	_, err := r.db.Exec(`DELETE FROM events WHERE id = $1`, id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23503" {
				return fmt.Errorf("event_repository: delete: %w", ErrForeignKeyViolation)
			}
		}
		return fmt.Errorf("event_repository: delete: %w", err)
	}
	return nil
}

func (r *RepositoryPostgres) Exists(id int) (bool, error) {
	row := r.db.QueryRow(`SELECT COUNT(*) FROM events WHERE id = $1`, id)
	var count int
	if err := row.Scan(&count); err != nil {
		return false, fmt.Errorf("event_repository: exists: %w", err)
	}
	return count > 0, nil
}

func (r *RepositoryPostgres) FindUpcoming() (*[]Event, error) {
	rows, err := r.db.Query(`SELECT id, description, "name", "date" FROM events WHERE "date" > NOW() ORDER BY "date" ASC`)
	if err != nil {
		return nil, fmt.Errorf("event_repository: find upcoming: %w", err)
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var ev Event
		if err := rows.Scan(&ev.Id, &ev.Description, &ev.Name, &ev.Date); err != nil {
			return nil, fmt.Errorf("event_repository: find upcoming: %w", err)
		}
		events = append(events, ev)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("event_repository: find upcoming: %w", err)
	}

	return &events, nil
}

func (r *RepositoryPostgres) FindCheckedUsers(e *Event) (*[]user.User, error) {
	rows, err := r.db.Query(`SELECT u.id, u.email, u.company_id, u.name FROM events_users eu JOIN users u ON eu.user_id = u.id WHERE eu.event_id = $1`)
	if err != nil {
		return nil, fmt.Errorf("event_repository: find checked users: %w", err)
	}

	var users []user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.Id, &u.Email, &u.Company.Id, &u.Name); err != nil {
			return nil, fmt.Errorf("event_repository: find checked users: %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("event_repository: find checked users: %w", err)
	}

	return &users, nil
}

func (r *RepositoryPostgres) CheckinUser(e *Event, u *user.User) error {
	_, err := r.db.Exec(`INSERT INTO events_users (user_id, event_id) VALUES ($1, $2)`, u.Id, e.Id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return fmt.Errorf("event_repository: checkin user: %w", ErrUniqueViolation)
			}
		}
		return fmt.Errorf("event_repository: checkin user: %w", err)
	}
	return nil
}
