# kafka topic 探活组件
## 工作机制简述
- 生产者：起3个实例，每个实例每隔10s生产一次消息，如果生产消息异常则推送异常消息到如流群【kafka风险通知】
- 消费者：起3个实例，每隔2分钟commit一次，如果检测到rebalance或者消费异常则推送异常消息到如流群【kafka风险通知】
# 使用方法
## topic 操作命令示例
无需授权topic
```

1. 创建topic
/home/work/local/kafka_online/bin/kafka-topics.sh --zookeeper '' --create --topic prober-topic1 --partitions 3 --replication-factor 3
2. 删除topic
/home/work/local/kafka_online/bin/kafka-topics.sh --zookeeper '' --delete --topic prober-topic1
Topic prober-topic1 is marked for deletion.
Note: This will have no impact if delete.topic.enable is not set to true.

```

## 部署
1. 部署路径：127.0.0.1:/tmp/kafka_prober
2. 批量启动kafka-cluster.txt下所有的生产者和消费这
默认走sasl协议，licai的集群未开启鉴权，所以需要单独启动

```
sh opt_kafka_prober.sh -m produce -o start
sh opt_kafka_prober.sh -m comsume -o start
# 检查
sh opt_kafka_prober.sh -m produce -o check
sh opt_kafka_prober.sh -m consume -o check
# 启动未开启鉴权的集群的探活服务
sh opt_kafka_prober.sh -m produce -c licai_online_dd -f /tmp/kafka-prober-noauth.yaml -o start
sh opt_kafka_prober.sh -m produce -c licai_online_bl -f /tmp/kafka-prober-noauth.yaml -o start
```
