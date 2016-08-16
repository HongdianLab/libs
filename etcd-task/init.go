package task

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/coreos/go-etcd/etcd"
)

var (
	logger           *log.Logger
	clientSingleton  *etcd.Client
	clientSingletonO = &sync.Once{}
)

func Client() *etcd.Client {
	clientSingletonO.Do(func() {
		hosts := []string{"http://localhost:4001"}
		if len(os.Getenv("ETCD_HOSTS")) != 0 {
			strhosts := os.Getenv("ETCD_HOSTS")
			hosts = strings.Split(strhosts, ",")
		}

		cacert := os.Getenv("ETCD_CACERT")
		tlskey := os.Getenv("ETCD_TLS_KEY")
		tlscert := os.Getenv("ETCD_TLS_CERT")

		ishttps := false
		if len(cacert) != 0 && len(tlskey) != 0 && len(tlscert) != 0 {
			ishttps = true
		}

		etcdhosts := make([]string, len(hosts))
		for _, host := range hosts {
			if ishttps {
				if !strings.Contains(host, "https://") {
					host = strings.Replace(host, "http", "https", 1)
				}
			}
			etcdhosts = append(etcdhosts, host)
		}
		if ishttps {
			c, err := etcd.NewTLSClient(etcdhosts, tlscert, tlskey, cacert)
			if err != nil {
				panic(err)
			}
			clientSingleton = c
		} else {
			clientSingleton = etcd.NewClient(etcdhosts)
		}
	})
	return clientSingleton
}

func init() {
	logger = log.New(os.Stderr, "[etcd-task]", log.LstdFlags)
}
