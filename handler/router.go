package handler

import "github.com/gin-gonic/gin"

func RegisterOrgHandlers(engine *gin.Engine, orgHandler *OrgHandler) {
	org := engine.Group("/api/v1/org")
	{
		org.POST("/create", orgHandler.CreateOrg)
		org.GET("/:id", orgHandler.FindOrgById)
		org.GET("/children/:parentID", orgHandler.FindChildren)
		org.GET("/trees", orgHandler.FindOrgTrees)
	}
}
