package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"testing"

	"github.com/emiliogozo/xyz-books/internal/util"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BookTestSuite struct {
	suite.Suite
}

func TestBookTestSuite(t *testing.T) {
	suite.Run(t, new(BookTestSuite))
}

func (ts *BookTestSuite) SetupTest() {
	err := util.DBMigrationUp(testConfig.MigrationSrc, testDBUrl)
	require.NoError(ts.T(), err, "db migration problem")
}

func (ts *BookTestSuite) TearDownTest() {
	err := util.DBMigrationDown(testConfig.MigrationSrc, testDBUrl)
	require.NoError(ts.T(), err, "reverse db migration problem")
}

func (ts *BookTestSuite) TestCreateBook() {
	createRandomBook(ts.T())
}

func (ts *BookTestSuite) TestGetBookByISBN() {
	t := ts.T()
	book := createRandomBook(ts.T())

	testCases := []struct {
		name        string
		arg         GetBookByISBNParams
		checkResult func(gotBook GetBookByISBNRow, err error)
	}{
		{
			name: "ISBN13",
			arg: GetBookByISBNParams{
				Isbn13: book.Isbn13,
			},
			checkResult: func(gotBook GetBookByISBNRow, err error) {
				require.NoError(t, err)
				requireBookEqual(t, book, gotBook.Book)
			},
		},
		{
			name: "ISBN10",
			arg: GetBookByISBNParams{
				Isbn10: book.Isbn10,
			},
			checkResult: func(gotBook GetBookByISBNRow, err error) {
				require.NoError(t, err)
				requireBookEqual(t, book, gotBook.Book)
			},
		},
		{
			name: "NotFound",
			arg: GetBookByISBNParams{
				Isbn13: sql.NullString{
					String: util.RandomNumericString(13),
					Valid:  true,
				},
			},
			checkResult: func(gotBook GetBookByISBNRow, err error) {
				require.Error(t, err)
				require.Empty(t, gotBook.Book)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			books, err := testStore.GetBookByISBN(context.Background(), tc.arg)
			tc.checkResult(books, err)
		})
	}
}

