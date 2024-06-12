package models

import "github.com/atsuyaourt/xyz-books/internal/util"

type Author struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
} //@name Author

type PaginatedAuthors = util.PaginatedList[Author] //@name PaginatedAuthors
