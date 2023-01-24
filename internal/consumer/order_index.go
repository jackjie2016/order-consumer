package consumer

import "time"

type OrderIndex struct {
	OrderId       string    `json:"order_id"`
	OrderIdSuffix string    `json:"order_id_suffix"`
	Names         []string  `json:"names"`
	ProductIds    []int64   `json:"product_ids"`
	Uid           int64     `json:"uid"`
	PayTime       time.Time `json:"pay_time"`
	PayType       string    `json:"pay_type"`
	RefundStatus  int       `json:"refund_status"` // '0 未退款 1 申请中 2 已退款'
	ShippingType  int       `json:"shipping_type"` //配送方式
	OrderStatus   int       `json:"order_status"`
	CreateTime    time.Time `json:"create_time"`
	UpdateTime    time.Time `json:"update_time"`
}

type OrderMsg struct {
	Operation string `json:"operation"`
	Status    int    `json:"_status"` //1-未付款 2-未发货 3-退款中 4-待收货 5-待评价 6-已完成 7-已退款
	CartInfo  []Cart `json:"cart_info"`
	OrderIndex
}

type Cart struct {
	Id                int64   `json:"id"`
	Uid               int64   `json:"uid"`
	Type              string  `json:"type"`
	ProductId         int64   `json:"product_id"`
	ProductAttrUnique string  `json:"product_attr_unique"`
	CartNum           int     `json:"cart_num"`
	CombinationId     int64   `json:"combination_id"`
	SeckillId         int64   `json:"seckill_id"`
	BargainId         int64   `json:"bargain_id"`
	CostPrice         float64 `json:"cost_price"`
	ProductInfo       Product `json:"product_info"`
	TruePrice         float64 `json:"true_price"`
	TrueStock         int     `json:"true_stock"`
	VipTruePrice      float64 `json:"vip_true_price"`
	Unique            string  `json:"unique"`
	IsReply           int     `json:"is_reply"`
}

type Product struct {
	Id        int64   `json:"id"`
	StoreName string  `json:"store_name"`
	CateId    int     `json:"cate_id"`
	Price     float64 `json:"price"`
	VipPrice  float64 `json:"vip_price"`
	OtPrice   float64 `json:"ot_price"`
}
