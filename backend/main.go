package main

import (
	"go-react/backend/config"
	"go-react/backend/database"
	"go-react/backend/pkg/redis"
	"go-react/backend/routes"
)

func main() {

	//load config .env
	config.LoadEnv()

	// inisialisasi Redis
	redis.Init()

	//inisialisasi database
	database.InitDB()

	//setup router
	r := routes.SetupRouter()

	//mulai server
	r.Run(":" + config.GetEnv("APP_PORT", "3000"))
}
