package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/emiliogozo/xyz-books/db/mocks"
	db "github.com/emiliogozo/xyz-books/db/sqlc"
	"github.com/emiliogozo/xyz-books/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateAuthorAPI(t *testing.T) {
	author := randomAuthor(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder, store *mockdb.MockStore)
	}{
		{
			name: "Default",
			body: gin.H{
				"first_name":  author.FirstName,
				"last_name":   author.LastName,
				"middle_name": author.MiddleName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAuthor(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.Author{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"first_name":  author.FirstName,
				"last_name":   author.LastName,
				"middle_name": author.MiddleName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAuthor(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.Author{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NoMiddle",
			body: gin.H{
				"first_name": author.FirstName,
				"last_name":  author.LastName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAuthor(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.Author{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "NoFirst",
			body: gin.H{
				"last_name":   author.LastName,
				"middle_name": author.MiddleName,
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NoLast",
			body: gin.H{
				"first_name":  author.FirstName,
				"middle_name": author.MiddleName,
			},
			buildStubs: func(store *mockdb.MockStore) {
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("%s/authors", server.config.APIBasePath)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(recorder, store)
		})
	}
}

func TestGetAuthorAPI(t *testing.T) {
	author := randomAuthor(t)

	testCases := []struct {
		name          string
		id            int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder, store *mockdb.MockStore)
	}{
		{
			name: "Default",
			id:   author.AuthorID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAuthor(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.Author{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InternalError",
			id:   author.AuthorID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAuthor(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.Author{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotFound",
			id:   author.AuthorID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAuthor(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.Author{}, db.ErrRecordNotFound)
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("%s/authors/%d", server.config.APIBasePath, tc.id)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(recorder, store)
		})
	}
}

func TestListAuthorsAPI(t *testing.T) {
	n := 10
	authors := make([]db.Author, n)
	for i := 0; i < n; i++ {
		authors[i] = randomAuthor(t)
	}

	testCases := []struct {
		name          string
		query         listAuthorsReq
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder, store *mockdb.MockStore)
	}{
		{
			name: "Default",
			query: listAuthorsReq{
				Page:    1,
				PerPage: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListAuthors(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(authors, nil)
				store.EXPECT().CountAuthors(mock.AnythingOfType("*gin.Context")).Return(int64(n), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InternalError",
			query: listAuthorsReq{
				Page:    1,
				PerPage: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListAuthors(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return([]db.Author{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPage",
			query: listAuthorsReq{
				Page:    -1,
				PerPage: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertNotCalled(t, "ListAuthors", mock.AnythingOfType("*gin.Context"), mock.Anything)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidLimit",
			query: listAuthorsReq{
				Page:    1,
				PerPage: -4,
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertNotCalled(t, "ListAuthors", mock.AnythingOfType("*gin.Context"), mock.Anything)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "EmptySlice",
			query: listAuthorsReq{
				Page:    1,
				PerPage: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListAuthors(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return([]db.Author{}, nil)
				store.EXPECT().CountAuthors(mock.AnythingOfType("*gin.Context")).Return(int64(n), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "CountInternalError",
			query: listAuthorsReq{
				Page:    1,
				PerPage: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListAuthors(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return([]db.Author{}, nil)
				store.EXPECT().CountAuthors(mock.AnythingOfType("*gin.Context")).Return(0, sql.ErrConnDone)
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("%s/authors", server.config.APIBasePath)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("page", fmt.Sprintf("%d", tc.query.Page))
			q.Add("per_page", fmt.Sprintf("%d", tc.query.PerPage))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(recorder, store)
		})
	}
}

func TestUpdateAuthorAPI(t *testing.T) {
	author := randomAuthor(t)
	updatedAuthor := randomAuthor(t)

	testCases := []struct {
		name          string
		id            int64
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder, store *mockdb.MockStore)
	}{
		{
			name: "Default",
			id:   author.AuthorID,
			body: gin.H{
				"first_name":  updatedAuthor.FirstName,
				"last_name":   updatedAuthor.LastName,
				"middle_name": updatedAuthor.MiddleName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateAuthor(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.Author{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InternalError",
			id:   updatedAuthor.AuthorID,
			body: gin.H{
				"first_name":  updatedAuthor.FirstName,
				"last_name":   updatedAuthor.LastName,
				"middle_name": updatedAuthor.MiddleName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateAuthor(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.Author{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotFound",
			id:   updatedAuthor.AuthorID,
			body: gin.H{
				"first_name":  updatedAuthor.FirstName,
				"last_name":   updatedAuthor.LastName,
				"middle_name": updatedAuthor.MiddleName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateAuthor(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.Author{}, db.ErrRecordNotFound)
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("%s/authors/%d", server.config.APIBasePath, tc.id)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(recorder, store)
		})
	}
}

func TestDeleteAuthorAPI(t *testing.T) {
	author := randomAuthor(t)

	testCases := []struct {
		name          string
		id            int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder, store *mockdb.MockStore)
	}{
		{
			name: "Default",
			id:   author.AuthorID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteAuthor(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			name: "InternalError",
			id:   author.AuthorID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteAuthor(mock.AnythingOfType("*gin.Context"), mock.Anything).
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("%s/authors/%d", server.config.APIBasePath, tc.id)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(recorder, store)
		})
	}
}

func randomAuthor(t *testing.T) db.Author {
	return db.Author{
		AuthorID:   util.RandomInt(1, 111),
		FirstName:  util.RandomString(12),
		LastName:   util.RandomString(16),
		MiddleName: util.RandomString(6),
	}
}
