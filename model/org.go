package model

import (
	"fmt"
	"time"
)

type SysOrg struct {
	ID         int64      `db:"id" json:"id"`                  // 组织ID
	Code       string     `db:"code" json:"code"`              // 组织编码. 4位一级. 0001,00010001,000100010001,以此类推
	Name       string     `db:"name" json:"name"`              // 组织名称
	NameAbbr   *string    `db:"name_abbr" json:"nameAbbr"`     // 组织名称简称
	Comment    *string    `db:"comment" json:"comment"`        // 组织备注
	ParentID   *int64     `db:"parent_id" json:"parentId"`     // 父级组织ID. 0表示没有父组织
	CreateTime *time.Time `db:"create_time" json:"createTime"` // 创建时间
	CreateBy   *string    `db:"create_by" json:"createBy"`     // 创建人
	UpdateTime *time.Time `db:"update_time" json:"updateTime"` // 信息更新时间
	UpdateBy   *string    `db:"update_by" json:"updateBy"`     // 信息更新人
	Children   []*SysOrg  `json:"children"`
}

func (o *SysOrg) ToString() string {
	return fmt.Sprintf("SysOrg{ID:%d, Code:%s, Name:%s, NameAbbr:%v, Comment:%v, ParentID:%v, CreateTime:%s, CreateBy:%v, UpdateTime:%s, UpdateBy:%v, Children:%v}",
		o.ID, o.Code, o.Name, o.NameAbbr, o.Comment, o.ParentID, o.CreateTime, o.CreateBy, o.UpdateTime, o.UpdateBy, o.Children)
}

func (target *SysOrg) MapNotNull(source *SysOrg) {
	if source.Name != "" {
		target.Name = source.Name
	}
	if source.NameAbbr != nil {
		target.NameAbbr = source.NameAbbr
	}
	if source.Comment != nil {
		target.Comment = source.Comment
	}
	if source.ParentID != nil {
		target.ParentID = source.ParentID
	}
}
