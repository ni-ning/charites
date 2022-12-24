package model

// Goods 商品表
type Goods struct {
	BaseModel // 嵌入默认7个字段

	GoodsId     int64
	CategoryId  int64
	BrandName   string
	Code        string
	Status      int8 // 小范围用int8
	Title       string
	MarketPrice int64
	Price       int64
	Brief       string
	HeadImgs    string
	Videos      string
	Detail      string
	ExtJson     string
}

func (Goods) TableName() string {
	return "shopping_goods"
}

// RoomGoods 直播间商品表
type RoomGoods struct {
	BaseModel // 嵌入默认7个字段

	RoomId    int64
	GoodsId   int64
	Weight    int64
	IsCurrent int8
}

func (RoomGoods) TableName() string {
	return "shopping_room_goods"
}
