package producer

import (
	"fmt"
	"kafka-prober/logger"
	"kafka-prober/util"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

func LoopProducer(clusterName string, splitBrokers []string, topic string, conf *sarama.Config, interval time.Duration) {
	var msg string
	brokers := strings.Join(splitBrokers, ",")
	// 模拟长连接每间隔10s生产一次消息
	syncProducer, err := sarama.NewSyncProducer(splitBrokers, conf)
	if err != nil {

		logger.Logger.Fatalf("failed to create producer:%s brokers=%s topic=%s",
			err, splitBrokers, topic)
	}
	for {
		nowFormat := time.Now().Format("2006-01-02 15:04:05")
		mesg := fmt.Sprintf("%s:%s", topic, nowFormat)

		partition, offset, err := syncProducer.SendMessage(&sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(mesg),
		})
		if err != nil {
			msg = fmt.Sprintf("failed to send message to topic,err:%s topic=%s", err, topic)
			logger.Logger.Warn(msg)
			util.SendAlertInfoToHi(clusterName, brokers, "produce", topic, "", msg)
		}
		logger.Logger.Infof("wrote message ok topic=%s partition=%d offset=%d",
			topic, partition, offset)
		time.Sleep(interval)
		// time.Sleep(10 * time.Second)
		// _ = syncProducer.Close()
		// logger.Logger.Println("Bye now !")
	}
	logger.Logger.Warnf("produce Bye now ! topic=%s", topic)
}
