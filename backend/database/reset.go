package database

import (
	"fmt"
	"log"

	"go-react-vue/backend/models"
	"gorm.io/gorm"
)

func ResetUsers(seedAfterReset bool) {
	fmt.Println("Starting database reset for users table...")

	err := DB.Exec("TRUNCATE TABLE users").Error
	if err != nil {
		// Fallback kalau TRUNCATE gagal (misal foreign key constraint)
		fmt.Println("TRUNCATE failed, falling back to DELETE + reset auto-increment")
		err = DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.User{}).Error
		if err != nil {
			log.Fatalf("Failed to delete all users: %v", err)
		}

		// Reset auto-increment (MySQL)
		DB.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
	}

	fmt.Println("All users deleted successfully!")

	if seedAfterReset {
		fmt.Println("Re-seeding users after reset...")
		Seed()
	} else {
		fmt.Println("Reset completed. No seeding performed.")
	}
}