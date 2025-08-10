// Excuse most of the comments in this file, it's for my understanding of channels and goroutines.

package main

import (
	"PingMeMaybe/libs/config"
	"PingMeMaybe/libs/db"
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"math/rand"
	"sync"
	"time"
)

const (
	TOTAL_USERS = 100000
	BATCH_SIZE  = 1000
	WORKERS     = 10
)

type UserSeed struct {
	Email                   string     `json:"email"`
	Username                string     `json:"username"`
	FirstName               string     `json:"first_name"`
	LastName                string     `json:"last_name"`
	IsPremiumUser           bool       `json:"is_premium_user"`
	SubscriptionTier        string     `json:"subscription_tier"`
	SubscriptionStartDate   *time.Time `json:"subscription_start_date"`
	SubscriptionEndDate     *time.Time `json:"subscription_end_date"`
	NotificationPreferences string     `json:"notification_preferences"`
	Timezone                string     `json:"timezone"`
	CreatedAt               time.Time  `json:"created_at"`
	LastLoginAt             *time.Time `json:"last_login_at"`
	IsActive                bool       `json:"is_active"`
}

func main() {
	// Initialize config and database
	config.LoadEnv(".")
	dbPool, err := db.InitDBPoolConn()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbPool.Close()

	log.Println("Starting user seeding process...")
	start := time.Now()

	seedUsers(dbPool)

	duration := time.Since(start)
	log.Printf("✅ Seeding completed in %v", duration)

	// Print some stats
	printStats(dbPool)
}

func seedUsers(dbPool *pgxpool.Pool) {
	// Create channels for work distribution
	// Here we set the limit as a buffer for the channel. Basically when the main thread push data to the channel, a worker in the channel immediately pcisk it up
	// If the buffer is full, main thread will block until a worker picks up the data and frees up space in the channel.
	// Very useful in this case as we dont want to overwhelm the postgres conn pool (each worker is performing a bulk write)
	userChan := make(chan []UserSeed, WORKERS)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < WORKERS; i++ {
		wg.Add(1)
		go worker(dbPool, userChan, &wg, i+1) // Go (literally) inside worker function to see how workers stay active until channel exists
	}

	// Generate and send batches
	totalBatches := TOTAL_USERS / BATCH_SIZE
	for batchNum := 0; batchNum < totalBatches; batchNum++ {
		batch := generateUserBatch(BATCH_SIZE, batchNum) // Has a 1000 user batch
		userChan <- batch                                // Push the batch to the channel, one of the workers will pick it up immediately and be responsible for the batch

		if batchNum%10 == 0 {
			// log progress every 10 batches
			log.Printf("Generated batch %d/%d", batchNum+1, totalBatches)
		}
	}

	close(userChan)
	wg.Wait()
}

func worker(dbPool *pgxpool.Pool, userChan <-chan []UserSeed, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for batch := range userChan { // This line here is a BLOCKING line, it will block unitl the channel has data or until channel is closed.
		if err := insertUserBatch(dbPool, batch); err != nil {
			// This db operation is a blocking operation, it will wait for the database to respond.
			// Yep sure we can make this non-blocking by adding "go" keyword, but then what will happen is that it will pick up the next batch immediately.
			// Sounds efficient, but this one worker could potentially have a 100 active db operations at the same time, which will overwhelm the database connection pool.
			log.Printf("❌ Worker %d failed to insert batch: %v", workerID, err)
		} else {
			log.Printf("✅ Worker %d inserted %d users", workerID, len(batch))
		}
	}

	// The range loop syntax above is equivalent to:
	//for {
	//	batch, ok := <-userChan  // Try to receive from channel
	//	if !ok {                 // If channel is closed
	//		break               // Exit the loop
	//	}
	//	// Process batch
	//}
}

func generateUserBatch(size int, batchNum int) []UserSeed {
	batch := make([]UserSeed, size)

	for i := 0; i < size; i++ {
		batch[i] = generateUser(batchNum*size + i)
	}

	return batch
}

