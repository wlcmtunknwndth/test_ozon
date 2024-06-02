package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/wlcmtunknwndth/test_ozon/internal/config"
)

type Storage struct {
	driver *sql.DB
}

func New(config *config.Database) (*Storage, error) {
	const op = "storage.postgres.New"

	connStr := fmt.Sprintf("postgres://%s:%s@postgres:%s/%s?sslmode=%s", config.DbUser, config.DbPass, config.Port, config.DbName, config.SslMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{driver: db}, nil
}

func (s *Storage) Close() error {
	return s.driver.Close()
}

func (s *Storage) Ping() error {
	return s.driver.Ping()
}
