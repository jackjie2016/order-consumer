package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/phper95/pkg/es"
	"gitee.com/phper95/pkg/mq"
	"gitee.com/phper95/pkg/strutil"
	"github.com/Shopify/sarama"
	"order-consumer/global"
)

var orderConsumer *mq.Consumer

func StartConsumer() {
	var err error
	orderConsumer, err = mq.StartKafkaConsumer(global.CONFIG.Kafka.Hosts,
		[]string{global.OrderTopic}, "order-consumer", nil, MsgHandler)
	if err != nil {
		panic(fmt.Sprintf("err %v,host %v", err, global.CONFIG.Kafka.Hosts))
	}
}
func MsgHandler(msg *sarama.ConsumerMessage) (bool, error) {
	mq.KafkaStdLogger.Printf("partion: %d ; offset : %d; msg : %s",
		msg.Partition, msg.Offset, string(msg.Value))
	orderMsg := OrderMsg{}
	err := json.Unmarshal(msg.Value, &orderMsg)
	if err != nil {
		//格式异常的数据，回到队列也不会解析成功
		global.LOG.Error("Unmarshal error", err, string(msg.Value))
		return true, nil
	}
	mq.KafkaStdLogger.Printf("product: %+v", orderMsg)
	orderIndex := orderMsg.OrderIndex
	orderIndex.OrderStatus = orderMsg.Status
	names := make([]string, 0)
	productIDs := make([]int64, 0)
	orderSuffixLen := 4

	if len(orderMsg.OrderId) >= orderSuffixLen {
		orderIndex.OrderIdSuffix = orderMsg.OrderId[len(orderMsg.OrderId)-4:]
	}

	for _, cart := range orderMsg.CartInfo {
		names = append(names, cart.ProductInfo.StoreName)
		productIDs = append(productIDs, cart.ProductInfo.Id)
	}
	orderIndex.Names = names
	orderIndex.ProductIds = productIDs

	esClient := es.GetClient(es.DefaultClient)
	switch orderMsg.Operation {
	case global.OperationCreate, global.OperationUpdate:
		esClient.BulkCreate(global.IndexName, orderIndex.OrderId,
			strutil.Int64ToString(orderIndex.Uid), orderIndex)

	case global.OperationDelete:
		err := esClient.Delete(context.Background(), global.IndexName, orderIndex.OrderId, strutil.Int64ToString(orderIndex.Uid))
		if err != nil {
			global.LOG.Error("DeleteRefresh error", err, "id", orderIndex.OrderId)
		}
	}

	return true, nil
}

func CloseConsumer() {
	if orderConsumer != nil {
		if err := orderConsumer.Close(); err != nil {
			global.LOG.Error("orderConsumer close error")
		}
	}
}
