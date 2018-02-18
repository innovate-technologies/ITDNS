package etcd3

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/coreos/etcd/pkg/transport"
	"github.com/innovate-technologies/ITDNS/cache"
	"github.com/innovate-technologies/ITDNS/config"
	"golang.org/x/net/context"
)

type etcdRecord struct {
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
}

type etcdMxRecord struct {
	Value []interface{} `json:"value"`
	TTL   int           `json:"ttl"`
}

// Client is the struct allowing to interact with the etcd3 storage backend
type Client struct {
	etcdAPI *etcd.Client
	cache   *cache.Cache
}

var cfg = config.GetConfig()

// New gives a new etcd3 client
func New(cache *cache.Cache) Client {
	returnObject := Client{}
	returnObject.cache = cache

	var etcdConfig = etcd.Config{
		Endpoints: cfg.Etcd3Config.Endpoints,
	}

	if cfg.Etcd3Config.Username != "" {
		etcdConfig.Username = cfg.Etcd3Config.Username
		etcdConfig.Password = cfg.Etcd3Config.Password
	}

	if cfg.Etcd3Config.CACert != "" {
		tlsInfo := transport.TLSInfo{
			TrustedCAFile: cfg.Etcd3Config.CACert,
		}
		tlsConfig, err := tlsInfo.ClientConfig()
		if err != nil {
			log.Fatal(err)
		}
		etcdConfig.TLS = tlsConfig
	}
	c, err := etcd.New(etcdConfig)
	if err != nil {
		panic(err)
	}
	returnObject.etcdAPI = c

	return returnObject
}

// CreateCache fetches all records and enters them into the cache
func (c *Client) CreateCache() {
	res, err := c.etcdAPI.Get(context.Background(), "/DNS/", etcd.WithPrefix())
	if err != nil {
		panic(err)
	}

	recordsToCache := map[string][]cache.Record{}

	for _, kv := range res.Kvs {
		domainParts := strings.Split(string(kv.Key), "/")
		domainName := domainParts[2]
		if _, ok := recordsToCache[domainName]; !ok {
			recordsToCache[domainName] = []cache.Record{}
		}
		recordsToCache[domainName] = append(recordsToCache[domainName], parseRecords(kv)...)
	}

	for domainName := range recordsToCache {
		c.cache.Add(domainName, recordsToCache[domainName])
	}
}

// Watch watches for new records to add them to the cache
func (c *Client) Watch() {
	chans := c.etcdAPI.Watch(context.Background(), "/DNS/", etcd.WithPrefix())
	for resp := range chans {
		if resp.Canceled {
			break // error found
		}
		for _, ev := range resp.Events {
			if ev.IsCreate() || ev.IsModify() {
				c.addToCache(ev.Kv, 0)
			}
		}
	}
	c.Watch()
}

func (c *Client) addToCache(record *mvccpb.KeyValue, retry int) {
	pathParts := strings.Split(string(record.Key), "/")
	res, err := c.etcdAPI.Get(context.Background(), strings.Join(pathParts[:3], "/")+"/", etcd.WithPrefix())
	if err != nil {
		time.Sleep(time.Second)
		retry++
		if retry > 1000 {
			return
		}
		c.addToCache(record, retry)
		return
	}

	newRecords := []cache.Record{}

	for _, kv := range res.Kvs {
		newRecords = append(newRecords, parseRecords(kv)...)
	}

	c.cache.Add(pathParts[2], newRecords)
}

func parseRecords(record *mvccpb.KeyValue) []cache.Record {
	pathParts := strings.Split(string(record.Key), "/")

	var etcdRecords []etcdRecord

	if pathParts[len(pathParts)-1] == "MX" {
		etcdRecords = []etcdRecord{}
		mxRecords := []etcdMxRecord{}
		json.Unmarshal(record.Value, &mxRecords)
		for _, record := range mxRecords {
			etcdRecords = append(etcdRecords, etcdRecord{
				TTL:   record.TTL,
				Value: fmt.Sprintf("%.0f", record.Value[0].(float64)) + " " + record.Value[1].(string),
			})
		}
	} else {
		etcdRecords = []etcdRecord{}
		json.Unmarshal(record.Value, &etcdRecords)
	}

	recordsToCache := []cache.Record{}

	for _, etcdRecord := range etcdRecords {
		cacheRecord := cache.NewRecord()
		cacheRecord.Qname = pathParts[len(pathParts)-2]
		cacheRecord.Qtype = pathParts[len(pathParts)-1]
		cacheRecord.Content = etcdRecord.Value
		cacheRecord.TTL = etcdRecord.TTL
		recordsToCache = append(recordsToCache, cacheRecord)
	}
	return recordsToCache
}
