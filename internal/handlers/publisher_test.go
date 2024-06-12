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
	"github.com/atsuyaourt/xyz-books/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreatePublisherAPI(t *testing.T) {
	publisher := randomPublisher(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder, store *mockdb.MockStore)
	}{
		{
			name: "Default",
			body: gin.H{
				"publisher_name": publisher.PublisherName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreatePublisher(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.Publisher{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"publisher_name": publisher.PublisherName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreatePublisher(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.Publisher{}, sql.ErrConnDone)
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
			router.POST("/publishers", handler.CreatePublisher)

			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/publishers"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			router.ServeHTTP(recorder, request)

			tc.checkResponse(recorder, store)
		})
	}
}

func TestGetPublisherAPI(t *testing.T) {
	publisher := randomPublisher(t)

	testCases := []struct {
		name          string
		id            int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder, store *mockdb.MockStore)
	}{
		{
			name: "Default",
			id:   publisher.PublisherID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetPublisher(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.Publisher{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InternalError",
			id:   publisher.PublisherID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetPublisher(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.Publisher{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotFound",
			id:   publisher.PublisherID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetPublisher(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.Publisher{}, db.ErrRecordNotFound)
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
			router.GET("/publishers/:id", handler.GetPublisher)

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/publishers/%d", tc.id)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			router.ServeHTTP(recorder, request)

			tc.checkResponse(recorder, store)
		})
	}
}

func TestListPublishersAPI(t *testing.T) {
	n := 10
	publishers := make([]db.Publisher, n)
	for i := range publishers {
		publishers[i] = randomPublisher(t)
	}

	testCases := []struct {
		name          string
		query         listPublishersReq
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder, store *mockdb.MockStore)
	}{
		{
			name: "Default",
			query: listPublishersReq{
				Page:    1,
				PerPage: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListPublishers(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(publishers, nil)
				store.EXPECT().CountPublishers(mock.AnythingOfType("*gin.Context")).Return(int64(n), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InternalError",
			query: listPublishersReq{
				Page:    1,
				PerPage: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListPublishers(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return([]db.Publisher{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPage",
			query: listPublishersReq{
				Page:    -1,
				PerPage: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertNotCalled(t, "ListPublishers", mock.AnythingOfType("*gin.Context"), mock.Anything)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidLimit",
			query: listPublishersReq{
				Page:    1,
				PerPage: 10000,
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertNotCalled(t, "ListPublishers", mock.AnythingOfType("*gin.Context"), mock.Anything)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "EmptySlice",
			query: listPublishersReq{
				Page:    1,
				PerPage: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListPublishers(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return([]db.Publisher{}, nil)
				store.EXPECT().CountPublishers(mock.AnythingOfType("*gin.Context")).Return(int64(n), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "CountInternalError",
			query: listPublishersReq{
				Page:    1,
				PerPage: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListPublishers(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return([]db.Publisher{}, nil)
				store.EXPECT().CountPublishers(mock.AnythingOfType("*gin.Context")).Return(0, sql.ErrConnDone)
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
			router.GET("/publishers", handler.ListPublishers)

			recorder := httptest.NewRecorder()

			url := "/publishers"
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

func TestUpdatePublisherAPI(t *testing.T) {
	publisher := randomPublisher(t)
	updatedPublisher := randomPublisher(t)

	testCases := []struct {
		name          string
		id            int64
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder, store *mockdb.MockStore)
	}{
		{
			name: "Default",
			id:   publisher.PublisherID,
			body: gin.H{
				"publisher_name": updatedPublisher.PublisherName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdatePublisher(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.Publisher{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InternalError",
			id:   publisher.PublisherID,
			body: gin.H{
				"publisher_name": updatedPublisher.PublisherName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdatePublisher(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.Publisher{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotFound",
			id:   publisher.PublisherID,
			body: gin.H{
				"publisher_name": updatedPublisher.PublisherName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdatePublisher(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(db.Publisher{}, db.ErrRecordNotFound)
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
			router.PUT("/publishers/:id", handler.UpdatePublisher)

			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/publishers/%d", tc.id)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			router.ServeHTTP(recorder, request)

			tc.checkResponse(recorder, store)
		})
	}
}

func TestDeletePublisherAPI(t *testing.T) {
	publisher := randomPublisher(t)

	testCases := []struct {
		name          string
		id            int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder, store *mockdb.MockStore)
	}{
		{
			name: "Default",
			id:   publisher.PublisherID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeletePublisher(mock.AnythingOfType("*gin.Context"), mock.Anything).
					Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, store *mockdb.MockStore) {
				store.AssertExpectations(t)
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			name: "InternalError",
			id:   publisher.PublisherID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeletePublisher(mock.AnythingOfType("*gin.Context"), mock.Anything).
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
			router.DELETE("/publishers/:id", handler.DeletePublisher)

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/publishers/%d", tc.id)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			router.ServeHTTP(recorder, request)

			tc.checkResponse(recorder, store)
		})
	}
}

func randomPublisher(t *testing.T) db.Publisher {
	return db.Publisher{
		PublisherID:   util.RandomInt(1, 111),
		PublisherName: util.RandomString(24),
	}
}
