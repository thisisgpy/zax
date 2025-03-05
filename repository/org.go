package repository

import (
	"fmt"
	"strings"
	"zax/model"

	"github.com/jmoiron/sqlx"
)

type OrgRepository struct {
	db *sqlx.DB
}

func NewOrgRepository(db *sqlx.DB) *OrgRepository {
	return &OrgRepository{db: db}
}

func (repo *OrgRepository) Insert(tx *sqlx.Tx, orgs []*model.SysOrg) error {
	sql := `
		INSERT INTO sys_org 
		(
			id, code, name, name_abbr, comment, parent_id, create_time, create_by
		) 
		VALUES 
		(
			:id, :code, :name, :name_abbr, :comment, :parent_id, :create_time, :create_by
		)
	`
	_, err := tx.NamedExec(sql, orgs)
	return err
}

func (repo *OrgRepository) UpdateSelective(tx *sqlx.Tx, org *model.SysOrg) error {
	var fields []string

	if org.Code != "" {
		fields = append(fields, "code = :code")
	}
	if org.Name != "" {
		fields = append(fields, "name = :name")
	}
	if org.NameAbbr != nil && *org.NameAbbr != "" {
		fields = append(fields, "name_abbr = :name_abbr")
	}
	if org.Comment != nil && *org.Comment != "" {
		fields = append(fields, "comment = :comment")
	}
	if org.ParentID != 0 {
		fields = append(fields, "parent_id = :parent_id")
	}
	if org.UpdateTime != nil {
		fields = append(fields, "update_time = :update_time")
	}
	if org.UpdateBy != nil && *org.UpdateBy != "" {
		fields = append(fields, "update_by = :update_by")
	}

	if len(fields) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE sys_org SET %s WHERE id = :id", strings.Join(fields, ", "))
	_, err := tx.NamedExec(query, org)
	return err
}

func (repo *OrgRepository) SelectById(id int64) (*model.SysOrg, error) {
	sql := `
		SELECT 
			id, code, name, name_abbr, comment, parent_id, create_time, create_by, update_time, update_by 
		FROM 
			sys_org 
		WHERE 
			id = ?
	`
	var org model.SysOrg
	err := repo.db.Get(&org, sql, id)
	return &org, err
}

func (repo *OrgRepository) SelectByParentID(parentID int64) ([]*model.SysOrg, error) {
	sql := `
		SELECT 
			id, code, name, name_abbr, comment, parent_id, create_time, create_by, update_time, update_by 
		FROM 
			sys_org 
		WHERE 
			parent_id = ?
		ORDER BY
			code ASC
	`
	var orgs []*model.SysOrg
	err := repo.db.Select(&orgs, sql, parentID)
	return orgs, err
}

func (repo *OrgRepository) SelectByCode(code string) (*model.SysOrg, error) {
	sql := `
		SELECT 
			id, code, name, name_abbr, comment, parent_id, create_time, create_by, update_time, update_by 
		FROM 
			sys_org 
		WHERE 
			code = ?
	`
	var org model.SysOrg
	err := repo.db.Get(&org, sql, code)
	return &org, err
}

func (repo *OrgRepository) SelectMaxCode(parentID int64) (string, error) {
	sql := `
		SELECT 
			MAX(code) 
		FROM 
			sys_org 
		WHERE 
			parent_id = ?
	`
	var code string
	err := repo.db.Get(&code, sql, parentID)
	return code, err
}

func (repo *OrgRepository) SelectDescendants(orgCode string) ([]*model.SysOrg, error) {
	sql := `
		SELECT 
			id, code, name, name_abbr, comment, parent_id, create_time, create_by, update_time, update_by 
		FROM 
			sys_org 
		WHERE 
			code LIKE CONCAT(?, '%') AND code != ?
		ORDER BY
			create_time ASC
	`
	var orgs []*model.SysOrg
	err := repo.db.Select(&orgs, sql, orgCode, orgCode)
	return orgs, err
}
