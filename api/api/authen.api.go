package api

import (
	"main/db"
	"main/interceptor"
	"main/model"

	// "time"
	// _ "time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SetupAuthenAPI(router *gin.RouterGroup) {
	authenAPI := router.Group("/auth")
	{
		authenAPI.POST("/login", login)
		authenAPI.POST("/register", register)
		authenAPI.POST("/me", interceptor.JwtVerify, getProfile)
	}
}

func login(c *gin.Context) {
	user := model.User{}

	if c.ShouldBind(&user) == nil {
		queryUser := model.User{}
		if err := db.GetDB().First(&queryUser, "email = ?", user.Email).Error; err != nil {
			c.JSON(400, gin.H{"status": "fail", "error": err,  "message": "อีเมลหรือรหัสผ่านผิดพลาด"})
		} else if checkPasswordHash(user.Password, queryUser.Password) == false {
			c.JSON(400, gin.H{"status": "fail", "message": "อีเมลหรือรหัสผ่านผิดพลาด"})
		} else {
			token := interceptor.JwtSign(queryUser)
			c.JSON(200, gin.H{
				"status": "ok", 
				"token": token,
				"result": queryUser,
			})
		}

	} else {
		c.JSON(401, gin.H{"status": "unable to bind data"})
	}
}

func register(c *gin.Context) {
	user := model.User{}
	if c.ShouldBind(&user) == nil {
		user.Password, _ = hashPassword(user.Password)
		if err := db.GetDB().Create(&user).Error; err != nil {
			c.JSON(400, gin.H{"status": "create fail", "error": err})
		} else {
			// Register Successful
			token := interceptor.JwtSign(user)
			c.JSON(200, gin.H{
				"status": "ok", 
				"token": token,
				"result": user,
			})
		}
	} else {
		c.JSON(401, gin.H{"status": "unable to bind data"})
	}
}

func getProfile(c *gin.Context) {
	user := model.User{}
	db.GetDB().First(&user, c.GetString("jwt_user_id"))
	c.JSON(200, gin.H{
				"status": "ok", 
				"result": user,
			})
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
