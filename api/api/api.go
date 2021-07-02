package api

import (
	"main/db"

	"github.com/gin-gonic/gin"
)

// Setup - call this method to setup routes
func Setup(router *gin.Engine) {

	db.SetupDB()

	apiRoute := router.Group("/api") 
	SetupAuthenAPI(apiRoute)
	SetupPartyAPI(apiRoute)
	// SetupTransactionAPI(router)
}
