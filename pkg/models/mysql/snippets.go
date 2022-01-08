package mysql

import (
	"database/sql"

	"github.com/yousifsabah0/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (model *SnippetModel) Create (title, content, expires string) (int, error) {
	return 1, nil
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