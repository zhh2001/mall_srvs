package model

type Inventory struct {
	BaseModel
	Goods   int32 `gorm:"type:int;index"`
	Stocks  int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int"` // 分布式锁的乐观锁
}

type InventoryNew struct {
	BaseModel
	Goods   int32 `gorm:"type:int;index"`
	Stocks  int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int"` // 分布式锁的乐观锁
	Freeze  int32 `gorm:"type:int"` // 冻结库存
}

//type InventoryHistory struct {
//	user  int32
//	goods int32
//	nums  int32
//	order int32 // 1. 表示库存是预扣减，幂等性。 2. 表示已支付
//}

type Delivery struct {
	Goods   int32  `gorm:"type:int;index"`
	Nums    int32  `gorm:"type:int"`
	OrderSn string `gorm:"type:varchar(200)"`
	Status  string `gorm:"type:varchar(200)"` // 1-表示等待支付 2-表示支付成功 3-失败
}
