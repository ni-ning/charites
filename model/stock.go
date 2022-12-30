package model

// Stock 库存表
type Stock struct {
	BaseModel // 嵌入默认7个字段

	GoodsId int64
	Num     int64
	Lock    int64
}

func (Stock) TableName() string {
	return "shopping_stock"
}

// StockRecord 库存记录表
type StockRecord struct {
	BaseModel // 嵌入默认7个字段

	OrderId int64
	GoodsId int64
	Num     int64
	Status  int64 // 1预扣减 2已回滚 3扣减
}

/*
1预扣减 创建订单和订单详情之前 扣减库存动作新增该记录
2已回滚
	(1)创建订单和订单详情失败，发送回滚消息 更新为该状态
	(2)超时未及时支付，发送回滚消息 更新为该状态
3扣减 成功支付 更新为该状态
*/

func (StockRecord) TableName() string {
	return "shopping_stock_record"
}
