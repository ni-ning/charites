create table charites.shopping_order(
    id bigint(20) unsigned not null primary key auto_increment comment '主键',
    created_at datetime not null default now() comment '创建时间',
    created_by varchar(64) default '' comment '创建人',
    updated_at datetime not null default now() comment '修改时间',
    updated_by varchar(64) default '' comment '修改人',
    `version` smallint(5) unsigned not null default 0 comment '乐观锁版本号',
    is_deleted tinyint(3) default 0 comment '是否删除 0未删除 1已删除',
    
    user_id bigint(20) unsigned not null default 0 comment '用户ID',
    order_id bigint(20) unsigned not null default 0 comment '订单ID',
    trade_id varchar(128)  not null default '' comment '交易单号',

    pay_channel tinyint(4) unsigned not null default 0 comment '支付方式',
    `status` int unsigned not null default 0 comment '订单状态:100创建订单/待支付 200已支付 300已关闭 400完成',

    pay_amount bigint(20) unsigned not null default 0 comment '支付金额(分)',
    pay_time datetime  comment '支付时间',

    receive_address varchar(128) not null default '' comment '收货地址',
    receive_name varchar(128) not null default '' comment '收货人',
    receive_phone varchar(11) not null default '' comment '收货电话',

    index(user_id),
    index(order_id)

)engine=innodb charset=utf8mb4 comment '订单表';

create table charites.shopping_order_detail(
    id bigint(20) unsigned not null primary key auto_increment comment '主键',
    created_at datetime not null default now() comment '创建时间',
    created_by varchar(64) default '' comment '创建人',
    updated_at datetime not null default now() comment '修改时间',
    updated_by varchar(64) default '' comment '修改人',
    `version` smallint(5) unsigned not null default 0 comment '乐观锁版本号',
    is_deleted tinyint(3) default 0 comment '是否删除 0未删除 1已删除',
    
    user_id bigint(20) unsigned not null default 0 comment '用户ID',
    order_id bigint(20) unsigned not null default 0 comment '订单ID',
    goods_id bigint(20) unsigned not null default 0 comment '商品ID',
    
    title varchar(255) not null comment '名称',
    market_price bigint(20) unsigned not null default 0 comment '市场价/划线价(分)',
    price bigint(20) unsigned not null default 0 comment '售价(分)',
    breif varchar(255) not null default '' comment '简介',
    head_imgs varchar(1024) not null default '' comment '头像 []',
    videos varchar(1024) not null default '' comment '视频介绍 []',
    detail varchar(2048) not null default '' comment '简介 [实际也是图片集合]',

    num bigint(20) unsigned not null default 0 comment '商品数量',
    pay_amount bigint(20) unsigned not null default 0 comment '支付金额(分)',
    pay_time datetime  comment '支付时间',

    index(user_id),
    index(order_id),
    index(goods_id)

)engine=innodb charset=utf8mb4 comment '订单详情表';