package database

import (
	"fmt"
	"log"
	"time"

	"go-react-vue/backend/models"
	"go-react-vue/backend/helpers"
)

// Seed runs database seeders
func Seed() {
	fmt.Println("Starting database seeding...")

	var users []models.User
	now := time.Now()

	// Generate 1000 users
	for i := 1; i <= 1000; i++ {
		username := fmt.Sprintf("username%d", i)
		email := fmt.Sprintf("email%d@gmail.com", i)

		// Cek apakah username atau email sudah ada
		var existingCount int64
		err := DB.Model(&models.User{}).
			Where("username = ? OR email = ?", username, email).
			Count(&existingCount).Error

		if err != nil {
			log.Printf("Error checking existing user %d: %v", i, err)
			continue
		}

		if existingCount > 0 {
			fmt.Printf("Skipped user%d (username or email already exists)\n", i)
			continue
		}

		user := models.User{
			Name:      fmt.Sprintf("user%d", i),
			Username:  username,
			Email:     email,
			Password:  helpers.HashPassword("password"),
			CreatedAt: now,
			UpdatedAt: now,
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		fmt.Println("No new users to seed (all already exist)")
		return
	}

	// Bulk insert hanya untuk user yang belum ada
	result := DB.Create(&users)
	if result.Error != nil {
		log.Fatalf("Failed to seed users: %v", result.Error)
	}

	fmt.Printf("Successfully seeded %d new users\n", len(users))
}