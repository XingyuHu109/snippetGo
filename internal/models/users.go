package models

import (
	"database/sql"
	"time"
)

// model of the user table
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// new model that wraps around a db connection pool
type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// if exist, return user id
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// check if a user exists with the specific id
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
