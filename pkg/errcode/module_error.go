package errcode

var (
	// 200xx tag
	ErrorGetTagListFail = NewError(20010001, "获取标签列表失败")

	// 300xx shopping
	ErrorNotFoundProduct = NewError(30010001, "没有找到商品记录")

	// 400xx stock
	ErrorDBOperateStock = NewError(40010001, "数据库操作错误")
	ErrorNotFoundStock  = NewError(40010002, "没有找到库存记录")
	ErrorNotEnoughStock = NewError(40010003, "商品库存不足")
	ErrorNeedRetryStock = NewError(40010004, "商品库存已更新, 请重试")

	ErrorRedisLockStock   = NewError(40010005, "RedisLock错误")
	ErrorRedisUnlockStock = NewError(40010006, "RedisUnlock错误")

	// 500xx order
	ErrorRPCOrderToGoods  = NewError(50010001, "订单请求商品微服务错误")
	ErrorRPCOrderToStock  = NewError(50010002, "订单请求库存微服务错误")
	ErrorCreateOrder      = NewError(50010003, "订单创建错误")
	ErrorCreateOrderDetal = NewError(50010004, "订单详情创建错误")

	ErrorNewTransactionProducer = NewError(50010005, "NewTransactionProducer Error")
	ErrorOrderEntityParam       = NewError(50010006, "OrderEntityParam Error")
	ErrorOrderCreate            = NewError(50010007, "创建订单失败")
)
