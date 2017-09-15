package config

import "os"
import "strings"

// Config contains the configuration used bu the system
type Config struct {
	Etcd3Config etcd3Config
	Etcd2Config etcd2Config
	SOAConfig   soaConfig
}

type etcd3Config struct {
	Username  string
	Password  string
	CACert    string
	Endpoints []string
}

type etcd2Config struct {
	Endpoints []string
}

type soaConfig struct {
}

var instance *Config

// GetConfig gives you the config for DJ
func GetConfig() *Config {
	if instance == nil {
		conf := readEnv()
		instance = &conf
	}
	return instance
}

func readEnv() Config {
	conf := Config{}
	if username := os.Getenv("ITDNS_ETCD3_USERNAME"); username != "" {
		conf.Etcd3Config.Username = username
	}
	if password := os.Getenv("ITDNS_ETCD3_PASSWORD"); password != "" {
		conf.Etcd3Config.Password = password
	}
	if ca := os.Getenv("ITDNS_ETCD3_CA"); ca != "" {
		conf.Etcd3Config.CACert = ca
	}
	if endpoints := os.Getenv("ITDNS_ETCD3_ENDPOINTS"); endpoints != "" {
		conf.Etcd3Config.Endpoints = strings.Split(endpoints, ",")
	}

	if endpoints := os.Getenv("ITDNS_ETCD2_ENDPOINTS"); endpoints != "" {
		conf.Etcd2Config.Endpoints = strings.Split(endpoints, ",")
	}

	return conf
}
