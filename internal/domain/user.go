package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	State       State
	First_Name  string
	Description string
}
