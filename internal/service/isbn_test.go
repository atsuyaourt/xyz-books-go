package service

import (
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
	"testing"

	mockhttp "github.com/atsuyaourt/xyz-books/internal/mocks/service"
	mockutil "github.com/atsuyaourt/xyz-books/internal/mocks/util"
	"github.com/atsuyaourt/xyz-books/internal/models"
	"github.com/atsuyaourt/xyz-books/internal/util"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newMockISBNService(t *testing.T, mClient *mockhttp.MockHTTPClient, mWriter *mockutil.MockWriter) *ISBNService {
	return &ISBNService{
		apiBasePath: "http://test.com/api",
		client:      mClient,
		csvWriter:   mWriter,
	}
}

func TestFetchBooks(t *testing.T) {
	mClient := mockhttp.NewMockHTTPClient(t)
	s := newMockISBNService(t, mClient, nil)
	books := loadBooksFromFile("test_books_missing_isbn.json")

	nextPage := 1
	for nextPage != 0 {
		data, res, err := mockGetFunc(nextPage, 3, books)
		mClient.EXPECT().Get(mock.Anything).Return(res, err).Once()
		nextPage = int(data.NextPage)
	}
	ch := make(chan models.Book)

	go s.fetchBooks(ch)

	var receivedBooks []models.Book
	for book := range ch {
		receivedBooks = append(receivedBooks, book)
	}

	mClient.AssertExpectations(t)
	require.Len(t, receivedBooks, 10)
}

func TestConvertISBN(t *testing.T) {
	inChan := make(chan models.Book)
	outChan := make(chan util.ISBN)

	s := newMockISBNService(t, nil, nil)

	go s.convertISBN(inChan, outChan)

	booksWithMissingISBN := loadBooksFromFile("test_books_missing_isbn.json")
	go func(books []models.Book) {
		for _, b := range books {
			inChan <- b
		}
		close(inChan)
	}(booksWithMissingISBN)

	books := loadBooksFromFile("test_books.json")
	slices.SortFunc(books, func(a, b models.Book) int { return cmp.Compare(a.ISBN13, b.ISBN13) })
	var actualISBNs []util.ISBN
	for isbn := range outChan {
		actualISBNs = append(actualISBNs, isbn)
		idx, found := slices.BinarySearchFunc(books, isbn, func(b models.Book, i util.ISBN) int { return cmp.Compare(b.ISBN13, i.ISBN13) })
		require.True(t, found)
		require.Positive(t, idx)
	}

	countMissing := 0
	for _, b := range booksWithMissingISBN {
		if len(b.ISBN10) != 10 || len(b.ISBN13) != 13 {
			countMissing++
		}
	}

	require.Len(t, actualISBNs, countMissing)
}

func TestUpdateISBN(t *testing.T) {
	tests := []struct {
		name       string
		inputISBNs []util.ISBN
		buildStubs func(mClient *mockhttp.MockHTTPClient)
		wantErrors int
	}{
		{
			name: "success",
			inputISBNs: []util.ISBN{
				{ISBN13: "9781234567890"},
				{ISBN13: "9780987654321"},
			},
			buildStubs: func(mClient *mockhttp.MockHTTPClient) {
				mClient.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				}, nil)
			},
			wantErrors: 0,
		},
		{
			name: "ServerError",
			inputISBNs: []util.ISBN{
				{ISBN13: "9781234567890"},
			},
			buildStubs: func(mClient *mockhttp.MockHTTPClient) {
				mClient.EXPECT().Do(mock.Anything).Return(&http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				}, nil)
			},
			wantErrors: 1,
		},
		{
			name: "RequestError",
			inputISBNs: []util.ISBN{
				{ISBN13: "9781234567890"},
			},
			buildStubs: func(mClient *mockhttp.MockHTTPClient) {
				mClient.EXPECT().Do(mock.Anything).Return(nil, fmt.Errorf("request error"))
			},
			wantErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mClient := mockhttp.NewMockHTTPClient(t)
			s := newMockISBNService(t, mClient, nil)

			tt.buildStubs(mClient)

			inChan := make(chan util.ISBN, len(tt.inputISBNs))
			outChan := make(chan error, len(tt.inputISBNs))

			for _, isbn := range tt.inputISBNs {
				inChan <- isbn
			}
			close(inChan)

			go s.updateISBN(inChan, outChan)

			var gotErrors int
			for err := range outChan {
				if err != nil {
					gotErrors++
				}
			}

			mClient.AssertExpectations(t)
			require.Equal(t, gotErrors, tt.wantErrors)
		})
	}
}

func TestAppendToCSV(t *testing.T) {
	inChan := make(chan util.ISBN)
	outChan := make(chan bool)

	mWriter := mockutil.NewMockWriter(t)
	s := newMockISBNService(t, nil, mWriter)

	isbns := []util.ISBN{
		{SourceType: util.ISBN13, ISBN13: "9781234567890"},
		{SourceType: util.ISBN10, ISBN10: "1234567890"},
	}
	for range isbns {
		mWriter.EXPECT().Write(mock.Anything).Return(nil)
	}
	mWriter.EXPECT().Flush().Return()
	mWriter.EXPECT().Error().Return(nil)

	go s.appendToCSV(inChan, outChan)

	go func(isbns []util.ISBN) {
		defer close(inChan)
		for _, isbn := range isbns {
			inChan <- isbn
		}
	}(isbns)

	successCount := 0
	for isSuccess := range outChan {
		if isSuccess {
			successCount++
		}
	}

	s.csvWriter.Flush()
	if err := s.csvWriter.Error(); err != nil {
		fmt.Println("Error closing CSV writer:", err)
	}

	mWriter.AssertExpectations(t)
	require.Len(t, isbns, successCount)
}

func mockGetFunc[T any](page, perPage int, items []T) (*util.PaginatedList[T], *http.Response, error) {
	totalItems := len(items)

	start := (page - 1) * perPage
	end := start + perPage
	if end > totalItems {
		end = totalItems
	}

	_data := util.NewPaginatedList(int32(page), int32(perPage), int32(totalItems), items[start:end])

	data, err := json.Marshal(_data)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return nil, &http.Response{
			StatusCode: 500,
			Body:       io.NopCloser(strings.NewReader("")),
			Header:     make(http.Header),
		}, err
	}

	reader := bytes.NewReader(data)

	return &_data, &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(reader),
		Header:     make(http.Header),
	}, nil
}

func loadBooksFromFile(fileName string) []models.Book {
	// Read the JSON file
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return nil
	}
	defer file.Close()

	// Read the file contents
	byteValue, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return nil
	}

	// Unmarshal the JSON into a slice of Book structs
	var books []models.Book
	err = json.Unmarshal(byteValue, &books)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return nil
	}

	return books
}
