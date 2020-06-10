package config

import (
	"time"

	"io/ioutil"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"

	yaml "gopkg.in/yaml.v2"
)

type Conf struct {
	MetadataFull             bool          `yaml:"metadataFull"`
	ClientID                 string        `yaml:"clientID"`
	NetSaslEnable            bool          `yaml:"netSaslEnable"`
	NetSaslHandshake         bool          `yaml:"netSaslHandshake"`
	ProducerReturnSuccesses  bool          `yaml:"producerReturnSuccesses"`
	ProducerRetryMax         int           `yaml:"producerRetryMax"`
	RequireAcks              int           `yaml:"requireAcks"`
	GroupOffsetsRetryMax     int           `yaml:"groupOffsetsRetryMax"`
	GroupSessionTimeout      time.Duration `yaml:"groupSessionTimeout"`
	GroupHeartbeatInterval   time.Duration `yaml:"groupHeartbeatInterval"`
	ConsumerReturnErrors     bool          `yaml:"consumerReturnErrors"`
	GroupReturnNotifications bool          `yaml:"groupReturnNotifications"`
}

var Vars Conf

func ParseConfigFile(configFileName string) error {
	content, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(content, &Vars)
	if err != nil {
		return err
	}
	return nil
}

// InitConfig returns a new configuration instance with sane defaults.
func InitConsumeConfig(UserName string, Passwd string, config Conf) *cluster.Config {
	c := &cluster.Config{
		Config: *sarama.NewConfig(),
	}
	c.Metadata.Full = config.MetadataFull
	c.Version = sarama.V0_10_0_0
	c.ClientID = config.ClientID
	c.Net.SASL.Enable = config.NetSaslEnable
	c.Net.SASL.User = UserName
	c.Net.SASL.Password = Passwd
	c.Net.SASL.Handshake = config.NetSaslHandshake
	c.Group.PartitionStrategy = cluster.StrategyRange
	c.Group.Offsets.Retry.Max = config.GroupOffsetsRetryMax
	c.Group.Offsets.Synchronization.DwellTime = c.Consumer.MaxProcessingTime
	c.Group.Session.Timeout = config.GroupSessionTimeout * time.Second
	c.Group.Heartbeat.Interval = config.GroupHeartbeatInterval * time.Second
	c.Consumer.Return.Errors = config.ConsumerReturnErrors
	c.Group.Return.Notifications = config.GroupReturnNotifications
	return c
}
func InitProduceConfig(UserName string, Passwd string, config Conf) *sarama.Config {
	c := sarama.NewConfig()
	c.Metadata.Full = config.MetadataFull
	c.Version = sarama.V0_10_0_0
	c.ClientID = config.ClientID
	c.Net.SASL.Enable = config.NetSaslEnable
	c.Net.SASL.User = UserName
	c.Net.SASL.Password = Passwd
	c.Net.SASL.Handshake = config.NetSaslHandshake
	c.Producer.Return.Successes = config.ProducerReturnSuccesses
	c.Producer.Retry.Max = config.ProducerRetryMax
	c.Producer.RequiredAcks = sarama.WaitForAll
	return c
}

// consumerConf := cluster.NewConfig()
