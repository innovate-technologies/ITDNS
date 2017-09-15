package lookup

import (
	"strings"

	"github.com/innovate-technologies/ITDNS/cache"
	"github.com/innovate-technologies/ITDNS/etcd2"
)

// Client lets you look up and cache DNS records
type Client struct {
	cache cache.Cache
	etcd2 etcd2.Etcd2Client
}

// New gived a new Lookup Client
func New() Client {
	returnObject := Client{}
	returnObject.cache = cache.New()
	returnObject.etcd2 = etcd2.New(&returnObject.cache)

	// initialize
	returnObject.createCache()
	go returnObject.watch()

	return returnObject
}

// createCache loads in all records into the cache for first usage
func (c *Client) createCache() {
	c.etcd2.CreateCache()
	// add etcd3 here to make sure it overrides
}

// watch allows to run etcd watchers to update the cache
func (c *Client) watch() {
	go c.etcd2.Watch()
}

// LookUp gives back all known records of a certain type for a domain
func (c *Client) LookUp(qtype, qname string) []cache.Record {
	qname = strings.TrimSuffix(qname, ".")
	if qtype == "SOA" {
		return c.sendSOA(qname)
	}
	records := c.cache.Get(strings.ToLower(qname))
	if len(records) <= 0 {
		records = c.lookUpLegacyInternal(strings.ToLower(qname))
	}
	if records == nil || len(records) <= 0 {
		return c.sendSOA(qname)
	}
	results := []cache.Record{}
	for _, record := range records {
		if record.Qtype == qtype || qtype == "ANY" {
			record.Qname = qname // back to WeIrDCaSE
			results = append(results, record)
		}
	}
	return results
}

// lookUpLegacyInternal converts the old name-int.domain.com to name.domain.com.int.domain.com
func (c *Client) lookUpLegacyInternal(qname string) []cache.Record {
	qnameParts := strings.Split(qname, ".")
	if len(qnameParts) > 5 && qnameParts[3] == "int" {
		records := c.cache.Get(qnameParts[0] + "-int." + strings.Join(qnameParts[len(qnameParts)-2:], "."))
		if len(records) > 0 {
			for i := range records {
				records[i].Qname = qname
			}
			return records
		}
	}
	return nil
}

func (c *Client) sendSOA(qname string) []cache.Record {
	record := cache.Record{
		Qname:    qname,
		Qtype:    "SOA",
		TTL:      10,
		Content:  "dns-par.shoutca.st. maartje.eyskens.me. 2016050400 7200 1800 1209600 7200", // to do: change this
		DomainID: -1,
	}
	return []cache.Record{record}
}
