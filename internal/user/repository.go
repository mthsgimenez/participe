package user

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrForeignKeyViolation = errors.New("foreign key constraint violated")
	ErrUniqueViolation     = errors.New("unique constraint violated")
)

type RepositoryPostgres struct {
	db *sql.DB
}

func NewRepositoryPostgres(db *sql.DB) *RepositoryPostgres {
	return &RepositoryPostgres{db}
}

func (r *RepositoryPostgres) FindById(id int) (*User, error) {
	user := &User{}
	row := r.db.QueryRow(`SELECT id, email, company_id, "name", "role", "password" FROM users WHERE id = $1`, id)
	if err := row.Scan(&user.Id, &user.Email, &user.Company.Id, &user.Name, &user.Role, &user.hash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user_repository: find by id: %w", ErrUserNotFound)
		}
		return nil, fmt.Errorf("user_repository: find by id: %w", err)
	}

	return user, nil
}

func (r *RepositoryPostgres) FindByEmail(email string) (*User, error) {
	user := &User{}
	row := r.db.QueryRow(`SELECT id, email, company_id, "name", "role", "password" FROM users WHERE email = $1`, email)
	if err := row.Scan(&user.Id, &user.Email, &user.Company.Id, &user.Name, &user.Role, &user.hash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user_repository: find by email: %w", ErrUserNotFound)
		}
		return nil, fmt.Errorf("user_repository: find by email: %w", err)
	}

	return user, nil
}

func (r *RepositoryPostgres) FindAll() (*[]User, error) {
	rows, err := r.db.Query(`SELECT id, email, company_id, "name", "role", "password" FROM users`)
	if err != nil {
		return nil, fmt.Errorf("user_repository: find all users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.Email, &u.Company.Id, &u.Name, &u.Role, &u.hash); err != nil {
			return nil, fmt.Errorf("user_repository: scan user in find all: %w", err)
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("user_repository: rows iteration in find all: %w", err)
	}

	return &users, nil
}

func (r *RepositoryPostgres) Insert(u *User) (*User, error) {
	row := r.db.QueryRow(`INSERT INTO users (email, company_id, "name", "role", "password") 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, email, company_id, "name", "role", "password"`,
		u.Email, u.Company.Id, u.Name, u.Role.String(), u.hash)

	var newUser User
	var role string
	if err := row.Scan(&newUser.Id, &newUser.Email, &newUser.Company.Id, &newUser.Name, &role, &newUser.hash); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" { // Unique violation
				return nil, fmt.Errorf("user_repository: insert user: %w", ErrUniqueViolation)
			}
		}
		return nil, fmt.Errorf("user_repository: insert user: %w", err)
	}

	newUser.Role = StringToUserRole(role)

	return &newUser, nil
}

func (r *RepositoryPostgres) Update(u *User) (*User, error) {
	row := r.db.QueryRow(`UPDATE users 
		SET email = $1, company_id = $2, "name" = $3, "role" = $4, "password" = $5 
		WHERE id = $6 
		RETURNING id, email, company_id, "name", "role", "password"`,
		u.Email, u.Company.Id, u.Name, u.Role.String(), u.hash, u.Id)

	var updatedUser User
	if err := row.Scan(&updatedUser.Id, &updatedUser.Email, &updatedUser.Company.Id, &updatedUser.Name, &updatedUser.Role, &updatedUser.hash); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" { // Unique violation
				return nil, fmt.Errorf("user_repository: update user: %w", ErrUniqueViolation)
			}
		}
		return nil, fmt.Errorf("user_repository: update user: %w", err)
	}

	return &updatedUser, nil
}

func (r *RepositoryPostgres) DeleteById(id int) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23503" { // Foreign key violation
				return fmt.Errorf("user_repository: delete user by id: %w", ErrForeignKeyViolation)
			}
		}
		return fmt.Errorf("user_repository: delete user by id: %w", err)
	}
	return nil
}

func (r *RepositoryPostgres) Exists(id int) (bool, error) {
	row := r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE id = $1`, id)
	var count int
	if err := row.Scan(&count); err != nil {
		return false, fmt.Errorf("user_repository: check if user exists: %w", err)
	}
	return count > 0, nil
}
