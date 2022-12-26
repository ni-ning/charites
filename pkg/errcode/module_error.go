package errcode

var (
	// 200xx tag
	ErrorGetTagListFail = NewError(20010001, "获取标签列表失败")

	// 300xx shopping
	ErrorNotFoundProduct = NewError(30010001, "没有找到商品记录")

	// 400xx stock
	ErrorDBOperateStock = NewError(40010001, "库存操作错误")
	ErrorNotFoundStock  = NewError(40010002, "没有找到库存记录")
	ErrorNotEnoughStock = NewError(40010003, "商品库存不足")
)
