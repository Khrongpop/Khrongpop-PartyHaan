package interceptor

import (
	"fmt"
	"main/model"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var secretKey = "87654321"

func JwtSign(payload model.User) string {
	atClaims := jwt.MapClaims{}

	// Payload begin
	atClaims["id"] = payload.ID
	atClaims["email"] = payload.Email

	date := time.Hour * 24 + 
			time.Minute * 60 + 
			time.Second * 60
			
	atClaims["exp"] = time.Now().Add(date * 7).Unix()
	// Payload end

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, _ := at.SignedString([]byte(secretKey))
	return token

}

func JwtVerify(c *gin.Context) {
	tokenString := strings.Split(c.Request.Header["Authorization"][0], " ")[1]
	fmt.Println(tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)	

		userID := fmt.Sprintf("%v", claims["id"])
		email := fmt.Sprintf("%v", claims["jwt_email"])

		c.Set("jwt_user_id", userID)
		c.Set("jwt_email", email)


		c.Next()
	} else {
		c.JSON(http.StatusOK, gin.H{"result": "fail", "message": "invalid token", "error": err})
		c.Abort()
	}
}
