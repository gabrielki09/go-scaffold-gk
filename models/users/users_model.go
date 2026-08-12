package usersmodel

import "time"

type UsersModel struct {
	UUID      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
	