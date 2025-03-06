package handler

import "github.com/gin-gonic/gin"

func RegisterOrgHandlers(engine *gin.Engine, orgHandler *OrgHandler) {
	org := engine.Group("/api/v1/org")
	{
		org.POST("/create", orgHandler.CreateOrg)
		org.POST("/update", orgHandler.UpdateOrg)
		org.GET("/:id", orgHandler.FindOrgById)
		org.GET("/children/:orgID", orgHandler.FindChildren)
		org.GET("/trees", orgHandler.FindOrgTrees)
		org.GET("/current", orgHandler.FindCurrentOrgTree)
		org.GET("/descendants/:orgID", orgHandler.FindDescendants)
	}
}
