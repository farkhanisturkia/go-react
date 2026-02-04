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
	// Cek apakah sudah ada data (opsional, supaya tidak duplicate setiap kali run
	var count int64
	DB.Model(&models.User{}).Count(&count)

	if count > 0 {
		fmt.Println("Seeding skipped: users table already has data")
		return
	}

	fmt.Println("Starting database seeding...")

	var users []models.User
	now := time.Now()

	// Generate 1000 users
	for i := 1; i <= 1000; i++ {
		user := models.User{
			Name:      fmt.Sprintf("user%d", i),
			Username:  fmt.Sprintf("username%d", i),
			Email:     fmt.Sprintf("email%d@gmail.com", i),
			Password:  helpers.HashPassword("password"),
			CreatedAt: now,
			UpdatedAt: now,
		}
		users = append(users, user)
	}

	// Bulk insert
	result := DB.Create(&users)
	if result.Error != nil {
		log.Fatalf("Failed to seed users: %v", result.Error)
	}

	fmt.Printf("Successfully seeded %d users\n", len(users))
}
