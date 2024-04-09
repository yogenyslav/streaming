package s3

type Config struct {
	Host      string `yaml:"host"`
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
	Ssl       bool   `yaml:"ssl"`
}
