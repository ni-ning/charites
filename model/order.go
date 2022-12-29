package model

import "time"

// Order 订单表
type Order struct {
	BaseModel // 嵌入默认7个字段

	UserId         int64
	OrderId        int64
	TradeId        string
	PayChannel     int8
	Status         int64
	PayAmount      int64
	PayTime        *time.Time
	ReceiveAddress string
	ReceiveName    string
	ReceivePhone   string
}

func (Order) TableName() string {
	return "shopping_order"
}

// OrderDetail 订单详情表
type OrderDetail struct {
	BaseModel // 嵌入默认7个字段

	UserId  int64
	OrderId int64
	GoodsId int64

	Title       string
	MarketPrice int64
	Price       int64
	Brief       string
	HeadImgs    string
	Videos      string
	Detail      string

	Num       int64
	PayAmount int64
	PayTime   *time.Time
}

func (OrderDetail) TableName() string {
	return "shopping_order_detail"
}

type OrderGoodsStockInfo struct {
	OrderId int64
	GoodsId int64
	Num     int64
}
