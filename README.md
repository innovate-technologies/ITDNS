# ITDNS
ITDNS is a PowerDNS backend server for etcd stored records listening on a unix socket. ITDNS fetches all records from an etcd2/3
service and caches them into memory for fast serving. Thanks to etcd watchers the program can keep the records up to date

## Configuration
Configuration via ITDNS is only done by envoirement variables.  
To enable etcd v2 set `ITDNS_ETCD2_ENDPOINTS` to a comma seperated list of endpoints for etcd v2.  
Configuration of etcd v3 allows more options:
* `ITDNS_ETCD3_ENDPOINTS` works the same as on v2
* `ITDNS_ETCD3_USERNAME` and `ITDNS_ETCD3_USERNAME` enable basic auth on etcd
* `ITDNS_ETCD3_CA` is the path to the CA certificate used by the etcd endpoints

## How to set up
*soon to be added*

## Etcd data structure
`/DNS/${domain}/${type}` with a JSON value of `[{"value":"127.0.0.1", "ttl":100}]`  
Note that the TTL is an integer and the value is always a string, for an MX record the value is an array where teh first element is a number for the priority and the second a string.  
The JSON array allows to set multiple records of a same type.

## ITDNSCTL
You can also use `itdnsctl` to add and remove entries to ITDNS. This utility will use all env vars used by etcdctl v3 to connect to etcd.