package handler

import (
	"strconv"
	"zax/service"
	"zax/util"

	"zax/model"

	"github.com/gin-gonic/gin"
)

type OrgHandler struct {
	orgService *service.OrgService
}

func NewOrgHandler(orgService *service.OrgService) *OrgHandler {
	return &OrgHandler{orgService: orgService}
}

// /org/create
func (h *OrgHandler) CreateOrg(c *gin.Context) {
	org := &model.SysOrg{}
	if err := c.ShouldBindJSON(&org); err != nil {
		Error(c, err.Error())
		return
	}
	username := c.GetString("username")
	org.CreateBy = &username
	org.CreateTime = util.Now()
	_, err := h.orgService.CreateOrg(org)
	if err != nil {
		Error(c, err.Error())
		return
	}
	Success(c, org)
}

// /org/update
func (h *OrgHandler) UpdateOrg(c *gin.Context) {
	org := &model.SysOrg{}
	if err := c.ShouldBindJSON(&org); err != nil {
		Error(c, err.Error())
		return
	}
	username := c.GetString("username")
	org.UpdateBy = &username
	org.UpdateTime = util.Now()
	_, err := h.orgService.UpdateOrg(org)
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

// /org/children/:orgID
func (h *OrgHandler) FindChildren(c *gin.Context) {
	orgID, _ := strconv.ParseInt(c.Param("orgID"), 10, 64)
	orgs, err := h.orgService.FindChildren(orgID)
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

// /org/current
func (h *OrgHandler) FindCurrentOrgTree(c *gin.Context) {
	orgID, _ := strconv.ParseInt(c.DefaultQuery("orgID", "0"), 10, 64)
	tree, err := h.orgService.FindCurrentOrgTree(orgID)
	if err != nil {
		Error(c, err.Error())
		return
	}
	Success(c, tree)
}

// /org/descendants/:orgID
func (h *OrgHandler) FindDescendants(c *gin.Context) {
	orgID, _ := strconv.ParseInt(c.Param("orgID"), 10, 64)
	descendants, err := h.orgService.FindDescendants(orgID)
	if err != nil {
		Error(c, err.Error())
		return
	}
	Success(c, descendants)
}
