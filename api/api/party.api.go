package api

import (
	"main/db"
	"main/model"
	"main/interceptor"

	"time"
	// _ "time"
	"strconv"
	"github.com/gin-gonic/gin"
)

func SetupPartyAPI(router *gin.RouterGroup) {
	partyAPI := router.Group("/party")
	{
		partyAPI.GET("/", interceptor.JwtVerify, list)
		partyAPI.POST("/", interceptor.JwtVerify, createParty)
		partyAPI.PATCH("/:id", interceptor.JwtVerify, editParty)
		partyAPI.DELETE("/:id/delete", interceptor.JwtVerify, deleteParty)

		partyAPI.POST("/join/:id", interceptor.JwtVerify, joinParty)
	}
}

func list(c *gin.Context) {
		// 
	// user := model.User{}
	// party := model.Party{}

	// db.GetDB().First(&party,3)
	// db.GetDB().First(&user,4)

	
	// var users []model.User
	
	//db.GetDB().Find(&parties)//.Model(&parties).Association("JoinUsers").Find(&user)
	// db.GetDB().Model(&party).Association("JoinUsers").Append([]model.User{user})
	// db.GetDB().Model(&party).Association("JoinUsers").Append(user)

	var parties []model.Party
	// db.GetDB().Preload("Parties").Find(&users)
	db.GetDB().Preload("JoinUsers").Preload("Owner").Find(&parties)

	c.JSON(200, gin.H{
			"status": "ok", 
			"result": parties,
		})
}

func createParty(c *gin.Context) {

	user := model.User{}
	party := model.Party{}
	request := model.PartyRequest{}

	db.GetDB().First(&user, c.GetString("jwt_user_id"))
	// party.Name = c.PostForm("name")
	// party.NumberOfMember, _ = strconv.ParseInt(c.PostForm("numberOfMember"), 10, 64)
	
	if c.ShouldBind(&request) == nil {
		party.Name = request.Name
		party.NumberOfMember, _ = strconv.ParseInt(request.NumberOfMember, 10, 64)
		party.CreatedAt = time.Now()
		party.Owner = user

		if err := db.GetDB().Create(&party).Error; err != nil {
			c.JSON(400, gin.H{"status": "fail", "message": "สร้างปาร์ตี้หารไม่สำเร็จ"})
		} else {

			db.GetDB().Model(&party).Association("Owner").Append(&user)
			c.JSON(200, gin.H{"status": "ok", "result": party, "message": "สร้างปาร์ตี้หารเรียบร้อยแล้ว"})
		}
	}  else {
		c.JSON(401, gin.H{"status": "unable to bind data"})
	}
}

func editParty(c *gin.Context) {
	party := model.Party{}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)
	party.ID = uint(id)
	party.Name = c.PostForm("name")
	party.NumberOfMember, _ = strconv.ParseInt(c.PostForm("numberOfMember"), 10, 64)

	if err := db.GetDB().Save(&party) .Error; err != nil {
		c.JSON(400, gin.H{"status": "fail", "message": "แก้ไขปาร์ตี้หารไม่สำเร็จ"})
	} else {
		c.JSON(200, gin.H{"status": "ok", "result": party, "message": "แก้ไขปาร์ตี้หารเรียบร้อยแล้ว"})
	}
}


func deleteParty(c *gin.Context) {
	if err := db.GetDB().Delete(&model.Party{}, c.Param("id")).Error; err != nil { 
		c.JSON(400, gin.H{"status": "fail", "message": "ลบปาร์ตี้หารไม่สำเร็จ"})
	} else {
		c.JSON(200, gin.H{"status": "ok", "message": "ลบปาร์ตี้หารเรียบร้อยแล้ว"})
	}
}

func joinParty(c *gin.Context) {
	user := model.User{}
	party := model.Party{}

	db.GetDB().First(&user, c.GetString("jwt_user_id"))
    db.GetDB().First(&party, c.Param("id"))
	db.GetDB().Preload("JoinUsers", c.GetString("jwt_user_id")).Find(&party)
	
	
	isJoin := db.GetDB().Model(&party).Where("id = ?", c.GetString("jwt_user_id")).Association("JoinUsers").Count() == 1

	if !isJoin {
		// join
		isNumberOfJoin := db.GetDB().Model(&party).Association("JoinUsers").Count()

		if isNumberOfJoin >= party.NumberOfMember {
			c.JSON(400, gin.H{"status": "fail", "message": "ปาร์ตี้หารนี้เต็มแล้ว"})
			return 
		}

		db.GetDB().Model(&party).Association("JoinUsers").Append([]model.User{user})
		c.JSON(200, gin.H{"status": "ok", "result": party, "message": "เข้าร่วมปาร์ตี้หารเรียบร้อยแล้ว"})
		
	} else {
		// leave
		db.GetDB().Model(&party).Association("JoinUsers").Delete([]model.User{user})
		c.JSON(200, gin.H{"status": "ok", "result": party, "message": "ยกเลิกการเข้าร่วมปาร์ตี้หารเรียบร้อยแล้ว"})
	}
	
}
