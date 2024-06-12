package models

import "github.com/atsuyaourt/xyz-books/internal/util"

type Publisher struct {
	PublisherName string `json:"publisher_name"`
} //@name Publisher

type PaginatedPublishers = util.PaginatedList[Publisher] //@name PaginatedPublishers
