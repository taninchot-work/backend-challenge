package background

import (
	"context"
	"github.com/taninchot-work/backend-challenge/internal/core/db"
	"go.mongodb.org/mongo-driver/v2/bson"
	"log"
	"time"
)

func StartUserCountLogger(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Println("Starting background user count logger...")
	for {
		select {
		case <-ticker.C:
			countCtx, cancelCount := context.WithTimeout(context.Background(), 5*time.Second) // timeout for database operations
			collection := db.GetDatabase().Collection("users")

			count, err := collection.CountDocuments(countCtx, bson.M{})
			cancelCount()

			if err != nil {
				log.Printf("Error counting users: %v", err)
			} else {
				log.Printf("Current number of users in DB: %d", count)
			}
		case <-ctx.Done():
			log.Println("Stopping background user count logger.")
			return
		}
	}
}
