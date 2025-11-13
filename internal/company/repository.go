package company

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

var (
	ErrCompanyNotFound     = errors.New("company not found")
	ErrForeignKeyViolation = errors.New("foreign key constraint violated")
	ErrUniqueViolation     = errors.New("unique constraint violated")
)

type RepositoryPostgres struct {
	db *sql.DB
}

func NewRepositoryPostgres(db *sql.DB) *RepositoryPostgres {
	return &RepositoryPostgres{db}
}

func (r *RepositoryPostgres) FindById(id int) (*Company, error) {
	cmp := &Company{}

	row := r.db.QueryRow(`SELECT * FROM companies WHERE id = $1`, id)
	if err := row.Scan(&cmp.Id, &cmp.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("company_repository: find by id: %w", ErrCompanyNotFound)
		}
		return nil, fmt.Errorf("company_repository: find by id: %w", err)
	}

	return cmp, nil
}

func (r *RepositoryPostgres) FindAll() (*[]Company, error) {
	rows, err := r.db.Query(`SELECT * FROM companies`)
	if err != nil {
		return nil, fmt.Errorf("company_repository: find all: %w", err)
	}
	defer rows.Close()

	var companies []Company
	for rows.Next() {
		var cmp Company
		if err := rows.Scan(&cmp.Id, &cmp.Name); err != nil {
			return nil, fmt.Errorf("company_repository: find all: %w", err)
		}
		companies = append(companies, cmp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("company_repository: find all: %w", err)
	}

	return &companies, nil
}

func (r *RepositoryPostgres) DeleteById(id int) error {
	_, err := r.db.Exec(`DELETE FROM companies WHERE id = $1`, id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23503" { // Foreign key violation
				return fmt.Errorf("company_repository: delete by id: %w", ErrForeignKeyViolation)
			}
		}
		return fmt.Errorf("company_repository: delete by id: %w", err)
	}
	return nil
}

func (r *RepositoryPostgres) Insert(cmp *Company) (*Company, error) {
	row := r.db.QueryRow(`INSERT INTO companies ("name") VALUES ($1) RETURNING *`, cmp.Name)

	var newCompany Company
	if err := row.Scan(&newCompany.Id, &newCompany.Name); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" { // Unique violation
				return nil, fmt.Errorf("company_repository: insert: %w", ErrUniqueViolation)
			}
		}
		return nil, fmt.Errorf("company_repository: insert: %w", err)
	}

	return &newCompany, nil
}

func (r *RepositoryPostgres) Update(cmp *Company) (*Company, error) {
	row := r.db.QueryRow(`UPDATE companies SET "name" = $1 WHERE id = $2 RETURNING *`, cmp.Name, cmp.Id)

	var updatedCompany Company
	if err := row.Scan(&updatedCompany.Id, &updatedCompany.Name); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" { // Unique violation
				return nil, fmt.Errorf("company_repository: update: %w", ErrUniqueViolation)
			}
		}
		return nil, fmt.Errorf("company_repository: update: %w", err)
	}

	return &updatedCompany, nil
}

func (r *RepositoryPostgres) Exists(id int) (bool, error) {
	row := r.db.QueryRow(`SELECT COUNT(*) FROM companies WHERE id = $1`, id)

	var count int
	if err := row.Scan(&count); err != nil {
		return false, fmt.Errorf("company_repository: exists: %w", err)
	}

	return count > 0, nil
}
