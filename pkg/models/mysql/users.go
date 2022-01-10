package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/yousifsabah0/snippetbox/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (model *UserModel) Insert (name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password),12)
	if err != nil {
		return err
	}

	stmt := "INSERT INTO users (name, email, password) VALUES (?, ?, ?)"
	_, err = model.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if  mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email"){
				return models.ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (model *UserModel) Authenticate (email, password string) (int, error) {
	var id int
	var hashedPassword string

	stmt := "SELECT id, password FROM users WHERE email = ?"
	row := model.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (model *UserModel) Get (id int) (*models.User,error) {
	return nil, nil
}