package model

import "time"

type SysOrg struct {
	ID         int64     `db:"id" json:"id"`                  // 组织ID
	Code       string    `db:"code" json:"code"`              // 组织编码. 4位一级. 0001,00010001,000100010001,以此类推
	Name       string    `db:"name" json:"name"`              // 组织名称
	NameAbbr   string    `db:"name_abbr" json:"nameAbbr"`     // 组织名称简称
	Comment    string    `db:"comment" json:"comment"`        // 组织备注
	ParentID   int64     `db:"parent_id" json:"parentId"`     // 父级组织ID. 0表示没有父组织
	IsDeleted  bool      `db:"is_deleted" json:"isDeleted"`   // 逻辑删除标记
	CreateTime time.Time `db:"create_time" json:"createTime"` // 创建时间
	CreateBy   string    `db:"create_by" json:"createBy"`     // 创建人
	UpdateTime time.Time `db:"update_time" json:"updateTime"` // 信息更新时间
	UpdateBy   string    `db:"update_by" json:"updateBy"`     // 信息更新人
}

type SysOrgTree struct {
	SysOrg
	Children []SysOrgTree `json:"children"`
}
