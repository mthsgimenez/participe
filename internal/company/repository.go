package company

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var (
	ErrCompanyNotFound     = errors.New("company not found")
	ErrInternal            = errors.New("something went wrong")
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
			return nil, ErrCompanyNotFound
		}

		return nil, ErrInternal
	}

	return cmp, nil
}

func (r *RepositoryPostgres) FindAll() (*[]Company, error) {
	rows, err := r.db.Query(`SELECT * FROM companies`)
	if err != nil {
		return nil, ErrInternal
	}
	defer rows.Close()

	var companies []Company
	for rows.Next() {
		var cmp Company
		if err := rows.Scan(&cmp.Id, &cmp.Name); err != nil {
			return nil, ErrInternal
		}
		companies = append(companies, cmp)
	}

	if err := rows.Err(); err != nil {
		return nil, ErrInternal
	}

	return &companies, nil
}

func (r *RepositoryPostgres) DeleteById(id int) error {
	_, err := r.db.Exec(`DELETE FROM companies WHERE id = $1`, id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23503" { // Foreign key violation. https://www.postgresql.org/docs/current/errcodes-appendix.html
				return ErrForeignKeyViolation
			}
		}

		return ErrInternal
	}
	return nil
}

func (r *RepositoryPostgres) Insert(cmp *Company) (*Company, error) {
	row := r.db.QueryRow(`INSERT INTO companies ("name") VALUES ($1) RETURNING *`, cmp.Name)

	var newCompany Company
	if err := row.Scan(&newCompany.Id, &newCompany.Name); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" { // Unique violation. https://www.postgresql.org/docs/current/errcodes-appendix.html
				return nil, ErrUniqueViolation
			}
		}

		return nil, ErrInternal
	}

	return &newCompany, nil
}

func (r *RepositoryPostgres) Update(cmp *Company) (*Company, error) {
	row := r.db.QueryRow(`UPDATE companies SET "name" = $1 WHERE id = $2 RETURNING *`, cmp.Name, cmp.Id)

	var updatedCompany Company
	if err := row.Scan(&updatedCompany.Id, &updatedCompany.Name); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" { // Unique violation. https://www.postgresql.org/docs/current/errcodes-appendix.html
				return nil, ErrUniqueViolation
			}
		}

		return nil, ErrInternal
	}

	return &updatedCompany, nil
}

func (r *RepositoryPostgres) Exists(id int) (bool, error) {
	row := r.db.QueryRow(`SELECT COUNT(*) FROM companies WHERE id = $1`, id)

	var count int
	if err := row.Scan(&count); err != nil {
		return false, ErrInternal
	}

	return count > 0, nil
}
