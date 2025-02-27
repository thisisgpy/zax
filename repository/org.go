package repository

import (
	"zax/model"

	"github.com/jmoiron/sqlx"
)

type OrgRepository struct {
	db *sqlx.DB
}

func NewOrgRepository(db *sqlx.DB) *OrgRepository {
	return &OrgRepository{db: db}
}

func (repo *OrgRepository) Insert(tx *sqlx.Tx, org *model.SysOrg) error {
	sql := `
		INSERT INTO sys_org 
		(
			id, code, name, name_abbr, comment, parent_id, is_deleted, create_time, create_by, update_time, update_by
		) 
		VALUES 
		(
			:id, :code, :name, :name_abbr, :comment, :parent_id, :is_deleted, :create_time, :create_by, :update_time, :update_by
		)
	`
	_, err := tx.NamedExec(sql, org)
	return err
}

func (repo *OrgRepository) BatchInsert(tx *sqlx.Tx, orgs []*model.SysOrg) error {
	sql := `
		INSERT INTO sys_org 
		(
			id, code, name, name_abbr, comment, parent_id, is_deleted, create_time, create_by, update_time, update_by
		) 
		VALUES 
		(
			:id, :code, :name, :name_abbr, :comment, :parent_id, :is_deleted, :create_time, :create_by, :update_time, :update_by
		)
	`
	_, err := tx.NamedExec(sql, orgs)
	return err
}

func (repo *OrgRepository) UpdateSelective(tx *sqlx.Tx, org *model.SysOrg) error {
	sql := `
        UPDATE sys_org
        SET update_time = :update_time
        {{if .Code}}        , code = :code{{end}}
        {{if .Name}}        , name = :name{{end}}
        {{if .NameAbbr}}    , name_abbr = :name_abbr{{end}}
        {{if .Comment}}     , comment = :comment{{end}}
        {{if .ParentId}}    , parent_id = :parent_id{{end}}
        {{if .IsDeleted}}   , is_deleted = :is_deleted{{end}}
        {{if .UpdateBy}}    , update_by = :update_by{{end}}
				{{if .UpdateTime}}  , update_time = :update_time{{end}}
        WHERE id = :id
    `
	_, err := tx.NamedExec(sql, org)
	return err
}

func (repo *OrgRepository) SelectById(id int64) (*model.SysOrg, error) {
	sql := `
		SELECT 
			id, code, name, name_abbr, comment, parent_id, is_deleted, create_time, create_by, update_time, update_by 
		FROM 
			sys_org 
		WHERE 
			id = :id AND is_deleted = false
	`
	var org model.SysOrg
	err := repo.db.Get(&org, sql, map[string]interface{}{"id": id})
	return &org, err
}
