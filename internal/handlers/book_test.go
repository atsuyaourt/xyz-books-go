package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	db "github.com/atsuyaourt/xyz-books/internal/db/sqlc"
	mockdb "github.com/atsuyaourt/xyz-books/internal/mocks/db"
	"github.com/atsuyaourt/xyz-books/internal/services"
	"github.com/atsuyaourt/xyz-books/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateBookAPI(t *testing.T) {
	book := randomBook(t)
	authors := make([]string, 3)
	for i := range authors {
		hasMiddleName := i%2 > 0
		authors[i] = fmt.Sprintf("%s %s", util.RandomString(16), util.RandomString(12))
		if hasMiddleName {
			authors[i] = fmt.Sprintf("%s %s", authors[i], util.RandomString(8))
		}
	}
	publisher := util.RandomString(12)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder, store *mockdb.MockStore)
	}{
		{
			name: "Default",
			body: gin.H{
				"book": gin.H{
					"title":            book.Title,
					"isbn13":           book.Isbn13.String,
					"isbn10":           book.Isbn10.String,
					"price":            book.Price,
					"publication_year": book.PublicationYear,
				},
				"authors":   authors,
				"publisher": publisher,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateBookTx(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.Book{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"book": gin.H{
					"title":            book.Title,
					"isbn13":           book.Isbn13.String,
					"isbn10":           book.Isbn10.String,
					"price":            book.Price,
					"publication_year": book.PublicationYear,
				},
				"authors":   authors,
				"publisher": publisher,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateBookTx(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.Book{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			store := mockdb.NewMockStore(t)
			tc.buildStubs(store)

			handler := newTestHandler(t, store)

			router := gin.Default()
			router.POST("/books", handler.CreateBook)

			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/books"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			router.ServeHTTP(recorder, request)

			tc.checkResponse(recorder, store)
		})
	}
}

func TestGetBookAPI(t *testing.T) {
	book := randomBook(t)
	authors := make([]string, 3)
	for i := range authors {
		hasMiddleName := i%2 > 0
		authors[i] = fmt.Sprintf("%s %s", util.RandomString(16), util.RandomString(12))
		if hasMiddleName {
			authors[i] = fmt.Sprintf("%s %s", authors[i], util.RandomString(8))
		}
	}
	authorNames := strings.Join(authors, ",")
	publisherName := util.RandomString(12)

	testCases := []struct {
		name          string
		isbn          string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder, store *mockdb.MockStore)
	}{
		{
			name: "Default",
			isbn: book.Isbn13.String,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetBookByISBN(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.GetBookByISBNRow{
						Book:          book,
						Authors:       authorNames,
						PublisherName: publisherName,
					}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidISBN13",
			isbn: "INVALIDISBN13",
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			isbn: book.Isbn13.String,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetBookByISBN(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.GetBookByISBNRow{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotFound",
			isbn: book.Isbn13.String,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetBookByISBN(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.GetBookByISBNRow{}, db.ErrRecordNotFound)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			store := mockdb.NewMockStore(t)
			tc.buildStubs(store)

			handler := newTestHandler(t, store)

			router := gin.Default()
			router.GET("/books/:isbn", handler.GetBook)

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/books/%s", tc.isbn)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			router.ServeHTTP(recorder, request)

			tc.checkResponse(recorder, store)
		})
	}
}

func TestListBooksAPI(t *testing.T) {
	n := 10
	books := make([]db.Book, n)
	for i := range books {
		books[i] = randomBook(t)
	}

	testCases := []struct {
		name          string
		query         services.ListBooksReq
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder, store *mockdb.MockStore)
	}{
		{
			name: "Default",
			query: services.ListBooksReq{
				Page:    1,
				PerPage: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListBooks(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return([]db.ListBooksRow{}, nil)
				store.EXPECT().CountBooks(mock.AnythingOfType("*gin.Context"), mock.Anything).Return(int64(n), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InternalError",
			query: services.ListBooksReq{
				Page:    1,
				PerPage: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListBooks(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return([]db.ListBooksRow{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPage",
			query: services.ListBooksReq{
				Page:    -1,
				PerPage: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertNotCalled(t, "ListBooks", mock.AnythingOfType("*gin.Context"), mock.Anything)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidLimit",
			query: services.ListBooksReq{
				Page:    1,
				PerPage: 10000,
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertNotCalled(t, "ListBooks", mock.AnythingOfType("*gin.Context"), mock.Anything)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "EmptySlice",
			query: services.ListBooksReq{
				Page:    1,
				PerPage: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListBooks(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return([]db.ListBooksRow{}, nil)
				store.EXPECT().CountBooks(mock.AnythingOfType("*gin.Context"), mock.Anything).Return(0, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "CountInternalError",
			query: services.ListBooksReq{
				Page:    1,
				PerPage: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListBooks(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return([]db.ListBooksRow{}, nil)
				store.EXPECT().CountBooks(mock.AnythingOfType("*gin.Context"), mock.Anything).Return(0, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			store := mockdb.NewMockStore(t)
			tc.buildStubs(store)

			handler := newTestHandler(t, store)

			router := gin.Default()
			router.GET("/books", handler.ListBooks)

			recorder := httptest.NewRecorder()

			url := "/books"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("page", fmt.Sprintf("%d", tc.query.Page))
			q.Add("per_page", fmt.Sprintf("%d", tc.query.PerPage))
			request.URL.RawQuery = q.Encode()

			router.ServeHTTP(recorder, request)

			tc.checkResponse(recorder, store)
		})
	}
}

func TestUpdateBookAPI(t *testing.T) {
	book := randomBook(t)
	updatedBook := randomBook(t)
	book2 := randomBook(t)

	testCases := []struct {
		name          string
		isbn          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder, store *mockdb.MockStore)
	}{
		{
			name: "Default",
			isbn: book.Isbn13.String,
			body: gin.H{
				"title": updatedBook.Title,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateBookByISBN(mock.AnythingOfType("*gin.Context"), mock.MatchedBy(func(arg db.UpdateBookByISBNParams) bool {
					return !arg.NewIsbn13.Valid && !arg.NewIsbn10.Valid
				})).
					Return(db.Book{}, nil)
				store.EXPECT().GetBookByISBN(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.GetBookByISBNRow{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "UpdateISBN13",
			isbn: book.Isbn13.String,
			body: gin.H{
				"title":  updatedBook.Title,
				"isbn13": updatedBook.Isbn13.String,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateBookByISBN(mock.AnythingOfType("*gin.Context"), mock.MatchedBy(func(arg db.UpdateBookByISBNParams) bool {
					return arg.NewIsbn13.Valid && arg.NewIsbn13.String == updatedBook.Isbn13.String && !arg.NewIsbn10.Valid
				})).
					Return(db.Book{}, nil)
				store.EXPECT().GetBookByISBN(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.GetBookByISBNRow{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "UpdateISBN10",
			isbn: book.Isbn13.String,
			body: gin.H{
				"title":  updatedBook.Title,
				"isbn10": updatedBook.Isbn10.String,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateBookByISBN(mock.AnythingOfType("*gin.Context"), mock.MatchedBy(func(arg db.UpdateBookByISBNParams) bool {
					return !arg.NewIsbn13.Valid && arg.NewIsbn10.Valid && arg.NewIsbn10.String == updatedBook.Isbn10.String
				})).
					Return(db.Book{}, nil)
				store.EXPECT().GetBookByISBN(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.GetBookByISBNRow{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "UpdateBothISBN13AndISBN10",
			isbn: book.Isbn13.String,
			body: gin.H{
				"title":  updatedBook.Title,
				"isbn13": updatedBook.Isbn13.String,
				"isbn10": updatedBook.Isbn10.String,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateBookByISBN(mock.AnythingOfType("*gin.Context"), mock.MatchedBy(func(arg db.UpdateBookByISBNParams) bool {
					return arg.NewIsbn13.Valid && arg.NewIsbn13.String == updatedBook.Isbn13.String && arg.NewIsbn10.Valid && arg.NewIsbn10.String == updatedBook.Isbn10.String
				})).
					Return(db.Book{}, nil)
				store.EXPECT().GetBookByISBN(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.GetBookByISBNRow{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "UpdateMismatchedISBN13AndISBN10",
			isbn: book.Isbn13.String,
			body: gin.H{
				"title":  updatedBook.Title,
				"isbn13": updatedBook.Isbn13.String,
				"isbn10": book2.Isbn10.String,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateBookByISBN(mock.AnythingOfType("*gin.Context"), mock.MatchedBy(func(arg db.UpdateBookByISBNParams) bool {
					return !arg.NewIsbn13.Valid && !arg.NewIsbn10.Valid
				})).
					Return(db.Book{}, nil)
				store.EXPECT().GetBookByISBN(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.GetBookByISBNRow{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidISBN13",
			isbn: "INVALIDISBN13",
			body: gin.H{
				"title": updatedBook.Title,
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			isbn: book.Isbn13.String,
			body: gin.H{
				"title": updatedBook.Title,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateBookByISBN(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.Book{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotFound",
			isbn: book.Isbn13.String,
			body: gin.H{
				"title": updatedBook.Title,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateBookByISBN(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.Book{}, db.ErrRecordNotFound)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			store := mockdb.NewMockStore(t)
			tc.buildStubs(store)

			handler := newTestHandler(t, store)

			router := gin.Default()
			router.PUT("/books/:isbn", handler.UpdateBook)

			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/books/%s", tc.isbn)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			router.ServeHTTP(recorder, request)

			tc.checkResponse(recorder, store)
		})
	}
}

func TestDeleteBookAPI(t *testing.T) {
	book := randomBook(t)

	testCases := []struct {
		name          string
		isbn          string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder, store *mockdb.MockStore)
	}{
		{
			name: "Default",
			isbn: book.Isbn13.String,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteBookByISBN(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			name: "InvalidISBN13",
			isbn: "INVALIDISBN13",
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			isbn: book.Isbn13.String,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteBookByISBN(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			store := mockdb.NewMockStore(t)
			tc.buildStubs(store)

			handler := newTestHandler(t, store)

			router := gin.Default()
			router.DELETE("/books/:isbn", handler.DeleteBook)

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/books/%s", tc.isbn)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			router.ServeHTTP(recorder, request)

			tc.checkResponse(recorder, store)
		})
	}
}

func randomBook(t *testing.T) db.Book {
	isbn := util.NewISBN(util.RandomISBN13())
	return db.Book{
		Title: util.RandomString(24),
		Isbn13: sql.NullString{
			String: isbn.ISBN13,
			Valid:  true,
		},
		Isbn10: sql.NullString{
			String: isbn.ISBN10,
			Valid:  true,
		},
		Price:           float64(util.RandomFloat(10.0, 1500.0)),
		PublicationYear: util.RandomInt(1000, 9999),
	}
}
