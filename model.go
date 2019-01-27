// model.go

package main

import (
	"database/sql"
	"fmt"
)

type user struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u *user) getUser(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT name, email FROM contact WHERE id=%d", u.ID)
	return db.QueryRow(statement).Scan(&u.Name, &u.Email)
}

func (u *user) updateUser(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE contact SET name='%s', email='%s' WHERE id=%d", u.Name, u.Email, u.ID)
	_, err := db.Exec(statement)
	return err
}

func (u *user) deleteUser(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM contact WHERE id=%d", u.ID)
	_, err := db.Exec(statement)
	return err
}

func (u *user) createUser(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO contact(name, email) VALUES('%s', '%s')", u.Name, u.Email)
	_, err := db.Exec(statement)

	if err != nil {
		return err
	}

	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&u.ID)

	if err != nil {
		return err
	}

	return nil
}

func getUsers(db *sql.DB, start, count int) ([]user, error) {
	statement := fmt.Sprintf("SELECT id, name, email FROM contact LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []user{}

	for rows.Next() {
		var u user
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
