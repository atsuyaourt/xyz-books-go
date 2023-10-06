package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		checkResult func(n Name)
	}{
		{
			name:  "Default",
			input: "John Doe",
			checkResult: func(n Name) {
				require.Equal(t, "John", n.FirstName)
				require.Equal(t, "Doe", n.LastName)
				require.Equal(t, "", n.MiddleName)
			},
		},
		{
			name:  "WithMiddleName",
			input: "John Cookie Doe",
			checkResult: func(n Name) {
				require.Equal(t, "John", n.FirstName)
				require.Equal(t, "Doe", n.LastName)
				require.Equal(t, "Cookie", n.MiddleName)
			},
		},
		{
			name:  "WithMiddleInitial",
			input: "John A. Doe",
			checkResult: func(n Name) {
				require.Equal(t, "John", n.FirstName)
				require.Equal(t, "Doe", n.LastName)
				require.Equal(t, "A.", n.MiddleName)
			},
		},
		{
			name:  "MixedCase",
			input: "joHn A. DOE",
			checkResult: func(n Name) {
				require.Equal(t, "John", n.FirstName)
				require.Equal(t, "Doe", n.LastName)
				require.Equal(t, "A.", n.MiddleName)
			},
		},
		{
			name:  "ExtraSpaces",
			input: "  John Doe  ",
			checkResult: func(n Name) {
				require.Equal(t, "John", n.FirstName)
				require.Equal(t, "Doe", n.LastName)
				require.Equal(t, "", n.MiddleName)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			n := new(Name)
			n.Parse(tc.input)
			tc.checkResult(*n)
		})
	}
}

func TestString(t *testing.T) {
	testCases := []struct {
		name        string
		input       Name
		checkResult func(s string)
	}{
		{
			name: "Default",
			input: Name{
				FirstName: "John",
				LastName:  "Doe",
			},
			checkResult: func(s string) {
				require.Equal(t, "John Doe", s)

			},
		},
		{
			name: "WithMiddleName",
			input: Name{
				FirstName:  "John",
				MiddleName: "Cookie",
				LastName:   "Doe",
			},
			checkResult: func(s string) {
				require.Equal(t, "John Cookie Doe", s)
			},
		},
		{
			name: "WithMiddleInitial",
			input: Name{
				FirstName:  "John",
				MiddleName: "C.",
				LastName:   "Doe",
			},
			checkResult: func(s string) {
				require.Equal(t, "John C. Doe", s)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			tc.checkResult(tc.input.String())
		})
	}
}

func TestNewName(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		checkResult func(n Name)
	}{
		{
			name:  "Default",
			input: "John Doe",
			checkResult: func(n Name) {
				require.Equal(t, "John", n.FirstName)
				require.Equal(t, "Doe", n.LastName)
				require.Equal(t, "", n.MiddleName)
			},
		},
		{
			name:  "WithMiddleName",
			input: "John Cookie Doe",
			checkResult: func(n Name) {
				require.Equal(t, "John", n.FirstName)
				require.Equal(t, "Doe", n.LastName)
				require.Equal(t, "Cookie", n.MiddleName)
			},
		},
		{
			name:  "WithMiddleInitial",
			input: "John A. Doe",
			checkResult: func(n Name) {
				require.Equal(t, "John", n.FirstName)
				require.Equal(t, "Doe", n.LastName)
				require.Equal(t, "A.", n.MiddleName)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			name := NewName(tc.input)
			tc.checkResult(*name)
		})
	}
}
