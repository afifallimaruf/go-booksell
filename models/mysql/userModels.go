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

func Authenticate(email, password string) (*models.User, error) {
	db, err := config.ConnectDB()
	if err != nil {
		return nil, err
	}

	// var id int
	// var hashedPassword []byte

	query := "SELECT id, hashed_password, role FROM users WHERE email = ? AND active = TRUE"

	row := db.QueryRow(query, email)
	user := &models.User{}
	err = row.Scan(&user.Id, &user.HashedPassword, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrInvalidCredentials
		} else {
			return nil, err
		}
	}

	if err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, models.ErrInvalidCredentials
		} else {
			return nil, err
		}
	}

	return user, nil
}
