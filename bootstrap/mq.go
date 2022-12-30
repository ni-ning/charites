package bootstrap

import (
	"charites/global"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func setupRocketMQ() (err error) {
	global.Producer, err = rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{global.RocketMQSetting.NameServer})),
		producer.WithRetry(2),
		producer.WithGroupName("order_srv"), // 生产者组
	)
	return
}
