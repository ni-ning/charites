create database if not exists charites charset utf8mb4 collate utf8mb4_general_ci;

create table charites.shopping_room_goods(
    id bigint(20) unsigned not null primary key auto_increment comment '主键',
    created_at datetime not null default now() comment '创建时间',
    created_by varchar(64) default '' comment '创建人',
    updated_at datetime not null default now() comment '修改时间',
    updated_by varchar(64) default '' comment '修改人',
    `version` smallint(5) unsigned not null default 0 comment '乐观锁版本号',
    is_deleted tinyint(3) default 0 comment '是否删除 0未删除 1已删除',
    
    goods_id bigint(20) unsigned not null default 0 comment '商品ID',
    room_id bigint(20) unsigned not null default 0 comment '直播间ID',
    weight bigint(20) not null default 1000 comment '排序权重',
    is_current tinyint(4) default 0 comment '是否当前讲解中 0否 1是',

    unique udx_goods_room_id(goods_id, room_id)

)engine=innodb charset=utf8mb4 comment '直播间商品表';