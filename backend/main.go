package main

import (
	"go-react-vue/backend/config"
	"go-react-vue/backend/database"
	"go-react-vue/backend/pkg/redis"
	"go-react-vue/backend/routes"
)

func main() {

	//load config .env
	config.LoadEnv()

	// inisialisasi Redis
	redis.Init()

	//inisialisasi database
	database.InitDB()

	// seeder
	database.Seed()

	//setup router
	r := routes.SetupRouter()

	//mulai server
	r.Run(":" + config.GetEnv("APP_PORT", "3000"))
}
