DROP TABLE IF EXISTS `sys_org`;

CREATE TABLE
  `sys_org` (
    `id` BIGINT (20) NOT NULL COMMENT '组织ID',
    `code` VARCHAR(64) NOT NULL COMMENT '组织编码. 4位一级. 0001,00010001,000100010001,以此类推',
    `name` VARCHAR(64) NOT NULL COMMENT '组织名称',
    `name_abbr` VARCHAR(64) NOT NULL COMMENT '组织名称简称',
    `comment` VARCHAR(128) COMMENT '组织备注',
    `parent_id` BIGINT (20) DEFAULT 0 COMMENT '父级组织ID. 0表示没有父组织',
    `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `create_by` VARCHAR(32) NOT NULL COMMENT '创建人',
    `update_time` DATETIME ON UPDATE CURRENT_TIMESTAMP COMMENT '信息更新时间',
    `update_by` VARCHAR(32) COMMENT '信息更新人',
    PRIMARY KEY `pk_org_id` (`id`)
  );