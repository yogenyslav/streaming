package minios3

type Config struct {
	Host      string         `yaml:"host"`
	AccessKey string         `yaml:"accessKey"`
	SecretKey string         `yaml:"secretKey"`
	Ssl       bool           `yaml:"ssl"`
	Buckets   []BucketConfig `yaml:"buckets"`
}

type BucketConfig struct {
	Name   string `yaml:"name"`
	Region string `yaml:"region"`
	Lock   bool   `yaml:"lock"`
}
