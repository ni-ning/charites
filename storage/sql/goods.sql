create database if not exists charites charset utf8mb4 collate utf8mb4_general_ci;

create table charites.shopping_goods(
    id bigint(20) unsigned not null primary key auto_increment comment '主键',
    created_at datetime not null default now() comment '创建时间',
    created_by varchar(64) default '' comment '创建人',
    updated_at datetime not null default now() comment '修改时间',
    updated_by varchar(64) default '' comment '修改人',
    `version` smallint(5) unsigned not null default 0 comment '乐观锁版本号',
    is_deleted tinyint(3) default 0 comment '是否删除 0未删除 1已删除',
    
    goods_id bigint(20) unsigned not null default 0 comment '商品ID',
    category_id bigint(20) unsigned not null default 0 comment '类目ID',
    brand_name varchar(255) not null comment '品牌名',
    code varchar(64) not null comment '码',
    status tinyint(4) unsigned not null default 0 comment '是否上架 0上架 1下架',
    title varchar(255) not null comment '名称',
    market_price bigint(20) unsigned not null default 0 comment '市场价/划线价(分)',
    price bigint(20) unsigned not null default 0 comment '售价(分)',
    
    breif varchar(255) not null default '' comment '简介',
    head_imgs varchar(1024) not null default '' comment '头像 []',
    videos varchar(1024) not null default '' comment '视频介绍 []',
    detail varchar(2048) not null default '' comment '简介 [实际也是图片集合]',
    ext_json varchar(2048) not null default '' comment '扩展字段 {}',
    
    unique udx_goods_goods_id(goods_id),
    index idx_goods_category_id(category_id)

)engine=innodb charset=utf8mb4 comment '商品表';