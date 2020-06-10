package consumer

import (
	"fmt"
	"kafka-prober/logger"
	"kafka-prober/util"
	"os"
	"os/signal"
	"strings"

	cluster "github.com/bsm/sarama-cluster"
)

func LoopConsumer(clusterName string, splitBrokers []string, topicList []string, group string, conf *cluster.Config) {
	//conver topic from string to list
	brokers := strings.Join(splitBrokers, ",")
	topics := strings.Join(topicList, ",")
	var msg string
	// topicList := make([]string, 1)
	// topicList[0] = topic
	// config := conf
	// config.Consumer.Return.Errors = true
	// config.Group.Return.Notifications = true
	consumer, err := cluster.NewConsumer(splitBrokers, group, topicList, conf)
	if err != nil {
		logger.Logger.Fatalf("consumer created error:%s groupID=%s", err, group)
	}
	logger.Logger.Infof("consumer created groupID=%s topic=%s ", group, topics)
	defer func() {
		if err := consumer.Close(); err != nil {
			logger.Logger.Warnf("consumer closed err:%s groupID=%s", err, group)
		}
	}()
	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	// consumer errors
	go func() {
		for err := range consumer.Errors() {
			msg = fmt.Sprintf("cosume Error:%s groupID=%s", err.Error(), group)
			util.SendAlertInfoToHi(clusterName, brokers, "consume", topics, group, msg)
			logger.Logger.Warn(msg)
		}
	}()
	// consumer notifications
	go func() {
		for ntf := range consumer.Notifications() {
			msg = fmt.Sprintf("Rebalanced:%+v groupID=%s", ntf, group)
			util.SendAlertInfoToHi(clusterName, brokers, "consume", topics, group, msg)
			logger.Logger.Warn(msg)
		}
	}()
	consumed := 0
ConsumerLoop:
	for {
		select {
		case info, ok := <-consumer.Messages():
			if ok {

				logger.Logger.Infof("Consumed message groupID=%s topic=%s partition=%d offset=%d key=%s value=%s",
					group, info.Topic, info.Partition, info.Offset, info.Key, info.Value)
				consumed++
			}
		case <-signals:
			break ConsumerLoop
		}
	}
	logger.Logger.Infof("Consumed time:%d groupID=%s", consumed, group)
}
