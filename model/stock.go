package model

// Stock 库存表
type Stock struct {
	BaseModel // 嵌入默认7个字段

	GoodsId int64
	Num     int64
}

func (Stock) TableName() string {
	return "shopping_stock"
}
