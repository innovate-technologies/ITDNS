# ITDNS
ITDNS is a PowerDNS backend server for etcd stored records listening on a unix socket. ITDNS fetches all records from an etcd2/3
service and caches them into memory for fast serving. Thanks to etcd watchers the program can keep the records up to date

## Configuration
*To be added in a future release*

## How to set up
*soon to be added*

## Etcdv2 data structure
`/DNS/${domain}/${type}` with a JSON value of `[{value:"127.0.0.1", ttl:100}]`  
Note that the TTL is an integer and the value is always a string, for an MX record the `${priority} ${hostname}` format is used.  
The JSON array allows to set multiple records of a same type.

## Etcdv3 data structure
Probably the same? to be decided
