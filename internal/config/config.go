package config

type Config struct {
	Port         string `yaml:"address" mapstructure:"address"`
	URL          string `yaml:"url" mapstructure:"url"`
	Endpoint     string `yaml:"endpoint" mapstructure:"endpoint"`
	Rate         int    `yaml:"rate" mapstructure:"rate"`
	BufferSize   int    `yaml:"bufferSize" mapstructure:"bufferSize"`
	OTLPEndpoint string `yaml:"otlpEndpoint" mapstructure:"otlpEndpoint"`
	ServiceName  string `yaml:"serviceName" mapstructure:"serviceName"`
	UpstreamURL  string `yaml:"upstreamUrl" mapstructure:"upstreamUrl"`
}

