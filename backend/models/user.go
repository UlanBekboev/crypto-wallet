package models

import "database/sql"

type User struct {
	ID       int    `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Name     sql.NullString `db:"name"`
	CreatedAt string `db:"created_at"`
}