func generateUser(index int) UserSeed {
	gofakeit.Seed(int64(index + time.Now().Nanosecond()))

	// Realistic distribution: 80% free, 15% pro, 5% enterprise
	tier := getSubscriptionTier()
	isPremium := tier != "free"

	user := UserSeed{
		Email:                   fmt.Sprintf("user%d@%s", index, gofakeit.DomainName()),
		Username:                fmt.Sprintf("%s_%d", gofakeit.Username(), index),
		FirstName:               gofakeit.FirstName(),
		LastName:                gofakeit.LastName(),
		IsPremiumUser:           isPremium,
		SubscriptionTier:        tier,
		NotificationPreferences: generateNotificationPrefs(),
		Timezone:                getRandomTimezone(),
		CreatedAt: gofakeit.DateRange(
			time.Now().AddDate(-2, 0, 0), // 2 years ago
			time.Now(),
		),
		IsActive: rand.Float32() < 0.95, // 95% active users
	}

	// Set subscription dates for premium users
	if isPremium {
		startDate := gofakeit.DateRange(
			time.Now().AddDate(-1, 0, 0), // 1 year ago
			time.Now().AddDate(0, -1, 0), // 1 month ago
		)
		user.SubscriptionStartDate = &startDate

		// Some subscriptions expire, some are active
		if rand.Float32() < 0.8 { // 80% have active subscriptions
			endDate := startDate.AddDate(1, 0, 0) // 1 year subscription
			user.SubscriptionEndDate = &endDate
		}
	}

	// Set last login for active users
	if user.IsActive && rand.Float32() < 0.8 {
		lastLogin := gofakeit.DateRange(
			time.Now().AddDate(0, 0, -30), // 30 days ago
			time.Now(),
		)
		user.LastLoginAt = &lastLogin
	}

	return user
}

func getSubscriptionTier() string {
	randomNo := rand.Float32()
	switch {
	case randomNo < 0.80:
		return "free"
	case randomNo < 0.95:
		return "pro"
	default:
		return "enterprise"
	}
}

func generateNotificationPrefs() string {
	preferences := []string{
		`{"email": true, "push": true, "sms": false}`,
		`{"email": true, "push": false, "sms": false}`,
		`{"email": false, "push": true, "sms": false}`,
		`{"email": true, "push": true, "sms": true}`,
		`{"email": false, "push": false, "sms": false}`,
	}
	return preferences[rand.Intn(len(preferences))]
}

func getRandomTimezone() string {
	timezones := []string{
		"UTC", "America/New_York", "America/Los_Angeles", "America/Chicago",
		"Europe/London", "Europe/Paris", "Europe/Berlin", "Asia/Tokyo",
		"Asia/Shanghai", "Asia/Kolkata", "Australia/Sydney", "America/Sao_Paulo",
	}
	return timezones[rand.Intn(len(timezones))]
}

func insertUserBatch(dbPool *pgxpool.Pool, users []UserSeed) error {
	ctx := context.Background()

	query := `
		INSERT INTO users (
			email, username, first_name, last_name, is_premium_user,
			subscription_tier, subscription_start_date, subscription_end_date,
			notification_preferences, timezone, created_at, last_login_at, is_active
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	batch := &pgx.Batch{}
	for _, user := range users {
		batch.Queue(query,
			user.Email,
			user.Username,
			user.FirstName,
			user.LastName,
			user.IsPremiumUser,
			user.SubscriptionTier,
			user.SubscriptionStartDate,
			user.SubscriptionEndDate,
			user.NotificationPreferences,
			user.Timezone,
			user.CreatedAt,
			user.LastLoginAt,
			user.IsActive,
		)
	}

	results := dbPool.SendBatch(ctx, batch) // BLOCKING network call
	defer results.Close()

	// Process all batch results
	for i := 0; i < len(users); i++ {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("failed to insert user %d: %w", i, err)
		}
	}

	return nil
}

func printStats(dbPool *pgxpool.Pool) {
	ctx := context.Background()

	// Total users
	var totalUsers int
	err := dbPool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&totalUsers)
	if err != nil {
		log.Printf("Error getting total users: %v", err)
		return
	}

	// premium users breakdown
	var freeUsers, proUsers, enterpriseUsers int
	dbPool.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE subscription_tier = 'free'").Scan(&freeUsers)
	dbPool.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE subscription_tier = 'pro'").Scan(&proUsers)
	dbPool.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE subscription_tier = 'enterprise'").Scan(&enterpriseUsers)

	log.Println("\n📊 Seeding Statistics:")
	log.Printf("Total Users: %d", totalUsers)
	log.Printf("Free Users: %d (%.1f%%)", freeUsers, float64(freeUsers)/float64(totalUsers)*100)
	log.Printf("Pro Users: %d (%.1f%%)", proUsers, float64(proUsers)/float64(totalUsers)*100)
	log.Printf("Enterprise Users: %d (%.1f%%)", enterpriseUsers, float64(enterpriseUsers)/float64(totalUsers)*100)
}
