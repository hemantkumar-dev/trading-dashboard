package api

import (
	"log"
	"net/http"
	"time"

	"trading-dashboard/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
)

type loginReq struct {
	Username string `json:"username"`
}

func RegisterAuthRoutes(r *gin.Engine, db *sqlx.DB, jwtSecret string) {
	r.POST("/login", func(c *gin.Context) {
		var req loginReq
		if err := c.BindJSON(&req); err != nil || req.Username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username required"})
			return
		}
		// On login: ensure user exists
		_, err := services.EnsureUser(db, req.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		// create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": req.Username,
			"exp": time.Now().Add(24 * time.Hour).Unix(),
		})
		ts, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
			return
		}
		c.JSON(200, gin.H{"token": ts})
	})
}

// middleware to validate token and set username
func JWTMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		log.Printf("Auth header: %s\n", auth)
		if len(auth) < 7 || auth[:7] != "Bearer " {
			log.Println("Missing or invalid Bearer token format")
			c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
			return
		}
		tkn := auth[7:]
		log.Printf("Token received: %s...\n", tkn[:20])
		log.Printf("JWT Secret used: %s\n", jwtSecret)

		token, err := jwt.Parse(tkn, func(t *jwt.Token) (interface{}, error) {
			// Verify signing method
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			log.Printf("JWT Parse error: %v\n", err)
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token", "details": err.Error()})
			return
		}

		if !token.Valid {
			log.Println("Token is not valid")
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		log.Printf("Token claims: %v\n", claims)
		if sub, ok := claims["sub"].(string); ok {
			c.Set("username", sub)
			c.Next()
			return
		}
		c.AbortWithStatusJSON(401, gin.H{"error": "invalid token claims"})
	}
}
