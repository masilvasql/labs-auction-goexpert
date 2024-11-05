package auction_test

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuctionCloserTestSuite struct {
	suite.Suite
	Repo *MockAuctionRepository
}

func (suite *AuctionCloserTestSuite) SetupTest() {
	suite.Repo = NewMockAuctionRepository()
	os.Setenv("AUCTION_DURATION", "10s")
	os.Setenv("CHECK_INTERVAL", "10s")

	err := suite.Repo.CreateAuction(context.Background(), &auction_entity.Auction{
		Id:          "auction1",
		ProductName: "Product 1",
		Category:    "Category 1",
		Description: "Description 1",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now().Add(-11 * time.Second),
	})

	assert.NoError(suite.T(), err)
}

func (suite *AuctionCloserTestSuite) TestStartAuctionCloser() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})

	go func() {
		suite.Repo.StartAuctionCloser(ctx)
		close(done)
	}()

	time.Sleep(12 * time.Second)
	cancel()

	<-done

	expiredAuctions, err := suite.Repo.GetExpiredAuctions(ctx)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), expiredAuctions, 1)
	assert.Equal(suite.T(), "auction1", expiredAuctions[0].Id)
}

func TestAuctionCloserTestSuite(t *testing.T) {
	suite.Run(t, new(AuctionCloserTestSuite))
}
