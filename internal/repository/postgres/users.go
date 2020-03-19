package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/isfonzar/slack-grand-race/pkg/domain"
	_ "github.com/lib/pq"
)

type (
	// Storage defines the struct that holds the connection to the database
	Storage struct {
		db *sql.DB
	}
)

var (
	QueryError = errors.New("could not run query")
	ScanError  = errors.New("could not Scan() row")
)

// NewUsersStorage returns a new database storage
func NewUsersStorage(db *sql.DB) *Storage {
	return &Storage{
		db,
	}
}

// Get gets a user from the database
func (s *Storage) Get(id string) (*domain.User, error) {
	query := `SELECT id, name, balance, is_active FROM users WHERE id = $1`
	row := s.db.QueryRow(
		query,
		id,
	)

	u := domain.User{}
	if err := row.Scan(&u.Id, &u.Name, &u.Balance, &u.IsActive); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &u, fmt.Errorf("%w : %v", ScanError, err)
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
		return fmt.Errorf("%w : %v", QueryError, err)
	}

	return nil
}

func (s *Storage) IncrementBalance(id string, inc int) error {
	query := `UPDATE users SET balance = balance + $1 WHERE id = $2`
	_, err := s.db.Exec(
		query,
		inc,
		id,
	)
	if err != nil {
		return fmt.Errorf("%w : %v", QueryError, err)
	}

	return nil
}

func (s *Storage) GetRanking() ([]domain.User, error) {
	query := `SELECT id, name, balance, is_active FROM users WHERE is_active = true ORDER BY balance DESC`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%w : %v", QueryError, err)
	}

	defer rows.Close()

	var users []domain.User

	for rows.Next() {
		u := domain.User{}
		if err := rows.Scan(&u.Id, &u.Name, &u.Balance, &u.IsActive); err != nil {
			return users, fmt.Errorf("%w : %v", ScanError, err)
		}
		users = append(users, u)
	}

	return users, nil
}
