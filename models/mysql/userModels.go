package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/afifallimaruf/go-booksell/config"
	"github.com/afifallimaruf/go-booksell/models"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func Insert(name, email, password string) error {
	db, err := config.ConnectDB()
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	query := "INSERT INTO users (name, email, hashed_password, created) VALUES (?, ?, ?, UTC_TIMESTAMP)"

	_, err = db.Exec(query, name, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}

		return err
	}

	return nil
}

func Authenticate(email, password string) (int, error) {
	db, err := config.ConnectDB()
	if err != nil {
		return 0, err
	}

	var id int
	var hashedPassword []byte

	query := "SELECT id, hashed_password FROM users WHERE email = ? AND active = TRUE"

	row := db.QueryRow(query, email)
	err = row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}
