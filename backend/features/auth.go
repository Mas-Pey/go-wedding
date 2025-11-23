package features

import (
	"database/sql"
	"go-wedding/domain"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func AuthRoute(r *gin.Engine, db *sql.DB) *gin.Engine {
	r.POST("/sign-up", func(ctx *gin.Context) {
		var signup domain.RequestSign

		if err := ctx.ShouldBindJSON(&signup); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request payload",
			})
			return
		}

		hashedPw, err := bcrypt.GenerateFromPassword([]byte(signup.Pw), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error hashing password => " + err.Error(),
			})
			return
		}

		query := "INSERT INTO users (fullname, email, pw) VALUES ($1, $2, $3) RETURNING id"
		row := db.QueryRow(query, signup.Fullname, signup.Email, string(hashedPw))
		var id int
		err = row.Scan(&id)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create user (email might already exist)",
			})
			return
		}

		ctx.JSON(200, gin.H{
			"message":  "Success add new user",
			"id":       id,
			"fullname": signup.Fullname,
			"email":    signup.Email,
		})
	})

	r.POST("/sign-in", func(ctx *gin.Context) {
		var signin domain.RequestSign
		if err := ctx.ShouldBindJSON(&signin); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request payload ",
			})
			return
		}

		query := "SELECT id, pw FROM users WHERE email = $1"
		row := db.QueryRow(query, signin.Email)
		var id int
		var password string
		err := row.Scan(&id, &password)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Incorrect email or password",
			})
			return
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(password), []byte(signin.Pw))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Incorrect email or password",
			})
			return
		}

		secretKey := os.Getenv("JWT_SECRET")
		claims := jwt.MapClaims{
			"sub": id,
			"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString([]byte(secretKey))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create token",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token":   tokenString,
		})
	})

	return r
}
