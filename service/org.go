package service

import (
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

func (service *OrgService) CreateOrg(org *model.SysOrg) (bool, error) {
	err := service.txHelper.RunTx(func(tx *sqlx.Tx) error {
		e := service.orgRepo.Insert(tx, org)
		if e != nil {
			service.logger.Errorf("新增组织失败, 组织ID:%d, 错误信息:%v", org.ID, e)
		} else {
			service.logger.Infof("新增组织成功, 组织ID:%d", org.ID)
		}
		return e
	})
	return err != nil, err
}
