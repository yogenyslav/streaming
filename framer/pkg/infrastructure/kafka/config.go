package kafka

type Config struct {
	Brokers []string      `yaml:"brokers"`
	Topics  []TopicConfig `yaml:"topics"`
}

type TopicConfig struct {
	Name       string `yaml:"name"`
	Partitions int32  `yaml:"partitions"`
}