func (ts *BookTestSuite) TestListBooks() {
	t := ts.T()

	f, err := os.Open("../sample/book.json")
	require.NoError(t, err)
	defer f.Close()

	var data []struct {
		Book struct {
			Title           string  `json:"title"`
			Isbn13          string  `json:"isbn13"`
			Isbn10          string  `json:"isbn10"`
			Price           float64 `json:"price"`
			PublicationYear int64   `json:"publication_year"`
			ImageUrl        string  `json:"image_url"`
			Edition         string  `json:"edition"`
			PublisherID     int64   `json:"publisher_id"`
		} `json:"book"`
		Authors   []string `json:"authors"`
		Publisher string   `json:"publisher"`
	}

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&data)
	require.NoError(t, err)

	ctx := context.Background()

	books := make([]Book, len(data))
	for i, d := range data {
		authors := make([]util.Name, len(d.Authors))
		for j, author := range d.Authors {
			authors[j] = *util.NewName(author)
		}

		books[i], err = testStore.CreateBookTx(ctx, CreateBookTxParams{
			Book: CreateBookParams{
				Title: d.Book.Title,
				Isbn13: sql.NullString{
					String: d.Book.Isbn13,
					Valid:  len(d.Book.Isbn13) == 13,
				},
				Isbn10: sql.NullString{
					String: d.Book.Isbn10,
					Valid:  len(d.Book.Isbn10) == 10,
				},
				Price:           d.Book.Price,
				PublicationYear: d.Book.PublicationYear,
				ImageUrl: sql.NullString{
					String: d.Book.ImageUrl,
					Valid:  len(d.Book.ImageUrl) > 0,
				},
				Edition: sql.NullString{
					String: d.Book.Edition,
					Valid:  len(d.Book.Edition) > 0,
				},
			},
			Authors:   authors,
			Publisher: d.Publisher,
		})
		require.NoError(t, err)
		require.NotEmpty(t, books[i])
	}

	testCases := []struct {
		name        string
		arg         ListBooksParams
		checkResult func(gotBooks []ListBooksRow, err error)
	}{
		{
			name: "Default",
			arg: ListBooksParams{
				Limit: int64(len(books)),
			},
			checkResult: func(gotBooks []ListBooksRow, err error) {
				require.NoError(t, err)
				require.Len(t, gotBooks, len(books))
			},
		},
		{
			name: "Limit",
			arg: ListBooksParams{
				Limit: int64(len(books) - 2),
			},
			checkResult: func(gotBooks []ListBooksRow, err error) {
				require.NoError(t, err)
				require.Len(t, gotBooks, len(books)-2)
			},
		},
		{
			name: "Offset",
			arg: ListBooksParams{
				Limit:  int64(len(books)),
				Offset: 2,
			},
			checkResult: func(gotBooks []ListBooksRow, err error) {
				require.NoError(t, err)
				require.Len(t, gotBooks, len(books)-2)
			},
		},
		{
			name: "Title",
			arg: ListBooksParams{
				Limit: int64(len(books)),
				Title: sql.NullString{
					String: "el",
					Valid:  true,
				},
			},
			checkResult: func(gotBooks []ListBooksRow, err error) {
				require.NoError(t, err)
				require.Len(t, gotBooks, 2)
			},
		},
		{
			name: "Author",
			arg: ListBooksParams{
				Limit: int64(len(books)),
				Author: sql.NullString{
					String: "anna",
					Valid:  true,
				},
			},
			checkResult: func(gotBooks []ListBooksRow, err error) {
				require.NoError(t, err)
				require.Len(t, gotBooks, 2)
			},
		},
		{
			name: "Publisher",
			arg: ListBooksParams{
				Limit: int64(len(books)),
				Publisher: sql.NullString{
					String: "gray",
					Valid:  true,
				},
			},
			checkResult: func(gotBooks []ListBooksRow, err error) {
				require.NoError(t, err)
				require.Len(t, gotBooks, 2)
			},
		},
		{
			name: "MinPrice",
			arg: ListBooksParams{
				Limit: int64(len(books)),
				MinPrice: sql.NullFloat64{
					Float64: 1100.0,
					Valid:   true,
				},
			},
			checkResult: func(gotBooks []ListBooksRow, err error) {
				require.NoError(t, err)
				require.Len(t, gotBooks, 3)
			},
		},
		{
			name: "MaxPrice",
			arg: ListBooksParams{
				Limit: int64(len(books)),
				MaxPrice: sql.NullFloat64{
					Float64: 1850.0,
					Valid:   true,
				},
			},
			checkResult: func(gotBooks []ListBooksRow, err error) {
				require.NoError(t, err)
				require.Len(t, gotBooks, 3)
			},
		},
		{
			name: "MinMaxPrice",
			arg: ListBooksParams{
				Limit: int64(len(books)),
				MinPrice: sql.NullFloat64{
					Float64: 900.0,
					Valid:   true,
				},
				MaxPrice: sql.NullFloat64{
					Float64: 5000.0,
					Valid:   true,
				},
			},
			checkResult: func(gotBooks []ListBooksRow, err error) {
				require.NoError(t, err)
				require.Len(t, gotBooks, 4)
			},
		},
		{
			name: "MinPublicationYear",
			arg: ListBooksParams{
				Limit: int64(len(books)),
				MinPublicationYear: sql.NullInt64{
					Int64: 2000,
					Valid: true,
				},
			},
			checkResult: func(gotBooks []ListBooksRow, err error) {
				require.NoError(t, err)
				require.Len(t, gotBooks, 4)
			},
		},
		{
			name: "MaxPublicationYear",
			arg: ListBooksParams{
				Limit: int64(len(books)),
				MaxPublicationYear: sql.NullInt64{
					Int64: 2018,
					Valid: true,
				},
			},
			checkResult: func(gotBooks []ListBooksRow, err error) {
				require.NoError(t, err)
				require.Len(t, gotBooks, 3)
			},
		},
		{
			name: "MinMaxPublicationYear",
			arg: ListBooksParams{
				Limit: int64(len(books)),
				MinPublicationYear: sql.NullInt64{
					Int64: 1999,
					Valid: true,
				},
				MaxPublicationYear: sql.NullInt64{
					Int64: 2005,
					Valid: true,
				},
			},
			checkResult: func(gotBooks []ListBooksRow, err error) {
				require.NoError(t, err)
				require.Len(t, gotBooks, 2)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			books, err := testStore.ListBooks(context.Background(), tc.arg)
			tc.checkResult(books, err)
		})
	}
}

