package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// func UserAuth(c *gin.Context) {
// 	// Get the cookie off req
// 	tokenString, err := c.Cookie("UserAuthorization")
// 	if err != nil {
// 		c.Redirect(http.StatusFound, "/user/login")
// 		c.Abort()
// 		return
// 	}

// 	// Decode/validate it
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return []byte(os.Getenv("SECRET")), nil
// 	})
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"message": "Invalid token",
// 		})
// 		c.Abort()
// 		return
// 	}

// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		// Check the exp
// 		if float64(time.Now().Unix()) > claims["exp"].(float64) {
// 			c.JSON(http.StatusUnauthorized, gin.H{
// 				"message": "Token expired",
// 			})
// 			c.Abort()
// 			return
// 		}

// 		// Attach user ID to the context
// 		c.Set("userid", claims["sub"])
// 		c.Next()
// 	} else {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"message": "Invalid token",
// 		})
// 		c.Abort()
// 		return
// 	}
// }

// func UserAuth(c *gin.Context) {
// 	tokenString, err := c.Cookie("UserAuthorization")
// 	if err != nil {
// 		log.Printf("Error retrieving cookie: %v", err)
// 		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization required"})
// 		c.Abort()
// 		return
// 	}

// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			log.Printf("Unexpected signing method: %v", token.Header["alg"])
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return []byte(os.Getenv("SECRET")), nil
// 	})

// 	if err != nil {
// 		log.Printf("Token parsing error: %v", err)
// 		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
// 		c.Abort()
// 		return
// 	}

// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		log.Printf("Claims: %v", claims)
// 		if float64(time.Now().Unix()) > claims["exp"].(float64) {
// 			c.JSON(http.StatusUnauthorized, gin.H{"message": "Token expired"})
// 			c.Abort()
// 			return
// 		}
// 		userID, ok := claims["sub"].(string)
// 		if !ok {
// 			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
// 			c.Abort()
// 			return
// 		}
// 		c.Set("userID", userID)
// 		c.Next()
// 	} else {
// 		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
// 		c.Abort()
// 	}
// }

func UserAuth(c *gin.Context) {
	tokenString, err := c.Cookie("UserAuthorization")
	if err != nil {
		log.Printf("Error retrieving cookie: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization required"})
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		log.Printf("Token parsing error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		c.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Printf("Claims: %v", claims)
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Token expired"})
			c.Abort()
			return
		}
		userIDStr, ok := claims["sub"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
			c.Abort()
			return
		}
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid User ID"})
			c.Abort()
			return
		}
		c.Set("userID", uint(userID))
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		c.Abort()
	}
}
