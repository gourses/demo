package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // postgres drive
	"gopkg.in/gorp.v1"
)

type Storage struct {
	db *gorp.DbMap
}

type Note struct {
	ID    int64  `db:"id, primarykey, autoincrement" json:"id"`
	Title string `db:"title,size:64" json:"title"`
	Body  string `db:"body,size:2048" json:"body"`
}

func New(connStr string) (*Storage, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open sql connection: %w", err)
	}

	dbmap := &gorp.DbMap{
		Db:      db,
		Dialect: gorp.PostgresDialect{},
	}

	dbmap.AddTableWithName(Note{}, "notes")
	if err := dbmap.CreateTablesIfNotExists(); err != nil {
		return nil, fmt.Errorf("failed to create missing tables: %w", err)
	}

	return &Storage{db: dbmap}, nil
}

func (s *Storage) GetNotes() ([]Note, error) {
	notes := make([]Note, 0)
	_, err := s.db.Select(&notes, "SELECT * FROM notes")
	if err != nil {
		return nil, fmt.Errorf("failed to get all notes: %w", err)
	}
	return notes, nil
}
