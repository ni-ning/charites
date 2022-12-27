create table charites.shopping_stock(
    id bigint(20) unsigned not null primary key auto_increment comment '主键',
    created_at datetime not null default now() comment '创建时间',
    created_by varchar(64) default '' comment '创建人',
    updated_at datetime not null default now() comment '修改时间',
    updated_by varchar(64) default '' comment '修改人',
    `version` smallint(5) unsigned not null default 0 comment '乐观锁版本号',
    is_deleted tinyint(3) default 0 comment '是否删除 0未删除 1已删除',
    
    goods_id bigint(20) unsigned not null default 0 comment '商品ID',
    num bigint(20) unsigned not null default 0 comment '库存数量',

    unique udx_stock_goods_id(goods_id)

)engine=innodb charset=utf8mb4 comment '库存表';