func (ts *BookTestSuite) TestUpdateBookByISBN() {
	var (
		oldBook    Book
		newTitle   string
		newPrice   float32
		newPubYear int32
	)

	t := ts.T()

	testCases := []struct {
		name        string
		buildArg    func() UpdateBookByISBNParams
		checkResult func(updatedBook Book, err error)
	}{
		{
			name: "OnlyTitle",
			buildArg: func() UpdateBookByISBNParams {
				oldBook = createRandomBook(t)
				newTitle = util.RandomString(16)
				return UpdateBookByISBNParams{
					Isbn13: oldBook.Isbn13,
					Title: sql.NullString{
						String: newTitle,
						Valid:  true,
					},
				}
			},
			checkResult: func(updatedBook Book, err error) {
				require.NoError(t, err)
				require.Equal(t, newTitle, updatedBook.Title)
				require.Equal(t, oldBook.Isbn13, updatedBook.Isbn13)
				require.Equal(t, oldBook.Isbn10, updatedBook.Isbn10)
				require.InDelta(t, oldBook.Price, updatedBook.Price, 0.001)
				require.Equal(t, oldBook.PublicationYear, updatedBook.PublicationYear)
				require.Equal(t, oldBook.ImageUrl, updatedBook.ImageUrl)
			},
		},
		{
			name: "OnlyPrice",
			buildArg: func() UpdateBookByISBNParams {
				oldBook = createRandomBook(t)
				newPrice = util.RandomFloat(1.0, 999.9)
				return UpdateBookByISBNParams{
					Isbn13: oldBook.Isbn13,
					Price: sql.NullFloat64{
						Float64: float64(newPrice),
						Valid:   true,
					},
				}
			},
			checkResult: func(updatedBook Book, err error) {
				require.NoError(t, err)
				require.InDelta(t, newPrice, updatedBook.Price, 0.001)
				require.Equal(t, oldBook.Title, updatedBook.Title)
				require.Equal(t, oldBook.Isbn13, updatedBook.Isbn13)
				require.Equal(t, oldBook.Isbn10, updatedBook.Isbn10)
				require.Equal(t, oldBook.PublicationYear, updatedBook.PublicationYear)
				require.Equal(t, oldBook.ImageUrl, updatedBook.ImageUrl)
			},
		},
		{
			name: "OnlyPulicationYear",
			buildArg: func() UpdateBookByISBNParams {
				oldBook = createRandomBook(t)
				newPubYear = int32(util.RandomInt(1111, 2000))
				return UpdateBookByISBNParams{
					Isbn13: oldBook.Isbn13,
					PublicationYear: sql.NullInt64{
						Int64: int64(newPubYear),
						Valid: true,
					},
				}
			},
			checkResult: func(updatedBook Book, err error) {
				require.NoError(t, err)
				require.Equal(t, int64(newPubYear), updatedBook.PublicationYear)
				require.Equal(t, oldBook.Title, updatedBook.Title)
				require.Equal(t, oldBook.Isbn13, updatedBook.Isbn13)
				require.Equal(t, oldBook.Isbn10, updatedBook.Isbn10)
				require.InDelta(t, oldBook.Price, updatedBook.Price, 0.001)
				require.Equal(t, oldBook.ImageUrl, updatedBook.ImageUrl)
			},
		},
		{
			name: "MutlipleFields",
			buildArg: func() UpdateBookByISBNParams {
				oldBook = createRandomBook(t)
				newTitle = util.RandomString(16)
				return UpdateBookByISBNParams{
					Isbn13: oldBook.Isbn13,
					Title: sql.NullString{
						String: newTitle,
						Valid:  true,
					},
					Price: sql.NullFloat64{
						Float64: float64(newPrice),
						Valid:   true,
					},
					PublicationYear: sql.NullInt64{
						Int64: int64(newPubYear),
						Valid: true,
					},
				}
			},
			checkResult: func(updatedBook Book, err error) {
				require.NoError(t, err)
				require.Equal(t, newTitle, updatedBook.Title)
				require.InDelta(t, newPrice, updatedBook.Price, 0.001)
				require.Equal(t, int64(newPubYear), updatedBook.PublicationYear)
				require.Equal(t, oldBook.Isbn13, updatedBook.Isbn13)
				require.Equal(t, oldBook.Isbn10, updatedBook.Isbn10)
				require.Equal(t, oldBook.ImageUrl, updatedBook.ImageUrl)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			updatedBook, err := testStore.UpdateBookByISBN(context.Background(), tc.buildArg())
			tc.checkResult(updatedBook, err)
		})
	}
}

