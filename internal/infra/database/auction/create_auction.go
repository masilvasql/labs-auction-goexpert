package auction

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,

	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	return nil
}

func (ar *AuctionRepository) StartAuctionCloser(ctx context.Context) {
	ticker := time.NewTicker(GetCheckInterval())
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			expiredAuctions, err := ar.GetExpiredAuctions(ctx)
			if err == nil {
				for _, auction := range expiredAuctions {
					err := ar.CloseAuction(ctx, auction.Id)
					if err != nil {
						logger.Error("Failed to close auction", err)
					} else {
						fmt.Printf("Auction %s closed\n", auction.Id)
					}
				}
			} else {
				logger.Error("Error fetching expired auctions", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (ar *AuctionRepository) GetExpiredAuctions(ctx context.Context) ([]AuctionEntityMongo, error) {
	expirationTime := time.Now().Add(-GetAuctionDuration()).Unix()
	filter := bson.M{"status": auction_entity.Active, "timestamp": bson.M{"$lt": expirationTime}}

	var expiredAuctions []AuctionEntityMongo
	cursor, err := ar.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &expiredAuctions); err != nil {
		return nil, err
	}
	return expiredAuctions, nil
}

func (ar *AuctionRepository) CloseAuction(ctx context.Context, auctionID string) error {
	filter := bson.M{"_id": auctionID}
	update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}

	_, err := ar.Collection.UpdateOne(ctx, filter, update)
	return err
}

func GetAuctionDuration() time.Duration {
	duration := os.Getenv("AUCTION_DURATION")
	interval, err := time.ParseDuration(duration)
	if err != nil {
		interval, _ = time.ParseDuration("10s")
	}
	fmt.Println("Auction duration: ", duration)
	return interval
}

func GetCheckInterval() time.Duration {
	interval := os.Getenv("CHECK_INTERVAL")
	duration, err := time.ParseDuration(interval)
	if err != nil {
		duration, _ = time.ParseDuration("2s")
	}
	fmt.Println("Check interval: ", interval)
	return duration
}
