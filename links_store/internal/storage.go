package internal

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Store struct {
	db *sql.DB
}

type Link struct {
	ID    int64  `json:"id"`
	URL   string `json:"url"`
	Label string `json:"label"`
}

func New() (*Store, error) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("DB_DSN env var is required, e.g. 'user:pass@tcp(mysql:3306)/db?parseTime=true'")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	const schema = `
	CREATE TABLE IF NOT EXISTS links (
    id    BIGINT PRIMARY KEY AUTO_INCREMENT,
    url   VARCHAR(2048) NOT NULL,
    label VARCHAR(255)  NOT NULL,
    UNIQUE KEY uq_url (url(191))
);`
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) List() ([]Link, error) {
	rows, err := s.db.Query(`SELECT id, url, label FROM links ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Link
	for rows.Next() {
		var l Link
		if err := rows.Scan(&l.ID, &l.URL, &l.Label); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func (s *Store) Get(id int64) (Link, bool, error) {
	var l Link
	err := s.db.QueryRow(`SELECT id, url, label FROM links WHERE id = ?`, id).
		Scan(&l.ID, &l.URL, &l.Label)
	if err == sql.ErrNoRows {
		return Link{}, false, nil
	}
	if err != nil {
		return Link{}, false, err
	}
	return l, true, nil
}

func (s *Store) Add(l Link) (int64, error) {
	res, err := s.db.Exec(
		`INSERT INTO links (url, label) VALUES (?, ?)`, l.URL, l.Label)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) Delete(id int64) error {
	_, err := s.db.Exec(`DELETE FROM links WHERE id = ?`, id)
	return err
}
