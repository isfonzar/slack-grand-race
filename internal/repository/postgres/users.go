package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type (
	// Storage defines the struct that holds the connection to the database
	Storage struct {
		db *sql.DB
	}

	User struct {
		Id       string
		Name     string
		Balance  int
		IsActive bool
	}
)

// NewUsersStorage returns a new database storage
func NewUsersStorage(db *sql.DB) *Storage {
	return &Storage{
		db,
	}
}

// Get gets a user from the database
func (s *Storage) Get(id string) (*User, error) {
	query := `SELECT id, name, balance, is_active FROM users WHERE id = $1`
	row := s.db.QueryRow(
		query,
		id,
	)

	u := User{}
	if err := row.Scan(&u.Id, &u.Name, &u.Balance, &u.IsActive); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &u, errors.Wrap(err, "Get() could not Scan() row")
	}

	return &u, nil
}

// Create adds a user to the database
func (s *Storage) Create(id, name string) error {
	query := `INSERT INTO users(id, name, balance, is_active) VALUES($1, $2, $3, $4)`
	_, err := s.db.Exec(
		query,
		id,
		name,
		0,
		true,
	)
	if err != nil {
		return errors.Wrap(err, "Create() could not exec query")
	}

	return nil
}
