package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"zax/model"
	"zax/repository"
	"zax/util"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type OrgService struct {
	logger   *zap.SugaredLogger
	txHelper *util.TxHelper
	orgRepo  *repository.OrgRepository
}

func NewOrgService(logger *zap.SugaredLogger, txHelper *util.TxHelper, orgRepo *repository.OrgRepository) *OrgService {
	return &OrgService{logger: logger, txHelper: txHelper, orgRepo: orgRepo}
}

// 新增组织
func (service *OrgService) CreateOrg(org *model.SysOrg) (bool, error) {
	// 计算 code
	orgCode, err := service.GenerateOrgCode(org.ParentID)
	if err != nil {
		return false, util.NewZaxError(err.Error())
	}
	org.Code = orgCode
	// 生成 ID
	// TODO
	err = service.txHelper.RunTx(func(tx *sqlx.Tx) error {
		e := service.orgRepo.Insert(tx, []*model.SysOrg{org})
		if e != nil {
			service.logger.Errorf("新增组织失败, 错误信息:%v. 数据:%s", e, org.ToString())
			return util.NewZaxError("新增组织失败")
		} else {
			service.logger.Infof("新增组织成功, 数据:%s", org.ToString())
		}
		return nil
	})
	return err != nil, err
}

// 更新组织信息
func (service *OrgService) UpdateOrg(org *model.SysOrg) (bool, error) {
	currentOrg, err := service.orgRepo.SelectById(org.ID)
	if err != nil {
		return false, util.NewZaxError(fmt.Sprintf("查询组织失败, 组织ID:%d", org.ID))
	}
	var preUpdateOrgs []*model.SysOrg
	if org.ParentID == currentOrg.ParentID {
		org.UpdateTime = time.Now()
		org.UpdateBy = "admin"
		preUpdateOrgs = append(preUpdateOrgs, org)
	} else {
		// 检查新的父组织是否存在
		if org.ParentID != 0 {
			parentOrg, err := service.orgRepo.SelectById(org.ParentID)
			if err != nil {
				return false, util.NewZaxError(fmt.Sprintf("查询父组织失败, 父组织ID:%d", org.ParentID))
			}
			// 检查循环引用
			if parentOrg.ID == org.ID {
				return false, util.NewZaxError(fmt.Sprintf("父组织不能是自己, 父组织ID:%d", org.ParentID))
			}
			// 检查父组织是否是自己的子组织
			if len(parentOrg.Code) > len(org.Code) && parentOrg.Code[:len(org.Code)] == org.Code {
				return false, util.NewZaxErrorf("父组织不能是自己的子组织.当前组织Code:%s, 目标父组织Code:%s", org.Code, parentOrg.Code)
			}
		}
		// 为组织生成新 code
		newOrgCode, err := service.GenerateOrgCode(org.ParentID)
		if err != nil {
			return false, util.NewZaxError(err.Error())
		}
		org.Code = newOrgCode
		org.UpdateTime = time.Now()
		org.UpdateBy = "admin"
		preUpdateOrgs = append(preUpdateOrgs, org)
		// 重新计算子孙组织的 code
		updatedDescendants, err := service.clacDescendantCode(currentOrg.Code, newOrgCode)
		if err != nil {
			return false, util.NewZaxError(err.Error())
		}
		preUpdateOrgs = append(preUpdateOrgs, updatedDescendants...)
	}
	// 更新所有发生变动的组织
	service.logger.Infof("组织ID:%d, 旧编码:%s, 新编码:%s.本次要更新的组织有 %d 个", org.ID, currentOrg.Code, org.Code, len(preUpdateOrgs))
	err = service.txHelper.RunTx(func(tx *sqlx.Tx) error {
		for _, org := range preUpdateOrgs {
			e := service.orgRepo.UpdateSelective(tx, org)
			if e != nil {
				service.logger.Errorf("更新组织失败, 错误信息:%v. 数据:%s", e, org.ToString())
				return util.NewZaxError(fmt.Sprintf("更新组织失败, 错误信息:%v. 数据:%s", e, org.ToString()))
			}
		}
		return nil
	})
	return err != nil, err
}

// 查询组织
func (service *OrgService) FindOrgById(id int64) (*model.SysOrg, error) {
	return service.orgRepo.SelectById(id)
}

// 查询子组织
func (service *OrgService) FindChildren(parentID int64) ([]*model.SysOrg, error) {
	return service.orgRepo.SelectByParentID(parentID)
}

// 查询组织树
func (service *OrgService) FindOrgTrees() ([]*model.SysOrg, error) {
	// 查找所有根组织
	rootOrgs, err := service.orgRepo.SelectByParentID(0)
	if err != nil {
		return nil, util.NewZaxError("查询根组织失败")
	}
	// 查找所有子孙组织
	var orgTrees []*model.SysOrg
	for _, rootOrg := range rootOrgs {
		descendants, err := service.findDescendants(rootOrg.ID)
		if err != nil {
			return nil, util.NewZaxError(err.Error())
		}
		rootOrg.Children = descendants
		orgTrees = append(orgTrees, rootOrg)
	}
	return orgTrees, nil
}

// 递归查找子孙组织
func (service *OrgService) findDescendants(parentID int64) ([]*model.SysOrg, error) {
	orgs, err := service.orgRepo.SelectByParentID(parentID)
	if err != nil {
		return nil, util.NewZaxError("查询子孙组织失败")
	}
	for i := range orgs {
		descendants, err := service.findDescendants(orgs[i].ID)
		if err != nil {
			return nil, util.NewZaxError("查询子孙组织失败")
		}
		orgs[i].Children = descendants
	}
	return orgs, nil
}

// 生成组织编码
func (service *OrgService) GenerateOrgCode(parentID int64) (string, error) {
	var parentOrgCode string
	if parentID != 0 {
		org, err := service.orgRepo.SelectById(parentID)
		if err != nil {
			service.logger.Errorf("查询父组织失败, parentID:%d, 错误信息:%v", parentID, err)
			return "", util.NewZaxError(fmt.Sprintf("查询父组织失败, parentID:%d", parentID))
		}
		parentOrgCode = org.Code
	}
	maxCode, err := service.orgRepo.SelectMaxCode(parentID)
	if err != nil {
		return fmt.Sprintf("%s0001", parentOrgCode), nil
	}
	var nextCode int
	if maxCode == "" {
		nextCode = 1
	} else {
		maxCodeSerial, _ := strconv.Atoi(maxCode[len(maxCode)-4:])
		nextCode = maxCodeSerial + 1
	}
	return fmt.Sprintf("%s%04d", parentOrgCode, nextCode), nil
}

// 更新子孙组织的编码
func (service *OrgService) clacDescendantCode(oldOrgCode string, newOrgCode string) ([]*model.SysOrg, error) {
	var updatedDescendants []*model.SysOrg
	descendants, err := service.orgRepo.SelectDescendants(oldOrgCode)
	if err != nil {
		return nil, util.NewZaxErrorf("查找子孙组织失败, 组织Code:%s", oldOrgCode)
	}
	for _, descendant := range descendants {
		descendantNewCode := strings.Replace(descendant.Code, oldOrgCode, newOrgCode, 1)
		descendant.Code = descendantNewCode
		descendant.UpdateTime = time.Now()
		descendant.UpdateBy = "admin"
		updatedDescendants = append(updatedDescendants, descendant)
	}
	return updatedDescendants, nil
}
