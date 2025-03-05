package service

import (
	"fmt"
	"strconv"
	"strings"
	"zax/model"
	"zax/repository"
	"zax/util"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type OrgInterface interface {
	CreateOrg(org *model.SysOrg) (bool, error)
	UpdateOrg(org *model.SysOrg) (bool, error)
	FindOrgById(id int64) (*model.SysOrg, error)
	FindChildren(parentID int64) ([]*model.SysOrg, error)
	FindOrgTrees(rootOrgID int64) ([]*model.SysOrg, error)
	FindCurrentOrgTree(orgID int64) (*model.SysOrg, error)
}

type OrgService struct {
	logger   *zap.SugaredLogger
	idGen    *util.Snowflake
	txHelper *util.TxHelper
	orgRepo  *repository.OrgRepository
}

func NewOrgService(logger *zap.SugaredLogger, idGen *util.Snowflake, txHelper *util.TxHelper, orgRepo *repository.OrgRepository) *OrgService {
	return &OrgService{logger: logger, idGen: idGen, txHelper: txHelper, orgRepo: orgRepo}
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
	org.ID = service.idGen.GenerateID()
	// 设置创建时间
	org.CreateTime = util.Now()
	// 设置创建人
	org.CreateBy = util.GetUser()
	// 执行事务
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
		org.UpdateTime = util.Now()
		org.UpdateBy = util.GetUser()
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
		org.UpdateTime = util.Now()
		org.UpdateBy = util.GetUser()
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
	org, err := service.orgRepo.SelectById(id)
	if err != nil {
		service.logger.Errorf("查询组织失败, 错误信息:%v. 组织ID:%d", err, id)
		return nil, util.NewZaxErrorf("组织不存在。ID:%d", id)
	}
	return org, nil
}

// 查询子组织
func (service *OrgService) FindChildren(parentID int64) ([]*model.SysOrg, error) {
	return service.orgRepo.SelectByParentID(parentID)
}

// 查询组织树. 如果 rootOrgID 为 0, 则查询所有组织树. 否则查询指定组织树
func (service *OrgService) FindOrgTrees(rootOrgID int64) ([]*model.SysOrg, error) {
	var rootOrgs []*model.SysOrg
	if rootOrgID == 0 {
		// 查找所有根组织
		orgs, err := service.orgRepo.SelectByParentID(0)
		if err != nil {
			return nil, util.NewZaxError("查询根组织失败")
		}
		rootOrgs = orgs
	} else {
		org, err := service.orgRepo.SelectById(rootOrgID)
		if err != nil {
			return nil, util.NewZaxErrorf("查询根组织失败, 根组织ID:%d", rootOrgID)
		}
		if org.ParentID != 0 {
			return nil, util.NewZaxErrorf("不是根组织, 组织ID:%d", rootOrgID)
		}
		rootOrgs = append(rootOrgs, org)
	}
	// 查找所有子孙组织
	var orgTrees []*model.SysOrg
	for _, rootOrg := range rootOrgs {
		descendants, err := service.FindDescendants(rootOrg.ID)
		if err != nil {
			return nil, util.NewZaxError(err.Error())
		}
		rootOrg.Children = descendants
		orgTrees = append(orgTrees, rootOrg)
	}
	return orgTrees, nil
}

// 查询指定组织所在的组织树
func (service *OrgService) FindCurrentOrgTree(orgID int64) (*model.SysOrg, error) {
	// 查询指定组织的 code
	currentOrg, err := service.orgRepo.SelectById(orgID)
	if err != nil {
		return nil, util.NewZaxErrorf("查询组织失败, 组织ID:%d", orgID)
	}
	// 确定当前组织的根组织 code
	rootCode := currentOrg.Code[0:4]
	// 查询根组织
	org, err := service.orgRepo.SelectByCode(rootCode)
	if err != nil {
		return nil, util.NewZaxErrorf("查询根组织失败, 根组织Code:%s", rootCode)
	}
	// 查询根组织所在的组织树
	trees, err := service.FindOrgTrees(org.ID)
	if err != nil {
		return nil, util.NewZaxError(err.Error())
	}
	return trees[0], nil
}

// 查找指定组织的所有子孙组织
func (service *OrgService) FindDescendants(orgID int64) ([]*model.SysOrg, error) {
	orgs, err := service.orgRepo.SelectByParentID(orgID)
	if err != nil {
		return nil, util.NewZaxError("查询子孙组织失败")
	}
	for i := range orgs {
		descendants, err := service.FindDescendants(orgs[i].ID)
		if err != nil {
			return nil, util.NewZaxError("查询子孙组织失败")
		}
		orgs[i].Children = descendants
	}
	return orgs, nil
}

// 生成组织编码. parentID 是当前组织的父组织ID，组织 code 一定是根据父级的 code 来生成.
func (service *OrgService) GenerateOrgCode(parentID int64) (string, error) {
	var parentOrgCode string
	// parentID 不为 0，表示当前要新增的组织不是顶级组织, 先查找父组织的编码
	if parentID != 0 {
		org, err := service.orgRepo.SelectById(parentID)
		if err != nil {
			service.logger.Errorf("查询父组织失败, parentID:%d, 错误信息:%v", parentID, err)
			return "", util.NewZaxError(fmt.Sprintf("查询父组织失败, parentID:%d", parentID))
		}
		parentOrgCode = org.Code
	}
	// 查找父组织的子组织中当前的最大编码
	maxCode, err := service.orgRepo.SelectMaxCode(parentID)
	// 找不到已使用的最大编码则说明当前要新增的组织是第一个子组织, 编码为父组织编码+0001
	if err != nil {
		return fmt.Sprintf("%s0001", parentOrgCode), nil
	}
	var nextCode int
	// 防御性判断
	if maxCode == "" {
		nextCode = 1
	} else {
		// 截取已使用最大编码的最后四位，并转换为整数
		maxCodeSerial, _ := strconv.Atoi(maxCode[len(maxCode)-4:])
		// 加 1 得到下一个编码
		nextCode = maxCodeSerial + 1
	}
	// 组合父级组织编码和下一个编码序号(不足4位前面补0)
	return fmt.Sprintf("%s%04d", parentOrgCode, nextCode), nil
}

// 更新当前组织的子孙组织的编码. oldOrgCode 为当前组织的旧编码, newOrgCode 为当前组织的新编码
func (service *OrgService) clacDescendantCode(oldOrgCode string, newOrgCode string) ([]*model.SysOrg, error) {
	var updatedDescendants []*model.SysOrg
	descendants, err := service.orgRepo.SelectDescendants(oldOrgCode)
	if err != nil {
		return nil, util.NewZaxErrorf("查找子孙组织失败, 组织Code:%s", oldOrgCode)
	}
	for _, descendant := range descendants {
		descendantNewCode := strings.Replace(descendant.Code, oldOrgCode, newOrgCode, 1)
		descendant.Code = descendantNewCode
		descendant.UpdateTime = util.Now()
		descendant.UpdateBy = util.GetUser()
		updatedDescendants = append(updatedDescendants, descendant)
	}
	return updatedDescendants, nil
}