func (ts *BookTestSuite) TestDeleteBookByISBN() {
	t := ts.T()
	books := make([]Book, 2)

	for i := range books {
		books[i] = createRandomBook(ts.T())
	}

	testCases := []struct {
		name string
		arg  DeleteBookByISBNParams
	}{
		{
			name: "ISBN13",
			arg: DeleteBookByISBNParams{
				Isbn13: books[0].Isbn13,
			},
		},
		{
			name: "ISBN10",
			arg: DeleteBookByISBNParams{
				Isbn10: books[1].Isbn10,
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			err := testStore.DeleteBookByISBN(ctx, tc.arg)
			require.NoError(t, err)
			book, err := testStore.GetBookByISBN(ctx, GetBookByISBNParams(tc.arg))
			require.Error(t, err)
			require.Empty(t, book)
		})
	}
}

func createRandomBook(t *testing.T) Book {
	isbn := util.NewISBN(util.RandomISBN13())
	publisher := createRandomPublisher(t)
	author := createRandomAuthor(t)

	arg := CreateBookParams{
		Title: util.RandomString(24),
		Isbn13: sql.NullString{
			String: isbn.ISBN13,
			Valid:  true,
		},
		Isbn10: sql.NullString{
			String: isbn.ISBN10,
			Valid:  true,
		},
		Price:           float64(util.RandomFloat(50.0, 999.9)),
		PublicationYear: util.RandomInt(1111, 2222),
		PublisherID:     publisher.PublisherID,
	}

	ctx := context.Background()

	book, err := testStore.CreateBook(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, book)

	err = testStore.CreateAuthorBookRel(ctx, CreateAuthorBookRelParams{
		AuthorID: author.AuthorID,
		BookID:   book.BookID,
	})
	require.NoError(t, err)

	requireBookEqual(t,
		Book{
			Title:           arg.Title,
			Isbn13:          arg.Isbn13,
			Isbn10:          arg.Isbn10,
			Price:           arg.Price,
			PublicationYear: arg.PublicationYear,
			PublisherID:     arg.PublisherID,
		}, book)

	return book
}

func requireBookEqual(t *testing.T, expected, actual Book) {
	require.Equal(t, expected.Title, actual.Title)
	require.Equal(t, expected.Isbn13, actual.Isbn13)
	require.Equal(t, expected.Isbn10, actual.Isbn10)
	require.InDelta(t, expected.Price, actual.Price, 0.001)
	require.Equal(t, expected.PublicationYear, actual.PublicationYear)
	require.Equal(t, expected.PublisherID, actual.PublisherID)
}
