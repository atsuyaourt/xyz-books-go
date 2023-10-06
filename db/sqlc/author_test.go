package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/emiliogozo/xyz-books/internal/util"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AuthorTestSuite struct {
	suite.Suite
}

func TestAuthorTestSuite(t *testing.T) {
	suite.Run(t, new(AuthorTestSuite))
}

func (ts *AuthorTestSuite) SetupTest() {
	err := util.DBMigrationUp(testConfig.MigrationSrc, testDBUrl)
	require.NoError(ts.T(), err, "db migration problem")
}

func (ts *AuthorTestSuite) TearDownTest() {
	err := util.DBMigrationDown(testConfig.MigrationSrc, testDBUrl)
	require.NoError(ts.T(), err, "reverse db migration problem")
}

func (ts *AuthorTestSuite) TestCreateAuthor() {
	createRandomAuthor(ts.T())
}
func (ts *AuthorTestSuite) TestGetAuthor() {
	t := ts.T()
	author := createRandomAuthor(ts.T())

	gotAuthor, err := testStore.GetAuthor(context.Background(), author.AuthorID)
	require.NoError(t, err)
	requireAuthorEqual(t, author, gotAuthor)
}

func (ts *AuthorTestSuite) TestGetAuthorByName() {
	t := ts.T()
	ctx := context.Background()
	author1 := createRandomAuthor(ts.T())
	if len(author1.MiddleName) == 0 {
		author1, _ = testStore.UpdateAuthor(
			ctx,
			UpdateAuthorParams{
				AuthorID: author1.AuthorID,
				MiddleName: sql.NullString{
					String: util.RandomString(6),
					Valid:  true,
				},
			})
	}
	author2 := createRandomAuthor(ts.T())
	if len(author2.MiddleName) > 0 {
		author2, _ = testStore.UpdateAuthor(
			ctx,
			UpdateAuthorParams{
				AuthorID: author2.AuthorID,
				MiddleName: sql.NullString{
					String: "",
					Valid:  false,
				},
			})
	}

	gotAuthor1, err := testStore.GetAuthorByName(
		ctx, GetAuthorByNameParams{
			FirstName:  author1.FirstName,
			LastName:   author1.LastName,
			MiddleName: author1.MiddleName,
		})
	require.NoError(t, err)
	requireAuthorEqual(t, author1, gotAuthor1)
	gotAuthor2, err := testStore.GetAuthorByName(
		ctx, GetAuthorByNameParams{
			FirstName:  author2.FirstName,
			LastName:   author2.LastName,
			MiddleName: author2.MiddleName,
		})
	require.NoError(t, err)
	requireAuthorEqual(t, author2, gotAuthor2)
}

func (ts *AuthorTestSuite) TestListAuthors() {
	t := ts.T()
	n := 10
	for i := 0; i < n; i++ {
		createRandomAuthor(t)
	}

	arg := ListAuthorsParams{
		Limit:  5,
		Offset: 5,
	}

	gotAuthors, err := testStore.ListAuthors(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, gotAuthors, 5)

	for i := range gotAuthors {
		require.NotEmpty(t, gotAuthors[i])
	}
}

