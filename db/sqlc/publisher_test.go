package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/emiliogozo/xyz-books/internal/util"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PublisherTestSuite struct {
	suite.Suite
}

func TestPublisherTestSuite(t *testing.T) {
	suite.Run(t, new(PublisherTestSuite))
}

func (ts *PublisherTestSuite) SetupTest() {
	err := util.DBMigrationUp(testConfig.MigrationSrc, testDBUrl)
	require.NoError(ts.T(), err, "db migration problem")
}

func (ts *PublisherTestSuite) TearDownTest() {
	err := util.DBMigrationDown(testConfig.MigrationSrc, testDBUrl)
	require.NoError(ts.T(), err, "reverse db migration problem")
}

func (ts *PublisherTestSuite) TestCreatePublisher() {
	createRandomPublisher(ts.T())
}
func (ts *PublisherTestSuite) TestGetPublisher() {
	t := ts.T()
	publisher := createRandomPublisher(ts.T())

	gotPublisher, err := testStore.GetPublisher(context.Background(), publisher.PublisherID)
	require.NoError(t, err)
	requirePublisherEqual(t, publisher, gotPublisher)
}

func (ts *PublisherTestSuite) TestGetPublisherByName() {
	t := ts.T()
	publisher := createRandomPublisher(ts.T())

	gotPublisher, err := testStore.GetPublisherByName(context.Background(), publisher.PublisherName)
	require.NoError(t, err)
	requirePublisherEqual(t, publisher, gotPublisher)
}

func (ts *PublisherTestSuite) TestListPublishers() {
	t := ts.T()
	n := 10
	for i := 0; i < n; i++ {
		createRandomPublisher(t)
	}

	arg := ListPublishersParams{
		Limit:  5,
		Offset: 5,
	}

	gotPublishers, err := testStore.ListPublishers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, gotPublishers, 5)

	for i := range gotPublishers {
		require.NotEmpty(t, gotPublishers[i])
	}
}

func (ts *PublisherTestSuite) TestUpdatePublisher() {
	t := ts.T()

	oldPublisher := createRandomPublisher(t)
	newName := util.RandomString(22)

	arg := UpdatePublisherParams{
		PublisherID: oldPublisher.PublisherID,
		PublisherName: sql.NullString{
			String: newName,
			Valid:  true,
		},
	}

	updatedPublisher, err := testStore.UpdatePublisher(context.Background(), arg)
	require.NoError(t, err)
	require.NotEqual(t, oldPublisher.PublisherName, updatedPublisher.PublisherName)
	require.Equal(t, newName, updatedPublisher.PublisherName)
}

func (ts *PublisherTestSuite) TestDeletePublisher() {
	t := ts.T()
	publisher := createRandomPublisher(t)

	err := testStore.DeletePublisher(context.Background(), publisher.PublisherID)
	require.NoError(t, err)

	gotPublisher, err := testStore.GetPublisher(context.Background(), publisher.PublisherID)
	require.Error(t, err)
	require.Empty(t, gotPublisher)
}

func createRandomPublisher(t *testing.T) Publisher {
	publisherName := util.RandomString(16)
	publisher, err := testStore.CreatePublisher(context.Background(), publisherName)
	require.NoError(t, err)
	require.NotEmpty(t, publisher)

	require.Equal(t, publisherName, publisher.PublisherName)

	return publisher
}

func requirePublisherEqual(t *testing.T, expected, actual Publisher) {
	require.Equal(t, expected.PublisherName, actual.PublisherName)
}
