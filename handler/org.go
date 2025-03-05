package handler

import (
	"strconv"
	"zax/service"

	"zax/model"

	"github.com/gin-gonic/gin"
)

type OrgHandler struct {
	orgService *service.OrgService
}

func NewOrgHandler(orgService *service.OrgService) *OrgHandler {
	return &OrgHandler{orgService: orgService}
}

func (h *OrgHandler) CreateOrg(c *gin.Context) {
	org := &model.SysOrg{}
	if err := c.ShouldBindJSON(&org); err != nil {
		Error(c, err.Error())
		return
	}
	_, err := h.orgService.CreateOrg(org)
	if err != nil {
		Error(c, err.Error())
		return
	}
	Success(c, org)
}

// /org/:id
func (h *OrgHandler) FindOrgById(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	org, err := h.orgService.FindOrgById(id)
	if err != nil {
		Error(c, err.Error())
		return
	}
	Success(c, org)
}

// /org/children/:parentID
func (h *OrgHandler) FindChildren(c *gin.Context) {
	parentID, _ := strconv.ParseInt(c.Param("parentID"), 10, 64)
	orgs, err := h.orgService.FindChildren(parentID)
	if err != nil {
		Error(c, err.Error())
		return
	}
	Success(c, orgs)
}

// /org/trees?rootOrgID=
func (h *OrgHandler) FindOrgTrees(c *gin.Context) {
	rootOrgID, _ := strconv.ParseInt(c.DefaultQuery("rootOrgID", "0"), 10, 64)
	trees, err := h.orgService.FindOrgTrees(rootOrgID)
	if err != nil {
		Error(c, err.Error())
		return
	}
	Success(c, trees)
}
