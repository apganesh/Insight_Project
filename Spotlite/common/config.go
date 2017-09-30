package insight

type config struct {
	KafkaBrokers      string   `json:"kafka_brokers"`
	CassandraClusters []string `json:"cassandra_clusters"`
	RedisClusters     string   `json:"redis_clusters"`
	KafkaTopics       string   `json:"kafka_topics"`
	KafkaGroup        string   `json:"kafka_groups"`
}

var DefaultConfig = config{
	KafkaBrokers:      "ec2-52-11-70-58.us-west-2.compute.amazonaws.com",
	CassandraClusters: []string{"127.0.0.1"},
	RedisClusters:     "127.0.0.1:6379",
	KafkaTopics:       "jaja",
	KafkaGroup:        "0",
}
