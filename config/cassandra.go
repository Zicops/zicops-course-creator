package config

import "os"

// struct for cassandra config
type Cassandra struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Keyspace string `yaml:"keyspace"`
}

// initialize cassandra config struct using env variables
func NewCassandraConfig() *Cassandra {
	return &Cassandra{
		Host:     getEnv("CASSANDRA_HOST", "127.0.0.1"),
		Port:     getEnv("CASSANDRA_PORT", "9042"),
		Username: getEnv("CASSANDRA_USERNAME", "cassandra"),
		Password: getEnv("CASSANDRA_PASSWORD", "cassandra"),
		Keyspace: getEnv("CASSANDRA_KEYSPACE", "coursez"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
