package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type ISBNType string

const (
	ISBN13 ISBNType = "ISBN13"
	ISBN10 ISBNType = "ISBN10"
)

type ISBN struct {
	ISBN13     string   `json:"isbn13"`
	ISBN10     string   `json:"isbn10"`
	SourceType ISBNType `json:"-"`
}

func NewISBN(input string) *ISBN {
	input = strings.ReplaceAll(input, " ", "")
	input = strings.ReplaceAll(input, "-", "")
	input = strings.ToUpper(input)

	isbn := new(ISBN)
	if isbn.IsValidISBN13(input) {
		isbn.ISBN13 = input
		isbn.ISBN10, _ = isbn13To10(isbn.ISBN13)
		isbn.SourceType = ISBN13
	} else if isbn.IsValidISBN10(input) {
		isbn.ISBN10 = input
		isbn.ISBN13, _ = isbn10To13(isbn.ISBN10)
		isbn.SourceType = ISBN10
	}
	return isbn
}

func (i ISBN) IsValidISBN13(input string) bool {
	if len(input) != 13 {
		return false
	}
	if match, _ := regexp.MatchString("^[0-9]+$", input); !match {
		return false
	}
	checkDigit := calculateISBN13CheckDigit(input[:12])
	return string(input[12]) == checkDigit
}

func (i ISBN) IsValidISBN10(input string) bool {
	if len(input) != 10 {
		return false
	}
	if match, _ := regexp.MatchString("^[0-9X]+$", input); !match {
		return false
	}
	checkDigit := calculateISBN10CheckDigit(input[:9])
	return string(input[9]) == checkDigit
}

func (i ISBN) String() string {
	if len(i.ISBN13) == 13 {
		return i.ISBN13
	}

	return ""
}

// RandomISBN13 generates a random ISBN-13.
func RandomISBN13() string {
	// The first three digits are typically the group or country identifier.
	// For simplicity, we'll use "978" as a common prefix.
	prefix := "978"

	// Generate nine random digits for the book identifier.
	bookIdentifier := RandomNumericString(9)

	// Calculate the check digit using the ISBN-13 algorithm.
	checkDigit := calculateISBN13CheckDigit(prefix + bookIdentifier)

	// Concatenate the parts to form the ISBN-13.
	return prefix + bookIdentifier + checkDigit
}

// RandomISBN10 generates a random ISBN-10.
func RandomISBN10() string {
	// Generate nine random digits for the book identifier.
	bookIdentifier := RandomNumericString(9)

	// Calculate the check digit using the ISBN-10 algorithm.
	checkDigit := calculateISBN10CheckDigit(bookIdentifier)

	// Concatenate the parts to form the ISBN-10.
	return bookIdentifier + checkDigit
}

// isbn13ToI10 converts an ISBN-13 to an ISBN-10.
func isbn13To10(isbn13 string) (string, error) {
	isbn13 = strings.ReplaceAll(isbn13, " ", "")
	isbn13 = strings.ReplaceAll(isbn13, "-", "")

	if len(isbn13) != 13 {
		return "", fmt.Errorf("invalid length")
	}

	partial := isbn13[3:12]
	checkDigit := calculateISBN10CheckDigit(partial)
	isbn10 := partial + checkDigit

	return isbn10, nil
}

// isbn10To13 converts an ISBN-10 to an ISBN-13.
func isbn10To13(isbn10 string) (string, error) {
	isbn10 = strings.ReplaceAll(isbn10, " ", "")
	isbn10 = strings.ReplaceAll(isbn10, "-", "")

	if len(isbn10) != 10 {
		return "", fmt.Errorf("invalid length")
	}

	partial := "978" + isbn10[:9]
	checkDigit := calculateISBN13CheckDigit(partial)
	isbn13 := partial + checkDigit

	return isbn13, nil
}

// calculateISBN13CheckDigit calculates the ISBN-13 check digit.
func calculateISBN13CheckDigit(digits string) string {
	sum := 0
	for i, s := range digits {
		val, _ := strconv.Atoi(string(s))
		if i%2 == 0 {
			sum += val
		} else {
			sum += val * 3
		}
	}
	checkDigit := 10 - (sum % 10)
	if checkDigit == 10 {
		checkDigit = 0
	}
	return strconv.Itoa(checkDigit)
}

// calculateISBN10CheckDigit calculates the ISBN-10 check digit.
func calculateISBN10CheckDigit(digits string) string {
	sum := 0
	weight := 10
	for _, s := range digits {
		digitVal, _ := strconv.Atoi(string(s))
		sum += digitVal * weight
		weight--
	}

	checkDigit := (11 - (sum % 11)) % 11
	if checkDigit == 10 {
		return "X"
	}

	return strconv.Itoa(checkDigit)
}
