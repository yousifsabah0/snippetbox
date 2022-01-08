package mysql

import (
	"database/sql"
	"errors"

	"github.com/yousifsabah0/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (model *SnippetModel) Create (title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippet (title, content, expires)
	VALUES(?, ?, DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := model.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	
	return int(id), nil
}

func (model *SnippetModel) FindOne (id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, expires, created_at FROM snippet WHERE id = ? AND expires > UTC_TIMESTAMP()`

	row := model.DB.QueryRow(stmt, id)
	snippet := &models.Snippet{}
	
	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Expires, &snippet.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotRecord
		} else {
			return nil, err
		}
	}

	return snippet, nil
}

func (model *SnippetModel) FindLatest () ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, expires, created_at FROM snippet WHERE expires > UTC_TIMESTAMP() ORDER BY created_at DESC LIMIT 10`
	rows, err := model.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*models.Snippet{}
	
	for rows.Next() {
		s := &models.Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Expires, &s.CreatedAt)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}