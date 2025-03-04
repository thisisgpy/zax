package handler

import "github.com/gin-gonic/gin"

func RegisterOrgHandlers(engine *gin.Engine, orgHandler *OrgHandler) {
	v1 := engine.Group("/v1")
	{
		v1.GET("/org/:id", orgHandler.FindOrgById)
		v1.GET("/org/children/:parentID", orgHandler.FindChildren)
		v1.GET("/org/trees", orgHandler.FindOrgTrees)
	}
}
