package mysql

import (
	"database/sql"

	"github.com/yousifsabah0/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (model *SnippetModel) Create (title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, expires) VALUES
					(?, ?, DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY`

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

func (model *SnippetModel) Find () ([]*models.Snippet, error) {
	return nil, nil
}

func (model *SnippetModel) FindOne (id int) (*models.Snippet, error) {
	return nil, nil
}

func (model *SnippetModel) FindLatest (id int) ([]*models.Snippet, error) {
	return nil, nil
}