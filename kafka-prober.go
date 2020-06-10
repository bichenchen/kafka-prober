package main

import (
	"flag"
	"fmt"
	"kafka-prober/config"
	"kafka-prober/consumer"
	"kafka-prober/global"
	"kafka-prober/logger"
	"kafka-prober/producer"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

func init() {
	sarama.Logger = log.New(os.Stdout, "[Sarama] ", log.LstdFlags)
}
func loadConfig() {
	global.ClusterName = flag.String("cluster", "", "The cluster name of kafka cluster")
	global.LogFileName = flag.String("log", "/tmp/kafka-prober-public-test.log", "the log file name of prober[use absolute path]")
	global.ConfigFileName = flag.String("config", "/tmp/kafka-prober.yaml", "the config file name of prober[use absolute path]")
	global.Brokers = flag.String("brokers", os.Getenv("KAFKA_PEERS"), "The Kafka brokers to connect to, as a comma separated list")
	global.UserName = flag.String("username", "test", "The SASL username")
	global.Passwd = flag.String("passwd", os.Getenv("KAFKA_PASSWORD"), "The SASL password")
	global.Topics = flag.String("topics", "prober-topic1", "The Kafka topic to use, as a comma separated list")
	global.ProduceInstanceNum = flag.Int("consumeNum", 3, "The instance number of single producer")
	global.ConsumeNum = flag.Int("produceNum", 3, "The number of consume")
	global.Mode = flag.String("mode", "produce", "Mode to run in: \"produce\" to produce, \"consume\" to consume")
	global.ProduceInvertal = flag.Int("interval", 10, "The interval of producer to create message")

	flag.Parse()
}

func main() {
	loadConfig()
	if *global.Brokers == "" {
		log.Fatalln("at least one broker is required")
	}
	splitBrokers := strings.Split(*global.Brokers, ",")
	splitTopics := strings.Split(*global.Topics, ",")
	if len(splitBrokers) == 0 || len(splitTopics) == 0 {
		log.Fatalln("Wrong format of brokers or topics")
	}

	if *global.UserName == "" {
		log.Fatalln("SASL username is required")
	}

	if *global.Passwd == "" {
		log.Fatalln("SASL password is required")
	}
	if *global.ClusterName == "" {
		log.Fatalln("cluster name is required")
	}
	if *global.ConfigFileName == "" {
		log.Fatalln("config file name is required")
	}
	err := config.ParseConfigFile(*global.ConfigFileName)
	if err != nil {
		log.Fatalf("parse config file error:%s", err)
	}
	consumerConf := config.InitConsumeConfig(*global.UserName, *global.Passwd, config.Vars)
	conf := config.InitProduceConfig(*global.UserName, *global.Passwd, config.Vars)
	logger.LogToFile(*global.LogFileName)

	if *global.Mode == "consume" {
		var groupID string
		for i := 1; i <= *global.ConsumeNum; i++ {
			groupID = fmt.Sprintf("prober-group-%d", i)
			go consumer.LoopConsumer(*global.ClusterName, splitBrokers, splitTopics, groupID, consumerConf)
		}

	} else if *global.Mode == "produce" {
		interval := time.Duration(*global.ProduceInvertal) * time.Second
		for i := 1; i <= *global.ProduceInstanceNum; i++ {
			for _, topic := range splitTopics {
				go producer.LoopProducer(*global.ClusterName, splitBrokers, topic, conf, interval)
			}
		}
	} else {
		log.Fatal("wrong mode, can only be consume or produce")
	}
	for {
		time.Sleep(time.Second)
	}
}
