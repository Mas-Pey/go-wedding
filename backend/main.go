package main

import (
	"go-wedding/features"
	"go-wedding/middleware"
	"go-wedding/utils"
	"net/http"

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

	protected := r.Group("/api")
	protected.Use(middleware.RequireAuth)
	protected.GET("/me", func(ctx *gin.Context) {
		userID, _ := ctx.Get("userID")
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Access dashboard successful",
			"user_id": userID,
		})
	})

	r.Run()
}
