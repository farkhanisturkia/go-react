package main

import (
	"os"
	"fmt"

	"go-react-vue/backend/config"
	"go-react-vue/backend/database"
	"go-react-vue/backend/pkg/redis"
	"go-react-vue/backend/routes"
	"go-react-vue/backend/cache"
)

func main() {
	// load config .env
	config.LoadEnv()

	// inisialisasi Redis
	redis.Init()

	//inisialisasi database
	database.InitDB()

	// reset cache
	cache.InvalidateUserListCache()
    cache.InvalidateAllSingleUserCache()

	// Cek argumen CLI
	if len(os.Args) > 1 {
		arg := os.Args[1]

		switch arg {
		case "reset":
			fmt.Println("Mode: RESET users table")
			database.ResetUsers(false) // reset tanpa seed ulang

		case "seed":
			fmt.Println("Mode: SEED only (tanpa reset)")
			database.Seed()

		case "reset-seed":
			fmt.Println("Mode: RESET + SEED ulang")
			database.ResetUsers(true) // true = seed setelah reset

		default:
			fmt.Printf("Argumen tidak dikenal: %s\n", arg)
			fmt.Println("Gunakan salah satu:")
			fmt.Println("  go run .               → normal run")
			fmt.Println("  go run . reset         → reset users lalu run")
			fmt.Println("  go run . seed          → seed ulang lalu run")
			fmt.Println("  go run . reset-seed    → reset + seed lalu run")
			os.Exit(1)
		}
	} else {
		fmt.Println("Mode: Normal run (tanpa reset/seed)")
	}

	//setup router
	r := routes.SetupRouter()

	//mulai server
	r.Run(":" + config.GetEnv("APP_PORT", "3000"))
}
