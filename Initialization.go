package cluster

import (
	"fmt"
	//"os"
	"encoding/json"
	//zmq "github.com/pebbe/zmq4"
	"io/ioutil"
	//"strconv"
)
func parseConfigFile(configFile string) Config {
	//fmt.Println("Inside parseConfigFile()")
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println("Error parsing the config file")
		panic(err)
	}
	var conf Config
	err = json.Unmarshal(content, &conf)
	if err != nil {
		fmt.Println("Error parsing the config file")
		panic(err)
	}
	return conf
}

func loadServer(serverId int, conf Config) Server {
	//fmt.Println("Inside loadServer()")
	var s Server
	s.pid = serverId

	s.peers = conf.getPeers(serverId)
	s.inbox = make(chan *Envelope, 100)
	s.outbox = make(chan *Envelope, 100)
	s.port = conf.getPort(serverId)
	s.peerInfo = conf.getPeerInfo(serverId)
	//fmt.Println(s.peerInfo)
	//fmt.Println(s.connections)
	s.stopInbox = make(chan bool, 1)
	s.stopOutbox = make(chan bool, 1)
	s.initializeSockets()
	go s.handleInbox()
	go s.handleOutbox()
	return s
}
func New(serverId int, configFile string) Server {
	//fmt.Println("Inside New()")
	conf := parseConfigFile(configFile)
	s := loadServer(serverId, conf)
	return s
}

func NewWithConfig(serverId int, config Config) Server {
	//fmt.Println("Inside NewWithConfig()")
	s := loadServer(serverId, config)
	return s
}