func (ts *AuthorTestSuite) TestUpdateAuthor() {
	var (
		oldAuthor Author
		newString string
	)
	t := ts.T()

	testCases := []struct {
		name        string
		buildArg    func() UpdateAuthorParams
		checkResult func(updatedAuthor Author, err error)
	}{
		{
			name: "OnlyFirstName",
			buildArg: func() UpdateAuthorParams {
				oldAuthor = createRandomAuthor(t)
				newString = util.RandomString(12)
				return UpdateAuthorParams{
					AuthorID: oldAuthor.AuthorID,
					FirstName: sql.NullString{
						String: newString,
						Valid:  true,
					},
				}
			},
			checkResult: func(updatedAuthor Author, err error) {
				require.NoError(t, err)
				require.NotEqual(t, oldAuthor.FirstName, updatedAuthor.FirstName)
				require.Equal(t, newString, updatedAuthor.FirstName)
				require.Equal(t, oldAuthor.LastName, updatedAuthor.LastName)
				require.Equal(t, oldAuthor.MiddleName, updatedAuthor.MiddleName)
			},
		},
		{
			name: "OnlyLastName",
			buildArg: func() UpdateAuthorParams {
				oldAuthor = createRandomAuthor(t)
				newString = util.RandomString(12)
				return UpdateAuthorParams{
					AuthorID: oldAuthor.AuthorID,
					LastName: sql.NullString{
						String: newString,
						Valid:  true,
					},
				}
			},
			checkResult: func(updatedAuthor Author, err error) {
				require.NoError(t, err)
				require.NotEqual(t, oldAuthor.LastName, updatedAuthor.LastName)
				require.Equal(t, newString, updatedAuthor.LastName)
				require.Equal(t, oldAuthor.FirstName, updatedAuthor.FirstName)
				require.Equal(t, oldAuthor.MiddleName, updatedAuthor.MiddleName)
			},
		},
		{
			name: "OnlyMiddleName",
			buildArg: func() UpdateAuthorParams {
				oldAuthor = createRandomAuthor(t)
				newString = util.RandomString(12)
				return UpdateAuthorParams{
					AuthorID: oldAuthor.AuthorID,
					MiddleName: sql.NullString{
						String: newString,
						Valid:  true,
					},
				}
			},
			checkResult: func(updatedAuthor Author, err error) {
				require.NoError(t, err)
				require.NotEqual(t, oldAuthor.MiddleName, updatedAuthor.MiddleName)
				require.Equal(t, newString, updatedAuthor.MiddleName)
				require.Equal(t, oldAuthor.FirstName, updatedAuthor.FirstName)
				require.Equal(t, oldAuthor.LastName, updatedAuthor.LastName)
			},
		},
		{
			name: "AllField",
			buildArg: func() UpdateAuthorParams {
				oldAuthor = createRandomAuthor(t)
				newString = util.RandomString(36)
				return UpdateAuthorParams{
					AuthorID: oldAuthor.AuthorID,
					FirstName: sql.NullString{
						String: newString[:12],
						Valid:  true,
					},
					LastName: sql.NullString{
						String: newString[12:24],
						Valid:  true,
					},
					MiddleName: sql.NullString{
						String: newString[24:36],
						Valid:  true,
					},
				}
			},
			checkResult: func(updatedAuthor Author, err error) {
				require.NoError(t, err)
				require.NotEqual(t, oldAuthor.FirstName, updatedAuthor.FirstName)
				require.NotEqual(t, oldAuthor.LastName, updatedAuthor.LastName)
				require.NotEqual(t, oldAuthor.MiddleName, updatedAuthor.MiddleName)
				require.Equal(t, newString[:12], updatedAuthor.FirstName)
				require.Equal(t, newString[12:24], updatedAuthor.LastName)
				require.Equal(t, newString[24:36], updatedAuthor.MiddleName)

			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			updatedAuthor, err := testStore.UpdateAuthor(context.Background(), tc.buildArg())
			tc.checkResult(updatedAuthor, err)
		})
	}
}

func (ts *AuthorTestSuite) TestDeleteAuthor() {
	t := ts.T()
	author := createRandomAuthor(t)

	err := testStore.DeleteAuthor(context.Background(), author.AuthorID)
	require.NoError(t, err)

	gotAuthor, err := testStore.GetAuthor(context.Background(), author.AuthorID)
	require.Error(t, err)
	require.Empty(t, gotAuthor)
}

func createRandomAuthor(t *testing.T) Author {
	hasMiddleName := util.RandomInt(0, 1) == 0
	arg := CreateAuthorParams{
		FirstName: util.RandomString(12),
		LastName:  util.RandomString(12),
	}
	if hasMiddleName {
		arg.MiddleName = util.RandomString(1)
	}

	author, err := testStore.CreateAuthor(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, author)

	requireAuthorEqual(t,
		Author{
			FirstName:  arg.FirstName,
			LastName:   arg.LastName,
			MiddleName: arg.MiddleName,
		}, author)

	return author
}

func requireAuthorEqual(t *testing.T, expected, actual Author) {
	require.Equal(t, expected.FirstName, actual.FirstName)
	require.Equal(t, expected.LastName, actual.LastName)
	require.Equal(t, expected.MiddleName, actual.MiddleName)
}
