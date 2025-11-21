package main

import (
	"go-wedding/features"
	"go-wedding/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db := utils.ConnectDB()
	defer db.Close()
	db.SetConnMaxIdleTime(0)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r = features.AuthRoute(r, db)
	r.Run()
}
