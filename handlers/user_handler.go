package handlers

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type bodyUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterUserHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body bodyUser
		c.Bind(&body)
		query := `SELECT username FROM users WHERE username=$1`
		row := db.QueryRow(query, body.Username)
		var username string

		err := row.Scan(&username)
		if err != nil {
			if err != sql.ErrNoRows {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "server error",
				})
				return
			}

		}
		hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.MinCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error occurred while generating hash",
			})
			return
		}
		if len(username) > 0 {
			c.JSON(http.StatusConflict, gin.H{
				"message": "username already exists",
			})
			return
		}

		query = `INSERT INTO users (username, password) VALUES ($1, $2)`
		_, err = db.Exec(query, body.Username, string(hash))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "server error",
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{

			"message": "user registered successfully",
		})
	}
}

func LoginUserHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body bodyUser
		c.Bind(&body)
		query := `SELECT username, password FROM users WHERE username=$1`
		row := db.QueryRow(query, body.Username)
		var username, password string
		err := row.Scan(&username, &password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "server error",
			})
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(body.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "wrong password",
			})
			return

		}
		if len(username) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "username does not exist",
			})
			return

		}
		err = godotenv.Load()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to load env file",
			})
			return
		}
		key := os.Getenv("secret_key")
		key_byte := []byte(key)
		token := jwt.New(jwt.SigningMethodHS256)
		token.Claims = jwt.MapClaims{
			"name": body.Username,
			"exp":  time.Now().Add(time.Hour).Unix(),
			"iat":  time.Now().Unix(),
		}

		tokenString, err := token.SignedString(key_byte)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "login successfull",
			"token":   tokenString,
		})
	}
}
