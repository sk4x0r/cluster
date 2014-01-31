package cluster

import(
	"fmt"
	//"os"
	"io/ioutil"
	"encoding/json"
	zmq "github.com/pebbe/zmq4"
	"strconv"
)


func speak(s string){
	fmt.Println(s)
}

func parseConfigFile(configFile string) (Config){
	//speak("Inside parseConfigFile()")
	content, err := ioutil.ReadFile(configFile)
	if err!=nil{
		fmt.Println("Error parsing the config file")
		panic(err)
	}
	var conf Config
	err=json.Unmarshal(content, &conf)
	if err!=nil{
		fmt.Println("Error parsing the config file")
		panic(err)
	}
	return conf
}

func loadServer(serverId int, conf Config)(Server){
	//speak("Inside loadServer()")
	var s Server
	s.pid=serverId
	
	s.peers=conf.getPeers(serverId)
	s.inbox=make(chan *Envelope,100)
	s.outbox=make(chan *Envelope,100)
	s.port=conf.getPort(serverId)
	s.peerInfo=conf.getPeerInfo(serverId)
	
	//create sockets for each peer
	s.connections=make(map[int] *zmq.Socket)
	for i:=range s.peers{
		peerId:=s.peers[i]
		
		sock, err :=zmq.NewSocket(zmq.REQ)
		if err!=nil{
			fmt.Println("Error creating socket",err)
		}
		st:="tcp://"+s.peerInfo[peerId].Ip+":"+strconv.Itoa(s.peerInfo[peerId].Port)
		sock.Connect(st)
		s.connections[peerId]=sock
		//s.peerInfo[peerId].soc= s.connections[peerId]
	}
	//fmt.Println(s.peerInfo)
	//fmt.Println(s.connections)
	go s.handleInbox()
	go s.handleOutbox()
	return s
}
func New(serverId int, configFile string) Server{
	//speak("Inside New()")
	conf:= parseConfigFile(configFile)
	s:=loadServer(serverId, conf)
	return s
}
