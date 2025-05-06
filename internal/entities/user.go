package entities

import (
    "gopkg.in/guregu/null.v4"
    "time"
)

type User struct {
    UserID     string      `db:"id" json:"id"`
    Email      string      `db:"email" json:"email"`
    FirstName  null.String `db:"first_name" json:"first_name"`
    FullName   null.String `db:"full_name" json:"full_name"`
    LastName   null.String `db:"last_name" json:"last_name"`
    DeletedDT  null.Time   `db:"deleted_dt" json:"deleted_dt"`
    InsertedDT time.Time   `db:"inserted_dt" json:"inserted_dt"`
    UpdatedDT  time.Time   `db:"updated_dt" json:"updated_dt"`
}
