package auction_test

import (
	"context"
	"errors"
	"fmt"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/infra/database/auction"
	"sync"
	"time"
)

type MockAuctionRepository struct {
	mu       sync.Mutex
	auctions map[string]*auction_entity.Auction
}

func NewMockAuctionRepository() *MockAuctionRepository {
	return &MockAuctionRepository{
		auctions: make(map[string]*auction_entity.Auction),
	}
}

func (m *MockAuctionRepository) CreateAuction(ctx context.Context, auctionEntity *auction_entity.Auction) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if auctionEntity.Id == "" {
		return errors.New("auction ID cannot be empty")
	}
	m.auctions[auctionEntity.Id] = auctionEntity
	return nil
}

var expiredAuctions []auction.AuctionEntityMongo

func (m *MockAuctionRepository) GetExpiredAuctions(ctx context.Context) ([]auction.AuctionEntityMongo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, au := range m.auctions {
		if au.Status == auction_entity.Active && time.Now().Unix() > au.Timestamp.Add(auction.GetAuctionDuration()).Unix() {
			for _, expiredAuction := range expiredAuctions {
				if expiredAuction.Id == id {
					continue
				}
			}

			expiredAuctions = append(expiredAuctions, auction.AuctionEntityMongo{id, au.ProductName, au.Category, au.Description, au.Condition, au.Status, au.Timestamp.Unix()})
			au.Status = auction_entity.Completed
			m.auctions[id] = au
			fmt.Printf("Auction %s expired\n", id)
		} else {
			fmt.Printf("Auction %v is still active\n", au)
		}
	}
	return expiredAuctions, nil
}

func (m *MockAuctionRepository) CloseAuction(ctx context.Context, auctionID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	auction, exists := m.auctions[auctionID]
	if !exists {
		return errors.New("auction not found")
	}
	auction.Status = auction_entity.Completed

	m.auctions[auctionID] = auction

	fmt.Printf("Auction %s closed\n", auctionID)
	return nil
}

func (m *MockAuctionRepository) StartAuctionCloser(ctx context.Context) {
	ticker := time.NewTicker(auction.GetCheckInterval())
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			expiredAuctions, err := m.GetExpiredAuctions(ctx)
			if err == nil {
				for _, auction := range expiredAuctions {
					err := m.CloseAuction(ctx, auction.Id)
					if err != nil {
						fmt.Println("Failed to close auction", err)
					} else {
						fmt.Printf("Auction %s closed\n", auction.Id)
					}
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
