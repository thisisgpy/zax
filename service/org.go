package service

import (
	"fmt"
	"strconv"
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

func (service *OrgService) UpdateOrgCode(orgId int64, parentId int64) (map[string]string, error) {
	return nil, nil
}

func (service *OrgService) CreateOrg(org *model.SysOrg) (bool, error) {
	orgCode, err := service.GenerateOrgCode(org.ParentID)
	if err != nil {
		return false, util.NewZaxError(err.Error())
	}
	org.Code = orgCode
	err = service.txHelper.RunTx(func(tx *sqlx.Tx) error {
		e := service.orgRepo.Insert(tx, []*model.SysOrg{org})
		if e != nil {
			service.logger.Errorf("新增组织失败, 错误信息:%v. 数据:%s", e, org.ToString())
			return util.NewZaxError("新增组织失败")
		} else {
			service.logger.Infof("新增组织成功, 数据:%s", org.ToString())
			return nil
		}
	})
	return err != nil, err
}

func (service *OrgService) UpdateOrg(org *model.SysOrg) (bool, error) {
	currentOrg, err := service.orgRepo.SelectById(org.ID)
	if err != nil {
		return false, util.NewZaxError(fmt.Sprintf("查询组织失败, 组织ID:%d", org.ID))
	}
	if org.Code != currentOrg.Code {
		// codeMap, err := service.UpdateOrgCode(org.ID, org.ParentID)
		// TODO
	}

	err = service.txHelper.RunTx(func(tx *sqlx.Tx) error {
		e := service.orgRepo.UpdateSelective(tx, org)
		if e != nil {
			service.logger.Errorf("更新组织失败, 错误信息:%v. 数据:%s", e, org.ToString())
			return util.NewZaxErrorf("更新组织失败, 组织ID:%d", org.ID)
		} else {
			service.logger.Infof("更新组织成功, 数据:%s", org.ToString())
			return nil
		}
	})
	return err != nil, err
}
