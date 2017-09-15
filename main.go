package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"log"

	"github.com/innovate-technologies/ITDNS/cache"
	"github.com/innovate-technologies/ITDNS/lookup"
)

type pdnsMessage struct {
	Method     string `json:"method"`
	Parameters struct {
		Path       string `json:"path"`
		Timeout    string `json:"timeout"`
		Local      string `json:"local"`
		Qname      string `json:"qname"`
		Qtype      string `json:"qtype"`
		RealRemote string `json:"real-remote"`
		Remote     string `json:"remote"`
		ZoneID     int    `json:"zone-id"`
	} `json:"parameters"`
}

type pdnsResult struct {
	Result []cache.Record `json:"result"`
}

var lookUpClient lookup.Client

func main() {
	fmt.Println("ITDNS etcd backend 2.0")

	lookUpClient = lookup.New()

	ln, err := net.Listen("unix", "/var/run/itdns")
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	os.Chmod("/var/run/itdns", 0777)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
			// handle error
		}
		go handleConnection(conn.(*net.UnixConn))
	}

}

func handleConnection(conn *net.UnixConn) {
	for {
		content, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Print(err)
			conn.Close()
			return
		}

		message := pdnsMessage{}
		json.Unmarshal([]byte(content), &message)

		if message.Method == "initialize" {
			conn.Write([]byte("{\"result\": true}"))
		} else if message.Method == "lookup" {
			records := lookUpClient.LookUp(message.Parameters.Qtype, message.Parameters.Qname)
			if records == nil || len(records) <= 0 {
				fmt.Println("send nil", message.Parameters.Qname, message.Parameters.Qtype)

				conn.Write([]byte("{\"result\": []}"))
				return
			}
			result := pdnsResult{
				Result: records,
			}
			resultString, _ := json.Marshal(result)
			conn.Write(resultString)
		} else { // all unimplemented
			conn.Write([]byte("{\"result\": false}"))
		}
	}
}
