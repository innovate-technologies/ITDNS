package etcd2

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	etcdClient "github.com/coreos/etcd/client"
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

type Etcd2Client struct {
	etcdAPI etcdClient.KeysAPI
	cache   *cache.Cache
}

var cfg = config.GetConfig()

func New(cache *cache.Cache) Etcd2Client {
	returnObject := Etcd2Client{}
	returnObject.cache = cache

	c, err := etcdClient.New(etcdClient.Config{
		Endpoints:               cfg.Etcd2Config.Endpoints,
		Transport:               etcdClient.DefaultTransport,
		HeaderTimeoutPerRequest: 10 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	returnObject.etcdAPI = etcdClient.NewKeysAPI(c)

	connecting := true
	for connecting == true {
		// Testing connection
		_, err := returnObject.etcdAPI.Get(context.Background(), "/", &etcdClient.GetOptions{})
		if err != nil {
			fmt.Println("Waiting on Etcd")
			time.Sleep(1 * time.Second)
		} else {
			connecting = false
		}
	}

	return returnObject
}

// CreateCache fetches all records and enters them into the cache
func (l *Etcd2Client) CreateCache() {
	res, err := l.etcdAPI.Get(context.Background(), "/DNS/", &etcdClient.GetOptions{Recursive: true})
	if err != nil {
		panic(err)
	}
	for _, domain := range res.Node.Nodes {
		if domain.Nodes != nil {
			domainParts := strings.Split(domain.Key, "/")
			domainName := domainParts[len(domainParts)-1]
			recordsToCache := []cache.Record{}
			for _, record := range domain.Nodes {
				recordsToCache = addRecords(record, recordsToCache)
			}
			l.cache.Add(domainName, recordsToCache)
		}
	}
}

// Watch watches for new records to add them to the cache
func (l *Etcd2Client) Watch() {
	w := l.etcdAPI.Watcher("/DNS/", &etcdClient.WatcherOptions{Recursive: true})
	for {
		r, err := w.Next(context.Background())
		if err != nil {
			go l.Watch()
			return
		}
		if pathParts := strings.Split(r.Node.Key, "/"); r.Action == "set" && len(pathParts) >= 3 {
			l.addToCache(pathParts, 0)
		}
	}
}

func (l *Etcd2Client) addToCache(pathParts []string, retry int) {
	recordsToCache := []cache.Record{}
	res, err := l.etcdAPI.Get(context.Background(), strings.Join(pathParts[:len(pathParts)-1], "/"), &etcdClient.GetOptions{})
	if err != nil || res.Node.Nodes == nil {
		time.Sleep(time.Second)
		retry++
		if retry > 1000 {
			return
		}
		l.addToCache(pathParts, retry)
		return
	}
	for _, record := range res.Node.Nodes {
		recordsToCache = addRecords(record, recordsToCache)
	}
	l.cache.Add(pathParts[2], recordsToCache)
}

func addRecords(record *etcdClient.Node, recordsToCache []cache.Record) []cache.Record {
	pathParts := strings.Split(record.Key, "/")

	var etcdRecords []etcdRecord

	if pathParts[len(pathParts)-1] == "MX" {
		etcdRecords = []etcdRecord{}
		mxRecords := []etcdMxRecord{}
		json.Unmarshal([]byte(record.Value), &mxRecords)
		for _, record := range mxRecords {
			etcdRecords = append(etcdRecords, etcdRecord{
				TTL:   record.TTL,
				Value: fmt.Sprintf("%.0f", record.Value[0].(float64)) + " " + record.Value[1].(string),
			})
		}
	} else {
		etcdRecords = []etcdRecord{}
		json.Unmarshal([]byte(record.Value), &etcdRecords)
	}

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
