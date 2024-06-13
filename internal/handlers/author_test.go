package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	db "github.com/atsuyaourt/xyz-books/internal/db/sqlc"
	mockdb "github.com/atsuyaourt/xyz-books/internal/mocks/db"
	"github.com/atsuyaourt/xyz-books/internal/services"
	"github.com/atsuyaourt/xyz-books/internal/util"
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

			handler := newTestHandler(t, store)

			router := gin.Default()
			router.POST("/authors", handler.CreateAuthor)

			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/authors"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			router.ServeHTTP(recorder, request)

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

			handler := newTestHandler(t, store)

			router := gin.Default()
			router.GET("/authors/:id", handler.GetAuthor)

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/authors/%d", tc.id)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			router.ServeHTTP(recorder, request)

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
		query         services.ListAuthorsReq
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder, store *mockdb.MockStore)
	}{
		{
			name: "Default",
			query: services.ListAuthorsReq{
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
			query: services.ListAuthorsReq{
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
			query: services.ListAuthorsReq{
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
			query: services.ListAuthorsReq{
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
			query: services.ListAuthorsReq{
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
			query: services.ListAuthorsReq{
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

			handler := newTestHandler(t, store)

			router := gin.Default()
			router.GET("/authors", handler.ListAuthors)

			recorder := httptest.NewRecorder()

			url := "/authors"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("page", fmt.Sprintf("%d", tc.query.Page))
			q.Add("per_page", fmt.Sprintf("%d", tc.query.PerPage))
			request.URL.RawQuery = q.Encode()

			router.ServeHTTP(recorder, request)

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

			handler := newTestHandler(t, store)

			router := gin.Default()
			router.PUT("/authors/:id", handler.UpdateAuthor)

			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/authors/%d", tc.id)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			router.ServeHTTP(recorder, request)

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

			handler := newTestHandler(t, store)

			router := gin.Default()
			router.DELETE("/authors/:id", handler.DeleteAuthor)

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/authors/%d", tc.id)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			router.ServeHTTP(recorder, request)

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
