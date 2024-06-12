package models

import "github.com/atsuyaourt/xyz-books/internal/util"

type Book struct {
	Title           string   `json:"title"`
	ISBN13          string   `json:"isbn13"`
	ISBN10          string   `json:"isbn10"`
	Price           float64  `json:"price"`
	PublicationYear int64    `json:"publication_year"`
	ImageUrl        string   `json:"image_url"`
	Edition         string   `json:"edition"`
	Authors         []string `json:"authors"`
	Publisher       string   `json:"publisher"`
} //@name Book

type PaginatedBooks = util.PaginatedList[Book] //@name PaginatedBooks
