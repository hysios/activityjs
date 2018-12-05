package model

var (
	disableCompute    = false
	upsidePropagation = false
)

// Time 自定义的时间结构
type Time int64

// ActivityItem 订单项
type ActivityItem struct {
	ID          *Identity
	ItemID      *Identity
	Product     Product
	Title       string
	Price       float64
	ProduceDate Time
}

// Product 产品
type Product struct {
	ID *Identity
	// Name         string
	Title        string
	Avatar       string
	DefaultPrice float64
}

// Gift 礼品
type Gift struct {
	OrderItem
}

func DisabledCompute(fn func()) {
	disableCompute = true
	fn()
	disableCompute = false
}

func UpsidePropagation(fn func()) {
	upsidePropagation = true
	fn()
	upsidePropagation = false
}
