package util

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Name struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
}

func NewName(input string) *Name {
	n := new(Name)
	n.Parse(input)
	return n
}

// parse returns a person's first, last and optional middle name or initial
func (n *Name) Parse(input string) Name {
	pattern := `^(\w+\.?)\s*(.*?)\s*\b(\w+)$`

	// Compile the regular expression.
	regex := regexp.MustCompile(pattern)

	// Find the matches in the input string.
	matches := regex.FindStringSubmatch(strings.TrimSpace(input))

	if len(matches) >= 3 {
		title := cases.Title(language.English, cases.Compact)
		n.FirstName = title.String(matches[1])
		n.MiddleName = title.String(matches[2])
		n.LastName = title.String(matches[3])
	}

	return *n
}

func (n Name) String() string {
	if len(n.MiddleName) > 0 {
		return fmt.Sprintf("%s %s %s", n.FirstName, n.MiddleName, n.LastName)
	}
	return fmt.Sprintf("%s %s", n.FirstName, n.LastName)
}

func (n Name) Valid() bool {
	return (len(n.FirstName) > 0) && (len(n.LastName) > 0)
}
