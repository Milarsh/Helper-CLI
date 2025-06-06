package internal

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Article struct {
	ID          int64     `json:"id"`
	LinkID      int64     `json:"link_id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	PublishedAt time.Time `json:"published_at"`
}

type Store struct{ db *sql.DB }

func New() (*Store, error) {
	dsn := os.Getenv("DB_DSN")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	const schema = `
	CREATE TABLE IF NOT EXISTS articles (
    id            BIGINT PRIMARY KEY AUTO_INCREMENT,
    link_id       BIGINT NOT NULL,
    title         TEXT    NOT NULL,
    url           TEXT NOT NULL,
    published_at  DATETIME NOT NULL,
    UNIQUE KEY uq_link_url (link_id, url(191)),
    FOREIGN KEY (link_id) REFERENCES links(id) ON DELETE CASCADE
);`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) List(from, to *time.Time, asc bool) ([]Article, error) {
	query := `SELECT id, link_id, title, url, published_at FROM articles WHERE 1=1`
	args := make([]any, 0, 2)

	if from != nil {
		query += " AND published_at >= ?"
		args = append(args, from)
	}
	if to != nil {
		query += " AND published_at <= ?"
		args = append(args, to)
	}
	order := "DESC"
	if asc {
		order = "ASC"
	}
	query += " ORDER BY published_at " + order

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Article
	for rows.Next() {
		var a Article
		if err := rows.Scan(&a.ID, &a.LinkID, &a.Title, &a.URL, &a.PublishedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}
