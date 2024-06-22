package middleware

import (
	"net/http"
	"saas-billing/config"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	App        string `json:"app"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Nickname   string `json:"nickname"`
	Name       string `json:"name"`
	Picture    string `json:"picture"`
	Locale     string `json:"locale"`
	UpdatedAt  string `json:"updated_at"`
	jwt.RegisteredClaims
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}

			return []byte(config.JWT_SECRET), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "message": err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "message": "Token is not valid"})
			c.Abort()
			return
		}

		// Set the user ID from the token in the context for further use
		claims := token.Claims.(*CustomClaims)
		userID := claims.Subject
		c.Set("user_id", userID)
		c.Set("user_token", tokenString)

		c.Next()
	}
}
