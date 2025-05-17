package model

type Inventory struct {
	BaseModel
	Goods   int32 `gorm:"type:int;index"`
	Stocks  int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int"` // 分布式锁的乐观锁
}

//type InventoryHistory struct {
//	user  int32
//	goods int32
//	nums  int32
//	order int32 // 1. 表示库存是预扣减，幂等性。 2. 表示已支付
//}
