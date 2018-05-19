package cmd

import (
	"os"
	"strings"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
)

func newClientFromEnv() (*etcd.Client, error) {
	config := etcd.Config{
		Endpoints: []string{"http://localhost:2379"},
	}

	if os.Getenv("ETCDCTL_ENDPOINTS") != "" {
		config.Endpoints = strings.Split(os.Getenv("ETCDCTL_ENDPOINTS"), ",")
	}

	if os.Getenv("ETCDCTL_CACERT") != "" {
		// TLS is set
		tlsInfo := transport.TLSInfo{
			TrustedCAFile: os.Getenv("ETCDCTL_CACERT"),
			CertFile:      os.Getenv("ETCDCTL_CERT"),
			KeyFile:       os.Getenv("ETCDCTL_KEY"),
		}
		tlsConfig, err := tlsInfo.ClientConfig()
		if err != nil {
			return nil, err
		}
		config.TLS = tlsConfig
	}

	if os.Getenv("ETCDCTL_USER") != "" {
		parts := strings.Split(os.Getenv("ETCDCTL_USER"), ":")
		config.Username = parts[0]
		if len(parts) > 1 {
			config.Password = strings.Join(parts[1:], ":")
		}
	}

	return etcd.New(config)
}
